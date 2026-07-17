import { api } from "./client";

export interface PublicProfileSummary {
  id: string;
  name: string;
  headline: string;
  avatar_url: string;
}

export interface Connection {
  id: string;
  user_a_id: string;
  user_b_id: string;
  status: "pending" | "accepted" | "declined" | "blocked";
  requested_by: string;
  created_at: string;
  responded_at?: string;
  updated_at: string;
  note?: string;
  source?: "search" | "profile_view" | "suggested" | "import";
  user: PublicProfileSummary; // Hydrated details of the OTHER user
}

export interface Suggestion {
  user: PublicProfileSummary;
  mutual_connection_count: number;
  reason: string;
}

export interface MutualConnectionsResponse {
  users: PublicProfileSummary[];
  total: number;
}

export const connectionsClient = {
  sendConnectionRequest: (targetUserId: string, note?: string, source: string = "profile_view") =>
    api.post<{ success: boolean; message: string }>("/connections/request", {
      target_user_id: targetUserId,
      note: note || undefined,
      source,
    }),

  acceptConnection: (connectionId: string) =>
    api.post<{ success: boolean; message: string }>(`/connections/${connectionId}/accept`, {}),

  declineConnection: (connectionId: string) =>
    api.post<{ success: boolean; message: string }>(`/connections/${connectionId}/decline`, {}),

  removeConnection: (connectionId: string) =>
    api.delete<{ success: boolean; message: string }>(`/connections/${connectionId}`),

  blockUser: (userId: string, reason?: string) =>
    api.post<{ success: boolean; message: string }>(`/connections/users/${userId}/block`, {
      reason: reason || undefined,
    }),

  unblockUser: (userId: string) =>
    api.delete<{ success: boolean; message: string }>(`/connections/users/${userId}/block`),

  getConnections: (page: number = 1, limit: number = 10) =>
    api.get<Connection[]>(`/connections?status=accepted&page=${page}&limit=${limit}`),

  getPendingRequests: (direction: "incoming" | "outgoing") =>
    api.get<Connection[]>(`/connections/pending?direction=${direction}`),

  getMutualConnections: (userId: string, limit: number = 10) =>
    api.get<MutualConnectionsResponse>(`/connections/mutual/${userId}?limit=${limit}`),

  getSuggestions: (limit: number = 10) =>
    api.get<Suggestion[]>(`/connections/suggestions?limit=${limit}`),

  getConnectionStatus: (userId: string) =>
    api.get<{ connection_id?: string; status: string; requested_by?: string }>(`/connections/status/${userId}`),
};
