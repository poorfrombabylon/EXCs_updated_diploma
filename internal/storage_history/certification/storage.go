package certification

import (
	"context"
	certificationDomain "excs_updater/internal/domain/certification"
	"excs_updater/internal/libdb"
	"github.com/Masterminds/squirrel"
	"time"
)

const tableHistoryCertification = "certification_history"

type CertificationStorage interface {
	CreateCertificationBunch(ctx context.Context, certificationList []certificationDomain.Certification) error
}

type storage struct {
	db libdb.DB
}

func NewHistoryCertificationStorage(db libdb.DB) CertificationStorage {
	return &storage{
		db: db,
	}
}

func (s *storage) CreateCertificationBunch(ctx context.Context, certificationList []certificationDomain.Certification) error {
	query := squirrel.Insert(tableHistoryCertification).
		Columns(
			"certification_id",
			"profile_id",
			"name",
			"authority",
			"license_number",
			"display_source",
			"url",
			"authority_linkedin_url",
			"start_date",
			"end_date",
			"certification_created_at",
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
			time.Now().In(time.UTC),
		)
	}

	err := s.db.Insert(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
