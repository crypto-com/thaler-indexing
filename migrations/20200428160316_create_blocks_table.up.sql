/* Avoid foreign key relationship to reduce cyclic dependency */
/* Create tables like block_committed_council_nodes to model such relationship */
CREATE TABLE blocks (
  height BIGINT,
  hash VARCHAR NOT NULL,
  time BIGINT NOT NULL,
  app_hash VARCHAR NOT NULL,
  committed_council_nodes JSONB NULL,
  UNIQUE(hash),
  PRIMARY KEY(height)
);