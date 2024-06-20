package service

import (
	"time"

	"github.com/google/uuid"
)

type Direction string

var (
	DirectionSent     Direction = "sent"
	DirectionReceived Direction = "received"
)

// Message is the model definition for what is considered a message
// contains the sender, target, content, timestamp and a footprint that *can* be nil if this is the first message.
type Message struct {
	Message   string     `json:"message"`
	UserID    *uuid.UUID `json:"userID,omitempty"`
	Direction Direction  `json:"direction"`
	Timestamp time.Time  `json:"timestamp,omitempty"`
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
	ID        *uuid.UUID `json:"id,omitempty"`
	Review    string     `json:"review"`
	Sentiment int        `json:"sentiment,omitempty"`
	Rating    int        `json:"rating"`
	ProductID *uuid.UUID `json:"product_id,omitempty"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Username  string     `json:"username,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
}
