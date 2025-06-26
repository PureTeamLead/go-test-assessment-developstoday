CREATE TABLE IF NOT EXISTS targets (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                       mission_id UUID NOT NULL UNIQUE,
                                       name VARCHAR(30) NOT NULL,
                                       country VARCHAR(50) NOT NULL,
                                       notes TEXT,
                                       state state_enum NOT NULL,
                                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                       updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                       FOREIGN KEY (mission_id) REFERENCES "missions" (id)
);

CREATE TRIGGER update_targets_updated_at
    BEFORE UPDATE ON "targets"
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();