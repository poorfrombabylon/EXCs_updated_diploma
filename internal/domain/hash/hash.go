package hash

import "time"

type Profile struct {
	FirstName           string
	LastName            string
	Country             *string
	City                *string
	State               *string
	Gender              *string
	Occupation          *string
	Summary             *string
	TelegramID          *string
	FacebookID          *string
	TwitterID           *string
	GithubID            *string
	AssociatedLinks     []string
	AssociatedUsernames []string
	PersonalEmails      []string
	PersonalNumbers     []string
}

type Experience struct {
	Position    *string
	CompanyName string
	Location    *string
	Description *string
	StartDate   *time.Time
	EndDate     *time.Time
}

type Education struct {
	FieldOfStudy             *string
	DegreeName               *string
	School                   string
	SchoolLinkedinProfileUrl *string
	Description              *string
	LogoUrl                  *string
	Grade                    *string
	ActivitiesAndSocieties   *string
	StartDate                *time.Time
	EndDate                  *time.Time
}

type Certification struct {
	Name                 string
	Authority            string
	LicenseNumber        *string
	DisplaySource        *string
	Url                  *string
	AuthorityLinkedinURL *string
	StartDate            *time.Time
	EndDate              *time.Time
}
