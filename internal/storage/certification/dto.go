package certification

import (
	certificationDomain "excs_updater/internal/domain/certification"
	profileDomain "excs_updater/internal/domain/profile"
	"github.com/google/uuid"
	"time"
)

type CertificationDTO struct {
	ID                   uuid.UUID  `db:"id"`
	ProfileID            uuid.UUID  `db:"profile_id"`
	Name                 string     `db:"name"`
	Authority            string     `db:"authority"`
	LicenseNumber        *string    `db:"license_number"`
	DisplaySource        *string    `db:"display_source"`
	Url                  *string    `db:"url"`
	AuthorityLinkedinURL *string    `db:"authority_linkedin_url"`
	StartDate            *time.Time `db:"start_date"`
	EndDate              *time.Time `db:"end_date"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
}

func NewCertificationFromDTO(dto CertificationDTO) certificationDomain.Certification {
	return certificationDomain.NewCertificationWithID(
		certificationDomain.CertificationID(dto.ID),
		profileDomain.ProfileID(dto.ProfileID),
		dto.Name,
		dto.Authority,
		dto.LicenseNumber,
		dto.DisplaySource,
		dto.Url,
		dto.AuthorityLinkedinURL,
		dto.StartDate,
		dto.EndDate,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
}

func NewCertificationListFromDTO(dtos []CertificationDTO) []certificationDomain.Certification {
	result := make([]certificationDomain.Certification, 0, len(dtos))

	for _, i := range dtos {
		result = append(result, NewCertificationFromDTO(i))
	}

	return result
}
