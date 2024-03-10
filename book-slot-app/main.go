package main

import (
	util "BookSlotApp/Util"
	"BookSlotApp/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//get the list of events
	var Events []handlers.Event

	for _, event := range Events {
		var channelRef handlers.ChannelsByEvents
		eventChan := make(chan util.Booking)
		errChan := make(chan error)
		channelRef.UserAction = eventChan
		channelRef.Res = errChan
		handlers.EventsAndChannelMap[event.EventID] = channelRef
		go util.UserAction(eventChan, errChan)
	}
	r := mux.NewRouter()
	r.HandleFunc("/home", handlers.HomePage).Methods("GET")
	r.HandleFunc("/book/{userID}", handlers.InitiateBooking).Methods("POST")
	r.HandleFunc("/cancel/{userID}/{slotID}", handlers.CancelBooking).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", r))
}
