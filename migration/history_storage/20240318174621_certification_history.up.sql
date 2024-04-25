CREATE TABLE certification_history
(
    id                       UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    certification_id         UUID         NOT NULL,
    profile_id               UUID         NOT NULL,
    name                     TEXT         NOT NULL,
    authority                TEXT         NOT NULL,
    license_number           TEXT,
    display_source           TEXT,
    url                      TEXT,
    authority_linkedin_url   TEXT,
    start_date               TIMESTAMP(0),
    end_date                 TIMESTAMP(0),
    certification_created_at TIMESTAMP(0) NOT NULL,
    created_at               TIMESTAMP(0) NOT NULL DEFAULT NOW()
);