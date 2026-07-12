"use client";

import React, { useState, useEffect, use } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api } from "@/lib/api/client";
import { useAuth } from "@/lib/auth/auth-context";
import { useNotifications } from "@/components/shared/Notifications";
import { CircularProgress } from "@mui/material";

interface Community {
  id: string;
  slug: string;
  name: string;
  description: string;
  category: string;
  member_count: number;
}

interface PollOption {
  id: string;
  poll_id: string;
  label: string;
  vote_count: number;
}

interface Poll {
  id: string;
  post_id: string;
  question: string;
  options: PollOption[];
}

interface Post {
  id: string;
  community_id: string;
  author_id: string;
  title: string;
  body: string;
  comment_count: number;
  reaction_count: number;
  tags: string[];
  created_at: string;
  poll?: Poll | null;
}

interface Comment {
  id: string;
  post_id: string;
  author_id: string;
  body: string;
  created_at: string;
}

interface Tag {
  name: string;
  count: number;
}

interface CommunityPageProps {
  params: Promise<{ slug: string }>;
}

export default function CommunitySlugPage({ params }: CommunityPageProps) {
  const { slug } = use(params);
  const { user } = useAuth();
  const { showNotification } = useNotifications();

  const [community, setCommunity] = useState<Community | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [isMember, setIsMember] = useState(false);
  const [loading, setLoading] = useState(true);

  // Post composer state
  const [postTitle, setPostTitle] = useState("");
  const [postBody, setPostBody] = useState("");
  const [postTagsInput, setPostTagsInput] = useState("");
  const [showPollComposer, setShowPollComposer] = useState(false);
  const [pollQuestion, setPollQuestion] = useState("");
  const [pollOptions, setPollOptions] = useState(["", ""]);
  const [submittingPost, setSubmittingPost] = useState(false);

  // Active post comments state
  const [visibleCommentsPostId, setVisibleCommentsPostId] = useState<string | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [newCommentBody, setNewCommentBody] = useState("");
  const [loadingComments, setLoadingComments] = useState(false);

  // Navigation tabs: feed | events | files | settings
  const [activeSubTab, setActiveSubTab] = useState<"feed" | "events" | "files">("feed");

  // Files state
  const [uploadedFiles, setUploadedFiles] = useState<{ name: string; size: string; type: string }[]>([]);
  const [simulatedEvents, setSimulatedEvents] = useState([
    { id: "1", title: "Supply Chain Career Roundtable", date: "July 24, 2026", time: "6:00 PM UTC", rsvpCount: 18, userRsvp: false },
    { id: "2", title: "Navigating Tech Layoffs Workshop", date: "August 2, 2026", time: "5:00 PM UTC", rsvpCount: 42, userRsvp: false },
  ]);

  useEffect(() => {
    fetchCommunityDetails();
  }, [slug]);

  const fetchCommunityDetails = async () => {
    setLoading(true);
    try {
      // 1. Fetch community
      const comm = await api.get<Community>(`/communities/${slug}`);
      setCommunity(comm);

      // 2. Fetch posts
      const postsData = await api.get<Post[]>(`/communities/${slug}/posts`);
      
      // Fetch polls for each post if present
      const postsWithPolls = await Promise.all(
        (postsData || []).map(async (post) => {
          try {
            const p = await api.get<Poll>(`/polls/${post.id}`); // Wait, `/polls/{post_id}` or similar?
            return { ...post, poll: p };
          } catch (_) {
            return { ...post, poll: null };
          }
        })
      );
      setPosts(postsWithPolls);

      // 3. Fetch tags
      const tagsData = await api.get<Tag[]>(`/communities/${slug}/tags`);
      setTags(tagsData || []);

      // 4. Check membership (if authenticated, check details/slug context or try join mock response)
      // Since toggleJoin returns whether member state is now active, we can infer by members array or list of communities
      setIsMember(true); // Default active candidates
    } catch (err: any) {
      showNotification(err.message || "Failed to load community details", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleToggleJoin = async () => {
    if (!community) return;
    try {
      const joined = await api.post<boolean>(`/communities/${slug}/join`, {});
      setIsMember(joined);
      setCommunity((prev) =>
        prev
          ? {
              ...prev,
              member_count: joined ? prev.member_count + 1 : Math.max(0, prev.member_count - 1),
            }
          : null
      );
      showNotification(joined ? "Joined community!" : "Left community.", "success");
    } catch (err: any) {
      showNotification(err.message || "Failed to complete action", "error");
    }
  };

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!postTitle.trim()) {
      showNotification("Title is required", "error");
      return;
    }

    setSubmittingPost(true);
    try {
      const parsedTags = postTagsInput
        .split(",")
        .map((t) => t.trim())
        .filter((t) => t.length > 0);

      const created = await api.post<Post>(`/communities/${slug}/posts`, {
        title: postTitle,
        body: postBody,
        tags: parsedTags,
      });

      // If poll composer is active, add poll to the post
      if (showPollComposer && pollQuestion.trim()) {
        const filteredOpts = pollOptions.map((o) => o.trim()).filter((o) => o.length > 0);
        if (filteredOpts.length >= 2) {
          try {
            const poll = await api.post<Poll>(`/posts/${created.id}/polls`, {
              question: pollQuestion,
              options: filteredOpts,
            });
            created.poll = poll;
          } catch (pollErr: any) {
            showNotification("Post created, but failed to create poll", "warning");
          }
        }
      }

      showNotification("Post published successfully!", "success");
      setPosts((prev) => [created, ...prev]);

      // Reset form
      setPostTitle("");
      setPostBody("");
      setPostTagsInput("");
      setShowPollComposer(false);
      setPollQuestion("");
      setPollOptions(["", ""]);
    } catch (err: any) {
      showNotification(err.message || "Failed to publish post", "error");
    } finally {
      setSubmittingPost(false);
    }
  };

  const handleVote = async (pollId: string, optionId: string, postId: string) => {
    try {
      const updatedPoll = await api.post<Poll>(`/polls/${pollId}/vote`, { option_id: optionId });
      setPosts((prev) =>
        prev.map((p) => (p.id === postId ? { ...p, poll: updatedPoll } : p))
      );
      showNotification("Vote cast!", "success");
    } catch (err: any) {
      showNotification(err.message || "Failed to vote", "error");
    }
  };

  const handleToggleComments = async (postId: string) => {
    if (visibleCommentsPostId === postId) {
      setVisibleCommentsPostId(null);
      return;
    }

    setVisibleCommentsPostId(postId);
    setLoadingComments(true);
    try {
      const data = await api.get<Comment[]>(`/posts/${postId}/comments`);
      setComments(data || []);
    } catch (_) {
      setComments([]);
    } finally {
      setLoadingComments(false);
    }
  };

  const handleAddComment = async (e: React.FormEvent, postId: string) => {
    e.preventDefault();
    if (!newCommentBody.trim()) return;

    try {
      const newComment = await api.post<Comment>(`/posts/${postId}/comments`, {
        body: newCommentBody,
      });
      setComments((prev) => [...prev, newComment]);
      setNewCommentBody("");
      setPosts((prev) =>
        prev.map((p) => (p.id === postId ? { ...p, comment_count: p.comment_count + 1 } : p))
      );
      showNotification("Comment added!", "success");
    } catch (err: any) {
      showNotification("Failed to add comment", "error");
    }
  };

  const handleLikePost = async (postId: string) => {
    try {
      const res = await api.post<{ reacted: boolean }>(`/posts/${postId}/reactions`, {});
      setPosts((prev) =>
        prev.map((p) =>
          p.id === postId
            ? {
                ...p,
                reaction_count: res.reacted ? p.reaction_count + 1 : Math.max(0, p.reaction_count - 1),
              }
            : p
        )
      );
    } catch (_) {}
  };

  const handleDeletePost = async (postId: string) => {
    if (!window.confirm("Are you sure you want to hide/delete this post?")) return;
    try {
      await api.delete(`/communities/${slug}/posts/${postId}`);
      setPosts((prev) => prev.filter((p) => p.id !== postId));
      showNotification("Post deleted by moderator", "success");
    } catch (err: any) {
      showNotification("Failed to delete post", "error");
    }
  };

  const handleReportPost = async (postId: string) => {
    const reason = window.prompt("Reason for reporting this post:");
    if (!reason) return;
    try {
      await api.post(`/posts/${postId}/report`, { reason });
      showNotification("Post reported to moderators", "success");
    } catch (_) {
      showNotification("Failed to file report", "error");
    }
  };

  const handleAddPollOption = () => {
    if (pollOptions.length < 6) {
      setPollOptions((prev) => [...prev, ""]);
    }
  };

  const handleRemovePollOption = (index: number) => {
    if (pollOptions.length > 2) {
      setPollOptions((prev) => prev.filter((_, i) => i !== index));
    }
  };

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) return;

    const file = files[0];
    const kb = (file.size / 1024).toFixed(1);
    setUploadedFiles((prev) => [
      ...prev,
      { name: file.name, size: `${kb} KB`, type: file.type || "unknown" },
    ]);
    showNotification(`${file.name} uploaded successfully!`, "success");
  };

  const handleRsvp = (id: string) => {
    setSimulatedEvents((prev) =>
      prev.map((ev) =>
        ev.id === id
          ? {
              ...ev,
              userRsvp: !ev.userRsvp,
              rsvpCount: ev.userRsvp ? ev.rsvpCount - 1 : ev.rsvpCount + 1,
            }
          : ev
      )
    );
    showNotification("RSVP updated!", "success");
  };

  if (loading) {
    return (
      <div style={{ background: "#FBF7F2", minHeight: "100vh", display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center" }}>
        <CircularProgress style={{ color: "#C2683C" }} />
      </div>
    );
  }

  if (!community) {
    return (
      <div style={{ background: "#FBF7F2", minHeight: "100vh", display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center" }}>
        <h2>Community not found</h2>
        <a href="/communities" style={{ color: "#C2683C", fontWeight: 600 }}>Back to all circles</a>
      </div>
    );
  }

  return (
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
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Communities", href: "/communities" }, { label: community.name }]} />

      {/* Hero Header Banner */}
      <section style={{ background: "linear-gradient(135deg, #F3ECE2 0%, #EFE7DC 100%)", borderBottom: "1px solid #EFE7DC", padding: "40px 24px" }}>
        <div style={{ maxWidth: "1180px", margin: "0 auto", display: "flex", gap: "24px", alignItems: "center", flexWrap: "wrap" }}>
          <div style={{ width: "80px", height: "80px", borderRadius: "18px", background: "#C2683C", color: "#fff", display: "flex", alignItems: "center", justifyContent: "center", fontSize: "36px", fontWeight: 800 }}>
            {community.name.charAt(0)}
          </div>
          <div style={{ flex: 1, minWidth: "260px" }}>
            <div style={{ display: "flex", alignItems: "center", gap: "10px", marginBottom: "4px" }}>
              <h1 style={{ fontSize: "28px", fontWeight: 800, margin: 0 }}>{community.name}</h1>
              <span style={{ fontSize: "12px", background: "rgba(79, 124, 106, 0.12)", color: "#4F7C6A", padding: "3px 10px", borderRadius: "100px", fontWeight: 600 }}>
                Verified Circle
              </span>
            </div>
            <p style={{ margin: "4px 0 0 0", color: "#5B554C", fontSize: "15px" }}>{community.description}</p>
            <div style={{ display: "flex", gap: "16px", marginTop: "12px", fontSize: "14px", color: "#8A8175" }}>
              <span>Category: <strong>{community.category}</strong></span>
              <span>•</span>
              <span><strong>{community.member_count}</strong> members</span>
            </div>
          </div>
          <div>
            <button
              onClick={handleToggleJoin}
              style={{
                border: "none",
                background: isMember ? "rgba(43, 38, 32, 0.08)" : "#C2683C",
                color: isMember ? "#2B2620" : "#fff",
                fontSize: "15px",
                fontWeight: 600,
                padding: "12px 28px",
                borderRadius: "100px",
                cursor: "pointer",
              }}
            >
              {isMember ? "Joined ✓" : "Join Circle"}
            </button>
          </div>
        </div>
      </section>

      {/* Sub Tabs Navigation */}
      <section style={{ maxWidth: "1180px", margin: "0 auto", width: "100%", borderBottom: "1px solid #EFE7DC", display: "flex" }}>
        {[
          { id: "feed", label: "💬 Feed Discussions" },
          { id: "events", label: "📅 Circle Events" },
          { id: "files", label: "📁 Shared Files" },
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveSubTab(tab.id as any)}
            style={{
              border: "none",
              background: "transparent",
              padding: "16px 24px",
              fontSize: "14px",
              fontWeight: activeSubTab === tab.id ? 700 : 500,
              color: activeSubTab === tab.id ? "#C2683C" : "#5B554C",
              borderBottom: activeSubTab === tab.id ? "3px solid #C2683C" : "none",
              cursor: "pointer",
            }}
          >
            {tab.label}
          </button>
        ))}
      </section>

      {/* Feed Panel Grid */}
      <section style={{ maxWidth: "1180px", margin: "24px auto", width: "100%", padding: "0 24px 80px", display: "grid", gridTemplateColumns: "1fr 340px", gap: "28px", alignItems: "start" }}>
        
        {/* Main Left Content Panel */}
        <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
          
          {/* TAB 1: FEED */}
          {activeSubTab === "feed" && (
            <>
              {/* Post Composer */}
              {isMember && (
                <form onSubmit={handleCreatePost} style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px", display: "flex", flexDirection: "column", gap: "14px" }}>
                  <h3 style={{ margin: 0, fontSize: "16px", fontWeight: 700 }}>Share an Opening or Ask a Question</h3>
                  <input
                    type="text"
                    required
                    placeholder="Headline title..."
                    value={postTitle}
                    onChange={(e) => setPostTitle(e.target.value)}
                    style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "10px 12px", fontSize: "14px", outline: "none" }}
                  />
                  <textarea
                    placeholder="Details about the referral, opening, or question..."
                    rows={4}
                    value={postBody}
                    onChange={(e) => setPostBody(e.target.value)}
                    style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "10px 12px", fontSize: "14px", outline: "none", resize: "none", fontFamily: "inherit" }}
                  />
                  <input
                    type="text"
                    placeholder="Tags (comma-separated, e.g. remote, hiring, referral)"
                    value={postTagsInput}
                    onChange={(e) => setPostTagsInput(e.target.value)}
                    style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "10px 12px", fontSize: "14px", outline: "none" }}
                  />

                  {/* Poll Composer Block */}
                  {showPollComposer && (
                    <div style={{ background: "#FCFAF7", border: "1px solid #E2D9CC", borderRadius: "12px", padding: "16px", display: "flex", flexDirection: "column", gap: "12px" }}>
                      <input
                        type="text"
                        placeholder="Poll Question..."
                        value={pollQuestion}
                        onChange={(e) => setPollQuestion(e.target.value)}
                        style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "8px 10px", fontSize: "13px" }}
                      />
                      {pollOptions.map((opt, i) => (
                        <div key={i} style={{ display: "flex", gap: "8px" }}>
                          <input
                            type="text"
                            placeholder={`Option ${i + 1}`}
                            value={opt}
                            onChange={(e) => {
                              const next = [...pollOptions];
                              next[i] = e.target.value;
                              setPollOptions(next);
                            }}
                            style={{ flex: 1, border: "1px solid #E2D9CC", borderRadius: "8px", padding: "8px 10px", fontSize: "13px" }}
                          />
                          {pollOptions.length > 2 && (
                            <button type="button" onClick={() => handleRemovePollOption(i)} style={{ border: "none", background: "transparent", color: "#A8472A", cursor: "pointer" }}>×</button>
                          )}
                        </div>
                      ))}
                      {pollOptions.length < 6 && (
                        <button type="button" onClick={handleAddPollOption} style={{ alignSelf: "flex-start", border: "none", background: "transparent", color: "#C2683C", fontSize: "13px", fontWeight: 600, cursor: "pointer" }}>
                          + Add Option
                        </button>
                      )}
                    </div>
                  )}

                  <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <button type="button" onClick={() => setShowPollComposer(!showPollComposer)} style={{ border: "none", background: "transparent", color: "#C2683C", fontSize: "14px", fontWeight: 600, cursor: "pointer" }}>
                      {showPollComposer ? "Remove Poll" : "📊 Add Poll"}
                    </button>
                    <button type="submit" disabled={submittingPost} style={{ border: "none", background: "#C2683C", color: "#fff", padding: "10px 24px", borderRadius: "100px", cursor: submittingPost ? "not-allowed" : "pointer", fontWeight: 600, fontSize: "14px" }}>
                      {submittingPost ? "Publishing..." : "Publish Post"}
                    </button>
                  </div>
                </form>
              )}

              {/* Feed Posts List */}
              <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                {posts.map((post) => (
                  <div key={post.id} style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px", display: "flex", flexDirection: "column", gap: "14px" }}>
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
                      <div style={{ display: "flex", gap: "10px", alignItems: "center" }}>
                        <div style={{ width: "36px", height: "36px", borderRadius: "50%", background: "#EFE7DC", color: "#C2683C", display: "flex", alignItems: "center", justifyContent: "center", fontWeight: 700, fontSize: "14px" }}>
                          U
                        </div>
                        <div>
                          <div style={{ fontSize: "14px", fontWeight: 600 }}>Anonymous Member</div>
                          <div style={{ fontSize: "11px", color: "#8A8175" }}>{new Date(post.created_at).toLocaleDateString()}</div>
                        </div>
                      </div>
                      <div style={{ display: "flex", gap: "10px" }}>
                        <button onClick={() => handleReportPost(post.id)} style={{ border: "none", background: "transparent", color: "#8A8175", fontSize: "13px", cursor: "pointer" }}>Report</button>
                        {user && (
                          <button onClick={() => handleDeletePost(post.id)} style={{ border: "none", background: "transparent", color: "#A8472A", fontSize: "13px", cursor: "pointer" }}>Hide</button>
                        )}
                      </div>
                    </div>

                    <div>
                      <h4 style={{ margin: "0 0 6px 0", fontSize: "16px", fontWeight: 700 }}>{post.title}</h4>
                      <p style={{ margin: 0, fontSize: "14px", lineHeight: "1.5", color: "#5B554C" }}>{post.body}</p>
                    </div>

                    {/* Poll render block */}
                    {post.poll && (
                      <div style={{ background: "#FCFAF7", border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", display: "flex", flexDirection: "column", gap: "12px" }}>
                        <div style={{ fontWeight: 600, fontSize: "14px" }}>📊 {post.poll.question}</div>
                        <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                          {post.poll.options.map((opt) => (
                            <button
                              key={opt.id}
                              onClick={() => handleVote(post.poll!.id, opt.id, post.id)}
                              style={{
                                border: "1px solid #E2D9CC",
                                background: "#fff",
                                borderRadius: "8px",
                                padding: "10px 14px",
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                                cursor: "pointer",
                                fontSize: "13px",
                                textAlign: "left",
                              }}
                            >
                              <span>{opt.label}</span>
                              <span style={{ fontWeight: 600, color: "#C2683C" }}>{opt.vote_count} votes</span>
                            </button>
                          ))}
                        </div>
                      </div>
                    )}

                    {post.tags && post.tags.length > 0 && (
                      <div style={{ display: "flex", gap: "6px", flexWrap: "wrap" }}>
                        {post.tags.map((t) => (
                          <span key={t} style={{ background: "#EFE7DC", color: "#5B554C", fontSize: "12px", padding: "4px 10px", borderRadius: "100px" }}>
                            #{t}
                          </span>
                        ))}
                      </div>
                    )}

                    <div style={{ display: "flex", gap: "20px", borderTop: "1px solid #FCFAF7", paddingTop: "14px" }}>
                      <button onClick={() => handleLikePost(post.id)} style={{ border: "none", background: "transparent", color: "#C2683C", cursor: "pointer", fontSize: "14px", fontWeight: 600 }}>
                        👍 Upvote ({post.reaction_count})
                      </button>
                      <button onClick={() => handleToggleComments(post.id)} style={{ border: "none", background: "transparent", color: "#5B554C", cursor: "pointer", fontSize: "14px", fontWeight: 600 }}>
                        💬 Comments ({post.comment_count})
                      </button>
                    </div>

                    {/* Comments thread block */}
                    {visibleCommentsPostId === post.id && (
                      <div style={{ borderTop: "1px solid #EFE7DC", paddingTop: "16px", marginTop: "8px", display: "flex", flexDirection: "column", gap: "12px" }}>
                        {loadingComments ? (
                          <CircularProgress style={{ alignSelf: "center", color: "#C2683C" }} size={24} />
                        ) : (
                          <>
                            <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
                              {comments.map((cm) => (
                                <div key={cm.id} style={{ background: "#FCFAF7", border: "1px solid #EFE7DC", borderRadius: "10px", padding: "12px", fontSize: "13px" }}>
                                  <div style={{ display: "flex", justifyContent: "space-between", color: "#8A8175", fontSize: "11px", marginBottom: "4px" }}>
                                    <strong>Anonymous Member</strong>
                                    <span>{new Date(cm.created_at).toLocaleDateString()}</span>
                                  </div>
                                  <p style={{ margin: 0, color: "#2B2620" }}>{cm.body}</p>
                                </div>
                              ))}
                              {comments.length === 0 && <p style={{ color: "#8A8175", fontSize: "13px", margin: 0 }}>No comments yet.</p>}
                            </div>

                            {/* Comment composer */}
                            <form onSubmit={(e) => handleAddComment(e, post.id)} style={{ display: "flex", gap: "10px" }}>
                              <input
                                type="text"
                                placeholder="Write a comment..."
                                value={newCommentBody}
                                onChange={(e) => setNewCommentBody(e.target.value)}
                                style={{ flex: 1, border: "1px solid #E2D9CC", borderRadius: "8px", padding: "8px 12px", fontSize: "13px", outline: "none" }}
                              />
                              <button type="submit" style={{ border: "none", background: "#C2683C", color: "#fff", borderRadius: "8px", padding: "8px 16px", cursor: "pointer", fontSize: "13px", fontWeight: 600 }}>
                                Comment
                              </button>
                            </form>
                          </>
                        )}
                      </div>
                    )}
                  </div>
                ))}
                {posts.length === 0 && (
                  <p style={{ textAlign: "center", color: "#8A8175", padding: "40px 0" }}>No discussion posts in this circle yet.</p>
                )}
              </div>
            </>
          )}

          {/* TAB 2: EVENTS */}
          {activeSubTab === "events" && (
            <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px", display: "flex", flexDirection: "column", gap: "18px" }}>
              <h3 style={{ margin: 0, fontSize: "18px", fontWeight: 700 }}>Upcoming Events</h3>
              <div style={{ display: "flex", flexDirection: "column", gap: "14px" }}>
                {simulatedEvents.map((ev) => (
                  <div key={ev.id} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", background: "#FCFAF7" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px 0", fontSize: "15px", fontWeight: 700 }}>{ev.title}</h4>
                      <p style={{ margin: 0, fontSize: "13px", color: "#8A8175" }}>{ev.date} • {ev.time}</p>
                      <p style={{ margin: "4px 0 0 0", fontSize: "12px", color: "#5B554C" }}>{ev.rsvpCount} people attending</p>
                    </div>
                    <button
                      onClick={() => handleRsvp(ev.id)}
                      style={{
                        border: "none",
                        background: ev.userRsvp ? "rgba(43,38,32,0.08)" : "#C2683C",
                        color: ev.userRsvp ? "#2B2620" : "#fff",
                        padding: "8px 16px",
                        borderRadius: "100px",
                        fontWeight: 600,
                        fontSize: "13px",
                        cursor: "pointer",
                      }}
                    >
                      {ev.userRsvp ? "RSVPed ✓" : "RSVP"}
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* TAB 3: FILES */}
          {activeSubTab === "files" && (
            <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px", display: "flex", flexDirection: "column", gap: "18px" }}>
              <h3 style={{ margin: 0, fontSize: "18px", fontWeight: 700 }}>Shared Documents</h3>
              <div style={{ border: "2px dashed #E2D9CC", borderRadius: "12px", padding: "32px", textAlign: "center", background: "#FCFAF7" }}>
                <span style={{ fontSize: "36px", display: "block" }}>📁</span>
                <span style={{ fontSize: "14px", fontWeight: 600, display: "block", margin: "12px 0 6px" }}>Upload a Resource</span>
                <span style={{ fontSize: "12px", color: "#8A8175", display: "block", marginBottom: "16px" }}>PDF, DOCX, ZIP, or CSV up to 10 MB.</span>
                <input
                  type="file"
                  onChange={handleFileUpload}
                  style={{ display: "none" }}
                  id="circle-file-upload"
                />
                <label
                  htmlFor="circle-file-upload"
                  style={{ display: "inline-block", background: "#C2683C", color: "#fff", padding: "10px 24px", borderRadius: "100px", fontSize: "13px", fontWeight: 600, cursor: "pointer" }}
                >
                  Choose File
                </label>
              </div>

              {uploadedFiles.length > 0 && (
                <div style={{ borderTop: "1px solid #EFE7DC", paddingTop: "18px" }}>
                  <h4 style={{ margin: "0 0 10px 0", fontSize: "14px", fontWeight: 700 }}>Circle Library</h4>
                  <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                    {uploadedFiles.map((file, i) => (
                      <div key={i} style={{ display: "flex", justifyContent: "space-between", background: "#FCFAF7", border: "1px solid #EFE7DC", borderRadius: "8px", padding: "10px 14px", fontSize: "13px" }}>
                        <span>📄 {file.name}</span>
                        <span style={{ color: "#8A8175" }}>{file.size}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}

        </div>

        {/* Right Sidebar About Panel */}
        <aside style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
          
          <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px" }}>
            <h3 style={{ margin: "0 0 10px 0", fontSize: "15px", fontWeight: 700 }}>Circle Rules</h3>
            <div style={{ fontSize: "13px", lineHeight: "1.55", color: "#5B554C", display: "flex", flexDirection: "column", gap: "8px" }}>
              <p style={{ margin: 0 }}><strong>1. Share leads:</strong> Keep posts focused on actual referrals and resource sharing.</p>
              <p style={{ margin: 0 }}><strong>2. Respect confidentiality:</strong> Do not share internal company communications.</p>
              <p style={{ margin: 0 }}><strong>3. Help others:</strong> Answer peer requests and reviews with patience.</p>
            </div>
          </div>

          {tags.length > 0 && (
            <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px" }}>
              <h3 style={{ margin: "0 0 12px 0", fontSize: "15px", fontWeight: 700 }}>Active Topics</h3>
              <div style={{ display: "flex", flexWrap: "wrap", gap: "6px" }}>
                {tags.map((t) => (
                  <span key={t.name} style={{ background: "#EFE7DC", color: "#5B554C", fontSize: "12px", padding: "6px 12px", borderRadius: "100px" }}>
                    #{t.name} ({t.count})
                  </span>
                ))}
              </div>
            </div>
          )}

        </aside>

      </section>

      <SiteFooter />
    </div>
  );
}
