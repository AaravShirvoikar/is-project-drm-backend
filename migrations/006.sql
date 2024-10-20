CREATE TABLE session_keys (
    id UUID PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    content_id UUID NOT NULL,
    key bytea NOT NULL,
    created_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
);
