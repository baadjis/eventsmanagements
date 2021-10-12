package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ticket struct {
	gorm.Model
	Owner      string `json:"owner"`
	Number     string `json:"number"`
	EventID    string `json:"event"`
	EventTitle string `json:"eventtitle"`
}

func createTicket(w http.ResponseWriter, r *http.Request) {
	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	var newTicket ticket

	//get request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, " enter data to create a ticket")
	}

	json.Unmarshal(reqBody, &newTicket)

	db.Create(&newTicket)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTicket)
}
func getAllTickets(w http.ResponseWriter, r *http.Request) {

	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	// get query parameters
	var eventID string = strings.ToLower(r.URL.Query().Get("eventid"))

	var tickets []ticket

	db.Find(&tickets)

	// filter by eventID
	if len(eventID) > 0 {
		for i, ticket := range tickets {
			if strings.EqualFold(ticket.EventID, eventID) {
				tickets = append(tickets[:i], tickets[i+1:]...)

			}

		}
	}

	json.NewEncoder(w).Encode(tickets)
}

func deleteTicket(w http.ResponseWriter, r *http.Request) {

	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	TicketID := mux.Vars(r)["id"]
	var Aticket ticket

	db.Where("ID=?", TicketID).Find(&Aticket)
	db.Delete(&Aticket)
}

func getOneTicket(w http.ResponseWriter, r *http.Request) {

	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	TicketID := mux.Vars(r)["id"]
	var Aticket event

	db.Where("ID=?", TicketID).Find(&Aticket)

	json.NewEncoder(w).Encode(Aticket)

}

func modifyTicket(w http.ResponseWriter, r *http.Request) {

	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	TicketID := mux.Vars(r)["id"]
	var updatedTicket ticket
	var Aticket ticket

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedTicket)
	db.Where("ID=?", TicketID).Find(&Aticket)

	db.Model(&Aticket).Updates(updatedTicket)

	json.NewEncoder(w).Encode(Aticket)

}
