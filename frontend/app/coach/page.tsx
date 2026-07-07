"use client";

import React, { useState, useRef, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  Sparkles,
  Send,
  BrainCircuit,
  Terminal,
  DollarSign,
  FileText,
  GraduationCap,
  CheckCircle,
  Cpu,
  Loader2,
  ArrowRight,
  TrendingUp,
  Briefcase,
} from "lucide-react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

interface Message {
  id: string;
  sender: "user" | "coach";
  text: string;
  timestamp: string;
}

const STARTER_PROMPTS = [
  {
    icon: BrainCircuit,
    title: "Layoff Explanation",
    desc: "Rehearse how to explain a career gap or layoff positively.",
    prompt:
      "I want to practice explaining my recent layoff to a recruiter without sounding nervous.",
  },
  {
    icon: DollarSign,
    title: "Salary Negotiation",
    desc: "Draft a word-for-word response for salary counters.",
    prompt:
      "Can you help me draft a counter-offer script for a Senior Operations role?",
  },
  {
    icon: Briefcase,
    title: "Interview Rehearsal",
    desc: "Practice mock questions for supply chain and engineering.",
    prompt:
      "Quiz me on S&OP or System Design questions for my upcoming interview.",
  },
  {
    icon: FileText,
    title: "Cold Referral Draft",
    desc: "Generate a short message to request a warm connection.",
    prompt:
      "Help me write a concise LinkedIn message to request a referral for an engineering role.",
  },
];

