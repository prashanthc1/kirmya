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
        <div className="p-3 bg-red-950/20 border border-red-500/25 rounded-xl text-red-200 text-xs leading-relaxed space-y-1.5">
          <p className="font-bold">⚠️ Warning: Blocking has immediate consequences:</p>
          <ul className="list-disc list-inside space-y-0.5">
            <li>Any existing active connection will be permanently removed.</li>
            <li>You will no longer be able to message each other.</li>
            <li>Profiles and activities will be hidden from each other.</li>
          </ul>
        </div>

        <div className="space-y-1">
          <label className="text-xs font-semibold text-gray-400">
            Reason for blocking (Optional)
          </label>
          <textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="Help us understand the issue (e.g. spam, inappropriate behavior)..."
            rows={3}
            disabled={isSubmitting}
            className="w-full p-3 text-sm rounded-xl bg-gray-900/60 border border-white/10 text-white placeholder-gray-500 focus:outline-none focus:border-red-500 focus:ring-1 focus:ring-red-500 resize-none transition-colors"
          />
        </div>

        <div className="flex justify-end gap-3 pt-2">
          <button
            type="button"
            onClick={onClose}
            disabled={isSubmitting}
            className="px-4 py-2 text-sm rounded-xl font-medium border border-white/10 hover:bg-white/5 transition-colors cursor-pointer text-gray-300"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isSubmitting}
            className="px-5 py-2 text-sm rounded-xl font-medium bg-red-600 hover:bg-red-700 transition-colors text-white disabled:opacity-50 cursor-pointer shadow-lg shadow-red-500/10"
          >
            {isSubmitting ? "Blocking..." : "Yes, Block User"}
          </button>
        </div>
      </form>
    </MuiModal>
  );
}
