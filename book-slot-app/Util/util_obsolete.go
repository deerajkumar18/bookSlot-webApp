package util

import (
	"database/sql"
	"fmt"
)

type DbQuery interface {
	Query() error
}

type Select struct {
	EventID int
	Rows    *sql.Row
}

func NewSelectType(eventID int) DbQuery {
	return &Select{EventID: eventID}
}

func (s *Select) Query() (err error) {
	db, err := ConnectDB()
	if err != nil {
		err = fmt.Errorf("unable to get the current status of slots availability , Err - %q", err)
		return
	}
	defer db.Close()

	s.Rows = db.QueryRow("Select SlotsAvailable from EventsInfo where EventID=?", s.EventID)
	return
}

type Insert struct {
	EventID int
	UserID  int
	SlotID  string
	Res     sql.Result
}

func NewInsertType(eventID, userID int, slotID string) DbQuery {
	return &Insert{EventID: eventID, UserID: userID, SlotID: slotID}
}

func (i *Insert) Query() (err error) {
	db, err := ConnectDB()
	if err != nil {
		return
	}
	defer db.Close()

	i.Res, err = db.Exec("insert into SlotsInfo values(?,?,now(),?)", i.EventID, i.UserID, i.SlotID)
	if err != nil {
		err = fmt.Errorf("unable to insert slot information . Err - %v", err)
		return
	}

	return

}

type Update struct {
	EventID    int
	SlotsCount int
	Res        sql.Result
}

func NewUpdateType(eventID int, slotsCount int) DbQuery {
	return &Update{EventID: eventID, SlotsCount: slotsCount}
}

func (u *Update) Query() (err error) {
	db, err := ConnectDB()
	if err != nil {
		err = fmt.Errorf("unable to update the slots availability , Err - %q", err)
		return
	}
	defer db.Close()

	u.Res, err = db.Exec("update EventsInfo set SlotsAvailable=? where EventID=?", u.SlotsCount, u.EventID)

	return

}

type Delete struct {
	SlotID string
	Res    sql.Result
}

func NewDeleteType(slotID string) DbQuery {
	return &Delete{SlotID: slotID}
}

func (d *Delete) Query() (err error) {
	db, err := ConnectDB()
	if err != nil {
		return
	}
	defer db.Close()

	d.Res, err = db.Exec("delete from SlotsInfo where SlotID=?", d.SlotID)
	if err != nil {
		err = fmt.Errorf("unable to Delete slot information . Err - %v", err)
		return
	}

	return

}

/*func BookSlotWorker(cs chan UserBookingReqPayload, results chan bool) {
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

}*/

/*func CancelSlotWorker(cs chan UserCancelSlotPayload, result chan bool) {

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
}*/

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
