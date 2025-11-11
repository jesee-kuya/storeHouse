-- Church members who contribute offerings
CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    notes TEXT,
    group UUID REFERENCES members_groups(id),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_members_phone ON members(phone_number);
CREATE INDEX idx_full_name ON members(last_name);

-- Full-text search index for member search
CREATE INDEX idx_members_fullname ON members USING gin(to_tsvector('english', first_name || ' ' || last_name || ' ' || COALESCE(other_names, '')));

COMMENT ON TABLE members IS 'Church members who make offerings and contributions';