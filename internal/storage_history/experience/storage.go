package experience

import (
	"context"
	experienceDomain "excs_updater/internal/domain/experience"
	"excs_updater/internal/libdb"
	"github.com/Masterminds/squirrel"
	"time"
)

const tableExperienceHistory = "experience_history"

type HistoryExperienceStorage interface {
	CreateExperienceBunch(ctx context.Context, experienceList []experienceDomain.Experience) error
}

type storage struct {
	db libdb.DB
}

func NewHistoryExperienceStorage(db libdb.DB) HistoryExperienceStorage {
	return &storage{
		db: db,
	}
}

func (s *storage) CreateExperienceBunch(ctx context.Context, experienceList []experienceDomain.Experience) error {
	query := squirrel.Insert(tableExperienceHistory).
		Columns(
			"experience_id",
			"profile_id",
			"position",
			"company_name",
			"location",
			"description",
			"start_date",
			"end_date",
			"experience_created_at",
			"created_at",
		).
		PlaceholderFormat(squirrel.Dollar)

	for _, i := range experienceList {
		query = query.Values(
			i.GetExperienceID().String(),
			i.GetProfileID().String(),
			i.GetPosition(),
			i.GetCompanyName(),
			i.GetLocation(),
			i.GetDescription(),
			i.GetStartDate(),
			i.GetEndDate(),
			i.GetCreatedAt(),
			time.Now().In(time.UTC),
		)
	}

	err := s.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
