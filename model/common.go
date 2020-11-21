package model

import "time"

type NotificationMessage struct {
	SourceAppId   string        `json:"source_app_id"`
	CreatedAt     time.Time     `json:"created_at"`
	Type          string	 `json:"type"`
	Destination[]   string	 `json:"destination"`
}


