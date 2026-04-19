-- +goose Up
CREATE TABLE games (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    name TEXT NOT NULL,
    created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_date TIMESTAMP,
    finished_date TIMESTAMP,
    played BOOLEAN NOT NULL DEFAULT false,
    finished BOOLEANNOT NULL DEFAULT false,
    description TEXT NOT NULL DEFAULT '',
    download_only BOOLEAN NOT NULL DEFAULT false,
    rating INT CHECK (
        rating >= 0
        AND rating <= 100
    ),
    release_date DATE
);

-- +goose Down
DROP TABLE games;