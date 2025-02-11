CREATE OR REPLACE PROCEDURE manage_user(user_name_input VARCHAR)
LANGUAGE plpgsql
AS $$
DECLARE
    user_exists BOOLEAN;
    is_user_deleted BOOLEAN;
BEGIN
    -- Check if user exists and if it is marked deleted
    SELECT EXISTS(SELECT 1 FROM Users WHERE user_name = user_name_input) 
    INTO user_exists;

    IF user_exists THEN
        SELECT is_deleted
        INTO is_user_deleted
        FROM Users WHERE user_name = user_name_input;

        IF is_user_deleted THEN 
            -- If user is deleted, update the record
            UPDATE Users
            SET is_deleted = FALSE, updated_at = CURRENT_TIMESTAMP
            WHERE user_name = user_name_input;
        END IF;
    ELSE 
        -- If user does not exist, insert a new record
        INSERT INTO Users (user_name) VALUES (user_name_input);
    END IF;
END $$;
