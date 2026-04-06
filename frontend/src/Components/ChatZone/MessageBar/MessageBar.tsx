import React, { useState } from "react";
import { IoSendOutline } from "react-icons/io5";
import { IoMdAttach } from "react-icons/io";
import { BsEmojiSmile } from "react-icons/bs";
import { useDispatch, useSelector } from "react-redux";
import {
  addMessage,
  selectCurrentChat,
  setMessages,
} from "../../../store/selectedChatSlice";
import { selectUser, selectToken } from "../../../store/userSlice";
import style from "./MessageBar.module.scss";

function MessageBar() {
  const [message, setMessage] = useState("");
  const [isSending, setIsSending] = useState(false);
  const dispatch = useDispatch();
  const currentChat = useSelector(selectCurrentChat);
  const currentUser = useSelector(selectUser);
  const token = useSelector(selectToken);

  const getToken = () => {
    return token || localStorage.getItem("token");
  };

  const sendMessageToServer = async (text: string) => {
    if (!currentChat) return null;

    const authToken = getToken();

    if (!authToken) {
      throw new Error("Токен авторизации отсутствует");
    }

    const response = await fetch(
      `http://localhost:8080/api/messages?chat_id=${currentChat.id}`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ text: text }),
      },
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "Failed to send message");
    }

    return await response.json();
  };

  const fetchMessages = async () => {
    if (!currentChat) return;

    const authToken = getToken();

    if (!authToken) return;

    const response = await fetch(
      `http://localhost:8080/api/messages?chat_id=${currentChat.id}&limit=50&offset=0`,
      {
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
      },
    );

    if (response.ok) {
      const data = await response.json();
      const messagesArray = Array.isArray(data) ? data : data.messages || [];
      const sortedMessages = messagesArray.sort(
        (a: any, b: any) =>
          new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
      );
      dispatch(setMessages(sortedMessages));
    }
  };

  const handleSend = async () => {
    if (!message.trim() || !currentChat || isSending) return;

    setIsSending(true);
    const messageText = message.trim();

    const tempMessage = {
      id: `temp-${Date.now()}`,
      chat_id: currentChat.id,
      user_id: currentUser?.id || "",
      text: messageText,
      reply_to: null,
      is_edited: false,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      deleted_at: null,
    };

    dispatch(addMessage(tempMessage));
    setMessage("");

    try {
      await sendMessageToServer(messageText);
      await fetchMessages();
    } catch (error) {
      console.error("Failed to send message:", error);
      setMessage(messageText);
    } finally {
      setIsSending(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey && !isSending) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className={style.messageBar}>
      <div className={style.attachments}>
        <button className={style.attachButton} title="Прикрепить файл">
          <IoMdAttach size={20} />
        </button>
        <button className={style.attachButton} title="Эмодзи">
          <BsEmojiSmile size={20} />
        </button>
      </div>

      <textarea
        className={style.messageInput}
        placeholder="Введите сообщение..."
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        onKeyPress={handleKeyPress}
        rows={1}
        disabled={isSending}
      />

      <button
        className={style.sendButton}
        onClick={handleSend}
        disabled={!message.trim() || isSending}
        title="Отправить"
      >
        <IoSendOutline size={20} />
      </button>
    </div>
  );
}

export default MessageBar;
