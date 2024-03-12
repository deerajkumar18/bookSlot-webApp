package util

import (
	"BookSlotApp/dao"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type UserBookingReqPayload struct {
	//sync.Mutex
	EventsDao dao.EventsInfoDao
	SlotsDao  dao.SlotsInfoDao
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
	EventsDao dao.EventsInfoDao
	SlotsDao  dao.SlotsInfoDao
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

/*func NewUserBookingReqPayload(userID, eventID int) Booking {
	uuid := uuid.New().String()
	slotId := FormatUUIDToSlotID(uuid, eventID, userID)
	return &UserBookingReqPayload{UserID: userID, EventID: eventID, SlotID: slotId}
}*/

func (obj *UserBookingReqPayload) BookAction() (err error) {

	//slotsOccupancyDetails := EventSlotUser{bookingInfo.UserID, bookingInfo.EventID, bookingInfo.SlotID, bookingInfo.SlotCount}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	slots, err := obj.EventsDao.GetSlotsAvailabilityByEventID(ctxTimeout)
	if err != nil {
		return err
	}

	if slots >= 1 {
		err := obj.SlotsDao.InsertSlotsInfo(ctxTimeout)
		if err != nil {
			err = fmt.Errorf("error in inserting the slotoccupancy details , ERR - %q", err)
			return err
		}
	}

	err = obj.EventsDao.UpdateSlotsAvailability("book", ctxTimeout)
	if err != nil {
		return err
	}
	return nil
}

/*func NewUserBookingCancelReqPayload(eventID, userID int, slotID string) Booking {
	return &UserCancelSlotPayload{EventID: eventID, UserID: userID, SlotID: slotID}
}*/

func (obj *UserCancelSlotPayload) BookAction() (err error) {

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = obj.SlotsDao.DeleteSlotsInfo(ctxTimeout)
	if err != nil {
		return

	}

	err = obj.EventsDao.UpdateSlotsAvailability("cancel", ctxTimeout)
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
