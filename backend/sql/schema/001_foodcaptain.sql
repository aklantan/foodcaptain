-- +goose Up
CREATE TABLE restraunts(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    created_at TIMESTAMP NOT NULL,
    uodated_at TIMESTAMP NOT NULL,
    restraunt_name TEXT NOT NULL,
    cuisine TEXT,


);

-- +goose Down 
DROP TABLE restraunts;