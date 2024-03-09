package handlers

import (
	util "BookSlotApp/Util"
)

var (
	//mutex               sync.Mutex
	EventSlotUserArr []util.EventSlotUser
	SlotID           chan int
	EventID          chan int
	ReqPayLoad       chan util.UserBookingReqPayload
	BookingResults   chan bool
	//cancelPayLoad       chan util.UserCancelSlotPayload
	CancellationResults chan bool
	EventUpdatePayload  chan util.UpdatePayload
	//b                   chan util.Booking
	EventsAndChannelMap map[int]ChannelsByEvents
)

type Event struct {
	EventID           int
	Eventname         string
	SlotsAvailability int
}

type ChannelsByEvents struct {
	UserAction chan util.Booking
	Res        chan error
}

func init() {
	//get the list of events
	var Events []Event

	for _, event := range Events {
		var channelRef ChannelsByEvents
		eventChan := make(chan util.Booking)
		errChan := make(chan error)
		channelRef.UserAction = eventChan
		channelRef.Res = errChan
		EventsAndChannelMap[event.EventID] = channelRef
		go util.UserAction(eventChan, errChan)
	}
}
