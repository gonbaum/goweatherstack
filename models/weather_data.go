package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type WeatherData struct {
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
}

type ApiResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		Temperature         float64  `json:"temperature"`
		WeatherDescriptions []string `json:"weather_descriptions"`
	} `json:"current"`
}

type ApiError struct {
	Code    int
	Type    string
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

func GetWeatherData(query string) (WeatherData, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("WEATHER_API_KEY")

	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY not found in .env file")
	}

	apiURL := fmt.Sprintf("http://api.weatherstack.com/current?access_key=%s&query=%s", apiKey, query)

	resp, err := http.Get(apiURL)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherData{}, errors.New("API request failed with status: " + resp.Status)
	}

	var apiResponse ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return WeatherData{}, err
	}

	if apiResponse.Error.Code != 0 {
		return WeatherData{}, &ApiError{Code: apiResponse.Error.Code, Type: apiResponse.Error.Type, Message: apiResponse.Error.Info}
	}

	weatherData := WeatherData{
		Temperature: apiResponse.Current.Temperature,
		Condition:   "",
	}

	if len(apiResponse.Current.WeatherDescriptions) > 0 {
		weatherData.Condition = apiResponse.Current.WeatherDescriptions[0]
	}

	return weatherData, nil
}
