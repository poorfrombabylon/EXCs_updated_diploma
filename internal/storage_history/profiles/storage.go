package profiles

import (
	"context"
	profileDomain "excs_updater/internal/domain/profile"
	"excs_updater/internal/libdb"
	"github.com/Masterminds/squirrel"
	"time"
)

const tableProfilesHistory = "profiles_history"

type ProfilesStorage interface {
	CreateHistoryProfile(ctx context.Context, profile profileDomain.Profile) error
}

type profileStorage struct {
	db libdb.DB
}

func NewHistoryProfileStorage(db libdb.DB) ProfilesStorage {
	return &profileStorage{
		db: db,
	}
}

func (p *profileStorage) CreateHistoryProfile(ctx context.Context, profile profileDomain.Profile) error {
	query := squirrel.Insert(tableProfilesHistory).
		Columns(
			"profile_id",
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
			"profile_created_at",
			"created_at",
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
			profile.GetCreatedAt(),
			time.Now().In(time.UTC),
		).
		PlaceholderFormat(squirrel.Dollar)

	err := p.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
