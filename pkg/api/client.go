package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL = "https://dns.hetzner.com/api/v1"
)

// Client represents a Hetzner DNS API client
type Client struct {
	apiToken   string
	httpClient *http.Client
}

// NewClient creates a new Hetzner DNS API client
func NewClient(apiToken string) *Client {
	return &Client{
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Zone represents a DNS zone
type Zone struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	TTL          int    `json:"ttl"`
	RecordsCount int    `json:"records_count"`
	// Other fields omitted for simplicity and to avoid unmarshaling issues
}

// Record represents a DNS record
type Record struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl,omitempty"`
	ZoneID   string `json:"zone_id"`
	Created  string `json:"created,omitempty"`
	Modified string `json:"modified,omitempty"`
}

// ZonesResponse represents the response from the zones endpoint
type ZonesResponse struct {
	Zones []Zone `json:"zones"`
	Meta  struct {
		Pagination struct {
			Page         int `json:"page"`
			PerPage      int `json:"per_page"`
			LastPage     int `json:"last_page"`
			TotalEntries int `json:"total_entries"`
		} `json:"pagination"`
	} `json:"meta"`
}

// RecordsResponse represents the response from the records endpoint
type RecordsResponse struct {
	Records []Record `json:"records"`
	Meta    struct {
		Pagination struct {
			Page         int `json:"page"`
			PerPage      int `json:"per_page"`
			LastPage     int `json:"last_page"`
			TotalEntries int `json:"total_entries"`
		} `json:"pagination"`
	} `json:"meta"`
}

// RecordResponse represents the response when creating or getting a single record
type RecordResponse struct {
	Record Record `json:"record"`
}

// ZoneResponse represents the response when getting a single zone
type ZoneResponse struct {
	Zone Zone `json:"zone"`
}

// GetZones retrieves all DNS zones
func (c *Client) GetZones() ([]Zone, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/zones", baseURL), nil)
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	// Use a map to avoid unmarshaling issues with unexpected fields
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract zones from the response
	zonesData, ok := result["zones"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format: zones field not found or not an array")
	}

	var zones []Zone
	for _, zoneData := range zonesData {
		zoneMap, ok := zoneData.(map[string]interface{})
		if !ok {
			continue
		}

		zone := Zone{}

		// Extract ID
		if id, ok := zoneMap["id"].(string); ok {
			zone.ID = id
		}

		// Extract Name
		if name, ok := zoneMap["name"].(string); ok {
			zone.Name = name
		}

		// Extract TTL
		if ttl, ok := zoneMap["ttl"].(float64); ok {
			zone.TTL = int(ttl)
		}

		// Extract RecordsCount
		if recordsCount, ok := zoneMap["records_count"].(float64); ok {
			zone.RecordsCount = int(recordsCount)
		}

		zones = append(zones, zone)
	}

	return zones, nil
}

// GetZoneIDByName retrieves a zone ID by its name
func (c *Client) GetZoneIDByName(name string) (string, error) {
	// Normalize the input name
	normalizedName := strings.ToLower(strings.TrimSuffix(name, "."))

	// Get all zones
	zones, err := c.GetZones()
	if err != nil {
		return "", err
	}

	// Look for exact match
	for _, zone := range zones {
		zoneName := strings.ToLower(strings.TrimSuffix(zone.Name, "."))
		if zoneName == normalizedName {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("zone with name '%s' not found", name)
}

// GetRecords retrieves all DNS records for a zone
func (c *Client) GetRecords(zoneID string) ([]Record, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/records?zone_id=%s", baseURL, zoneID), nil)
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	// Use a map to avoid unmarshaling issues
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract records from the response
	recordsData, ok := result["records"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format: records field not found or not an array")
	}

	var records []Record
	for _, recordData := range recordsData {
		recordMap, ok := recordData.(map[string]interface{})
		if !ok {
			continue
		}

		record := Record{}

		// Extract ID
		if id, ok := recordMap["id"].(string); ok {
			record.ID = id
		}

		// Extract Type
		if recordType, ok := recordMap["type"].(string); ok {
			record.Type = recordType
		}

		// Extract Name
		if name, ok := recordMap["name"].(string); ok {
			record.Name = name
		}

		// Extract Value
		if value, ok := recordMap["value"].(string); ok {
			record.Value = value
		}

		// Extract TTL
		if ttl, ok := recordMap["ttl"].(float64); ok {
			record.TTL = int(ttl)
		}

		// Extract ZoneID
		if zoneID, ok := recordMap["zone_id"].(string); ok {
			record.ZoneID = zoneID
		}

		records = append(records, record)
	}

	return records, nil
}

// CreateRecord creates a new DNS record
func (c *Client) CreateRecord(record Record) (*Record, error) {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/records", baseURL), bytes.NewBuffer(recordJSON))
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

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	// Use a map to avoid unmarshaling issues
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Extract record from the response
	recordData, ok := result["record"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format: record field not found or not an object")
	}

	createdRecord := &Record{}

	// Extract ID
	if id, ok := recordData["id"].(string); ok {
		createdRecord.ID = id
	}

	// Extract Type
	if recordType, ok := recordData["type"].(string); ok {
		createdRecord.Type = recordType
	}

	// Extract Name
	if name, ok := recordData["name"].(string); ok {
		createdRecord.Name = name
	}

	// Extract Value
	if value, ok := recordData["value"].(string); ok {
		createdRecord.Value = value
	}

	// Extract TTL
	if ttl, ok := recordData["ttl"].(float64); ok {
		createdRecord.TTL = int(ttl)
	}

	// Extract ZoneID
	if zoneID, ok := recordData["zone_id"].(string); ok {
		createdRecord.ZoneID = zoneID
	}

	return createdRecord, nil
}

// UpdateRecord updates an existing DNS record
func (c *Client) UpdateRecord(record Record) (*Record, error) {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/records/%s", baseURL, record.ID), bytes.NewBuffer(recordJSON))
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

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	// Use a map to avoid unmarshaling issues
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Extract record from the response
	recordData, ok := result["record"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format: record field not found or not an object")
	}

	updatedRecord := &Record{}

	// Extract ID
	if id, ok := recordData["id"].(string); ok {
		updatedRecord.ID = id
	}

	// Extract Type
	if recordType, ok := recordData["type"].(string); ok {
		updatedRecord.Type = recordType
	}

	// Extract Name
	if name, ok := recordData["name"].(string); ok {
		updatedRecord.Name = name
	}

	// Extract Value
	if value, ok := recordData["value"].(string); ok {
		updatedRecord.Value = value
	}

	// Extract TTL
	if ttl, ok := recordData["ttl"].(float64); ok {
		updatedRecord.TTL = int(ttl)
	}

	// Extract ZoneID
	if zoneID, ok := recordData["zone_id"].(string); ok {
		updatedRecord.ZoneID = zoneID
	}

	return updatedRecord, nil
}

// DeleteRecord deletes a DNS record
func (c *Client) DeleteRecord(recordID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/records/%s", baseURL, recordID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Auth-API-Token", c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	return nil
}
