
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL UNIQUE,
    bookmark_count INT DEFAULT 0,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE URLs (
    url_id SERIAL PRIMARY KEY,
    url_address VARCHAR(256) NOT NULL UNIQUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE UserURLs (
    user_url_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    url_id INT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users (user_id),
    FOREIGN KEY (url_id) REFERENCES URLs (url_id),
    UNIQUE (user_id, url_id)
);

-- Users table trigger function
CREATE OR REPLACE FUNCTION set_timestamp_users()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_users
BEFORE UPDATE ON Users
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_users();

-- URLs table trigger function
CREATE OR REPLACE FUNCTION set_timestamp_urls()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_urls
BEFORE UPDATE ON URLs
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_urls();

-- UserURLs table trigger function
CREATE OR REPLACE FUNCTION set_timestamp_userurls()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_userurls
BEFORE UPDATE ON UserURLs
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_userurls();
