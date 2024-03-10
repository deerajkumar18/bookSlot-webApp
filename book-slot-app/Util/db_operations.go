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

func GetSlotsAvailabilityByEventID(eventID int, ctx context.Context) (rows *sql.Row, err error) {
	db, err := ConnectDB()
	if err != nil {
		err = fmt.Errorf("unable to get the current status of slots availability , Err - %q", err)
		return
	}
	defer db.Close()

	rows = db.QueryRow("Select SlotsAvailable from EventsInfo where EventID=?", eventID)

	if ctx.Err() == context.DeadlineExceeded {
		return nil, errors.New("GetSlotsAvailabilityByEventID: Context Deadline Exceeded")
	}

	return

}

func UpdateSlotsAvailability(eventID int, updateCountVal int, ctx context.Context) (res sql.Result, err error) {
	db, err := ConnectDB()
	if err != nil {
		err = fmt.Errorf("unable to update the slots availability , Err - %q", err)
		return
	}
	defer db.Close()

	res, err = db.Exec("update EventsInfo set SlotsAvailable=? where EventID=?", updateCountVal, eventID)
	if ctx.Err() == context.DeadlineExceeded {
		return nil, errors.New("UpdateSlotsAvailability: Context Deadline Exceeded")
	}

	return

}

func InsertSlotsInfo(slotID string, eventID int, userID int, ctx context.Context) (res sql.Result, err error) {
	db, err := ConnectDB()
	if err != nil {
		return
	}
	defer db.Close()

	res, err = db.Exec("insert into SlotsInfo values(?,?,now(),?)", eventID, userID, slotID)
	if err != nil {
		err = fmt.Errorf("unable to insert slot information . Err - %v", err)
		return
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, errors.New("InsertSlotsInfo: Context Deadline Exceeded")
	}

	return
}

func DeleteSlotsInfo(slotID string, ctx context.Context) (res sql.Result, err error) {
	db, err := ConnectDB()
	if err != nil {
		return
	}
	defer db.Close()

	res, err = db.Exec("delete from SlotsInfo where SlotID=?", slotID)
	if err != nil {
		err = fmt.Errorf("unable to Delete slot information . Err - %v", err)
		return
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, errors.New("DeleteSlotsInfo: Context Deadline Exceeded")
	}

	return

}
