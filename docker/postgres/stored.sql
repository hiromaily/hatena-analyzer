CREATE OR REPLACE PROCEDURE public.bulk_insert_urls(_urls text[], _categories text[])
LANGUAGE plpgsql
AS $$
DECLARE
    i INT;
BEGIN
    FOR i IN 1 .. array_length(_urls, 1) LOOP
        INSERT INTO URLs (url_address, category_code)
        VALUES (_urls[i], _categories[i])
        ON CONFLICT (url_address) DO NOTHING;
    END LOOP;
END;
$$;
