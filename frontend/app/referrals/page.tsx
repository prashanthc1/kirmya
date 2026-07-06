"use client";

import React, { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { 
  Compass, 
  CheckCircle2, 
  Clock, 
  XCircle, 
  Send, 
  PlusCircle, 
  Loader2,
  Building,
  ArrowRight
} from "lucide-react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api, ApiError } from "@/lib/api/client";

interface Referral {
  id: string;
  company: string;
  message: string;
  status: string;
  outcome: string;
  created_at: string;
}

export default function ReferralsPage() {
  const [referrals, setReferrals] = useState<Referral[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // New Request Form states
  const [showForm, setShowForm] = useState(false);
  const [companyName, setCompanyName] = useState("");
  const [referralMessage, setReferralMessage] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const fetchReferrals = async () => {
    try {
      const data = await api.get<{ referrals: Referral[] }>("/referrals/outgoing");
      setReferrals(data?.referrals ?? []);
    } catch (err) {
      setError(
        err instanceof ApiError
          ? err.message
          : "Could not load your referrals. Please try again."
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReferrals();
  }, []);

  const handleCreateReferral = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!companyName.trim() || submitting) return;
    setSubmitting(true);
    try {
      await api.post("/referrals", {
        company: companyName,
        message: referralMessage,
      });
      setCompanyName("");
      setReferralMessage("");
      setShowForm(false);
      await fetchReferrals();
    } catch (err) {
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  };

  const getStatusConfig = (status: string) => {
    switch (status.toLowerCase()) {
      case "accepted":
        return {
          bg: "bg-emerald-500/10 border-emerald-500/20 text-emerald-500",
          icon: CheckCircle2,
          label: "Accepted"
        };
      case "declined":
        return {
          bg: "bg-destructive/10 border-destructive/20 text-destructive",
          icon: XCircle,
          label: "Declined"
        };
      default:
        return {
          bg: "bg-amber-500/10 border-amber-500/20 text-amber-600 dark:text-amber-400",
          icon: Clock,
          label: "Pending"
        };
    }
  };

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Referrals" }]} />

      <main className="flex-grow max-w-4xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8 space-y-6">
        
        {/* Header section */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
          <div className="space-y-1">
            <span className="text-xs font-bold uppercase tracking-widest text-primary">Warm Connections</span>
            <h1 className="text-3xl font-extrabold tracking-tight">Referral Requests</h1>
            <p className="text-sm text-muted-foreground">Track the referral requests you&apos;ve sent and where each one stands.</p>
          </div>

          <button
            onClick={() => setShowForm(!showForm)}
            className="px-5 py-2.5 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-xs font-bold shadow-sm shadow-blue-500/10 flex items-center gap-1.5 shrink-0"
          >
            <PlusCircle className="h-4 w-4" />
            Request Referral
          </button>
        </div>

        {/* Form panel */}
        <AnimatePresence>
          {showForm && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              className="overflow-hidden"
            >
              <form onSubmit={handleCreateReferral} className="bg-card border border-border/80 p-6 rounded-3xl space-y-4 shadow-sm">
                <h3 className="text-sm font-bold">New Referral Request</h3>
                
                <div className="grid grid-cols-1 gap-4">
                  <div className="space-y-1.5">
                    <label className="text-xs font-semibold text-muted-foreground block">Target Company</label>
                    <input
                      type="text"
                      placeholder="e.g. Stripe, Linear..."
                      value={companyName}
                      onChange={(e) => setCompanyName(e.target.value)}
                      required
                      className="w-full px-4 py-2 rounded-full border border-border/80 bg-background text-foreground text-sm focus:outline-none focus:ring-1 focus:ring-primary"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <label className="text-xs font-semibold text-muted-foreground block">Message / Note to employee</label>
                    <textarea
                      placeholder="Explain your fit, linking your resume or portfolio..."
                      value={referralMessage}
                      onChange={(e) => setReferralMessage(e.target.value)}
                      rows={3}
                      className="w-full px-4 py-3 rounded-2xl border border-border/80 bg-background text-foreground text-sm focus:outline-none focus:ring-1 focus:ring-primary"
                    />
                  </div>
                </div>

                <div className="flex gap-2 justify-end">
                  <button
                    type="button"
                    onClick={() => setShowForm(false)}
                    className="px-4 py-2 rounded-full border border-border hover:bg-secondary text-xs font-bold"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={submitting || !companyName.trim()}
                    className="px-5 py-2 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-xs font-bold shadow-sm flex items-center gap-1"
                  >
                    {submitting && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
                    Submit Request
                  </button>
                </div>
              </form>
            </motion.div>
          )}
        </AnimatePresence>

        {error && (
          <div className="p-4 bg-destructive/10 border border-destructive/20 rounded-2xl text-destructive text-xs font-medium">
            {error}
          </div>
        )}

        {/* Content list */}
        <div className="space-y-4">
          {loading ? (
            <div className="flex flex-col items-center justify-center py-20 gap-2">
              <Loader2 className="h-6 w-6 text-primary animate-spin" />
              <span className="text-xs text-muted-foreground">Checking referral pipelines...</span>
            </div>
          ) : referrals.length === 0 ? (
            <div className="text-center py-16 border border-dashed border-border/80 rounded-3xl p-8 bg-secondary/15 space-y-3">
              <Compass className="h-10 w-10 text-muted-foreground mx-auto" />
              <h3 className="text-base font-bold">No referrals requested</h3>
              <p className="text-xs text-muted-foreground max-w-sm mx-auto">
                Request connections at hiring companies to skip the screening queue. You can send a request from the button above.
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-3">
              {referrals.map((ref) => {
                const conf = getStatusConfig(ref.status);
                const Icon = conf.icon;
                return (
                  <div key={ref.id} className="bg-card border border-border/60 p-5 rounded-2xl flex flex-col gap-3 hover:border-border transition-all">
                    <div className="flex items-center justify-between gap-4">
                      <div className="flex items-center gap-2.5">
                        <div className="h-8 w-8 rounded-lg bg-primary/10 border border-primary/20 flex items-center justify-center text-primary shrink-0">
                          <Building className="h-4.5 w-4.5" />
                        </div>
                        <span className="text-sm font-bold text-foreground">{ref.company}</span>
                      </div>
                      
                      <span className={`px-3 py-1 rounded-full text-xs font-semibold border flex items-center gap-1 shrink-0 ${conf.bg}`}>
                        <Icon className="h-3.5 w-3.5" />
                        {conf.label}
                      </span>
                    </div>

                    {ref.message && (
                      <p className="text-xs text-muted-foreground leading-relaxed pl-1">
                        {ref.message}
                      </p>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
