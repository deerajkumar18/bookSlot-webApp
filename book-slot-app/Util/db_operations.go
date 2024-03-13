package util

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/bookslotdb")
	if err != nil {
		err = fmt.Errorf("unable to Connect to the Database. Err - %v ", err)
		return
	}
	return
}

/*func QueryWorker(queryChan chan DbQuery, m sync.Mutex) {
	m.Lock()
	defer m.Lock()
	query := <-queryChan
	query.Query()

}*/

type EventsInfoTbl struct {
	EventID        int
	EventName      string
	SlotsAvailable int
}

func (e *EventsInfoTbl) GetSlotsAvailabilityByEventID(ctx context.Context) (SlotsAvailable int, err error) {
	db, err := ConnectDB()
	if err != nil {
		err = fmt.Errorf("unable to get the current status of slots availability , Err - %q", err)
		return
	}
	defer db.Close()

	rows := db.QueryRow("Select SlotsAvailable from EventsInfo where EventID=?", e.EventID)

	if ctx.Err() == context.DeadlineExceeded {
		return SlotsAvailable, errors.New("GetSlotsAvailabilityByEventID: Context Deadline Exceeded")
	}

	rows.Scan(&SlotsAvailable)

	return

}

func (e *EventsInfoTbl) UpdateSlotsAvailability(requestType string, ctx context.Context) (err error) {
	db, err := ConnectDB()
	if err != nil {
		err = fmt.Errorf("unable to update the slots availability , Err - %q", err)
		return
	}
	defer db.Close()

	slots, err := e.GetSlotsAvailabilityByEventID(ctx)
	if err != nil {
		return err
	}

	if requestType == "book" {
		slots--
	}

	if requestType == "cancel" {
		slots++
	}

	res, err := db.Exec("update EventsInfo set SlotsAvailable=? where EventID=?", slots, e.EventID)
	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("UpdateSlotsAvailability: Context Deadline Exceeded")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		err = fmt.Errorf("Update slots availability failed for event id - %q", e.EventID)
		return err
	}

	return

}

type SlotsInfoTbl struct {
	EventID int
	UserID  int
	SlotID  string
}

func (s *SlotsInfoTbl) InsertSlotsInfo(ctx context.Context) (err error) {
	db, err := ConnectDB()
	if err != nil {
		return
	}
	defer db.Close()

	res, err := db.Exec("insert into SlotsInfo values(?,?,now(),?)", s.EventID, s.UserID, s.SlotID)
	if err != nil {
		err = fmt.Errorf("unable to insert slot information . Err - %v", err)
		return
	}

	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("InsertSlotsInfo: Context Deadline Exceeded")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		err = fmt.Errorf("insert failed for userid - %q", s.UserID)
		return
	}

	return
}

func (s *SlotsInfoTbl) DeleteSlotsInfo(ctx context.Context) (err error) {
	db, err := ConnectDB()
	if err != nil {
		return
	}
	defer db.Close()

	res, err := db.Exec("delete from SlotsInfo where SlotID=?", s.SlotID)
	if err != nil {
		err = fmt.Errorf("unable to Delete slot information . Err - %v", err)
		return
	}

	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("DeleteSlotsInfo: Context Deadline Exceeded")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		err = fmt.Errorf("Delete failed for slot id - %q", s.SlotID)
		return
	}

	return

}
