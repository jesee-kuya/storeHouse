-- Transactions head
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_ref VARCHAR(20),
    transaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
    transaction_type VARCHAR(50) NOT NULL
        CHECK (transaction_type IN ('receipts', 'withdrawal', 'expenses', 'transfer')), 
    amount NUMERIC(10, 2) NOT NULL,
    notes TEXT,
    debit_account UUID REFERENCES accounts(id) NOT NULL,
    member UUID REFERENCES members(id),
    created_by UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE transactions IS 'Transactions head';