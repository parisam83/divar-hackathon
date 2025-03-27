Create TABLE travel_metrics(
    id SERIAL PRIMARY KEY,
    distance INT NOT NULL,
    duration INT NOT NULL,
    origin_id VARCHAR(255) NOT NULL,
    destionation_id INTEGER NOT NULL,
    FOREIGN KEY(origin_id) REFERENCES posts(post_id),
    FOREIGN KEY(destionation_id) REFERENCES poi(id)
);