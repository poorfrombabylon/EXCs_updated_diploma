package profile

import (
	"context"
	hashDomain "excs_updater/internal/domain/hash"
	profileDomain "excs_updater/internal/domain/profile"
	redisQueueDomain "excs_updater/internal/domain/queues"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mitchellh/hashstructure/v2"
)

type ProfileService interface {
	ProcessUpdatedProfile(ctx context.Context, newProfile profileDomain.Profile) error
	GetBunchProfilesForUpdate(ctx context.Context, page uint64) ([]profileDomain.Profile, error)
	GetProfileByLinkedInID(ctx context.Context, linkedinID string) (profileDomain.Profile, error)
	SendLinkedinID(ctx context.Context, profile redisQueueDomain.ProfileToUpdate) error
	UpdateProfileFlags(ctx context.Context, profile profileDomain.Profile) error
	UpdateFullProfile(ctx context.Context, profile profileDomain.Profile) error
	SubscribeToUpdatedProfilesChan(_ context.Context) *redis.PubSub
	GetProfileIDByLinkedinID(ctx context.Context, linkedinID string) (profileDomain.ProfileID, error)
}

type ProfileStorage interface {
	GetProfileByLinkedInID(ctx context.Context, linkedinID string) (profileDomain.Profile, error)
	UpdateFullProfile(ctx context.Context, profile profileDomain.Profile) error
	GetBunchProfilesForUpdate(ctx context.Context, page uint64) ([]profileDomain.Profile, error)
	UpdateProfileFlags(ctx context.Context, profile profileDomain.Profile) error
	GetProfileIDByLinkedinID(ctx context.Context, linkedinID string) (profileDomain.ProfileID, error)
}

type ProfileSendRedis interface {
	SendLinkedinID(ctx context.Context, profile redisQueueDomain.ProfileToUpdate) error
}

type ProfileReceiveRedis interface {
	SubscribeToSubscriberChannel(_ context.Context) *redis.PubSub
}

type HistoryProfileStorage interface {
	CreateHistoryProfile(ctx context.Context, profile profileDomain.Profile) error
}

type service struct {
	profileStorage        ProfileStorage
	historyProfileStorage HistoryProfileStorage
	profileSendRedis      ProfileSendRedis
	profileReceiveRedis   ProfileReceiveRedis
}

func NewProfileService(
	profileStorage ProfileStorage,
	historyProfileStorage HistoryProfileStorage,
	profileSendRedis ProfileSendRedis,
	profileReceiveRedis ProfileReceiveRedis,
) ProfileService {
	return &service{
		profileStorage:        profileStorage,
		historyProfileStorage: historyProfileStorage,
		profileSendRedis:      profileSendRedis,
		profileReceiveRedis:   profileReceiveRedis,
	}
}

