"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { create } from "zustand";
import { connectionsClient, Connection, Suggestion } from "@/lib/api/connections";

// Zustand store for tracking optimistic connection status overrides
interface ConnectionsStore {
  statusOverrides: Record<string, "none" | "pending_outgoing" | "pending_incoming" | "accepted" | "blocked">;
  setStatusOverride: (userId: string, status: "none" | "pending_outgoing" | "pending_incoming" | "accepted" | "blocked") => void;
  clearOverrides: () => void;
}

export const useConnectionsStore = create<ConnectionsStore>((set) => ({
  statusOverrides: {},
  setStatusOverride: (userId, status) =>
    set((state) => ({
      statusOverrides: { ...state.statusOverrides, [userId]: status },
    })),
  clearOverrides: () => set({ statusOverrides: {} }),
}));

export function useConnections(page: number = 1, limit: number = 10) {
  return useQuery<Connection[]>({
    queryKey: ["connections", "accepted", page, limit],
    queryFn: () => connectionsClient.getConnections(page, limit),
  });
}

export function usePendingRequests(direction: "incoming" | "outgoing") {
  return useQuery<Connection[]>({
    queryKey: ["connections", "pending", direction],
    queryFn: () => connectionsClient.getPendingRequests(direction),
  });
}

export function useMutualConnections(userId: string, limit: number = 10) {
  return useQuery({
    queryKey: ["connections", "mutual", userId, limit],
    queryFn: () => connectionsClient.getMutualConnections(userId, limit),
    enabled: !!userId,
  });
}

export function useSuggestions(limit: number = 10) {
  return useQuery<Suggestion[]>({
    queryKey: ["connections", "suggestions", limit],
    queryFn: () => connectionsClient.getSuggestions(limit),
  });
}

export function useConnectionStatus(userId: string) {
  return useQuery({
    queryKey: ["connections", "status", userId],
    queryFn: () => connectionsClient.getConnectionStatus(userId),
    enabled: !!userId,
  });
}

export function useSendConnectionRequest() {
  const queryClient = useQueryClient();
  const setStatusOverride = useConnectionsStore((s) => s.setStatusOverride);

  return useMutation({
    mutationFn: ({ targetUserId, note, source }: { targetUserId: string; note?: string; source?: string }) =>
      connectionsClient.sendConnectionRequest(targetUserId, note, source),
    onMutate: async ({ targetUserId }) => {
      // Optimistically set the status to pending_outgoing
      setStatusOverride(targetUserId, "pending_outgoing");
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["connections", "pending", "outgoing"] });
      queryClient.invalidateQueries({ queryKey: ["connections", "suggestions"] });
      queryClient.invalidateQueries({ queryKey: ["connections", "status", variables.targetUserId] });
    },
    onError: (err, variables) => {
      // Rollback on error
      setStatusOverride(variables.targetUserId, "none");
    },
  });
}

export function useAcceptConnection() {
  const queryClient = useQueryClient();
  const setStatusOverride = useConnectionsStore((s) => s.setStatusOverride);

  return useMutation({
    mutationFn: ({ connectionId, targetUserId }: { connectionId: string; targetUserId: string }) =>
      connectionsClient.acceptConnection(connectionId),
    onMutate: async ({ targetUserId }) => {
      setStatusOverride(targetUserId, "accepted");
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["connections", "pending"] });
      queryClient.invalidateQueries({ queryKey: ["connections", "accepted"] });
      queryClient.invalidateQueries({ queryKey: ["profile", variables.targetUserId] });
      queryClient.invalidateQueries({ queryKey: ["connections", "suggestions"] });
      queryClient.invalidateQueries({ queryKey: ["connections", "status", variables.targetUserId] });
    },
    onError: (err, variables) => {
      setStatusOverride(variables.targetUserId, "pending_incoming");
    },
  });
}

export function useDeclineConnection() {
  const queryClient = useQueryClient();
  const setStatusOverride = useConnectionsStore((s) => s.setStatusOverride);

  return useMutation({
    mutationFn: ({ connectionId, targetUserId }: { connectionId: string; targetUserId: string }) =>
      connectionsClient.declineConnection(connectionId),
    onMutate: async ({ targetUserId }) => {
      setStatusOverride(targetUserId, "none");
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["connections", "pending"] });
      queryClient.invalidateQueries({ queryKey: ["connections", "suggestions"] });
      queryClient.invalidateQueries({ queryKey: ["connections", "status", variables.targetUserId] });
    },
    onError: (err, variables) => {
      setStatusOverride(variables.targetUserId, "pending_incoming");
    },
  });
}

export function useRemoveConnection() {
  const queryClient = useQueryClient();
  const setStatusOverride = useConnectionsStore((s) => s.setStatusOverride);

  return useMutation({
    mutationFn: ({ connectionId, targetUserId }: { connectionId: string; targetUserId: string }) =>
      connectionsClient.removeConnection(connectionId),
    onMutate: async ({ targetUserId }) => {
      setStatusOverride(targetUserId, "none");
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["connections", "accepted"] });
      queryClient.invalidateQueries({ queryKey: ["profile", variables.targetUserId] });
      queryClient.invalidateQueries({ queryKey: ["connections", "suggestions"] });
      queryClient.invalidateQueries({ queryKey: ["connections", "status", variables.targetUserId] });
    },
    onError: (err, variables) => {
      setStatusOverride(variables.targetUserId, "accepted");
    },
  });
}

export function useBlockUser() {
  const queryClient = useQueryClient();
  const setStatusOverride = useConnectionsStore((s) => s.setStatusOverride);

  return useMutation({
    mutationFn: ({ targetUserId, reason }: { targetUserId: string; reason?: string }) =>
      connectionsClient.blockUser(targetUserId, reason),
    onMutate: async ({ targetUserId }) => {
      setStatusOverride(targetUserId, "blocked");
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["connections"] });
      queryClient.invalidateQueries({ queryKey: ["connections", "suggestions"] });
      queryClient.invalidateQueries({ queryKey: ["profile", variables.targetUserId] });
      queryClient.invalidateQueries({ queryKey: ["connections", "status", variables.targetUserId] });
    },
    onError: (err, variables) => {
      // Revert if error
      setStatusOverride(variables.targetUserId, "none");
    },
  });
}

export function useUnblockUser() {
  const queryClient = useQueryClient();
  const setStatusOverride = useConnectionsStore((s) => s.setStatusOverride);

  return useMutation({
    mutationFn: (targetUserId: string) => connectionsClient.unblockUser(targetUserId),
    onMutate: async (targetUserId) => {
      setStatusOverride(targetUserId, "none");
    },
    onSuccess: (data, targetUserId) => {
      queryClient.invalidateQueries({ queryKey: ["connections"] });
      queryClient.invalidateQueries({ queryKey: ["connections", "suggestions"] });
      queryClient.invalidateQueries({ queryKey: ["profile", targetUserId] });
      queryClient.invalidateQueries({ queryKey: ["connections", "status", targetUserId] });
    },
  });
}
