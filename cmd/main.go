package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"excs_updater/internal/config"
	"excs_updater/internal/domain"
	certificationDomain "excs_updater/internal/domain/certification"
	educationDomain "excs_updater/internal/domain/education"
	experienceDomain "excs_updater/internal/domain/experience"
	profileDomain "excs_updater/internal/domain/profile"
	"excs_updater/internal/domain/queues"
	"excs_updater/internal/libdb"
	"excs_updater/internal/service"
	"excs_updater/internal/storage"
	history_storage "excs_updater/internal/storage_history"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
	"log"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal("error while init config")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	historyPsqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.HistoryPostgres.Host, cfg.HistoryPostgres.Port, cfg.HistoryPostgres.User, cfg.HistoryPostgres.Password, cfg.HistoryPostgres.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("failed to connect to database:", err.Error())
	} else {
		log.Println("connected to db")
	}

	historyDB, err := sql.Open("postgres", historyPsqlInfo)
	if err != nil {
		log.Fatal("failed to connect to history database:", err.Error())
	} else {
		log.Println("connected to db")
	}

	dbx := sqlx.NewDb(db, "pgx")
	libDBWrapper := libdb.NewSQLXDB(dbx)

	historyDBx := sqlx.NewDb(historyDB, "pgx")
	historyLibDBWrapper := libdb.NewSQLXDB(historyDBx)

	storageRegistry, err := storage.NewStorageRegistry(ctx, libDBWrapper, *cfg)
	if err != nil {
		log.Fatal("failed to init storageRegistry:", err.Error())
	}

	historyStorageRegistry, err := history_storage.NewHistoryStorageRegistry(historyLibDBWrapper)
	if err != nil {
		log.Fatal("failed to init HistoryStorageRegistry:", err.Error())
	}

	serviceRegistry := service.NewServiceRegistry(storageRegistry, historyStorageRegistry)

	go startJobUpdate(ctx, serviceRegistry)
	startJobSend(ctx, serviceRegistry)
}

func startJobSend(
	ctx context.Context,
	serviceRegistry *service.Services,
) error {
	totalProfilesSent := 0
	iter := 0
	offset := uint64(0)
	for {
		iter += 1
		fmt.Println("iteration:", iter)

		profiles, err := serviceRegistry.ProfileService.GetBunchProfilesForUpdate(ctx, offset)
		if err != nil {
			return err
		}

		if len(profiles) == 0 || profiles == nil {
			log.Println("no profiles to check")
			time.Sleep(5 * time.Hour)
		}

		for _, i := range profiles {
			region := strings.ToLower(*i.GetCountry())

			link := fmt.Sprintf("https://%s.linkedin.com/in/%s", region, i.GetLinkedInID())

			redisInfo := queues.ProfileToUpdate{
				FirstName:   "",
				LastName:    "",
				ProfileId:   "",
				ProfileLink: link,
			}

			err = serviceRegistry.ProfileService.SendLinkedinID(ctx, redisInfo)
			if err != nil {
				log.Println("err while send:", err.Error())
			}

			i.Check()

			err = serviceRegistry.ProfileService.UpdateProfileFlags(ctx, i)
			if err != nil {
				log.Printf("while update profile flags %s after send err: %s", i.GetProfileID().String(), err.Error())
			}

			totalProfilesSent += 1

			log.Println("sent to scraper profiles:", totalProfilesSent, i.GetProfileID().String())
		}
	}

	return nil
}

