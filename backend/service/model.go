package service

import (
	"time"

	"github.com/google/uuid"
)

// Message is the model definition for what is considered a message
// contains the sender, target, content, timestamp and a footprint that *can* be nil if this is the first message.
type Message struct {
	Content   string
	Sender    uuid.UUID
	Target    uuid.UUID
	Timestamp time.Time
	// Footprint is the list of the previous messages exchanged between the sender and the target
	FootPrint []Message
}

// GetOtherParticipant takes a UUID as an argument, which represents the ID of one participant in a message.
// If the provided ID matches the sender of the message, it returns the target's ID.
// If the provided ID does not match the sender, it returns the sender's ID.
// This function is useful for determining the other participant in a message given one participant's ID.
func (m Message) GetOtherParticipant(id uuid.UUID) uuid.UUID {
	if m.Sender == id {
		return m.Target
	}
	return m.Sender
}
