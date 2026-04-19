package mongo

import (
	"context"
	"time"

	"github.com/utm/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TelemetryRepository struct {
	telemetry *mongo.Collection
	sessions  *mongo.Collection
}

func NewTelemetryRepository(db *mongo.Database) *TelemetryRepository {
	return &TelemetryRepository{
		telemetry: db.Collection("telemetry"),
		sessions:  db.Collection("flight_sessions"),
	}
}

func (r *TelemetryRepository) Insert(ctx context.Context, data *model.Telemetry) error {
	data.Timestamp = time.Now()
	_, err := r.telemetry.InsertOne(ctx, data)
	return err
}

func (r *TelemetryRepository) FindBySession(ctx context.Context, sessionID string) ([]*model.Telemetry, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
	cursor, err := r.telemetry.Find(ctx, bson.M{"session_id": sessionID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Telemetry
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *TelemetryRepository) FindActiveSessions(ctx context.Context) ([]*model.FlightSession, error) {
	cursor, err := r.sessions.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.FlightSession
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *TelemetryRepository) StartSession(ctx context.Context, session *model.FlightSession) error {
	session.StartedAt = time.Now()
	session.IsActive = true
	_, err := r.sessions.InsertOne(ctx, session)
	return err
}

func (r *TelemetryRepository) EndSession(ctx context.Context, sessionID string) error {
	now := time.Now()
	_, err := r.sessions.UpdateOne(ctx,
		bson.M{"session_id": sessionID},
		bson.M{"$set": bson.M{"is_active": false, "ended_at": now}},
	)
	return err
}
