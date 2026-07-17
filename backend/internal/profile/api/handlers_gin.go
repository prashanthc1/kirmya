package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"workspace-app/internal/common"
	"workspace-app/internal/connections"
	"workspace-app/internal/platform/cache"
	"workspace-app/internal/profile/application"
	"workspace-app/internal/profile/domain"
)

type GinHandler struct {
	svc   *application.Service
	db    *sql.DB
	cache cache.Cache
}

func NewGinHandler(svc *application.Service, db *sql.DB, c cache.Cache) *GinHandler {
	return &GinHandler{
		svc:   svc,
		db:    db,
		cache: c,
	}
}

// adaptMiddleware converts a standard http.Handler middleware into a Gin-compatible middleware
func adaptMiddleware(auth func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			c.Request = r
			c.Next()
		})
		auth(next).ServeHTTP(c.Writer, c.Request)
		if !called {
			c.Abort()
		}
	}
}

// RegisterGinRoutes mounts the profile endpoints onto the Gin engine
func RegisterGinRoutes(r *gin.Engine, db *sql.DB, auth func(http.Handler) http.Handler, svc *application.Service, c cache.Cache) {
	h := NewGinHandler(svc, db, c)
	authGin := adaptMiddleware(auth)
	optAuthGin := optionalAuth(auth)

	// Me endpoint (authenticated)
	r.GET("/api/profile/me", authGin, h.GetMe)

	// Public profile endpoint (no auth required, but optional token check)
	r.GET("/api/profile/public/:userId", optAuthGin, h.GetPublicProfile)

	// PATCH endpoints (authenticated)
	g := r.Group("/api/profile/me")
	g.Use(authGin)
	{
		g.PATCH("/basic-info", h.PatchBasicInfo)
		g.PATCH("/summary", h.PatchSummary)
		g.PATCH("/experience", h.PatchExperience)
		g.PATCH("/education", h.PatchEducation)
		g.PATCH("/skills", h.PatchSkills)
		g.PATCH("/certifications", h.PatchCertifications)
		g.PATCH("/projects", h.PatchProjects)
		g.PATCH("/achievements", h.PatchAchievements)
		g.PATCH("/contact", h.PatchContact)
		g.PATCH("/preferences", h.PatchPreferences)
	}
}

// optionalAuth adapts the auth middleware but allows anonymous access when no Authorization header is present
func optionalAuth(auth func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			c.Request = r
			c.Next()
		})
		hasAuth := c.Request.Header.Get("Authorization") != ""
		if hasAuth {
			auth(next).ServeHTTP(c.Writer, c.Request)
			if !called {
				c.Abort()
			}
		} else {
			c.Next()
		}
	}
}

// GetMe handles GET /api/profile/me (decrypted own profile)
func (h *GinHandler) GetMe(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(p))
}

// GetPublicProfile handles GET /api/profile/public/:userId
func (h *GinHandler) GetPublicProfile(c *gin.Context) {
	targetID := c.Param("userId")
	ctx := c.Request.Context()

	// 1. Check if user is deactivated or deleted
	var status string
	err := h.db.QueryRowContext(ctx, `
		SELECT status FROM users WHERE id = $1 AND deleted_at IS NULL
	`, targetID).Scan(&status)
	if err == sql.ErrNoRows || status == "deactivated" {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "PROFILE_NOT_FOUND", "message": "profile not found"}})
		return
	}

	// 2. Fetch the published profile snapshot
	p, err := h.svc.GetPublished(ctx, targetID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "PROFILE_NOT_FOUND", "message": "profile not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	// 3. Determine connection status if viewer is logged in
	viewerID := common.UserIDFromContext(ctx)
	isConnection := false
	viewerCanMessage := false

	if viewerID != "" && viewerID != targetID {
		var err error
		isConnection, err = connections.CanMessage(ctx, h.db, viewerID, targetID)
		if err == nil && isConnection {
			viewerCanMessage = true
		}
	} else if viewerID == targetID {
		isConnection = true
		viewerCanMessage = true
	}

	// 4. Check profile level visibility
	visProfile := p.Privacy.FieldVisibility["profile"]
	if visProfile == "private" && viewerID != targetID {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "PROFILE_PRIVATE", "message": "profile is private"}})
		return
	}
	if visProfile == "connections_only" && !isConnection {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "profile is connections-only"}})
		return
	}

	// 5. Increment views if viewer != target and not deduped within 1 hour
	if viewerID != targetID {
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		actorKey := viewerID
		if actorKey == "" {
			actorKey = ip
		}
		dedupeKey := "view:dedupe:" + actorKey + ":" + targetID
		if _, ok := h.cache.Get(ctx, dedupeKey); !ok {
			h.cache.Set(ctx, dedupeKey, []byte("1"), 1*time.Hour)
			go func() {
				var actorPtr *string
				if viewerID != "" {
					actorPtr = &viewerID
				}
				_ = h.svc.RecordView(context.Background(), targetID, actorPtr, ip, ua)
			}()
		}
	}

	// 6. Map and filter public profile fields
	resp := h.filterPublicProfile(p, isConnection, viewerCanMessage)
	c.JSON(http.StatusOK, resp)
}

