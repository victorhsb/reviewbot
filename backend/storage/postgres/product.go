package postgres

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/victorhsb/review-bot/backend/service"
	"github.com/victorhsb/review-bot/backend/storage/postgres/sqlc"
)

// UpsertProduct saves a new product to the database. if the product already exists (and it conflicts with the ID) it'll update it instead.
func (c *client) UpsertProduct(ctx context.Context, prod service.Product) (*uuid.UUID, error) {
	id, err := sqlc.New(c.conn).UpsertProduct(ctx, sqlc.UpsertProductParams{Title: prod.Title, ID: *prod.ID})
	if err != nil {
		return nil, fmt.Errorf("could not save product; %w", err)
	}

	return &id, nil
}

// GetProduct is a method that retrieves a product from the database.
// It takes a context and a product ID as parameters.
// The product ID is of type uuid.UUID which is the unique identifier of the product.
// It returns a product of type service.Product and an error if the operation fails.
func (c *client) GetProduct(ctx context.Context, id uuid.UUID) (service.Product, error) {
	p, err := sqlc.New(c.conn).GetProduct(ctx, id)
	if err != nil {
		return service.Product{}, fmt.Errorf("could not get product; %w", err)
	}

	reviews := make([]service.ProductReview, 0)
	err = json.NewDecoder(bytes.NewReader(p.JsonbAgg)).Decode(&reviews)
	if err != nil {
		return service.Product{}, fmt.Errorf("could not decode product reviews; %w", err)
	}

	return service.Product{
		Title: p.Title,
		ID:    &p.ID,
	}, nil
}

// ListProducts returns a paginated list of products with their reviews.
func (c *client) ListProduct(ctx context.Context, limit int64, offset int64) ([]*service.Product, error) {
	ps, err := sqlc.New(c.conn).ListProducts(ctx, sqlc.ListProductsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, fmt.Errorf("could not list products; %w", err)
	}

	products := make([]*service.Product, len(ps))
	for i, p := range ps {
		reviews := make([]service.ProductReview, 0)
		if string(p.JsonbAgg) != "[null]" { // manually skip jsonb_agg null values for the sake of simplicity
			err = json.NewDecoder(bytes.NewReader(p.JsonbAgg)).Decode(&reviews)
			if err != nil {
				return nil, fmt.Errorf("could not decode product reviews; %w", err)
			}
		}

		products[i] = &service.Product{
			Title:   p.Title,
			ID:      &p.ID,
			Reviews: reviews,
		}
	}

	return products, nil
}

func (c *client) SaveProductReview(ctx context.Context, pr service.ProductReview) (service.ProductReview, error) {
	rev, err := sqlc.New(c.conn).SaveProductReview(ctx, sqlc.SaveProductReviewParams{
		ProductID: *pr.ProductID,
		UserID:    pr.UserID,
		Rating:    int32(pr.Rating),
		Sentiment: int32(pr.Sentiment),
		Review:    pr.Review,
	})

	return service.ProductReview{
		ID:        &rev.ID,
		ProductID: &rev.ProductID,
		UserID:    rev.UserID,
		Review:    rev.Review.String,
		Rating:    int(rev.Rating.Int32),
		Sentiment: int(rev.Sentiment.Int32),
	}, err
}

func (c *client) UpdateProductReview(ctx context.Context, pr service.ProductReview) error {
	return sqlc.New(c.conn).UpdateProductReview(ctx, sqlc.UpdateProductReviewParams{
		ID:        *pr.ID,
		Rating:    int32(pr.Rating),
		Review:    pr.Review,
		Sentiment: int32(pr.Sentiment),
	})
}

func (c *client) GetProductReview(ctx context.Context, id uuid.UUID) (*service.ProductReview, error) {
	rev, err := sqlc.New(c.conn).GetProductReview(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not get product review; %w", err)
	}

	return &service.ProductReview{
		ID:        &rev.ID,
		ProductID: &rev.ProductID,
		UserID:    rev.UserID,
		Review:    rev.Review.String,
		Rating:    int(rev.Rating.Int32),
		Sentiment: int(rev.Sentiment.Int32),
	}, nil
}
