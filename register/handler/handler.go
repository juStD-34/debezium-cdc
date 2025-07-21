package handler

import (
	"net/http"
	"register/models"
	"register/pkg/logger"
	"register/service"

	"github.com/gin-gonic/gin"
)

type CDCHandler interface {
	RegisterConnector(c *gin.Context)
	ListConnectors(c *gin.Context)
	GetConnectorStatus(c *gin.Context)
	DeleteConnector(c *gin.Context)
}
type cDCHandler struct {
	service service.CDCRegistrationService
	logger  logger.Logger
}

func NewCDCHandler(service service.CDCRegistrationService, logger logger.Logger) CDCHandler {
	return &cDCHandler{
		service: service,
		logger:  logger,
	}
}

func (h *cDCHandler) RegisterConnector(c *gin.Context) {
	var req models.RegisterConnectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request payload", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Registering connector", logger.String("connector_name", req.ConnectorName))

	response, err := h.service.RegisterConnector(req)
	if err != nil {
		h.logger.Error("Failed to register connector", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Connector registered successfully", logger.String("connector_name", req.ConnectorName))
	c.JSON(http.StatusCreated, response)
}

func (h *cDCHandler) ListConnectors(c *gin.Context) {
	h.logger.Info("Listing connectors")

	response, err := h.service.ListConnectors()
	if err != nil {
		h.logger.Error("Failed to list connectors", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *cDCHandler) GetConnectorStatus(c *gin.Context) {
	connectorName := c.Param("name")

	h.logger.Info("Getting connector status", logger.String("connector_name", connectorName))

	status, err := h.service.GetConnectorStatus(connectorName)
	if err != nil {
		h.logger.Error("Failed to get connector status", logger.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *cDCHandler) DeleteConnector(c *gin.Context) {
	connectorName := c.Param("name")

	h.logger.Info("Deleting connector", logger.String("connector_name", connectorName))

	if err := h.service.DeleteConnector(connectorName); err != nil {
		h.logger.Error("Failed to delete connector", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Connector deleted successfully", logger.String("connector_name", connectorName))
	c.JSON(http.StatusOK, gin.H{"message": "Connector deleted successfully"})
}