export default function CoachPage() {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: "1",
      sender: "coach",
      text: "Hello! I am your Kirmya AI Career Coach, trained specifically in transition support, salary negotiation, and interview confidence. Choose a topic below or type anything to get started.",
      timestamp: new Date().toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      }),
    },
  ]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, loading]);

  const handleSendMessage = (textToSend: string) => {
    if (!textToSend.trim() || loading) return;

    const userMessage: Message = {
      id: Math.random().toString(),
      sender: "user",
      text: textToSend,
      timestamp: new Date().toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      }),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInput("");
    setLoading(true);

    // Mock AI response delay
    setTimeout(() => {
      let coachReply = "";
      const textLower = textToSend.toLowerCase();

      if (textLower.includes("layoff")) {
        coachReply =
          "Layoffs are a business metric, not a personal report card. When a recruiter asks, keep it brief and positive: 'My role was impacted by restructuring alongside X% of the workforce. What that gave me was an opportunity to focus fully on honing my skills in Y.' Let's try to draft your specific pitch. Tell me, what was your previous role?";
      } else if (
        textLower.includes("salary") ||
        textLower.includes("counter")
      ) {
        coachReply =
          "The key to salary negotiation is gratitude followed by a firm, clear range. Try: 'Thank you so much for the offer, I'm thrilled about the team. Given my experience in X and matching market indicators, I was hoping to land closer to $185,000. Is there flexibility to make that adjustment?' Never apologize for asking for your market worth.";
      } else if (
        textLower.includes("quiz") ||
        textLower.includes("interview") ||
        textLower.includes("s&op")
      ) {
        coachReply =
          "Excellent, let's start a mock interview. Imagine I am the hiring manager. Question 1: 'Can you describe a time you resolved a major constraint in your team's weekly workflow under pressure?' Answer as you would in an interview.";
      } else {
        coachReply =
          "That is a great starting point. To give you the most tailored strategy, tell me a bit more about the job you are targeting or upload your resume in the Resume section. What is the next milestone in your search?";
      }

      const coachMessage: Message = {
        id: Math.random().toString(),
        sender: "coach",
        text: coachReply,
        timestamp: new Date().toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
        }),
      };

      setMessages((prev) => [...prev, coachMessage]);
      setLoading(false);
    }, 1800);
  };

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav
        breadcrumb={[{ label: "Home", href: "/" }, { label: "Coach" }]}
      />

      <main className="flex-grow max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8 flex flex-col lg:flex-row gap-6 overflow-hidden">
        {/* Left Side: Coach Info / Prompts */}
        <div className="lg:w-[320px] lg:flex-none flex flex-col gap-6">
          <div className="space-y-2">
            <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-primary/10 text-primary border border-primary/20">
              <Sparkles className="h-3.5 w-3.5" />
              AI Career Operating System
            </div>
            <h1 className="text-2xl font-extrabold tracking-tight">
              AI Career Coach
            </h1>
            <p className="text-sm text-muted-foreground leading-relaxed">
              Your on-demand partner for interview training, cover letter
              generation, and salary negotiation. Always private, always active.
            </p>
          </div>

          {/* Quick starterm prompt cards */}
          <div className="space-y-3">
            <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest block">
              Quick Starters
            </span>
            <div className="grid grid-cols-1 gap-2.5">
              {STARTER_PROMPTS.map((prompt, idx) => {
                const Icon = prompt.icon;
                return (
                  <div
                    key={idx}
                    onClick={() => handleSendMessage(prompt.prompt)}
                    className="p-3.5 bg-card border border-border/60 rounded-2xl cursor-pointer hover:bg-secondary/40 hover:border-border transition-all flex flex-col gap-1.5 group text-left"
                  >
                    <div className="flex items-center gap-2">
                      <div className="h-6 w-6 rounded-lg bg-primary/10 border border-primary/20 flex items-center justify-center text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-all shrink-0">
                        <Icon className="h-3.5 w-3.5" />
                      </div>
                      <span className="text-xs font-bold group-hover:text-primary transition-colors">
                        {prompt.title}
                      </span>
                    </div>
                    <p className="text-[10.5px] text-muted-foreground leading-normal">
                      {prompt.desc}
                    </p>
                  </div>
                );
              })}
            </div>
          </div>
        </div>

        {/* Right Side: Interactive Chat Panel */}
        <div className="flex-grow flex flex-col bg-card border border-border/80 rounded-3xl overflow-hidden shadow-sm h-[calc(100vh-160px)]">
          {/* Header */}
          <div className="p-4 md:px-6 border-b border-border/40 flex items-center justify-between">
            <div className="flex items-center gap-2.5">
              <div className="h-9 w-9 rounded-full bg-primary/15 border border-primary/20 flex items-center justify-center text-primary shrink-0">
                <Sparkles className="h-4.5 w-4.5" />
              </div>
              <div>
                <span className="text-sm font-bold block">
                  Kirmya Career Co-pilot
                </span>
                <span className="text-[10px] text-muted-foreground block">
                  Always online &bull; Secure AES Encryption
                </span>
              </div>
            </div>
            <div className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
          </div>

          {/* Messages display area */}
          <div className="flex-1 overflow-y-auto p-6 space-y-4">
            <AnimatePresence initial={false}>
              {messages.map((msg) => (
                <motion.div
                  key={msg.id}
                  initial={{ opacity: 0, y: 12, scale: 0.98 }}
                  animate={{ opacity: 1, y: 0, scale: 1 }}
                  transition={{ duration: 0.2 }}
                  className={`flex flex-col max-w-[85%] sm:max-w-[70%] ${
                    msg.sender === "user"
                      ? "ml-auto items-end"
                      : "mr-auto items-start"
                  }`}
                >
                  <div
                    className={`p-4 rounded-3xl text-sm leading-relaxed ${
                      msg.sender === "user"
                        ? "bg-primary text-primary-foreground rounded-tr-none shadow-sm shadow-blue-500/5"
                        : "bg-secondary/40 text-foreground rounded-tl-none border border-border/40"
                    }`}
                  >
                    {msg.text}
                  </div>
                  <span className="text-[9px] text-muted-foreground mt-1 px-1">
                    {msg.timestamp}
                  </span>
                </motion.div>
              ))}

              {loading && (
                <motion.div
                  initial={{ opacity: 0, y: 8 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="flex items-center gap-2 mr-auto bg-secondary/40 border border-border/40 p-4 rounded-3xl rounded-tl-none text-sm text-muted-foreground max-w-[70%]"
                >
                  <Loader2 className="h-4 w-4 animate-spin text-primary shrink-0" />
                  <span>Coach is analyzing market responses...</span>
                </motion.div>
              )}
            </AnimatePresence>
            <div ref={messagesEndRef} />
          </div>

          {/* Input control row */}
          <div className="p-4 border-t border-border/40 bg-card">
            <div className="relative">
              <input
                type="text"
                placeholder="Ask Kirmya Coach about job status, negotiation scripts, or mock prep..."
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleSendMessage(input)}
                className="w-full pl-4 pr-12 py-3 rounded-full border border-border/80 bg-secondary/15 text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary text-sm shadow-sm"
              />
              <button
                onClick={() => handleSendMessage(input)}
                disabled={loading || !input.trim()}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-2 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground transition-all disabled:opacity-40 disabled:hover:bg-primary shadow-sm"
              >
                <Send className="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
