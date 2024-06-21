package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type MessageService struct {
	storage MessageStorage

	processor MessageProcessor
}

func NewMessageService(storage MessageStorage, processor MessageProcessor) *MessageService {
	return &MessageService{storage: storage, processor: processor}
}

// Make sure MessageService implements the Messager interfaces
var _ Messager = (*MessageService)(nil)

// GetUserByID returns the user with the given ID
func (m *MessageService) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	return m.storage.GetUserByID(ctx, id)
}

// ListMessagesByUserID returns a list of messages sent or received by a specific user
func (m *MessageService) ListMessagesByUserID(ctx context.Context, id uuid.UUID) ([]Message, error) {
	return m.storage.ListMessagesByUserID(ctx, id)
}

// SaveMessage persists the message to the storage and emits any events to further process the message
func (m *MessageService) SaveMessage(ctx context.Context, msg Message) error {
	if msg.Direction == "" {
		msg.Direction = DirectionSent
	}

	err := m.storage.SaveMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("could not save message; %w", err)
	}

	if m.processor != nil {
		return m.processor.ProcessMessage(ctx, msg)
	}
	return nil
}

// ListUsers returns a list of all users
func (m *MessageService) ListUsers(ctx context.Context) ([]User, error) {
	return m.storage.ListUsers(ctx)
}
