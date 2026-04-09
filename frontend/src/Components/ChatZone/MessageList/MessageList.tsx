import React, { useRef, useEffect, useState, useCallback } from "react";
import { useSelector, useDispatch } from "react-redux";
import {
  selectChatMessages,
  selectCurrentChat,
  setMessages,
  setLoading,
} from "../../../store/selectedChatSlice";
import { selectUser, selectToken } from "../../../store/userSlice";
import style from "./MessageList.module.scss";

interface Message {
  id: string;
  chat_id: string;
  user_id: string;
  text: string;
  reply_to: string | null;
  is_edited: boolean;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}

function MessageList() {
  const messages = useSelector(selectChatMessages);
  const currentChat = useSelector(selectCurrentChat);
  const currentUser = useSelector(selectUser);
  const token = useSelector(selectToken);
  const dispatch = useDispatch();
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);

  const isFetchingRef = useRef(false);

  const getToken = useCallback(() => {
    return token || localStorage.getItem("token");
  }, [token]);

  const fetchMessages = useCallback(async () => {
    if (!currentChat || isFetchingRef.current) return;

    isFetchingRef.current = true;
    setIsLoadingMessages(true);
    dispatch(setLoading(true));

    const authToken = getToken();

    if (!authToken) {
      console.error("No token available");
      setIsLoadingMessages(false);
      dispatch(setLoading(false));
      isFetchingRef.current = false;
      return;
    }

    try {
      const url = `http://localhost:8080/api/messages?chat_id=${currentChat.id}&limit=50&offset=0`;
      console.log("Fetching messages for chat:", currentChat.id);

      const response = await fetch(url, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
      });

      if (response.status === 204) {
        dispatch(setMessages([]));
        return;
      }

      if (response.ok) {
        const data = await response.json();
        const messagesArray = Array.isArray(data) ? data : data.messages || [];
        const sortedMessages = messagesArray.sort(
          (a: Message, b: Message) =>
            new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
        );
        dispatch(setMessages(sortedMessages));
      } else if (response.status === 401) {
        console.error("Unauthorized - please login again");
      } else {
        console.error("Failed to fetch messages:", response.status);
      }
    } catch (error) {
      console.error("Error fetching messages:", error);
    } finally {
      setIsLoadingMessages(false);
      dispatch(setLoading(false));
      isFetchingRef.current = false;
    }
  }, [currentChat, getToken, dispatch]);

  useEffect(() => {
    if (currentChat) {
      fetchMessages();
    }

    return () => {
      dispatch(setMessages([]));
    };
  }, [currentChat?.id, fetchMessages, dispatch]);

  useEffect(() => {
    if (!currentChat) return;

    const interval = setInterval(() => {
      fetchMessages();
    }, 5000);

    return () => clearInterval(interval);
  }, [currentChat?.id, fetchMessages]);

  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  };

  const formatDate = (timestamp: string) => {
    const date = new Date(timestamp);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);

    if (date.toDateString() === today.toDateString()) {
      return "Сегодня";
    } else if (date.toDateString() === yesterday.toDateString()) {
      return "Вчера";
    } else {
      return date.toLocaleDateString();
    }
  };

  const groupMessagesByDate = () => {
    const groups: { [key: string]: Message[] } = {};

    messages.forEach((message) => {
      const date = formatDate(message.created_at);
      if (!groups[date]) {
        groups[date] = [];
      }
      groups[date].push(message);
    });

    return groups;
  };

  if (!currentChat) {
    return (
      <div className={style.emptyState}>
        <p>Выберите чат для начала общения</p>
      </div>
    );
  }

  if (isLoadingMessages && messages.length === 0) {
    return (
      <div className={style.emptyState}>
        <div className={style.spinner}></div>
        <p>Загрузка сообщений...</p>
      </div>
    );
  }

  const messageGroups = groupMessagesByDate();

  return (
    <div className={style.messageList}>
      <div className={style.messagesContainer}>
        {Object.entries(messageGroups).map(([date, dateMessages]) => (
          <div key={date} className={style.dateGroup}>
            <div className={style.dateDivider}>
              <span>{date}</span>
            </div>
            {dateMessages.map((message) => {
              const isOwn = message.user_id === currentUser?.id;

              return (
                <div
                  key={message.id}
                  className={`${style.message} ${isOwn ? style.own : style.other}`}
                >
                  <div className={style.messageBubble}>
                    <p className={style.messageText}>{message.text}</p>
                    <span className={style.timestamp}>
                      {formatTime(message.created_at)}
                      {message.is_edited && " (ред.)"}
                    </span>
                  </div>
                </div>
              );
            })}
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>
    </div>
  );
}

export default MessageList;
