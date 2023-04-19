-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    text TEXT,
    tags TEXT[],
    due TEXT,
    store_id INT,
    FOREIGN KEY (store_id) REFERENCES stores (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
