import { api } from "./client";

export interface Profile {
  user_id: string;
  headline: string;
  about: string;
  photo_url: string;
  bio: string;
  location: string;
  website: string;
  version: number;

  // Core Identity
  pronouns: string;
  career_status: "actively_looking" | "open_to_opportunities" | "employed_exploring" | "career_break" | "";

  // Career Recovery
  transition_reason?: "layoff" | "sabbatical" | "caregiving" | "health" | "upskilling" | "relocation" | "other" | "";
  target_comeback_timeline: string;
  supports_needed: ("referrals" | "mentorship" | "interview_prep" | "skill_building")[];

  // Mobility & Preferences
  open_to_remote: boolean;
  open_to_relocation: boolean;
  relocation_locations: string[];
  desired_roles: string[];
  desired_industries: string[];
  employment_type: "full_time" | "contract" | "freelance" | "part_time" | "";
  salary_min?: number;
  salary_max?: number;
  salary_currency?: string;
  salary_visible: boolean;
  work_mode: "remote" | "hybrid" | "onsite" | "";
  availability_date: string;
  notice_period: string;

  // Trust & Verification
  referral_eligible: boolean;
  email_verified: boolean;
  phone_verified: boolean;
  linkedin_verified: boolean;
  id_verified: boolean;

  // AI Coach
  career_narrative: string;
  coaching_metadata: string;

  // Work Auth
  work_auth_status: "citizen" | "resident" | "needs_sponsorship" | "visit_visa" | "";
  passport_nationality: string;
  driving_license_bool: boolean;
  driving_license_type: string;

  // Communication & Accessibility
  preferred_contact_channel: "email" | "whatsapp" | "in_app" | "";
  accessibility_needs?: string;
  video_intro_url: string;

  // Mentorship
  willing_to_mentor: boolean;

  // Calculated Fields
  avg_response_time_hours: number;
  profile_completeness_score: number;
  last_active_at: string;

  // Consent
  background_check_consent: boolean;
  background_check_consent_at: string;

  // Alerts
  job_alert_frequency: "instant" | "daily" | "weekly" | "";
  job_alert_channel: "email" | "whatsapp" | "push" | "";

  // Privacy Visibility
  visibility_profile: "public" | "recruiters_only" | "mentors_only" | "private";
  visibility_salary: "public" | "recruiters_only" | "mentors_only" | "private";
  visibility_transition_reason: "public" | "recruiters_only" | "mentors_only" | "private";
  visibility_experience: "public" | "recruiters_only" | "mentors_only" | "private";
  visibility_education: "public" | "recruiters_only" | "mentors_only" | "private";
  visibility_certifications: "public" | "recruiters_only" | "mentors_only" | "private";
  visibility_skills: "public" | "recruiters_only" | "mentors_only" | "private";
  visibility_portfolio: "public" | "recruiters_only" | "mentors_only" | "private";
  visibility_references: "public" | "recruiters_only" | "mentors_only" | "private";

  // Nested Collections
  experiences: WorkExperience[];
  educations: Education[];
  certifications: Certification[];
  skills: ProfileSkill[];
  languages: Language[];
  portfolio: PortfolioLink[];
  endorsements?: Endorsement[];
  references?: Reference[];
}

export interface WorkExperience {
  id?: string;
  title: string;
  company: string;
  location: string;
  employment_type: string;
  start_date: string;
  end_date: string;
  is_current: boolean;
  description: string;
  achievements: string[];
}

export interface Education {
  id?: string;
  school: string;
  degree: string;
  field_of_study: string;
  start_date: string;
  end_date: string;
  grade: string;
  description: string;
}

export interface Certification {
  id?: string;
  name: string;
  issuer: string;
  issue_date: string;
  expiry_date: string;
  credential_id: string;
  credential_url: string;
}

export interface ProfileSkill {
  name: string;
  proficiency_level: string;
  endorsed_count: number;
}

export interface Language {
  name: string;
  proficiency: string;
}

export interface PortfolioLink {
  id?: string;
  platform: string;
  url: string;
}

export interface Endorsement {
  id?: string;
  to_user_id: string;
  from_user_id: string;
  relationship: string;
  text: string;
  created_at?: string;
}

export interface Reference {
  id?: string;
  name: string;
  relationship: string;
  contact_info: string;
  permission_to_contact: boolean;
}

export interface ConsentLog {
  id?: string;
  consent_type: "background_check" | "data_sharing";
  target_entity: string;
  consented: boolean;
  ip_address?: string;
  user_agent?: string;
  created_at?: string;
}

// API functions
export const profileClient = {
  getMe: () => api.get<Profile>("/profiles/me"),
  
  updateMe: (data: Partial<Profile>) => api.put<Profile>("/profiles/me", data),
  
  getByID: (id: string) => api.get<Profile>(`/profiles/${id}`),

  // Experiences
  addExperience: (exp: WorkExperience) => api.post<Profile>("/profiles/me/experiences", exp),
  updateExperience: (id: string, exp: WorkExperience) => api.put<Profile>(`/profiles/me/experiences/${id}`, exp),
  deleteExperience: (id: string) => api.delete<Profile>(`/profiles/me/experiences/${id}`),

  // Education
  addEducation: (edu: Education) => api.post<Profile>("/profiles/me/educations", edu),
  updateEducation: (id: string, edu: Education) => api.put<Profile>(`/profiles/me/educations/${id}`, edu),
  deleteEducation: (id: string) => api.delete<Profile>(`/profiles/me/educations/${id}`),

  // Certifications
  addCertification: (cert: Certification) => api.post<Profile>("/profiles/me/certifications", cert),
  updateCertification: (id: string, cert: Certification) => api.put<Profile>(`/profiles/me/certifications/${id}`, cert),
  deleteCertification: (id: string) => api.delete<Profile>(`/profiles/me/certifications/${id}`),

  // Skills
  setSkills: (skills: ProfileSkill[]) => api.put<Profile>("/profiles/me/skills", { skills }),

  // Languages
  setLanguages: (languages: Language[]) => api.put<Profile>("/profiles/me/languages", { languages }),

  // Portfolio
  setPortfolio: (portfolio: PortfolioLink[]) => api.put<Profile>("/profiles/me/portfolio", { portfolio }),

  // Endorsements
  addEndorsement: (endorsement: Omit<Endorsement, "id" | "from_user_id">) => 
    api.post<Profile>("/profiles/me/endorsements", endorsement),

  // References
  addReference: (ref: Reference) => api.post<Profile>("/profiles/me/references", ref),
  updateReference: (id: string, ref: Reference) => api.put<Profile>(`/profiles/me/references/${id}`, ref),
  deleteReference: (id: string) => api.delete<Profile>(`/profiles/me/references/${id}`),

  // Consent log
  addConsent: (consent: Omit<ConsentLog, "id">) => api.post<void>("/profiles/me/consent", consent),
};
