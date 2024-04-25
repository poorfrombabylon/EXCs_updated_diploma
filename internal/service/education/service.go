package education

import (
	"context"
	educationDomain "excs_updater/internal/domain/education"
	hashDomain "excs_updater/internal/domain/hash"
	profileDomain "excs_updater/internal/domain/profile"
	"fmt"
	"github.com/mitchellh/hashstructure/v2"
)

type EducationService interface {
	ProcessUpdatedEducation(context.Context, profileDomain.ProfileID, []educationDomain.Education) error
}

type EducationStorage interface {
	GetEducationByProfileID(ctx context.Context, profileID profileDomain.ProfileID) ([]educationDomain.Education, error)
	CreateEducationBunch(ctx context.Context, educationList []educationDomain.Education) error
}

type HistoryEducationStorage interface {
	CreateEducationBunch(ctx context.Context, educationList []educationDomain.Education) error
}

type service struct {
	educationStorage        EducationStorage
	historyEducationStorage HistoryEducationStorage
}

func NewEducationService(
	educationStorage EducationStorage,
	historyEducationStorage HistoryEducationStorage,
) EducationService {
	return &service{
		educationStorage:        educationStorage,
		historyEducationStorage: historyEducationStorage,
	}
}

func (s *service) ProcessUpdatedEducation(
	ctx context.Context,
	profileID profileDomain.ProfileID,
	updatedEducation []educationDomain.Education,
) error {
	initialEducation, err := s.educationStorage.GetEducationByProfileID(ctx, profileID)
	if err != nil {
		return err
	}

	lenInitialEducation := len(initialEducation)
	lenUpdatedEducation := len(updatedEducation)

	if lenInitialEducation >= lenUpdatedEducation {
		return nil
	}

	initialEduHashes := map[uint64]educationDomain.Education{}
	updatedEduHashes := map[uint64]educationDomain.Education{}

	for _, i := range initialEducation {
		hashStruct := hashDomain.Education{
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
		}

		hash, err := hashstructure.Hash(hashStruct, hashstructure.FormatV2, nil)
		if err != nil {
			return fmt.Errorf("err while get hash for initial education %s: %s", i.GetEducationID().String(), err.Error())
		}

		initialEduHashes[hash] = i
	}

	for _, i := range updatedEducation {
		hashStruct := hashDomain.Education{
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
		}

		hash, err := hashstructure.Hash(hashStruct, hashstructure.FormatV2, nil)
		if err != nil {
			return fmt.Errorf("err while get hash for updated education %s: %s", i.GetEducationID().String(), err.Error())
		}

		updatedEduHashes[hash] = i
	}

	eduToAdd := make([]educationDomain.Education, 0, lenInitialEducation)

	for k, v := range updatedEduHashes {
		if _, ok := initialEduHashes[k]; !ok {
			eduToAdd = append(eduToAdd, v)
		}
	}

	if lenInitialEducation != 0 {
		err = s.historyEducationStorage.CreateEducationBunch(ctx, initialEducation)
		if err != nil {
			fmt.Println(fmt.Errorf("while create history education err: %s", err.Error()))
		}
	}

	err = s.educationStorage.CreateEducationBunch(ctx, eduToAdd)
	if err != nil {
		return fmt.Errorf("while create updated education err: %s", err.Error())
	}

	return nil
}
