-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id UUID not null,
    expires_at timestamp not null,
    revoked_at timestamp,

    CONSTRAINT fk_users 
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
