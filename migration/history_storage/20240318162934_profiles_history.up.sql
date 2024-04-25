CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE profiles_history
(
    id                 UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    profile_id         UUID         NOT NULL,
    first_name         TEXT         NOT NULL,
    last_name          TEXT         NOT NULL,
    country            TEXT,
    city               TEXT,
    state              TEXT,
    gender             TEXT,
    occupation         TEXT,
    summary            TEXT,
    linkedin_id        TEXT         NOT NULL,
    is_blurred         BOOL,
    www_is_blurred     BOOL,
    profile_created_at TIMESTAMP(0) NOT NULL,
    created_at         TIMESTAMP(0) NOT NULL DEFAULT NOW()
);