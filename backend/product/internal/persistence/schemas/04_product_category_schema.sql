CREATE TABLE "product_categories" (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    business_id uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by uuid NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by uuid,
    deleted_at timestamptz,
    deleted_by uuid,
    -- FOREIGN KEY (business_id) REFERENCES "businesses" (id) ON DELETE CASCADE,
    -- FOREIGN KEY (created_by) REFERENCES "users" (id) ON DELETE CASCADE,
    -- FOREIGN KEY (updated_by) REFERENCES "users" (id) ON DELETE CASCADE,
    -- FOREIGN KEY (deleted_by) REFERENCES "users" (id) ON DELETE CASCADE,
    PRIMARY KEY (id)
);