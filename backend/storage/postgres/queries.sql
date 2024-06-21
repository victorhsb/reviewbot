-- name: ListMessagesByUser :many
SELECT message, created_at, direction, created_at FROM messages WHERE user_id = $1::uuid ORDER BY created_at;

-- name: SaveMessage :exec
INSERT INTO messages (direction, user_id, message) VALUES ($1, $2, $3);

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 limit 1;

-- name: ListUsers :many
select * from users;

-- name: SaveUser :one
INSERT INTO users (username) VALUES ($1) RETURNING *;

-- name: UpsertProduct :one
INSERT INTO products (id, title) values (@id, @title) ON CONFLICT (id) DO UPDATE SET title = @title WHERE products.id = @id returning id;

-- name: GetProduct :one
SELECT p.*, jsonb_agg(r.*) FROM products p LEFT JOIN product_reviews r ON p.id = r.product_id WHERE p.id = $1 GROUP BY p.id;

-- name: ListProducts :many
SELECT 
    p.*,
    jsonb_agg(u.*)
FROM products p LEFT JOIN (
    select r.product_id, r.review, r.sentiment, r.rating, u.username
    from product_reviews r 
    left join users u 
        on r.user_id = u.id) u 
        on p.id = u.product_id 
GROUP BY p.id 
LIMIT sqlc.arg('limit')::bigint 
OFFSET sqlc.arg('offset')::bigint;

-- name: GetProductReview :one
select * from product_reviews where id = $1 limit 1;

-- name: GetProductReviews :many
SELECT * FROM product_reviews WHERE product_id = $1;

-- name: SaveProductReview :one
INSERT INTO product_reviews (product_id, user_id, rating, sentiment, review) VALUES ($1, $2, @rating::integer, @sentiment::integer, @review::text) RETURNING *;

-- name: UpdateProductReview :exec
UPDATE product_reviews SET rating = @rating::integer, sentiment = @sentiment::integer, review = @review::text WHERE id = @id::uuid;
