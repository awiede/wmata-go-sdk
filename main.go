package main

import (
	"flag"
	"github.com/awiede/where-is-wmata/wmata"
	"log"
)

func main() {

	wmataKey := flag.String("wmata_key", "", "API key used to access WMATA API")
	metroLine := flag.String("metro_line", "", "The metro line to get information for")

	flag.Parse()

	if *wmataKey == "" {
		log.Fatalf("flag: wmata_key is required")
	}

	wmataService := wmata.Service{APIKey: *wmataKey}

	stations, err := wmataService.GetStationsByLine(*metroLine)

	if err != nil {
		log.Fatalf("unable to retrieve station information, got error: %s", err)
	}

	log.Println(stations)

}
