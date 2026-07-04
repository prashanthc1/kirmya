package api

import "workspace-app/internal/profile/domain"

type profileResponse struct {
	UserID         string              `json:"user_id"`
	Headline       string              `json:"headline"`
	About          string              `json:"about"`
	PhotoURL       string              `json:"photo_url"`
	Bio            string              `json:"bio"`
	Location       string              `json:"location"`
	Website        string              `json:"website"`
	Experiences    []experienceDTO     `json:"experiences"`
	Educations     []educationDTO      `json:"educations"`
	Certifications []certificationDTO  `json:"certifications"`
	Skills         []string            `json:"skills"`
	Languages      []languageDTO       `json:"languages"`
	Portfolio      []portfolioLinkDTO  `json:"portfolio"`
}

type experienceDTO struct {
	ID             string `json:"id,omitempty"`
	Title          string `json:"title"`
	Company        string `json:"company"`
	Location       string `json:"location"`
	EmploymentType string `json:"employment_type"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	IsCurrent      bool   `json:"is_current"`
	Description    string `json:"description"`
}

type educationDTO struct {
	ID           string `json:"id,omitempty"`
	School       string `json:"school"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"field_of_study"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Grade        string `json:"grade"`
	Description  string `json:"description"`
}

type certificationDTO struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name"`
	Issuer        string `json:"issuer"`
	IssueDate     string `json:"issue_date"`
	ExpiryDate    string `json:"expiry_date"`
	CredentialID  string `json:"credential_id"`
	CredentialURL string `json:"credential_url"`
}

type languageDTO struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"`
}

type portfolioLinkDTO struct {
	ID    string `json:"id,omitempty"`
	Label string `json:"label"`
	URL   string `json:"url"`
}

type updateScalarsRequest struct {
	Headline string `json:"headline"`
	About    string `json:"about"`
	PhotoURL string `json:"photo_url"`
	Bio      string `json:"bio"`
	Location string `json:"location"`
	Website  string `json:"website"`
}

type skillsRequest struct {
	Skills []string `json:"skills"`
}

type languagesRequest struct {
	Languages []languageDTO `json:"languages"`
}

type portfolioRequest struct {
	Portfolio []portfolioLinkDTO `json:"portfolio"`
}

// --- mapping ---

func toResponse(p *domain.Profile) profileResponse {
	r := profileResponse{
		UserID: p.UserID, Headline: p.Headline, About: p.About, PhotoURL: p.PhotoURL,
		Bio: p.Bio, Location: p.Location, Website: p.Website, Skills: p.Skills,
		Experiences: make([]experienceDTO, 0, len(p.Experiences)),
		Educations:  make([]educationDTO, 0, len(p.Educations)),
		Certifications: make([]certificationDTO, 0, len(p.Certifications)),
		Languages: make([]languageDTO, 0, len(p.Languages)),
		Portfolio: make([]portfolioLinkDTO, 0, len(p.Portfolio)),
	}
	if r.Skills == nil {
		r.Skills = []string{}
	}
	for _, e := range p.Experiences {
		r.Experiences = append(r.Experiences, experienceDTO(e))
	}
	for _, e := range p.Educations {
		r.Educations = append(r.Educations, educationDTO(e))
	}
	for _, c := range p.Certifications {
		r.Certifications = append(r.Certifications, certificationDTO(c))
	}
	for _, l := range p.Languages {
		r.Languages = append(r.Languages, languageDTO(l))
	}
	for _, l := range p.Portfolio {
		r.Portfolio = append(r.Portfolio, portfolioLinkDTO(l))
	}
	return r
}

func (d experienceDTO) toDomain() domain.WorkExperience       { return domain.WorkExperience(d) }
func (d educationDTO) toDomain() domain.Education              { return domain.Education(d) }
func (d certificationDTO) toDomain() domain.Certification      { return domain.Certification(d) }