func startJobUpdate(
	ctx context.Context,
	serviceRegistry *service.Services,
) error {
	diffExperience := 0
	diffEduсation := 0
	diffCertification := 0
	totalAccsReceived := 0

	subscriber := serviceRegistry.ProfileService.SubscribeToUpdatedProfilesChan(ctx)
	ch := subscriber.Channel()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				errClose := subscriber.Close()
				if errClose != nil {
					return errClose
				}
				fmt.Println("ctx cancel")
			case m := <-ch:
				totalAccsReceived += 1
				res := queues.FullProfileInfo{}

				fmt.Println(m.Payload)
				fmt.Println()

				err := json.Unmarshal([]byte(m.Payload), &res)
				if err != nil {
					fmt.Println("err while json unmarshal", err.Error())
				}

				profileID, err := serviceRegistry.ProfileService.GetProfileIDByLinkedinID(ctx, res.PublicIdentifier)
				if err != nil {
					fmt.Println("err while get profileID", err.Error())
					continue
				}

				newProfileData := profileDomain.NewProfileWithID(
					profileID,
					res.FirstName,
					res.LastName,
					res.Country,
					res.City,
					res.State,
					res.Gender,
					res.Occupation,
					res.Summary,
					res.PublicIdentifier,
					res.IsBlurred,
					res.WwwIsBlurred,
					domain.NewFullModel(),
				)

				if res.Extra != nil {
					newProfileData.SetGithubProfileID(res.Extra.GithubProfileID)
					newProfileData.SetTwitterProfileID(res.Extra.TwitterProfileID)
					newProfileData.SetFacebookProfileID(res.Extra.FacebookProfileID)
					newProfileData.SetTelegramProfileID(res.Extra.TelegramProfileID)
					newProfileData.SetAssociatedLinks(res.Extra.AssociatedLinks)
					newProfileData.SetAssociatedUsernames(res.Extra.AssociatedUsernames)
					newProfileData.SetPersonalEmails(res.PersonalEmails)
					newProfileData.SetPersonalNumbers(res.PersonalNumbers)
				}

				err = serviceRegistry.ProfileService.ProcessUpdatedProfile(ctx, newProfileData)
				if err != nil {
					fmt.Printf("err while ProcessUpdatedProfile for profile id %s: %s", profileID.String(), err.Error())
					continue
				}

				newExperienceList := convertExperience(res.Experiences, profileID)

				err = serviceRegistry.ExperienceService.ProcessUpdatedExperience(ctx, profileID, newExperienceList)
				if err != nil {
					fmt.Printf("err while ProcessUpdatedExperience for profile id %s: %s", profileID.String(), err.Error())
					continue
				}

				newEducationList := convertEducation(res.Education, profileID)

				err = serviceRegistry.EducationService.ProcessUpdatedEducation(ctx, profileID, newEducationList)
				if err != nil {
					fmt.Printf("err while ProcessUpdatedEducation for profile id %s: %s", profileID.String(), err.Error())
					continue
				}

				newCertificationList := convertCertification(res.Certifications, profileID)

				err = serviceRegistry.CertificationService.ProcessUpdatedCertification(ctx, profileID, newCertificationList)
				if err != nil {
					fmt.Printf("err while ProcessUpdatedCertification for profile id %s: %s", profileID.String(), err.Error())
					continue
				}

				log.Printf("\ndiffExperience: %v\ndiffEducation: %v\ndiffCertification: %v\ntotalAccsReceived: %v\n",
					diffExperience, diffEduсation, diffCertification, totalAccsReceived)
			}
		}
	})

	err := group.Wait()
	log.Println("finish")
	if err != nil {
		return err
	}

	log.Printf("----------TOTAL----------\ndiffExperience: %v\ndiffEducation: %v\ndiffCertification: %v\ntotalAccsReceived: %v\n",
		diffExperience, diffEduсation, diffCertification, totalAccsReceived)

	return nil
}

func convertExperience(expDTO []queues.Experience, profileID profileDomain.ProfileID) []experienceDomain.Experience {
	result := make([]experienceDomain.Experience, 0, len(expDTO))

	for _, e := range expDTO {
		experience := experienceDomain.NewExperience(
			profileID,
			e.Position,
			e.CompanyName,
			e.Location,
			e.Description,
			e.CompanyLinkedinProfileUrl,
			getDate(e.StartDate),
			getDate(e.EndDate),
		)

		result = append(result, experience)
	}

	return result
}

func convertEducation(eduDTO []queues.Education, profileID profileDomain.ProfileID) []educationDomain.Education {
	result := make([]educationDomain.Education, 0, len(eduDTO))

	for _, i := range eduDTO {
		education := educationDomain.NewEducation(
			profileID,
			i.FieldOfStudy,
			i.DegreeName,
			i.School,
			i.SchoolLinkedinProfileUrl,
			i.Description,
			i.LogoUrl,
			i.Grade,
			i.ActivitiesAndSocieties,
			getDate(i.StartDate),
			getDate(i.EndDate),
		)

		result = append(result, education)
	}

	return result
}

func convertCertification(certDTO []queues.Certification, profileID profileDomain.ProfileID) []certificationDomain.Certification {
	result := make([]certificationDomain.Certification, 0, len(certDTO))

	for _, i := range certDTO {
		certification := certificationDomain.NewCertification(
			profileID,
			i.Name,
			i.Authority,
			i.LicenseNumber,
			i.DisplaySource,
			i.Url,
			i.AuthorityLinkedinURL,
			getDate(i.StartDate),
			getDate(i.EndDate),
		)

		result = append(result, certification)
	}

	return result
}

func getDate(date *queues.Date) *time.Time {
	if date != nil && date.Day != nil && date.Month != nil && date.Year != nil {
		day := strconv.Itoa(*date.Day)
		if len(day) == 1 {
			day = "0" + day
		}

		month := strconv.Itoa(*date.Month)
		if len(month) == 1 {
			month = "0" + month
		}

		res, err := time.Parse("2006-01-02", fmt.Sprintf("%d-%v-%v", *date.Year, month, day))
		if err != nil {
			log.Println("err while parsing date:", err.Error())
		}

		return &res
	}

	return nil
}