// filterPublicProfile applies server-side visibility rules to construct a safe PublicProfileResponse
func (h *GinHandler) filterPublicProfile(p *domain.Profile, isConnection, viewerCanMessage bool) map[string]any {
	// Construct map dynamically to fully omit keys (rather than setting to null)
	resp := make(map[string]any)

	resp["user_id"] = p.UserID
	resp["is_connection"] = isConnection
	resp["viewer_can_message"] = viewerCanMessage

	// Identity / Contact Details (sensitive, only visible to connections if public/connections)
	visProfile := p.Privacy.FieldVisibility["profile"]
	if visProfile == "public" || (visProfile == "connections_only" && isConnection) {
		resp["headline"] = p.Identity.Headline
		resp["about"] = p.Identity.Bio
		resp["photo_url"] = p.Identity.PhotoURL
		resp["cover_url"] = p.Identity.CoverURL
		resp["bio"] = p.Identity.Bio
		resp["location"] = p.Identity.Location
		resp["website"] = p.Identity.SocialLinks.Website
		resp["version"] = p.Version
		resp["pronouns"] = p.Identity.VisaStatus
		resp["career_status"] = p.Identity.Availability

		// Sensitive Identity fields: direct email/phone/address, work authorizations, visa, driving license
		if isConnection {
			resp["email"] = p.Identity.Email
			resp["phone"] = p.Identity.Phone
			resp["address"] = p.Identity.Address
			resp["work_auth_status"] = p.Identity.WorkAuthorization
			resp["passport_nationality"] = p.Identity.Nationality
			resp["driving_license_bool"] = p.Verification.IdentityVerified
			resp["driving_license_type"] = p.Identity.VisaStatus
			resp["preferred_contact_channel"] = p.Identity.PreferredContactChannel
			resp["video_intro_url"] = p.Identity.CoverURL
		}
	}

	// Experiences (governed by visibility_experience)
	visExp := p.Privacy.FieldVisibility["experience"]
	if visExp == "public" || (visExp == "connections_only" && isConnection) {
		resp["experiences"] = toExperiencesDTOs(p.Experiences)
	}

	// Educations (governed by visibility_education)
	visEdu := p.Privacy.FieldVisibility["education"]
	if visEdu == "public" || (visEdu == "connections_only" && isConnection) {
		resp["educations"] = toEducationsDTOs(p.Educations)
	}

	// Certifications (governed by visibility_certifications)
	visCert := p.Privacy.FieldVisibility["certifications"]
	if visCert == "public" || (visCert == "connections_only" && isConnection) {
		resp["certifications"] = toCertificationsDTOs(p.Certifications)
	}

	// Skills (governed by visibility_skills)
	visSkills := p.Privacy.FieldVisibility["skills"]
	if visSkills == "public" || (visSkills == "connections_only" && isConnection) {
		resp["skills"] = toSkillsDTOs(p.Skills)
	}

	// Projects (governed by visibility_portfolio)
	visPortfolio := p.Privacy.FieldVisibility["portfolio"]
	if visPortfolio == "public" || (visPortfolio == "connections_only" && isConnection) {
		resp["portfolio"] = toProjectsDTOs(p.Projects)
	}

	// Summary (Section 2 - default to public profile visibility rules)
	if visProfile == "public" || (visProfile == "connections_only" && isConnection) {
		resp["summary"] = map[string]any{
			"executive_summary":        p.Summary.ExecutiveSummary,
			"career_objectives":        p.Summary.CareerObjectives,
			"career_highlights":        p.Summary.CareerHighlights,
			"industries":               p.Summary.Industries,
			"functional_areas":         p.Summary.FunctionalAreas,
			"personal_brand_statement": p.Summary.PersonalBrandStatement,
			"elevator_pitch":           p.Summary.ElevatorPitch,
		}
	}

	// Preferences (salary only connection-gated if salary visibility != private)
	visSalary := p.Privacy.FieldVisibility["salary"]
	if isConnection && (visSalary == "public" || visSalary == "connections_only") {
		resp["salary_min"] = p.Preferences.SalaryMin
		resp["salary_max"] = p.Preferences.SalaryMax
		resp["salary_currency"] = p.Preferences.SalaryCurrency
	}

	// References (references visibility)
	visRef := p.Privacy.FieldVisibility["references"]
	if isConnection && (visRef == "public" || visRef == "connections_only") {
		// Mocked array if needed
		resp["references"] = []any{}
	}

	return resp
}

