CREATE TABLE IF NOT EXISTS queues
(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    target VARCHAR NOT NULL, -- target table: USER, ...
    target_id INTEGER NOT NULL,
    trigger VARCHAR NOT NULL, -- ONLINE, ACTION, ETC.
    payload text NOT NULL, -- json item
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT NULL
);