//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"workspace-app/internal/messaging/domain"
	"workspace-app/internal/testsupport"
)

func TestMessagingRepository_ConversationLifecycle(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	// Seed users for participants
	userA := testsupport.InsertUser(t, db, "usera@cb.test", "User A")
	userB := testsupport.InsertUser(t, db, "userb@cb.test", "User B")

	// 1. Create Conversation
	c := &domain.Conversation{
		Type:           "direct",
		Title:          "Chat between A and B",
		CreatedBy:      userA,
		ParticipantIDs: []string{userA, userB},
	}
	if err := repo.CreateConversation(ctx, c); err != nil {
		t.Fatalf("create conversation: %v", err)
	}
	if c.ID == "" {
		t.Fatal("expected conversation ID to be set")
	}

	// 2. Get Conversation
	got, err := repo.GetConversation(ctx, c.ID)
	if err != nil {
		t.Fatalf("get conversation: %v", err)
	}
	if got.Type != c.Type || got.Title != c.Title || got.CreatedBy != c.CreatedBy {
		t.Fatalf("conversation fields mismatch: expected %+v, got %+v", c, got)
	}
	if len(got.ParticipantIDs) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(got.ParticipantIDs))
	}

	// 3. Find Direct
	dirID, exists, err := repo.FindDirect(ctx, userA, userB)
	if err != nil {
		t.Fatalf("find direct: %v", err)
	}
	if !exists || dirID != c.ID {
		t.Fatalf("expected direct conversation to exist with id %s, got exists=%t id=%s", c.ID, exists, dirID)
	}

	// 4. IsParticipant
	isPart, err := repo.IsParticipant(ctx, c.ID, userA)
	if err != nil {
		t.Fatalf("is participant: %v", err)
	}
	if !isPart {
		t.Fatal("expected userA to be a participant")
	}

	// 5. Pin and Archive
	if err := repo.PinConversation(ctx, c.ID, userA, true); err != nil {
		t.Fatalf("pin conversation: %v", err)
	}
	if err := repo.ArchiveConversation(ctx, c.ID, userA, true); err != nil {
		t.Fatalf("archive conversation: %v", err)
	}

	pDetail, err := repo.GetParticipantDetail(ctx, c.ID, userA)
	if err != nil {
		t.Fatalf("get participant detail: %v", err)
	}
	if !pDetail.IsPinned || !pDetail.IsArchived {
		t.Fatalf("expected pinned and archived, got pinned=%t archived=%t", pDetail.IsPinned, pDetail.IsArchived)
	}

	// 6. List Conversations
	list, err := repo.ListConversations(ctx, userA)
	if err != nil {
		t.Fatalf("list conversations: %v", err)
	}
	if len(list) != 1 || list[0].ID != c.ID || !list[0].IsPinned || !list[0].IsArchived {
		t.Fatalf("unexpected listed conversations: %+v", list)
	}
}

func TestMessagingRepository_MessageLifecycle(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	userA := testsupport.InsertUser(t, db, "usera2@cb.test", "User A2")
	userB := testsupport.InsertUser(t, db, "userb2@cb.test", "User B2")

	c := &domain.Conversation{
		Type:           "direct",
		CreatedBy:      userA,
		ParticipantIDs: []string{userA, userB},
	}
	if err := repo.CreateConversation(ctx, c); err != nil {
		t.Fatalf("create conversation: %v", err)
	}

	// 1. Add Message
	m1 := &domain.Message{
		ConversationID: c.ID,
		SenderID:       userA,
		Content:        "Hello User B",
		ContentType:    "text",
	}
	if err := repo.AddMessage(ctx, m1); err != nil {
		t.Fatalf("add message: %v", err)
	}
	if m1.ID == "" {
		t.Fatal("expected message ID to be set")
	}

	// 2. Fetch Last Message
	last, err := repo.GetLastMessage(ctx, c.ID)
	if err != nil {
		t.Fatalf("get last message: %v", err)
	}
	if last == nil || last.ID != m1.ID || last.Content != m1.Content {
		t.Fatalf("unexpected last message: %+v", last)
	}

	// 3. Unread Count Check
	unread, err := repo.GetUnreadCount(ctx, c.ID, userB, nil)
	if err != nil {
		t.Fatalf("get unread count: %v", err)
	}
	if unread != 1 {
		t.Fatalf("expected 1 unread message for userB, got %d", unread)
	}

	// User B reads the conversation
	if err := repo.MarkRead(ctx, c.ID, userB); err != nil {
		t.Fatalf("mark read: %v", err)
	}

	// Retrieve conversation list for User B to get the updated LastReadAt time.
	listConv, err := repo.ListConversations(ctx, userB)
	if err != nil {
		t.Fatalf("list conversations for userB: %v", err)
	}
	if len(listConv) != 1 {
		t.Fatalf("expected 1 conversation for userB, got %d", len(listConv))
	}
	lastReadAt := listConv[0].LastReadAt
	if lastReadAt == nil {
		t.Fatal("expected LastReadAt to be set after mark read")
	}

	// Unread count relative to last read time
	unreadAfterRead, err := repo.GetUnreadCount(ctx, c.ID, userB, lastReadAt)
	if err != nil {
		t.Fatalf("get unread count after read: %v", err)
	}
	if unreadAfterRead != 0 {
		t.Fatalf("expected 0 unread messages after mark read, got %d", unreadAfterRead)
	}

	// 4. Message Statuses
	status := &domain.MessageStatus{
		MessageID:       m1.ID,
		UserID:          userB,
		Status:          "read",
		StatusUpdatedAt: time.Now(),
	}
	if err := repo.SetMessageStatus(ctx, status); err != nil {
		t.Fatalf("set message status: %v", err)
	}

	statuses, err := repo.GetMessageStatuses(ctx, m1.ID)
	if err != nil {
		t.Fatalf("get message statuses: %v", err)
	}
	if len(statuses) != 1 || statuses[0].Status != "read" {
		t.Fatalf("expected 1 read status, got %+v", statuses)
	}

	// 5. Update/Edit Message
	m1.Content = "Hello User B (edited)"
	now := time.Now()
	m1.EditedAt = &now
	if err := repo.UpdateMessage(ctx, m1); err != nil {
		t.Fatalf("update message: %v", err)
	}

	gotMsg, err := repo.GetMessage(ctx, m1.ID)
	if err != nil {
		t.Fatalf("get message: %v", err)
	}
	if gotMsg.Content != m1.Content || gotMsg.EditedAt == nil {
		t.Fatalf("message edit mismatch: %+v", gotMsg)
	}

	// 6. List Messages
	list, err := repo.ListMessages(ctx, c.ID, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(list) != 1 || list[0].ID != m1.ID {
		t.Fatalf("unexpected listed messages: %+v", list)
	}
}
