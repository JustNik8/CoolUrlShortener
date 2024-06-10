CREATE TABLE IF NOT EXISTS url_events
(
    long_url   String,
    short_url  String,
    event_time TIMESTAMP,
    event_type Enum8('create' = 1, 'follow' = 2)
)
ENGINE = MergeTree
ORDER BY (event_time, short_url);