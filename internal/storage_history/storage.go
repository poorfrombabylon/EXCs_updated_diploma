package storage_history

import (
	"excs_updater/internal/libdb"
	history_certification "excs_updater/internal/storage_history/certification"
	history_education "excs_updater/internal/storage_history/education"
	history_experience "excs_updater/internal/storage_history/experience"
	history_profiles "excs_updater/internal/storage_history/profiles"
)

type HistoryStorages struct {
	ProfileStorage       history_profiles.ProfilesStorage
	ExperienceStorage    history_experience.HistoryExperienceStorage
	EducationStorage     history_education.EducationStorage
	CertificationStorage history_certification.CertificationStorage
}

func NewHistoryStorageRegistry(db libdb.DB) (*HistoryStorages, error) {
	return &HistoryStorages{
		ProfileStorage:       history_profiles.NewHistoryProfileStorage(db),
		ExperienceStorage:    history_experience.NewHistoryExperienceStorage(db),
		EducationStorage:     history_education.NewHistoryEducationStorage(db),
		CertificationStorage: history_certification.NewHistoryCertificationStorage(db),
	}, nil
}
