package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

const baseURL = "http://127.0.0.1:8088/api/v1"

type userResponse struct {
	UserID string `json:"user_id"`
}

type userRequest struct {
	Balance int `json:"balance"`
}

type orderRequest struct {
	UserID       string   `json:"user_id"`
	Content      string   `json:"content"`
	Destinations []string `json:"destinations"`
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: loadtest <scenario1|scenario2> [args]")
	}

	switch os.Args[1] {
	case "scenario1":
		scenario1()
	case "scenario2":
		scenario2()
	default:
		log.Fatal("Invalid scenario. Use scenario1 or scenario2")
	}
}

// Scenario 1: Continuous flow test (create user -> create order)
func scenario1() {
	fs := flag.NewFlagSet("scenario1", flag.ExitOnError)
	balance := fs.Int("balance", 1000, "User balance")
	destCount := fs.Int("destinations", 10, "Number of destinations")
	rate := fs.Int("rate", 10, "Requests per second")
	duration := fs.String("duration", "30s", "Test duration")
	fs.Parse(os.Args[2:])

	testDuration, err := time.ParseDuration(*duration)
	if err != nil {
		log.Fatalf("Invalid duration: %v", err)
	}

	var wg sync.WaitGroup
	stop := time.After(testDuration)
	ticker := time.NewTicker(time.Second / time.Duration(*rate))
	defer ticker.Stop()

	var successCount, failureCount int64
	start := time.Now()

loop:
	for {
		select {
		case <-stop:
			break loop
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				userID, err := createUser(*balance)
				if err != nil {
					log.Printf("User creation failed: %v", err)
					failureCount++
					return
				}
				if err := createOrder(userID, *destCount); err != nil {
					log.Printf("Order creation failed: %v", err)
					failureCount++
					return
				}
				successCount++
			}()
		}
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("\nScenario 1 Results:")
	fmt.Printf("Success: %d, Failures: %d\n", successCount, failureCount)
	fmt.Printf("Elapsed: %s, RPS: %.2f\n",
		elapsed, float64(successCount+failureCount)/elapsed.Seconds())
}

// Scenario 2: Bulk creation test (create N users -> create orders for all)
func scenario2() {
	fs := flag.NewFlagSet("scenario2", flag.ExitOnError)
	userCount := fs.Int("users", 100, "Number of users")
	balance := fs.Int("balance", 1000, "User balance")
	destCount := fs.Int("destinations", 50, "Number of destinations")
	fs.Parse(os.Args[2:])

	// Create users
	userIDs := make([]string, 0, *userCount)
	for i := 0; i < *userCount; i++ {
		userID, err := createUser(*balance)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i+1, err)
			continue
		}
		userIDs = append(userIDs, userID)
	}

	fmt.Printf("Created %d users\n", len(userIDs))

	// Create orders concurrently
	var wg sync.WaitGroup
	start := time.Now()
	errCh := make(chan error, len(userIDs))

	for _, userID := range userIDs {
		wg.Add(1)
		go func(uid string) {
			defer wg.Done()
			if err := createOrder(uid, *destCount); err != nil {
				errCh <- fmt.Errorf("user %s: %w", uid, err)
			}
		}(userID)
	}

	wg.Wait()
	close(errCh)
	elapsed := time.Since(start)

	// Report errors
	var errorCount int
	for err := range errCh {
		log.Println(err)
		errorCount++
	}

	successCount := len(userIDs) - errorCount
	fmt.Printf("\nScenario 2 Results:\n")
	fmt.Printf("Users: %d, Orders: %d\n", len(userIDs), *destCount)
	fmt.Printf("Successful orders: %d\n", successCount)
	fmt.Printf("Failed orders: %d\n", errorCount)
	fmt.Printf("Total time: %s\n", elapsed)
	fmt.Printf("Throughput: %.2f orders/sec\n",
		float64(successCount)/elapsed.Seconds())
}

func createUser(balance int) (string, error) {
	url := baseURL + "/user"
	reqBody := userRequest{Balance: balance}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("status: %d", resp.StatusCode)
	}

	var userResp userResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return "", err
	}

	return userResp.UserID, nil
}

func createOrder(userID string, destCount int) error {
	url := baseURL + "/order"
	destinations := generateDestinations(destCount)

	reqBody := orderRequest{
		UserID:       userID,
		Content:      "hello arcs",
		Destinations: destinations,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}
	return nil
}

func generateDestinations(count int) []string {
	rand.Seed(time.Now().UnixNano())
	destinations := make([]string, count)
	used := make(map[string]bool)

	for i := 0; i < count; {
		num := rand.Intn(1000)
		dest := fmt.Sprintf("%11d", num)

		if !used[dest] {
			used[dest] = true
			destinations[i] = dest
			i++
		}
	}
	return destinations
}
