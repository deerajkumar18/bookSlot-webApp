package main

import (
	"concurrency/operations"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/home", operations.HomePage).Methods("GET")
	r.HandleFunc("/book/{eventID}/{userID}", operations.InitiateBooking).Methods("POST")
	r.HandleFunc("/cancel/{eventID}/{userID}/{slotID}", operations.CancelBooking).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", r))
}
