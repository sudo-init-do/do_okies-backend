-- Payout destinations a.k.a. bank accounts for withdrawals
CREATE TABLE payout_destinations (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  bank_code       TEXT NOT NULL,
  account_number  TEXT NOT NULL,
  account_name    TEXT NOT NULL,
  is_default      BOOLEAN NOT NULL DEFAULT FALSE,
  status          TEXT NOT NULL DEFAULT 'active',
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_payout_destinations_user ON payout_destinations(user_id, created_at DESC);

-- Wallet accounting
CREATE TABLE transactions (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  idempotency_key  TEXT UNIQUE,
  kind             TEXT NOT NULL,      -- e.g. deposit|withdrawal_reserve|withdrawal_refund|gift|order
  amount           BIGINT NOT NULL,    -- minor units (kobo)
  currency         TEXT NOT NULL DEFAULT 'NGN',
  metadata         JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_kind ON transactions(kind, created_at DESC);

CREATE TABLE ledger_entries (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tx_id      UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
  wallet_id  UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
  direction  TEXT NOT NULL CHECK (direction IN ('debit','credit')),
  amount     BIGINT NOT NULL
);

CREATE INDEX idx_ledger_wallet ON ledger_entries(wallet_id);
CREATE INDEX idx_ledger_tx ON ledger_entries(tx_id);

-- Payouts (withdrawals)
CREATE TABLE payouts (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  destination_id  UUID NOT NULL REFERENCES payout_destinations(id) ON DELETE RESTRICT,
  amount          BIGINT NOT NULL,
  currency        TEXT NOT NULL DEFAULT 'NGN',
  status          TEXT NOT NULL DEFAULT 'pending',  -- pending|approved|succeeded|rejected|failed
  reference       TEXT UNIQUE,                       -- idempotent/business ref
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_payouts_user ON payouts(user_id, created_at DESC);
