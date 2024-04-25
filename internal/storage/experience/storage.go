package experience

import (
	"context"
	experienceDomain "excs_updater/internal/domain/experience"
	profileDomain "excs_updater/internal/domain/profile"
	"excs_updater/internal/libdb"
	"excs_updater/internal/storage/profiles"
	"github.com/Masterminds/squirrel"
	"time"
)

const experienceTable = "experience"

type ExperienceStorage interface {
	CreateExperience(ctx context.Context, experience experienceDomain.Experience) error
	CreateExperienceBunch(ctx context.Context, experience []experienceDomain.Experience) error
	GetExperienceByProfileID(ctx context.Context, profileID profileDomain.ProfileID) ([]experienceDomain.Experience, error)
	GetExperienceByLinkedinID(ctx context.Context, linkedinID string) ([]experienceDomain.Experience, error)
	DeleteExperienceByID(ctx context.Context, id experienceDomain.ExperienceID) error
	UpdateExperienceByID(ctx context.Context, id experienceDomain.ExperienceID) error
	GetOldBunchData(ctx context.Context) ([]ExperienceDTO, error)
}

type experienceStorage struct {
	db libdb.DB
}

func NewExperience(db libdb.DB) ExperienceStorage {
	return &experienceStorage{
		db: db,
	}
}

func (e *experienceStorage) CreateExperience(ctx context.Context, experience experienceDomain.Experience) error {
	query := squirrel.Insert(experienceTable).
		Columns(
			"id",
			"profile_id",
			"experience",
			"position",
			"company_name",
			"location",
			"description",
			"start_date",
			"end_date",
			"created_at",
			"updated_at",
		).
		Values(
			experience.GetExperienceID().String(),
			experience.GetProfileID().String(),
			nil,
			experience.GetPosition(),
			experience.GetCompanyName(),
			experience.GetLocation(),
			experience.GetDescription(),
			experience.GetStartDate(),
			experience.GetEndDate(),
			experience.GetCreatedAt(),
			experience.GetUpdatedAt(),
		).
		PlaceholderFormat(squirrel.Dollar)

	err := e.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (e *experienceStorage) CreateExperienceBunch(ctx context.Context, experience []experienceDomain.Experience) error {
	query := squirrel.Insert(experienceTable).
		Columns(
			"id",
			"profile_id",
			"experience",
			"position",
			"company_name",
			"location",
			"description",
			"start_date",
			"end_date",
			"created_at",
			"updated_at",
		).
		PlaceholderFormat(squirrel.Dollar)

	for _, i := range experience {
		query = query.Values(
			i.GetExperienceID().String(),
			i.GetProfileID().String(),
			nil,
			i.GetPosition(),
			i.GetCompanyName(),
			i.GetLocation(),
			i.GetDescription(),
			i.GetStartDate(),
			i.GetEndDate(),
			i.GetCreatedAt(),
			i.GetUpdatedAt(),
		)
	}

	err := e.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (e *experienceStorage) GetExperienceByLinkedinID(ctx context.Context, linkedinID string) ([]experienceDomain.Experience, error) {
	query := squirrel.Select(
		"e.id as id",
		"e.profile_id as profile_id",
		"e.experience as experience",
		"e.position as position",
		"e.company_name as company_name",
		"e.location as location",
		"e.description as description",
		"e.start_date as start_date",
		"e.end_date as end_date",
		"e.created_at as created_at",
		"e.updated_at as updated_at",
	).
		From(experienceTable + " as e").
		Join(profiles.ProfilesTable + " as p ON p.id = e.profile_id").
		Where(squirrel.Eq{"p.linkedin_id": linkedinID}).
		PlaceholderFormat(squirrel.Dollar)

	var result []ExperienceDTO

	err := e.db.Select(ctx, query, &result)
	if err != nil {
		return nil, err
	}

	return NewExperienceListFromDTO(result), nil
}

func (e *experienceStorage) GetOldBunchData(ctx context.Context) ([]ExperienceDTO, error) {
	query := squirrel.Select(
		"id",
		"profile_id",
		"experience",
		"position",
		"company_name",
		"location",
		"description",
		"start_date",
		"end_date",
		"created_at",
		"updated_at",
	).
		From(experienceTable).
		Where(squirrel.Eq{"company_name": ""}).
		Limit(1000).
		PlaceholderFormat(squirrel.Dollar)

	var result []ExperienceDTO

	err := e.db.Select(ctx, query, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e *experienceStorage) GetExperienceByProfileID(
	ctx context.Context,
	profileID profileDomain.ProfileID,
) ([]experienceDomain.Experience, error) {
	query := squirrel.Select(
		"id",
		"profile_id",
		"experience",
		"position",
		"company_name",
		"location",
		"description",
		"start_date",
		"end_date",
		"created_at",
		"updated_at",
	).
		From(experienceTable).
		Where(squirrel.Eq{"profile_id": profileID.String()}).
		PlaceholderFormat(squirrel.Dollar)

	var result []ExperienceDTO

	err := e.db.Select(ctx, query, &result)
	if err != nil {
		return nil, err
	}

	return NewExperienceListFromDTO(result), nil
}

func (e *experienceStorage) UpdateExperienceByID(ctx context.Context, id experienceDomain.ExperienceID) error {
	query := squirrel.Update(experienceTable).
		Set("updated_at", time.Now().In(time.UTC)).
		Where(squirrel.Eq{"id": id.String()}).
		PlaceholderFormat(squirrel.Dollar)

	return e.db.Update(ctx, query)
}

func (e *experienceStorage) DeleteExperienceByID(ctx context.Context, id experienceDomain.ExperienceID) error {
	query := squirrel.Delete(experienceTable).
		Where(squirrel.Eq{"id": id.String()}).
		PlaceholderFormat(squirrel.Dollar)

	return e.db.Delete(ctx, query)
}
