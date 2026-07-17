"use client";

import React from "react";
import { useSuggestions } from "@/hooks/useConnections";
import ConnectButton from "./ConnectButton";
import { Sparkles, Users } from "lucide-react";

interface SuggestionsCarouselProps {
  limit?: number;
}

export default function SuggestionsCarousel({ limit = 10 }: SuggestionsCarouselProps) {
  const { data: suggestions, isLoading } = useSuggestions(limit);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <div className="h-6 w-48 bg-white/5 rounded animate-pulse" />
        </div>
        <div className="flex gap-4 overflow-x-auto pb-4 scrollbar-thin scrollbar-thumb-white/5">
          {Array.from({ length: 4 }).map((_, idx) => (
            <div
              key={idx}
              className="w-[240px] flex-shrink-0 p-5 rounded-2xl bg-white/5 border border-white/5 space-y-4"
            >
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-full bg-white/5 animate-pulse" />
                <div className="space-y-2 flex-1">
                  <div className="h-4 bg-white/5 rounded animate-pulse w-3/4" />
                  <div className="h-3 bg-white/5 rounded animate-pulse w-1/2" />
                </div>
              </div>
              <div className="h-8 bg-white/5 rounded-xl animate-pulse" />
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (!suggestions || suggestions.length === 0) {
    return null; // Don't show if there are no suggestions
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2">
        <Sparkles className="w-5 h-5 text-orange-400" />
        <h3 className="text-base font-bold font-sora text-white">Recommended Connections</h3>
      </div>

      <div className="flex gap-4 overflow-x-auto pb-4 scrollbar-thin scrollbar-thumb-white/10 scroll-smooth">
        {suggestions.map((s) => (
          <div
            key={s.user.id}
            className="w-[250px] flex-shrink-0 flex flex-col justify-between p-5 rounded-2xl bg-[#0D1B2A]/60 border border-white/5 backdrop-blur-md hover:border-white/10 hover:bg-[#0D1B2A]/80 transition-all shadow-xl group"
          >
            <div className="space-y-3">
              <div className="flex items-start justify-between gap-2">
                <div className="w-12 h-12 rounded-full overflow-hidden bg-gray-800 border border-white/10 flex-shrink-0">
                  {s.user.avatar_url ? (
                    <img
                      src={s.user.avatar_url}
                      alt={s.user.name}
                      className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center bg-gradient-to-tr from-gray-700 to-gray-800 text-base font-bold text-white uppercase">
                      {s.user.name.slice(0, 1)}
                    </div>
                  )}
                </div>

                {s.mutual_connection_count > 0 && (
                  <div className="flex items-center gap-1 py-1 px-2 rounded-lg bg-orange-500/10 border border-orange-500/20 text-orange-400 text-[10px] font-bold">
                    <Users className="w-3 h-3" />
                    {s.mutual_connection_count} mutual
                  </div>
                )}
              </div>

              <div className="min-w-0">
                <h4 className="text-sm font-bold text-white truncate leading-tight">
                  {s.user.name}
                </h4>
                <p className="text-xs text-gray-400 line-clamp-2 mt-1 leading-normal min-h-[32px]">
                  {s.user.headline || "Practitioner"}
                </p>
              </div>
            </div>

            <div className="mt-4 pt-3 border-t border-white/5 flex items-center justify-between gap-2">
              <span className="text-[10px] font-medium text-gray-500 truncate max-w-[100px]">
                {s.reason}
              </span>
              <ConnectButton
                targetUserId={s.user.id}
                targetUserName={s.user.name}
                currentConnectionStatus="none"
              />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
