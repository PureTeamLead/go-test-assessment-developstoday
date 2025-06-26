CREATE TYPE "state_enum" AS enum('started', 'completed');

CREATE TABLE IF NOT EXISTS missions (
                                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                        cat_id UUID NOT NULL UNIQUE,
                                        state state_enum NOT NULL DEFAULT 'started',
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                        FOREIGN KEY (cat_id) REFERENCES "cats" (id)
);

CREATE TRIGGER update_missions_updated_at
    BEFORE UPDATE ON "missions"
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();