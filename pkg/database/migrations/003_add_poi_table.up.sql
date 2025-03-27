CREATE TYPE poi_type AS ENUM ('subway', 'hospital', 'mall');

-- Create the poi table
CREATE TABLE poi (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type poi_type NOT NULL, -- 'subway', 'hospital', 'mall' 
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP
);
