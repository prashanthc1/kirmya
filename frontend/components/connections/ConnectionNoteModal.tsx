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
        <p className="text-sm text-muted-foreground">
          Adding a note about why you want to connect makes it more likely they accept.
        </p>

        <div className="relative">
          <textarea
            value={note}
            onChange={(e) => setNote(e.target.value.slice(0, 300))}
            placeholder="Type your message here..."
            rows={4}
            disabled={isSubmitting}
            className="w-full p-3 text-sm rounded-xl bg-secondary/15 border border-border text-foreground placeholder:text-muted-foreground focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary resize-none transition-colors"
          />
          <div
            className={`absolute bottom-2 right-3 text-xs ${
              note.length >= 280 ? "text-destructive font-bold" : "text-muted-foreground/80"
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
            className="px-4 py-2 text-sm rounded-xl font-medium border border-border hover:bg-secondary/40 transition-colors cursor-pointer text-foreground"
          >
            Skip & Send
          </button>
          <button
            type="submit"
            disabled={isSubmitting || note.length > 300}
            className="px-5 py-2 text-sm rounded-xl font-medium bg-primary hover:bg-primary/95 transition-colors text-primary-foreground disabled:opacity-50 cursor-pointer shadow-sm"
          >
            {isSubmitting ? "Sending..." : "Send Note"}
          </button>
        </div>
      </form>
    </MuiModal>
  );
}
