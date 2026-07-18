"use client";

import React, { useState, useEffect, useRef } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import AuthGuard from "@/components/shared/AuthGuard";
import { api, getAccessToken } from "@/lib/api/client";
import { useAuth } from "@/lib/auth/auth-context";
import { useNotifications } from "@/components/shared/Notifications";
import { Paperclip, Pin, Archive, Trash2, Send, Loader2 } from "lucide-react";

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
      <div className="min-h-screen bg-background text-foreground flex flex-col">
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Inbox" }]} />

        {/* Chat Interface Grid */}
        <div className="flex-grow max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8 flex flex-col min-h-0">
          <div className="bg-card border border-border/80 rounded-3xl overflow-hidden grid grid-cols-1 md:grid-cols-3 h-[calc(100vh-180px)] min-h-[500px] shadow-sm">
            
            {/* Left Conversations Sidebar */}
            <div className="border-r border-border/85 flex flex-col min-h-0 bg-card col-span-1">
              <div className="p-5 border-b border-border/85 flex-shrink-0">
                <h1 className="text-lg font-extrabold mb-3 tracking-tight">Messages</h1>
                <div className="flex items-center gap-2 bg-muted/50 border border-border/60 rounded-xl px-3 py-2">
                  <span className="text-muted-foreground text-sm font-semibold select-none">⌕</span>
                  <input
                    placeholder="Search conversations"
                    value={filterQuery}
                    onChange={(e) => setFilterQuery(e.target.value)}
                    className="border-none outline-none bg-transparent text-sm text-foreground w-full placeholder:text-muted-foreground"
                  />
                </div>
              </div>

              {/* List of active threads */}
              <div className="flex-grow overflow-y-auto min-h-0 divide-y divide-border/60">
                {loadingList ? (
                  <div className="flex justify-center items-center py-12">
                    <Loader2 className="h-6 w-6 text-primary animate-spin" />
                  </div>
                ) : (
                  <div className="divide-y divide-border/65">
                    {filteredConvs.map((c) => (
                      <div
                        key={c.id}
                        onClick={() => setActiveConvId(c.id)}
                        className={`p-4 flex gap-3 cursor-pointer transition-all border-l-4 ${
                          activeConvId === c.id
                            ? "bg-primary/5 border-primary"
                            : "bg-transparent border-transparent hover:bg-muted/30"
                        }`}
                      >
                        <div className="w-10 h-10 rounded-full bg-primary/10 border border-primary/20 text-primary flex items-center justify-center font-bold text-sm flex-shrink-0">
                          {c.title?.charAt(0) || "P"}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="flex justify-between items-baseline gap-2">
                            <span className="font-semibold text-sm truncate">{c.title}</span>
                            <span className="text-[10px] text-muted-foreground whitespace-nowrap">
                              {new Date(c.updated_at).toLocaleDateString()}
                            </span>
                          </div>
                          <div className="text-xs text-muted-foreground truncate mt-1">
                            {c.last_message_preview || "No messages yet."}
                          </div>
                        </div>
                      </div>
                    ))}
                    {filteredConvs.length === 0 && (
                      <div className="text-center py-12 text-sm text-muted-foreground">
                        No conversations yet.
                      </div>
                    )}
                  </div>
                )}
              </div>
            </div>

            {/* Right Chat Window */}
            <div className="col-span-2 flex flex-col min-h-0 bg-muted/10">
              {activeConvId ? (
                <>
                  {/* Active Header */}
                  <div className="p-4 border-b border-border/85 flex items-center justify-between bg-card flex-shrink-0">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-full bg-primary/10 border border-primary/20 text-primary flex items-center justify-center font-bold text-sm">
                        {activeConv?.title?.charAt(0) || "P"}
                      </div>
                      <div className="min-w-0">
                        <div className="font-bold text-sm truncate">{activeConv?.title}</div>
                        <div className="text-xs text-muted-foreground">Active Connection</div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() => handlePin(activeConvId, !!activeConv?.is_pinned)}
                        className={`p-2 rounded-full transition-colors ${
                          activeConv?.is_pinned
                            ? "text-primary bg-primary/10"
                            : "text-muted-foreground hover:bg-muted"
                        }`}
                        title={activeConv?.is_pinned ? "Unpin Chat" : "Pin Chat"}
                      >
                        <Pin size={16} />
                      </button>
                      <button
                        onClick={() => handleArchive(activeConvId, !!activeConv?.is_archived)}
                        className={`p-2 rounded-full transition-colors ${
                          activeConv?.is_archived
                            ? "text-primary bg-primary/10"
                            : "text-muted-foreground hover:bg-muted"
                        }`}
                        title={activeConv?.is_archived ? "Unarchive Chat" : "Archive Chat"}
                      >
                        <Archive size={16} />
                      </button>
                    </div>
                  </div>

                  {/* Messages Stream */}
                  <div className="flex-grow overflow-y-auto min-h-0 p-6 flex flex-col gap-4 bg-muted/5">
                    {loadingMessages ? (
                      <div className="flex justify-center items-center h-full">
                        <Loader2 className="h-6 w-6 text-primary animate-spin" />
                      </div>
                    ) : (
                      <>
                        {messages.map((m) => {
                          const isMe = m.sender_id === user?.id;
                          return (
                            <div
                              key={m.id}
                              className={`max-w-[70%] rounded-2xl px-4 py-3 text-sm shadow-sm border ${
                                isMe
                                  ? "self-end bg-primary text-primary-foreground border-primary/10 rounded-tr-none"
                                  : "self-start bg-card text-foreground border-border/80 rounded-tl-none"
                              }`}
                            >
                              <div className="leading-relaxed break-words">{m.body}</div>
                              <div
                                className={`flex items-center gap-2 mt-1.5 text-[10px] ${
                                  isMe ? "text-primary-foreground/75 justify-end" : "text-muted-foreground"
                                }`}
                              >
                                <span>
                                  {new Date(m.created_at).toLocaleTimeString([], {
                                    hour: "2-digit",
                                    minute: "2-digit",
                                  })}
                                </span>
                                {isMe && !m.deleted_at && (
                                  <button
                                    onClick={() => handleDeleteMessage(m.id)}
                                    className="opacity-75 hover:opacity-100 transition-opacity ml-1 cursor-pointer"
                                    title="Delete Message"
                                  >
                                    <Trash2 size={12} />
                                  </button>
                                )}
                              </div>
                            </div>
                          );
                        })}
                        {isPartnerTyping && (
                          <div className="self-start flex gap-2 items-center text-xs text-muted-foreground italic bg-muted/50 px-3 py-1.5 rounded-full">
                            <Loader2 className="h-3 w-3 text-muted-foreground animate-spin" />
                            <span>Typing...</span>
                          </div>
                        )}
                        <div ref={messagesEndRef} />
                      </>
                    )}
                  </div>

                  {/* Composer Panel Footer */}
                  <form
                    onSubmit={handleSendMessage}
                    className="p-4 border-t border-border/85 bg-card flex-shrink-0 flex gap-3 items-center"
                  >
                    <input
                      type="file"
                      id="inbox-file-upload"
                      onChange={handleFileUpload}
                      className="hidden"
                    />
                    <label
                      htmlFor="inbox-file-upload"
                      className="cursor-pointer text-muted-foreground hover:text-foreground transition-colors p-2 rounded-xl hover:bg-muted/50 flex-shrink-0"
                      title="Attach File"
                    >
                      <Paperclip size={18} />
                    </label>
                    <input
                      placeholder="Write a message…"
                      value={composedText}
                      onChange={handleComposerChange}
                      className="flex-grow bg-muted/30 border border-border/60 focus:border-primary/50 rounded-full px-4 py-2.5 text-sm text-foreground outline-none transition-all placeholder:text-muted-foreground"
                    />
                    <button
                      type="submit"
                      disabled={sendingMessage || !composedText.trim()}
                      className="bg-primary text-primary-foreground hover:bg-primary/95 transition-colors px-5 py-2.5 rounded-full font-semibold text-sm flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed flex-shrink-0"
                    >
                      <span>Send</span>
                      <Send size={14} />
                    </button>
                  </form>
                </>
              ) : (
                <div className="flex-1 flex flex-col justify-center items-center text-muted-foreground p-8 bg-muted/5">
                  <span className="text-4xl mb-4 select-none">💬</span>
                  <h3 className="text-base font-bold text-foreground mb-1">Your Messages</h3>
                  <p className="text-xs text-center max-w-sm">
                    Select a conversation from the sidebar or start a new chat from the Network Center.
                  </p>
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
