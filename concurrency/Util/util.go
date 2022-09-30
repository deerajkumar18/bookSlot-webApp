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

func FormatUUIDToSlotID(uuid string, eventID, userID string) string {
	uuid = strings.ReplaceAll(uuid, "-", "")
	uuid = uuid[:(len(uuid) / 2)]
	return uuid + eventID + userID
}

func Worker(cs chan UserBookingReqPayload, results chan EventSlotUser) {
	bookingInfo := <-cs
	slotsOccupancyDetails := EventSlotUser{bookingInfo.UserID, bookingInfo.EventID, bookingInfo.SlotID, bookingInfo.SlotCount}

	res, err := InsertSlotsInfo(slotsOccupancyDetails.SlotID, slotsOccupancyDetails.EventID, slotsOccupancyDetails.UserID)
	if err != nil {
		log.Println("Error in inserting the slotoccupancy details , ERR -", err)
	}
	if res != nil {
		if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
			log.Println("Insert failed for userid - ", slotsOccupancyDetails.UserID)
		}

		updateRes, err := UpdateSlotsAvailability(slotsOccupancyDetails)
		if err != nil {
			log.Println("Error in updating the sloutcount , ERR -", err)
		}

		if rowsAffected, _ := updateRes.RowsAffected(); rowsAffected != 1 {
			log.Println("update failed for userid - ", slotsOccupancyDetails.UserID)
		}
	} else {
		log.Println("User not found , ID - ", slotsOccupancyDetails.UserID)
	}

	results <- slotsOccupancyDetails

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
