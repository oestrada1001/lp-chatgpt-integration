package database

import "database/sql"

func Connection() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/lp_db")
	if err != nil {
		panic(err)
	}
	return db
}
