CREATE TABLE education_history
(
    id                          UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    education_id                UUID         NOT NULL,
    profile_id                  UUID         NOT NULL,
    field_of_study              TEXT,
    degree_name                 TEXT,
    school                      TEXT         NOT NULL default '',
    school_linkedin_profile_url TEXT,
    description                 TEXT,
    logo_url                    TEXT,
    grade                       TEXT,
    activities_and_societies    TEXT,
    start_date                  TIMESTAMP(0),
    end_date                    TIMESTAMP(0),
    education_created_at        TIMESTAMP(0) NOT NULL,
    created_at                  TIMESTAMP(0) NOT NULL DEFAULT NOW()
);