package service

import (
	"context"

	"github.com/google/uuid"
)

type MessageService struct {
	storage MessageStorage
}

func NewMessageService(storage MessageStorage) *MessageService {
	return &MessageService{storage: storage}
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
	return m.storage.SaveMessage(ctx, msg)
}

// ListUsers returns a list of all users
func (m *MessageService) ListUsers(ctx context.Context) ([]User, error) {
	return m.storage.ListUsers(ctx)
}
