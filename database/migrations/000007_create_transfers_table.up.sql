-- Account transfers
CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    particulars VARCHAR(255) NOT NULL,
    credit_account UUID NOT NULL REFERENCES accounts(id),
    amount NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE transfers IS
'Account transfers';
