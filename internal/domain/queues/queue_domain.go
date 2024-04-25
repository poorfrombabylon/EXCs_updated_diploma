package queues

type ProfileToUpdate struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	ProfileId   string `json:"profile_id"`
	ProfileLink string `json:"profile_link"`
}

type FullProfileInfo struct {
	CaptchaMeet      int             `json:"captcha_meet"`
	PublicIdentifier string          `json:"public_identifier"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	Country          *string         `json:"country,omitempty"`
	City             *string         `json:"city,omitempty"`
	State            *string         `json:"state,omitempty"`
	Occupation       *string         `json:"occupation,omitempty"`
	Summary          *string         `json:"summary,omitempty"`
	Gender           *string         `json:"gender,omitempty"`
	IsBlurred        *bool           `json:"is_blurred,omitempty"`
	WwwIsBlurred     *bool           `json:"www_is_blurred,omitempty"`
	Experiences      []Experience    `json:"experiences"`
	Education        []Education     `json:"education"`
	Certifications   []Certification `json:"certifications"`
	Extra            *Extra          `json:"extra,omitempty"`
	PersonalEmails   []string        `json:"personal_emails"`
	PersonalNumbers  []string        `json:"personal_numbers"`
}

type Date struct {
	Day   *int `json:"day,omitempty"`
	Month *int `json:"month,omitempty"`
	Year  *int `json:"year,omitempty"`
}

type Experience struct {
	Position                  *string `json:"position,omitempty"`
	CompanyName               string  `json:"company_name"`
	StartDate                 *Date   `json:"start_date,omitempty"`
	EndDate                   *Date   `json:"end_date,omitempty"`
	Location                  *string `json:"location,omitempty"`
	Description               *string `json:"description,omitempty"`
	CompanyLinkedinProfileUrl *string `json:"company_linkedin_profile_url,omitempty"`
}

type Education struct {
	StartDate                *Date   `json:"start_date,omitempty"`
	EndDate                  *Date   `json:"end_date,omitempty"`
	FieldOfStudy             *string `json:"field_of_study,omitempty"`
	DegreeName               *string `json:"degree_name,omitempty"`
	School                   string  `json:"school"`
	SchoolLinkedinProfileUrl *string `json:"school_linkedin_profile_url,omitempty"`
	Description              *string `json:"description,omitempty"`
	LogoUrl                  *string `json:"logo_url,omitempty"`
	Grade                    *string `json:"grade,omitempty"`
	ActivitiesAndSocieties   *string `json:"activities_and_societies,omitempty"`
}

type Certification struct {
	StartDate            *Date   `json:"start_date,omitempty"`
	EndDate              *Date   `json:"end_date,omitempty"`
	Name                 string  `json:"name"`
	LicenseNumber        *string `json:"license_number,omitempty"`
	DisplaySource        *string `json:"display_source,omitempty"`
	Authority            string  `json:"authority"`
	Url                  *string `json:"url,omitempty"`
	AuthorityLinkedinURL *string `json:"authority_linkedin_url,omitempty"`
}

type Extra struct {
	GithubProfileID     *string  `json:"github_profile_id,omitempty"`
	TwitterProfileID    *string  `json:"twitter_profile_id,omitempty"`
	FacebookProfileID   *string  `json:"facebook_profile_id,omitempty"`
	TelegramProfileID   *string  `json:"telegram_profile_id,omitempty"`
	AssociatedLinks     []string `json:"associated_links,omitempty"`
	AssociatedUsernames []string `json:"associated_usernames,omitempty"`
}
