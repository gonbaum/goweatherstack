// main.go
// This is the entry point of your application. It will set up the HTTP server, define routes, and start the server. You can keep your main logic here or delegate it to separate files.
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

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	// This is where you would implement your weather handling logic
	// For now, let's just send a simple response
	fmt.Fprintln(w, "Weather API: Sunny and warm!")
}
