package handlers

import (
	"encoding/json"
	"net/http"
	"weather-app/models"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("query")

	weatherData, err := models.GetWeatherData(query)
	if err != nil {
		apiErr, ok := err.(*models.ApiError)
		if ok {
			errorResponse := struct {
				ErrorCode    int    `json:"error_code"`
				ErrorType    string `json:"error_type"`
				ErrorMessage string `json:"error_message"`
			}{
				ErrorCode:    apiErr.Code,
				ErrorType:    apiErr.Type,
				ErrorMessage: apiErr.Message,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				// Handle encoding error (unlikely)
				http.Error(w, "Error encoding response", http.StatusInternalServerError)
			}
		} else {
			// Handle other errors (not related to the API)
			http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		}
		return
	}

	response := struct {
		Temperature float64 `json:"temperature"`
		Condition   string  `json:"condition"`
	}{
		Temperature: weatherData.Temperature,
		Condition:   weatherData.Condition,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
