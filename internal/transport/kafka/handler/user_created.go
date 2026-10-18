package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"marketing/internal/transport/kafka/message"
)

type UserCreatedPayload struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

type UserCreatedHandler struct {
	logger Logger
}

func NewUserCreatedHandler(
	logger Logger,
) *UserCreatedHandler {
	return &UserCreatedHandler{
		logger: logger,
	}
}

func (uch *UserCreatedHandler) Handle(
	ctx context.Context,
	msg message.Message,
) error {
	var payload UserCreatedPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("bad payload: %w", err)
	}

	uch.logger.Info(
		ctx,
		"user created message got",
		nil,
	)

	return nil
}
