-- Receipts/ Offerings
CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES transactions(id) NOT NULL,
    income_account UUID REFERENCES accounts(id) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
);

COMMENT ON TABLE receipts IS 'receipts';
