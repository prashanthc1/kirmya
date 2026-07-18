import {
  Profile as BaseProfile,
  WorkExperience,
  Education,
  Certification,
  ProfileSkill,
  Language,
  PortfolioLink,
  Reference,
} from "@/lib/api/profile";

export type {
  WorkExperience,
  Education,
  Certification,
  ProfileSkill,
  Language,
  PortfolioLink,
  Reference,
};

export interface ProjectItem {
  id: string;
  title: string;
  description: string;
  cover_image?: string;
  repository_url?: string;
  live_demo_url?: string;
  video_url?: string;
  screenshots?: string[];
  technologies: string[];
  timeline: string;
  team_size: number;
  metrics?: string;
  awards?: string;
}

export interface AchievementItem {
  id: string;
  title: string;
  issuer_or_org: string;
  date: string;
  category:
    | "award"
    | "publication"
    | "patent"
    | "conference"
    | "hackathon"
    | "open_source"
    | "volunteer"
    | "competition"
    | "leadership";
  description: string;
  url?: string;
}

export interface ResumeVersion {
  id: string;
  name: string;
  uploaded_at: string;
  file_size: string;
  is_primary: boolean;
  ats_score: number;
  readability_score: number;
  file_url: string;
}

export interface CoverLetter {
  id: string;
  title: string;
  recipient_company: string;
  recipient_role: string;
  content: string;
  last_modified: string;
}

export interface NetworkConnection {
  id: string;
  name: string;
  avatar_url?: string;
  headline: string;
  type: "connection" | "follower" | "following" | "mentor" | "recruiter";
  company?: string;
  is_verified?: boolean;
}

export interface ProfileAnalytics {
  profile_views: number;
  recruiter_searches: number;
  ats_score: number;
  resume_downloads: number;
  portfolio_views: number;
  search_ranking: number;
  keyword_ranking: string[];
  interview_rate: number;
  referral_rate: number;
  application_success: number;
}

export interface ExtendedProfile extends BaseProfile {
  // Section 1: Additional Identity fields
  full_name?: string;
  preferred_name?: string;
  cover_banner?: string;
  nationality?: string;
  timezone?: string;
  phone?: string;
  email?: string;
  linkedin_url?: string;
  github_url?: string;
  behance_url?: string;
  dribbble_url?: string;
  medium_url?: string;
  stackoverflow_url?: string;
  googlescholar_url?: string;
  researchgate_url?: string;
  orcid_id?: string;
  calendar_link?: string;
  qr_code_url?: string;
  visa_status?: string;

  // Section 2: Professional Summary details
  years_experience?: number;
  industry?: string;
  career_goals?: string;
  interests?: string[];
  strengths?: string[];
  personal_brand?: string;
  elevator_pitch?: string;

  // Section 3: Work Experience with details
  experiences: (WorkExperience & {
    id: string;
    employment_type?: string;
    location_type?: "remote" | "hybrid" | "onsite";
    achievements?: string[];
    kpis?: string[];
    technologies?: string[];
    skills_used?: string[];
    attachments?: string[];
    recommendations?: string[];
  })[];

  // Section 4: Education with details
  educations: (Education & {
    id: string;
    major?: string;
    gpa?: string;
    honors?: string;
    activities?: string;
    research?: string;
    projects?: string;
    thesis?: string;
  })[];

  // Section 5: Skills with metrics
  skills: (ProfileSkill & {
    years_experience?: number;
    last_used?: string;
    verification_status?: "verified" | "pending" | "unverified";
    recruiter_demand?: "low" | "medium" | "high";
    market_trend?: "rising" | "stable" | "declining";
  })[];

  // Section 6: Projects
  projects: ProjectItem[];

  // Section 7: Certifications & Licenses
  certifications: (Certification & {
    id: string;
    skills_covered?: string[];
  })[];

  // Section 8: Achievements
  achievements_list: AchievementItem[];

  // Section 9: Resumes & Documents
  resumes: ResumeVersion[];
  cover_letters: CoverLetter[];

  // Section 10: Career Preferences
  // Section 11: Verification & Trust Statuses
  verification_status?: {
    email: "verified" | "unverified";
    phone: "verified" | "unverified";
    gov_id: "verified" | "pending" | "unverified";
    employment: "verified" | "pending" | "unverified";
    education: "verified" | "pending" | "unverified";
    skills: "verified" | "pending" | "unverified";
    certifications: "verified" | "pending" | "unverified";
  };
  trust_score?: number; // 0-100
  recruiter_badge?: boolean;

  // Section 12: Networking List
  network: NetworkConnection[];

  // Section 13: Analytics
  analytics?: ProfileAnalytics;

  // Section 14: Privacy & Security
  privacy_settings?: {
    profile_visibility:
      "public" | "recruiters_only" | "connections_only" | "private";
    anonymous_mode: boolean;
    hide_salary: boolean;
    hide_employer: boolean;
    search_indexing: boolean;
    blocked_users: string[];
    mfa_enabled: boolean;
    active_sessions: {
      device: string;
      location: string;
      last_active: string;
    }[];
    api_tokens: {
      id: string;
      name: string;
      created_at: string;
      last_used: string;
    }[];
  };
}
