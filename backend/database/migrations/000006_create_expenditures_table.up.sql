-- Expenses and Withdrowals from accounts
CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES transactions(id) NOT NULL,
    perticulars VARCHAR(255) NOT NULL,
    bank_account UUID REFERENCES accounts(id) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
);

COMMENT ON TABLE Expenditure IS 'expenses and withdrawals from accounts';