func (s *service) ProcessUpdatedProfile(ctx context.Context, newProfile profileDomain.Profile) error {
	initialProfile, err := s.profileStorage.GetProfileByLinkedInID(ctx, newProfile.GetLinkedInID())
	if err != nil {
		return fmt.Errorf("ProcessUpdatedProfile: err while get old data: %s", err.Error())
	}

	npHashData := hashDomain.Profile{
		FirstName:           newProfile.GetFirstName(),
		LastName:            newProfile.GetLastName(),
		Country:             newProfile.GetCountry(),
		City:                newProfile.GetCity(),
		State:               newProfile.GetCity(),
		Gender:              newProfile.GetGender(),
		Occupation:          newProfile.GetOccupation(),
		Summary:             newProfile.GetSummary(),
		TelegramID:          newProfile.GetTelegramProfileID(),
		FacebookID:          newProfile.GetFacebookProfileID(),
		TwitterID:           newProfile.GetTwitterProfileID(),
		GithubID:            newProfile.GetGithubProfileID(),
		AssociatedLinks:     newProfile.GetAssociatedLinks(),
		AssociatedUsernames: newProfile.GetAssociatedUsernames(),
		PersonalEmails:      newProfile.GetPersonalEmails(),
		PersonalNumbers:     newProfile.GetPersonalNumbers(),
	}

	ipHashData := hashDomain.Profile{
		FirstName:           initialProfile.GetFirstName(),
		LastName:            initialProfile.GetLastName(),
		Country:             initialProfile.GetCountry(),
		City:                initialProfile.GetCity(),
		State:               initialProfile.GetCity(),
		Gender:              initialProfile.GetGender(),
		Occupation:          initialProfile.GetOccupation(),
		Summary:             initialProfile.GetSummary(),
		TelegramID:          initialProfile.GetTelegramProfileID(),
		FacebookID:          initialProfile.GetFacebookProfileID(),
		TwitterID:           initialProfile.GetTwitterProfileID(),
		GithubID:            initialProfile.GetGithubProfileID(),
		AssociatedLinks:     initialProfile.GetAssociatedLinks(),
		AssociatedUsernames: initialProfile.GetAssociatedUsernames(),
		PersonalEmails:      initialProfile.GetPersonalEmails(),
		PersonalNumbers:     initialProfile.GetPersonalNumbers(),
	}

	newProfileHash, _ := hashstructure.Hash(npHashData, hashstructure.FormatV2, nil)
	initialProfileHash, _ := hashstructure.Hash(ipHashData, hashstructure.FormatV2, nil)

	if newProfileHash != initialProfileHash {
		resultProfile := initialProfile
		resultProfile.SetIsBlurred(newProfile.GetIsBlurred())
		resultProfile.SetWwwIsBlurred(newProfile.GetWwwIsBlurred())

		if newFirstName := newProfile.GetFirstName(); newFirstName != "" && newFirstName != initialProfile.GetFirstName() {
			resultProfile.SetFirstName(newFirstName)
		}

		if newLastName := newProfile.GetLastName(); newLastName != "" && newLastName != initialProfile.GetLastName() {
			resultProfile.SetLastName(newLastName)
		}

		if newCountry := newProfile.GetCountry(); newCountry != nil && newCountry != initialProfile.GetCountry() {
			resultProfile.SetCountry(newCountry)
		}

		if newCity := newProfile.GetCity(); newCity != nil && newCity != initialProfile.GetCity() {
			resultProfile.SetCity(newCity)
		}

		if newState := newProfile.GetState(); newState != nil && newState != initialProfile.GetState() {
			resultProfile.SetState(newState)
		}

		if newGender := newProfile.GetGender(); newGender != nil && newGender != initialProfile.GetGender() {
			resultProfile.SetGender(newGender)
		}

		if newOccupation := newProfile.GetOccupation(); newOccupation != nil && newOccupation != initialProfile.GetOccupation() {
			resultProfile.SetOccupation(newOccupation)
		}

		if newSummary := newProfile.GetSummary(); newSummary != nil && newSummary != initialProfile.GetSummary() {
			resultProfile.SetSummary(newSummary)
		}

		if newTelegramID := newProfile.GetTelegramProfileID(); newTelegramID != nil && newTelegramID != initialProfile.GetTelegramProfileID() {
			resultProfile.SetTelegramProfileID(newTelegramID)
		}

		if newFacebookID := newProfile.GetFacebookProfileID(); newFacebookID != nil && newFacebookID != initialProfile.GetFacebookProfileID() {
			resultProfile.SetFacebookProfileID(newFacebookID)
		}

		if newTwitterID := newProfile.GetTwitterProfileID(); newTwitterID != nil && newTwitterID != initialProfile.GetTwitterProfileID() {
			resultProfile.SetTwitterProfileID(newTwitterID)
		}

		if newGithubID := newProfile.GetGithubProfileID(); newGithubID != nil && newGithubID != initialProfile.GetGithubProfileID() {
			resultProfile.SetGithubProfileID(newGithubID)
		}

		if newPersonalEmails := newProfile.GetPersonalEmails(); newPersonalEmails != nil && len(newPersonalEmails) > len(initialProfile.GetPersonalEmails()) {
			resultProfile.SetPersonalEmails(newPersonalEmails)
		}

		if newPersonalNumbers := newProfile.GetPersonalNumbers(); newPersonalNumbers != nil && len(newPersonalNumbers) > len(initialProfile.GetPersonalNumbers()) {
			resultProfile.SetPersonalNumbers(newPersonalNumbers)
		}

		if newAssociatedLinks := newProfile.GetAssociatedLinks(); newAssociatedLinks != nil && len(newAssociatedLinks) > len(initialProfile.GetAssociatedLinks()) {
			resultProfile.SetAssociatedLinks(newAssociatedLinks)
		}

		if newAssociatedUsernames := newProfile.GetAssociatedUsernames(); newAssociatedUsernames != nil && len(newAssociatedUsernames) > len(initialProfile.GetAssociatedUsernames()) {
			resultProfile.SetAssociatedUsernames(newAssociatedUsernames)
		}

		err = s.historyProfileStorage.CreateHistoryProfile(ctx, initialProfile)
		if err != nil {
			fmt.Println(fmt.Errorf("ProcessUpdatedProfile: err while create history profile: %s", err.Error()))
		}

		resultProfile.Update()

		err = s.profileStorage.UpdateFullProfile(ctx, resultProfile)
		if err != nil {
			return fmt.Errorf("ProcessUpdatedProfile: err while update resultProfile: %s", err.Error())
		}
	}

	return nil
}

func (s *service) GetBunchProfilesForUpdate(ctx context.Context, page uint64) ([]profileDomain.Profile, error) {
	return s.profileStorage.GetBunchProfilesForUpdate(ctx, page)
}

func (s *service) GetProfileByLinkedInID(ctx context.Context, linkedinID string) (profileDomain.Profile, error) {
	return s.profileStorage.GetProfileByLinkedInID(ctx, linkedinID)
}

func (s *service) SendLinkedinID(ctx context.Context, profile redisQueueDomain.ProfileToUpdate) error {
	return s.profileSendRedis.SendLinkedinID(ctx, profile)
}

func (s *service) SubscribeToUpdatedProfilesChan(ctx context.Context) *redis.PubSub {
	return s.profileReceiveRedis.SubscribeToSubscriberChannel(ctx)
}

func (s *service) UpdateProfileFlags(ctx context.Context, profile profileDomain.Profile) error {
	return s.profileStorage.UpdateProfileFlags(ctx, profile)
}

func (s *service) UpdateFullProfile(ctx context.Context, profile profileDomain.Profile) error {
	return s.profileStorage.UpdateFullProfile(ctx, profile)
}

func (s *service) GetProfileIDByLinkedinID(ctx context.Context, linkedinID string) (profileDomain.ProfileID, error) {
	return s.profileStorage.GetProfileIDByLinkedinID(ctx, linkedinID)
}
