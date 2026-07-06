"use client";

import React, { useState, useEffect, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Sparkles, CloudLightning, ShieldCheck, Check, Clock, AlertCircle } from "lucide-react";
import { ExtendedProfile } from "./types";
import { MOCK_PROFILE } from "./mockData";
import { profileClient } from "@/lib/api/profile";
import ProfileLeftSidebar from "./ProfileLeftSidebar";
import ProfileCenterWorkspace from "./ProfileCenterWorkspace";
import ProfileAiAssistant from "./ProfileAiAssistant";
import ProfileOnboarding from "./ProfileOnboarding";

export default function ProfileWorkspace() {
  const [profile, setProfile] = useState<ExtendedProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeSectionId, setActiveSectionId] = useState("identity");
  const [saveStatus, setSaveStatus] = useState<"idle" | "saving" | "saved" | "error">("idle");
  const [isCloudSynced, setIsCloudSynced] = useState(false);
  const [isOnboardingOpen, setIsOnboardingOpen] = useState(false);

  // Undo/Redo Stacks
  const [historyStack, setHistoryStack] = useState<ExtendedProfile[]>([]);
  const [redoStack, setRedoStack] = useState<ExtendedProfile[]>([]);
  const saveTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const triggerAutosave = React.useCallback((dataToSave: ExtendedProfile) => {
    setSaveStatus("saving");
    
    if (saveTimeoutRef.current) {
      clearTimeout(saveTimeoutRef.current);
    }

    saveTimeoutRef.current = setTimeout(async () => {
      try {
        // Calculate completeness score
        const filledFields = Object.keys(dataToSave).filter(k => !!(dataToSave as unknown as Record<string, unknown>)[k]).length;
        const totalFields = 20; // estimate
        dataToSave.profile_completeness_score = Math.min(98, Math.round((filledFields / totalFields) * 100));

        // Save locally
        localStorage.setItem("kirmya_profile_data", JSON.stringify(dataToSave));
        
        // Attempt cloud sync if authenticated
        if (isCloudSynced) {
          await profileClient.updateMe(dataToSave);
        }
        
        setSaveStatus("saved");
      } catch (err) {
        console.error("Autosave sync failed", err);
        setSaveStatus("error");
      }
    }, 700);
  }, [isCloudSynced]);

  // Push state to history stack for Undo
  const pushToHistory = React.useCallback((currentState: ExtendedProfile) => {
    setHistoryStack(prev => [...prev.slice(-19), currentState]); // Limit to 20 states
    setRedoStack([]); // Clear redo stack on new action
  }, []);

  const handleUndo = React.useCallback(() => {
    if (historyStack.length === 0 || !profile) return;
    const prev = historyStack[historyStack.length - 1];
    setHistoryStack(prevHistory => prevHistory.slice(0, -1));
    setRedoStack(prevRedo => [...prevRedo, profile]);
    
    setProfile(prev);
    triggerAutosave(prev);
  }, [historyStack, profile, triggerAutosave]);

  const handleRedo = React.useCallback(() => {
    if (redoStack.length === 0 || !profile) return;
    const next = redoStack[redoStack.length - 1];
    setRedoStack(prevRedo => prevRedo.slice(0, -1));
    setHistoryStack(prevHistory => [...prevHistory, profile]);
    
    setProfile(next);
    triggerAutosave(next);
  }, [redoStack, profile, triggerAutosave]);

  // Load profile data
  useEffect(() => {
    let active = true;
    const fetchProfile = async () => {
      try {
        setLoading(true);
        // Try calling the actual Kirmya backend
        const response = await profileClient.getMe();
        if (active && response) {
          // Merge response into extended structure
          const extended = { ...MOCK_PROFILE, ...response } as ExtendedProfile;
          setProfile(extended);
          setIsCloudSynced(true);
          setSaveStatus("saved");
        }
      } catch (err) {
        console.warn("Kirmya API offline or unauthorized. Falling back to Local Storage flow.", err);
        // Offline / Local storage fallback
        if (active) {
          const stored = localStorage.getItem("kirmya_profile_data");
          if (stored) {
            setProfile(JSON.parse(stored));
          } else {
            // Seed default mockup profile
            setProfile(MOCK_PROFILE);
            localStorage.setItem("kirmya_profile_data", JSON.stringify(MOCK_PROFILE));
          }
          setIsCloudSynced(false);
          setSaveStatus("saved");
        }
      } finally {
        if (active) setLoading(false);
      }
    };

    fetchProfile();
    return () => {
      active = false;
    };
  }, []);

  // Keyboard shortcut listener (Undo/Redo)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const isCmdOrCtrl = e.metaKey || e.ctrlKey;
      if (isCmdOrCtrl && e.key.toLowerCase() === "z") {
        e.preventDefault();
        handleUndo();
      } else if (isCmdOrCtrl && e.key.toLowerCase() === "y") {
        e.preventDefault();
        handleRedo();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [handleUndo, handleRedo]);

  const handleUpdateField = (updatedFields: Partial<ExtendedProfile>) => {
    if (!profile) return;
    pushToHistory(profile);

    const nextState = { ...profile, ...updatedFields } as ExtendedProfile;
    setProfile(nextState);
    triggerAutosave(nextState);
  };

  const handleApplyAiChange = (sectionId: string, updatedFields: Partial<ExtendedProfile>) => {
    if (!profile) return;
    pushToHistory(profile);

    const nextState = { ...profile, ...updatedFields } as ExtendedProfile;
    setProfile(nextState);
    triggerAutosave(nextState);
  };

  const handleReset = () => {
    if (!window.confirm("Are you sure you want to reset your profile to default settings? This will clear current changes.")) return;
    if (profile) pushToHistory(profile);
    setProfile(MOCK_PROFILE);
    triggerAutosave(MOCK_PROFILE);
  };

  const handleScrollToSection = (sectionId: string) => {
    const el = document.getElementById(`section-${sectionId}`);
    if (el) {
      el.scrollIntoView({ behavior: "smooth", block: "center" });
      setActiveSectionId(sectionId);
    }
  };

  const handleOnboardingComplete = (data: Partial<ExtendedProfile>) => {
    handleUpdateField(data);
  };

  if (loading || !profile) {
    return (
      <div className="max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-10 space-y-6">
        <div className="h-10 bg-secondary/50 rounded-2xl animate-pulse w-1/4" />
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          <div className="lg:col-span-3 h-96 bg-secondary/35 rounded-3xl animate-pulse" />
          <div className="lg:col-span-6 space-y-4">
            <div className="h-40 bg-secondary/35 rounded-3xl animate-pulse" />
            <div className="h-40 bg-secondary/35 rounded-3xl animate-pulse" />
          </div>
          <div className="lg:col-span-3 h-96 bg-secondary/35 rounded-3xl animate-pulse" />
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8 space-y-6 relative z-10">
      
      {/* Save status header indicator */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 border-b border-border/40 pb-5">
        <div>
          <h1 className="text-xl md:text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
            Kirmya Profile Workspace
            <span className="text-[10px] uppercase font-bold tracking-widest text-primary bg-primary/10 px-2 py-0.5 rounded-full">
              AI Copilot Active
            </span>
          </h1>
          <p className="text-xs text-muted-foreground mt-0.5">Redesign and optimize your professional identity with live suggestions.</p>
        </div>

        {/* Save Pill */}
        <div className="flex items-center gap-3">
          <AnimatePresence mode="wait">
            {saveStatus === "saving" && (
              <motion.div 
                initial={{ opacity: 0, y: -5 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: 5 }}
                className="bg-secondary text-muted-foreground border border-border/80 px-3.5 py-1.5 rounded-full text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5 shadow-sm"
              >
                <Clock className="h-3.5 w-3.5 animate-spin text-primary" />
                Saving changes...
              </motion.div>
            )}
            {saveStatus === "saved" && (
              <motion.div 
                initial={{ opacity: 0, y: -5 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: 5 }}
                className="bg-emerald-500/5 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 px-3.5 py-1.5 rounded-full text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5 shadow-sm"
              >
                <Check className="h-3.5 w-3.5" />
                {isCloudSynced ? "Synced to cloud" : "Autosaved locally"}
              </motion.div>
            )}
            {saveStatus === "error" && (
              <motion.div 
                initial={{ opacity: 0, y: -5 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: 5 }}
                className="bg-destructive/5 text-destructive border border-destructive/20 px-3.5 py-1.5 rounded-full text-[10px] font-bold tracking-wider uppercase flex items-center gap-1.5 shadow-sm"
              >
                <AlertCircle className="h-3.5 w-3.5" />
                Save error (Offline)
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>

      {/* Main Grid Layout */}
      <div className="flex flex-col lg:flex-row gap-8 items-start">
        
        {/* Left Sidebar */}
        <ProfileLeftSidebar
          profile={profile}
          canUndo={historyStack.length > 0}
          canRedo={redoStack.length > 0}
          onUndo={handleUndo}
          onRedo={handleRedo}
          onReset={handleReset}
          onOpenOnboarding={() => setIsOnboardingOpen(true)}
          onScrollToSection={handleScrollToSection}
        />

        {/* Center Workspace */}
        <ProfileCenterWorkspace
          profile={profile}
          activeSectionId={activeSectionId}
          setActiveSectionId={setActiveSectionId}
          onUpdateField={handleUpdateField}
        />

        {/* Right AI Assistant */}
        <ProfileAiAssistant
          profile={profile}
          activeSectionId={activeSectionId}
          onApplyAiChange={handleApplyAiChange}
        />
      </div>

      {/* Onboarding Dialog */}
      <AnimatePresence>
        {isOnboardingOpen && (
          <ProfileOnboarding
            isOpen={isOnboardingOpen}
            onClose={() => setIsOnboardingOpen(false)}
            onComplete={handleOnboardingComplete}
          />
        )}
      </AnimatePresence>
    </div>
  );
}
