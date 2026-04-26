import React, { useRef, useEffect, useState, useCallback } from "react";
import { useSelector, useDispatch } from "react-redux";
import {
  selectChatMessages,
  selectCurrentChat,
  setMessages,
  setLoading,
  updateCurrentChatLastMessage,
} from "../../../store/selectedChatSlice";
import { updateChatLastMessage } from "../../../store/chatSlice";
import { selectUser, selectToken } from "../../../store/userSlice";
import { BsFileEarmark, BsDownload } from "react-icons/bs";
import style from "./MessageList.module.scss";

interface Message {
  id: string;
  chat_id: string;
  user_id: string;
  type: "text" | "image" | "file";
  text: string;
  file_url?: string;
  file_name?: string;
  file_size?: number;
  mime_type?: string;
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
  const [lightboxImage, setLightboxImage] = useState<string | null>(null);

  const isFetchingRef = useRef(false);
  const lastChatIdRef = useRef<string | null>(null);

  const getToken = useCallback(() => {
    return token || localStorage.getItem("token");
  }, [token]);

  const fetchMessages = useCallback(async () => {
    if (!currentChat) return;
    if (isFetchingRef.current) return;

    // Если это тот же чат, не перезагружаем
    if (lastChatIdRef.current === currentChat.id && messages.length > 0) {
      return;
    }

    isFetchingRef.current = true;
    dispatch(setLoading(true));

    const authToken = getToken();

    if (!authToken) {
      console.error("Токен не найден");
      dispatch(setLoading(false));
      isFetchingRef.current = false;
      return;
    }

    try {
      const url = `http://localhost:8080/api/messages?chat_id=${currentChat.id}&limit=50&offset=0`;

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

        // Обновляем последнее сообщение
        if (sortedMessages.length > 0) {
          const lastMessage = sortedMessages[sortedMessages.length - 1];
          dispatch(updateCurrentChatLastMessage({ lastMessage }));
          dispatch(
            updateChatLastMessage({
              chatId: currentChat.id,
              lastMessage: lastMessage,
            }),
          );
        }

        lastChatIdRef.current = currentChat.id;
      } else {
        console.error("Ошибка загрузки сообщений:", response.status);
      }
    } catch (error) {
      console.error("Error fetching messages:", error);
    } finally {
      dispatch(setLoading(false));
      isFetchingRef.current = false;
    }
  }, [currentChat, getToken, dispatch, messages.length]);

  // Загружаем сообщения только при смене чата
  useEffect(() => {
    if (currentChat) {
      // Сбрасываем флаги при смене чата
      isFetchingRef.current = false;
      fetchMessages();
    }

    return () => {
      // Очищаем сообщения при закрытии чата
      dispatch(setMessages([]));
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentChat?.id]); // Только при изменении ID чата

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

  const formatFileSize = (bytes?: number): string => {
    if (!bytes) return "";
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
    return (bytes / (1024 * 1024)).toFixed(1) + " MB";
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

  const renderMessageContent = (message: Message) => {
    switch (message.type) {
      case "image":
        return (
          <div className={style.imageMessage}>
            <img
              src={`http://localhost:8080${message.file_url}`}
              alt={message.file_name || "Изображение"}
              onClick={() =>
                setLightboxImage(`http://localhost:8080${message.file_url}`)
              }
              loading="lazy"
            />
            {message.text && message.text !== message.file_name && (
              <p className={style.imageCaption}>{message.text}</p>
            )}
          </div>
        );

      case "file":
        return (
          <div className={style.fileMessage}>
            <a
              href={`http://localhost:8080${message.file_url}`}
              download={message.file_name}
              className={style.fileDownload}
              target="_blank"
              rel="noopener noreferrer"
            >
              <BsFileEarmark size={32} />
              <div className={style.fileDetails}>
                <span className={style.fileName}>{message.file_name}</span>
                <span className={style.fileSize}>
                  {formatFileSize(message.file_size)}
                </span>
              </div>
              <BsDownload size={20} className={style.downloadIcon} />
            </a>
            {message.text && message.text !== message.file_name && (
              <p className={style.fileCaption}>{message.text}</p>
            )}
          </div>
        );

      default:
        return <p className={style.messageText}>{message.text}</p>;
    }
  };

  // Проверяем загрузку
  const isLoading = useSelector((state: any) => state.selectedChat.isLoading);

  if (!currentChat) {
    return (
      <div className={style.emptyState}>
        <p>Выберите чат для начала общения</p>
      </div>
    );
  }

  if (isLoading && messages.length === 0) {
    return (
      <div className={style.emptyState}>
        <div className={style.spinner}></div>
        <p>Загрузка сообщений...</p>
      </div>
    );
  }

  const messageGroups = groupMessagesByDate();

  return (
    <>
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
                      {renderMessageContent(message)}
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

      {lightboxImage && (
        <div className={style.lightbox} onClick={() => setLightboxImage(null)}>
          <img src={lightboxImage} alt="Просмотр изображения" />
          <button
            className={style.closeLightbox}
            onClick={() => setLightboxImage(null)}
          >
            ✕
          </button>
        </div>
      )}
    </>
  );
}

export default MessageList;
