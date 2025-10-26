-- Transaction Accounts
CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_name VARCHAR(100) NOT NULL,
    local_share NUMERIC(5,4) CHECK (discount >= 0 AND discount <= 1),
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_account_name ON members(phone_number);

-- Full-text search index for account_name search
CREATE INDEX idx_account_name ON Accounts.

COMMENT ON TABLE Accounts IS 'Transaction Accounts';