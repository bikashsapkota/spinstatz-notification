package model

import "time"

type NotificationMessage struct {
	SourceAppId string    `json:"source_app_id"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
	Destination []string  `json:"destination"`

	UserId               int               `json:"user_id"`
	Message              string            `json:"message"`
	Href                 string            `json:"href"`
	ReferenceId          int               `json:"reference_id"`
	Subject              string            `json:"subject"`
	Image                string            `json:"image"`
	Seen                 bool              `json:"seen"`
	OneSignalToken       string            `json:"one_signal_token"`
	Data                 map[string]string `json:"data"`
	EmailData            map[string]string `json:"email_data"`
	BulkNotificationData map[string]string `json:"bulk_notification_data"`
}
