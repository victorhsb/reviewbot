package service

import (
	"context"

	"github.com/google/uuid"
)

// MessageWriter defines the interface definition for the message writer
type MessageWriter interface {
	// SaveMessage persists the message to the storage and emits any events to further process the message
	SaveMessage(context.Context, Message) error
}

// MessageReader defines the interface definition for the message reader
type MessageReader interface {
	// GetMessagesByParticipant returns a list of messages sent or received by a specific user
	GetMessagesByParticipant(context.Context, uuid.UUID) ([]Message, error)
}

// MessageStorage defines the interface definition for the message persistence layer
type MessageStorage interface {
	SaveMessage(context.Context, Message) error

	ListMessagesByParticipant(context.Context, uuid.UUID) ([]Message, error)
}
