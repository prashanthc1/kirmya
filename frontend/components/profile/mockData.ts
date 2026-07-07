import { ExtendedProfile } from "./types";

export const MOCK_PROFILE: ExtendedProfile = {
  user_id: "user_marcus_hale_102",
  headline: "Operations Director · Supply Chain & Logistics Network Architect",
  about:
    "Operations leader with 22 years steadying complex supply chains through growth, restructuring, and two downturns. I've owned cost bases north of $400M, rebuilt distribution networks under pressure, and kept service levels high when budgets were not. I lead calmly, hire well, and make the unglamorous call when it's the right one.",
  photo_url: "", // Will generate a visual initial or use a beautiful SVG placeholder
  bio: "Operations Director & Supply Chain Executive building resilient networks, optimizing distribution channels, and coaching high-performing logistics teams.",
  location: "Denver, CO",
  website: "https://marcushale.dev",
  version: 1,

  // Core Identity
  pronouns: "He/Him",
  career_status: "actively_looking",
  transition_reason: "layoff",
  target_comeback_timeline: "Immediate (Next 30 days)",
  supports_needed: ["referrals", "interview_prep", "skill_building"],

  // Mobility
  open_to_remote: true,
  open_to_relocation: true,
  relocation_locations: [
    "Austin, TX",
    "Seattle, WA",
    "Chicago, IL",
    "Dallas, TX",
  ],
  desired_roles: [
    "Director of Operations",
    "VP of Supply Chain",
    "Director of Logistics",
    "Global Operations Lead",
  ],
  desired_industries: [
    "Logistics & Supply Chain",
    "E-commerce & Retail",
    "SaaS & Enterprise Technology",
    "Manufacturing",
  ],
  employment_type: "full_time",
  salary_min: 185000,
  salary_max: 230000,
  salary_currency: "USD",
  salary_visible: false,
  work_mode: "hybrid",
  availability_date: "2026-07-15",
  notice_period: "Immediate",

  // Verification & Trust
  referral_eligible: true,
  email_verified: true,
  phone_verified: true,
  linkedin_verified: true,
  id_verified: true,
  email: "marcus.hale@kirmya.net",
  phone: "+1 (303) 555-0148",

  // AI Coach & Narrative
  career_narrative:
    "Marcus is a veteran supply chain executive who recently parted ways with Cascade Freight following a corporate restructuring. Rather than treating this as a setback, he is targeting a high-impact role in high-growth e-commerce where his 22 years of scaling networks and consolidating operations can yield significant margins.",
  coaching_metadata: JSON.stringify({
    last_assessment: "2026-07-06",
    focus_areas: [
      "ATS Keyword alignment for tech supply chains",
      "Narrative repositioning for remote leadership",
    ],
  }),

  // Work Auth
  work_auth_status: "citizen",
  passport_nationality: "United States",
  driving_license_bool: true,
  driving_license_type: "Class A Commercial & Standard",

  // Communication & Accessibility
  preferred_contact_channel: "in_app",
  accessibility_needs: "",
  video_intro_url: "https://assets.kirmya.ai/videos/marcus_intro.mp4",

  // Mentorship
  willing_to_mentor: true,

  // Calculated Fields
  avg_response_time_hours: 1.4,
  profile_completeness_score: 88,
  last_active_at: new Date().toISOString(),
  background_check_consent: true,
  background_check_consent_at: "2026-06-15T10:00:00Z",

  // Alerts
  job_alert_frequency: "daily",
  job_alert_channel: "email",

  // Privacy Visibility
  visibility_profile: "public",
  visibility_salary: "recruiters_only",
  visibility_transition_reason: "recruiters_only",
  visibility_experience: "public",
  visibility_education: "public",
  visibility_certifications: "public",
  visibility_skills: "public",
  visibility_portfolio: "public",
  visibility_references: "recruiters_only",

  // Social URLs
  preferred_name: "Marcus",
  cover_banner:
    "https://images.unsplash.com/photo-1578575437130-527eed3abbec?w=1200&auto=format&fit=crop&q=80",
  linkedin_url: "https://linkedin.com/in/marcushale-ops",
  github_url: "https://github.com/mhale-supplychain",
  behance_url: "",
  dribbble_url: "",
  medium_url: "https://medium.com/@marcushale-logistics",
  stackoverflow_url: "",
  googlescholar_url: "",
  researchgate_url: "",
  orcid_id: "",
  calendar_link: "https://cal.com/marcushale-ops/30min",
  qr_code_url: "",
  visa_status: "N/A - US Citizen",

  // Section 2: Summary metrics & AI versions
  years_experience: 22,
  industry: "Logistics & Supply Chain",
  career_goals:
    "To lead global operations for an enterprise logistics or high-growth e-commerce organization, applying network optimization models to cut unit costs by 15-20% while building resilient, human-first operations culture.",
  interests: [
    "Automation & Robotics",
    "Operations Research",
    "Mentorship",
    "Sustainability in Shipping",
    "Mountain Hiking",
  ],
  strengths: [
    "P&L Management",
    "Strategic Network Planning",
    "Warehouse Automation",
    "S&OP Modeling",
    "Leadership & Talent Development",
  ],
  personal_brand: "The Unshakeable Network Optimizer",
  elevator_pitch:
    "I scale logistics networks from chaotic bottlenecks into predictable margin-drivers. Over the last 20 years, I've consolidated supply chains, introduced automation that shaved 18% off unit costs, and kept teams aligned through intense organizational restructuring. I make networks run smoothly so business strategy can execute cleanly.",

  // Section 3: Work Experience
  experiences: [
    {
      id: "exp_1",
      title: "Operations Director",
      company: "Cascade Freight",
      location: "Denver, CO",
      employment_type: "Full-Time",
      location_type: "hybrid",
      start_date: "2014-04",
      end_date: "2025-05",
      is_current: false,
      description:
        "Ran a 120-person operations organization across 14 distribution centers handling regional freight operations and mid-mile consolidation. Led the network consolidation that reduced overall cost-per-unit by 18% while lifting delivery consistency to 99.2%. Managed a $42M annual budget.",
      achievements: [
        "Consolidated 14 distribution facilities into 8 regional hubs, saving $6.4M in annual lease and labor overhead.",
        "Introduced automated sorting conveyor technology that increased throughput capacity by 35% during Q4 peak seasons.",
        "Maintained an industry-leading 99.2% on-time delivery metric over a 3-year running period.",
      ],
      kpis: [
        "18% reduction in cost-per-unit",
        "99.2% On-Time Delivery Rate",
        "$6.4M annual overhead savings",
      ],
      technologies: [
        "SAP Extended Warehouse Management (EWM)",
        "Tableau",
        "MercuryGate TMS",
        "Python (Pulp Linear Programming)",
      ],
      skills_used: [
        "Network Strategy",
        "S&OP Planning",
        "P&L Ownership",
        "Warehouse Automation",
        "Labor Management Systems",
      ],
      attachments: ["freight_consolidation_case_study.pdf"],
      recommendations: [
        "'Marcus is the calmest and most analytical operations leader I've worked with. He turned our supply chain from a constant fire drill into a competitive advantage.' - Sarah Jenkins, VP of Logistics",
      ],
    },
    {
      id: "exp_2",
      title: "Senior Operations Manager",
      company: "Northstar Distribution",
      location: "Salt Lake City, UT",
      employment_type: "Full-Time",
      location_type: "onsite",
      start_date: "2007-08",
      end_date: "2014-03",
      is_current: false,
      description:
        "Scaled regional fulfillment from 2 to 6 facilities through a high-growth period, establishing the standard S&OP pipeline still in use today. Directed facility layout design, labor scheduling, and carrier negotiations for Western US region.",
      achievements: [
        "Designed and launched 4 new fulfillment centers on-time and under budget, totaling 1.2M sq ft of storage.",
        "Renegotiated LTL (Less-Than-Truckload) carrier contracts, yielding an immediate 12% margin improvement on outbound shipping.",
        "Spearheaded safety compliance program that reduced OSHA reportable incidents by 42% in 18 months.",
      ],
      kpis: [
        "12% outbound shipping cost reduction",
        "42% decrease in OSHA incidents",
        "1.2M sq ft expansion executed",
      ],
      technologies: [
        "Manhattan Associates WMS",
        "Excel (VBA modeling)",
        "Oracle ERP",
      ],
      skills_used: [
        "Carrier Negotiation",
        "Facility Design",
        "Contract Negotiation",
        "Safety Compliance",
        "OSHA Regulations",
      ],
      attachments: [],
      recommendations: [],
    },
    {
      id: "exp_3",
      title: "Operations Manager",
      company: "Meridian Logistics",
      location: "Phoenix, AZ",
      employment_type: "Full-Time",
      location_type: "onsite",
      start_date: "2003-01",
      end_date: "2007-07",
      is_current: false,
      description:
        "Began on the warehouse floor and worked up to managing a regional operations hub of 60 active fulfillment staff. Directed shipping/receiving, inventory auditing, and shift supervisor teams.",
      achievements: [
        "Promoted twice within 2 years based on outstanding performance in inventory accuracy.",
        "Achieved 99.98% inventory accuracy over 12 consecutive months by implementing localized cycle-counting protocols.",
        "Coached and promoted 5 warehouse associates into supervisory positions.",
      ],
      kpis: [
        "99.98% inventory accuracy rate",
        "2 promotions within 24 months",
        "5 staff members promoted to leadership",
      ],
      technologies: ["RedPrairie WMS", "Radio Frequency (RF) scanners"],
      skills_used: [
        "Inventory Auditing",
        "Cycle Counting",
        "Team Coaching",
        "Shift Management",
      ],
      attachments: [],
      recommendations: [],
    },
  ],

  // Section 4: Education
  educations: [
    {
      id: "edu_1",
      school: "Colorado State University",
      degree: "Master of Science",
      field_of_study: "Supply Chain Management & Operations Research",
      start_date: "2000-09",
      end_date: "2002-06",
      grade: "3.85 GPA",
      description:
        "Specialized in quantitative analysis, linear programming optimization, and network graph theories. Graduate Teaching Assistant in Stochastic Models.",
      major: "Operations Research",
      gpa: "3.85",
      honors: "Summa Cum Laude, Supply Chain Scholarship recipient",
      activities: "APICS Student Chapter President, Operations Research Club",
      research:
        "Co-authored paper on dynamic routing algorithms under variable demand constraints.",
      projects:
        "Heuristic solver for vehicle routing problems using custom genetic algorithms.",
      thesis:
        "Dynamic Network Allocation: Simulating Facility Location Under Volatile Carrier Capacity",
    },
    {
      id: "edu_2",
      school: "University of Colorado Boulder",
      degree: "Bachelor of Science",
      field_of_study: "Business Administration (Operations Management)",
      start_date: "1996-09",
      end_date: "2000-06",
      grade: "3.62 GPA",
      description:
        "Foundation in finance, business logistics, statistics, and organizational behavior.",
      major: "Operations Management",
      gpa: "3.62",
      honors: "Dean's List (6 semesters)",
      activities: "Intramural Soccer, Business Council",
      research: "",
      projects: "",
      thesis: "",
    },
  ],

  // Section 5: Skills
  skills: [
    {
      name: "Network Strategy",
      proficiency_level: "Expert",
      endorsed_count: 14,
      years_experience: 15,
      last_used: "2025",
      verification_status: "verified",
      recruiter_demand: "high",
      market_trend: "rising",
    },
    {
      name: "S&OP Planning",
      proficiency_level: "Expert",
      endorsed_count: 18,
      years_experience: 12,
      last_used: "2025",
      verification_status: "verified",
      recruiter_demand: "high",
      market_trend: "stable",
    },
    {
      name: "P&L Ownership",
      proficiency_level: "Expert",
      endorsed_count: 11,
      years_experience: 11,
      last_used: "2025",
      verification_status: "verified",
      recruiter_demand: "high",
      market_trend: "stable",
    },
    {
      name: "Carrier Negotiation",
      proficiency_level: "Expert",
      endorsed_count: 9,
      years_experience: 14,
      last_used: "2024",
      verification_status: "verified",
      recruiter_demand: "medium",
      market_trend: "stable",
    },
    {
      name: "Cost Reduction",
      proficiency_level: "Expert",
      endorsed_count: 22,
      years_experience: 18,
      last_used: "2025",
      verification_status: "verified",
      recruiter_demand: "high",
      market_trend: "rising",
    },
    {
      name: "Warehouse Automation",
      proficiency_level: "Intermediate",
      endorsed_count: 6,
      years_experience: 6,
      last_used: "2025",
      verification_status: "pending",
      recruiter_demand: "high",
      market_trend: "rising",
    },
    {
      name: "Fulfillment Logistics",
      proficiency_level: "Expert",
      endorsed_count: 15,
      years_experience: 20,
      last_used: "2025",
      verification_status: "verified",
      recruiter_demand: "high",
      market_trend: "stable",
    },
    {
      name: "Linear Programming",
      proficiency_level: "Intermediate",
      endorsed_count: 4,
      years_experience: 8,
      last_used: "2025",
      verification_status: "unverified",
      recruiter_demand: "medium",
      market_trend: "rising",
    },
    {
      name: "SAP Extended Warehouse Management",
      proficiency_level: "Intermediate",
      endorsed_count: 5,
      years_experience: 7,
      last_used: "2025",
      verification_status: "verified",
      recruiter_demand: "medium",
      market_trend: "stable",
    },
    {
      name: "Labor Management Systems",
      proficiency_level: "Advanced",
      endorsed_count: 7,
      years_experience: 10,
      last_used: "2025",
      verification_status: "verified",
      recruiter_demand: "medium",
      market_trend: "stable",
    },
  ],

  // Section 6: Projects
  projects: [
    {
      id: "proj_1",
      title: "Cascade Freight Hub Consolidation Initiative",
      description:
        "A complete overhaul and mathematical consolidation of Marcus's distribution network. We simulated cost models, renegotiated leases, and moved 120 FTEs to regional hubs with 0 shipping outages.",
      cover_image:
        "https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?w=800&auto=format&fit=crop&q=60",
      repository_url: "https://github.com/mhale-supplychain/cascade-hub-solver",
      live_demo_url: "https://cascade-consolidator.kirmya.app",
      video_url: "",
      screenshots: [
        "https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?w=800&auto=format&fit=crop&q=60",
        "https://images.unsplash.com/photo-1553413719-875871274712?w=800&auto=format&fit=crop&q=60",
      ],
      technologies: [
        "Python",
        "PuLP Solver",
        "Tableau Dashboard API",
        "AWS ECS",
      ],
      timeline: "9 Months (2022 - 2023)",
      team_size: 14,
      metrics: "$6.4M annual overhead reduced, 18% unit cost reduction",
      awards: "Cascade Board Excellence Award 2023",
    },
    {
      id: "proj_2",
      title: "Automated Conveyor & Sortation Integration",
      description:
        "Spearheaded design and vendor procurement for a multi-tier conveyor sorter. Unified conveyor PLC data with WMS to feed sorting algorithms that optimize carrier assignment at the barcode scanner stage.",
      cover_image:
        "https://images.unsplash.com/photo-1563986768609-322da13575f3?w=800&auto=format&fit=crop&q=60",
      repository_url: "",
      live_demo_url: "",
      video_url: "",
      screenshots: [],
      technologies: [
        "RedPrairie PLC integrations",
        "Cognex Barcode Systems",
        "Java Spring Boot",
      ],
      timeline: "12 Months (2020)",
      team_size: 8,
      metrics: "35% increase in throughput capacity, -40% sorting errors",
      awards: "Supply Chain Innovator Honoree 2021",
    },
  ],

  // Section 7: Certifications
  certifications: [
    {
      id: "cert_1",
      name: "Certified Supply Chain Professional (CSCP)",
      issuer: "ASCM (Association for Supply Chain Management)",
      issue_date: "2006-03-12",
      expiry_date: "2027-03-12",
      credential_id: "CSCP-104928-MH",
      credential_url: "https://ascm.org/verify/CSCP-104928-MH",
      skills_covered: [
        "S&OP Planning",
        "Network Strategy",
        "Carrier Negotiation",
      ],
    },
    {
      id: "cert_2",
      name: "Lean Six Sigma Black Belt",
      issuer: "Institute of Industrial and Systems Engineers",
      issue_date: "2010-11-20",
      expiry_date: "",
      credential_id: "LSSBB-93821",
      credential_url: "https://iise.org/verify/blackbelt/93821",
      skills_covered: [
        "Cost Reduction",
        "Fulfillment Logistics",
        "Cycle Counting",
      ],
    },
  ],

  // Section 8: Achievements
  achievements_list: [
    {
      id: "ach_1",
      title: "National Supply Chain Leadership Award",
      issuer_or_org: "Logistics Weekly Association",
      date: "2023-10-14",
      category: "award",
      description:
        "Awarded to 3 national directors who displayed exemplary network adaptability and team coaching during logistics carrier strikes.",
      url: "https://logisticsweekly.com/awards-2023-winners",
    },
    {
      id: "ach_2",
      title: "Consolidated Warehouse Network Placement Solver",
      issuer_or_org: "US Patent and Trademark Office",
      date: "2018-05-22",
      category: "patent",
      description:
        "Patent US-93821-B2: System and algorithm for placement of hubs in networks under volatile fuel surcharge vectors.",
      url: "https://patents.google.com/patent/US93821B2",
    },
    {
      id: "ach_3",
      title: "Keynote Speaker: Building Resilient Distribution Networks",
      issuer_or_org: "Rocky Mountain Logistics Symposium",
      date: "2024-02-18",
      category: "conference",
      description:
        "Delivered keynote speech on the realities of distribution center consolidation to an audience of 600 supply chain leaders.",
      url: "",
    },
  ],

  // Section 9: Resumes & Documents
  resumes: [
    {
      id: "res_1",
      name: "Marcus_Hale_Operations_Director_2026.pdf",
      uploaded_at: "2026-06-15T09:30:00Z",
      file_size: "245 KB",
      is_primary: true,
      ats_score: 94,
      readability_score: 88,
      file_url: "#",
    },
    {
      id: "res_2",
      name: "Marcus_Hale_Technical_Supply_Chain_CV.pdf",
      uploaded_at: "2026-06-20T14:15:00Z",
      file_size: "260 KB",
      is_primary: false,
      ats_score: 89,
      readability_score: 92,
      file_url: "#",
    },
  ],
  cover_letters: [
    {
      id: "cl_1",
      title: "Standard Executive Cover Letter",
      recipient_company: "General Supply Partners",
      recipient_role: "VP of Global Logistics",
      content:
        "Dear Hiring Team,\n\nI am writing to express my interest in the VP of Global Logistics role. With 22 years of scaling warehouse capacities, negotiating carrier contracts, and consolidating regional networks, I have consistently driven double-digit cost improvements. At Cascade Freight, I consolidated our footprint from 14 to 8 facilities, reducing overall cost-per-unit by 18% while lifting delivery reliability to 99.2%.\n\nI look forward to discussing how I can bring this same standard of operational excellence to General Supply Partners.\n\nSincerely,\nMarcus Hale",
      last_modified: "2026-07-01T11:00:00Z",
    },
  ],

  // Section 11: Verification Detail
  verification_status: {
    email: "verified",
    phone: "verified",
    gov_id: "verified",
    employment: "verified",
    education: "verified",
    skills: "verified",
    certifications: "verified",
  },
  trust_score: 95,
  recruiter_badge: true,

  // Section 12: Network Contacts
  network: [
    {
      id: "conn_1",
      name: "Sarah Jenkins",
      headline: "VP of Global Logistics at Cascade Freight",
      type: "mentor",
      company: "Cascade Freight",
      is_verified: true,
    },
    {
      id: "conn_2",
      name: "David Kim",
      headline: "Principal Operations Researcher at Amazon",
      type: "connection",
      company: "Amazon",
      is_verified: true,
    },
    {
      id: "conn_3",
      name: "Emily Watson",
      headline: "Executive Talent Acquisition Lead",
      type: "recruiter",
      company: "Apex Supply Recruiting",
      is_verified: true,
    },
    {
      id: "conn_4",
      name: "Robert Chen",
      headline: "Director of Warehouse Automation",
      type: "following",
      company: "Symbotic",
      is_verified: false,
    },
  ],

  // Section 13: Analytics Dashboard
  analytics: {
    profile_views: 1420,
    recruiter_searches: 384,
    ats_score: 94,
    resume_downloads: 122,
    portfolio_views: 405,
    search_ranking: 12,
    keyword_ranking: [
      "Supply Chain Strategy",
      "Warehouse Automation",
      "S&OP Planning",
      "Cost Reduction",
      "Network Design",
    ],
    interview_rate: 18.5,
    referral_rate: 64.0,
    application_success: 78.5,
  },

  // Section 14: Privacy Config
  privacy_settings: {
    profile_visibility: "public",
    anonymous_mode: false,
    hide_salary: true,
    hide_employer: false,
    search_indexing: true,
    blocked_users: ["user_block_1", "user_block_2"],
    mfa_enabled: true,
    active_sessions: [
      {
        device: "MacBook Pro · macOS Chrome",
        location: "Denver, CO",
        last_active: "Active Now",
      },
      {
        device: "iPhone 15 Pro · iOS Safari",
        location: "Denver, CO",
        last_active: "2 hours ago",
      },
    ],
    api_tokens: [
      {
        id: "tok_1",
        name: "Resume Parser CLI Link",
        created_at: "2026-05-10",
        last_used: "2026-07-05",
      },
    ],
  },
  languages: [
    { name: "English", proficiency: "Native or Bilingual" },
    { name: "Spanish", proficiency: "Professional Working" },
  ],
  portfolio: [
    { platform: "LinkedIn", url: "https://linkedin.com/in/marcushale-ops" },
    { platform: "GitHub", url: "https://github.com/mhale-supplychain" },
    {
      platform: "Personal Blog",
      url: "https://medium.com/@marcushale-logistics",
    },
  ],
};
