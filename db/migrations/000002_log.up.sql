CREATE TABLE
IF NOT EXISTS log
(
    id INTEGER NOT NULL PRIMARY KEY,
    metric_id INTEGER NOT NULL,
    value int NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(metric_id) REFERENCES metric(id),
    UNIQUE(metric_id,timestamp)
);

INSERT INTO log(id, metric_id, value, timestamp) SELECT NULL, metric_id, value, timestamp FROM ct;

DROP TABLE ct;
