package education

import (
	"context"
	educationDomain "excs_updater/internal/domain/education"
	"excs_updater/internal/libdb"
	"github.com/Masterminds/squirrel"
	"time"
)

const tableHistoryEducation = "education_history"

type EducationStorage interface {
	CreateEducationBunch(ctx context.Context, educationList []educationDomain.Education) error
}

type storage struct {
	db libdb.DB
}

func NewHistoryEducationStorage(db libdb.DB) EducationStorage {
	return &storage{
		db: db,
	}
}

func (s *storage) CreateEducationBunch(ctx context.Context, educationList []educationDomain.Education) error {
	query := squirrel.Insert(tableHistoryEducation).
		Columns(
			"education_id",
			"profile_id",
			"field_of_study",
			"degree_name",
			"school",
			"school_linkedin_profile_url",
			"description",
			"logo_url",
			"grade",
			"activities_and_societies",
			"start_date",
			"end_date",
			"education_created_at",
			"created_at",
		).
		PlaceholderFormat(squirrel.Dollar)

	for _, i := range educationList {
		query = query.Values(
			i.GetEducationID().String(),
			i.GetProfileID().String(),
			i.GetFieldOfStudy(),
			i.GetDegreeName(),
			i.GetSchool(),
			i.GetSchoolLinkedinProfileUrl(),
			i.GetDescription(),
			i.GetLogoUrl(),
			i.GetGrade(),
			i.GetActivitiesAndSocieties(),
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
