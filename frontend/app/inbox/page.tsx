"use client";

import React, { useState, useEffect, useRef } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import AuthGuard from "@/components/shared/AuthGuard";
import { api, getAccessToken } from "@/lib/api/client";
import { useAuth } from "@/lib/auth/auth-context";
import { useNotifications } from "@/components/shared/Notifications";
import { CircularProgress } from "@mui/material";
import { Paperclip, Pin, Archive, Trash2, Send, CornerDownLeft, Smile } from "lucide-react";

interface Conversation {
  id: string;
  type: string;
  title: string;
  participants: string[];
  updated_at: string;
  unread_count: number;
  last_message_preview?: string;
  is_pinned?: boolean;
  is_archived?: boolean;
}

interface Message {
  id: string;
  sender_id: string;
  body: string;
  content_type: string;
  created_at: string;
  edited_at?: string;
  deleted_at?: string;
}

export default function InboxPage() {
  const { user } = useAuth();
  const { showNotification } = useNotifications();

  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [activeConvId, setActiveConvId] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [loadingList, setLoadingList] = useState(true);
  const [loadingMessages, setLoadingMessages] = useState(false);

  // Message composer
  const [composedText, setComposedText] = useState("");
  const [sendingMessage, setSendingMessage] = useState(false);

  // Search filter
  const [filterQuery, setFilterQuery] = useState("");

  // Presence & Typing indicators
  const [isPartnerTyping, setIsPartnerTyping] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const typingTimerRef = useRef<NodeJS.Timeout | null>(null);

  // Read URL query parameter if user was redirected from Network/Connections page
  useEffect(() => {
    if (typeof window !== "undefined") {
      const searchParams = new URLSearchParams(window.location.search);
      const convId = searchParams.get("convId");
      if (convId) {
        setActiveConvId(convId);
      }
    }
  }, []);

  // Fetch conversations list on mount
  useEffect(() => {
    fetchConversations();
  }, []);

  // Fetch messages when active conversation changes & set up WebSocket connection
  useEffect(() => {
    if (!activeConvId) return;
    fetchMessages(activeConvId);

    // Initialize real-time WebSocket connection
    const token = getAccessToken();
    if (token) {
      const wsProto = window.location.protocol === "https:" ? "wss:" : "ws:";
      // Check if we are in local development vs production proxy
      const wsHost = window.location.host;
      const socket = new WebSocket(`${wsProto}//${wsHost}/api/v1/ws?token=${encodeURIComponent(token)}`);
      socketRef.current = socket;

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.kind === "message" && data.conversation_id === activeConvId) {
            setMessages((prev) => {
              if (prev.some((m) => m.id === data.id)) return prev;
              return [
                ...prev,
                {
                  id: data.id,
                  sender_id: data.sender_id,
                  body: data.body,
                  content_type: data.content_type || "text",
                  created_at: data.created_at || new Date().toISOString(),
                },
              ];
            });
            // Mark read
            api.post(`/conversations/${activeConvId}/read`, {}).catch(() => {});
          } else if (data.kind === "typing" && data.conversation_id === activeConvId && data.sender_id !== user?.id) {
            setIsPartnerTyping(true);
            setTimeout(() => setIsPartnerTyping(false), 3000);
          }
        } catch (_) {}
      };

      socket.onclose = () => {
        socketRef.current = null;
      };

      return () => {
        socket.close();
      };
    }
  }, [activeConvId]);

  // Scroll to bottom of message list on updates
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, isPartnerTyping]);

  const fetchConversations = async () => {
    setLoadingList(true);
    try {
      const data = await api.get<{ conversations: Conversation[] }>("/conversations");
      setConversations(data.conversations || []);
    } catch (err: any) {
      showNotification("Failed to load inbox", "error");
    } finally {
      setLoadingList(false);
    }
  };

  const fetchMessages = async (convId: string) => {
    setLoadingMessages(true);
    try {
      const data = await api.get<{ messages: Message[] }>(`/conversations/${convId}/messages`);
      setMessages(data.messages || []);
      // Mark read
      await api.post(`/conversations/${convId}/read`, {});
    } catch (_) {
      setMessages([]);
    } finally {
      setLoadingMessages(false);
    }
  };

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!composedText.trim() || !activeConvId) return;

    setSendingMessage(true);
    try {
      const sent = await api.post<Message>(`/conversations/${activeConvId}/messages`, {
        body: composedText,
        content_type: "text",
      });

      setMessages((prev) => [...prev, sent]);
      setComposedText("");

      // Update last message in the sidebar
      setConversations((prev) =>
        prev.map((c) => (c.id === activeConvId ? { ...c, last_message_preview: sent.body } : c))
      );
    } catch (err: any) {
      showNotification("Failed to send message", "error");
    } finally {
      setSendingMessage(false);
    }
  };

  const handleComposerChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setComposedText(e.target.value);

    // Send typing notification to WebSocket
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      if (typingTimerRef.current) clearTimeout(typingTimerRef.current);
      socketRef.current.send(JSON.stringify({ type: "typing", conversation_id: activeConvId }));
      typingTimerRef.current = setTimeout(() => {}, 2000);
    }
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0 || !activeConvId) return;

    const file = files[0];
    const formData = new FormData();
    formData.append("body", `Shared a file: ${file.name}`);
    formData.append("content_type", "file");
    formData.append("attachment", file);

    try {
      showNotification("Uploading attachment...", "info");
      const sent = await api.post<Message>(`/conversations/${activeConvId}/messages`, {
        body: formData,
      });
      setMessages((prev) => [...prev, sent]);
      showNotification("File shared successfully!", "success");
    } catch (err: any) {
      showNotification("Failed to upload file", "error");
    }
  };

  const handlePin = async (convId: string, currentPinStatus: boolean) => {
    try {
      await api.post(`/conversations/${convId}/pin`, { pin: !currentPinStatus });
      setConversations((prev) =>
        prev.map((c) => (c.id === convId ? { ...c, is_pinned: !currentPinStatus } : c))
      );
      showNotification(currentPinStatus ? "Unpinned conversation." : "Pinned conversation!", "success");
    } catch (_) {}
  };

  const handleArchive = async (convId: string, currentArchiveStatus: boolean) => {
    try {
      await api.post(`/conversations/${convId}/archive`, { archive: !currentArchiveStatus });
      setConversations((prev) =>
        prev.map((c) => (c.id === convId ? { ...c, is_archived: !currentArchiveStatus } : c))
      );
      showNotification(currentArchiveStatus ? "Restored conversation." : "Archived conversation.", "success");
    } catch (_) {}
  };

  const handleDeleteMessage = async (messageID: string) => {
    if (!window.confirm("Are you sure you want to delete this message?")) return;
    try {
      await api.delete(`/conversations/${activeConvId}/messages/${messageID}`);
      setMessages((prev) =>
        prev.map((m) => (m.id === messageID ? { ...m, body: "This message was deleted.", deleted_at: new Date().toISOString() } : m))
      );
      showNotification("Message deleted.", "success");
    } catch (_) {}
  };

  // Filter threads
  const filteredConvs = conversations.filter((c) =>
    c.title?.toLowerCase().includes(filterQuery.toLowerCase())
  );

  const activeConv = conversations.find((c) => c.id === activeConvId);

  return (
    <AuthGuard>
      <div
        style={{
          background: "#FBF7F2",
          fontFamily: "'Public Sans', sans-serif",
          color: "#2B2620",
          minHeight: "100vh",
          display: "flex",
          flexDirection: "column",
        }}
      >
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Inbox" }]} />

        {/* Chat Interface Grid */}
        <div style={{ flex: 1, maxWidth: "1240px", width: "100%", margin: "0 auto", padding: "clamp(20px,3vw,32px) 40px", minHeight: 0 }}>
          <div style={{ background: "#fff", border: "1px solid #EFE7DC", borderRadius: "22px", overflow: "hidden", display: "grid", gridTemplateColumns: "340px 1fr", height: "calc(100vh - 150px)", minHeight: "520px" }}>
            
            {/* Left Conversations Sidebar */}
            <div style={{ borderRight: "1px solid #EFE7DC", display: "flex", flexDirection: "column", minHeight: 0 }}>
              <div style={{ padding: "22px 22px 16px", flex: "none" }}>
                <h1 style={{ fontWeight: 800, fontSize: "22px", margin: "0 0 14px 0" }}>Messages</h1>
                <div style={{ display: "flex", alignItems: "center", gap: "10px", background: "#F3ECE2", borderRadius: "10px", padding: "10px 14px" }}>
                  <span style={{ color: "#8A8175", fontSize: "16px" }}>⌕</span>
                  <input
                    placeholder="Search conversations"
                    value={filterQuery}
                    onChange={(e) => setFilterQuery(e.target.value)}
                    style={{ border: "none", outline: "none", background: "transparent", fontSize: "14px", color: "#2B2620", width: "100%", fontFamily: "inherit" }}
                  />
                </div>
              </div>

              {/* List of active threads */}
              <div style={{ flex: 1, overflowY: "auto", minHeight: 0 }}>
                {loadingList ? (
                  <div style={{ display: "flex", justifyContent: "center", alignItems: "center", padding: "32px" }}>
                    <CircularProgress size={24} style={{ color: "#C2683C" }} />
                  </div>
                ) : (
                  <div>
                    {filteredConvs.map((c) => (
                      <div
                        key={c.id}
                        onClick={() => setActiveConvId(c.id)}
                        style={{
                          padding: "16px 22px",
                          borderBottom: "1px solid #F3ECE2",
                          cursor: "pointer",
                          display: "flex",
                          gap: "13px",
                          background: activeConvId === c.id ? "#FBF7F2" : "#ffffff",
                          borderLeft: activeConvId === c.id ? "4px solid #C2683C" : "4px solid transparent",
                          position: "relative",
                        }}
                      >
                        <div style={{ width: "44px", height: "44px", borderRadius: "50%", background: "#4F7C6A", color: "#fff", display: "flex", alignItems: "center", justifyContent: "center", fontWeight: 700, fontSize: "15px" }}>
                          {c.title?.charAt(0) || "P"}
                        </div>
                        <div style={{ flex: 1, minWidth: 0 }}>
                          <div style={{ display: "flex", justifyContent: "space-between", gap: "8px" }}>
                            <span style={{ fontWeight: 600, fontSize: "15px" }}>{c.title}</span>
                            <span style={{ fontSize: "12px", color: "#A89C8A" }}>
                              {new Date(c.updated_at).toLocaleDateString()}
                            </span>
                          </div>
                          <div style={{ fontSize: "13px", color: "#5B554C", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis", marginTop: "4px" }}>
                            {c.last_message_preview || "No messages yet."}
                          </div>
                        </div>
                      </div>
                    ))}
                    {filteredConvs.length === 0 && (
                      <div style={{ textAlign: "center", padding: "40px 20px", color: "#8A8175", fontSize: "14px" }}>
                        No conversations yet.
                      </div>
                    )}
                  </div>
                )}
              </div>
            </div>

            {/* Right Chat Window */}
            <div style={{ display: "flex", flexDirection: "column", minHeight: 0 }}>
              {activeConvId ? (
                <>
                  {/* Active Header */}
                  <div style={{ padding: "18px 26px", borderBottom: "1px solid #EFE7DC", display: "flex", alignItems: "center", gap: "14px", flex: "none", background: "#ffffff" }}>
                    <div style={{ width: "44px", height: "44px", borderRadius: "50%", background: "#4F7C6A", color: "#fff", display: "flex", alignItems: "center", justifyContent: "center", fontWeight: 700, fontSize: "15px" }}>
                      {activeConv?.title?.charAt(0) || "P"}
                    </div>
                    <div style={{ flex: 1, minWidth: 0 }}>
                      <div style={{ fontWeight: 600, fontSize: "16px" }}>{activeConv?.title}</div>
                      <div style={{ fontSize: "13px", color: "#8A8175" }}>Active Connection</div>
                    </div>
                    <div style={{ display: "flex", gap: "12px" }}>
                      <button onClick={() => handlePin(activeConvId, !!activeConv?.is_pinned)} style={{ border: "none", background: "transparent", color: activeConv?.is_pinned ? "#C2683C" : "#8A8175", cursor: "pointer" }}>
                        <Pin size={18} />
                      </button>
                      <button onClick={() => handleArchive(activeConvId, !!activeConv?.is_archived)} style={{ border: "none", background: "transparent", color: activeConv?.is_archived ? "#C2683C" : "#8A8175", cursor: "pointer" }}>
                        <Archive size={18} />
                      </button>
                    </div>
                  </div>

                  {/* Messages Stream */}
                  <div style={{ flex: 1, overflowY: "auto", minHeight: 0, padding: "26px", display: "flex", flexDirection: "column", gap: "16px", background: "#FCFAF7" }}>
                    {loadingMessages ? (
                      <div style={{ display: "flex", justifyContent: "center", alignItems: "center", height: "100%" }}>
                        <CircularProgress style={{ color: "#C2683C" }} />
                      </div>
                    ) : (
                      <>
                        {messages.map((m) => {
                          const isMe = m.sender_id === user?.id;
                          return (
                            <div
                              key={m.id}
                              style={{
                                alignSelf: isMe ? "flex-end" : "flex-start",
                                maxWidth: "74%",
                                background: isMe ? "#C2683C" : "#ffffff",
                                color: isMe ? "#ffffff" : "#2B2620",
                                border: isMe ? "none" : "1px solid #EFE7DC",
                                borderRadius: isMe ? "16px 16px 4px 16px" : "16px 16px 16px 4px",
                                padding: "14px 18px",
                                position: "relative",
                                boxShadow: "0 2px 8px rgba(43,38,32,0.02)",
                              }}
                            >
                              <div style={{ fontSize: "15px", lineHeight: 1.55 }}>
                                {m.body}
                              </div>
                              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginTop: "6px", fontSize: "11px", color: isMe ? "rgba(255,255,255,0.7)" : "#A89C8A" }}>
                                <span>{new Date(m.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</span>
                                {isMe && !m.deleted_at && (
                                  <button onClick={() => handleDeleteMessage(m.id)} style={{ border: "none", background: "transparent", color: "rgba(255,255,255,0.8)", cursor: "pointer", marginLeft: "10px" }}>
                                    <Trash2 size={12} />
                                  </button>
                                )}
                              </div>
                            </div>
                          );
                        })}
                        {isPartnerTyping && (
                          <div style={{ alignSelf: "flex-start", background: "transparent", padding: "6px 12px", fontSize: "13px", color: "#8A8175", display: "flex", gap: "6px", alignItems: "center" }}>
                            <CircularProgress size={12} style={{ color: "#8A8175" }} />
                            <span>Typing...</span>
                          </div>
                        )}
                        <div ref={messagesEndRef} />
                      </>
                    )}
                  </div>

                  {/* Composer Panel Footer */}
                  <form onSubmit={handleSendMessage} style={{ padding: "16px 22px", borderTop: "1px solid #EFE7DC", flex: "none", display: "flex", gap: "12px", alignItems: "center", background: "#ffffff" }}>
                    <input
                      type="file"
                      id="inbox-file-upload"
                      onChange={handleFileUpload}
                      style={{ display: "none" }}
                    />
                    <label htmlFor="inbox-file-upload" style={{ cursor: "pointer", color: "#8A8175" }}>
                      <Paperclip size={20} />
                    </label>
                    <input
                      placeholder="Write a message…"
                      value={composedText}
                      onChange={handleComposerChange}
                      style={{ flex: 1, border: "1px solid #E2D9CC", borderRadius: "100px", padding: "13px 20px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7", fontFamily: "inherit" }}
                    />
                    <button
                      type="submit"
                      disabled={sendingMessage || !composedText.trim()}
                      style={{ border: "none", background: "#C2683C", color: "#fff", padding: "13px 26px", borderRadius: "100px", cursor: "pointer", fontWeight: 600, fontSize: "15px", display: "flex", alignItems: "center", gap: "6px" }}
                    >
                      <span>Send</span>
                      <Send size={16} />
                    </button>
                  </form>
                </>
              ) : (
                <div style={{ flex: 1, display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", color: "#8A8175", padding: "32px", background: "#FCFAF7" }}>
                  <span style={{ fontSize: "48px" }}>💬</span>
                  <h3 style={{ margin: "16px 0 6px 0", fontSize: "18px", fontWeight: 700, color: "#2B2620" }}>Your Messages</h3>
                  <p style={{ margin: 0, fontSize: "14px" }}>Select a conversation from the sidebar or start a new chat from the Network Center.</p>
                </div>
              )}
            </div>

          </div>
        </div>

        <SiteFooter />
      </div>
    </AuthGuard>
  );
}
