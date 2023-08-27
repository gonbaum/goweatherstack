package main

import (
	"fmt"
	"net/http"
	"weather-app/handlers"
)

func main() {
	port := 8080

	http.HandleFunc("/weather", handlers.WeatherHandler)

	fmt.Printf("Server started at :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
