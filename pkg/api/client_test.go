package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testClient extends the Client with a custom baseURL for testing
type testClient struct {
	*Client
	testBaseURL string
}

// Override methods that use baseURL to use the test server URL instead
func (c *testClient) GetZones() ([]Zone, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/zones", c.testBaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Auth-API-Token", c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	var zonesResp ZonesResponse
	if err := json.NewDecoder(resp.Body).Decode(&zonesResp); err != nil {
		return nil, err
	}

	return zonesResp.Zones, nil
}

func (c *testClient) GetZoneIDByName(name string) (string, error) {
	zones, err := c.GetZones()
	if err != nil {
		return "", err
	}

	for _, zone := range zones {
		if zone.Name == name {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("zone with name %s not found", name)
}

func (c *testClient) CreateRecord(record Record) (*Record, error) {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/records", c.testBaseURL), bytes.NewBuffer(recordJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Auth-API-Token", c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	var recordResp RecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&recordResp); err != nil {
		return nil, err
	}

	return &recordResp.Record, nil
}

func (c *testClient) UpdateRecord(record Record) (*Record, error) {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/records/%s", c.testBaseURL, record.ID), bytes.NewBuffer(recordJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Auth-API-Token", c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	var recordResp RecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&recordResp); err != nil {
		return nil, err
	}

	return &recordResp.Record, nil
}

func (c *testClient) DeleteRecord(recordID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/records/%s", c.testBaseURL, recordID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Auth-API-Token", c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	return nil
}

// newTestClient creates a new test client with the given server
func newTestClient(server *httptest.Server) *testClient {
	client := NewClient("test-token")
	client.httpClient = server.Client()
	return &testClient{
		Client:      client,
		testBaseURL: server.URL,
	}
}

func setupTestServer(t *testing.T, path string, statusCode int, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		if r.URL.Path != path {
			t.Errorf("Expected path %s, got %s", path, r.URL.Path)
		}

		// Check auth header
		authHeader := r.Header.Get("Auth-API-Token")
		if authHeader != "test-token" {
			t.Errorf("Expected Auth-API-Token header 'test-token', got '%s'", authHeader)
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		// Write response body
		if response != nil {
			json.NewEncoder(w).Encode(response)
		}
	}))
}

func TestNewClient(t *testing.T) {
	client := NewClient("test-token")
	if client.apiToken != "test-token" {
		t.Errorf("Expected apiToken to be 'test-token', got '%s'", client.apiToken)
	}
	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}
}

func TestGetZones(t *testing.T) {
	// Setup test server
	zonesResponse := ZonesResponse{
		Zones: []Zone{
			{
				ID:           "zone1",
				Name:         "example.com",
				TTL:          3600,
				RecordsCount: 5,
			},
			{
				ID:           "zone2",
				Name:         "example.org",
				TTL:          7200,
				RecordsCount: 3,
			},
		},
	}

	server := setupTestServer(t, "/zones", http.StatusOK, zonesResponse)
	defer server.Close()

	// Create test client
	client := newTestClient(server)

	// Call the method
	zones, err := client.GetZones()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check results
	if len(zones) != 2 {
		t.Errorf("Expected 2 zones, got %d", len(zones))
	}
	if zones[0].ID != "zone1" || zones[0].Name != "example.com" {
		t.Errorf("Zone data doesn't match expected values")
	}
}

func TestGetZoneIDByName(t *testing.T) {
	// Setup test server
	zonesResponse := ZonesResponse{
		Zones: []Zone{
			{
				ID:           "zone1",
				Name:         "example.com",
				TTL:          3600,
				RecordsCount: 5,
			},
			{
				ID:           "zone2",
				Name:         "example.org",
				TTL:          7200,
				RecordsCount: 3,
			},
		},
	}

	server := setupTestServer(t, "/zones", http.StatusOK, zonesResponse)
	defer server.Close()

	// Create test client
	client := newTestClient(server)

	// Test finding an existing zone
	zoneID, err := client.GetZoneIDByName("example.com")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if zoneID != "zone1" {
		t.Errorf("Expected zone ID 'zone1', got '%s'", zoneID)
	}

	// Test finding a non-existent zone
	_, err = client.GetZoneIDByName("nonexistent.com")
	if err == nil {
		t.Error("Expected error for non-existent zone, got nil")
	}
}

func TestCreateRecord(t *testing.T) {
	// Setup test record
	record := Record{
		Type:   "A",
		Name:   "test",
		Value:  "192.168.1.1",
		ZoneID: "zone1",
	}

	// Setup expected response
	expectedRecord := record
	expectedRecord.ID = "record1"
	recordResponse := RecordResponse{
		Record: expectedRecord,
	}

	// Setup test server
	server := setupTestServer(t, "/records", http.StatusOK, recordResponse)
	defer server.Close()

	// Create test client
	client := newTestClient(server)

	// Call the method
	createdRecord, err := client.CreateRecord(record)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check results
	if createdRecord.ID != "record1" {
		t.Errorf("Expected record ID 'record1', got '%s'", createdRecord.ID)
	}
	if createdRecord.Type != "A" || createdRecord.Name != "test" || createdRecord.Value != "192.168.1.1" {
		t.Errorf("Record data doesn't match expected values")
	}
}

func TestUpdateRecord(t *testing.T) {
	// Setup test record
	record := Record{
		ID:     "record1",
		Type:   "A",
		Name:   "test",
		Value:  "192.168.1.2", // Updated value
		ZoneID: "zone1",
	}

	// Setup expected response
	recordResponse := RecordResponse{
		Record: record,
	}

	// Setup test server
	server := setupTestServer(t, "/records/record1", http.StatusOK, recordResponse)
	defer server.Close()

	// Create test client
	client := newTestClient(server)

	// Call the method
	updatedRecord, err := client.UpdateRecord(record)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check results
	if updatedRecord.ID != "record1" {
		t.Errorf("Expected record ID 'record1', got '%s'", updatedRecord.ID)
	}
	if updatedRecord.Value != "192.168.1.2" {
		t.Errorf("Expected updated value '192.168.1.2', got '%s'", updatedRecord.Value)
	}
}

func TestDeleteRecord(t *testing.T) {
	// Setup test server
	server := setupTestServer(t, "/records/record1", http.StatusOK, nil)
	defer server.Close()

	// Create test client
	client := newTestClient(server)

	// Call the method
	err := client.DeleteRecord("record1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
