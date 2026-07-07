import { api } from "./client";

export interface Connection {
  id: string;
  requester_id: string;
  receiver_id: string;
  status: "pending" | "accepted" | "rejected" | "blocked" | "declined";
  origin:
    | "mentorship_match"
    | "referral_request"
    | "job_application"
    | "manual_request";
  created_at: string;
  updated_at: string;
  responded_at?: string;
  requester_name?: string;
  requester_headline?: string;
  requester_photo_url?: string;
  receiver_name?: string;
  receiver_headline?: string;
  receiver_photo_url?: string;
}

export interface ConnectionStatusResponse {
  status: "pending" | "accepted" | "rejected" | "blocked" | "declined" | "";
  requester_id?: string;
  origin?: string;
}

export const networkClient = {
  sendRequest: (receiverID: string, origin: string = "manual_request") =>
    api.post<Connection>("/network/requests", {
      receiver_id: receiverID,
      origin,
    }),

  acceptRequest: (id: string) =>
    api.put<{ status: string }>(`/network/requests/${id}/accept`, {}),

  rejectRequest: (id: string) =>
    api.put<{ status: string }>(`/network/requests/${id}/reject`, {}),

  blockUser: (targetID: string) =>
    api.post<{ status: string }>("/network/block", { target_id: targetID }),

  unconnect: (userID: string) =>
    api.delete<{ status: string }>(`/network/connections/${userID}`),

  getConnections: () => api.get<Connection[]>("/network/connections"),

  getIncomingRequests: () =>
    api.get<Connection[]>("/network/requests/incoming"),

  getConnectionStatus: (userID: string) =>
    api.get<ConnectionStatusResponse>(`/network/status/${userID}`),
};
