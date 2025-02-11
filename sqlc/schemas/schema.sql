
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL UNIQUE,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE URLs (
    url_id SERIAL PRIMARY KEY,
    url_address VARCHAR(256) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE UserURLs (
    user_url_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    url_id INT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES Users (user_id),
    FOREIGN KEY (url_id) REFERENCES URLs (url_id)
);
