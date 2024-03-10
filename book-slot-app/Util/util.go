package util

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserBookingReqPayload struct {
	//sync.Mutex
	UserID  int
	EventID int
	SlotID  string
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

type Booking interface {
	BookAction() error
}

func FormatUUIDToSlotID(uuid string, eventID, userID int) string {
	eID := strconv.Itoa(eventID)
	uID := strconv.Itoa(userID)
	uuid = strings.ReplaceAll(uuid, "-", "")
	uuid = uuid[:(len(uuid) / 2)]
	return uuid + eID + uID
}

func NewUserBookingReqPayload(userID, eventID int) Booking {
	uuid := uuid.New().String()
	slotId := FormatUUIDToSlotID(uuid, eventID, userID)
	return &UserBookingReqPayload{UserID: userID, EventID: eventID, SlotID: slotId}
}

func (obj *UserBookingReqPayload) BookAction() (err error) {

	//slotsOccupancyDetails := EventSlotUser{bookingInfo.UserID, bookingInfo.EventID, bookingInfo.SlotID, bookingInfo.SlotCount}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := GetSlotsAvailabilityByEventID(obj.EventID, ctxTimeout)
	if err != nil {
		return err
	}
	var slots int
	rows.Scan(&slots)

	if slots >= 1 {

		res, err := InsertSlotsInfo(obj.SlotID, obj.EventID, obj.UserID, ctxTimeout)
		if err != nil {
			err = fmt.Errorf("error in inserting the slotoccupancy details , ERR - %q", err)
			return err
		}
		if res != nil {
			if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
				err = fmt.Errorf("insert failed for userid - %q", obj.UserID)
				return err
			}

		} else {
			err = fmt.Errorf("user not found , ID - %q", obj.UserID)
			return err
		}
	}

	_, err = UpdateSlotsAvailability(obj.EventID, slots-1, ctxTimeout)
	if err != nil {
		return err
	}
	return nil
}

func NewUserBookingCancelReqPayload(eventID, userID int, slotID string) Booking {
	return &UserCancelSlotPayload{EventID: eventID, UserID: userID, SlotID: slotID}
}

func (obj *UserCancelSlotPayload) BookAction() (err error) {

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = DeleteSlotsInfo(obj.SlotID, ctxTimeout)
	if err != nil {
		return

	}

	rows, err := GetSlotsAvailabilityByEventID(obj.EventID, ctxTimeout)
	if err != nil {
		return err
	}
	var slots int
	rows.Scan(&slots)

	_, err = UpdateSlotsAvailability(obj.EventID, slots+1, ctxTimeout)
	if err != nil {
		return err
	}

	return

}

func UserAction(b chan Booking, res chan error) {
	for obj := range b {
		err := obj.BookAction()
		res <- err
	}

	res <- nil
}
