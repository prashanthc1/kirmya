-- Migration 026: Create connections system

-- Drop legacy table user_connections cascade
DROP TABLE IF EXISTS user_connections CASCADE;

-- 1. connections
CREATE TABLE IF NOT EXISTS connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_a_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_b_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('pending','accepted','declined','blocked')),
    requested_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    responded_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT uq_pair UNIQUE (user_a_id, user_b_id),
    CONSTRAINT chk_order CHECK (user_a_id < user_b_id)
);

CREATE INDEX IF NOT EXISTS idx_connections_user_a ON connections(user_a_id);
CREATE INDEX IF NOT EXISTS idx_connections_user_b ON connections(user_b_id);
CREATE INDEX IF NOT EXISTS idx_connections_status ON connections(status);

-- 2. connection_requests_meta
CREATE TABLE IF NOT EXISTS connection_requests_meta (
    connection_id UUID PRIMARY KEY REFERENCES connections(id) ON DELETE CASCADE,
    note VARCHAR(300),
    source TEXT CHECK (source IN ('search','profile_view','suggested','import'))
);

-- 3. blocks
CREATE TABLE IF NOT EXISTS blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blocker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT uq_blocker_blocked UNIQUE (blocker_id, blocked_id)
);

CREATE INDEX IF NOT EXISTS idx_blocks_blocker ON blocks(blocker_id);
CREATE INDEX IF NOT EXISTS idx_blocks_blocked ON blocks(blocked_id);

-- 4. connection_counts
CREATE TABLE IF NOT EXISTS connection_counts (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    connection_count INT DEFAULT 0,
    pending_incoming_count INT DEFAULT 0,
    pending_outgoing_count INT DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT now()
);
