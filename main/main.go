package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	baseURL = "https://roads.googleapis.com/v1/snapToRoads"
)

type Location struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

type SnappedPoint struct {
	Location Location `json:"location"`
}

type SnapToRoadsResponse struct {
	Points []SnappedPoint `json:"snappedPoints"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		fmt.Println("GOOGLE_MAPS_API_KEY not set in the environment")
		return
	}

	// Replace these coordinates with your desired points
	coordinates := "-35.27801,149.12958|-35.28032,149.12907|-35.28099,149.12929|-35.28144,149.12984|-35.28194,149.13003|-35.28282,149.12956|-35.28302,149.12881|-35.28473,149.12836"

	// Make the API request
	resp, err := http.Get(fmt.Sprintf("%s?path=%s&interpolate=true&key=%s", baseURL, coordinates, apiKey))
	if err != nil {
		fmt.Println("Error making API request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse the response JSON
	var snapToRoadsResp SnapToRoadsResponse
	err = json.Unmarshal(body, &snapToRoadsResp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Extract the snapped points and print them
	for _, point := range snapToRoadsResp.Points {
		fmt.Printf("Lat: %.6f, Lng: %.6f\n", point.Location.Lat, point.Location.Lng)
	}
}
