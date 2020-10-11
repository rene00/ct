CREATE TABLE
IF NOT EXISTS log_comment
(
    log_id INTEGER NOT NULL,
    comment text NOT NULL,
    UNIQUE(log_id),
    FOREIGN KEY(log_id) REFERENCES log(id)
);
