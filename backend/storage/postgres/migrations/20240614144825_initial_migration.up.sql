BEGIN;
    CREATE TABLE products (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        title TEXT NOT NULL, -- product title

        created_at TIMESTAMPZ DEFAULT now()
    );

    CREATE TABLE product_reviews (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        product_id UUID REFERENCES product(id) NOT NULL,
        rating INT CHECK (rating >= 1 AND rating <= 5), -- product rating
        sentiment INT CHECK (sentiment >= -1 AND sentiment <= 1), -- sentiment analysis result
        comment TEXT,

        created_at TIMESTAMPZ DEFAULT now()
    );

    CREATE TABLE users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        username TEXT NOT NULL UNIQUE, -- unique username
        role TEXT NOT NULL, -- user role

        created_at TIMESTAMPZ DEFAULT now()
    );

    CREATE TABLE messages (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        receiver_id UUID REFERENCES users(id), -- receiver user id
        sender_id UUID REFERENCES users(id), -- sender user id
        message TEXT NOT NULL,

        created_at TIMESTAMPZ DEFAULT now()
    );
COMMIT;
