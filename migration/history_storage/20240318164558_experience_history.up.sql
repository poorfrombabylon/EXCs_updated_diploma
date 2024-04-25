CREATE TABLE experience_history
(
    id                    UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    experience_id         UUID         NOT NULL,
    profile_id            UUID         NOT NULL,
    position              TEXT,
    company_name          TEXT         NOT NULL DEFAULT '',
    location              TEXT,
    description           TEXT,
    start_date            TIMESTAMP(0),
    end_date              TIMESTAMP(0),
    experience_created_at TIMESTAMP(0) NOT NULL,
    created_at            TIMESTAMP(0) NOT NULL DEFAULT NOW()
);