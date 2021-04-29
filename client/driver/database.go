package driver

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"gitlab.jiangxingai.com/poseidon/client/interfaces"
	log "k8s.io/klog"
)

var Db *sql.DB

//初始化sqlite
func Setup() {
	db, err := sql.Open("sqlite3", "./poseidon.db")
	if err != nil {
		panic(err)
	}
	Db = db

	initTable := `CREATE TABLE log (
    	uid INTEGER PRIMARY KEY AUTOINCREMENT,
      	timestamp INTEGER NULL,
     	eventType VARCHAR(64) NULL,
     	description VARCHAR(64)  NULL);`

	_, err = Db.Exec(initTable)
	if err != nil {
		log.Error(err)
	}
}

func QueryLog(starTime int64, endTime int64, eventType string, limit int, offset int) ([]interfaces.LogMsg, int, error) {
	var data interfaces.LogMsg
	var res []interfaces.LogMsg
	var rows *sql.Rows
	var count int
	var err, errS error
	if eventType != "" {
		rows, err = Db.Query("SELECT * FROM log where eventType = ? AND timestamp <= ? AND timestamp >= ? ORDER BY timestamp DESC limit ? offset ?",
			eventType, endTime, starTime, limit, offset)
		errS = Db.QueryRow("SELECT COUNT(*) FROM log where eventType = ? AND timestamp <= ? AND timestamp >= ? ",
			eventType, endTime, starTime).Scan(&count)
	} else {
		rows, err = Db.Query("SELECT * FROM log where timestamp <= ? AND timestamp >= ? ORDER BY timestamp DESC  limit ? offset ? ",
			endTime, starTime, limit, offset)
		errS = Db.QueryRow("SELECT COUNT(*) FROM log where timestamp <= ? AND timestamp >= ? ",
			endTime, starTime).Scan(&count)
	}
	if err != nil || errS != nil {
		return res, count, err
	}
	for rows.Next() {
		err = rows.Scan(&data.Uid, &data.Timestamp, &data.EventType, &data.Description)
		if err != nil {
			return res, count, err
		}
		res = append(res, data)
	}
	rows.Close()
	return res, count, nil
}

func InsertLog(data interfaces.LogMsg) error {
	stmt, err := Db.Prepare("INSERT INTO log(timestamp, eventType, description) values(?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data.Timestamp, data.EventType, data.Description)
	if err != nil {
		return err
	}
	return nil
}

func DeleteLog(id int) error {
	stmt, err := Db.Prepare("delete from log where uid=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateLog(description string, id int) error {
	stmt, err := Db.Prepare("update log set description=? where uid=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(description, id)
	if err != nil {
		return err
	}
	return nil
}