BEGIN TRANSACTION;

CREATE TABLE flight_log_new (
    key TEXT PRIMARY KEY,
    value TEXT,
    last_seen DATETIME
);

INSERT INTO flight_log_new (key, value, last_seen)
SELECT key, value, last_seen FROM flight_log;

DROP TABLE flight_log;

ALTER TABLE flight_log_new RENAME TO flight_log;

COMMIT;
