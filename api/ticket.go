package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ticket struct {
	gorm.Model
	Owner      string
	Number     string
	EventID    string
	EventTitle string
}

func createTicket(w http.ResponseWriter, r *http.Request) {
	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	var newTicket ticket

	//get request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
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

	eventID := mux.Vars(r)["eventid"]
	var tickets []ticket

	db.Where("EventID=?", eventID).Find(&tickets)
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

func changeTicket(w http.ResponseWriter, r *http.Request) {

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
