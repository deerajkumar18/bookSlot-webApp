package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	util "BookSlotApp/Util"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to eve'N'ts slots booking page")
}

func InitiateBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.Atoi(vars["eventID"])
	userID, _ := strconv.Atoi(vars["userID"])

	reqProcessingChan := EventsAndChannelMap[eventID]
	obj := &util.UserBookingReqPayload{
		EventsDao: &util.EventsInfoTbl{
			EventID: eventID,
		},
		SlotsDao: &util.SlotsInfoTbl{
			UserID:  userID,
			SlotID:  util.FormatUUIDToSlotID(uuid.New().String(), eventID, userID),
			EventID: eventID,
		},
	}
	reqProcessingChan.UserAction <- obj

	if err := <-reqProcessingChan.Res; err != nil {
		fmt.Fprintf(w, "Unable to finish the Booking process")
	}
}

func CancelBooking(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eventID, _ := strconv.Atoi(vars["eventID"])
	userID, _ := strconv.Atoi(vars["userID"])
	slotID := vars["slotID"]

	reqProcessingChan := EventsAndChannelMap[eventID]
	obj := &util.UserCancelSlotPayload{
		EventsDao: &util.EventsInfoTbl{
			EventID: eventID,
		},
		SlotsDao: &util.SlotsInfoTbl{
			UserID:  userID,
			SlotID:  slotID,
			EventID: eventID,
		},
	}

	reqProcessingChan.UserAction <- obj

	if err := <-reqProcessingChan.Res; err != nil {
		fmt.Fprintf(w, "Unable to finish the cancellation process")
	}

}
