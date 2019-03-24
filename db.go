package main

import (
	"database/sql"
	"fmt"
	"log"
)

// connectToDB returns a pointer to an sql database writing to the database
func connectToDB() *sql.DB {
	//connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", DBUSER, DBNAME, DBSSLMODE)
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBHOST, DBPORT, DBUSER, DBPASSWD, DBNAME)
	db := dbConnect(connStr)
	return db
}

// dbConnect connects to a PostgreSQL database
func dbConnect(connStr string) *sql.DB {
	// connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("[ E ] connection: %v", err)
	}

	return db
}
