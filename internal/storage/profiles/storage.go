package profiles

import (
	"context"
	"encoding/json"
	profileDomain "excs_updater/internal/domain/profile"
	"excs_updater/internal/libdb"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const ProfilesTable = "profiles"

type ProfilesStorage interface {
	CreateProfile(ctx context.Context, profile profileDomain.Profile) error
	GetBunchProfilesForUpdate(ctx context.Context, page uint64) ([]profileDomain.Profile, error)
	GetProfileByLinkedInID(ctx context.Context, linkedinID string) (profileDomain.Profile, error)
	DeleteProfileByID(ctx context.Context, profileID profileDomain.ProfileID) error
	UpdateProfileFlags(ctx context.Context, profile profileDomain.Profile) error
	UpdateFullProfile(ctx context.Context, profile profileDomain.Profile) error
	GetProfileIDByLinkedinID(ctx context.Context, linkedinID string) (profileDomain.ProfileID, error)
}

type profileStorage struct {
	db libdb.DB
}

func NewProfileStorage(db libdb.DB) ProfilesStorage {
	return &profileStorage{
		db: db,
	}
}

func (p *profileStorage) CreateProfile(ctx context.Context, profile profileDomain.Profile) error {
	jsonAssociatedLinks, err := json.Marshal(profile.GetAssociatedLinks())
	if err != nil {
		return fmt.Errorf("err while marshal assoc links for %s: %s", profile.GetProfileID().String(), err.Error())
	}

	jsonAssociatedUsernames, err := json.Marshal(profile.GetAssociatedUsernames())
	if err != nil {
		return fmt.Errorf("err while marshal assoc usernames for %s: %s", profile.GetProfileID().String(), err.Error())
	}

	jsonPersonalEmails, err := json.Marshal(profile.GetPersonalEmails())
	if err != nil {
		return fmt.Errorf("err while marshal personal emails for %s: %s", profile.GetProfileID().String(), err.Error())
	}

	jsonPersonalNumbers, err := json.Marshal(profile.GetPersonalNumbers())
	if err != nil {
		return fmt.Errorf("err while marshal personal numbers for %s: %s", profile.GetProfileID().String(), err.Error())
	}

	query := squirrel.Insert(ProfilesTable).
		Columns(
			"id",
			"first_name",
			"last_name",
			"country",
			"city",
			"state",
			"gender",
			"occupation",
			"summary",
			"linkedin_id",
			"is_blurred",
			"www_is_blurred",
			"github_profile_id",
			"twitter_profile_id",
			"facebook_profile_id",
			"telegram_profile_id",
			"associated_links",
			"associated_usernames",
			"personal_emails",
			"personal_numbers",
			"created_at",
			"updated_at",
			"last_checked_at",
		).
		Values(
			profile.GetProfileID().String(),
			profile.GetFirstName(),
			profile.GetLastName(),
			profile.GetCountry(),
			profile.GetCity(),
			profile.GetState(),
			profile.GetGender(),
			profile.GetOccupation(),
			profile.GetSummary(),
			profile.GetLinkedInID(),
			profile.GetIsBlurred(),
			profile.GetWwwIsBlurred(),
			profile.GetGithubProfileID(),
			profile.GetTwitterProfileID(),
			profile.GetFacebookProfileID(),
			profile.GetTelegramProfileID(),
			jsonAssociatedLinks,
			jsonAssociatedUsernames,
			jsonPersonalEmails,
			jsonPersonalNumbers,
			profile.GetCreatedAt(),
			profile.GetUpdatedAt(),
			profile.GetLastCheckedAt(),
		).
		PlaceholderFormat(squirrel.Dollar)

	err = p.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (p *profileStorage) GetBunchProfilesForUpdate(ctx context.Context, page uint64) ([]profileDomain.Profile, error) {
	query := squirrel.Select(
		"p.id as id",
		"p.first_name as first_name",
		"p.last_name as last_name",
		"p.country as country",
		"p.city as city",
		"p.state as state",
		"p.gender as gender",
		"p.occupation as occupation",
		"p.summary as summary",
		"p.linkedin_id as linkedin_id",
		"p.is_blurred as is_blurred",
		"p.www_is_blurred as www_is_blurred",
		"p.github_profile_id as github_profile_id",
		"p.twitter_profile_id as twitter_profile_id",
		"p.facebook_profile_id as facebook_profile_id",
		"p.telegram_profile_id as telegram_profile_id",
		"p.associated_links as associated_links",
		"p.associated_usernames as associated_usernames",
		"p.personal_emails as personal_emails",
		"p.personal_numbers as personal_numbers",
		"p.created_at as created_at",
		"p.updated_at as updated_at",
		"p.last_checked_at as last_checked_at",
	).
		From(ProfilesTable + " as p").
		Where(squirrel.Expr("p.last_checked_at < current_date - interval '30' day")).
		Limit(500).
		Offset(page).
		PlaceholderFormat(squirrel.Dollar)

	var result []ProfileDTO

	err := p.db.Select(ctx, query, &result)
	if err != nil {
		return nil, err
	}

	return NewProfilesFromDTOList(result)
}

func (p *profileStorage) GetProfileByLinkedInID(ctx context.Context, linkedinID string) (profileDomain.Profile, error) {
	query := squirrel.Select(
		"id",
		"first_name",
		"last_name",
		"country",
		"city",
		"state",
		"gender",
		"occupation",
		"summary",
		"linkedin_id",
		"is_blurred",
		"www_is_blurred",
		"github_profile_id",
		"twitter_profile_id",
		"facebook_profile_id",
		"telegram_profile_id",
		"associated_links",
		"associated_usernames",
		"personal_emails",
		"personal_numbers",
		"created_at",
		"updated_at",
		"last_checked_at",
	).
		From(ProfilesTable).
		Where(squirrel.Eq{"linkedin_id": linkedinID}).
		PlaceholderFormat(squirrel.Dollar)

	var res ProfileDTO

	err := p.db.Get(ctx, query, &res)
	if err != nil {
		return profileDomain.Profile{}, err
	}

	return NewProfileFromDTO(res)
}

func (p *profileStorage) GetProfileIDByLinkedinID(ctx context.Context, linkedinID string) (profileDomain.ProfileID, error) {
	query := squirrel.Select(
		"id",
	).
		From(ProfilesTable).
		Where(squirrel.Eq{"linkedin_id": linkedinID}).
		PlaceholderFormat(squirrel.Dollar)

	var result uuid.UUID

	err := p.db.Get(ctx, query, &result)
	if err != nil {
		return profileDomain.ProfileID{}, err
	}

	return profileDomain.ProfileID(result), nil
}

func (p *profileStorage) UpdateProfileFlags(ctx context.Context, profile profileDomain.Profile) error {
	query := squirrel.Update(ProfilesTable).
		Set("is_blurred", profile.GetIsBlurred()).
		Set("www_is_blurred", profile.GetWwwIsBlurred()).
		Set("updated_at", profile.GetUpdatedAt()).
		Set("last_checked_at", profile.GetLastCheckedAt()).
		Where(squirrel.Eq{"id": profile.GetProfileID().String()}).
		PlaceholderFormat(squirrel.Dollar)

	err := p.db.Update(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (p *profileStorage) UpdateFullProfile(ctx context.Context, profile profileDomain.Profile) error {
	jsonAssociatedLinks, err := json.Marshal(profile.GetAssociatedLinks())
	if err != nil {
		return fmt.Errorf("err while marshal assoc links for %s: %s", profile.GetProfileID().String(), err.Error())
	}

	jsonAssociatedUsernames, err := json.Marshal(profile.GetAssociatedUsernames())
	if err != nil {
		return fmt.Errorf("err while marshal assoc usernames for %s: %s", profile.GetProfileID().String(), err.Error())
	}

	jsonPersonalEmails, err := json.Marshal(profile.GetPersonalEmails())
	if err != nil {
		return fmt.Errorf("err while marshal personal emails for %s: %s", profile.GetProfileID().String(), err.Error())
	}

	jsonPersonalNumbers, err := json.Marshal(profile.GetPersonalNumbers())
	if err != nil {
		return fmt.Errorf("err while marshal personal numbers for %s: %s", profile.GetProfileID().String(), err.Error())
	}

	query := squirrel.Update(ProfilesTable).
		Set("first_name", profile.GetFirstName()).
		Set("last_name", profile.GetLastName()).
		Set("country", profile.GetCountry()).
		Set("city", profile.GetCity()).
		Set("state", profile.GetState()).
		Set("gender", profile.GetGender()).
		Set("occupation", profile.GetOccupation()).
		Set("summary", profile.GetSummary()).
		Set("is_blurred", profile.GetIsBlurred()).
		Set("www_is_blurred", profile.GetWwwIsBlurred()).
		Set("github_profile_id", profile.GetGithubProfileID()).
		Set("twitter_profile_id", profile.GetTwitterProfileID()).
		Set("facebook_profile_id", profile.GetFacebookProfileID()).
		Set("telegram_profile_id", profile.GetTelegramProfileID()).
		Set("associated_usernames", jsonAssociatedUsernames).
		Set("associated_links", jsonAssociatedLinks).
		Set("personal_numbers", jsonPersonalNumbers).
		Set("personal_emails", jsonPersonalEmails).
		Set("updated_at", profile.GetUpdatedAt()).
		Set("last_checked_at", profile.GetLastCheckedAt()).
		Where(squirrel.Eq{"id": profile.GetProfileID().String()}).
		PlaceholderFormat(squirrel.Dollar)

	err = p.db.Update(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (p *profileStorage) DeleteProfileByID(ctx context.Context, profileID profileDomain.ProfileID) error {
	query := squirrel.Delete(ProfilesTable).
		Where(squirrel.Eq{"id": profileID.String()}).
		PlaceholderFormat(squirrel.Dollar)

	return p.db.Delete(ctx, query)
}
