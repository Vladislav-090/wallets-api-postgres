CREATE TABLE IF NOT EXISTS transfers(
    id BIGSERIAL PRIMARY KEY,
    from_wallet_id BIGINT NOT NULL 
        REFERENCES wallets(id),
    to_wallet_id BIGINT NOT NULL 
        REFERENCES wallets(id),
    amount NUMERIC(20, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (from_wallet_id <> to_wallet_id)
);