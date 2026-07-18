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
        <div className="h-10 w-full bg-muted/80 rounded-xl animate-pulse" />
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.from({ length: 6 }).map((_, idx) => (
            <div key={idx} className="h-40 bg-card border border-border/60 rounded-3xl animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Search Bar */}
      <div className="relative">
        <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search connections by name or headline..."
          className="w-full pl-11 pr-4 py-3 text-sm rounded-xl bg-secondary/15 border border-border text-foreground placeholder:text-muted-foreground focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-all shadow-inner"
        />
      </div>

      {!connections || connections.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-center rounded-3xl bg-card border border-border/60 shadow-sm">
          <Users className="w-12 h-12 text-muted-foreground/60 mb-3" />
          <h4 className="text-sm font-bold text-foreground">No connections yet</h4>
          <p className="text-xs text-muted-foreground mt-1 max-w-[280px]">
            Grow your network by accepting pending requests or checking out recommendations on your dashboard.
          </p>
        </div>
      ) : filteredConnections.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-center rounded-3xl bg-card border border-border/60 shadow-sm">
          <Search className="w-10 h-10 text-muted-foreground/60 mb-3" />
          <h4 className="text-sm font-bold text-foreground">No matching connections</h4>
          <p className="text-xs text-muted-foreground mt-1">
            We couldn&apos;t find anyone matching &quot;{searchQuery}&quot;.
          </p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredConnections.map((c) => {
              const u = c.user;
              return (
                <div
                  key={c.id}
                  className="p-5 rounded-3xl bg-card border border-border/80 hover:border-primary/30 transition-all shadow-sm hover:shadow flex flex-col justify-between"
                >
                  <div className="flex items-start gap-3 min-w-0">
                    <div className="w-12 h-12 rounded-full overflow-hidden bg-secondary border border-border flex-shrink-0">
                      {u.avatar_url ? (
                        <img src={u.avatar_url} alt={u.name} className="w-full h-full object-cover" />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center bg-gradient-to-tr from-primary/10 to-primary/20 text-sm font-bold text-primary uppercase">
                          {u.name.slice(0, 1)}
                        </div>
                      )}
                    </div>
                    <div className="min-w-0 flex-1">
                      <h4 className="text-sm font-bold text-foreground truncate leading-tight">{u.name}</h4>
                      <p className="text-xs text-muted-foreground line-clamp-2 mt-1 leading-relaxed min-h-[36px]">
                        {u.headline || "Practitioner"}
                      </p>
                    </div>
                  </div>

                  <div className="mt-4 pt-4 border-t border-border/40 flex justify-end">
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
                className="px-4 py-2 text-xs font-semibold rounded-xl border border-border hover:bg-secondary/40 transition-all text-foreground disabled:opacity-30 cursor-pointer flex items-center gap-1.5"
              >
                <ArrowLeft className="w-4 h-4" />
                Previous
              </button>
              <span className="text-xs font-medium text-muted-foreground">Page {page}</span>
              <button
                onClick={() => setPage((p) => p + 1)}
                disabled={!hasNextPage}
                className="px-4 py-2 text-xs font-semibold rounded-xl border border-border hover:bg-secondary/40 transition-all text-foreground disabled:opacity-30 cursor-pointer flex items-center gap-1.5"
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
