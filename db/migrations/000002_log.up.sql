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
