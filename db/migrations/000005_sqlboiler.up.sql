ALTER TABLE config RENAME TO old_config;

CREATE TABLE
IF NOT EXISTS config 
(
    id INTEGER NOT NULL PRIMARY KEY,
    metric_id INTEGER NOT NULL,
    opt text NOT NULL,
    val text NOT NULL,
    UNIQUE(metric_id, opt)
        ON CONFLICT REPLACE,
    FOREIGN KEY(metric_id) REFERENCES metric(id)
);

INSERT INTO config SELECT NULL, metric_id, opt, val FROM old_config;

DROP TABLE old_config;

ALTER TABLE log_comment RENAME TO old_log_comment;

CREATE TABLE
IF NOT EXISTS log_comment
(
    id INTEGER NOT NULL PRIMARY KEY,
    log_id INTEGER NOT NULL,
    comment text NOT NULL,
    UNIQUE(log_id),
    FOREIGN KEY(log_id) REFERENCES log(id)
);

INSERT INTO log_comment SELECT NULL, log_id, comment FROM old_log_comment;

DROP TABLE old_log_comment;
