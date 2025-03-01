CREATE TABLE IF NOT EXISTS accounts (
  id bigint PRIMARY KEY,
  owner varchar(100) NOT NULL,
  balance bigint NOT NULL,
  currency varchar(10) NOT NULL,
  created_at timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS entries (
  id bigint PRIMARY KEY,
  account_id bigint NOT NULL,
  CONSTRAINT FK_account_id FOREIGN KEY (account_id) REFERENCES accounts(id),
  amount bigint NOT NULL,
  created_at timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS transfers (
  id bigint PRIMARY KEY,
  from_account_id bigint NOT NULL,
  CONSTRAINT FK_from_acc_id FOREIGN KEY (from_account_id) REFERENCES accounts(id),
  to_account_id bigint NOT NULL,
  CONSTRAINT FK_to_acc_id FOREIGN KEY (to_account_id) REFERENCES accounts(id),
  amount bigint NOT NULL,
  created_at timestamp NOT NULL DEFAULT (now())
);


CREATE INDEX accounts_owner ON accounts (owner);

CREATE INDEX entries_account_id ON entries (account_id);

CREATE INDEX transfers_from_account_id ON transfers (from_account_id);

CREATE INDEX transfers_to_account_id ON transfers (to_account_id);

CREATE INDEX transfers_from_to_account_id ON transfers (from_account_id, to_account_id);