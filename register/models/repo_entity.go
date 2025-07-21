package models

import "time"

type Connector struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ConnectorName string    `gorm:"uniqueIndex" json:"connector_name"`
	DatabaseType  string    `json:"database_type"`
	DatabaseHost  string    `json:"database_host"`
	DatabasePort  int       `json:"database_port"`
	DatabaseName  string    `json:"database_name"`
	TopicPrefix   string    `json:"topic_prefix"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
