package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"register/models"
	"time"
)

// Register a new connector
func (s *cDCRegistrationService) RegisterConnector(req models.RegisterConnectorRequest) (*models.ConnectorResponse, error) {
	log.Printf("Registering connector: %s for %s database", req.ConnectorName, req.DatabaseType)

	// Build connector configuration based on database type
	config, err := s.buildConnectorConfig(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build connector config: %w", err)
	}

	// Check if connector already exists
	exists, err := s.connectorExists(req.ConnectorName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if connector exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("connector %s already exists", req.ConnectorName)
	}

	// Create connector via Kafka Connect REST API
	if err := s.createConnector(config); err != nil {
		return nil, fmt.Errorf("failed to create connector: %w", err)
	}

	// Wait a moment and check status
	time.Sleep(2 * time.Second)
	status, err := s.getConnectorStatus(req.ConnectorName)
	if err != nil {
		log.Printf("Warning: Failed to get connector status: %v", err)
	}

	response := &models.ConnectorResponse{
		ConnectorName: req.ConnectorName,
		Status:        "created",
		Config:        flattenConfig(config["config"].(map[string]interface{})),
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	if status != nil {
		response.Status = status.Connector.State
	}

	return response, nil
}

// List all connectors
func (s *cDCRegistrationService) ListConnectors() (*models.ListConnectorsResponse, error) {
	url := fmt.Sprintf("%s/connectors", s.kafkaConnectURL)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get connectors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kafka connect returned status: %d", resp.StatusCode)
	}

	var connectors []string
	if err := json.NewDecoder(resp.Body).Decode(&connectors); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &models.ListConnectorsResponse{Connectors: connectors}, nil
}

// Get connector status
func (s *cDCRegistrationService) GetConnectorStatus(connectorName string) (*models.ConnectorStatus, error) {
	return s.getConnectorStatus(connectorName)
}

// Delete connector
func (s *cDCRegistrationService) DeleteConnector(connectorName string) error {
	url := fmt.Sprintf("%s/connectors/%s", s.kafkaConnectURL, connectorName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete connector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to delete connector, status: %d", resp.StatusCode)
	}

	log.Printf("Connector %s deleted successfully", connectorName)
	return nil
}
