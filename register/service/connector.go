package service

import (
	"fmt"
	"go.uber.org/zap"
	"register/models"
	"time"
)

// Register a new connector
func (s *cDCRegistrationService) RegisterConnector(req models.RegisterConnectorRequest) (*models.ConnectorResponse, error) {
	s.log.Info("Registering connector: %s for %s database", zap.Any("connector", req.ConnectorName), zap.Any("db", req.DatabaseType))

	// Build connector configuration based on database type
	config, err := s.buildConnectorConfig(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build connector config: %w", err)
	}

	// Create connector via Kafka Connect REST API
	var createResp interface{}
	createURL := fmt.Sprintf("%s/connectors", s.cfg.ConnectorUrl)
	if err := s.client.Post(createURL, config, &createResp); err != nil {
		return nil, fmt.Errorf("failed to create connector: %w", err)
	}

	// Wait briefly and check status
	time.Sleep(2 * time.Second)
	status, err := s.getConnectorStatus(req.ConnectorName)
	if err != nil {
		s.log.Warn("Failed to get connector status after creation: %v", zap.Any("err", err))
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
	url := fmt.Sprintf("%s/connectors", s.cfg.ConnectorUrl)

	var connectors []string
	if err := s.client.Get(url, &connectors); err != nil {
		return nil, fmt.Errorf("failed to get connectors: %w", err)
	}

	return &models.ListConnectorsResponse{
		Connectors: connectors,
	}, nil
}

// Get connector status
func (s *cDCRegistrationService) GetConnectorStatus(connectorName string) (*models.ConnectorStatus, error) {
	url := fmt.Sprintf("%s/connectors/%s/status", s.cfg.ConnectorUrl, connectorName)

	var status models.ConnectorStatus
	if err := s.client.Get(url, &status); err != nil {
		return nil, fmt.Errorf("failed to get status for connector %s: %w", connectorName, err)
	}

	return &status, nil
}

// Delete connector
func (s *cDCRegistrationService) DeleteConnector(connectorName string) error {
	url := fmt.Sprintf("%s/connectors/%s", s.cfg.ConnectorUrl, connectorName)

	if err := s.client.Delete(url); err != nil {
		return fmt.Errorf("failed to delete connector %s: %w", connectorName, err)
	}

	s.log.Info("Connector %s deleted successfully", zap.String("connector", connectorName))
	return nil
}
