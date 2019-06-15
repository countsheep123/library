CREATE TABLE IF NOT EXISTS users (
    id         VARCHAR(20)  PRIMARY KEY,
    created_at BIGINT       NOT NULL,
    updated_at BIGINT       NOT NULL,
    email      VARCHAR(320) NOT NULL UNIQUE,
    name       VARCHAR(255) NOT NULL,
    company    VARCHAR(255),
    is_admin   BOOLEAN      NOT NULL,

    CONSTRAINT timestamp_check CHECK(updated_at >= created_at)
);

CREATE TABLE IF NOT EXISTS books (
    id         VARCHAR(20)    PRIMARY KEY,
    created_at BIGINT         NOT NULL,
    updated_at BIGINT         NOT NULL,
    title      TEXT           NOT NULL,
    isbn       VARCHAR(17)    UNIQUE,
    publisher  VARCHAR(255),
    pubdate    BIGINT,
    authors    VARCHAR(255)[],
    cover_url  TEXT,

    CONSTRAINT timestamp_check CHECK(updated_at >= created_at)
);

CREATE TABLE IF NOT EXISTS book_labels (
    id         VARCHAR(20) PRIMARY KEY,
    created_at BIGINT      NOT NULL,
    updated_at BIGINT      NOT NULL,
    book_id    VARCHAR(20) NOT NULL,
    user_id    VARCHAR(20) NOT NULL,
    label      TEXT        NOT NULL,

    CONSTRAINT timestamp_check CHECK(updated_at >= created_at),
    UNIQUE (book_id, label),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS book_recommenders (
    id         VARCHAR(20) PRIMARY KEY,
    created_at BIGINT      NOT NULL,
    updated_at BIGINT      NOT NULL,
    book_id    VARCHAR(20) NOT NULL,
    user_id    VARCHAR(20) NOT NULL,

    CONSTRAINT timestamp_check CHECK(updated_at >= created_at),
    UNIQUE (book_id, user_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS marks (
    id         VARCHAR(20)  PRIMARY KEY,
    created_at BIGINT       NOT NULL,
    updated_at BIGINT       NOT NULL,
    user_id    VARCHAR(20)  NOT NULL,
    name       VARCHAR(255) NOT NULL,
    url        TEXT,

    CONSTRAINT timestamp_check CHECK(updated_at >= created_at)
);

CREATE TABLE IF NOT EXISTS locations (
    id         VARCHAR(20)  PRIMARY KEY,
    created_at BIGINT       NOT NULL,
    updated_at BIGINT       NOT NULL,
    name       VARCHAR(255) NOT NULL,

    CONSTRAINT timestamp_check CHECK(updated_at >= created_at)
);

CREATE TABLE IF NOT EXISTS stocks (
    id           VARCHAR(20) PRIMARY KEY,
    created_at   BIGINT      NOT NULL,
    updated_at   BIGINT      NOT NULL,
    is_available BOOLEAN     NOT NULL,
    book_id      VARCHAR(20) NOT NULL,
    user_id      VARCHAR(20) NOT NULL,
    mark_id      VARCHAR(20) NOT NULL,
    location_id  VARCHAR(20) NOT NULL,

    CONSTRAINT timestamp_check CHECK(updated_at >= created_at),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE,
    FOREIGN KEY (mark_id) REFERENCES marks (id) ON DELETE CASCADE,
    FOREIGN KEY (location_id) REFERENCES locations (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS records (
    id          VARCHAR(20) PRIMARY KEY,
    lent_at     BIGINT      NOT NULL,
    returned_at BIGINT,
    user_id     VARCHAR(20) NOT NULL,
    stock_id    VARCHAR(20) NOT NULL,

    CONSTRAINT timestamp_check CHECK(returned_at > lent_at),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (stock_id) REFERENCES stocks (id) ON DELETE CASCADE
);
