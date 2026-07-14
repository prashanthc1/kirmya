"use client";

import React, { useState } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { Mail, MessageSquare, Send, CheckCircle2 } from "lucide-react";

export default function ContactPage() {
  const [email, setEmail] = useState("");
  const [message, setMessage] = useState("");
  const [sent, setSent] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim() || !message.trim()) return;
    setSent(true);
  };

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Contact Us" }]} />

      <main className="flex-grow max-w-lg mx-auto px-4 sm:px-6 py-12 w-full space-y-6">
        <div className="flex items-center gap-3 border-b border-border/40 pb-4">
          <div className="h-10 w-10 bg-primary/10 border border-primary/20 rounded-xl flex items-center justify-center text-primary">
            <Mail className="h-5 w-5" />
          </div>
          <div>
            <h1 className="text-xl font-black tracking-tight">Contact Support</h1>
            <p className="text-[10px] text-muted-foreground">Send a message to the Kirmya core team.</p>
          </div>
        </div>

        {sent ? (
          <div className="bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 p-6 rounded-2xl text-center space-y-3">
            <CheckCircle2 className="h-8 w-8 text-emerald-500 mx-auto" />
            <div>
              <h3 className="text-xs font-bold">Message Sent!</h3>
              <p className="text-[10px] text-muted-foreground mt-1">
                Thank you for reaching out. We will get back to you at {email} within 24 hours.
              </p>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="bg-card border border-border/60 p-6 rounded-3xl shadow-sm space-y-4">
            <div className="space-y-1.5">
              <label htmlFor="email" className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground block">
                Your Email Address
              </label>
              <input
                id="email"
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@example.com"
                className="w-full px-3.5 py-2 rounded-xl border border-border/60 bg-secondary/15 placeholder:text-muted-foreground text-xs focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
              />
            </div>

            <div className="space-y-1.5">
              <label htmlFor="message" className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground block">
                Message Content
              </label>
              <textarea
                id="message"
                required
                rows={5}
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder="Type your question or feedback..."
                className="w-full px-3.5 py-2 rounded-xl border border-border/60 bg-secondary/15 placeholder:text-muted-foreground text-xs focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
              />
            </div>

            <button
              type="submit"
              className="w-full py-2 bg-primary hover:bg-primary/95 text-primary-foreground text-xs font-bold rounded-full transition-all shadow-sm flex items-center justify-center gap-1.5 cursor-pointer"
            >
              <Send className="h-3.5 w-3.5" />
              Send Message
            </button>
          </form>
        )}
      </main>

      <SiteFooter />
    </div>
  );
}
