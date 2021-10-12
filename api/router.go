package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func CreateEventsRouter() {
	router := mux.NewRouter().StrictSlash(true)

	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PATCH", "DELETE"})

	//home
	router.HandleFunc("/", homeLink)
	//events

	router.HandleFunc("/events", IsAuthorized(createEvent)).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getSingleEvent).Methods("GET")
	router.HandleFunc("/events/{id}", IsAuthorized(updateEvent)).Methods("PATCH")
	router.HandleFunc("/events/{id}", IsAuthorized(deleteEvent)).Methods("DELETE")

	//tikcets
	router.HandleFunc("/ticket", IsAuthorized(createTicket)).Methods("POST")
	router.HandleFunc("/tickets", IsAuthorized(getAllTickets)).Methods("GET")
	router.HandleFunc("/tickets/{id}", IsAuthorized(getOneTicket)).Methods("GET")
	router.HandleFunc("/tickets/{id}", IsAuthorized(modifyTicket)).Methods("PATCH")
	router.HandleFunc("/tickets/{id}", IsAuthorized(deleteTicket)).Methods("DELETE")

	// users

	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/logout", Logout).Methods("POST")
	router.HandleFunc("/users/{email}", IsAuthorized(GetUser)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
