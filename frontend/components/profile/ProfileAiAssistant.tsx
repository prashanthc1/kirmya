"use client";

import React, { useState, useEffect } from "react";
import {
  Sparkles,
  Send,
  Loader2,
  RefreshCw,
  Check,
  ArrowRight,
  TrendingUp,
  DollarSign,
  BrainCircuit,
  FileText,
  CheckCircle2,
} from "lucide-react";
import { ExtendedProfile } from "./types";

interface ProfileAiAssistantProps {
  profile: ExtendedProfile;
  activeSectionId: string;
  onApplyAiChange: (
    sectionId: string,
    updatedFields: Partial<ExtendedProfile>,
  ) => void;
}

export default function ProfileAiAssistant({
  profile,
  activeSectionId,
  onApplyAiChange,
}: ProfileAiAssistantProps) {
  const [query, setQuery] = useState("");
  const [isGenerating, setIsGenerating] = useState(false);
  const [activeTab, setActiveTab] = useState<"chat" | "insights" | "optimize">(
    "chat",
  );
  const [aiSuggestions, setAiSuggestions] = useState<string | null>(null);
  const [appliedText, setAppliedText] = useState(false);

  // Track messages in chat
  const [messages, setMessages] = useState<
    Array<{
      role: "user" | "assistant";
      content: string;
      actions?: {
        type: string;
        label: string;
        data: Partial<ExtendedProfile>;
      } | null;
    }>
  >([
    {
      role: "assistant",
      content: `Hello ${profile.preferred_name || "Marcus"}! I'm your Kirmya Career Assistant. I've audited your profile and calculated an **88% completeness score**.

Here is what I recommend doing next:
1. Re-format your **Cascade Freight** experience using the **STAR format**.
2. Add your **personal branding statement** to your professional summary.
3. Link your calendar for recruiter booking.
`,
    },
  ]);

  // Handle active section change to suggest relevant context
  useEffect(() => {
    setAppliedText(false);
    if (!activeSectionId) return;

    if (activeSectionId === "experience") {
      setAiSuggestions(`**STAR Bullet Points Suggestion for Cascade Freight**
- **Situation**: Handled mid-mile logistics network across 14 distribution centers.
- **Task**: Consolidate distribution network to cut operations cost.
- **Action**: Spearheaded mathematical location-placement algorithms to merge 14 facilities into 8 hubs, and renegotiated truck contracts.
- **Result**: Reduced overall cost-per-unit by **18%** and yielded **$6.4M** in annual overhead savings.

Would you like to rewrite this achievement in your profile?`);
    } else if (activeSectionId === "summary") {
      setAiSuggestions(`**AI Personal Branding Elevator Pitch**
"Supply chain architect with 22 years of operations leadership. I specialize in turning complex logistics networks from cost-centers into margin-drivers through linear programming and automation. I've saved $6.4M in annual overhead and maintained 99.2% on-time delivery. Let me build your next-generation fulfillment pipeline."`);
    } else {
      setAiSuggestions(null);
    }
  }, [activeSectionId]);

  const handleSend = (textToSend?: string) => {
    const finalQuery = textToSend || query;
    if (!finalQuery.trim() || isGenerating) return;

    // Add user message
    setMessages((prev) => [...prev, { role: "user", content: finalQuery }]);
    if (!textToSend) setQuery("");

    setIsGenerating(true);

    // Simulate AI response
    setTimeout(() => {
      let response = "";
      let actions = null;

      if (
        finalQuery.toLowerCase().includes("salary") ||
        finalQuery.toLowerCase().includes("market")
      ) {
        response = `Based on your role as a **Director of Operations** in **Denver, CO** with **22 years of experience**:

- **Predictive Salary Range**: $195,000 – $240,000 Base
- **Market Demand Score**: Very High (8.8/10)
- **Recruiter Interest**: 42 searches in your area this week.

**High Demand Keywords**: *SAP EWM, S&OP Optimization, Contract Negotiation, Warehouse Robotics*. You have 3/4. Let's add *Warehouse Robotics* to your skills to boost search visibility.`;
      } else if (
        finalQuery.toLowerCase().includes("star") ||
        finalQuery.toLowerCase().includes("rewrite")
      ) {
        response = `I have rewritten your first achievement at **Cascade Freight** to follow the STAR methodology:

*"Spearheaded supply chain network optimization that consolidated 14 distribution facilities into 8 regional hubs, reducing overall cost-per-unit by 18% and saving $6.4M in annual lease and labor overhead."*`;
        actions = {
          type: "experience_rewrite",
          label: "Apply to Cascade Freight",
          data: {
            experiences: profile.experiences.map((e) =>
              e.id === "exp_1"
                ? {
                    ...e,
                    achievements: [
                      "Spearheaded supply chain network optimization that consolidated 14 distribution facilities into 8 regional hubs, reducing overall cost-per-unit by 18% and saving $6.4M in annual lease and labor overhead.",
                      ...e.achievements.slice(1),
                    ],
                  }
                : e,
            ),
          },
        };
      } else {
        response = `I've analyzed your profile sections. To optimize your ATS score, consider adding more quantitative metrics. Under **Education**, I noticed you have research experience. Highlighting the linear programming algorithms in your thesis will attract tech-enabled logistics firms.`;
      }

      setMessages((prev) => [
        ...prev,
        { role: "assistant", content: response, actions },
      ]);
      setIsGenerating(false);
    }, 1200);
  };

  const handleApplyAction = (action: {
    type: string;
    label: string;
    data: Partial<ExtendedProfile>;
  }) => {
    if (action.type === "experience_rewrite") {
      onApplyAiChange("experience", action.data);
      setAppliedText(true);

      setMessages((prev) => [
        ...prev,
        {
          role: "assistant",
          content:
            "✅ Successfully applied the STAR bullet optimization to your Cascade Freight experience!",
        },
      ]);
    }
  };

  const applyContextSuggestion = () => {
    if (activeSectionId === "experience") {
      // Modify first experience
      const updated = {
        experiences: profile.experiences.map((e) =>
          e.id === "exp_1"
            ? {
                ...e,
                achievements: [
                  "Spearheaded mathematical location-placement algorithms to merge 14 facilities into 8 hubs, reducing overall cost-per-unit by 18% and saving $6.4M in annual lease and labor overhead.",
                ],
              }
            : e,
        ),
      };
      onApplyAiChange("experience", updated);
      setAppliedText(true);
    } else if (activeSectionId === "summary") {
      // Modify summary elevator pitch
      onApplyAiChange("summary", {
        elevator_pitch:
          "Supply chain architect with 22 years of operations leadership. I specialize in turning complex logistics networks from cost-centers into margin-drivers through linear programming and automation. I've saved $6.4M in annual overhead and maintained 99.2% on-time delivery. Let me build your next-generation fulfillment pipeline.",
      });
      setAppliedText(true);
    }
  };

  return (
    <aside className="w-full lg:w-90 flex-shrink-0 lg:sticky lg:top-24 space-y-6 self-start">
      <div className="bg-card border border-border/80 rounded-3xl p-6 shadow-sm flex flex-col h-[650px]">
        {/* Header */}
        <div className="flex items-center justify-between pb-4 border-b border-border/40 shrink-0">
          <div className="flex items-center gap-2">
            <div className="h-7 w-7 rounded-lg bg-primary/10 flex items-center justify-center text-primary">
              <BrainCircuit className="h-4 w-4" />
            </div>
            <div>
              <h3 className="text-xs font-bold text-foreground">
                Kirmya Career Assistant
              </h3>
              <p className="text-[10px] text-muted-foreground">Always active</p>
            </div>
          </div>
          <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
        </div>

        {/* Tab Selection */}
        <div className="flex border-b border-border/40 py-2.5 shrink-0 gap-1">
          {[
            { id: "chat", label: "Assistant Chat" },
            { id: "optimize", label: "ATS & Key Optimizer" },
            { id: "insights", label: "Market Insights" },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() =>
                setActiveTab(tab.id as "chat" | "insights" | "optimize")
              }
              className={`flex-1 text-center py-1.5 rounded-lg text-[10px] font-bold tracking-wide uppercase transition-all cursor-pointer ${
                activeTab === tab.id
                  ? "bg-primary/5 text-primary border border-primary/20"
                  : "text-muted-foreground hover:bg-secondary hover:text-foreground border border-transparent"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Tab Content Panels */}
        <div className="flex-grow overflow-y-auto py-4 space-y-4 pr-1">
          {activeTab === "chat" && (
            <>
              {/* Message History */}
              <div className="space-y-4 text-xs">
                {messages.map((m, i) => (
                  <div
                    key={i}
                    className={`flex flex-col ${m.role === "user" ? "items-end" : "items-start"}`}
                  >
                    <div
                      className={`p-3.5 rounded-2xl max-w-[88%] leading-relaxed ${
                        m.role === "user"
                          ? "bg-primary text-primary-foreground rounded-tr-none"
                          : "bg-secondary/60 text-foreground border border-border/40 rounded-tl-none"
                      }`}
                    >
                      <p className="whitespace-pre-wrap">{m.content}</p>

                      {/* Action buttons embedded in message */}
                      {m.actions && !appliedText && (
                        <button
                          onClick={() =>
                            m.actions && handleApplyAction(m.actions)
                          }
                          className="mt-3 w-full bg-primary hover:bg-primary/95 text-primary-foreground font-bold py-2 px-3 rounded-xl flex items-center justify-center gap-1.5 transition-all cursor-pointer"
                        >
                          <span>{m.actions.label}</span>
                          <ArrowRight className="h-3.5 w-3.5" />
                        </button>
                      )}
                    </div>
                  </div>
                ))}

                {isGenerating && (
                  <div className="flex justify-start">
                    <div className="bg-secondary/60 border border-border/40 p-3 rounded-2xl rounded-tl-none flex items-center gap-2">
                      <Loader2 className="h-3.5 w-3.5 animate-spin text-primary" />
                      <span className="text-[10px] font-semibold text-muted-foreground">
                        Drafting recommendations...
                      </span>
                    </div>
                  </div>
                )}
              </div>
            </>
          )}

          {activeTab === "optimize" && (
            <div className="space-y-4 text-xs">
              <div className="bg-primary/5 border border-primary/10 p-4 rounded-2xl space-y-2">
                <h5 className="font-bold flex items-center gap-1.5 text-primary">
                  <Sparkles className="h-4 w-4" />
                  Contextual Optimizer
                </h5>
                <p className="text-muted-foreground leading-normal">
                  Select any section card on the center workspace to run
                  automatic optimizations.
                </p>
                {activeSectionId ? (
                  <div className="bg-card border border-border/60 p-2.5 rounded-xl text-[10px] font-bold text-foreground">
                    ⚡ Currently inspecting:{" "}
                    <span className="text-primary capitalize">
                      {activeSectionId}
                    </span>
                  </div>
                ) : (
                  <div className="bg-card border border-border/60 p-2.5 rounded-xl text-[10px] font-semibold text-muted-foreground">
                    Click a section card to begin
                  </div>
                )}
              </div>

              {aiSuggestions ? (
                <div className="border border-border/80 p-4 rounded-2xl space-y-3 bg-secondary/35">
                  <p className="font-semibold text-foreground border-b border-border/40 pb-1.5 capitalize">
                    Suggestions for {activeSectionId}
                  </p>
                  <p className="text-muted-foreground whitespace-pre-wrap leading-relaxed">
                    {aiSuggestions}
                  </p>

                  {!appliedText ? (
                    <button
                      onClick={applyContextSuggestion}
                      className="w-full bg-primary hover:bg-primary/95 text-primary-foreground font-bold py-2 px-3 rounded-xl flex items-center justify-center gap-1.5 transition-all cursor-pointer"
                    >
                      <Check className="h-4 w-4" />
                      Apply to Profile
                    </button>
                  ) : (
                    <div className="bg-emerald-500/10 border border-emerald-500/20 p-2.5 rounded-xl text-[10px] font-bold text-emerald-600 dark:text-emerald-400 flex items-center gap-1.5 justify-center">
                      <CheckCircle2 className="h-4 w-4" />
                      Applied Successfully!
                    </div>
                  )}
                </div>
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <p>
                    Choose &quot;Professional Summary&quot; or &quot;Work
                    Experience&quot; in the main center panel to generate quick
                    copy.
                  </p>
                </div>
              )}
            </div>
          )}

          {activeTab === "insights" && (
            <div className="space-y-4 text-xs">
              {/* Analytics card */}
              <div className="bg-card border border-border/80 p-4 rounded-2xl space-y-3 shadow-sm">
                <h5 className="font-bold flex items-center gap-1 text-foreground">
                  <TrendingUp className="h-4 w-4 text-primary" />
                  Denver Logistics Market
                </h5>
                <div className="grid grid-cols-2 gap-3 pt-2 text-center">
                  <div className="bg-secondary/40 border border-border/40 p-2 rounded-xl">
                    <p className="text-[10px] text-muted-foreground font-bold uppercase tracking-wider">
                      Median Comp
                    </p>
                    <p className="text-sm font-black text-foreground mt-0.5">
                      $215,000
                    </p>
                  </div>
                  <div className="bg-secondary/40 border border-border/40 p-2 rounded-xl">
                    <p className="text-[10px] text-muted-foreground font-bold uppercase tracking-wider">
                      Demand Level
                    </p>
                    <p className="text-sm font-black text-emerald-500 mt-0.5">
                      High (+12%)
                    </p>
                  </div>
                </div>
              </div>

              {/* Recommended certifications */}
              <div className="bg-card border border-border/80 p-4 rounded-2xl space-y-3 shadow-sm">
                <h5 className="font-bold text-foreground">
                  Top Recruiter Search Keywords
                </h5>
                <div className="flex flex-wrap gap-1.5">
                  {[
                    "S&OP Modeling",
                    "Warehouse Layouts",
                    "Contract Negotiation",
                    "WMS Architect",
                    "Linear Programming",
                    "OSHA Compliance",
                  ].map((kw) => (
                    <span
                      key={kw}
                      className="px-2 py-1 bg-secondary text-muted-foreground rounded-md text-[10px] font-semibold border border-border/40"
                    >
                      {kw}
                    </span>
                  ))}
                </div>
              </div>

              {/* Salary Suggestion trigger */}
              <div className="bg-gradient-to-r from-blue-600/10 to-indigo-600/10 border border-blue-500/20 p-4 rounded-2xl space-y-2">
                <h5 className="font-bold text-primary flex items-center gap-1">
                  <DollarSign className="h-4 w-4" />
                  AI Salary Benchmarking
                </h5>
                <p className="text-muted-foreground text-[11px] leading-relaxed">
                  Run a custom market analytics report comparing your
                  qualifications against open requisitions at major logistics
                  and e-commerce companies.
                </p>
                <button
                  onClick={() =>
                    handleSend(
                      "Analyze my target market salary and recruiter demand",
                    )
                  }
                  className="w-full text-center py-2 bg-primary hover:bg-primary/95 text-primary-foreground font-bold rounded-xl transition-all cursor-pointer"
                >
                  Generate Report
                </button>
              </div>
            </div>
          )}
        </div>

        {/* Query Input Box */}
        <div className="pt-3 border-t border-border/40 shrink-0">
          <div className="flex bg-secondary/80 border border-border/80 hover:border-primary/40 focus-within:border-primary rounded-2xl px-3.5 py-2.5 items-center gap-2.5 transition-all">
            <input
              type="text"
              placeholder="Ask AI Copilot (e.g. 'Optimize my resume'...)"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSend()}
              className="bg-transparent border-none text-xs text-foreground placeholder-muted-foreground outline-none flex-grow"
            />
            <button
              onClick={() => handleSend()}
              className="h-7 w-7 rounded-xl bg-primary hover:bg-primary/95 text-primary-foreground flex items-center justify-center shrink-0 transition-colors cursor-pointer"
            >
              <Send className="h-3.5 w-3.5" />
            </button>
          </div>

          {/* Quick chip queries */}
          <div className="flex gap-1.5 mt-2 overflow-x-auto py-1.5 no-scrollbar">
            {[
              "Optimize for ATS",
              "Rewrite Cascade Freight",
              "Salary prediction",
            ].map((chip) => (
              <button
                key={chip}
                onClick={() => handleSend(chip)}
                className="px-2.5 py-1 bg-secondary text-muted-foreground border border-border/40 hover:bg-primary/5 hover:text-primary rounded-lg text-[9px] font-bold tracking-wide uppercase shrink-0 transition-all cursor-pointer"
              >
                {chip}
              </button>
            ))}
          </div>
        </div>
      </div>
    </aside>
  );
}
