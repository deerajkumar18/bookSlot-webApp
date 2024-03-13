package concurrency_test

import (
	util "BookSlotApp/Util"
	"BookSlotApp/handlers"
	"context"
	"errors"
	"testing"
)

type testCase struct {
	testName string
	action   string
	eventID  []int
	userID   []int
	slotID   []string
}

type mockEventDao struct {
	EventID int
}

func (e *mockEventDao) GetSlotsAvailabilityByEventID(ctx context.Context) (SlotsAvailable int, err error) {
	if e.EventID == 100 {
		return SlotsAvailable, errors.New("no slots available")
	}
	return
}
func (e *mockEventDao) UpdateSlotsAvailability(requestType string, ctx context.Context) (err error) {
	return
}

type mockSlotsInfoDao struct {
	EventID int
	UserID  int
	SlotID  string
}

func (s *mockSlotsInfoDao) InsertSlotsInfo(ctx context.Context) (err error) {
	if s.UserID == 1001 {
		return errors.New("slot confirmation record entry failed")
	}
	return
}
func (s *mockSlotsInfoDao) DeleteSlotsInfo(ctx context.Context) (err error) {
	if s.SlotID == "aabb" {
		return errors.New("error occured while deleting record")
	}
	return
}

func TestBookAction(t *testing.T) {

	tests := []testCase{
		{
			//positive scenario - booking a slot
			testName: "booking a slot when there is an availability for the event",
			action:   "book",
			eventID:  []int{1},
			userID:   []int{1},
		},
		{
			//negative scenario - booking a slot
			testName: "booking a slot when there is no availability for the event",
			action:   "book",
			eventID:  []int{100},
			userID:   []int{1},
		},
		{
			//negative scenario - booking a slot
			testName: "booking a slot when there is an availability for the event , but internally failed due to DB error",
			action:   "book",
			eventID:  []int{1},
			userID:   []int{1001},
		},
		{
			//negative scenario - cancellation
			testName: "cancelling a slot when there is an availability for the event , but internally failed due to DB error",
			action:   "book",
			eventID:  []int{1},
			userID:   []int{1},
			slotID:   []string{"aabb"},
		},
		{
			//Concurrent -booking requests
			testName: "Two Booking requests on same event",
			action:   "book",
			eventID:  []int{1, 2},
			userID:   []int{1, 2},
		},
		{
			//Concurrent -booking requests
			testName: "Two Booking requests on different events",
			action:   "book",
			eventID:  []int{1, 2},
			userID:   []int{1, 2},
		},
		{
			//Concurrent -cancellation requests
			testName: "Two cancellation requests on same event",
			action:   "cancel",
			eventID:  []int{1, 2},
			userID:   []int{1, 2},
		},
		{
			//Concurrent -cancellation requests
			testName: "Two cancellation requests on different event",
			action:   "cancel",
			eventID:  []int{1, 2},
			userID:   []int{1, 2},
		},
	}

	for _, test := range tests {
		for i := 0; i < len(test.eventID); i++ {
			var channelRef handlers.ChannelsByEvents
			eventChan := make(chan util.Booking)
			errChan := make(chan error)
			channelRef.UserAction = eventChan
			channelRef.Res = errChan
			handlers.EventsAndChannelMap[test.eventID[i]] = channelRef
			go util.UserAction(eventChan, errChan)
		}

		for i, _ := range test.eventID {
			go func(i int) {
				obj1 := &util.UserBookingReqPayload{
					EventsDao: &mockEventDao{},
					SlotsDao:  &mockSlotsInfoDao{},
				}

				err := obj1.BookAction()
				if err != nil {
					t.Errorf("action %s failed due to err - %v", test.action, err)
				}
			}(i)
		}

	}

}
