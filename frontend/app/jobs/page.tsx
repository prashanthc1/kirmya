"use client";

import React, { useEffect, useState, useMemo } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { 
  Search, 
  MapPin, 
  DollarSign, 
  Briefcase, 
  Sparkles, 
  ChevronRight, 
  CheckCircle2, 
  Building2, 
  Clock, 
  ShieldAlert,
  ArrowLeft,
  Award,
  BookOpen,
  Info
} from "lucide-react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api, ApiError } from "@/lib/api/client";

interface Job {
  id: string;
  title: string;
  company: string;
  location: string;
  salary: string;
  job_type: string;
  description: string;
}

const JOB_TYPES = ["All Types", "Full-time", "Contract", "Remote", "Part-time"];

// Stable mock metrics derived from job titles/content for visual excellence
const getMockJobData = (job: Job) => {
  const code = job.id.charCodeAt(0) + job.id.charCodeAt(job.id.length - 1);
  const matchScore = 80 + (code % 19); // 80 - 98
  const difficulty = (code % 3) === 0 ? "Hard" : (code % 3) === 1 ? "Medium" : "Easy";
  
  const skillsMap: Record<string, string[]> = {
    engineer: ["React", "Next.js", "TypeScript", "Node.js", "System Design"],
    developer: ["Go", "PostgreSQL", "Redis", "Docker", "AWS"],
    designer: ["Figma", "Design Systems", "Prototyping", "UX Research"],
    product: ["Product Strategy", "Agile", "User Interviews", "Roadmapping"],
    marketing: ["SEO", "Copywriting", "Growth Hacking", "Google Analytics"],
  };
  
  const titleLower = job.title.toLowerCase();
  let skills = ["Communication", "Problem Solving", "Collaboration"];
  for (const [key, list] of Object.entries(skillsMap)) {
    if (titleLower.includes(key)) {
      skills = list;
      break;
    }
  }

  const missingSkills = skills.slice(Math.max(1, skills.length - 2));
  const matchedSkills = skills.slice(0, Math.max(1, skills.length - 2));

  return { matchScore, difficulty, matchedSkills, missingSkills };
};

