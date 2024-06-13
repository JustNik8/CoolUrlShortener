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