"use client";

import React from "react";
import Link from "next/link";
import { 
  GraduationCap, 
  Users, 
  Calendar, 
  Sparkles, 
  HeartHandshake, 
  ChevronRight, 
  Award,
  Video,
  MessagesSquare,
  BookmarkCheck
} from "lucide-react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

const STATS = [
  { metric: "2,100+", label: "Active mentors" },
  { metric: "Free", label: "For mentees, always" },
  { metric: "9,800", label: "Sessions this year" },
  { metric: "4.9★", label: "Average rating" }
];

const STEPS = [
  {
    step: "1",
    title: "Match on what matters",
    desc: "Tell us your industry field and what you're working through. We pair you with a mentor who's successfully navigated the same situation."
  },
  {
    step: "2",
    title: "Meet on your terms",
    desc: "Book a single strategic session or a recurring arrangement. Meet via video, phone call, or async text — whatever fits your recovery schedule."
  },
  {
    step: "3",
    title: "Move forward with a plan",
    desc: "Leave each session with concrete action items — a refined resume direction, a practiced layoff response pitch, and clear next steps."
  }
];

const MENTOR_TOPICS = [
  { icon: FileTextIcon, title: "Resume & Profile Review", desc: "Make your actual experience pop." },
  { icon: TargetIcon, title: "Interview Practice", desc: "Interactive mock rounds & feedback." },
  { icon: MoveIcon, title: "Career Path Pivots", desc: "Translate skills into a new niche." },
  { icon: HandshakeIcon, title: "Offer Negotiation", desc: "Know your market worth & request it." }
];

function FileTextIcon() { return <Award className="h-5 w-5 text-primary" />; }
function TargetIcon() { return <Sparkles className="h-5 w-5 text-primary" />; }
function MoveIcon() { return <HeartHandshake className="h-5 w-5 text-primary" />; }
function HandshakeIcon() { return <Users className="h-5 w-5 text-primary" />; }

