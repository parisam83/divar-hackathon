-- Create the tokens table
CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    user_token TEXT NOT NULL,
    post_token TEXT NOT NULL,
    user_token_expiry TIMESTAMP WITH TIME ZONE NOT NULL,
    post_token_expiry TIMESTAMP WITH TIME ZONE NOT NULL,
    is_processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_token, post_token)
);