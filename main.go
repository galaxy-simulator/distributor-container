package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"

	"git.darknebu.la/GalaxySimulator/db-actions"
)

var (
	db                *sql.DB
	currentStarBuffer int64 = 1
	idBufferChannel         = make(chan int64, 100000000)

	// local http api
	port string

	// define the parameters needed to connect to the database
	DBHOST    = "postgres.docker.localhost"
	DBPORT    = 5432
	DBUSER    = "postgres"
	DBPASSWD  = ""
	DBNAME    = "postgres"
	DBSSLMODE = "disable"
)

func listOfStarIDs(treeindex int64) []int64 {
	log.Printf("getting a list of stars using the treeindex %d", treeindex)

	// get a list of all stars inside of the treeindex
	listofstars := db_actions.GetListOfStarIDsTimestep(db, treeindex)

	return listofstars
}

func fillStarIdBufferChannel() {
	log.Println("Filling the idBufferChannel")

	// get a list of ids using the currentStarBuffer value
	// the currentStarBuffer value is a counter keeping track of which galaxy is going to
	// be inserted into the idBufferChannel next
	listOfStarIDs := listOfStarIDs(currentStarBuffer)
	log.Printf("len(listOfStarIDs: %d)", len(listOfStarIDs))

	// insert all the ids from the list of ids into the idBufferChannel
	for _, id := range listOfStarIDs {
		idBufferChannel <- id
	}

	// increase the currentStarBuffer counter
	currentStarBuffer += 1
}

func getFlags() {
	// get the port on which the service should be hosted and the url of the database
	flag.StringVar(&port, "port", "8080", "port used to host the service")
	flag.Parse()
	log.Println("[ ] Done loading the flags")
}

func pingDB() {
	// ping the db
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("[ ] Done Pinging the DB")
}

func getEnvironmentVariables() {
	// get the data that should be used to connect to the database
	DBHOST = os.Getenv("DBURL")
	if DBHOST == "" {
		DBHOST = "postgresql.docker.localhost"
	}

	DBUSER = os.Getenv("DBUSER")
	if DBUSER == "" {
		DBUSER = "postgres"
	}

	DBPASSWD = os.Getenv("DBPASSWD")
	if DBPASSWD == "" {
		DBPASSWD = ""
	}

	DBPORT, _ := strconv.ParseInt(os.Getenv("DBPORT"), 10, 64)
	if DBPORT == 0 {
		DBPORT = 5432
	}

	DBNAME = os.Getenv("DBNAME")
	if DBNAME == "" {
		DBNAME = "postgres"
	}

	log.Printf("DBURL: %s", DBHOST)
	log.Printf("DBUSER: %s", DBUSER)
	log.Printf("DBPASSWD: %s", DBPASSWD)
	log.Printf("DBPORT: %d", DBPORT)
	log.Printf("DBPROJECTNAME: %s", DBNAME)
	log.Printf("frontend port: %s", port)
}

func main() {
	getFlags()
	getEnvironmentVariables()

	db = connectToDB()
	db.SetMaxOpenConns(75)
	pingDB()

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/distributor", distributorHandler).Methods("GET")

	log.Printf("[ ] Distributor on localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
