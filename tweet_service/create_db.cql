CREATE TABLE tweet_by_user (
    posted_by text,
    id timeuuid,
    text text,
    retweet boolean,
    original_posted_by text,
    PRIMARY KEY ((posted_by), id)
)
    WITH CLUSTERING ORDER BY (id DESC);

CREATE INDEX tweet_by_user_id ON tweet_by_user(id);

CREATE TABLE likes (
    id timeuuid,
    username text,
    PRIMARY KEY ((id), username)
);

CREATE TABLE feed_by_user (
    username text,
    id timeuuid,
    posted_by text,
    text text,
    retweet boolean,
    original_posted_by text,
    PRIMARY KEY ((username), id)
);
    WITH CLUSTERING ORDER BY (id DESC);