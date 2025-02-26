-- Create the metro_stations table
CREATE TABLE metro_stations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL
);

-- Create the users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone_number TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

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

-- Create the results table
CREATE TABLE results (
    id SERIAL PRIMARY KEY,
    token_id INT REFERENCES tokens(id),
    metro_station_id INT REFERENCES metro_stations(id),
    transportation_cost NUMERIC,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);