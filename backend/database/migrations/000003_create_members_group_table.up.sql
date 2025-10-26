-- Church members Grops or estates
CREATE TABLE members_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_name VARCHAR(50) NOT NULL,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_group_name ON members_groups(group_name);

-- Full-text search index for group search
CREATE INDEX idx_group_name ON members_groups USING gin(to_tsvector('english', group_name));

-- Full-text search index for member search
CREATE INDEX idx_members_fullname ON members_groups USING gin(to_tsvector('english', group_name));

COMMENT ON TABLE members IS 'members_groups: Church members Grops or estates';