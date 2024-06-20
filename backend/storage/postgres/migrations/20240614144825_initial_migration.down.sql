BEGIN;
    DROP TABLE messages CASCADE;
    DROP TABLE users CASCADE;
    DROP TABLE product_reviews CASCADE;
    DROP TABLE products CASCADE;
    DROP TYPE direction;
COMMIT;
