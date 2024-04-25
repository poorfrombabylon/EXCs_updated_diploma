package profile

import (
	"excs_updater/internal/domain"
	"github.com/google/uuid"
)

type ProfileID uuid.UUID

func (a ProfileID) String() string {
	return uuid.UUID(a).String()
}

type Profile struct {
	id                  ProfileID
	firstName           string
	lastName            string
	country             *string
	city                *string
	state               *string
	gender              *string
	occupation          *string
	summary             *string
	linkedInID          string
	isBlurred           *bool
	wwwIsBlurred        *bool
	githubProfileID     *string
	twitterProfileID    *string
	facebookProfileID   *string
	telegramProfileID   *string
	associatedLinks     []string
	associatedUsernames []string
	personalEmails      []string
	personalNumbers     []string
	domain.FullModel
}

func NewProfile(
	firstName string,
	lastName string,
	country *string,
	city *string,
	state *string,
	gender *string,
	occupation *string,
	summary *string,
	linkedInID string,
	isBlurred *bool,
	wwwIsBlurred *bool,
	githubProfileID *string,
	twitterProfileID *string,
	facebookProfileID *string,
	telegramProfileID *string,
) Profile {
	return Profile{
		id:                ProfileID(uuid.New()),
		firstName:         firstName,
		lastName:          lastName,
		country:           country,
		city:              city,
		state:             state,
		gender:            gender,
		occupation:        occupation,
		summary:           summary,
		linkedInID:        linkedInID,
		isBlurred:         isBlurred,
		wwwIsBlurred:      wwwIsBlurred,
		githubProfileID:   githubProfileID,
		twitterProfileID:  twitterProfileID,
		facebookProfileID: facebookProfileID,
		telegramProfileID: telegramProfileID,
		FullModel:         domain.NewFullModel(),
	}
}

func NewProfileWithID(
	id ProfileID,
	firstName string,
	lastName string,
	country *string,
	city *string,
	state *string,
	gender *string,
	occupation *string,
	summary *string,
	linkedInID string,
	isBlurred *bool,
	wwwIsBlurred *bool,
	model domain.FullModel,
) Profile {
	return Profile{
		id:           id,
		firstName:    firstName,
		lastName:     lastName,
		country:      country,
		city:         city,
		state:        state,
		gender:       gender,
		occupation:   occupation,
		summary:      summary,
		linkedInID:   linkedInID,
		isBlurred:    isBlurred,
		wwwIsBlurred: wwwIsBlurred,
		FullModel:    model,
	}
}

func (p Profile) GetProfileID() ProfileID {
	return p.id
}

func (p Profile) GetFirstName() string {
	return p.firstName
}

func (p *Profile) SetFirstName(firstName string) {
	p.firstName = firstName
}

func (p Profile) GetLastName() string {
	return p.lastName
}

func (p *Profile) SetLastName(lastName string) {
	p.lastName = lastName
}

func (p Profile) GetCountry() *string {
	return p.country
}

func (p *Profile) SetCountry(country *string) {
	p.country = country
}

func (p Profile) GetCity() *string {
	return p.city
}

func (p *Profile) SetCity(city *string) {
	p.city = city
}

func (p Profile) GetState() *string {
	return p.state
}

func (p *Profile) SetState(state *string) {
	p.state = state
}

func (p Profile) GetGender() *string {
	return p.gender
}

func (p *Profile) SetGender(gender *string) {
	p.gender = gender
}

func (p Profile) GetOccupation() *string {
	return p.occupation
}

func (p *Profile) SetOccupation(occupation *string) {
	p.occupation = occupation
}

func (p Profile) GetSummary() *string {
	return p.summary
}

func (p *Profile) SetSummary(summary *string) {
	p.summary = summary
}

func (p Profile) GetLinkedInID() string {
	return p.linkedInID
}

func (p Profile) GetIsBlurred() *bool {
	return p.isBlurred
}

func (p *Profile) SetIsBlurred(isBlurred *bool) {
	p.isBlurred = isBlurred
}

func (p Profile) GetWwwIsBlurred() *bool {
	return p.wwwIsBlurred
}

func (p *Profile) SetWwwIsBlurred(wwwIsBlurred *bool) {
	p.wwwIsBlurred = wwwIsBlurred
}

func (p Profile) GetGithubProfileID() *string {
	return p.githubProfileID
}

func (p *Profile) SetGithubProfileID(githubID *string) {
	p.githubProfileID = githubID
}

func (p Profile) GetTwitterProfileID() *string {
	return p.twitterProfileID
}

func (p *Profile) SetTwitterProfileID(twitterID *string) {
	p.twitterProfileID = twitterID
}

func (p Profile) GetFacebookProfileID() *string {
	return p.facebookProfileID
}

func (p *Profile) SetFacebookProfileID(facebookID *string) {
	p.facebookProfileID = facebookID
}

func (p Profile) GetTelegramProfileID() *string {
	return p.telegramProfileID
}

func (p *Profile) SetTelegramProfileID(telegramID *string) {
	p.telegramProfileID = telegramID
}

func (p Profile) GetAssociatedLinks() []string {
	return p.associatedLinks
}

func (p *Profile) SetAssociatedLinks(links []string) {
	p.associatedLinks = links
}

func (p Profile) GetAssociatedUsernames() []string {
	return p.associatedUsernames
}

func (p *Profile) SetAssociatedUsernames(usernames []string) {
	p.associatedUsernames = usernames
}

func (p Profile) GetPersonalEmails() []string {
	return p.personalEmails
}

func (p *Profile) SetPersonalEmails(emails []string) {
	p.personalEmails = emails
}

func (p Profile) GetPersonalNumbers() []string {
	return p.personalNumbers
}

func (p *Profile) SetPersonalNumbers(numbers []string) {
	p.personalNumbers = numbers
}
