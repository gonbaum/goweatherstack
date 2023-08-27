// models/weather_data.go
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
	// Replace with your actual access key and query parameters
	apiKey := os.Getenv("WEATHER_API_KEY")

	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY not found in .env file")
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("http://api.weatherstack.com/current?access_key=%s&query=%s", apiKey, query)

	// Make a request to the API
	resp, err := http.Get(apiURL)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	// Check the HTTP status code of the API response
	if resp.StatusCode != http.StatusOK {
		return WeatherData{}, errors.New("API request failed with status: " + resp.Status)
	}

	// Decode the API response
	// 1- var apiResponse ApiResponse: We declare a variable named apiResponse of type ApiResponse, which is the struct we've defined earlier to represent the structure of the API response.
	var apiResponse ApiResponse
	// .Decode(&apiResponse) attempts to decode the JSON data read from the response body and store it in the apiResponse variable. The & before apiResponse is used to pass a pointer to the variable so that the JSON decoder can populate it.
	// err will hold any error that might occur during the decoding process. If there's an error, we handle it in the following block of code.
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return WeatherData{}, err
	}

	// Check if the API response contains an error message
	if apiResponse.Error.Code != 0 {
		return WeatherData{}, &ApiError{Code: apiResponse.Error.Code, Type: apiResponse.Error.Type, Message: apiResponse.Error.Info}
	}

	// Create a WeatherData struct from the API response
	weatherData := WeatherData{
		Temperature: apiResponse.Current.Temperature,
		Condition:   "",
	}

	if len(apiResponse.Current.WeatherDescriptions) > 0 {
		weatherData.Condition = apiResponse.Current.WeatherDescriptions[0]
	}

	return weatherData, nil
}
