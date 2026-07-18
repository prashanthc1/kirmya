"use client";

import React, { useState } from "react";
import MuiModal from "@/components/shared/MuiModal";

interface BlockConfirmDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: (reason: string) => void;
  isSubmitting: boolean;
  userName: string;
}

export default function BlockConfirmDialog({
  open,
  onClose,
  onConfirm,
  isSubmitting,
  userName,
}: BlockConfirmDialogProps) {
  const [reason, setReason] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onConfirm(reason);
  };

  return (
    <MuiModal
      open={open}
      onClose={onClose}
      title={`Block ${userName}?`}
      maxWidth="xs"
    >
      <form onSubmit={handleSubmit} className="space-y-4 pt-2">
        <div className="p-3 bg-destructive/10 border border-destructive/20 rounded-xl text-destructive text-xs leading-relaxed space-y-1.5">
          <p className="font-bold">⚠️ Warning: Blocking has immediate consequences:</p>
          <ul className="list-disc list-inside space-y-0.5">
            <li>Any existing active connection will be permanently removed.</li>
            <li>You will no longer be able to message each other.</li>
            <li>Profiles and activities will be hidden from each other.</li>
          </ul>
        </div>

        <div className="space-y-1">
          <label className="text-xs font-semibold text-muted-foreground">
            Reason for blocking (Optional)
          </label>
          <textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="Help us understand the issue (e.g. spam, inappropriate behavior)..."
            rows={3}
            disabled={isSubmitting}
            className="w-full p-3 text-sm rounded-xl bg-secondary/15 border border-border text-foreground placeholder:text-muted-foreground focus:outline-none focus:border-destructive focus:ring-1 focus:ring-destructive resize-none transition-colors"
          />
        </div>

        <div className="flex justify-end gap-3 pt-2">
          <button
            type="button"
            onClick={onClose}
            disabled={isSubmitting}
            className="px-4 py-2 text-sm rounded-xl font-medium border border-border hover:bg-secondary/40 transition-colors cursor-pointer text-foreground"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isSubmitting}
            className="px-5 py-2 text-sm rounded-xl font-medium bg-destructive hover:bg-destructive/95 transition-colors text-destructive-foreground disabled:opacity-50 cursor-pointer shadow-sm"
          >
            {isSubmitting ? "Blocking..." : "Yes, Block User"}
          </button>
        </div>
      </form>
    </MuiModal>
  );
}
