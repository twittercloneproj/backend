CREATE TABLE tweet_by_user (
    posted_by text,
    id timeuuid,
    text text,
    PRIMARY KEY ((posted_by), id)
)
    WITH CLUSTERING ORDER BY (id DESC);

--     CREATE INDEX tweet_id on tweet_by_user(id);