package main

import (
	"concurrency/operations"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/Home", operations.HomePage).Methods("GET")
	r.HandleFunc("/Book/{eventID}/{userID}", operations.InitiateBooking).Methods("GET")
	r.HandleFunc("/Book/{eventID}/{userID}/{slotID}", operations.InitiateBooking).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", r))
}