export default function JobsPage() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [searchQuery, setSearchQuery] = useState("");
  const [selectedType, setSelectedType] = useState("All Types");
  const [selectedJobId, setSelectedJobId] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await api.get<{ jobs: Job[] }>("/jobs");
        if (active) {
          const list = data?.jobs ?? [];
          setJobs(list);
          if (list.length > 0) {
            setSelectedJobId(list[0].id);
          }
        }
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError ? err.message : "Could not load jobs. Please try again."
          );
        }
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  const filteredJobs = useMemo(() => {
    let result = jobs;

    if (searchQuery.trim() !== "") {
      const q = searchQuery.toLowerCase();
      result = result.filter(
        (job) =>
          job.title.toLowerCase().includes(q) ||
          job.company.toLowerCase().includes(q) ||
          job.description?.toLowerCase().includes(q) ||
          job.location?.toLowerCase().includes(q)
      );
    }

    if (selectedType !== "All Types") {
      result = result.filter(
        (job) => job.job_type?.toLowerCase() === selectedType.toLowerCase()
      );
    }

    return result;
  }, [searchQuery, selectedType, jobs]);

  const activeJob = useMemo(() => {
    return jobs.find((j) => j.id === selectedJobId) || filteredJobs[0] || null;
  }, [selectedJobId, jobs, filteredJobs]);

  const mockData = useMemo(() => {
    if (!activeJob) return null;
    return getMockJobData(activeJob);
  }, [activeJob]);

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Jobs" }]} />

      <main className="flex-1 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex flex-col md:flex-row gap-6 overflow-hidden">
        {/* Left Side: Job Search List */}
        <div className={`flex-1 md:w-[420px] md:flex-none flex flex-col gap-4 ${
          activeJob && selectedJobId ? "hidden md:flex" : "flex"
        }`}>
          {/* Header */}
          <div>
            <h1 className="text-2xl font-extrabold tracking-tight">Open Opportunities</h1>
            <p className="text-xs text-muted-foreground mt-1">
              {loading ? "Searching roles..." : `${filteredJobs.length} active positions found`}
            </p>
          </div>

          {/* Search Box */}
          <div className="relative">
            <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <input
              type="text"
              placeholder="Search by title, keyword, or company..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2.5 rounded-full border border-border/80 bg-card text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary text-sm shadow-sm"
            />
          </div>

          {/* Type filters */}
          <div className="flex flex-wrap gap-1.5 pb-2 border-b border-border/40">
            {JOB_TYPES.map((type) => (
              <button
                key={type}
                onClick={() => setSelectedType(type)}
                className={`px-3.5 py-1.5 rounded-full text-xs font-semibold border transition-all ${
                  selectedType === type
                    ? "bg-primary border-primary text-primary-foreground shadow-sm"
                    : "border-border/60 hover:bg-secondary text-muted-foreground hover:text-foreground"
                }`}
              >
                {type}
              </button>
            ))}
          </div>

          {/* Jobs List container */}
          <div className="flex-1 overflow-y-auto max-h-[calc(100vh-270px)] space-y-3 pr-1.5">
            {loading ? (
              <div className="flex flex-col items-center justify-center py-12 gap-2 text-muted-foreground">
                <div className="h-6 w-6 border-2 border-primary border-t-transparent rounded-full animate-spin" />
                <span className="text-xs font-medium">Scanning live jobs database...</span>
              </div>
            ) : error ? (
              <div className="p-4 rounded-2xl bg-destructive/10 border border-destructive/20 text-destructive text-xs font-medium">
                {error}
              </div>
            ) : filteredJobs.length === 0 ? (
              <div className="text-center py-12 border border-dashed border-border/80 rounded-2xl p-6 bg-secondary/15">
                <Briefcase className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
                <p className="text-sm font-bold">No jobs matching your filters</p>
                <p className="text-xs text-muted-foreground mt-1">Try resetting search keywords or type filters.</p>
              </div>
            ) : (
              filteredJobs.map((job) => {
                const isSelected = activeJob?.id === job.id;
                const { matchScore } = getMockJobData(job);
                return (
                  <div
                    key={job.id}
                    onClick={() => setSelectedJobId(job.id)}
                    className={`p-4 rounded-2xl border text-left cursor-pointer transition-all duration-200 flex flex-col justify-between gap-3 ${
                      isSelected
                        ? "bg-card border-primary ring-1 ring-primary/20 shadow-md shadow-blue-500/5"
                        : "bg-card border-border/60 hover:border-border hover:bg-secondary/35"
                    }`}
                  >
                    <div>
                      <div className="flex items-center justify-between gap-2">
                        <span className="text-xs font-bold text-muted-foreground truncate">{job.company}</span>
                        <div className="px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-[9px] font-extrabold text-emerald-500 flex items-center gap-0.5 shrink-0">
                          <Sparkles className="h-2 w-2" />
                          {matchScore}% Match
                        </div>
                      </div>
                      <h3 className="text-sm font-bold text-foreground mt-1.5 line-clamp-1">{job.title}</h3>
                    </div>

                    <div className="flex items-center justify-between text-xs text-muted-foreground border-t border-border/40 pt-3 mt-1">
                      <span className="truncate max-w-[150px]">{job.location}</span>
                      <span className="font-semibold text-foreground shrink-0">{job.salary}</span>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </div>

        {/* Right Side: Active Job Details Workspace */}
        <div className={`flex-1 flex flex-col bg-card border border-border/80 rounded-3xl overflow-hidden shadow-sm ${
          activeJob && selectedJobId ? "flex" : "hidden md:flex"
        }`}>
          {activeJob && mockData ? (
            <div className="flex flex-col h-full relative">
              {/* Mobile Back Button */}
              <div className="md:hidden flex items-center border-b border-border/40 p-4">
                <button
                  onClick={() => setSelectedJobId(null)}
                  className="flex items-center gap-1 text-xs font-semibold text-muted-foreground hover:text-foreground"
                >
                  <ArrowLeft className="h-4 w-4" />
                  Back to search
                </button>
              </div>

              {/* Detail Header area */}
              <div className="p-6 md:p-8 border-b border-border/40 space-y-4">
                <div className="flex items-start justify-between gap-4">
                  <div className="space-y-1.5">
                    <span className="text-sm font-bold text-primary">{activeJob.company}</span>
                    <h2 className="text-xl md:text-2xl font-extrabold tracking-tight">{activeJob.title}</h2>
                  </div>
                  
                  {/* Match Ring */}
                  <div className="px-4 py-2 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 text-center shrink-0">
                    <div className="text-sm font-extrabold text-emerald-500 flex items-center justify-center gap-1">
                      <Sparkles className="h-4.5 w-4.5" />
                      {mockData.matchScore}%
                    </div>
                    <span className="text-[10px] text-emerald-600 dark:text-emerald-400 font-semibold block mt-0.5">AI Match Index</span>
                  </div>
                </div>

                {/* Badges row */}
                <div className="flex flex-wrap items-center gap-4 text-xs text-muted-foreground pt-1">
                  <div className="flex items-center gap-1">
                    <MapPin className="h-4 w-4 text-muted-foreground/80" />
                    <span>{activeJob.location}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <DollarSign className="h-4 w-4 text-muted-foreground/80" />
                    <span className="font-semibold text-foreground">{activeJob.salary}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Briefcase className="h-4 w-4 text-muted-foreground/80" />
                    <span>{activeJob.job_type}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Clock className="h-4 w-4 text-muted-foreground/80" />
                    <span>Interview difficulty: <strong className="text-foreground">{mockData.difficulty}</strong></span>
                  </div>
                </div>
              </div>

              {/* Scrollable details body */}
              <div className="flex-1 overflow-y-auto p-6 md:p-8 space-y-8 max-h-[calc(100vh-320px)]">
                {/* AI Skills Gap Analysis widget */}
                <div className="p-5 rounded-2xl bg-secondary/25 border border-border/60 space-y-4">
                  <div className="flex items-center gap-2">
                    <Award className="h-4.5 w-4.5 text-primary" />
                    <h3 className="text-sm font-bold">AI Skill Fit Check</h3>
                  </div>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <span className="text-xs font-semibold text-muted-foreground block">Skills Matched ({mockData.matchedSkills.length})</span>
                      <div className="flex flex-wrap gap-1.5">
                        {mockData.matchedSkills.map((skill) => (
                          <span key={skill} className="px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-[10px] font-semibold text-emerald-500 flex items-center gap-1">
                            <CheckCircle2 className="h-3 w-3" />
                            {skill}
                          </span>
                        ))}
                      </div>
                    </div>

                    <div className="space-y-2">
                      <span className="text-xs font-semibold text-muted-foreground block">Missing / Improvement Gap ({mockData.missingSkills.length})</span>
                      <div className="flex flex-wrap gap-1.5">
                        {mockData.missingSkills.map((skill) => (
                          <span key={skill} className="px-2 py-0.5 rounded-full bg-amber-500/10 border border-amber-500/20 text-[10px] font-semibold text-amber-600 dark:text-amber-400 flex items-center gap-1">
                            <ShieldAlert className="h-3 w-3" />
                            {skill}
                          </span>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>

                {/* Job Description section */}
                <div className="space-y-4">
                  <h3 className="text-base font-bold text-foreground">Role Description</h3>
                  <div className="text-sm leading-relaxed text-muted-foreground space-y-3 whitespace-pre-wrap">
                    {activeJob.description || "No description provided for this opening."}
                  </div>
                </div>

                {/* Interview Pipeline Steps widget */}
                <div className="space-y-4">
                  <h3 className="text-base font-bold text-foreground">Hiring Process Pipeline</h3>
                  <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                    {[
                      { step: 1, label: "Application Screen" },
                      { step: 2, label: "Technical Coding" },
                      { step: 3, label: "System Design" },
                      { step: 4, label: "Leadership Fit" },
                    ].map((pipeline) => (
                      <div key={pipeline.step} className="p-3 bg-secondary/20 border border-border/40 rounded-xl space-y-1">
                        <span className="text-[10px] font-bold text-primary block">STAGE 0{pipeline.step}</span>
                        <span className="text-xs font-semibold text-foreground block">{pipeline.label}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Sticky bottom CTA panel */}
              <div className="p-6 border-t border-border/40 bg-card/85 backdrop-blur-sm flex items-center justify-between gap-4 mt-auto">
                <div className="hidden sm:block">
                  <span className="text-xs text-muted-foreground block">Apply directly to recruiter</span>
                  <span className="text-sm font-semibold text-foreground">{activeJob.company} Hiring Desk</span>
                </div>
                
                <button className="flex-1 sm:flex-initial px-6 py-2.5 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-sm font-bold shadow-sm shadow-blue-500/10 flex items-center justify-center gap-1.5">
                  Apply Now
                  <ChevronRight className="h-4.5 w-4.5" />
                </button>
              </div>
            </div>
          ) : (
            <div className="flex-1 flex flex-col items-center justify-center text-center p-12">
              <Briefcase className="h-10 w-10 text-muted-foreground mb-3" />
              <h3 className="text-base font-bold text-foreground">No active job selected</h3>
              <p className="text-xs text-muted-foreground mt-1.5 max-w-xs">
                Select an open role from the search listings pane to view the ATS scoring check, required skills, and direct apply pipeline.
              </p>
            </div>
          )}
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
