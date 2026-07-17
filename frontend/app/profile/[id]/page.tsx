"use client";

import React, { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  User,
  MapPin,
  Clock,
  CheckCircle,
  UserPlus,
  UserCheck,
  X,
  Loader2,
  ArrowLeft,
  Mail,
  Calendar,
  Sparkles,
  Bookmark,
  Layers,
  GraduationCap,
  Globe,
} from "lucide-react";
import { motion } from "framer-motion";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { profileClient, Profile, PublicProfileResponse } from "@/lib/api/profile";
import { networkClient, ConnectionStatusResponse } from "@/lib/api/network";
import { ApiError } from "@/lib/api/client";
import { useConnectionStatus } from "@/hooks/useConnections";
import ConnectButton from "@/components/connections/ConnectButton";
import MutualConnectionsStrip from "@/components/connections/MutualConnectionsStrip";

export default function OtherProfilePage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const [profile, setProfile] = useState<PublicProfileResponse | null>(null);
  const [currentUserID, setCurrentUserID] = useState<string | null>(null);
  const { data: connectionStatus } = useConnectionStatus(id);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    let active = true;
    (async () => {
      try {
        const me = await profileClient.getMe();
        if (active) setCurrentUserID(me.user_id);

        if (me.user_id === id) {
          router.replace("/profile");
          return;
        }

        const data = await profileClient.getPublicProfile(id);
        if (active) setProfile(data);
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError
              ? err.message
              : "Could not load profile. It might be private or not exist.",
          );
        }
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [id, router]);

  if (loading) {
    return (
      <div className="min-h-screen bg-background text-foreground flex flex-col">
        <SiteNav
          breadcrumb={[{ label: "Home", href: "/" }, { label: "Profile" }]}
        />
        <div className="flex-grow flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 className="h-8 w-8 text-primary animate-spin" />
          <span className="text-sm font-semibold text-muted-foreground">
            Fetching professional timeline...
          </span>
        </div>
        <SiteFooter />
      </div>
    );
  }

  if (error || !profile) {
    return (
      <div className="min-h-screen bg-background text-foreground flex flex-col">
        <SiteNav
          breadcrumb={[{ label: "Home", href: "/" }, { label: "Profile" }]}
        />
        <main className="flex-grow max-w-lg mx-auto w-full px-4 py-20">
          <div className="p-6 bg-destructive/10 border border-destructive/20 rounded-3xl text-destructive text-center space-y-3">
            <p className="text-sm font-bold">
              {error || "Profile could not be loaded."}
            </p>
            <button
              onClick={() => router.back()}
              className="px-4 py-2 bg-destructive text-destructive-foreground rounded-full text-xs font-bold"
            >
              Go Back
            </button>
          </div>
        </main>
        <SiteFooter />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav
        breadcrumb={[
          { label: "Home", href: "/" },
          { label: "Profiles", href: "/search?type=user" },
          { label: profile.headline || "Profile" },
        ]}
      />

      <main className="flex-grow max-w-4xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8 space-y-6">
        {/* Core Profile Header */}
        <div className="bg-card border border-border/80 p-6 md:p-8 rounded-3xl shadow-sm space-y-6">
          <div className="flex flex-col md:flex-row items-center md:items-start gap-6">
            <div className="h-20 w-20 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-black text-2xl uppercase select-none shrink-0">
              {profile.headline?.charAt(0) || "P"}
            </div>

            <div className="flex-grow space-y-2 text-center md:text-left">
              <div className="flex flex-col md:flex-row items-center gap-2 justify-center md:justify-start">
                <h1 className="text-xl md:text-2xl font-extrabold tracking-tight">
                  {profile.headline || "Kirmya Professional"}
                </h1>
                {profile.career_status && (
                  <span className="px-2.5 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-[10px] font-extrabold text-emerald-500 uppercase tracking-widest">
                    {profile.career_status.replace("_", " ")}
                  </span>
                )}
              </div>

              <p className="text-sm font-semibold text-muted-foreground">
                {profile.headline || "Career Transition & Development"}
              </p>

              <p className="text-xs text-muted-foreground leading-relaxed max-w-xl mx-auto md:mx-0">
                {profile.bio ||
                  "Active community member focused on professional career comeback."}
              </p>

              <div className="flex flex-wrap items-center justify-center md:justify-start gap-2 pt-2 text-xs">
                {profile.location && (
                  <span className="px-3 py-1 bg-secondary text-muted-foreground rounded-full border border-border/40 flex items-center gap-1">
                    <MapPin className="h-3.5 w-3.5" />
                    {profile.location}
                  </span>
                )}
                {profile.open_to_remote && (
                  <span className="px-3 py-1 bg-secondary text-muted-foreground rounded-full border border-border/40">
                    💻 Open to Remote
                  </span>
                )}
                {profile.willing_to_mentor && (
                  <span className="px-3 py-1 bg-secondary text-muted-foreground rounded-full border border-border/40">
                    🤝 Willing to Mentor
                  </span>
                )}
              </div>
            </div>
          </div>

          {/* Action Row */}
          <div className="border-t border-border/40 pt-4 flex items-center justify-between gap-4 flex-wrap">
            <span className="text-[11px] font-bold text-muted-foreground uppercase tracking-widest">
              Messaging Gate status
            </span>
            <div className="flex items-center gap-3">
              <MutualConnectionsStrip userId={id} />
              <ConnectButton
                targetUserId={id}
                targetUserName={profile.headline || "User"}
                currentConnectionStatus={
                  connectionStatus?.status === "pending"
                    ? (connectionStatus.requested_by === currentUserID ? "pending_outgoing" : "pending_incoming")
                    : (connectionStatus?.status || "none") as any
                }
                connectionId={connectionStatus?.connection_id}
              />
            </div>
          </div>
        </div>

        {/* Details Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 items-start">
          {/* Left Column: Timeline & About */}
          <div className="md:col-span-2 space-y-8">
            {/* About */}
            {profile.about && (
              <div className="space-y-3">
                <h2 className="text-lg font-bold text-foreground">About</h2>
                <p className="text-sm leading-relaxed text-muted-foreground whitespace-pre-wrap">
                  {profile.about}
                </p>
              </div>
            )}

            {/* Experience timeline */}
            {profile.experiences && profile.experiences.length > 0 && (
              <div className="space-y-4">
                <h2 className="text-lg font-bold text-foreground">
                  Experience
                </h2>
                <div className="space-y-6">
                  {profile.experiences.map((exp, index) => (
                    <div key={exp.id || index} className="flex gap-4">
                      <div className="flex flex-col items-center shrink-0">
                        <div
                          className={`h-4 w-4 rounded-full border-2 ${
                            index === 0
                              ? "border-primary bg-primary"
                              : "border-border bg-card"
                          }`}
                        />
                        {index < (profile.experiences?.length || 0) - 1 && (
                          <div className="w-[1px] bg-border/80 flex-grow my-1" />
                        )}
                      </div>
                      <div className="space-y-1 pb-2">
                        <h3 className="text-sm font-bold text-foreground">
                          {exp.title}
                        </h3>
                        <p className="text-xs text-muted-foreground">
                          {exp.company} &bull; {exp.start_date} &ndash;{" "}
                          {exp.is_current ? "Present" : exp.end_date}
                        </p>
                        {exp.description && (
                          <p className="text-xs text-muted-foreground/95 leading-relaxed mt-2">
                            {exp.description}
                          </p>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Education timeline */}
            {profile.educations && profile.educations.length > 0 && (
              <div className="space-y-4">
                <h2 className="text-lg font-bold text-foreground">Education</h2>
                <div className="space-y-4">
                  {profile.educations.map((edu, index) => (
                    <div key={edu.id || index} className="space-y-1">
                      <h3 className="text-sm font-bold text-foreground">
                        {edu.school}
                      </h3>
                      <p className="text-xs text-muted-foreground font-medium">
                        {edu.degree} &bull; {edu.field_of_study}
                      </p>
                      <span className="text-[10px] text-muted-foreground block">
                        {edu.start_date} &ndash; {edu.end_date}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Right Column: Sidebar */}
          <div className="space-y-6">
            {/* Skills snap */}
            {profile.skills && profile.skills.length > 0 && (
              <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm space-y-4">
                <h3 className="text-sm font-bold flex items-center gap-1.5">
                  <Layers className="h-4.5 w-4.5 text-primary" />
                  Skills
                </h3>
                <div className="flex flex-wrap gap-1.5">
                  {profile.skills.map((sk) => (
                    <span
                      key={sk.name}
                      className="px-2.5 py-1 bg-secondary text-muted-foreground border border-border/40 rounded-full text-[10px] font-semibold"
                    >
                      {sk.name}
                    </span>
                  ))}
                </div>
              </div>
            )}

            {/* Languages */}
            {profile.languages && profile.languages.length > 0 && (
              <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm space-y-3">
                <h3 className="text-sm font-bold flex items-center gap-1.5">
                  <Globe className="h-4.5 w-4.5 text-primary" />
                  Languages
                </h3>
                <div className="space-y-2 text-xs">
                  {profile.languages.map((l) => (
                    <div
                      key={l.name}
                      className="flex justify-between items-center"
                    >
                      <span className="font-semibold">{l.name}</span>
                      <span className="text-muted-foreground">
                        {l.proficiency}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Contact details card (connections only) */}
            {profile.is_connection && (profile.email || profile.phone || profile.address) && (
              <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm space-y-4">
                <h3 className="text-sm font-bold flex items-center gap-1.5 text-primary">
                  <Mail className="h-4.5 w-4.5" />
                  Contact Information (Connected)
                </h3>
                <div className="space-y-2 text-xs">
                  {profile.email && (
                    <div className="flex flex-col gap-0.5">
                      <span className="text-[10px] text-muted-foreground uppercase font-bold tracking-wider">Email</span>
                      <a href={`mailto:${profile.email}`} className="text-primary hover:underline">{profile.email}</a>
                    </div>
                  )}
                  {profile.phone && (
                    <div className="flex flex-col gap-0.5">
                      <span className="text-[10px] text-muted-foreground uppercase font-bold tracking-wider">Phone</span>
                      <span className="text-foreground">{profile.phone}</span>
                    </div>
                  )}
                  {profile.address && (
                    <div className="flex flex-col gap-0.5">
                      <span className="text-[10px] text-muted-foreground uppercase font-bold tracking-wider">Address</span>
                      <span className="text-foreground">{profile.address}</span>
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
