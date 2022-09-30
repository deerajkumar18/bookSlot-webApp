package operations

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	util "concurrency/Util"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to VaxciNation slots booking page")
	//fmt.Fprintf(w, "Number of Slots available currently - %d", slots)

}

func InitiateBooking(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	vars := mux.Vars(r)
	eventID, _ := strconv.Atoi(vars["eventID"])
	userID, _ := strconv.Atoi(vars["userID"])
	var userReq util.UserBookingReqPayload
	//util.EventSlotUser{UserID: userID, EventID: eventID, SlotID: nil, NoOfSlotsAvailable: nil}
	userReq.EventID = eventID
	userReq.UserID = userID
	fmt.Println("Req from User ID - ", userID, "for Event ID - ", eventID)
	rows, err := util.GetSlotsAvailabilityByEventID(eventID)
	if err != nil {
		log.Println(err)
		return
	}
	var slots int
	rows.Scan(&slots)
	fmt.Println("Total no of slots available - ", slots)

	//slotID := make(chan int)

	if slots >= 1 {
		fmt.Println("Starting booking process ")
		uuid := uuid.New().String()
		slotId := util.FormatUUIDToSlotID(uuid, vars["eventID"], vars["userID"])
		fmt.Println(slotId)
		userReq.SlotCount = slots
		userReq.SlotID = slotId
		fmt.Println(userReq)
		//SlotID <- slotId
		ReqPayLoad <- userReq

		slotBookingResult := <-Results

		//close(slotID)
		fmt.Println("Done for the User - ", userID, slotBookingResult)
	} else {
		util.SlotsClosed()
	}

}

/*func MonitorBookings(cs chan string) {
	//wg.Wait()
	close(cs)
}

func PrintBookings(cs <-chan string, done chan<- bool) {
	for i := range cs {
		fmt.Println(i)
	}

	done <- true
}*/