// PATCH /api/profile/me/basic-info
type patchBasicInfoRequest struct {
	PreferredName           *string               `json:"preferred_name"`
	TimeZone                *string               `json:"timezone"`
	Nationality             *string               `json:"nationality"`
	Headline                *string               `json:"headline"`
	Bio                     *string               `json:"bio"`
	PhotoURL                *string               `json:"photo_url"`
	CoverURL                *string               `json:"cover_url"`
	WorkAuthorization       *string               `json:"work_authorization"`
	PreferredContactChannel *string               `json:"preferred_contact_channel"`
	VideoIntroURL           *string               `json:"video_intro_url"`
	Location                *string               `json:"location"`
	Languages               *[]domain.LanguageItem `json:"languages"`
}

func (h *GinHandler) PatchBasicInfo(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchBasicInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	if req.PreferredName != nil {
		p.Identity.PreferredName = *req.PreferredName
	}
	if req.TimeZone != nil {
		p.Identity.TimeZone = *req.TimeZone
	}
	if req.Nationality != nil {
		p.Identity.Nationality = *req.Nationality
	}
	if req.Headline != nil {
		p.Identity.Headline = *req.Headline
	}
	if req.Bio != nil {
		p.Identity.Bio = *req.Bio
	}
	if req.PhotoURL != nil {
		p.Identity.PhotoURL = *req.PhotoURL
	}
	if req.CoverURL != nil {
		p.Identity.CoverURL = *req.CoverURL
	}
	if req.WorkAuthorization != nil {
		p.Identity.WorkAuthorization = *req.WorkAuthorization
	}
	if req.PreferredContactChannel != nil {
		p.Identity.PreferredContactChannel = *req.PreferredContactChannel
	}
	if req.VideoIntroURL != nil {
		p.Identity.CoverURL = *req.VideoIntroURL
	}
	if req.Location != nil {
		p.Identity.Location = *req.Location
	}
	if req.Languages != nil {
		p.Identity.Languages = *req.Languages
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Identity: &p.Identity,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

// PATCH /api/profile/me/summary
type patchSummaryRequest struct {
	ExecutiveSummary       *string   `json:"executive_summary"`
	CareerObjectives       *string   `json:"career_objectives"`
	CareerHighlights       *[]string `json:"career_highlights"`
	Industries             *[]string `json:"industries"`
	FunctionalAreas        *[]string `json:"functional_areas"`
	PersonalBrandStatement *string   `json:"personal_brand_statement"`
	ElevatorPitch          *string   `json:"elevator_pitch"`
}

func (h *GinHandler) PatchSummary(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	if req.ExecutiveSummary != nil {
		p.Summary.ExecutiveSummary = *req.ExecutiveSummary
	}
	if req.CareerObjectives != nil {
		p.Summary.CareerObjectives = *req.CareerObjectives
	}
	if req.CareerHighlights != nil {
		p.Summary.CareerHighlights = *req.CareerHighlights
	}
	if req.Industries != nil {
		p.Summary.Industries = *req.Industries
	}
	if req.FunctionalAreas != nil {
		p.Summary.FunctionalAreas = *req.FunctionalAreas
	}
	if req.PersonalBrandStatement != nil {
		p.Summary.PersonalBrandStatement = *req.PersonalBrandStatement
	}
	if req.ElevatorPitch != nil {
		p.Summary.ElevatorPitch = *req.ElevatorPitch
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Summary: &p.Summary,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

// PATCH /api/profile/me/experience
type patchExperienceRequest struct {
	Experience []domain.WorkExperience `json:"experience"`
}

func (h *GinHandler) PatchExperience(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Experiences: &req.Experience,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

// PATCH /api/profile/me/education
type patchEducationRequest struct {
	Education []domain.Education `json:"education"`
}

func (h *GinHandler) PatchEducation(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchEducationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Educations: &req.Education,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

// PATCH /api/profile/me/skills
type patchSkillsRequest struct {
	Skills []domain.SkillItem `json:"skills"`
}

func (h *GinHandler) PatchSkills(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchSkillsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Skills: &req.Skills,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

// PATCH /api/profile/me/certifications
type patchCertificationsRequest struct {
	Certifications []domain.CertificationItem `json:"certifications"`
}

func (h *GinHandler) PatchCertifications(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchCertificationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Certifications: &req.Certifications,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

// PATCH /api/profile/me/projects
type patchProjectsRequest struct {
	Projects []domain.ProjectItem `json:"projects"`
}

func (h *GinHandler) PatchProjects(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchProjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Projects: &req.Projects,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

// PATCH /api/profile/me/achievements
type patchAchievementsRequest struct {
	Achievements []domain.AchievementItem `json:"achievements"`
}

func (h *GinHandler) PatchAchievements(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchAchievementsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Achievements: &req.Achievements,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

// PATCH /api/profile/me/contact
type patchContactRequest struct {
	Email   *string `json:"email"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`
}

func (h *GinHandler) PatchContact(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	if req.Email != nil {
		if *req.Email != "" && !strings.Contains(*req.Email, "@") {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": "invalid email format"}})
			return
		}
		p.Identity.Email = *req.Email
	}
	if req.Phone != nil {
		p.Identity.Phone = *req.Phone
	}
	if req.Address != nil {
		p.Identity.Address = *req.Address
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Identity: &p.Identity,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

// PATCH /api/profile/me/preferences
type patchPreferencesRequest struct {
	JobPreferences     *jobPreferencesDTO     `json:"job_preferences"`
	SalaryExpectations *salaryExpectationsDTO `json:"salary_expectations"`
	VisibilitySettings *map[string]string     `json:"visibility_settings"`
}

type jobPreferencesDTO struct {
	DesiredRoles           *[]string `json:"desired_roles"`
	DesiredIndustries      *[]string `json:"desired_industries"`
	EmploymentTypes        *[]string `json:"employment_types"`
	NoticePeriod           *string   `json:"notice_period"`
	RemotePreference       *string   `json:"remote_preference"`
	OpenToRelocation       *bool     `json:"open_to_relocation"`
	PreferredCountries     *[]string `json:"preferred_countries"`
	PreferredCities        *[]string `json:"preferred_cities"`
	TravelWillingness      *string   `json:"travel_willingness"`
	CompanySizePreferences *[]string `json:"company_size_preferences"`
}

type salaryExpectationsDTO struct {
	SalaryMin      *int    `json:"salary_min"`
	SalaryMax      *int    `json:"salary_max"`
	SalaryCurrency *string `json:"salary_currency"`
}

func (h *GinHandler) PatchPreferences(c *gin.Context) {
	uid := common.UserIDFromContext(c.Request.Context())
	p, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		h.writeError(c, err)
		return
	}

	var req patchPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}

	if req.JobPreferences != nil {
		jp := req.JobPreferences
		if jp.DesiredRoles != nil {
			p.Preferences.DesiredRoles = *jp.DesiredRoles
		}
		if jp.DesiredIndustries != nil {
			p.Preferences.DesiredIndustries = *jp.DesiredIndustries
		}
		if jp.NoticePeriod != nil {
			p.Preferences.NoticePeriod = *jp.NoticePeriod
		}
		if jp.RemotePreference != nil {
			p.Preferences.RemotePreference = *jp.RemotePreference
		}
		if jp.OpenToRelocation != nil {
			p.Preferences.OpenToRelocation = *jp.OpenToRelocation
		}
		if jp.PreferredCountries != nil {
			p.Preferences.PreferredCountries = *jp.PreferredCountries
		}
		if jp.PreferredCities != nil {
			p.Preferences.PreferredCities = *jp.PreferredCities
		}
		if jp.TravelWillingness != nil {
			p.Preferences.TravelWillingness = *jp.TravelWillingness
		}
		if jp.CompanySizePreferences != nil {
			p.Preferences.CompanySizePreferences = *jp.CompanySizePreferences
		}
	}

	if req.SalaryExpectations != nil {
		se := req.SalaryExpectations
		if se.SalaryMin != nil {
			p.Preferences.SalaryMin = *se.SalaryMin
		}
		if se.SalaryMax != nil {
			p.Preferences.SalaryMax = *se.SalaryMax
		}
		if se.SalaryCurrency != nil {
			p.Preferences.SalaryCurrency = *se.SalaryCurrency
		}
	}

	if req.VisibilitySettings != nil {
		for k, v := range *req.VisibilitySettings {
			p.Privacy.FieldVisibility[k] = v
		}
	}

	updated, err := h.svc.UpdateProfile(c.Request.Context(), uid, p.Version, domain.AggregateUpdate{
		Preferences: &p.Preferences,
		Privacy:     &p.Privacy,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(updated))
}

func (h *GinHandler) writeError(c *gin.Context, err error) {
	code := "VALIDATION_ERROR"
	msg := err.Error()
	status := http.StatusBadRequest

	if errors.Is(err, domain.ErrNotFound) {
		code = "PROFILE_NOT_FOUND"
		msg = "profile not found"
		status = http.StatusNotFound
	} else if errors.Is(err, domain.ErrOptimisticLock) {
		code = "CONFLICT_ERROR"
		msg = "profile was modified by another request; reload and retry"
		status = http.StatusConflict
	}

	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": msg,
		},
	})
}

// Helpers for slice-to-DTO conversion

func toExperiencesDTOs(exps []domain.WorkExperience) []experienceDTO {
	out := make([]experienceDTO, 0, len(exps))
	for _, e := range exps {
		out = append(out, experienceDTO{
			ID:             e.ID,
			Title:          e.Position,
			Company:        e.Company,
			Location:       e.Location,
			EmploymentType: e.EmploymentType,
			StartDate:      e.StartDate.Format("2006-01-02"),
			EndDate:        e.EndDate.Format("2006-01-02"),
			IsCurrent:      e.IsCurrent,
			Description:    e.Responsibilities,
			Achievements:   e.Achievements,
		})
	}
	return out
}

func toEducationsDTOs(edus []domain.Education) []educationDTO {
	out := make([]educationDTO, 0, len(edus))
	for _, e := range edus {
		out = append(out, educationDTO{
			ID:           e.ID,
			School:       e.Institution,
			Degree:       e.Degree,
			FieldOfStudy: e.FieldOfStudy,
			StartDate:    "",
			EndDate:      e.GraduationDate.Format("2006-01-02"),
			Grade:        "",
			Description:  e.Thesis,
		})
	}
	return out
}

func toCertificationsDTOs(certs []domain.CertificationItem) []certificationDTO {
	out := make([]certificationDTO, 0, len(certs))
	for _, c := range certs {
		out = append(out, certificationDTO{
			ID:            c.ID,
			Name:          c.Name,
			Issuer:        c.Issuer,
			IssueDate:     c.IssueDate.Format("2006-01-02"),
			ExpiryDate:    c.ExpirationDate.Format("2006-01-02"),
			CredentialID:  c.CredentialID,
			CredentialURL: c.VerificationURL,
		})
	}
	return out
}

func toSkillsDTOs(skills []domain.SkillItem) []skillDTO {
	out := make([]skillDTO, 0, len(skills))
	for _, s := range skills {
		out = append(out, skillDTO{
			Name:             s.Name,
			ProficiencyLevel: s.Level,
		})
	}
	return out
}

func toProjectsDTOs(projs []domain.ProjectItem) []portfolioLinkDTO {
	out := make([]portfolioLinkDTO, 0, len(projs))
	for _, p := range projs {
		out = append(out, portfolioLinkDTO{
			ID:       p.ID,
			Platform: p.Title,
			URL:      p.LiveDemoURL,
		})
	}
	return out
}
