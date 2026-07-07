"use client";

import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Check, ArrowLeft, Upload, Loader2, ArrowRight } from "lucide-react";
import { ExtendedProfile } from "./types";
import MuiModal from "@/components/shared/MuiModal";

interface ProfileOnboardingProps {
  isOpen: boolean;
  onClose: () => void;
  onComplete: (data: Partial<ExtendedProfile>) => void;
}

export default function ProfileOnboarding({
  isOpen,
  onClose,
  onComplete,
}: ProfileOnboardingProps) {
  const [step, setStep] = useState(1);
  const [formData, setFormData] = useState({
    fullName: "Marcus Hale",
    preferredName: "Marcus",
    headline: "Operations Director · Supply Chain Architect",
    desiredRole: "Director of Operations",
    salaryMin: 180000,
    workMode: "hybrid",
    resumeFile: null as File | null,
    isParsing: false,
  });

  const handleNext = () => {
    if (step < 3) {
      setStep(step + 1);
    } else {
      // Completed onboarding
      onComplete({
        headline: formData.headline,
        preferred_name: formData.preferredName,
        desired_roles: [formData.desiredRole],
        salary_min: formData.salaryMin,
        work_mode: formData.workMode as "remote" | "hybrid" | "onsite" | "",
        // Mocking resume upload on complete
        resumes: formData.resumeFile
          ? [
              {
                id: "res_uploaded_" + Date.now(),
                name: formData.resumeFile.name,
                uploaded_at: new Date().toISOString(),
                file_size: `${Math.round(formData.resumeFile.size / 1024)} KB`,
                is_primary: true,
                ats_score: 87,
                readability_score: 91,
                file_url: "#",
              },
            ]
          : [],
      });
      onClose();
    }
  };

  const handleBack = () => {
    if (step > 1) setStep(step - 1);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      setFormData((prev) => ({ ...prev, resumeFile: file, isParsing: true }));

      // Simulate AI Parsing
      setTimeout(() => {
        setFormData((prev) => ({
          ...prev,
          isParsing: false,
          headline: "Operations Director & Supply Chain Leader",
          desiredRole: "Director of Logistics",
        }));
      }, 1800);
    }
  };

  const modalActions = (
    <div className="flex items-center justify-between w-full">
      <button
        onClick={handleBack}
        disabled={step === 1}
        className={`px-4 py-2 rounded-xl text-xs font-bold border border-border flex items-center gap-1 cursor-pointer transition-all ${
          step === 1
            ? "opacity-30 cursor-not-allowed"
            : "hover:bg-secondary text-foreground"
        }`}
      >
        <ArrowLeft className="h-3.5 w-3.5" />
        Back
      </button>

      <button
        onClick={handleNext}
        disabled={formData.isParsing}
        className="px-5 py-2 bg-primary hover:bg-primary/95 text-primary-foreground text-xs font-bold rounded-xl shadow-md flex items-center gap-1.5 cursor-pointer transition-all"
      >
        {step === 3 ? "Complete" : "Continue"}
        <ArrowRight className="h-3.5 w-3.5" />
      </button>
    </div>
  );

  return (
    <MuiModal
      open={isOpen}
      onClose={onClose}
      title="AI Career Setup"
      actions={modalActions}
      maxWidth="sm"
    >
      <div className="flex flex-col gap-5">
        {/* Step Indicator */}
        <div className="flex px-4 py-2.5 bg-secondary/50 border border-border/40 rounded-xl text-xs gap-3">
          {[
            { n: 1, label: "Basics" },
            { n: 2, label: "Preferences" },
            { n: 3, label: "Resume Upload" },
          ].map((s) => (
            <div key={s.n} className="flex items-center gap-1.5 flex-1">
              <div
                className={`h-5 w-5 rounded-full flex items-center justify-center font-bold text-[10px] ${
                  step === s.n
                    ? "bg-primary text-primary-foreground"
                    : step > s.n
                      ? "bg-emerald-500 text-white"
                      : "bg-border text-muted-foreground"
                }`}
              >
                {step > s.n ? <Check className="h-3 w-3" /> : s.n}
              </div>
              <span
                className={`font-semibold ${step === s.n ? "text-foreground" : "text-muted-foreground"}`}
              >
                {s.label}
              </span>
            </div>
          ))}
        </div>

        {/* Content Body */}
        <div className="min-h-[280px] relative z-10">
          <AnimatePresence mode="wait">
            {step === 1 && (
              <motion.div
                key="step1"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-4"
              >
                <div>
                  <h3 className="text-sm font-bold tracking-tight text-foreground">
                    Welcome to Kirmya
                  </h3>
                  <p className="text-xs text-muted-foreground">
                    Let&apos;s build your professional identity. Start with your
                    basics.
                  </p>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="col-span-2 sm:col-span-1">
                    <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1.5">
                      Full Name
                    </label>
                    <input
                      type="text"
                      value={formData.fullName}
                      onChange={(e) =>
                        setFormData((prev) => ({
                          ...prev,
                          fullName: e.target.value,
                        }))
                      }
                      className="w-full bg-secondary/40 border border-border hover:border-primary/40 focus:border-primary rounded-xl px-4 py-2 text-sm outline-none transition-all"
                    />
                  </div>
                  <div className="col-span-2 sm:col-span-1">
                    <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1.5">
                      Preferred Name
                    </label>
                    <input
                      type="text"
                      value={formData.preferredName}
                      onChange={(e) =>
                        setFormData((prev) => ({
                          ...prev,
                          preferredName: e.target.value,
                        }))
                      }
                      className="w-full bg-secondary/40 border border-border hover:border-primary/40 focus:border-primary rounded-xl px-4 py-2 text-sm outline-none transition-all"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1.5">
                    Professional Headline
                  </label>
                  <input
                    type="text"
                    value={formData.headline}
                    onChange={(e) =>
                      setFormData((prev) => ({
                        ...prev,
                        headline: e.target.value,
                      }))
                    }
                    className="w-full bg-secondary/40 border border-border hover:border-primary/40 focus:border-primary rounded-xl px-4 py-2 text-sm outline-none transition-all"
                  />
                  <p className="text-[10px] text-muted-foreground mt-1">
                    AI suggestion: Keep it action-oriented showing domain
                    expertise.
                  </p>
                </div>
              </motion.div>
            )}

            {step === 2 && (
              <motion.div
                key="step2"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-4"
              >
                <div>
                  <h3 className="text-sm font-bold tracking-tight text-foreground">
                    Target Role &amp; Preferences
                  </h3>
                  <p className="text-xs text-muted-foreground">
                    What positions and work styles are you targeting next?
                  </p>
                </div>
                <div>
                  <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1.5">
                    Desired Role
                  </label>
                  <input
                    type="text"
                    value={formData.desiredRole}
                    onChange={(e) =>
                      setFormData((prev) => ({
                        ...prev,
                        desiredRole: e.target.value,
                      }))
                    }
                    className="w-full bg-secondary/40 border border-border hover:border-primary/40 focus:border-primary rounded-xl px-4 py-2 text-sm outline-none transition-all"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1.5">
                      Min Annual Salary (USD)
                    </label>
                    <input
                      type="number"
                      value={formData.salaryMin}
                      onChange={(e) =>
                        setFormData((prev) => ({
                          ...prev,
                          salaryMin: parseInt(e.target.value) || 0,
                        }))
                      }
                      className="w-full bg-secondary/40 border border-border hover:border-primary/40 focus:border-primary rounded-xl px-4 py-2 text-sm outline-none transition-all"
                    />
                  </div>
                  <div>
                    <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1.5">
                      Work Preference
                    </label>
                    <select
                      value={formData.workMode}
                      onChange={(e) =>
                        setFormData((prev) => ({
                          ...prev,
                          workMode: e.target.value,
                        }))
                      }
                      className="w-full bg-secondary/40 border border-border focus:border-primary rounded-xl px-4 py-2 text-sm outline-none transition-all"
                    >
                      <option value="remote">Remote</option>
                      <option value="hybrid">Hybrid</option>
                      <option value="onsite">Onsite</option>
                    </select>
                  </div>
                </div>
              </motion.div>
            )}

            {step === 3 && (
              <motion.div
                key="step3"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                className="space-y-4"
              >
                <div>
                  <h3 className="text-sm font-bold tracking-tight text-foreground">
                    Import Resume &amp; Documents
                  </h3>
                  <p className="text-xs text-muted-foreground">
                    Upload your current CV. Kirmya AI will parse and
                    pre-populate your work timeline.
                  </p>
                </div>

                <div className="border-2 border-dashed border-border hover:border-primary/50 rounded-2xl p-6 flex flex-col items-center justify-center gap-3 bg-secondary/20 transition-all relative">
                  <input
                    type="file"
                    id="resume-upload"
                    className="absolute inset-0 opacity-0 cursor-pointer"
                    onChange={handleFileChange}
                    accept=".pdf,.docx"
                  />
                  <div className="h-12 w-12 rounded-full bg-primary/5 flex items-center justify-center text-primary">
                    <Upload className="h-6 w-6" />
                  </div>
                  <div className="text-center">
                    <p className="text-sm font-semibold">
                      {formData.resumeFile
                        ? formData.resumeFile.name
                        : "Drag & drop your resume here"}
                    </p>
                    <p className="text-xs text-muted-foreground mt-1">
                      Supports PDF or DOCX up to 10 MB
                    </p>
                  </div>
                </div>

                {formData.isParsing && (
                  <div className="bg-primary/5 border border-primary/20 p-4 rounded-xl flex items-center gap-3 text-xs text-primary">
                    <Loader2 className="h-4 w-4 animate-spin shrink-0" />
                    <span>
                      Analyzing document sections, extracting entities, and
                      calculating ATS readiness...
                    </span>
                  </div>
                )}

                {formData.resumeFile && !formData.isParsing && (
                  <div className="bg-emerald-500/10 border border-emerald-500/20 p-4 rounded-xl flex items-center gap-3 text-xs text-emerald-600 dark:text-emerald-400">
                    <Check className="h-4 w-4 shrink-0" />
                    <span>
                      Resume parsed! Headline and desired role suggested. Click
                      Complete to apply changes.
                    </span>
                  </div>
                )}
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </MuiModal>
  );
}
