package mongo

import (
	"context"
	"time"

	"github.com/utm/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type NotificationRepository struct {
	col *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{col: db.Collection("notifications")}
}

func (r *NotificationRepository) Insert(ctx context.Context, n *model.Notification) error {
	n.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, n)
	return err
}

func (r *NotificationRepository) FindByRecipient(ctx context.Context, recipientID string, unreadOnly bool) ([]*model.Notification, error) {
	filter := bson.M{
		"$or": bson.A{
			bson.M{"recipient_id": recipientID},
			bson.M{"recipient_id": "broadcast"},
		},
	}
	if unreadOnly {
		filter["is_read"] = false
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(50)
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Notification
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = r.col.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"is_read": true, "read_at": now}},
	)
	return err
}
