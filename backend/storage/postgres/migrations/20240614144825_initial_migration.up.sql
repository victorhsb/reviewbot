BEGIN;
    CREATE TABLE products (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        title TEXT NOT NULL,

        created_at TIMESTAMPTZ DEFAULT now()
    );

    CREATE TABLE users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        username TEXT NOT NULL UNIQUE,

        created_at TIMESTAMPTZ DEFAULT now()
    );

    CREATE TABLE product_reviews (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        product_id UUID REFERENCES products(id) NOT NULL,
        user_id UUID REFERENCES users(id),
        rating INT CHECK (rating >= 1 AND rating <= 5),
        sentiment INT CHECK (sentiment >= -1 AND sentiment <= 1),
        review TEXT,

        created_at TIMESTAMPTZ DEFAULT now()
    );

    CREATE TYPE direction AS ENUM ('sent', 'received');

    CREATE TABLE messages (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
        direction direction NOT NULL,
        message TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT now()
    );
COMMIT;
