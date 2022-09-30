package operations

import (
	util "concurrency/Util"
	"fmt"
	"sync"
)

var mutex sync.Mutex

var EventSlotUserArr []util.EventSlotUser

var SlotID chan int

var ReqPayLoad chan util.UserBookingReqPayload

var Results chan util.EventSlotUser

func init() {

	workers := 10
	ReqPayLoad = make(chan util.UserBookingReqPayload, workers)
	Results = make(chan util.EventSlotUser, workers)
	//slotId := make(chan int, workers)
	//SlotID = make(chan int, workers)
	for i := 0; i < workers; i++ {
		go util.Worker(ReqPayLoad, Results)
	}
	fmt.Println("Workers created")
	//go util.MonitorResults(results)
	fmt.Println("Monitor reults created")
}
