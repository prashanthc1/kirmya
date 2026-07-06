import { api, getAccessToken } from "./client";
import { config } from "@/lib/config";

export interface Conversation {
  id: string;
  type: "direct" | "group";
  title?: string;
  participants: string[];
  updated_at: string;
  unread_count: number;
  last_message_preview?: string;
  is_pinned: boolean;
  is_archived: boolean;
}

export interface Message {
  id: string;
  sender_id: string;
  body: string;
  content_type: "text" | "image" | "file" | "system";
  created_at: string;
  edited_at?: string;
  deleted_at?: string;
}

export const messagingClient = {
  getConversations: () =>
    api.get<{ conversations: Conversation[] }>("/conversations"),

  startConversation: (participantIDs: string[], title: string = "") =>
    api.post<Conversation>("/conversations", { participant_ids: participantIDs, title }),

  getMessages: (conversationID: string, searchQuery?: string) => {
    const path = `/conversations/${conversationID}/messages` + (searchQuery ? `?q=${encodeURIComponent(searchQuery)}` : "");
    return api.get<{ messages: Message[] }>(path);
  },

  sendMessage: (conversationID: string, body: string, contentType: string = "text") =>
    api.post<Message>(`/conversations/${conversationID}/messages`, { body, content_type: contentType }),

  sendAttachment: (conversationID: string, formData: FormData) =>
    api.post<Message>(`/conversations/${conversationID}/messages`, undefined, {
      body: formData,
      headers: {
        // Fetch/XHR sets boundary automatically if header Content-Type is omitted
      },
    }),

  deleteMessage: (conversationID: string, messageID: string) =>
    api.delete<{ deleted: boolean }>(`/conversations/${conversationID}/messages/${messageID}`),

  markRead: (conversationID: string) =>
    api.post<{ read: boolean }>(`/conversations/${conversationID}/read`),

  sendTypingIndicator: (conversationID: string) =>
    api.post<{ ok: boolean }>(`/conversations/${conversationID}/typing`),

  archiveConversation: (conversationID: string, archive: boolean) =>
    api.post<{ archived: boolean }>(`/conversations/${conversationID}/archive`, { archive }),

  pinConversation: (conversationID: string, pin: boolean) =>
    api.post<{ pinned: boolean }>(`/conversations/${conversationID}/pin`, { pin }),
};

export interface WSEvent {
  kind: "message" | "typing" | "read";
  conversation_id: string;
  sender_id?: string;
  reader_id?: string;
  id?: string;
  body?: string;
  content_type?: string;
  created_at?: string;
  at?: string;
}

export interface WebSocketCallbacks {
  onMessage?: (event: WSEvent) => void;
  onTyping?: (event: WSEvent) => void;
  onRead?: (event: WSEvent) => void;
  onOpen?: () => void;
  onClose?: (code: number, reason: string) => void;
  onError?: (err: Event) => void;
}

export class KirmyaWebSocketClient {
  private socket: WebSocket | null = null;
  private reconnectTimeout: number | null = null;
  private reconnectDelay = 1000;
  private shouldReconnect = true;
  private pingInterval: number | null = null;

  constructor(private callbacks: WebSocketCallbacks) {}

  public connect(): void {
    if (typeof window === "undefined") return;

    this.shouldReconnect = true;
    const token = getAccessToken();
    if (!token) {
      console.warn("[ws] no access token available; delaying connection");
      this.reconnectLater();
      return;
    }

    const wsBase = config.apiBase.replace(/^http/, "ws").replace(/\/+$/, "");
    const wsUrl = `${wsBase}/api/v1/ws?token=${encodeURIComponent(token)}`;

    console.log("[ws] connecting...");
    this.socket = new WebSocket(wsUrl);

    this.socket.onopen = () => {
      console.log("[ws] connected");
      this.reconnectDelay = 1000; // Reset backoff
      this.callbacks.onOpen?.();
      this.startPing();
    };

    this.socket.onclose = (event) => {
      console.log(`[ws] closed: code=${event.code}, reason=${event.reason}`);
      this.stopPing();
      this.callbacks.onClose?.(event.code, event.reason);
      if (this.shouldReconnect) {
        this.reconnectLater();
      }
    };

    this.socket.onerror = (err) => {
      console.error("[ws] error:", err);
      this.callbacks.onError?.(err);
    };

    this.socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data);
        if (payload.type === "ping" || payload.type === "pong") {
          return;
        }

        const ev = payload as WSEvent;
        switch (ev.kind) {
          case "message":
            this.callbacks.onMessage?.(ev);
            break;
          case "typing":
            this.callbacks.onTyping?.(ev);
            break;
          case "read":
            this.callbacks.onRead?.(ev);
            break;
        }
      } catch (err) {
        console.error("[ws] parsing message failed:", err);
      }
    };
  }

  public sendTyping(conversationID: string): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ type: "typing", conversation_id: conversationID }));
    }
  }

  public sendPing(): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ type: "ping" }));
    }
  }

  public close(): void {
    this.shouldReconnect = false;
    this.stopPing();
    if (this.reconnectTimeout) {
      window.clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }

  private reconnectLater(): void {
    if (this.reconnectTimeout) return;

    console.log(`[ws] reconnecting in ${this.reconnectDelay}ms`);
    this.reconnectTimeout = window.setTimeout(() => {
      this.reconnectTimeout = null;
      this.connect();
    }, this.reconnectDelay);

    // Exponential backoff capped at 30 seconds
    this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
  }

  private startPing(): void {
    this.stopPing();
    this.pingInterval = window.setInterval(() => {
      this.sendPing();
    }, 30000);
  }

  private stopPing(): void {
    if (this.pingInterval) {
      window.clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }
}
