ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_messages_reply;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_messages_user;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_messages_chat;
ALTER TABLE participants DROP CONSTRAINT IF EXISTS fk_participants_user;
ALTER TABLE participants DROP CONSTRAINT IF EXISTS fk_participants_chat;
ALTER TABLE chats DROP CONSTRAINT IF EXISTS fk_chats_creator;