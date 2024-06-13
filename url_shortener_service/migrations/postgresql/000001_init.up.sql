CREATE TABLE IF NOT EXISTS "url_data"
(
    "id"         BIGINT      NOT NULL,
    "short_url"  VARCHAR(10) NOT NULL,
    "long_url"   TEXT        NOT NULL,
    "created_at" DATE        NOT NULL
);
ALTER TABLE
    "url_data"
    ADD PRIMARY KEY ("id");