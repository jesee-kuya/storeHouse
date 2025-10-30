-- Transaction Accounts
CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_name VARCHAR(100) NOT NULL,
    account_type VARCHAR(50) NOT NULL 
        CHECK (account_type IN ('Bank', 'Expense', 'Income', 'Asset', 'liability')),
    local_share NUMERIC(5,4) CHECK (discount >= 0 AND discount <= 1),
    notes TEXT,
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_account_name ON members(phone_number);


