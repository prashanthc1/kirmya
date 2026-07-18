"use client";

import React, { useState, useEffect } from "react";
import {
  User,
  Briefcase,
  GraduationCap,
  Code,
  FolderGit,
  Award,
  ShieldCheck,
  Mail,
  Settings,
  ShieldAlert,
  ChevronDown,
  ChevronUp,
  Plus,
  Trash2,
  Eye,
  EyeOff,
  Sparkles,
  Check,
  Download,
  Calendar,
  ExternalLink,
  Users,
  BarChart3,
  HelpCircle,
  FileText,
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import {
  ExtendedProfile,
  ProjectItem,
  AchievementItem,
  WorkExperience,
  Education,
  Certification,
  ProfileSkill,
  Language,
  PortfolioLink,
} from "./types";

interface ProfileCenterWorkspaceProps {
  profile: ExtendedProfile;
  activeSectionId: string;
  setActiveSectionId: (id: string) => void;
  onUpdateField: (updatedFields: Partial<ExtendedProfile>) => void;
}

// Common dial codes; the value is stored as a full E.164 string (+<cc><number>).
const DIAL_CODES: [string, string][] = [
  ["+1", "US/CA"],
  ["+44", "UK"],
  ["+91", "IN"],
  ["+61", "AU"],
  ["+49", "DE"],
  ["+33", "FR"],
  ["+81", "JP"],
  ["+86", "CN"],
  ["+971", "AE"],
  ["+65", "SG"],
  ["+27", "ZA"],
  ["+55", "BR"],
  ["+7", "RU"],
  ["+34", "ES"],
  ["+39", "IT"],
];

// E.164: leading '+', country digit 1-9, then up to 14 more digits.
export const isValidE164 = (v: string) => /^\+[1-9]\d{7,14}$/.test(v);

function PhoneField({
  value,
  onChange,
}: {
  value: string;
  onChange: (v: string) => void;
}) {
  const dial =
    DIAL_CODES.find(([d]) => value.startsWith(d))?.[0] || "+1";
  const national = (value.startsWith(dial)
    ? value.slice(dial.length)
    : value.replace(/^\+/, "")
  ).replace(/\D/g, "");
  const valid = isValidE164(value);
  const compose = (d: string, n: string) => onChange(d + n.replace(/\D/g, ""));

  return (
    <div>
      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
        Phone (International)
      </label>
      <div className="flex gap-2">
        <select
          value={dial}
          onChange={(e) => compose(e.target.value, national)}
          className="bg-secondary/40 border border-border focus:border-primary rounded-xl px-2 py-2 outline-none shrink-0"
        >
          {DIAL_CODES.map(([d, label]) => (
            <option key={d} value={d}>
              {d} {label}
            </option>
          ))}
        </select>
        <input
          type="tel"
          inputMode="numeric"
          value={national}
          placeholder="9876543210"
          onChange={(e) => compose(dial, e.target.value)}
          className={`w-full bg-secondary/40 border rounded-xl px-3.5 py-2 outline-none ${
            value && !valid
              ? "border-destructive focus:border-destructive"
              : "border-border focus:border-primary"
          }`}
        />
      </div>
      <p
        className={`mt-1 text-[10px] font-semibold ${
          !value
            ? "text-muted-foreground/60"
            : valid
              ? "text-emerald-600 dark:text-emerald-400"
              : "text-destructive"
        }`}
      >
        {!value
          ? "Select country code and enter your number"
          : valid
            ? `✓ Valid — stored as ${value}`
            : "Enter a valid international number (7–15 digits)"}
      </p>
    </div>
  );
}

export default function ProfileCenterWorkspace({
  profile,
  activeSectionId,
  setActiveSectionId,
  onUpdateField,
}: ProfileCenterWorkspaceProps) {
  // Track open state for all 15 cards
  const [expandedSections, setExpandedSections] = useState<
    Record<string, boolean>
  >({
    identity: true,
    summary: false,
    experience: false,
    education: false,
    skills: false,
    projects: false,
    certifications: false,
    achievements: false,
    resume: false,
    preferences: false,
    verification: false,
    networking: false,
    analytics: false,
    privacy: false,
    aicoach: false,
  });

  // Track edit modes for individual sections
  const [editingSections, setEditingSections] = useState<
    Record<string, boolean>
  >({});

  const toggleSection = (sectionId: string) => {
    setExpandedSections((prev) => ({
      ...prev,
      [sectionId]: !prev[sectionId],
    }));
    setActiveSectionId(sectionId);
  };

  const toggleEditMode = (sectionId: string) => {
    setEditingSections((prev) => ({
      ...prev,
      [sectionId]: !prev[sectionId],
    }));
    setActiveSectionId(sectionId);
  };

  const handleInputChange = (field: keyof ExtendedProfile, value: unknown) => {
    onUpdateField({ [field]: value });
  };

  // Sections definition for modular rendering
  return (
    <div className="flex-grow space-y-6 max-w-3xl">
      {/* ---------------- 1. Identity & Personal Information ---------------- */}
      <div
        id="section-identity"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "identity"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("identity")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("identity")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-blue-500/10 flex items-center justify-center text-blue-500">
              <User className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                1. Identity &amp; Personal Info
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Contact, Pronouns, &amp; Work Auth
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleEditMode("identity");
              }}
              className="px-3.5 py-1.5 border border-border/80 hover:bg-secondary rounded-full text-[10px] font-bold tracking-wider uppercase cursor-pointer"
            >
              {editingSections.identity ? "Done" : "Edit inline"}
            </button>
            {expandedSections.identity ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.identity && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-6">
                {editingSections.identity ? (
                  <div className="grid grid-cols-2 gap-4 text-xs">
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Full Name
                      </label>
                      <input
                        type="text"
                        value={profile.full_name || ""}
                        onChange={(e) =>
                          handleInputChange("full_name", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Preferred Name
                      </label>
                      <input
                        type="text"
                        value={profile.preferred_name || ""}
                        onChange={(e) =>
                          handleInputChange("preferred_name", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      />
                    </div>
                    <div className="col-span-2">
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Professional Headline
                      </label>
                      <input
                        type="text"
                        value={profile.headline || ""}
                        onChange={(e) =>
                          handleInputChange("headline", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      />
                    </div>
                    <div className="col-span-2">
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Bio Description
                      </label>
                      <textarea
                        rows={3}
                        value={profile.bio || ""}
                        onChange={(e) =>
                          handleInputChange("bio", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none resize-none"
                      />
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Location
                      </label>
                      <input
                        type="text"
                        value={profile.location || ""}
                        onChange={(e) =>
                          handleInputChange("location", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      />
                    </div>
                    <PhoneField
                      value={profile.phone || ""}
                      onChange={(v) => handleInputChange("phone", v)}
                    />
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Email
                      </label>
                      <input
                        type="email"
                        value={profile.email || ""}
                        onChange={(e) =>
                          handleInputChange("email", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Visa/Work Status
                      </label>
                      <input
                        type="text"
                        value={profile.visa_status || ""}
                        onChange={(e) =>
                          handleInputChange("visa_status", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        LinkedIn Link
                      </label>
                      <input
                        type="text"
                        value={profile.linkedin_url || ""}
                        onChange={(e) =>
                          handleInputChange("linkedin_url", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        GitHub Link
                      </label>
                      <input
                        type="text"
                        value={profile.github_url || ""}
                        onChange={(e) =>
                          handleInputChange("github_url", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      />
                    </div>
                  </div>
                ) : (
                  <div className="space-y-4 text-xs">
                    <p className="text-sm font-medium text-foreground">
                      {profile.bio}
                    </p>
                    <div className="grid grid-cols-2 sm:grid-cols-3 gap-y-3.5 gap-x-4 border-t border-border/40 pt-4 text-muted-foreground">
                      <div>
                        <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                          Preferred Name
                        </span>
                        <span className="font-semibold text-foreground">
                          {profile.preferred_name} ({profile.pronouns})
                        </span>
                      </div>
                      <div>
                        <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                          Email Address
                        </span>
                        <span className="font-semibold text-foreground">
                          {profile.email}
                        </span>
                      </div>
                      <div>
                        <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                          Mobile Phone
                        </span>
                        <span className="font-semibold text-foreground">
                          {profile.phone}
                        </span>
                      </div>
                      <div>
                        <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                          Visa Status
                        </span>
                        <span className="font-semibold text-foreground">
                          {profile.visa_status}
                        </span>
                      </div>
                      <div>
                        <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                          LinkedIn Profile
                        </span>
                        <a
                          href={profile.linkedin_url}
                          target="_blank"
                          className="font-semibold text-primary hover:underline flex items-center gap-0.5"
                        >
                          LinkedIn <ExternalLink className="h-3 w-3" />
                        </a>
                      </div>
                      <div>
                        <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                          GitHub Port
                        </span>
                        <a
                          href={profile.github_url}
                          target="_blank"
                          className="font-semibold text-primary hover:underline flex items-center gap-0.5"
                        >
                          GitHub <ExternalLink className="h-3 w-3" />
                        </a>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 2. Professional Summary ---------------- */}
      <div
        id="section-summary"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "summary"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("summary")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("summary")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-emerald-500/10 flex items-center justify-center text-emerald-500">
              <FileText className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                2. Professional Summary
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Elevator Pitch &amp; Brand Statement
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleEditMode("summary");
              }}
              className="px-3.5 py-1.5 border border-border/80 hover:bg-secondary rounded-full text-[10px] font-bold tracking-wider uppercase cursor-pointer"
            >
              {editingSections.summary ? "Done" : "Edit inline"}
            </button>
            {expandedSections.summary ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.summary && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-6">
                {editingSections.summary ? (
                  <div className="space-y-4 text-xs">
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Career Overview
                      </label>
                      <textarea
                        rows={3}
                        value={profile.about || ""}
                        onChange={(e) =>
                          handleInputChange("about", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none resize-none"
                      />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                          Personal Brand Title
                        </label>
                        <input
                          type="text"
                          value={profile.personal_brand || ""}
                          onChange={(e) =>
                            handleInputChange("personal_brand", e.target.value)
                          }
                          className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                        />
                      </div>
                      <div>
                        <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                          Industry Sector
                        </label>
                        <input
                          type="text"
                          value={profile.industry || ""}
                          onChange={(e) =>
                            handleInputChange("industry", e.target.value)
                          }
                          className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                        />
                      </div>
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Elevator Pitch (30s Intro)
                      </label>
                      <textarea
                        rows={3}
                        value={profile.elevator_pitch || ""}
                        onChange={(e) =>
                          handleInputChange("elevator_pitch", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none resize-none"
                      />
                    </div>
                  </div>
                ) : (
                  <div className="space-y-4 text-xs">
                    <div>
                      <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block mb-1">
                        Career Overview
                      </span>
                      <p className="text-foreground leading-relaxed">
                        {profile.about}
                      </p>
                    </div>
                    <div>
                      <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block mb-1">
                        Personal Brand &amp; Pitch
                      </span>
                      <div className="p-4 bg-secondary/30 border border-border/40 rounded-2xl italic leading-relaxed text-muted-foreground">
                        &quot;{profile.elevator_pitch}&quot;
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 3. Work Experience ---------------- */}
      <div
        id="section-experience"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "experience"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("experience")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("experience")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-purple-500/10 flex items-center justify-center text-purple-500">
              <Briefcase className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                3. Work Experience
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Career History &amp; STAR KPIs
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleEditMode("experience");
              }}
              className="px-3.5 py-1.5 border border-border/80 hover:bg-secondary rounded-full text-[10px] font-bold tracking-wider uppercase cursor-pointer"
            >
              {editingSections.experience ? "Done" : "Edit timeline"}
            </button>
            {expandedSections.experience ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.experience && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-6">
                {editingSections.experience ? (
                  <div className="space-y-6 text-xs">
                    {profile.experiences.map((exp, i) => (
                      <div
                        key={exp.id || i}
                        className="p-4 bg-secondary/20 border border-border/60 rounded-2xl space-y-4"
                      >
                        <div className="flex justify-between items-center">
                          <h4 className="font-bold text-foreground">
                            Job #{i + 1}
                          </h4>
                          <button
                            onClick={() => {
                              const updated = profile.experiences.filter(
                                (e) => e.id !== exp.id,
                              );
                              handleInputChange("experiences", updated);
                            }}
                            className="text-destructive hover:underline font-bold"
                          >
                            Remove
                          </button>
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Company
                            </label>
                            <input
                              type="text"
                              value={exp.company}
                              onChange={(e) => {
                                const list = [...profile.experiences];
                                list[i].company = e.target.value;
                                handleInputChange("experiences", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Position Title
                            </label>
                            <input
                              type="text"
                              value={exp.title}
                              onChange={(e) => {
                                const list = [...profile.experiences];
                                list[i].title = e.target.value;
                                handleInputChange("experiences", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Start Date
                            </label>
                            <input
                              type="text"
                              value={exp.start_date}
                              onChange={(e) => {
                                const list = [...profile.experiences];
                                list[i].start_date = e.target.value;
                                handleInputChange("experiences", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              End Date
                            </label>
                            <input
                              type="text"
                              value={exp.is_current ? "Present" : exp.end_date}
                              onChange={(e) => {
                                const list = [...profile.experiences];
                                list[i].end_date = e.target.value;
                                handleInputChange("experiences", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                          <div className="col-span-2">
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Achievements (STAR points, one per line)
                            </label>
                            <textarea
                              rows={3}
                              value={exp.achievements?.join("\n") || ""}
                              onChange={(e) => {
                                const list = [...profile.experiences];
                                list[i].achievements =
                                  e.target.value.split("\n");
                                handleInputChange("experiences", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none resize-none"
                            />
                          </div>
                        </div>
                      </div>
                    ))}
                    <button
                      onClick={() => {
                        const newExp: WorkExperience & {
                          id: string;
                          achievements: string[];
                        } = {
                          id: "exp_new_" + Date.now(),
                          company: "New Company",
                          title: "New Role",
                          location: "Denver, CO",
                          employment_type: "Full-Time",
                          start_date: "2025-01",
                          end_date: "Present",
                          is_current: true,
                          description: "",
                          achievements: [],
                        };
                        handleInputChange("experiences", [
                          ...profile.experiences,
                          newExp,
                        ]);
                      }}
                      className="w-full border border-dashed border-border hover:border-primary text-primary font-bold py-2.5 rounded-xl flex items-center justify-center gap-1 cursor-pointer"
                    >
                      <Plus className="h-4 w-4" /> Add Experience
                    </button>
                  </div>
                ) : (
                  <div className="space-y-8 text-xs">
                    {profile.experiences.map((exp, idx) => (
                      <div key={exp.id || idx} className="flex gap-4 group">
                        <div className="flex flex-col items-center shrink-0">
                          <div className="h-6 w-6 rounded-full bg-secondary border border-border flex items-center justify-center font-bold text-[10px] text-muted-foreground group-hover:border-primary group-hover:text-primary transition-all">
                            {idx + 1}
                          </div>
                          {idx < profile.experiences.length - 1 && (
                            <div className="w-[1.5px] bg-border/50 flex-grow my-1.5" />
                          )}
                        </div>

                        <div className="space-y-3 pb-2 flex-grow">
                          <div>
                            <h4 className="text-sm font-bold text-foreground">
                              {exp.title}
                            </h4>
                            <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-0.5">
                              <span className="font-bold text-foreground">
                                {exp.company}
                              </span>
                              <span>&bull;</span>
                              <span>
                                {exp.start_date} –{" "}
                                {exp.is_current ? "Present" : exp.end_date}
                              </span>
                              <span>&bull;</span>
                              <span className="capitalize">
                                {exp.location_type || "hybrid"}
                              </span>
                            </div>
                          </div>
                          <p className="text-muted-foreground leading-relaxed">
                            {exp.description}
                          </p>

                          {/* Achievements list */}
                          {exp.achievements && exp.achievements.length > 0 && (
                            <div className="space-y-1.5">
                              <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                                Key Achievements &amp; KPIs
                              </span>
                              <ul className="list-disc list-inside space-y-1 text-muted-foreground pl-1">
                                {exp.achievements.map(
                                  (ach: string, ai: number) => (
                                    <li key={ai} className="leading-relaxed">
                                      <span className="text-foreground font-medium">
                                        {ach}
                                      </span>
                                    </li>
                                  ),
                                )}
                              </ul>
                            </div>
                          )}

                          {/* Tech stack tags */}
                          {exp.technologies && exp.technologies.length > 0 && (
                            <div className="flex flex-wrap gap-1.5 pt-1">
                              {exp.technologies.map((t: string) => (
                                <span
                                  key={t}
                                  className="px-2 py-0.5 bg-secondary text-muted-foreground border border-border/50 rounded-md text-[9px] font-semibold"
                                >
                                  {t}
                                </span>
                              ))}
                            </div>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 4. Education ---------------- */}
      <div
        id="section-education"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "education"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("education")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("education")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-indigo-500/10 flex items-center justify-center text-indigo-500">
              <GraduationCap className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                4. Education
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Degrees, Research, &amp; Theses
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleEditMode("education");
              }}
              className="px-3.5 py-1.5 border border-border/80 hover:bg-secondary rounded-full text-[10px] font-bold tracking-wider uppercase cursor-pointer"
            >
              {editingSections.education ? "Done" : "Edit timeline"}
            </button>
            {expandedSections.education ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.education && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-6">
                {editingSections.education ? (
                  <div className="space-y-6 text-xs">
                    {profile.educations.map((edu, i) => (
                      <div
                        key={edu.id || i}
                        className="p-4 bg-secondary/20 border border-border/60 rounded-2xl space-y-4"
                      >
                        <div className="grid grid-cols-2 gap-4">
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Institution
                            </label>
                            <input
                              type="text"
                              value={edu.school}
                              onChange={(e) => {
                                const list = [...profile.educations];
                                list[i].school = e.target.value;
                                handleInputChange("educations", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Degree
                            </label>
                            <input
                              type="text"
                              value={edu.degree}
                              onChange={(e) => {
                                const list = [...profile.educations];
                                list[i].degree = e.target.value;
                                handleInputChange("educations", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              GPA / Grade
                            </label>
                            <input
                              type="text"
                              value={edu.grade}
                              onChange={(e) => {
                                const list = [...profile.educations];
                                list[i].grade = e.target.value;
                                handleInputChange("educations", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Thesis Topic
                            </label>
                            <input
                              type="text"
                              value={edu.thesis || ""}
                              onChange={(e) => {
                                const list = [...profile.educations];
                                list[i].thesis = e.target.value;
                                handleInputChange("educations", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                        </div>
                        <div className="flex justify-end">
                          <button
                            onClick={() => {
                              const updated = profile.educations.filter(
                                (_, idx) => idx !== i,
                              );
                              handleInputChange("educations", updated);
                            }}
                            className="text-destructive hover:underline font-bold"
                          >
                            Remove
                          </button>
                        </div>
                      </div>
                    ))}
                    <button
                      onClick={() => {
                        const newEdu: Education & { id: string } = {
                          id: "edu_new_" + Date.now(),
                          school: "New Institution",
                          degree: "Degree",
                          field_of_study: "Field of Study",
                          start_date: "2019-09",
                          end_date: "2023-06",
                          grade: "",
                          description: "",
                        };
                        handleInputChange("educations", [
                          ...profile.educations,
                          newEdu,
                        ]);
                      }}
                      className="w-full border border-dashed border-border hover:border-primary text-primary font-bold py-2.5 rounded-xl flex items-center justify-center gap-1 cursor-pointer"
                    >
                      <Plus className="h-4 w-4" /> Add Education
                    </button>
                  </div>
                ) : (
                  <div className="space-y-6 text-xs">
                    {profile.educations.map((edu, idx) => (
                      <div key={edu.id || idx} className="space-y-2">
                        <div className="flex justify-between items-start">
                          <div>
                            <h4 className="text-sm font-bold text-foreground">
                              {edu.school}
                            </h4>
                            <p className="text-xs text-muted-foreground font-semibold mt-0.5">
                              {edu.degree} &bull; {edu.field_of_study}
                            </p>
                          </div>
                          <span className="text-[10px] font-bold text-primary bg-primary/5 px-2 py-0.5 rounded-full">
                            {edu.grade}
                          </span>
                        </div>
                        <p className="text-muted-foreground leading-relaxed">
                          {edu.description}
                        </p>
                        {edu.thesis && (
                          <div className="bg-secondary/40 border border-border/40 p-3 rounded-xl text-muted-foreground">
                            <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block mb-0.5">
                              Master&apos;s Thesis
                            </span>
                            <span className="font-semibold text-foreground italic">
                              &quot;{edu.thesis}&quot;
                            </span>
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 5. Skills & Expertise ---------------- */}
      <div
        id="section-skills"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "skills"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("skills")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("skills")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-emerald-500/10 flex items-center justify-center text-emerald-500">
              <Code className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                5. Skills &amp; Expertise
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Verified Tech stack &amp; Soft skills
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleEditMode("skills");
              }}
              className="px-3.5 py-1.5 border border-border/80 hover:bg-secondary rounded-full text-[10px] font-bold tracking-wider uppercase cursor-pointer"
            >
              {editingSections.skills ? "Done" : "Edit skills"}
            </button>
            {expandedSections.skills ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.skills && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-4">
                {editingSections.skills ? (
                  <div className="space-y-4 text-xs">
                    <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider">
                      Comma Separated Skills
                    </label>
                    <input
                      type="text"
                      value={profile.skills.map((s) => s.name).join(", ")}
                      onChange={(e) => {
                        const splitted = e.target.value
                          .split(",")
                          .map((s) => s.trim())
                          .filter(Boolean);
                        const list = splitted.map((name) => {
                          const existing = profile.skills.find(
                            (x) => x.name.toLowerCase() === name.toLowerCase(),
                          );
                          return (
                            existing || {
                              name,
                              proficiency_level: "Intermediate",
                              endorsed_count: 0,
                            }
                          );
                        });
                        handleInputChange("skills", list);
                      }}
                      className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                    />
                  </div>
                ) : (
                  <div className="flex flex-wrap gap-2 text-xs">
                    {profile.skills.map((skill) => (
                      <span
                        key={skill.name}
                        className={`px-3 py-1.5 border rounded-full font-semibold flex items-center gap-1.5 ${
                          skill.verification_status === "verified"
                            ? "bg-emerald-500/5 border-emerald-500/20 text-emerald-600 dark:text-emerald-400"
                            : "bg-secondary text-muted-foreground border-border/40"
                        }`}
                      >
                        {skill.name}
                        {skill.verification_status === "verified" && (
                          <Check className="h-3 w-3 fill-emerald-500/10" />
                        )}
                        <span className="text-[9px] opacity-60 font-bold">
                          ({skill.proficiency_level})
                        </span>
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 6. Projects & Portfolio ---------------- */}
      <div
        id="section-projects"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "projects"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("projects")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("projects")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-orange-500/10 flex items-center justify-center text-orange-500">
              <FolderGit className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                6. Projects &amp; Portfolio
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Project Gallery &amp; Live Demos
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleEditMode("projects");
              }}
              className="px-3.5 py-1.5 border border-border/80 hover:bg-secondary rounded-full text-[10px] font-bold tracking-wider uppercase cursor-pointer"
            >
              {editingSections.projects ? "Done" : "Edit gallery"}
            </button>
            {expandedSections.projects ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.projects && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-6">
                {editingSections.projects ? (
                  <div className="space-y-6 text-xs">
                    {profile.projects.map((proj, i) => (
                      <div
                        key={proj.id}
                        className="p-4 bg-secondary/20 border border-border/60 rounded-2xl space-y-4"
                      >
                        <div className="grid grid-cols-2 gap-4">
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Project Title
                            </label>
                            <input
                              type="text"
                              value={proj.title}
                              onChange={(e) => {
                                const list = [...profile.projects];
                                list[i].title = e.target.value;
                                handleInputChange("projects", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                          <div>
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Metrics &amp; ROI
                            </label>
                            <input
                              type="text"
                              value={proj.metrics || ""}
                              onChange={(e) => {
                                const list = [...profile.projects];
                                list[i].metrics = e.target.value;
                                handleInputChange("projects", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none"
                            />
                          </div>
                          <div className="col-span-2">
                            <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                              Description
                            </label>
                            <textarea
                              rows={2}
                              value={proj.description}
                              onChange={(e) => {
                                const list = [...profile.projects];
                                list[i].description = e.target.value;
                                handleInputChange("projects", list);
                              }}
                              className="w-full bg-card border border-border rounded-xl px-3.5 py-2 outline-none resize-none"
                            />
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
                    {profile.projects.map((proj) => (
                      <div
                        key={proj.id}
                        className="border border-border/80 rounded-2xl overflow-hidden bg-card flex flex-col justify-between shadow-sm"
                      >
                        {proj.cover_image && (
                          <div className="h-32 overflow-hidden bg-secondary">
                            <img
                              src={proj.cover_image}
                              alt={proj.title}
                              className="w-full h-full object-cover"
                            />
                          </div>
                        )}
                        <div className="p-4 space-y-2 flex-grow">
                          <h4 className="font-bold text-foreground">
                            {proj.title}
                          </h4>
                          <p className="text-muted-foreground leading-normal line-clamp-3">
                            {proj.description}
                          </p>
                          {proj.metrics && (
                            <div className="text-[10px] text-emerald-500 bg-emerald-500/5 px-2 py-1 rounded-md font-bold inline-block border border-emerald-500/10 mt-1">
                              📈 Impact: {proj.metrics}
                            </div>
                          )}
                        </div>
                        <div className="p-4 border-t border-border/40 bg-secondary/20 flex gap-3 text-primary font-bold">
                          {proj.repository_url && (
                            <a
                              href={proj.repository_url}
                              target="_blank"
                              className="flex items-center gap-0.5"
                            >
                              Repo <ExternalLink className="h-3.5 w-3.5" />
                            </a>
                          )}
                          {proj.live_demo_url && (
                            <a
                              href={proj.live_demo_url}
                              target="_blank"
                              className="flex items-center gap-0.5"
                            >
                              Demo <ExternalLink className="h-3.5 w-3.5" />
                            </a>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 7. Certifications & Licenses ---------------- */}
      <div
        id="section-certifications"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "certifications"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("certifications")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("certifications")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-violet-500/10 flex items-center justify-center text-violet-500">
              <Award className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                7. Certifications &amp; Licenses
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Credential IDs &amp; Verification URLs
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleEditMode("certifications");
              }}
              className="px-3.5 py-1.5 border border-border/80 hover:bg-secondary rounded-full text-[10px] font-bold tracking-wider uppercase cursor-pointer"
            >
              {editingSections.certifications ? "Done" : "Edit certs"}
            </button>
            {expandedSections.certifications ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.certifications && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-4">
                {editingSections.certifications ? (
                  <div className="space-y-4 text-xs">
                    {profile.certifications.map((cert, i) => (
                      <div
                        key={cert.id}
                        className="p-3 bg-secondary/20 border border-border/60 rounded-xl grid grid-cols-2 gap-3"
                      >
                        <div className="col-span-2">
                          <input
                            type="text"
                            value={cert.name}
                            placeholder="Cert name"
                            onChange={(e) => {
                              const list = [...profile.certifications];
                              list[i].name = e.target.value;
                              handleInputChange("certifications", list);
                            }}
                            className="w-full bg-card border border-border rounded-xl px-3 py-1.5 outline-none font-bold"
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="space-y-4 text-xs">
                    {profile.certifications.map((cert) => (
                      <div
                        key={cert.id}
                        className="flex justify-between items-start p-3 bg-secondary/35 border border-border/40 rounded-2xl"
                      >
                        <div>
                          <h4 className="font-bold text-foreground">
                            {cert.name}
                          </h4>
                          <p className="text-[11px] text-muted-foreground mt-0.5">
                            {cert.issuer} &bull; ID: {cert.credential_id}
                          </p>
                        </div>
                        {cert.credential_url && (
                          <a
                            href={cert.credential_url}
                            target="_blank"
                            className="text-primary hover:underline font-bold flex items-center gap-0.5 shrink-0"
                          >
                            Verify <ExternalLink className="h-3 w-3" />
                          </a>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 8. Achievements ---------------- */}
      <div
        id="section-achievements"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "achievements"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("achievements")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("achievements")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-yellow-500/10 flex items-center justify-center text-yellow-500">
              <Award className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                8. Achievements
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Awards, Patents, &amp; Publications
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleEditMode("achievements");
              }}
              className="px-3.5 py-1.5 border border-border/80 hover:bg-secondary rounded-full text-[10px] font-bold tracking-wider uppercase cursor-pointer"
            >
              {editingSections.achievements ? "Done" : "Edit list"}
            </button>
            {expandedSections.achievements ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.achievements && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-4">
                {editingSections.achievements ? (
                  <div className="space-y-4 text-xs">
                    {profile.achievements_list?.map((ach, i) => (
                      <div
                        key={ach.id}
                        className="p-3 bg-secondary/20 border border-border/60 rounded-xl"
                      >
                        <input
                          type="text"
                          value={ach.title}
                          onChange={(e) => {
                            const list = [...profile.achievements_list];
                            list[i].title = e.target.value;
                            handleInputChange("achievements_list", list);
                          }}
                          className="w-full bg-card border border-border rounded-xl px-3 py-1.5 outline-none font-bold"
                        />
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="space-y-4 text-xs">
                    {profile.achievements_list?.map((ach) => (
                      <div
                        key={ach.id}
                        className="p-3.5 bg-secondary/35 border border-border/40 rounded-2xl space-y-1.5"
                      >
                        <div className="flex justify-between items-start">
                          <span className="text-[9px] bg-primary/10 text-primary font-bold px-2 py-0.5 rounded-full uppercase tracking-wider">
                            {ach.category}
                          </span>
                          <span className="text-[10px] text-muted-foreground font-semibold">
                            {ach.date}
                          </span>
                        </div>
                        <h4 className="font-bold text-foreground leading-snug">
                          {ach.title}
                        </h4>
                        <p className="text-muted-foreground leading-normal">
                          {ach.description}
                        </p>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 9. Resume & Documents ---------------- */}
      <div
        id="section-resume"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "resume"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("resume")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("resume")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-rose-500/10 flex items-center justify-center text-rose-500">
              <FileText className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                9. Resume &amp; Documents
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Resumes &amp; Tailored Cover Letters
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {expandedSections.resume ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.resume && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-4 text-xs">
                <div className="space-y-3">
                  <h4 className="font-bold text-foreground">
                    Active Resume Versions
                  </h4>
                  {profile.resumes?.map((res) => (
                    <div
                      key={res.id}
                      className="flex items-center justify-between p-3.5 bg-secondary/35 border border-border/40 rounded-2xl"
                    >
                      <div className="space-y-1">
                        <p className="font-bold text-foreground flex items-center gap-1.5">
                          {res.name}
                          {res.is_primary && (
                            <span className="bg-emerald-500/10 text-emerald-600 px-1.5 py-0.5 rounded text-[8px] font-bold">
                              PRIMARY
                            </span>
                          )}
                        </p>
                        <p className="text-[10px] text-muted-foreground">
                          Uploaded{" "}
                          {new Date(res.uploaded_at).toLocaleDateString()}{" "}
                          &bull; {res.file_size}
                        </p>
                      </div>

                      {/* ATS Score widget */}
                      <div className="flex items-center gap-4 shrink-0">
                        <div className="text-right">
                          <span className="text-[9px] text-muted-foreground font-bold block uppercase tracking-wider">
                            ATS Score
                          </span>
                          <span className="font-black text-sm text-primary">
                            {res.ats_score}%
                          </span>
                        </div>
                        <button
                          onClick={() => window.print()}
                          className="h-8 w-8 rounded-xl bg-card border border-border flex items-center justify-center hover:bg-secondary cursor-pointer transition-colors"
                        >
                          <Download className="h-4 w-4 text-muted-foreground" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 10. Career Preferences ---------------- */}
      <div
        id="section-preferences"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "preferences"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("preferences")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("preferences")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-pink-500/10 flex items-center justify-center text-pink-500">
              <Settings className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                10. Career Preferences
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Salary Targets &amp; Mobility settings
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleEditMode("preferences");
              }}
              className="px-3.5 py-1.5 border border-border/80 hover:bg-secondary rounded-full text-[10px] font-bold tracking-wider uppercase cursor-pointer"
            >
              {editingSections.preferences ? "Done" : "Edit settings"}
            </button>
            {expandedSections.preferences ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.preferences && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-6">
                {editingSections.preferences ? (
                  <div className="grid grid-cols-2 gap-4 text-xs">
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Target Minimum Salary (USD)
                      </label>
                      <input
                        type="number"
                        value={profile.salary_min || 0}
                        onChange={(e) =>
                          handleInputChange(
                            "salary_min",
                            parseInt(e.target.value) || 0,
                          )
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                        Work Mode
                      </label>
                      <select
                        value={profile.work_mode || ""}
                        onChange={(e) =>
                          handleInputChange("work_mode", e.target.value)
                        }
                        className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-3.5 py-2 outline-none"
                      >
                        <option value="remote">Remote</option>
                        <option value="hybrid">Hybrid</option>
                        <option value="onsite">Onsite</option>
                      </select>
                    </div>
                  </div>
                ) : (
                  <div className="grid grid-cols-2 sm:grid-cols-3 gap-4 text-xs">
                    <div>
                      <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                        Desired Roles
                      </span>
                      <span className="font-semibold text-foreground">
                        {profile.desired_roles?.join(", ")}
                      </span>
                    </div>
                    <div>
                      <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                        Target Salary Range
                      </span>
                      <span className="font-semibold text-foreground">
                        ${profile.salary_min?.toLocaleString()} - $
                        {profile.salary_max?.toLocaleString()}{" "}
                        {profile.salary_currency}
                      </span>
                    </div>
                    <div>
                      <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                        Work Mode
                      </span>
                      <span className="font-semibold text-foreground capitalize">
                        {profile.work_mode || "Hybrid"}
                      </span>
                    </div>
                  </div>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 11. Verification & Trust ---------------- */}
      <div
        id="section-verification"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "verification"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("verification")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("verification")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-teal-500/10 flex items-center justify-center text-teal-500">
              <ShieldCheck className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                11. Verification &amp; Trust
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Credibility ratings &amp; Trust Score
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {expandedSections.verification ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.verification && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-6 text-xs">
                {/* Stripe verification simulation card */}
                <div className="bg-gradient-to-r from-teal-500/10 to-emerald-500/10 border border-teal-500/20 p-5 rounded-2xl flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                  <div className="space-y-1">
                    <h4 className="font-extrabold text-teal-600 dark:text-teal-400 flex items-center gap-1.5">
                      <ShieldCheck className="h-4.5 w-4.5" /> Verified Recruiter
                      Badge
                    </h4>
                    <p className="text-muted-foreground">
                      Kirmya screens and verifies employment history directly
                      with payroll partners.
                    </p>
                  </div>
                  <button className="px-4 py-2 bg-teal-500 hover:bg-teal-600 text-white font-bold rounded-xl shadow-md transition-all cursor-pointer shrink-0">
                    Verify via Stripe Identity
                  </button>
                </div>

                <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 text-center">
                  {[
                    { label: "Email Address", val: "Verified", active: true },
                    { label: "Phone Number", val: "Verified", active: true },
                    {
                      label: "Employment History",
                      val: "Verified",
                      active: true,
                    },
                    {
                      label: "Degrees & Education",
                      val: "Verified",
                      active: true,
                    },
                  ].map((ver) => (
                    <div
                      key={ver.label}
                      className="bg-secondary/40 border border-border/40 p-3 rounded-2xl"
                    >
                      <span className="text-[9px] text-muted-foreground font-bold uppercase tracking-wider block">
                        {ver.label}
                      </span>
                      <span className="font-black text-xs text-emerald-500 flex items-center justify-center gap-1 mt-1">
                        <Check className="h-3.5 w-3.5" /> {ver.val}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 12. Networking ---------------- */}
      <div
        id="section-networking"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "networking"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("networking")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("networking")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-blue-500/10 flex items-center justify-center text-blue-500">
              <Users className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                12. Networking &amp; Contacts
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Connections, Mentors, &amp; Referrals
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {expandedSections.networking ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.networking && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-4 text-xs">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  {profile.network?.map((c) => (
                    <div
                      key={c.id}
                      className="border border-border/60 p-3.5 rounded-2xl flex items-center gap-3 bg-secondary/20"
                    >
                      <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center font-bold text-primary">
                        {c.name.charAt(0)}
                      </div>
                      <div className="space-y-0.5">
                        <p className="font-bold text-foreground flex items-center gap-1">
                          {c.name}
                        </p>
                        <p className="text-[10px] text-muted-foreground line-clamp-1">
                          {c.headline}
                        </p>
                        <span className="text-[8px] bg-secondary text-muted-foreground/80 font-bold px-1.5 py-0.5 rounded capitalize">
                          {c.type}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 13. Analytics Dashboard ---------------- */}
      <div
        id="section-analytics"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "analytics"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("analytics")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("analytics")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-violet-500/10 flex items-center justify-center text-violet-500">
              <BarChart3 className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                13. Analytics Dashboard
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Weekly Profile Impressions &amp; Search rankings
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {expandedSections.analytics ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.analytics && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-6 text-xs">
                {/* Premium analytical charts simulation */}
                <div className="grid grid-cols-3 gap-4 text-center">
                  {[
                    {
                      label: "Profile Views",
                      val: profile.analytics?.profile_views,
                      change: "+14% vs last week",
                    },
                    {
                      label: "Search Appearances",
                      val: profile.analytics?.recruiter_searches,
                      change: "+8% vs last week",
                    },
                    {
                      label: "Resume Downloads",
                      val: profile.analytics?.resume_downloads,
                      change: "+24% vs last week",
                    },
                  ].map((an) => (
                    <div
                      key={an.label}
                      className="border border-border/60 p-4 rounded-2xl bg-secondary/15"
                    >
                      <span className="text-[9px] text-muted-foreground font-bold uppercase tracking-wider block">
                        {an.label}
                      </span>
                      <span className="font-black text-lg text-foreground block mt-1">
                        {an.val}
                      </span>
                      <span className="text-[9px] text-emerald-500 font-bold block mt-1">
                        {an.change}
                      </span>
                    </div>
                  ))}
                </div>

                {/* SVG Mock Chart */}
                <div className="bg-secondary/35 border border-border/40 p-5 rounded-2xl space-y-3">
                  <span className="font-bold text-[9px] uppercase tracking-wider text-muted-foreground/60 block">
                    Recruiter Views Trend (Last 7 Days)
                  </span>
                  <div className="h-24 flex items-end justify-between gap-4 px-4 pt-4">
                    {[12, 18, 15, 24, 30, 28, 42].map((v, i) => (
                      <div
                        key={i}
                        className="flex-1 flex flex-col items-center gap-1.5"
                      >
                        <div
                          style={{ height: `${(v / 45) * 100}%` }}
                          className="w-full bg-primary/70 hover:bg-primary rounded-t-lg transition-all"
                        />
                        <span className="text-[9px] text-muted-foreground font-bold">
                          Day {i + 1}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 14. Privacy & Security ---------------- */}
      <div
        id="section-privacy"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "privacy"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("privacy")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("privacy")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-zinc-500/10 flex items-center justify-center text-zinc-500">
              <Mail className="h-4.5 w-4.5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground">
                14. Privacy &amp; Security
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Visibility controls &amp; Active sessions
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {expandedSections.privacy ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.privacy && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-4 text-xs">
                <div className="space-y-3.5">
                  {[
                    {
                      label: "Profile Visibility",
                      desc: "Allow recruiters to view your timeline publicly",
                      key: "visibility_profile",
                    },
                    {
                      label: "Anonymous Mode",
                      desc: "Browse other user profiles without leaving footprints",
                      key: "anonymous",
                    },
                    {
                      label: "Hide Target Salary",
                      desc: "Hide expected salary field on the public listing",
                      key: "hide_salary",
                    },
                  ].map((item) => (
                    <div
                      key={item.label}
                      className="flex justify-between items-center p-3 bg-secondary/35 border border-border/40 rounded-2xl"
                    >
                      <div>
                        <p className="font-bold text-foreground">
                          {item.label}
                        </p>
                        <p className="text-[10px] text-muted-foreground mt-0.5">
                          {item.desc}
                        </p>
                      </div>
                      <div className="w-10 h-6 bg-primary rounded-full relative cursor-pointer">
                        <div className="w-4.5 h-4.5 bg-card rounded-full absolute top-0.75 right-0.75" />
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* ---------------- 15. AI Career Roadmaps ---------------- */}
      <div
        id="section-aicoach"
        className={`bg-card border rounded-3xl transition-all ${
          activeSectionId === "aicoach"
            ? "border-primary/60 shadow-md"
            : "border-border/80"
        }`}
        onClick={() => setActiveSectionId("aicoach")}
      >
        <div
          className="flex justify-between items-center p-6 cursor-pointer"
          onClick={() => toggleSection("aicoach")}
        >
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-gradient-to-r from-blue-600/10 to-indigo-600/10 flex items-center justify-center text-primary">
              <Sparkles className="h-4.5 w-4.5 animate-pulse" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-foreground flex items-center gap-1">
                15. AI Career Roadmap
              </h3>
              <p className="text-[10px] text-muted-foreground uppercase tracking-wider font-semibold">
                Weekly Growth Reports &amp; Skill gap audit
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {expandedSections.aicoach ? (
              <ChevronUp className="h-4 w-4 text-muted-foreground" />
            ) : (
              <ChevronDown className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        </div>

        <AnimatePresence>
          {expandedSections.aicoach && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: "auto", opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              className="overflow-hidden border-t border-border/40"
            >
              <div className="p-6 space-y-4 text-xs">
                <div className="border border-primary/20 bg-primary/5 p-4 rounded-2xl space-y-3 leading-relaxed text-muted-foreground">
                  <h4 className="font-extrabold text-foreground flex items-center gap-1.5 text-primary">
                    <Sparkles className="h-4 w-4 animate-spin text-primary" />{" "}
                    Active Career Narrative
                  </h4>
                  <p>{profile.career_narrative}</p>
                </div>

                <div className="p-4 bg-secondary/35 border border-border/40 rounded-2xl space-y-3">
                  <h4 className="font-bold text-foreground">
                    Identified Skill Gap Focus Areas
                  </h4>
                  <div className="space-y-2">
                    {[
                      {
                        gap: "Warehouse Robotics Automation",
                        detail:
                          "Market demand is up 42%. Adding this closes 3 critical recruiter query gaps.",
                      },
                      {
                        gap: "Supply Chain Financial Auditing",
                        detail: "Strengthens VP level resume indexing by 14%.",
                      },
                    ].map((g, i) => (
                      <div key={i} className="flex gap-2.5 items-start">
                        <span className="h-4 w-4 rounded-full bg-amber-500/10 text-amber-500 font-bold text-[9px] flex items-center justify-center shrink-0 mt-0.5">
                          !
                        </span>
                        <div>
                          <p className="font-bold text-foreground">{g.gap}</p>
                          <p className="text-[10px] text-muted-foreground">
                            {g.detail}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}
