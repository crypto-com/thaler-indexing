CREATE TABLE block_rewards (
    block_height BIGINT,
    minted NUMERIC(20) NOT NULL,
    PRIMARY KEY(block_height)
)