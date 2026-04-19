package mongo

import (
	"context"

	"github.com/utm/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CommandLogRepository struct {
	col *mongo.Collection
}

func NewCommandLogRepository(db *mongo.Database) *CommandLogRepository {
	return &CommandLogRepository{col: db.Collection("control_command_logs")}
}

func (r *CommandLogRepository) Insert(ctx context.Context, entry *model.ControlCommandLog) error {
	_, err := r.col.InsertOne(ctx, entry)
	return err
}
