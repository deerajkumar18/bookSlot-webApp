package util

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

var m sync.Mutex

type UserBookingReqPayload struct {
	//sync.Mutex
	UserID    int
	EventID   int
	SlotID    string
	SlotCount int
}

type EventSlotUser struct {
	//sync.Mutex
	UserID             int
	EventID            int
	SlotID             string
	NoOfSlotsAvailable int
}

type UserCancelSlotPayload struct {
	//sync.Mutex
	UserID  int
	EventID int
	SlotID  string
}

type UpdatePayload struct {
	EventID        int
	AvailableSlots int
}

func FormatUUIDToSlotID(uuid string, eventID, userID string) string {
	uuid = strings.ReplaceAll(uuid, "-", "")
	uuid = uuid[:(len(uuid) / 2)]
	return uuid + eventID + userID
}

func BookSlotWorker(cs chan UserBookingReqPayload, results chan bool) {
	bookingInfo := <-cs
	slotsOccupancyDetails := EventSlotUser{bookingInfo.UserID, bookingInfo.EventID, bookingInfo.SlotID, bookingInfo.SlotCount}

	res, err := InsertSlotsInfo(slotsOccupancyDetails.SlotID, slotsOccupancyDetails.EventID, slotsOccupancyDetails.UserID)
	if err != nil {
		log.Println("Error in inserting the slotoccupancy details , ERR -", err)
		results <- false
	}
	if res != nil {
		if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
			log.Println("Insert failed for userid - ", slotsOccupancyDetails.UserID)
			results <- false
			return
		}

	} else {
		log.Println("User not found , ID - ", slotsOccupancyDetails.UserID)
		results <- false

	}

	results <- true

}

func CancelSlotWorker(cs chan UserCancelSlotPayload, result chan bool) {

	slotsInfo := <-cs

	res, err := DeleteSlotsInfo(slotsInfo.SlotID)
	if err != nil {
		log.Println(err)
		result <- false

	}

	log.Println(res)
	result <- true

}

func UpdateSlotInfoWorker(cs chan UpdatePayload) {

	EventSlotUpdateInfo := <-cs
	updateRes, err := UpdateSlotsAvailability(EventSlotUpdateInfo.EventID, EventSlotUpdateInfo.AvailableSlots)
	if err != nil {
		log.Println("Error in updating the sloutcount , ERR -", err)
		return
	}

	if rowsAffected, _ := updateRes.RowsAffected(); rowsAffected != 1 {
		log.Println("Event update failed for the event - ", EventSlotUpdateInfo.EventID)
	}
}

/*func MonitorResults(results <-chan EventSlotUser) {
	SlotsInfo := <-results
	fmt.Println("Monitor results", SlotsInfo)
	/*res, err := UpdateSlotsAvailability(SlotsInfo)
	if err != nil {
		fmt.Println("Error in updating the sloutcount , ERR - %v", err)
		return
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		fmt.Println("update for userid %d failed", SlotsInfo.UserID)
		return
	}

}*/

func SlotsClosed() {
	fmt.Println("Slots closed")
}
