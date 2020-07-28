CREATE TABLE
IF NOT EXISTS metric
(
    id INTEGER NOT NULL PRIMARY KEY,
    name text
);

CREATE TABLE
IF NOT EXISTS ct
(
    id INTEGER NOT NULL PRIMARY KEY,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    metric_id INTEGER NOT NULL,
    value int NOT NULL,
    FOREIGN KEY(metric_id) REFERENCES metric(id)
);

CREATE TABLE
IF NOT EXISTS config 
(
    metric_id INTEGER NOT NULL,
    opt text NOT NULL,
    val text NOT NULL,
    UNIQUE(metric_id, opt)
        ON CONFLICT REPLACE,
    FOREIGN KEY(metric_id) REFERENCES metric(id)
);
