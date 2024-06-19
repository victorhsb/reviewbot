package service

import (
	"context"

	"github.com/google/uuid"
)

type ProductServiceConfig struct {
	DefaultListingLimit int64
}

type ProductService struct {
	storage ProductStorage

	config ProductServiceConfig
}

func NewProductService(storage ProductStorage, cfg ProductServiceConfig) ProductReviewer {
	return &ProductService{storage: storage, config: cfg}
}

// SaveProduct saves a new product to the database.
func (p *ProductService) SaveProduct(ctx context.Context, prod Product) (*uuid.UUID, error) {
	id, err := p.storage.UpsertProduct(ctx, prod)
	if err != nil {
		return nil, err
	}

	return id, nil
}

// GetProduct receives a context and a product ID as parameters and returns the product (or an error if anything goes wrong)
func (p *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (Product, error) {
	return p.storage.GetProduct(ctx, id)
}

// ListProduct returns a paginated list of products with their reviews.
// It accepts a context, a limit and a page (1-indexed) as parameters and returns a list of products
// and an error if the operation fails
func (p *ProductService) ListProducts(ctx context.Context, limit int64, page int64) ([]*Product, error) {
	if limit <= 0 {
		limit = p.config.DefaultListingLimit
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	return p.storage.ListProduct(ctx, limit, offset)
}
