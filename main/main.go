package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joho/godotenv"
)

type Location struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

type Street []Location

type GeoJSONFeature struct {
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read the contents of index.html file
		indexHTML, err := ioutil.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Error reading index.html", http.StatusInternalServerError)
			return
		}

		// Serve the index.html file as the response
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHTML)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/get_streets", func(w http.ResponseWriter, r *http.Request) {
		// Load and parse your Bristol GeoJSON file
		bristolGeoJSON, err := ioutil.ReadFile("bristol.geojson")
		if err != nil {
			http.Error(w, "Error reading Bristol GeoJSON file", http.StatusInternalServerError)
			fmt.Println("Error reading Bristol GeoJSON file:", err)
			return
		}

		// Parse the GeoJSON data into a struct
		var geoJSON struct {
			Features []GeoJSONFeature `json:"features"`
		}
		if err := json.Unmarshal(bristolGeoJSON, &geoJSON); err != nil {
			http.Error(w, "Error parsing GeoJSON", http.StatusInternalServerError)
			fmt.Println("Error parsing GeoJSON:", err)
			return
		}

		// Extract streets and send them to the front end
		var streets []Street
		for _, feature := range geoJSON.Features {
			var street Street
			for _, coordsList := range feature.Geometry.Coordinates {
				for _, coords := range coordsList {
					location := Location{Lat: coords[1], Lng: coords[0]}
					street = append(street, location)
				}
			}
			streets = append(streets, street)
		}

		// Convert streets to JSON and send to front end
		streetsJSON, err := json.Marshal(streets)
		if err != nil {
			http.Error(w, "Error encoding streets to JSON", http.StatusInternalServerError)
			fmt.Println("Error encoding streets to JSON:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(streetsJSON)
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
