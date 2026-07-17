"use client";

import React, { useState } from "react";
import MuiModal from "@/components/shared/MuiModal";

interface ConnectionNoteModalProps {
  open: boolean;
  onClose: () => void;
  onSend: (note: string) => void;
  isSubmitting: boolean;
}

export default function ConnectionNoteModal({
  open,
  onClose,
  onSend,
  isSubmitting,
}: ConnectionNoteModalProps) {
  const [note, setNote] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (note.length <= 300) {
      onSend(note);
    }
  };

  return (
    <MuiModal
      open={open}
      onClose={onClose}
      title="Add a Personal Note (Optional)"
      maxWidth="xs"
    >
      <form onSubmit={handleSubmit} className="space-y-4 pt-2">
        <p className="text-sm text-gray-400">
          Adding a note about why you want to connect makes it more likely they accept.
        </p>

        <div className="relative">
          <textarea
            value={note}
            onChange={(e) => setNote(e.target.value.slice(0, 300))}
            placeholder="Type your message here..."
            rows={4}
            disabled={isSubmitting}
            className="w-full p-3 text-sm rounded-xl bg-gray-900/60 border border-white/10 text-white placeholder-gray-500 focus:outline-none focus:border-orange-500 focus:ring-1 focus:ring-orange-500 resize-none transition-colors"
          />
          <div
            className={`absolute bottom-2 right-3 text-xs ${
              note.length >= 280 ? "text-orange-500 font-bold" : "text-gray-500"
            }`}
          >
            {note.length} / 300
          </div>
        </div>

        <div className="flex justify-end gap-3 pt-2">
          <button
            type="button"
            onClick={onClose}
            disabled={isSubmitting}
            className="px-4 py-2 text-sm rounded-xl font-medium border border-white/10 hover:bg-white/5 transition-colors cursor-pointer text-gray-300"
          >
            Skip & Send
          </button>
          <button
            type="submit"
            disabled={isSubmitting || note.length > 300}
            className="px-5 py-2 text-sm rounded-xl font-medium bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 transition-colors text-white disabled:opacity-50 cursor-pointer shadow-lg shadow-orange-500/10"
          >
            {isSubmitting ? "Sending..." : "Send Note"}
          </button>
        </div>
      </form>
    </MuiModal>
  );
}
