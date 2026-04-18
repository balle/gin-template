-- +goose Up
CREATE TABLE games (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    name text NOT NULL,
    created_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_date timestamp,
    finished_date timestamp,
    played boolean NOT NULL DEFAULT false,
    finished boolean NOT NULL DEFAULT false,
    description text NOT NULL DEFAULT '',
    download_only boolean NOT NULL DEFAULT false,
    rating int CHECK (
        rating >= 0
        AND rating <= 100
    ),
    release_date date
);

-- +goose Down
DROP TABLE games;