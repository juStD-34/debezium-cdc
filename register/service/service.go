package service

import (
	"register/config"
	"register/models"
	"register/pkg/http"
	"register/pkg/logger"
)

type CDCRegistrationService interface {
	RegisterConnector(req models.RegisterConnectorRequest) (*models.ConnectorResponse, error)
	ListConnectors() (*models.ListConnectorsResponse, error)
	GetConnectorStatus(connectorName string) (*models.ConnectorStatus, error)
	DeleteConnector(connectorName string) error
}

type cDCRegistrationService struct {
	cfg    *config.Config
	log    logger.Logger
	client http.HTTPClient
}

func NewCDCRegistrationService(cfg *config.Config, log logger.Logger, c http.HTTPClient) CDCRegistrationService {
	return &cDCRegistrationService{
		cfg:    cfg,
		log:    log,
		client: c,
	}
}
