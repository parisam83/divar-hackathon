-- Create the results table
CREATE TABLE results (
    id SERIAL PRIMARY KEY,
    token_id INT REFERENCES tokens(id),
    metro_station_id INT REFERENCES metro_stations(id),
    transportation_cost NUMERIC,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
