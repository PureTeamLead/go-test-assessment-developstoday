CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS cats (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    name VARCHAR(30) NOT NULL,
                                    experience INT,
                                    breed TEXT NOT NULL,
                                    salary INT NOT NULL,
                                    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX "name_idx" ON "cats" ("name");

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_cats_updated_at
    BEFORE UPDATE ON "cats"
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();