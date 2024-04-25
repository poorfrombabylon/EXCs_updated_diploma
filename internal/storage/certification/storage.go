package certification

import (
	"context"
	certificationDomain "excs_updater/internal/domain/certification"
	profileDomain "excs_updater/internal/domain/profile"
	"excs_updater/internal/libdb"
	"github.com/Masterminds/squirrel"
)

const certificationTable = "certification"

type CertificationStorage interface {
	CreateCertification(ctx context.Context, certification certificationDomain.Certification) error
	CreateCertificationBunch(ctx context.Context, certificationList []certificationDomain.Certification) error
	GetCertsByProfileID(ctx context.Context, profileID profileDomain.ProfileID) ([]certificationDomain.Certification, error)
}

type certificationStorage struct {
	db libdb.DB
}

func NewCertificationStorage(db libdb.DB) CertificationStorage {
	return &certificationStorage{db: db}
}

func (c *certificationStorage) CreateCertification(ctx context.Context, certification certificationDomain.Certification) error {
	query := squirrel.Insert(certificationTable).
		Columns(
			"id",
			"profile_id",
			"name",
			"authority",
			"license_number",
			"display_source",
			"url",
			"authority_linkedin_url",
			"start_date",
			"end_date",
			"created_at",
		).
		Values(
			certification.GetCertificationID().String(),
			certification.GetProfileID().String(),
			certification.GetName(),
			certification.GetAuthority(),
			certification.GetLicenseNumber(),
			certification.GetDisplaySource(),
			certification.GetUrl(),
			certification.GetAuthorityLinkedinURL(),
			certification.GetStartDate(),
			certification.GetEndDate(),
			certification.GetCreatedAt(),
		).
		PlaceholderFormat(squirrel.Dollar)

	err := c.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (c *certificationStorage) CreateCertificationBunch(ctx context.Context, certificationList []certificationDomain.Certification) error {
	query := squirrel.Insert(certificationTable).
		Columns(
			"id",
			"profile_id",
			"name",
			"authority",
			"license_number",
			"display_source",
			"url",
			"authority_linkedin_url",
			"start_date",
			"end_date",
			"created_at",
		).
		PlaceholderFormat(squirrel.Dollar)

	for _, i := range certificationList {
		query = query.Values(
			i.GetCertificationID().String(),
			i.GetProfileID().String(),
			i.GetName(),
			i.GetAuthority(),
			i.GetLicenseNumber(),
			i.GetDisplaySource(),
			i.GetUrl(),
			i.GetAuthorityLinkedinURL(),
			i.GetStartDate(),
			i.GetEndDate(),
			i.GetCreatedAt(),
		)
	}

	err := c.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (c *certificationStorage) GetCertsByProfileID(ctx context.Context, profileID profileDomain.ProfileID) ([]certificationDomain.Certification, error) {
	query := squirrel.Select(
		"id",
		"profile_id",
		"name",
		"authority",
		"license_number",
		"display_source",
		"url",
		"authority_linkedin_url",
		"start_date",
		"end_date",
		"created_at",
		"updated_at",
	).
		From(certificationTable).
		Where(squirrel.Eq{"profile_id": profileID.String()}).
		PlaceholderFormat(squirrel.Dollar)

	var res []CertificationDTO

	err := c.db.Select(ctx, query, &res)
	if err != nil {
		return nil, err
	}

	return NewCertificationListFromDTO(res), nil
}
