package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	apiURL         = "https://utilisocial.io/datacapable/v2/p/NES/map/events"
	defaultPort    = "8080"
	requestTimeout = 10 * time.Second
)

// OutageEvent represents an event from the NES outage API
type OutageEvent struct {
	ID              int     `json:"id"`
	StartTime       int64   `json:"startTime"`
	LastUpdatedTime int64   `json:"lastUpdatedTime"`
	Title           string  `json:"title"`
	NumPeople       int     `json:"numPeople"`
	Status          string  `json:"status"`
	Cause           string  `json:"cause"`
	Identifier      string  `json:"identifier"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}

// HealthResponse represents the JSON response from the health endpoint
type HealthResponse struct {
	Status     string   `json:"status"`
	Message    string   `json:"message,omitempty"`
	EventCount int      `json:"event_count,omitempty"`
	Checks     []Check  `json:"checks"`
	Timestamp  string   `json:"timestamp"`
}

// Check represents an individual health check result
type Check struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// validateStatusFields checks if the events have the required status fields
func validateStatusFields(events []OutageEvent) error {
	if len(events) == 0 {
		// Empty array is valid - it means no current outages
		return nil
	}

	// Check that at least one event has the required fields populated
	for _, event := range events {
		// Status field must exist (can be empty string but the field must be present)
		// Since we're using Go structs, if JSON parsing succeeded, the fields exist
		// We verify the data makes sense by checking if Status is a non-empty string
		if event.Status == "" {
			return fmt.Errorf("event %d has empty status field", event.ID)
		}
	}

	return nil
}

// checkAPIHealth performs the health check against the NES API
func checkAPIHealth() HealthResponse {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    []Check{},
	}

	// Check 1: Can we reach the API?
	apiCheck := Check{Name: "api_reachable", Status: "pass"}

	client := &http.Client{Timeout: requestTimeout}
	resp, err := client.Get(apiURL)
	if err != nil {
		apiCheck.Status = "fail"
		apiCheck.Error = fmt.Sprintf("failed to reach API: %v", err)
		response.Checks = append(response.Checks, apiCheck)
		response.Status = "unhealthy"
		response.Message = "Cannot reach NES API"
		return response
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiCheck.Status = "fail"
		apiCheck.Error = fmt.Sprintf("API returned status %d", resp.StatusCode)
		response.Checks = append(response.Checks, apiCheck)
		response.Status = "unhealthy"
		response.Message = fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode)
		return response
	}

	response.Checks = append(response.Checks, apiCheck)

	// Check 2: Can we parse the JSON response?
	parseCheck := Check{Name: "json_parseable", Status: "pass"}

	var events []OutageEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		parseCheck.Status = "fail"
		parseCheck.Error = fmt.Sprintf("failed to parse JSON: %v", err)
		response.Checks = append(response.Checks, parseCheck)
		response.Status = "unhealthy"
		response.Message = "Failed to parse API response as JSON"
		return response
	}

	response.Checks = append(response.Checks, parseCheck)
	response.EventCount = len(events)

	// Check 3: Do the events have the required status fields?
	statusFieldsCheck := Check{Name: "status_fields_present", Status: "pass"}

	if err := validateStatusFields(events); err != nil {
		statusFieldsCheck.Status = "fail"
		statusFieldsCheck.Error = err.Error()
		response.Checks = append(response.Checks, statusFieldsCheck)
		response.Status = "unhealthy"
		response.Message = "Status fields validation failed"
		return response
	}

	response.Checks = append(response.Checks, statusFieldsCheck)
	response.Message = fmt.Sprintf("All checks passed. Found %d outage events.", len(events))

	return response
}

// healthHandler handles requests to the /health endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	healthResult := checkAPIHealth()

	w.Header().Set("Content-Type", "application/json")

	if healthResult.Status == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(healthResult)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.HandleFunc("/health", healthHandler)

	log.Printf("Health check server starting on port %s", port)
	log.Printf("Health endpoint available at http://localhost:%s/health", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
