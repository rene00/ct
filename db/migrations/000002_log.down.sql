CREATE TABLE
IF NOT EXISTS ct
(
    id INTEGER NOT NULL PRIMARY KEY,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    metric_id INTEGER NOT NULL,
    value int NOT NULL,
    FOREIGN KEY(metric_id) REFERENCES metric(id)
);

INSERT INTO ct(id, metric_id, value, timestamp) SELECT NULL, metric_id, value, timestamp FROM log;

DROP TABLE IF EXISTS log;
