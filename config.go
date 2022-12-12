package techtrain_api

import "database/sql"

func ConnectToDb() (*sql.DB, error) {
	db, errdb := sql.Open("mysql", "root:$A!c6i7xC!93%@FZ@localhost:3306/techtrain_db")
	if errdb != nil{
		return nil, errdb
	}
	err := db.Ping()
	return db, err
}
