import { api } from "./client";

export interface Connection {
  id: string;
  requester_id: string;
  receiver_id: string;
  status: "pending" | "accepted" | "rejected";
  created_at: string;
  updated_at: string;
  requester_name?: string;
  requester_headline?: string;
  requester_photo_url?: string;
  receiver_name?: string;
  receiver_headline?: string;
  receiver_photo_url?: string;
}

export interface ConnectionStatusResponse {
  status: "pending" | "accepted" | "rejected" | "";
  requester_id?: string;
}

export const networkClient = {
  sendRequest: (receiverID: string) =>
    api.post<Connection>("/network/requests", { receiver_id: receiverID }),

  acceptRequest: (id: string) =>
    api.put<{ status: string }>(`/network/requests/${id}/accept`, {}),

  rejectRequest: (id: string) =>
    api.put<{ status: string }>(`/network/requests/${id}/reject`, {}),

  getConnections: () =>
    api.get<Connection[]>("/network/connections"),

  getIncomingRequests: () =>
    api.get<Connection[]>("/network/requests/incoming"),

  getConnectionStatus: (userID: string) =>
    api.get<ConnectionStatusResponse>(`/network/status/${userID}`),
};
