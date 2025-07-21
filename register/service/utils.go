package service

import (
	"fmt"
	"register/models"
	"strings"
)

func (s *cDCRegistrationService) buildConnectorConfig(req models.RegisterConnectorRequest) (map[string]interface{}, error) {
	// Set defaults
	if req.SnapshotMode == "" {
		req.SnapshotMode = "initial"
	}
	if req.ServerID == 0 {
		req.ServerID = 184054 // Default server ID
	}

	var config map[string]interface{}

	switch strings.ToLower(string(req.DatabaseType)) {
	case "mysql":
		config = map[string]interface{}{
			"name": req.ConnectorName,
			"config": map[string]interface{}{
				"connector.class":       "io.debezium.connector.mysql.MySqlConnector",
				"database.hostname":     req.DatabaseHost,
				"database.port":         fmt.Sprintf("%d", req.DatabasePort),
				"database.user":         req.Username,
				"database.password":     req.Password,
				"database.server.id":    fmt.Sprintf("%d", req.ServerID),
				"topic.prefix":          req.TopicPrefix,
				"database.include.list": req.DatabaseName,
				"table.include.list":    s.formatTableList(req.DatabaseName, req.Tables),
				"schema.history.internal.kafka.bootstrap.servers": "kafka:9092",
				"schema.history.internal.kafka.topic":             fmt.Sprintf("schemahistory.%s", req.DatabaseName),
				"snapshot.mode":                                   req.SnapshotMode,
				"decimal.handling.mode":                           "string",
				"time.precision.mode":                             "connect",
				"include.schema.changes":                          "true",
			},
		}

	case "postgresql":
		config = map[string]interface{}{
			"name": req.ConnectorName,
			"config": map[string]interface{}{
				"connector.class":       "io.debezium.connector.postgresql.PostgreSqlConnector",
				"database.hostname":     req.DatabaseHost,
				"database.port":         fmt.Sprintf("%d", req.DatabasePort),
				"database.user":         req.Username,
				"database.password":     req.Password,
				"database.dbname":       req.DatabaseName,
				"topic.prefix":          req.TopicPrefix,
				"table.include.list":    s.formatTableList("public", req.Tables), // PostgreSQL uses public schema by default
				"plugin.name":           "pgoutput",
				"snapshot.mode":         req.SnapshotMode,
				"decimal.handling.mode": "string",
				"time.precision.mode":   "connect",
			},
		}

	default:
		return nil, fmt.Errorf("unsupported database type: %s", req.DatabaseType)
	}

	// Add custom transforms if provided
	if req.Transforms != nil && len(req.Transforms) > 0 {
		configMap := config["config"].(map[string]interface{})
		for key, value := range req.Transforms {
			configMap[key] = value
		}
	}

	return config, nil
}

func (s *cDCRegistrationService) formatTableList(database string, tables []string) string {
	var formattedTables []string
	for _, table := range tables {
		if !strings.Contains(table, ".") {
			formattedTables = append(formattedTables, fmt.Sprintf("%s.%s", database, table))
		} else {
			formattedTables = append(formattedTables, table)
		}
	}
	return strings.Join(formattedTables, ",")
}

func (s *cDCRegistrationService) getConnectorStatus(connectorName string) (*models.ConnectorStatus, error) {
	url := fmt.Sprintf("%s/connectors/%s/status", s.cfg.ConnectorUrl, connectorName)

	var status models.ConnectorStatus
	if err := s.client.Get(url, &status); err != nil {
		return nil, fmt.Errorf("failed to get connector status for %s: %w", connectorName, err)
	}

	return &status, nil
}

func flattenConfig(config map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range config {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
