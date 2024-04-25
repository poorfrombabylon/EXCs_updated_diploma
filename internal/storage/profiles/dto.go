package profiles

import (
	"encoding/json"
	"excs_updater/internal/domain"
	profileDomain "excs_updater/internal/domain/profile"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ProfileDTO struct {
	Id                  uuid.UUID `db:"id"`
	FirstName           string    `db:"first_name"`
	LastName            string    `db:"last_name"`
	Country             *string   `db:"country"`
	City                *string   `db:"city"`
	State               *string   `db:"state"`
	Gender              *string   `db:"gender"`
	Occupation          *string   `db:"occupation"`
	Summary             *string   `db:"summary"`
	LinkedinID          string    `db:"linkedin_id"`
	IsBlurred           *bool     `db:"is_blurred"`
	WwwIsBlurred        *bool     `db:"www_is_blurred"`
	GithubProfileID     *string   `db:"github_profile_id"`
	TwitterProfileID    *string   `db:"twitter_profile_id"`
	FacebookProfileID   *string   `db:"facebook_profile_id"`
	TelegramProfileID   *string   `db:"telegram_profile_id"`
	AssociatedLinks     []byte    `db:"associated_links"`
	AssociatedUsernames []byte    `db:"associated_usernames"`
	PersonalEmails      []byte    `db:"personal_emails"`
	PersonalNumbers     []byte    `db:"personal_numbers"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
	LastCheckedAt       time.Time `db:"last_checked_at"`
}

func NewProfileFromDTO(dto ProfileDTO) (profileDomain.Profile, error) {
	model := domain.NewFullModelFrom(dto.CreatedAt, dto.UpdatedAt, dto.LastCheckedAt)

	result := profileDomain.NewProfileWithID(
		profileDomain.ProfileID(dto.Id),
		dto.FirstName,
		dto.LastName,
		dto.Country,
		dto.City,
		dto.State,
		dto.Gender,
		dto.Occupation,
		dto.Summary,
		dto.LinkedinID,
		dto.IsBlurred,
		dto.WwwIsBlurred,
		model,
	)

	result.SetGithubProfileID(dto.GithubProfileID)
	result.SetTwitterProfileID(dto.TwitterProfileID)
	result.SetFacebookProfileID(dto.FacebookProfileID)
	result.SetTelegramProfileID(dto.FacebookProfileID)

	if dto.AssociatedLinks != nil {
		var associatedLinks []string

		err := json.Unmarshal(dto.AssociatedLinks, &associatedLinks)
		if err != nil {
			return profileDomain.Profile{}, fmt.Errorf("err while unmarshal associatedLinks for %s: %s", result.GetProfileID().String(), err.Error())
		}

		result.SetAssociatedLinks(associatedLinks)
	}

	if dto.AssociatedUsernames != nil {
		var associatedUsernames []string

		err := json.Unmarshal(dto.AssociatedUsernames, &associatedUsernames)
		if err != nil {
			return profileDomain.Profile{}, fmt.Errorf("err while unmarshal associatedUsernames for %s: %s", result.GetProfileID().String(), err.Error())
		}

		result.SetAssociatedUsernames(associatedUsernames)
	}

	if dto.PersonalEmails != nil {
		var personalEmails []string

		err := json.Unmarshal(dto.PersonalEmails, &personalEmails)
		if err != nil {
			return profileDomain.Profile{}, fmt.Errorf("err while unmarshal personalEmails for %s: %s", result.GetProfileID().String(), err.Error())
		}

		result.SetPersonalEmails(personalEmails)
	}

	if dto.PersonalNumbers != nil {
		var personalNumbers []string

		err := json.Unmarshal(dto.PersonalNumbers, &personalNumbers)
		if err != nil {
			return profileDomain.Profile{}, fmt.Errorf("err while unmarshal personalNumbers for %s: %s", result.GetProfileID().String(), err.Error())
		}

		result.SetPersonalNumbers(personalNumbers)
	}

	return result, nil
}

func NewProfilesFromDTOList(dto []ProfileDTO) ([]profileDomain.Profile, error) {
	result := make([]profileDomain.Profile, 0, len(dto))

	for _, i := range dto {
		p, err := NewProfileFromDTO(i)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}
