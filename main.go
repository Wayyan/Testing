package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	// ==========================================
	// CONFIGURATION (Defaults)
	// You can change these defaults here, or override them via terminal
	// ==========================================
	baseURL := flag.String("url", "http://localhost", "The base API URL")
	appID := flag.String("app-id", "-", "The App ID for the header")
	secretKey := flag.String("secret-key", "-", "The Secret Key for the header")
	totalRequestsPtr := flag.Int("total-requests", 35, "The total number of requests")

	// Parse the flags from the terminal
	flag.Parse()

	totalRequests := *totalRequestsPtr

	// ==========================================
	// SCRIPT LOGIC
	// ==========================================
	ids := make([]int, totalRequests)
	for i := 0; i < totalRequests; i++ {
		ids[i] = i + 1
	}

	// Shuffle the IDs
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	// Reusable HTTP client
	client := &http.Client{Timeout: 15 * time.Second}

	fmt.Printf("Starting script targeting: %s\n", *baseURL)
	fmt.Printf("Using App-ID: %s\n\n", *appID)
	fmt.Printf("Total requests: %d\n\n", totalRequests)

	for index, id := range ids {
		// Construct the URL using the variable
		targetURL := fmt.Sprintf("%s%d", *baseURL, id)

		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			fmt.Printf("Error creating request for ID %d: %v\n", id, err)
			continue
		}

		// Add headers using the variables
		req.Header.Set("app-id", *appID)
		req.Header.Set("secret-key", *secretKey)

		// Make the API Call
		fmt.Printf("[%d/%d] Fetching ID: %d... ", index+1, totalRequests, id)
		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("FAILED (%v)\n", err)
		} else {
			fmt.Printf("SUCCESS (Status: %s)\n", resp.Status)
			_, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
		}

		// Sleep for a random time between 1 and 10 seconds
		sleepSeconds := rand.Intn(6) + 1
		fmt.Printf("--> Waiting %d seconds before next call...\n\n", sleepSeconds)
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}

	fmt.Printf("All %d requests completed!", totalRequests)
}
