package service

import (
	"time"

	"github.com/google/uuid"
)

// Message is the model definition for what is considered a message
// contains the sender, target, content, timestamp and a footprint that *can* be nil if this is the first message.
type Message struct {
	Content   string     `json:"content"`
	Sender    *uuid.UUID `json:"sender,omitempty"`
	Receiver  *uuid.UUID `json:"receiver,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
	// Footprint is the list of the previous messages exchanged between the sender and the target
	FootPrint []Message `json:"footPrint"`
}

// User defines the modelling to the user type
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}
