package internal

import "database/sql"

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	return db
}
