
ALTER TABLE chats 
    ADD CONSTRAINT fk_chats_creator 
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE participants 
    ADD CONSTRAINT fk_participants_chat 
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_participants_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE messages 
    ADD CONSTRAINT fk_messages_chat 
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_messages_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_messages_reply 
    FOREIGN KEY (reply_to) REFERENCES messages(id) ON DELETE SET NULL;