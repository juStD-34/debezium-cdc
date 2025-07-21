package service

import (
	"register/models"
	"register/pkg/http"
	"time"
)

type CDCRegistrationService interface {
	RegisterConnector(req models.RegisterConnectorRequest) (*models.ConnectorResponse, error)
	ListConnectors() (*models.ListConnectorsResponse, error)
	GetConnectorStatus(connectorName string) (*models.ConnectorStatus, error)
	DeleteConnector(connectorName string) error
}

type cDCRegistrationService struct {
	kafkaConnectURL string
	httpClient      *http.HTTPClient
}

func NewCDCRegistrationService(kafkaConnectURL string) CDCRegistrationService {
	return &cDCRegistrationService{
		kafkaConnectURL: kafkaConnectURL,
		httpClient:      &http.Client{Timeout: 30 * time.Second},
	}
}
