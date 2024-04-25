package storage

import (
	"context"
	"excs_updater/internal/config"
	"excs_updater/internal/libdb"
	"excs_updater/internal/storage/captcha"
	"excs_updater/internal/storage/certification"
	"excs_updater/internal/storage/education"
	"excs_updater/internal/storage/experience"
	"excs_updater/internal/storage/profiles"
	"excs_updater/internal/storage/redis/redis_to_update_receive"
	"excs_updater/internal/storage/redis/redis_to_update_send"
)

type Storages struct {
	EducationStorage     education.EducationStorage
	ExperienceStorage    experience.ExperienceStorage
	ProfilesStorage      profiles.ProfilesStorage
	CaptchaStorage       captcha.CaptchaStorage
	CertificationStorage certification.CertificationStorage
	RedisToUpdateSend    redis_to_update_send.RedisToUpdateSend
	RedisToUpdateReceive redis_to_update_receive.RedisToUpdateReceive
}

func NewStorageRegistry(_ context.Context, db libdb.DB, cfg config.Config) (*Storages, error) {
	redisToSend, err := redis_to_update_send.NewRedisToUpdateSend(cfg.RedisProfilesToUpdate)
	if err != nil {
		return nil, err
	}

	redisToReceive, err := redis_to_update_receive.NewRedisToUpdateReceive(cfg.RedisUpdatedProfiles)
	if err != nil {
		return nil, err
	}

	return &Storages{
		EducationStorage:     education.NewEducationStorage(db),
		ExperienceStorage:    experience.NewExperience(db),
		ProfilesStorage:      profiles.NewProfileStorage(db),
		CaptchaStorage:       captcha.NewCaptchaStorage(db),
		CertificationStorage: certification.NewCertificationStorage(db),
		RedisToUpdateSend:    redisToSend,
		RedisToUpdateReceive: redisToReceive,
	}, nil
}
