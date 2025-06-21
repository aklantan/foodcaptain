-- +goose Up
CREATE TABLE restaurants(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    restaurant_name TEXT NOT NULL,
    cuisine TEXT
);

-- +goose Down 
DROP TABLE restaurants;