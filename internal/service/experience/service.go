package experience

import (
	"context"
	experienceDomain "excs_updater/internal/domain/experience"
	hashDomain "excs_updater/internal/domain/hash"
	profileDomain "excs_updater/internal/domain/profile"
	"fmt"
	"github.com/mitchellh/hashstructure/v2"
)

type ExperienceService interface {
	GetExperienceByProfileID(context.Context, profileDomain.ProfileID) ([]experienceDomain.Experience, error)
	GetExperienceByLinkedinID(context.Context, string) ([]experienceDomain.Experience, error)
	ProcessUpdatedExperience(context.Context, profileDomain.ProfileID, []experienceDomain.Experience) error
}

type ExperienceStorage interface {
	GetExperienceByProfileID(ctx context.Context, profileID profileDomain.ProfileID) ([]experienceDomain.Experience, error)
	GetExperienceByLinkedinID(context.Context, string) ([]experienceDomain.Experience, error)
	CreateExperienceBunch(ctx context.Context, experience []experienceDomain.Experience) error
}

type HistoryExperienceStorage interface {
	CreateExperienceBunch(ctx context.Context, experienceList []experienceDomain.Experience) error
}

type service struct {
	experienceStorage        ExperienceStorage
	historyExperienceStorage HistoryExperienceStorage
}

func NewExperienceService(
	experienceStorage ExperienceStorage,
	historyExperienceStorage HistoryExperienceStorage,
) ExperienceService {
	return &service{
		experienceStorage:        experienceStorage,
		historyExperienceStorage: historyExperienceStorage,
	}
}

func (s *service) GetExperienceByProfileID(
	ctx context.Context,
	profileID profileDomain.ProfileID,
) ([]experienceDomain.Experience, error) {
	return s.experienceStorage.GetExperienceByProfileID(ctx, profileID)
}

func (s *service) GetExperienceByLinkedinID(
	ctx context.Context,
	linkedinID string,
) ([]experienceDomain.Experience, error) {
	return s.experienceStorage.GetExperienceByLinkedinID(ctx, linkedinID)
}

func (s *service) ProcessUpdatedExperience(
	ctx context.Context,
	profileID profileDomain.ProfileID,
	updatedExperience []experienceDomain.Experience,
) error {
	initialExperience, err := s.experienceStorage.GetExperienceByProfileID(ctx, profileID)
	if err != nil {
		return err
	}

	lenInitialExp := len(initialExperience)
	lenUpdatedExp := len(updatedExperience)

	if lenInitialExp >= lenUpdatedExp {
		return nil
	}

	initialExpHashes := map[uint64]experienceDomain.Experience{}
	updatedExpHashes := map[uint64]experienceDomain.Experience{}

	for _, i := range initialExperience {
		hashStruct := hashDomain.Experience{
			i.GetPosition(),
			i.GetCompanyName(),
			i.GetLocation(),
			i.GetDescription(),
			i.GetStartDate(),
			i.GetEndDate(),
		}

		hash, err := hashstructure.Hash(hashStruct, hashstructure.FormatV2, nil)
		if err != nil {
			return fmt.Errorf("err while get hash for initial experience %s: %s", i.GetExperienceID().String(), err.Error())
		}

		initialExpHashes[hash] = i
	}

	for _, i := range updatedExperience {
		hashStruct := hashDomain.Experience{
			i.GetPosition(),
			i.GetCompanyName(),
			i.GetLocation(),
			i.GetDescription(),
			i.GetStartDate(),
			i.GetEndDate(),
		}

		hash, err := hashstructure.Hash(hashStruct, hashstructure.FormatV2, nil)
		if err != nil {
			return fmt.Errorf("err while get hash for new experience %s: %s", i.GetExperienceID().String(), err.Error())
		}

		updatedExpHashes[hash] = i
	}

	expToAdd := make([]experienceDomain.Experience, 0, lenInitialExp)

	for k, v := range updatedExpHashes {
		if _, ok := initialExpHashes[k]; !ok {
			expToAdd = append(expToAdd, v)
		}
	}

	if lenInitialExp != 0 {
		err = s.historyExperienceStorage.CreateExperienceBunch(ctx, initialExperience)
		if err != nil {
			fmt.Println(fmt.Errorf("while CreateExperienceBunch history err: %s", err.Error()))
		}
	}

	err = s.experienceStorage.CreateExperienceBunch(ctx, expToAdd)
	if err != nil {
		return fmt.Errorf("while CreateExperienceBunch updated err: %s", err.Error())
	}

	return nil
}
