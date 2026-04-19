package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ControlCommandLog lưu lại mỗi lần FE gửi lệnh điều khiển.
// Dùng để theo dõi luồng xử lý realtime.
type ControlCommandLog struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Action      string        `bson:"action"        json:"action"`
	Params      string        `bson:"params"        json:"params"`       // raw JSON string
	ReceivedAt  time.Time     `bson:"received_at"   json:"received_at"`  // thời điểm BE nhận
	PublishedAt time.Time     `bson:"published_at"  json:"published_at"` // thời điểm publish lên MQTT
	MQTTTopic   string        `bson:"mqtt_topic"    json:"mqtt_topic"`
}
