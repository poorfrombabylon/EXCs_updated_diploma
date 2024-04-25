package education

import (
	"context"
	educationDomain "excs_updater/internal/domain/education"
	profileDomain "excs_updater/internal/domain/profile"
	"excs_updater/internal/libdb"
	"github.com/Masterminds/squirrel"
	"time"
)

const EducationTable = "education"

type EducationStorage interface {
	CreateEducation(ctx context.Context, education educationDomain.Education) error
	CreateEducationBunch(ctx context.Context, educationList []educationDomain.Education) error
	GetEducationByProfileID(ctx context.Context, profileID profileDomain.ProfileID) ([]educationDomain.Education, error)
	UpdateEducationByID(ctx context.Context, id educationDomain.EducationID) error
	GetOldBunchData(ctx context.Context) ([]EducationDTO, error)
	DeleteEducationByID(ctx context.Context, id educationDomain.EducationID) error
}

type educationStorage struct {
	db libdb.DB
}

func NewEducationStorage(db libdb.DB) EducationStorage {
	return &educationStorage{
		db: db,
	}
}

func (e *educationStorage) CreateEducation(ctx context.Context, education educationDomain.Education) error {
	query := squirrel.Insert(EducationTable).
		Columns(
			"id",
			"profile_id",
			"education",
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
			"created_at",
		).
		Values(
			education.GetEducationID().String(),
			education.GetProfileID().String(),
			nil,
			education.GetFieldOfStudy(),
			education.GetDegreeName(),
			education.GetSchool(),
			education.GetSchoolLinkedinProfileUrl(),
			education.GetDescription(),
			education.GetLogoUrl(),
			education.GetGrade(),
			education.GetActivitiesAndSocieties(),
			education.GetStartDate(),
			education.GetEndDate(),
			education.GetCreatedAt(),
		).
		PlaceholderFormat(squirrel.Dollar)

	err := e.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (e *educationStorage) GetEducationByProfileID(ctx context.Context, profileID profileDomain.ProfileID) ([]educationDomain.Education, error) {
	query := squirrel.Select(
		"id",
		"profile_id",
		"education",
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
		"created_at",
	).
		From(EducationTable).
		Where(squirrel.Eq{"profile_id": profileID.String()}).
		PlaceholderFormat(squirrel.Dollar)

	var res []EducationDTO

	err := e.db.Select(ctx, query, &res)
	if err != nil {
		return nil, err
	}

	return NewEducationListFromDTO(res), nil
}

func (e *educationStorage) CreateEducationBunch(ctx context.Context, educationList []educationDomain.Education) error {
	query := squirrel.Insert(EducationTable).
		Columns(
			"id",
			"profile_id",
			"education",
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
			"created_at",
		).
		PlaceholderFormat(squirrel.Dollar)

	for _, i := range educationList {
		query = query.Values(
			i.GetEducationID().String(),
			i.GetProfileID().String(),
			nil,
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
		)
	}

	err := e.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (e *educationStorage) GetOldBunchData(ctx context.Context) ([]EducationDTO, error) {
	query := squirrel.Select(
		"id",
		"profile_id",
		"education",
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
		"created_at",
	).
		From(EducationTable).
		Where(squirrel.Eq{"school": ""}).
		Limit(1000).
		PlaceholderFormat(squirrel.Dollar)

	var res []EducationDTO

	err := e.db.Select(ctx, query, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *educationStorage) UpdateEducationByID(ctx context.Context, id educationDomain.EducationID) error {
	query := squirrel.Update(EducationTable).
		Set("updated_at", time.Now().In(time.UTC)).
		Where(squirrel.Eq{"id": id.String()}).
		PlaceholderFormat(squirrel.Dollar)

	return e.db.Update(ctx, query)
}

func (e *educationStorage) DeleteEducationByID(ctx context.Context, id educationDomain.EducationID) error {
	query := squirrel.Delete(EducationTable).
		Where(squirrel.Eq{"id": id.String()}).
		PlaceholderFormat(squirrel.Dollar)

	return e.db.Delete(ctx, query)
}
