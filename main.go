package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/jdotcurs/go-weather-app/internal/weather"
)

func main() {
	service := weather.NewOpenMeteoService()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
		if err != nil {
			http.Error(w, "Invalid latitude", http.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
		if err != nil {
			http.Error(w, "Invalid longitude", http.StatusBadRequest)
			return
		}

		timezone := r.FormValue("timezone")
		if timezone == "" {
			timezone = "auto"
		}

		forecast, err := service.GetHourlyForecast(lat, lon, timezone)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching weather data: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(forecast)
	})

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
