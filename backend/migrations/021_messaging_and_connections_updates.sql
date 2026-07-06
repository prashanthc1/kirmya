-- Migration 021: Update connections and messaging schemas

-- 1. Update user_connections
ALTER TABLE user_connections ADD COLUMN responded_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE user_connections ADD COLUMN origin VARCHAR(50) DEFAULT 'manual_request';

-- 2. Update conversations
ALTER TABLE conversations ADD COLUMN type VARCHAR(20) NOT NULL DEFAULT 'direct';
ALTER TABLE conversations ADD COLUMN last_message_at TIMESTAMP WITH TIME ZONE;

-- Update type based on is_group
UPDATE conversations SET type = 'group' WHERE is_group = true;
ALTER TABLE conversations DROP COLUMN is_group;

-- 3. Update conversation_participants
ALTER TABLE conversation_participants ADD COLUMN left_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE conversation_participants ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'member';
ALTER TABLE conversation_participants ADD COLUMN is_archived BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE conversation_participants ADD COLUMN is_pinned BOOLEAN NOT NULL DEFAULT false;

-- 4. Update messages
ALTER TABLE messages ADD COLUMN content TEXT NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN content_type VARCHAR(20) NOT NULL DEFAULT 'text';
ALTER TABLE messages ADD COLUMN edited_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE messages ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

-- Copy existing body to content, then drop body
UPDATE messages SET content = body WHERE content = '';
ALTER TABLE messages DROP COLUMN body;

-- 5. Create message_statuses
CREATE TABLE IF NOT EXISTS message_statuses (
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'sent', -- sent, delivered, read
    status_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (message_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_message_statuses_user ON message_statuses(user_id);

-- 6. Update content_reports target_type check constraint
ALTER TABLE content_reports DROP CONSTRAINT IF EXISTS content_reports_target_type_check;
ALTER TABLE content_reports ADD CONSTRAINT content_reports_target_type_check CHECK (target_type IN ('post','comment','user','message','connection'));
