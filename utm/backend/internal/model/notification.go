package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationType string

const (
	NotifTypeAuthApproved NotificationType = "auth_approved" // Cấp phép thành công
	NotifTypeAuthRejected NotificationType = "auth_rejected" // Cấp phép bị từ chối
	NotifTypeAuthExpiring NotificationType = "auth_expiring" // Sắp hết hạn phép
	NotifTypeZoneAlert    NotificationType = "zone_alert"    // Cảnh báo vùng bay
	NotifTypeWeatherAlert NotificationType = "weather_alert" // Cảnh báo thời tiết
	NotifTypeSystemInfo   NotificationType = "system_info"   // Thông báo hệ thống
	NotifTypeViolation    NotificationType = "violation"     // Vi phạm phát hiện
)

type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "low"
	PriorityMedium   NotificationPriority = "medium"
	PriorityHigh     NotificationPriority = "high"
	PriorityCritical NotificationPriority = "critical"
)

// Notification - thông báo hệ thống (lưu MongoDB)
type Notification struct {
	ID          bson.ObjectID        `bson:"_id,omitempty" json:"id"`
	RecipientID string               `bson:"recipient_id" json:"recipient_id"` // userID hoặc "broadcast"
	Type        NotificationType     `bson:"type" json:"type"`
	Priority    NotificationPriority `bson:"priority" json:"priority"`
	Title       string               `bson:"title" json:"title"`
	Body        string               `bson:"body" json:"body"`
	Metadata    map[string]any       `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IsRead      bool                 `bson:"is_read" json:"is_read"`
	ReadAt      *time.Time           `bson:"read_at,omitempty" json:"read_at,omitempty"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	ExpiresAt   *time.Time           `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
}
