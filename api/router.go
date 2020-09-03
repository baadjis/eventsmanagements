package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func createEventsRouter() {
	router := mux.NewRouter().StrictSlash(true)

	//home
	router.HandleFunc("/", homeLink)
	//events

	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	//tikcets
	router.HandleFunc("/ticket", createTicket).Methods("POST")
	router.HandleFunc("/tickets/eventid", getAllTickets).Methods("GET")
	router.HandleFunc("/tickets/{id}", getOneTicket).Methods("GET")
	router.HandleFunc("/tickets/{id}", changeTicket).Methods("PATCH")
	router.HandleFunc("/tickets/{id}", deleteTicket).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
