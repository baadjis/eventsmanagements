package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type event struct {
	gorm.Model

	Title       string `json:"title"`
	Description string `json:"description"`
	City        string `json:"city"`
	Address     string `json:"address"`
	StartDate   string `json:"startdate"`
	EndDate     string `json:"enddate"`
	Creator     string `json:"creator"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createEvent(w http.ResponseWriter, r *http.Request) {

	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	var newEvent event

	//get request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)

	db.Create(&newEvent)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func getSingleEvent(w http.ResponseWriter, r *http.Request) {
	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	eventID := mux.Vars(r)["id"]
	var Anevent event
	db.Where("ID=?", eventID).Find(&Anevent)

	json.NewEncoder(w).Encode(Anevent)

}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	// get query parameters
	var title string = strings.ToLower(r.URL.Query().Get("title"))
	var city string = strings.Title(r.URL.Query().Get("city"))
	var description string = strings.ToLower(r.URL.Query().Get("description"))
	var start string = r.URL.Query().Get("startdate")

	//filter events  by city
	var events []event
	db.Where(&event{City: city}).Find(&events)

	// filter by title
	if len(title) > 0 {
		for i, ev := range events {
			if !(strings.Contains(strings.ToLower(ev.Title), title) || strings.Contains(title, strings.ToLower(ev.Title))) {
				events = append(events[:i], events[i+1:]...)

			}

		}
	}
	// filter by description content
	if len(description) > 0 {

		for i, ev := range events {
			if !(strings.Contains(strings.ToLower(ev.Description), description) || strings.Contains(description, strings.ToLower(ev.Description))) {
				events = append(events[:i], events[i+1:]...)

			}

		}
	}
	//filter by date
	if len(start) > 0 {
		startdate, _ := time.Parse("2006-01-02", start)
		for i, ev := range events {
			evdate, _ := time.Parse("2006-01-02", ev.StartDate)
			if evdate != startdate {
				events = append(events[:i], events[i+1:]...)

			}

		}

	}

	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {

	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	eventID := mux.Vars(r)["id"]

	var updatedEvent event
	var Anevent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)
	db.Where("ID=?", eventID).Find(&Anevent)

	db.Model(&Anevent).Updates(updatedEvent)

	json.NewEncoder(w).Encode(Anevent)

}

func deleteEvent(w http.ResponseWriter, r *http.Request) {

	//open databse
	var db *gorm.DB = openDataBase()
	defer db.Close()

	eventID := mux.Vars(r)["id"]

	var Anevent event

	db.Where("ID=?", eventID).Find(&Anevent)
	db.Delete(&Anevent)
}
