"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import AuthGuard from "@/components/shared/AuthGuard";
import { api } from "@/lib/api/client";
import { Bell, CheckCircle2, ChevronRight, Loader2, MailOpen } from "lucide-react";

interface Notification {
  id: string;
  type: string;
  title: string;
  body?: string;
  link?: string;
  read: boolean;
  created_at: string;
}

function NotificationsContent() {
  const router = useRouter();
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState(false);

  const fetchNotifications = async () => {
    try {
      const data = await api.get<{ notifications: Notification[]; unread: number }>("/notifications");
      setNotifications(data?.notifications ?? []);
      setUnreadCount(data?.unread ?? 0);
    } catch (err: any) {
      setError(err.message || "Failed to load notifications.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNotifications();
  }, []);

  const handleMarkAllRead = async () => {
    if (actionLoading) return;
    setActionLoading(true);
    try {
      await api.post("/notifications/read-all", {});
      setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
      setUnreadCount(0);
    } catch (err) {
      console.error("Failed to mark all as read:", err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleNotificationClick = async (notif: Notification) => {
    if (!notif.read) {
      try {
        await api.post(`/notifications/${notif.id}/read`, {});
        setNotifications((prev) =>
          prev.map((n) => (n.id === notif.id ? { ...n, read: true } : n))
        );
        setUnreadCount((prev) => Math.max(0, prev - 1));
      } catch (err) {
        console.error("Failed to mark notification as read:", err);
      }
    }
    if (notif.link) {
      router.push(notif.link);
    }
  };

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/dashboard" }, { label: "Notifications" }]} />

      <main className="flex-grow max-w-3xl mx-auto px-4 sm:px-6 py-8 w-full">
        <div className="space-y-6">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl font-black tracking-tight flex items-center gap-2">
                <Bell className="h-6 w-6 text-primary" />
                Notifications
              </h1>
              <p className="text-xs text-muted-foreground mt-1">
                You have {unreadCount} unread notification{unreadCount !== 1 ? "s" : ""}.
              </p>
            </div>

            {unreadCount > 0 && (
              <button
                onClick={handleMarkAllRead}
                disabled={actionLoading}
                className="self-start sm:self-auto px-4 py-1.5 rounded-full border border-border hover:bg-secondary text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer"
              >
                <MailOpen className="h-3.5 w-3.5 text-muted-foreground" />
                Mark all as read
              </button>
            )}
          </div>

          {loading ? (
            <div className="flex flex-col items-center justify-center py-20 gap-3">
              <Loader2 className="h-8 w-8 text-primary animate-spin" />
              <span className="text-xs font-semibold text-muted-foreground">Loading notifications...</span>
            </div>
          ) : error ? (
            <div className="p-6 bg-destructive/10 border border-destructive/20 text-destructive rounded-2xl text-xs font-semibold">
              {error}
            </div>
          ) : notifications.length === 0 ? (
            <div className="text-center py-16 border border-dashed border-border/60 rounded-3xl p-8 bg-secondary/15 space-y-3">
              <div className="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center mx-auto text-primary">
                <Bell className="h-6 w-6" />
              </div>
              <h3 className="text-sm font-bold">All caught up!</h3>
              <p className="text-xs text-muted-foreground max-w-sm mx-auto">
                No new notifications at this time. We will let you know when something important happens.
              </p>
            </div>
          ) : (
            <div className="bg-card border border-border/60 rounded-2xl divide-y divide-border/40 overflow-hidden shadow-sm">
              {notifications.map((notif) => (
                <div
                  key={notif.id}
                  onClick={() => handleNotificationClick(notif)}
                  className={`p-4 transition-all flex gap-3 cursor-pointer items-start hover:bg-secondary/20 ${
                    !notif.read ? "bg-secondary/10" : ""
                  }`}
                >
                  <div className="shrink-0 mt-0.5">
                    {!notif.read ? (
                      <div className="h-2 w-2 rounded-full bg-primary animate-pulse" />
                    ) : (
                      <div className="h-2 w-2 rounded-full bg-transparent" />
                    )}
                  </div>
                  <div className="flex-grow space-y-0.5">
                    <p className={`text-xs ${!notif.read ? "font-bold text-foreground" : "text-foreground/80"}`}>
                      {notif.title}
                    </p>
                    {notif.body && <p className="text-[11px] text-muted-foreground leading-relaxed">{notif.body}</p>}
                    <p className="text-[9px] text-muted-foreground/60">
                      {new Date(notif.created_at).toLocaleDateString([], { month: "short", day: "numeric" })} at{" "}
                      {new Date(notif.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                    </p>
                  </div>
                  {notif.link && (
                    <ChevronRight className="h-4 w-4 text-muted-foreground/45 shrink-0 self-center" />
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}

export default function NotificationsPage() {
  return (
    <AuthGuard>
      <NotificationsContent />
    </AuthGuard>
  );
}
