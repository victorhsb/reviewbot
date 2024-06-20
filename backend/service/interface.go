package service

import (
	"context"

	"github.com/google/uuid"
)

// Messager is an interface that defines the methods for the message service
type Messager interface {
	// SaveMessage persists the message to the storage and emits any events to further process the message
	SaveMessage(context.Context, Message) error
	// ListMessagesByUserID returns a list of messages sent or received by a specific user
	ListMessagesByUserID(context.Context, uuid.UUID) ([]Message, error)
	// GetUserByID returns the user with the given ID
	GetUserByID(context.Context, uuid.UUID) (User, error)
	// ListUsers returns a list of all users
	ListUsers(context.Context) ([]User, error)
}

// ProductReviewer is an interface that defines the methods for the product service
// that range from managing products to registering their reviews
type ProductReviewer interface {
	SaveProduct(context.Context, Product) (*uuid.UUID, error)
	GetProduct(context.Context, uuid.UUID) (Product, error)
	ListProducts(context.Context, int64, int64) ([]*Product, error)

	// SaveProductReview(context.Context, ProductReview) error
	// GetProductReviews(context.Context, uuid.UUID) ([]ProductReview, error)
	// UpdateProductReviewSentiment(context.Context, uuid.UUID, int64) error
}

// Storage groups by all the storage interfaces along with utility methods
type Storage interface {
	MessageStorage
	ProductStorage

	Migrate() error
	Ping(context.Context) error
}

// ProductStorage defines the interface definition for the product&reviews persistence layer
type ProductStorage interface {
	UpsertProduct(context.Context, Product) (*uuid.UUID, error)
	GetProduct(context.Context, uuid.UUID) (Product, error)
	ListProduct(context.Context, int64, int64) ([]*Product, error)
}

// MessageStorage defines the interface definition for the message persistence layer
type MessageStorage interface {
	SaveMessage(context.Context, Message) error
	ListMessagesByUserID(context.Context, uuid.UUID) ([]Message, error)
	GetUserByID(context.Context, uuid.UUID) (User, error)
	ListUsers(context.Context) ([]User, error)
}
