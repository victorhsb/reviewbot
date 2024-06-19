package service

import (
	"time"

	"github.com/google/uuid"
)

// Message is the model definition for what is considered a message
// contains the sender, target, content, timestamp and a footprint that *can* be nil if this is the first message.
type Message struct {
	Message   string     `json:"message"`
	Sender    *uuid.UUID `json:"sender,omitempty"`
	Target    *uuid.UUID `json:"target,omitempty"`
	Timestamp time.Time  `json:"timestamp,omitempty"`
	// Footprint is the list of the previous messages exchanged between the sender and the target
	FootPrint []Message `json:"footPrint,omitempty"`
}

// User defines the modelling to the user type
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}
