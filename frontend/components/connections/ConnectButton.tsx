"use client";

import React, { useState } from "react";
import { useNotifications } from "@/components/shared/Notifications";
import {
  useSendConnectionRequest,
  useAcceptConnection,
  useDeclineConnection,
  useRemoveConnection,
  useBlockUser,
  useConnectionsStore,
} from "@/hooks/useConnections";
import ConnectionNoteModal from "./ConnectionNoteModal";
import BlockConfirmDialog from "./BlockConfirmDialog";
import { api } from "@/lib/api/client";
import { UserCheck, Clock, UserX, ShieldAlert, ChevronDown, MessageSquare } from "lucide-react";

interface ConnectButtonProps {
  targetUserId: string;
  targetUserName?: string;
  currentConnectionStatus: "none" | "pending_outgoing" | "pending_incoming" | "accepted" | "blocked" | "";
  connectionId?: string; // If known, useful for accept/decline/remove mutations
}

export default function ConnectButton({
  targetUserId,
  targetUserName = "User",
  currentConnectionStatus,
  connectionId: initialConnectionId,
}: ConnectButtonProps) {
  const { showNotification } = useNotifications();
  const statusOverrides = useConnectionsStore((s) => s.statusOverrides);
  const [showNoteModal, setShowNoteModal] = useState(false);
  const [showBlockModal, setShowBlockModal] = useState(false);
  const [showDropdown, setShowDropdown] = useState(false);

  // Determine current effective status (favor override if present)
  const effectiveStatus = (statusOverrides[targetUserId] !== undefined
    ? statusOverrides[targetUserId]
    : currentConnectionStatus || "none") as "none" | "pending_outgoing" | "pending_incoming" | "accepted" | "blocked";

  // Mutations
  const sendMutation = useSendConnectionRequest();
  const acceptMutation = useAcceptConnection();
  const declineMutation = useDeclineConnection();
  const removeMutation = useRemoveConnection();
  const blockMutation = useBlockUser();

  const isMutating =
    sendMutation.isPending ||
    acceptMutation.isPending ||
    declineMutation.isPending ||
    removeMutation.isPending ||
    blockMutation.isPending;

  if (effectiveStatus === "blocked") {
    return null; // Hidden entirely
  }

  const handleSend = (note: string) => {
    setShowNoteModal(false);
    sendMutation.mutate(
      { targetUserId, note },
      {
        onSuccess: () => showNotification("Connection request sent!", "success"),
        onError: (err: any) => showNotification(err.message || "Failed to send request", "error"),
      }
    );
  };

  const handleAccept = () => {
    if (!initialConnectionId) {
      showNotification("Cannot process request: Connection ID missing", "error");
      return;
    }
    acceptMutation.mutate(
      { connectionId: initialConnectionId, targetUserId },
      {
        onSuccess: () => showNotification("Connection request accepted!", "success"),
        onError: (err: any) => showNotification(err.message || "Failed to accept connection", "error"),
      }
    );
  };

  const handleDecline = () => {
    if (!initialConnectionId) {
      showNotification("Cannot process request: Connection ID missing", "error");
      return;
    }
    declineMutation.mutate(
      { connectionId: initialConnectionId, targetUserId },
      {
        onSuccess: () => showNotification("Connection request declined.", "success"),
        onError: (err: any) => showNotification(err.message || "Failed to decline connection", "error"),
      }
    );
  };

  const handleRemove = () => {
    if (!initialConnectionId) {
      showNotification("Cannot process request: Connection ID missing", "error");
      return;
    }
    if (!window.confirm(`Are you sure you want to remove your connection with ${targetUserName}?`)) return;
    removeMutation.mutate(
      { connectionId: initialConnectionId, targetUserId },
      {
        onSuccess: () => showNotification("Connection removed.", "success"),
        onError: (err: any) => showNotification(err.message || "Failed to remove connection", "error"),
      }
    );
  };

  const handleBlockConfirm = (reason: string) => {
    setShowBlockModal(false);
    blockMutation.mutate(
      { targetUserId, reason },
      {
        onSuccess: () => showNotification(`${targetUserName} has been blocked.`, "success"),
        onError: (err: any) => showNotification(err.message || "Failed to block user", "error"),
      }
    );
  };

  const handleStartChat = async () => {
    try {
      const conv = await api.post<{ id: string }>("/conversations", {
        participant_ids: [targetUserId],
        title: targetUserName,
      });
      window.location.href = `/inbox?convId=${conv.id}`;
    } catch (err: any) {
      showNotification("Failed to start conversation", "error");
    }
  };

  switch (effectiveStatus) {
    case "none":
      return (
        <>
          <button
            onClick={() => setShowNoteModal(true)}
            disabled={isMutating}
            className="px-4 py-2 text-sm rounded-xl font-semibold bg-primary hover:bg-primary/95 active:scale-[0.98] text-primary-foreground shadow-sm cursor-pointer transition-all disabled:opacity-50"
          >
            Connect
          </button>

          <ConnectionNoteModal
            open={showNoteModal}
            onClose={() => handleSend("")} // skips & sends without note
            onSend={handleSend}
            isSubmitting={isMutating}
          />
        </>
      );

    case "pending_outgoing":
      return (
        <button
          disabled
          className="px-4 py-2 text-sm rounded-xl font-medium bg-secondary/40 border border-border text-muted-foreground flex items-center gap-1.5 cursor-not-allowed"
        >
          <Clock className="w-4 h-4" />
          Pending
        </button>
      );

    case "pending_incoming":
      return (
        <div className="flex gap-2">
          <button
            onClick={handleAccept}
            disabled={isMutating}
            className="px-3.5 py-1.5 text-xs rounded-xl font-semibold bg-primary hover:bg-primary/95 text-primary-foreground cursor-pointer active:scale-95 transition-transform"
          >
            Accept
          </button>
          <button
            onClick={handleDecline}
            disabled={isMutating}
            className="px-3.5 py-1.5 text-xs rounded-xl font-medium border border-border text-foreground hover:bg-secondary/40 cursor-pointer active:scale-95 transition-transform"
          >
            Decline
          </button>
        </div>
      );

    case "accepted":
      return (
        <div className="relative inline-block text-left">
          <div className="flex gap-1">
            <button
              onClick={handleStartChat}
              className="px-3 py-1.5 text-xs rounded-l-xl font-medium bg-secondary/20 hover:bg-secondary/40 border border-r-0 border-border text-primary flex items-center gap-1.5 cursor-pointer transition-colors"
            >
              <MessageSquare className="w-3.5 h-3.5" />
              Message
            </button>
            <button
              onClick={() => setShowDropdown(!showDropdown)}
              disabled={isMutating}
              className="px-2.5 py-1.5 rounded-r-xl bg-secondary/20 hover:bg-secondary/40 border border-border text-muted-foreground cursor-pointer transition-colors"
            >
              <ChevronDown className="w-3.5 h-3.5" />
            </button>
          </div>

          {showDropdown && (
            <>
              <div
                className="fixed inset-0 z-10"
                onClick={() => setShowDropdown(false)}
              />
              <div className="absolute right-0 mt-2 w-44 rounded-xl bg-card border border-border shadow-md z-20 overflow-hidden">
                <button
                  onClick={() => {
                    setShowDropdown(false);
                    handleRemove();
                  }}
                  className="w-full text-left px-4 py-2.5 text-xs font-medium text-destructive hover:bg-secondary/40 flex items-center gap-2 cursor-pointer transition-colors"
                >
                  <UserX className="w-4 h-4" />
                  Remove Connection
                </button>
                <button
                  onClick={() => {
                    setShowDropdown(false);
                    setShowBlockModal(true);
                  }}
                  className="w-full text-left px-4 py-2.5 text-xs font-medium text-destructive hover:bg-destructive/10 border-t border-border/40 flex items-center gap-2 cursor-pointer transition-colors"
                >
                  <ShieldAlert className="w-4 h-4" />
                  Block User
                </button>
              </div>
            </>
          )}

          <BlockConfirmDialog
            open={showBlockModal}
            onClose={() => setShowBlockModal(false)}
            onConfirm={handleBlockConfirm}
            isSubmitting={isMutating}
            userName={targetUserName}
          />
        </div>
      );

    default:
      return null;
  }
}
