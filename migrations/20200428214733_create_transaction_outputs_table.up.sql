CREATE TABLE transaction_outputs (
  txid VARCHAR,
  index INTEGER NOT NULL,
  spent_at_txid VARCHAR NULL,
  PRIMARY KEY(txid, index),
  FOREIGN KEY(spent_at_txid) REFERENCES activities(txid)
);