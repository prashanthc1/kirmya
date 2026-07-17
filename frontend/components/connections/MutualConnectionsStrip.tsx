"use client";

import React, { useState } from "react";
import { useMutualConnections } from "@/hooks/useConnections";
import MuiModal from "@/components/shared/MuiModal";
import ConnectButton from "./ConnectButton";
import { Users } from "lucide-react";

interface MutualConnectionsStripProps {
  userId: string;
}

export default function MutualConnectionsStrip({ userId }: MutualConnectionsStripProps) {
  const [showModal, setShowModal] = useState(false);
  const { data, isLoading } = useMutualConnections(userId);

  if (isLoading || !data || data.total === 0) {
    return null; // Don't render if loading, error, or no mutual connections
  }

  const { users, total } = data;
  const avatarsToShow = users.slice(0, 3);

  return (
    <>
      <div
        onClick={() => setShowModal(true)}
        className="flex items-center gap-2 cursor-pointer hover:opacity-90 active:scale-[0.99] transition-all py-1.5 px-3.5 rounded-full bg-white/5 border border-white/5 backdrop-blur-sm w-fit"
        role="button"
        aria-label={`${total} mutual connections`}
      >
        <div className="flex -space-x-2">
          {avatarsToShow.map((u, i) => (
            <div
              key={u.id}
              className="w-6 h-6 rounded-full border border-gray-900 bg-gray-800 overflow-hidden flex-shrink-0"
              style={{ zIndex: 3 - i }}
            >
              {u.avatar_url ? (
                <img
                  src={u.avatar_url}
                  alt={u.name}
                  className="w-full h-full object-cover"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center bg-gradient-to-tr from-gray-700 to-gray-800 text-[10px] font-bold text-white uppercase">
                  {u.name.slice(0, 1)}
                </div>
              )}
            </div>
          ))}
        </div>

        <span className="text-xs font-semibold text-gray-300 hover:text-orange-400 transition-colors">
          {total} mutual connection{total > 1 ? "s" : ""}
        </span>
      </div>

      <MuiModal
        open={showModal}
        onClose={() => setShowModal(false)}
        title="Mutual Connections"
        maxWidth="xs"
      >
        <div className="space-y-4 pt-2 max-h-[380px] overflow-y-auto pr-1">
          {users.map((u) => (
            <div
              key={u.id}
              className="flex items-center justify-between p-3 rounded-xl bg-white/5 border border-white/5 hover:border-white/10 transition-all"
            >
              <div className="flex items-center gap-3 min-w-0 mr-2">
                <div className="w-10 h-10 rounded-full overflow-hidden bg-gray-800 flex-shrink-0 border border-white/5">
                  {u.avatar_url ? (
                    <img
                      src={u.avatar_url}
                      alt={u.name}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center bg-gradient-to-tr from-gray-700 to-gray-800 text-sm font-bold text-white uppercase">
                      {u.name.slice(0, 1)}
                    </div>
                  )}
                </div>
                <div className="min-w-0">
                  <h4 className="text-sm font-bold text-white truncate leading-tight">
                    {u.name}
                  </h4>
                  <p className="text-xs text-gray-400 truncate mt-0.5">
                    {u.headline || "Practitioner"}
                  </p>
                </div>
              </div>

              {/* Muted/Mouthpiece ConnectButton for the mutual user */}
              <ConnectButton
                targetUserId={u.id}
                targetUserName={u.name}
                currentConnectionStatus="accepted" // They are already connected to $1 (the viewer)
              />
            </div>
          ))}
        </div>
      </MuiModal>
    </>
  );
}
