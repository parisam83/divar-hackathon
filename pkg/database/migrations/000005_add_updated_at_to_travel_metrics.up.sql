--  Add a new column to the travel_metrics table
ALTER TABLE travel_metrics
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;