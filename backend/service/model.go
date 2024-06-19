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

// Product defines the modelling to the product type
type Product struct {
	ID      *uuid.UUID      `json:"id"`
	Title   string          `json:"title"`
	Reviews []ProductReview `json:"reviews"`
}

// ProductReview defines the modelling to the product review type
type ProductReview struct {
	ID        uuid.UUID `json:"id"`
	Review    string    `json:"review"`
	Sentiment int       `json:"sentiment,omitempty"`
	Rating    int       `json:"rating"`
	ProductID uuid.UUID `json:"product_id"`
	UserID    uuid.UUID `json:"user_id"`
}
