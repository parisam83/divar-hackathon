-- Create the poi table
CREATE TABLE poi (
    id SERIAL PRIMARY KEY,
    post_token TEXT NOT NULL REFERENCES posts(post_token),
    place_id INT REFERENCES places(id),
    distance NUMERIC,
    duration NUMERIC,
    snapp_cost NUMERIC,
    tapsi_cost NUMERIC,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
