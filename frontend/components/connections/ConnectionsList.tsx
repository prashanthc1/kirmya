"use client";

import React, { useState, useMemo } from "react";
import { useConnections } from "@/hooks/useConnections";
import ConnectButton from "./ConnectButton";
import { Users, Search, ArrowLeft, ArrowRight } from "lucide-react";

export default function ConnectionsList() {
  const [page, setPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");
  const limit = 9; // Grid of 3x3

  const { data: connections, isLoading } = useConnections(page, limit);

  // Local filtering for responsive search experience
  const filteredConnections = useMemo(() => {
    if (!connections) return [];
    return connections.filter((c) => {
      const nameMatch = c.user.name.toLowerCase().includes(searchQuery.toLowerCase());
      const headlineMatch = c.user.headline.toLowerCase().includes(searchQuery.toLowerCase());
      return nameMatch || headlineMatch;
    });
  }, [connections, searchQuery]);

  const hasNextPage = connections && connections.length === limit;
  const hasPrevPage = page > 1;

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="h-10 w-full bg-white/5 rounded-xl animate-pulse" />
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.from({ length: 6 }).map((_, idx) => (
            <div key={idx} className="h-40 bg-white/5 rounded-2xl animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Search Bar */}
      <div className="relative">
        <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search connections by name or headline..."
          className="w-full pl-11 pr-4 py-3 text-sm rounded-xl bg-[#0D1B2A]/60 border border-white/5 text-white placeholder-gray-500 focus:outline-none focus:border-orange-500 focus:ring-1 focus:ring-orange-500 transition-all shadow-inner"
        />
      </div>

      {!connections || connections.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-center rounded-2xl bg-[#0D1B2A]/40 border border-white/5 backdrop-blur-md">
          <Users className="w-12 h-12 text-gray-600 mb-3" />
          <h4 className="text-sm font-bold text-gray-400">No connections yet</h4>
          <p className="text-xs text-gray-500 mt-1 max-w-[280px]">
            Grow your network by accepting pending requests or checking out recommendations on your dashboard.
          </p>
        </div>
      ) : filteredConnections.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-center rounded-2xl bg-[#0D1B2A]/40 border border-white/5 backdrop-blur-md">
          <Search className="w-10 h-10 text-gray-600 mb-3" />
          <h4 className="text-sm font-bold text-gray-400">No matching connections</h4>
          <p className="text-xs text-gray-500 mt-1">
            We couldn't find anyone matching "{searchQuery}".
          </p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {filteredConnections.map((c) => {
              const u = c.user;
              return (
                <div
                  key={c.id}
                  className="p-5 rounded-2xl bg-[#0D1B2A]/60 border border-white/5 backdrop-blur-md hover:border-white/10 hover:bg-[#0D1B2A]/80 transition-all shadow-xl flex flex-col justify-between"
                >
                  <div className="flex items-start gap-3 min-w-0">
                    <div className="w-12 h-12 rounded-full overflow-hidden bg-gray-800 border border-white/10 flex-shrink-0">
                      {u.avatar_url ? (
                        <img src={u.avatar_url} alt={u.name} className="w-full h-full object-cover" />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center bg-gradient-to-tr from-gray-700 to-gray-800 text-base font-bold text-white uppercase">
                          {u.name.slice(0, 1)}
                        </div>
                      )}
                    </div>
                    <div className="min-w-0 flex-1">
                      <h4 className="text-sm font-bold text-white truncate leading-tight">{u.name}</h4>
                      <p className="text-xs text-gray-400 line-clamp-2 mt-1 leading-relaxed min-h-[36px]">
                        {u.headline || "Practitioner"}
                      </p>
                    </div>
                  </div>

                  <div className="mt-4 pt-4 border-t border-white/5 flex justify-end">
                    <ConnectButton
                      targetUserId={u.id}
                      targetUserName={u.name}
                      currentConnectionStatus="accepted"
                      connectionId={c.id}
                    />
                  </div>
                </div>
              );
            })}
          </div>

          {/* Pagination Controls */}
          {(hasPrevPage || hasNextPage) && (
            <div className="flex items-center justify-between pt-4">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={!hasPrevPage}
                className="px-4 py-2 text-xs font-semibold rounded-xl border border-white/5 hover:bg-white/5 transition-all text-gray-300 disabled:opacity-30 cursor-pointer flex items-center gap-1.5"
              >
                <ArrowLeft className="w-4 h-4" />
                Previous
              </button>
              <span className="text-xs font-medium text-gray-500">Page {page}</span>
              <button
                onClick={() => setPage((p) => p + 1)}
                disabled={!hasNextPage}
                className="px-4 py-2 text-xs font-semibold rounded-xl border border-white/5 hover:bg-white/5 transition-all text-gray-300 disabled:opacity-30 cursor-pointer flex items-center gap-1.5"
              >
                Next
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
