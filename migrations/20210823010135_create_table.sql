-- +goose Up
-- +goose StatementBegin
CREATE TABLE contact
(
    id      SERIAL PRIMARY KEY,
    user_id INT,
    type    INT,
    text    TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE contact;
-- +goose StatementEnd
