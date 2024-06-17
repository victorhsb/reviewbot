BEGIN;
    CREATE TABLE products (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        title TEXT NOT NULL,

        created_at TIMESTAMPTZ DEFAULT now()
    );

    CREATE TABLE product_reviews (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        product_id UUID REFERENCES product(id) NOT NULL,
        rating INT CHECK (rating >= 1 AND rating <= 5),
        sentiment INT CHECK (sentiment >= -1 AND sentiment <= 1),
        comment TEXT,

        created_at TIMESTAMPTZ DEFAULT now()
    );

    CREATE TABLE users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        username TEXT NOT NULL UNIQUE,
        role TEXT NOT NULL,

        created_at TIMESTAMPTZ DEFAULT now()
    );

    CREATE TABLE messages (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        receiver_id UUID NOT NULL,
        sender_id UUID NOT NULL,
        message TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT now(),

        CONSTRAINT fk_message_receiver FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
        CONSTRAINT fk_message_sender FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
    );
COMMIT;
