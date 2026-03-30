CREATE TABLE IF NOT EXISTS chats(
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    avatar_path VARCHAR(255),
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_chats_type ON chats(type);
CREATE INDEX idx_chats_created_by ON chats(created_by);
CREATE INDEX idx_chats_created_at ON chats(created_at);