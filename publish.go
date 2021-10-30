package pqstream

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func PublishEvent(ctx context.Context, db *sql.DB, topic string, payload map[string]interface{}) error {
	data, err := json.Marshal(Event{
		ID:        uuid.New().String(),
		Topic:     topic,
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	// TODO consider using prepared statement here.
	_, err = db.ExecContext(
		ctx,
		fmt.Sprintf("notify eventstream, '%s'", string(data)),
	)
	return err
}

func MustPublishEvent(ctx context.Context, db *sql.DB, topic string, payload map[string]interface{}) {
	err := PublishEvent(ctx, db, topic, payload)
	if err != nil {
		panic(fmt.Errorf("Failed to publish event: %s\n", err.Error()))
	}
}
