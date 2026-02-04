-- Expenses and withdrawals from accounts
CREATE TABLE expenditure (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    particulars VARCHAR(255) NOT NULL,
    bank_account UUID NOT NULL REFERENCES accounts(id),
    amount NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE expenditure IS
'Expenses and withdrawals from accounts';
