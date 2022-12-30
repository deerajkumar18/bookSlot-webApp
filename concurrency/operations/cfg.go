package operations

import (
	util "concurrency/Util"
	"fmt"
	"sync"
)

var (
	mutex               sync.Mutex
	EventSlotUserArr    []util.EventSlotUser
	SlotID              chan int
	EventID             chan int
	ReqPayLoad          chan util.UserBookingReqPayload
	BookingResults      chan bool
	cancelPayLoad       chan util.UserCancelSlotPayload
	CancellationResults chan bool
	EventUpdatePayload  chan util.UpdatePayload
)

func init() {

	workers := 2
	ReqPayLoad = make(chan util.UserBookingReqPayload, workers)
	BookingResults = make(chan bool, 1)
	cancelPayLoad = make(chan util.UserCancelSlotPayload, workers)
	CancellationResults = make(chan bool, 1)
	EventUpdatePayload = make(chan util.UpdatePayload, workers)
	//slotId := make(chan int, workers)
	//SlotID = make(chan int, workers)
	for i := 0; i < workers; i++ {
		go util.BookSlotWorker(ReqPayLoad, BookingResults)
		go util.CancelSlotWorker(cancelPayLoad, CancellationResults)
		go util.UpdateSlotInfoWorker(EventUpdatePayload)
	}
	fmt.Println("Workers created")
	//go util.MonitorResults(results)
	//fmt.Println("Monitor reults created")
}
