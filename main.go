package main

import (
	"encoding/json"
	"flag"
	"github.com/awiede/wmata-go-sdk/wmata"
	"log"
	"net/http"
	"time"
)

func main() {
	wmataKey := flag.String("wmata_key", "", "API key used to access WMATA API")

	flag.Parse()

	if *wmataKey == "" {
		log.Fatalf("flag: wmata_key is required")
	}

	wmataService := wmata.Service{
		APIKey: *wmataKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	http.HandleFunc("/StationList", logRequestMiddleware(getStationInfoHandler(&wmataService)))
	http.HandleFunc("/GetTrainPredictions", logRequestMiddleware(getTrainPredictionsHandler(&wmataService)))

	serverAddress := ":8080"

	log.Printf("Launching Server at %s", serverAddress)
	http.ListenAndServe(serverAddress, nil)

}

func logRequestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("[%s] [%s]", request.Method, request.URL)
		next.ServeHTTP(writer, request)
	}
}

func getStationInfoHandler(service *wmata.Service) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		line := request.URL.Query().Get("MetroLine")

		stations, err := service.StationList(line)

		if err != nil {
			log.Printf("error retrieving station information: %s", err)
			http.Error(writer, "error processing request", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(stations)

		if err != nil {
			log.Printf("error marshalling json: %s", err)
			http.Error(writer, "error processing request", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(response)

	}
}

func getTrainPredictionsHandler(service *wmata.Service) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		stationCode := request.URL.Query().Get("StationCode")

		trains, err := service.GetTrainPredictionsByStation(stationCode)

		if err != nil {
			log.Printf("error retrieving train predictions: %s", err)
			http.Error(writer, "error processing request", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(trains)

		if err != nil {
			log.Printf("error marshalling json: %s", err)
			http.Error(writer, "error processing request", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(response)

	}
}
