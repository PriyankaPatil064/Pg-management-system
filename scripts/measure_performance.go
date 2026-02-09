package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func measureEndpoint(url string, iterations int) {
	fmt.Printf("Measuring endpoint: %s (%d iterations)\n", url, iterations)
	var totalDuration time.Duration
	var minDuration time.Duration
	var maxDuration time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("  Iteration %d: Error - %v\n", i+1, err)
			continue
		}
		resp.Body.Close()

		if i == 0 || duration < minDuration {
			minDuration = duration
		}
		if duration > maxDuration {
			maxDuration = duration
		}
		totalDuration += duration
		fmt.Printf("  Iteration %d: %v\n", i+1, duration)
	}

	avgDuration := totalDuration / time.Duration(iterations)
	fmt.Printf("\nResults for %s:\n", url)
	fmt.Printf("  Average: %v\n", avgDuration)
	fmt.Printf("  Minimum: %v\n", minDuration)
	fmt.Printf("  Maximum: %v\n", maxDuration)

	if avgDuration < 200*time.Millisecond {
		fmt.Println("  STATUS: PASS (< 200ms)")
	} else {
		fmt.Println("  STATUS: FAIL (> 200ms)")
	}
	fmt.Println("--------------------------------------------------")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	baseURL := "http://localhost:" + port

	endpoints := []string{
		"/health",
		"/rooms",
		"/guests",
	}

	fmt.Println("Starting Performance Measurement...")
	fmt.Println("==================================================")

	for _, endpoint := range endpoints {
		measureEndpoint(baseURL+endpoint, 5)
	}
}
