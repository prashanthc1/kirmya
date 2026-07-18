"use client";

import React, { useState } from "react";
import { usePendingRequests, useDeclineConnection } from "@/hooks/useConnections";
import ConnectButton from "./ConnectButton";
import { useNotifications } from "@/components/shared/Notifications";
import { Inbox, Send, MessageSquareText, Calendar, X } from "lucide-react";

export default function PendingRequestsPanel() {
  const { showNotification } = useNotifications();
  const [activeTab, setActiveTab] = useState<"incoming" | "outgoing">("incoming");

  // Fetch pending requests
  const { data: incoming, isLoading: loadingIncoming } = usePendingRequests("incoming");
  const { data: outgoing, isLoading: loadingOutgoing } = usePendingRequests("outgoing");

  const cancelMutation = useDeclineConnection(); // Cancelling outgoing is essentially deleting/declining it

  const handleCancelRequest = (connectionId: string, targetUserId: string) => {
    if (!window.confirm("Are you sure you want to cancel this connection request?")) return;
    cancelMutation.mutate(
      { connectionId, targetUserId },
      {
        onSuccess: () => showNotification("Connection request cancelled.", "success"),
        onError: (err: any) => showNotification(err.message || "Failed to cancel request", "error"),
      }
    );
  };

  const currentList = activeTab === "incoming" ? incoming : outgoing;
  const isLoading = activeTab === "incoming" ? loadingIncoming : loadingOutgoing;

  return (
    <div className="rounded-3xl bg-card border border-border/80 overflow-hidden shadow-sm">
      {/* Tabs */}
      <div className="flex border-b border-border/40 bg-secondary/10">
        <button
          onClick={() => setActiveTab("incoming")}
          className={`flex-1 py-4 text-xs font-bold flex items-center justify-center gap-2 cursor-pointer transition-all border-b-2 ${
            activeTab === "incoming"
              ? "border-primary text-primary bg-secondary/20"
              : "border-transparent text-muted-foreground hover:text-foreground hover:bg-secondary/10"
          }`}
        >
          <Inbox className="w-4 h-4" />
          Received ({incoming?.length || 0})
        </button>
        <button
          onClick={() => setActiveTab("outgoing")}
          className={`flex-1 py-4 text-xs font-bold flex items-center justify-center gap-2 cursor-pointer transition-all border-b-2 ${
            activeTab === "outgoing"
              ? "border-primary text-primary bg-secondary/20"
              : "border-transparent text-muted-foreground hover:text-foreground hover:bg-secondary/10"
          }`}
        >
          <Send className="w-4 h-4" />
          Sent ({outgoing?.length || 0})
        </button>
      </div>

      <div className="p-6">
        {isLoading ? (
          <div className="space-y-4">
            {Array.from({ length: 2 }).map((_, idx) => (
              <div key={idx} className="flex gap-4 p-4 rounded-xl bg-secondary/10 border border-border/60 animate-pulse">
                <div className="w-12 h-12 rounded-full bg-muted" />
                <div className="flex-1 space-y-2">
                  <div className="h-4 bg-muted rounded w-1/4" />
                  <div className="h-3 bg-muted rounded w-1/2" />
                </div>
              </div>
            ))}
          </div>
        ) : !currentList || currentList.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-10 text-center">
            <Inbox className="w-10 h-10 text-muted-foreground/60 mb-3" />
            <h4 className="text-sm font-bold text-foreground">No pending requests</h4>
            <p className="text-xs text-muted-foreground mt-1 max-w-[240px]">
              {activeTab === "incoming"
                ? "When people request to connect with you, they will appear here."
                : "Any connection requests you send will list here until they respond."}
            </p>
          </div>
        ) : (
          <div className="divide-y divide-border/40">
            {currentList.map((c) => {
              const u = c.user;
              return (
                <div key={c.id} className="py-4 first:pt-0 last:pb-0 flex flex-col gap-3">
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex gap-3 min-w-0">
                      <div className="w-12 h-12 rounded-full overflow-hidden bg-secondary border border-border flex-shrink-0">
                        {u.avatar_url ? (
                          <img src={u.avatar_url} alt={u.name} className="w-full h-full object-cover" />
                        ) : (
                          <div className="w-full h-full flex items-center justify-center bg-gradient-to-tr from-primary/10 to-primary/20 text-sm font-bold text-primary uppercase">
                            {u.name.slice(0, 1)}
                          </div>
                        )}
                      </div>
                      <div className="min-w-0">
                        <h4 className="text-sm font-bold text-foreground truncate leading-tight">{u.name}</h4>
                        <p className="text-xs text-muted-foreground truncate mt-0.5">{u.headline || "Practitioner"}</p>
                        <div className="flex items-center gap-1.5 text-[10px] text-muted-foreground/80 mt-1.5">
                          <Calendar className="w-3.5 h-3.5" />
                          Requested {new Date(c.created_at).toLocaleDateString()}
                        </div>
                      </div>
                    </div>

                    {/* Actions */}
                    {activeTab === "incoming" ? (
                      <ConnectButton
                        targetUserId={u.id}
                        targetUserName={u.name}
                        currentConnectionStatus="pending_incoming"
                        connectionId={c.id}
                      />
                    ) : (
                      <button
                        onClick={() => handleCancelRequest(c.id, u.id)}
                        className="px-3 py-1.5 text-xs rounded-xl font-medium border border-destructive/20 text-destructive hover:bg-destructive/10 cursor-pointer flex items-center gap-1 transition-colors"
                      >
                        <X className="w-3.5 h-3.5" />
                        Cancel
                      </button>
                    )}
                  </div>

                  {/* Attachment note */}
                  {c.note && (
                    <div className="ml-15 p-3 rounded-xl bg-primary/5 border border-primary/10 text-foreground text-xs flex gap-2 items-start leading-relaxed">
                      <MessageSquareText className="w-4 h-4 text-primary flex-shrink-0 mt-0.5" />
                      <div>
                        <span className="font-bold text-primary">Message note: </span>
                        {c.note}
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
