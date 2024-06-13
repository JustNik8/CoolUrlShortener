CREATE TABLE IF NOT EXISTS url_events
(
    long_url   String,
    short_url  String,
    event_time TIMESTAMP,
    event_type Enum8('create' = 1, 'follow' = 2)
)
    ENGINE = Kafka SETTINGS
        kafka_broker_list = 'kafka1:9092',
        kafka_topic_list = 'events',
        kafka_group_name = 'group1',
        kafka_format = 'JSONEachRow';

CREATE TABLE url_events_counter
(
    long_url     String,
    short_url    String,
    follow_count Int64,
    create_count Int64
) ENGINE = SummingMergeTree((follow_count, create_count))
      ORDER BY (long_url, short_url);

CREATE MATERIALIZED VIEW url_events_counter_mv TO url_events_counter AS
SELECT long_url,
       short_url,
       SUM(if(event_type == 'follow', 1, 0)) as follow_count,
       SUM(if(event_type == 'create', 1, 0)) as create_count
FROM url_events
GROUP BY long_url, short_url