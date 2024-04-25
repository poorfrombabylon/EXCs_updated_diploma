package certification

import (
	"context"
	certificationDomain "excs_updater/internal/domain/certification"
	hashDomain "excs_updater/internal/domain/hash"
	profileDomain "excs_updater/internal/domain/profile"
	"fmt"
	"github.com/mitchellh/hashstructure/v2"
)

type CertificationService interface {
	ProcessUpdatedCertification(context.Context, profileDomain.ProfileID, []certificationDomain.Certification) error
}

type CertificationStorage interface {
	GetCertsByProfileID(ctx context.Context, profileID profileDomain.ProfileID) ([]certificationDomain.Certification, error)
	CreateCertificationBunch(ctx context.Context, certificationList []certificationDomain.Certification) error
}

type HistoryCertificationStorage interface {
	CreateCertificationBunch(ctx context.Context, certificationList []certificationDomain.Certification) error
}

type service struct {
	certificationStorage        CertificationStorage
	historyCertificationStorage HistoryCertificationStorage
}

func NewCertificationService(
	certificationStorage CertificationStorage,
	historyCertificationStorage HistoryCertificationStorage,
) CertificationService {
	return &service{
		certificationStorage:        certificationStorage,
		historyCertificationStorage: historyCertificationStorage,
	}
}

func (s *service) ProcessUpdatedCertification(
	ctx context.Context,
	profileID profileDomain.ProfileID,
	updatedCertification []certificationDomain.Certification,
) error {
	initialCertification, err := s.certificationStorage.GetCertsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}

	lenInitialCert := len(initialCertification)
	lenUpdatedCert := len(updatedCertification)

	if lenInitialCert >= lenUpdatedCert {
		return nil
	}

	initialCertHashes := map[uint64]certificationDomain.Certification{}
	updatedCertHashes := map[uint64]certificationDomain.Certification{}

	for _, i := range initialCertification {
		hashStruct := hashDomain.Certification{
			i.GetName(),
			i.GetAuthority(),
			i.GetLicenseNumber(),
			i.GetDisplaySource(),
			i.GetUrl(),
			i.GetAuthorityLinkedinURL(),
			i.GetStartDate(),
			i.GetEndDate(),
		}

		hash, err := hashstructure.Hash(hashStruct, hashstructure.FormatV2, nil)
		if err != nil {
			return fmt.Errorf("err while get hash for initial certification %s: %s", i.GetCertificationID().String(), err.Error())
		}

		initialCertHashes[hash] = i
	}

	for _, i := range updatedCertification {
		hashStruct := hashDomain.Certification{
			i.GetName(),
			i.GetAuthority(),
			i.GetLicenseNumber(),
			i.GetDisplaySource(),
			i.GetUrl(),
			i.GetAuthorityLinkedinURL(),
			i.GetStartDate(),
			i.GetEndDate(),
		}

		hash, err := hashstructure.Hash(hashStruct, hashstructure.FormatV2, nil)
		if err != nil {
			return fmt.Errorf("err while get hash for updated certification %s: %s", i.GetCertificationID().String(), err.Error())
		}

		updatedCertHashes[hash] = i
	}

	certToAdd := make([]certificationDomain.Certification, 0, lenInitialCert)

	for k, v := range updatedCertHashes {
		if _, ok := initialCertHashes[k]; !ok {
			certToAdd = append(certToAdd, v)
		}
	}

	if lenInitialCert != 0 {
		err = s.historyCertificationStorage.CreateCertificationBunch(ctx, initialCertification)
		if err != nil {
			fmt.Println(fmt.Errorf("while create history certs err: %s", err.Error()))
		}
	}

	err = s.certificationStorage.CreateCertificationBunch(ctx, certToAdd)
	if err != nil {
		return fmt.Errorf("while create udpated certs err: %s", err.Error())
	}

	return nil
}