export default function MentorshipPage() {
  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col relative overflow-hidden">
      {/* Background Glow */}
      <div className="absolute top-[-5%] right-[-10%] w-[500px] h-[500px] rounded-full bg-emerald-500/5 blur-[120px] pointer-events-none" />

      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Mentorship" }]} />

      {/* Hero Section */}
      <section className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-16 md:py-24 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
          
          {/* Left Column: Heading */}
          <div className="lg:col-span-7 space-y-6">
            <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20">
              <HeartHandshake className="h-3.5 w-3.5" />
              Community Mentorship
            </div>
            
            <h1 className="text-4xl sm:text-5xl md:text-6xl font-extrabold tracking-tight leading-[1.05]">
              No one should job-hunt alone.
            </h1>
            
            <p className="text-base sm:text-lg text-muted-foreground max-w-xl leading-relaxed">
              Get matched with an experienced industry professional who has navigated a career layoff and come out stronger — or pay it forward by guiding someone through theirs.
            </p>

            <div className="flex flex-wrap items-center gap-3 pt-2">
              <Link
                href="/mentors"
                className="px-6 py-3 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-sm font-bold shadow-lg shadow-blue-500/10 flex items-center gap-1.5 group"
              >
                Find a Mentor
                <ChevronRight className="h-4 w-4 group-hover:translate-x-1 transition-transform" />
              </Link>
              <Link
                href="/sign-in"
                className="px-6 py-3 rounded-full border border-border hover:bg-secondary text-sm font-bold"
              >
                Become a Mentor
              </Link>
            </div>
          </div>

          {/* Right Column: Featured Mentor Card */}
          <div className="lg:col-span-5 bg-card border border-border/80 p-6 rounded-3xl shadow-lg relative overflow-hidden">
            <div className="flex items-center gap-4 mb-4">
              <div className="h-12 w-12 rounded-full bg-emerald-500/15 border border-emerald-500/20 flex items-center justify-center text-emerald-600 dark:text-emerald-400 font-extrabold text-sm select-none shrink-0">
                RG
              </div>
              <div className="space-y-0.5">
                <h3 className="text-sm font-bold">Rosa Guerrero</h3>
                <p className="text-xs text-muted-foreground">Sr. Logistics Manager &bull; 17 years exp</p>
              </div>
              <span className="ml-auto px-2.5 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-[9px] font-extrabold text-emerald-500 uppercase tracking-wider">
                Mentor
              </span>
            </div>
            
            <p className="text-base font-semibold leading-relaxed text-foreground italic mb-4">
              &ldquo;I had three mentors when I was laid off in 2009. This is how I pass that forward.&rdquo;
            </p>

            <div className="flex flex-wrap gap-1.5">
              {["Career pivots", "Interview prep", "Offer negotiation"].map((tag) => (
                <span key={tag} className="text-[10px] px-2.5 py-0.5 bg-secondary text-muted-foreground rounded-full border border-border/40 font-semibold">
                  {tag}
                </span>
              ))}
            </div>
          </div>

        </div>
      </section>

      {/* Stats counter panel */}
      <section className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 pb-16">
        <div className="bg-card border border-border/60 p-8 rounded-3xl grid grid-cols-2 md:grid-cols-4 gap-6 text-center shadow-sm">
          {STATS.map((stat, idx) => (
            <div key={idx} className="space-y-1">
              <span className="text-3xl font-black tracking-tight text-primary block">{stat.metric}</span>
              <span className="text-xs text-muted-foreground font-semibold block">{stat.label}</span>
            </div>
          ))}
        </div>
      </section>

      {/* 3 Step Process */}
      <section className="py-20 bg-secondary/25 border-y border-border/40">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-2xl mx-auto mb-16 space-y-2">
            <span className="text-xs font-bold uppercase tracking-widest text-primary">How it works</span>
            <h2 className="text-3xl font-extrabold tracking-tight">Structured coaching in three steps</h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {STEPS.map((step) => (
              <div key={step.step} className="bg-card border border-border/50 p-8 rounded-3xl space-y-4 shadow-sm hover:border-border transition-all">
                <div className="h-10 w-10 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-extrabold text-sm">
                  {step.step}
                </div>
                <h3 className="text-lg font-bold text-foreground">{step.title}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{step.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* What mentors help with */}
      <section className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-20">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
          {/* Left Column: Summary */}
          <div className="lg:col-span-5 space-y-4">
            <span className="text-xs font-bold uppercase tracking-widest text-primary">Focus Areas</span>
            <h2 className="text-3xl font-extrabold tracking-tight">Practical guidance from people who have been there</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              Our mentors are not career coaches — they are active industry professionals who have stood in your shoes. Their advice is tactical, transparent, and built on real experience.
            </p>
          </div>

          {/* Right Column: Grid */}
          <div className="lg:col-span-7 grid grid-cols-1 sm:grid-cols-2 gap-4">
            {MENTOR_TOPICS.map((topic, idx) => {
              const Icon = topic.icon;
              return (
                <div key={idx} className="bg-card border border-border/60 p-6 rounded-2xl space-y-2 hover:border-border transition-all shadow-sm">
                  <Icon />
                  <h3 className="text-sm font-bold text-foreground pt-1">{topic.title}</h3>
                  <p className="text-xs text-muted-foreground leading-normal">{topic.desc}</p>
                </div>
              );
            })}
          </div>
        </div>
      </section>

      {/* Call to action card */}
      <section className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 pb-20">
        <div className="bg-gradient-to-br from-emerald-500 to-teal-600 dark:from-emerald-600 dark:to-teal-700 p-8 md:p-12 rounded-3xl text-center space-y-6 text-white shadow-xl shadow-emerald-500/5">
          <h2 className="text-3xl md:text-4xl font-extrabold tracking-tight">Find an internal advisor today</h2>
          <p className="text-sm md:text-base text-emerald-100 max-w-md mx-auto leading-relaxed">
            Browse our database of 2,100+ mentors who have walked this road. Your initial alignment session is completely free.
          </p>
          <Link
            href="/mentors"
            className="inline-flex px-8 py-3 rounded-full bg-white hover:bg-emerald-50 text-slate-900 text-sm font-bold shadow-sm transition-all"
          >
            Browse Mentors List
          </Link>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
