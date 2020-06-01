CREATE TABLE staking_accounts (
  address VARCHAR,
  nonce BIGINT NOT NULL DEFAULT 0,
  bonded NUMERIC(20) NOT NULL DEFAULT 0,
  unbonded NUMERIC(20) NOT NULL DEFAULT 0,
  unbonded_from BIGINT NULL,
  jailed_until BIGINT NULL,
  punishment_kind PUNISHMENT_KIND NULL,
  current_council_node_id INTEGER NULL,
  PRIMARY KEY(address),
  FOREIGN KEY(current_council_node_id) REFERENCES council_nodes(id)
);