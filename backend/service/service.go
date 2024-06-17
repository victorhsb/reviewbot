package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type MessageService struct {
	repo MessageRepository
}

// Make sure MessageService implements the MessageWriter and MessageReader interfaces
var _ MessageWriter = (*MessageService)(nil)
var _ MessageReader = (*MessageService)(nil)

// GetMessagesByParticipant returns a list of messages sent or received by a specific user
func (m *MessageService) GetMessagesByParticipant(ctx context.Context, id uuid.UUID) (map[string][]Message, error) {
	messages, err := m.repo.ListMessagesByParticipant(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not list messages by participant: %w", err)
	}

	messageMap := make(map[string][]Message)
	for _, msg := range messages {
		contact := msg.GetOtherParticipant(id).String()
		messageMap[contact] = append(messageMap[contact], msg)
	}

	return messageMap, nil
}

// SaveMessage persists the message to the storage and emits any events to further process the message
func (m *MessageService) SaveMessage(ctx context.Context, msg Message) error {
	return m.repo.SaveMessage(ctx, msg)
}
