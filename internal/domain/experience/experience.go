package experience

import (
	"excs_updater/internal/domain"
	profileDomain "excs_updater/internal/domain/profile"
	"github.com/google/uuid"
	"time"
)

type ExperienceID uuid.UUID

func (e ExperienceID) String() string {
	return uuid.UUID(e).String()
}

type Experience struct {
	id                        ExperienceID
	profileID                 profileDomain.ProfileID
	position                  *string
	companyName               string
	location                  *string
	description               *string
	companyLinkedinProfileUrl *string
	startDate                 *time.Time
	endDate                   *time.Time
	domain.Model
}

func NewExperience(
	profileID profileDomain.ProfileID,
	position *string,
	companyName string,
	location *string,
	description *string,
	companyLinkedinProfileUrl *string,
	startDate *time.Time,
	endDate *time.Time,
) Experience {
	return Experience{
		id:                        ExperienceID(uuid.New()),
		profileID:                 profileID,
		position:                  position,
		companyName:               companyName,
		location:                  location,
		description:               description,
		companyLinkedinProfileUrl: companyLinkedinProfileUrl,
		startDate:                 startDate,
		endDate:                   endDate,
		Model:                     domain.NewModel(),
	}
}

func NewExperienceWithID(
	id ExperienceID,
	profileID profileDomain.ProfileID,
	position *string,
	companyName string,
	location *string,
	description *string,
	companyLinkedinProfileUrl *string,
	startDate *time.Time,
	endDate *time.Time,
	model domain.Model,
) Experience {
	return Experience{
		id:                        id,
		profileID:                 profileID,
		position:                  position,
		companyName:               companyName,
		location:                  location,
		description:               description,
		companyLinkedinProfileUrl: companyLinkedinProfileUrl,
		startDate:                 startDate,
		endDate:                   endDate,
		Model:                     model,
	}
}

func (e Experience) GetExperienceID() ExperienceID {
	return e.id
}

func (e Experience) GetProfileID() profileDomain.ProfileID {
	return e.profileID
}

func (e Experience) GetPosition() *string {
	return e.position
}

func (e Experience) GetCompanyName() string {
	return e.companyName
}

func (e Experience) GetLocation() *string {
	return e.location
}

func (e Experience) GetDescription() *string {
	return e.description
}

func (e Experience) GetStartDate() *time.Time {
	return e.startDate
}

func (e Experience) GetEndDate() *time.Time {
	return e.endDate
}

func (e Experience) GetCompanyLinkedinProfileUrl() *string {
	return e.companyLinkedinProfileUrl
}
