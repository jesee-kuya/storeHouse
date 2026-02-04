-- Church members groups or estates
CREATE TABLE members_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_name VARCHAR(50) NOT NULL,
    notes TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- B-tree index for exact / prefix lookups
CREATE INDEX idx_members_groups_group_name
    ON members_groups(group_name);

-- Full-text search index
CREATE INDEX idx_members_groups_search
    ON members_groups
    USING gin(to_tsvector('english', group_name));

COMMENT ON TABLE members_groups IS
'Church member groups / estates';
