-- name: ListMessagesByParticipant :many
SELECT * FROM messages WHERE sender_id = $1::uuid OR receiver_id = $1::uuid ORDER BY created_at;

-- name: SaveMessage :exec
INSERT INTO messages (receiver_id, sender_id, message) VALUES ($1, $2, $3);

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: SaveUser :one
INSERT INTO users (username) VALUES ($1) RETURNING *;

-- name: UpsertProduct :one
INSERT INTO products (id, title) values (@id, @title) ON CONFLICT (id) DO UPDATE SET title = @title WHERE products.id = @id returning id;

-- name: GetProduct :one
SELECT p.*, json_agg(r.*) FROM products p LEFT JOIN product_reviews r ON p.id = r.product_id WHERE p.id = $1 GROUP BY p.id;

-- name: ListProducts :many
SELECT p.*, json_agg(r.*) FROM products p LEFT JOIN product_reviews r ON p.id = r.product_id GROUP BY p.id LIMIT sqlc.arg('limit')::bigint OFFSET sqlc.arg('offset')::bigint;

-- name: GetProductReviews :many
SELECT * FROM product_reviews WHERE product_id = $1;

-- name: SaveProductReview :one
INSERT INTO product_reviews (product_id, user_id, rating, sentiment, review) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateProductSentiment :exec
UPDATE product_reviews SET sentiment = $1 WHERE id = $2;
