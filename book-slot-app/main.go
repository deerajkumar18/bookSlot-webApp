package main

import (
	"BookSlotApp/operations"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/home", operations.HomePage).Methods("GET")
	r.HandleFunc("/book/{userID}", operations.InitiateBooking).Methods("POST")
	r.HandleFunc("/cancel/{userID}/{slotID}", operations.CancelBooking).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", r))
}
