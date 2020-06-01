CREATE TABLE council_nodes (
  id BIGSERIAL,
  name VARCHAR NOT NULL,
  security_contact VARCHAR NULL,
  pubkey_type PUBKEY_TYPE NOT NULL,
  pubkey VARCHAR NOT NULL,
  address VARCHAR NOT NULL,
  created_at_block_height BIGINT NOT NULL,
  last_left_at_block_height BIGINT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY(created_at_block_height) REFERENCES blocks(height),
  FOREIGN KEY(last_left_at_block_height) REFERENCES blocks(height)
);