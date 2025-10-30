-- Account Transfers
CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES transactions(id) NOT NULL,
    perticulars VARCHAR(255) NOT NULL,
    credit_account UUID REFERENCES accounts(id) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
);

COMMENT ON TABLE transfers IS 'Acoount Transfers';
