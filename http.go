package main

import (
	"fmt"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	responseString := `
<html>
	<head>
		<title>Distributor Container</title>
	</head>
	<body>
		<h1>Distributor</h1>
		<p>
			<a href="/distributor">Distributor</a>
		</p>
	</body>
</html>`
	_, _ = fmt.Fprintf(w, responseString)
}

func distributorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("The distributorHandler was accessed")

	// if the starIdBufferChannel is not filled yet, fill it
	if len(idBufferChannel) == 0 {
		log.Println("The idBufferChannel is empty, fetching new stars")
		fillStarIdBufferChannel()
	}

	// get a single id from the idBufferChannel
	log.Println("Getting an id from the idBufferChannel")
	id := <-idBufferChannel
	log.Println("Done...")

	// return the id using the http.ResponseWriter w
	_, _ = fmt.Fprintf(w, "%d", id)

	log.Printf("Done providing a starID (%d) from the StarBufferHandler", id)

}
