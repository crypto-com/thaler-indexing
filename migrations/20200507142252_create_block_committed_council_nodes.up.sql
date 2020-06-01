CREATE TABLE block_committed_council_nodes (
    block_height BIGINT,
    council_node_id INT,
    signature VARCHAR,
    is_proposer BOOL,
    PRIMARY KEY(block_height, council_node_id),
    FOREIGN KEY(block_height) REFERENCES blocks(height),
    FOREIGN KEY(council_node_id) REFERENCES council_nodes(id)
)