package models

type Database string

const (
	MYSQL    Database = "mysql"
	POSTGRES Database = "postgres"
)

// Request models
type RegisterConnectorRequest struct {
	ConnectorName string            `json:"connector_name" binding:"required"`
	DatabaseType  Database          `json:"database_type" binding:"required"` // mysql, postgresql
	DatabaseHost  string            `json:"database_host" binding:"required"`
	DatabasePort  int               `json:"database_port" binding:"required"`
	DatabaseName  string            `json:"database_name" binding:"required"`
	Username      string            `json:"username" binding:"required"`
	Password      string            `json:"password" binding:"required"`
	TopicPrefix   string            `json:"topic_prefix" binding:"required"`
	Tables        []string          `json:"tables" binding:"required"`
	SnapshotMode  string            `json:"snapshot_mode,omitempty"` // initial, never, when_needed
	ServerID      int               `json:"server_id,omitempty"`     // for MySQL
	Transforms    map[string]string `json:"transforms,omitempty"`    // custom transforms
}

// Response models
type ConnectorResponse struct {
	ConnectorName string            `json:"connector_name"`
	Status        string            `json:"status"`
	Config        map[string]string `json:"config"`
	CreatedAt     string            `json:"created_at"`
}

type ConnectorStatus struct {
	Name      string `json:"name"`
	Connector struct {
		State    string `json:"state"`
		WorkerID string `json:"worker_id"`
	} `json:"connector"`
	Tasks []struct {
		ID       int    `json:"id"`
		State    string `json:"state"`
		WorkerID string `json:"worker_id"`
	} `json:"tasks"`
}

type ListConnectorsResponse struct {
	Connectors []string `json:"connectors"`
}
