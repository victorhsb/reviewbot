package service

import (
	"context"

	"github.com/google/uuid"
)

type MessageService struct {
	storage MessageStorage
}

func New(storage MessageStorage) *MessageService {
	return &MessageService{storage: storage}
}

// Make sure MessageService implements the MessageWriter and MessageReader interfaces
var _ MessageWriter = (*MessageService)(nil)
var _ MessageReader = (*MessageService)(nil)

// GetMessagesByParticipant returns a list of messages sent or received by a specific user
func (m *MessageService) GetMessagesByParticipant(ctx context.Context, id uuid.UUID) ([]Message, error) {
	return m.storage.ListMessagesByParticipant(ctx, id)
}

// SaveMessage persists the message to the storage and emits any events to further process the message
func (m *MessageService) SaveMessage(ctx context.Context, msg Message) error {
	return m.storage.SaveMessage(ctx, msg)
}
