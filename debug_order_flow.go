package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Testing order creation flow...")

	// Test data
	orderData := map[string]interface{}{
		"user_id": 123,
		"args": []map[string]interface{}{
			{
				"product_id": 1001,
				"quantity":   2,
				"price":      4.99,
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(orderData)
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	fmt.Printf("Sending request: %s\n", string(jsonData))

	// Create HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/order", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read response
	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)
	if err != nil && err != io.EOF {
		log.Fatal("Error reading response:", err)
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response body: %s\n", string(body[:n]))

	// Wait a bit for the event to be processed
	fmt.Println("Waiting 5 seconds for event processing...")
	time.Sleep(5 * time.Second)

	fmt.Println("Test completed. Check the logs for debug information.")
}
