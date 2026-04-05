DROP TRIGGER IF EXISTS add_creator_after_chat_insert ON chats;
DROP TRIGGER IF EXISTS update_chat_last_message_trigger ON messages;
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;
DROP TRIGGER IF EXISTS update_chats_updated_at ON chats;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP FUNCTION IF EXISTS add_creator_to_participants() CASCADE;
DROP FUNCTION IF EXISTS update_chat_last_message() CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;