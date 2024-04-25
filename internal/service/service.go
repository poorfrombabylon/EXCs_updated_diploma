package service

import (
	"excs_updater/internal/service/certification"
	"excs_updater/internal/service/education"
	"excs_updater/internal/service/experience"
	"excs_updater/internal/service/profile"
	"excs_updater/internal/storage"
	history_storage "excs_updater/internal/storage_history"
)

type Services struct {
	ProfileService       profile.ProfileService
	ExperienceService    experience.ExperienceService
	EducationService     education.EducationService
	CertificationService certification.CertificationService
}

func NewServiceRegistry(storages *storage.Storages, historyStorages *history_storage.HistoryStorages) *Services {
	return &Services{
		ProfileService:       profile.NewProfileService(storages.ProfilesStorage, historyStorages.ProfileStorage, storages.RedisToUpdateSend, storages.RedisToUpdateReceive),
		ExperienceService:    experience.NewExperienceService(storages.ExperienceStorage, historyStorages.ExperienceStorage),
		EducationService:     education.NewEducationService(storages.EducationStorage, historyStorages.EducationStorage),
		CertificationService: certification.NewCertificationService(storages.CertificationStorage, historyStorages.CertificationStorage),
	}
}
