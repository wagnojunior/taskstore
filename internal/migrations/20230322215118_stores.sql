-- +goose Up
-- +goose StatementBegin
CREATE TABLE stores (
    id SERIAL PRIMARY KEY,
    name TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stores;
-- +goose StatementEnd
