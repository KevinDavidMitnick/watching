package graph

import "time"

type MetricEndpoint struct {
	ID         uint `gorm:"primary_key"`
	Metric     string
	EndpointID int
	Ts         int
	TCreate    time.Time
	TModify    time.Time
}

func (MetricEndpoint) TableName() string {
	return "metric_endpoint"
}

