import React, { useRef, useEffect, useState, useCallback } from "react";
import { useSelector, useDispatch } from "react-redux";
import {
  selectChatMessages,
  selectCurrentChat,
  setMessages,
  setLoading,
  updateCurrentChatLastMessage,
  updateMessage,
  updateCurrentChat,
} from "../../../store/selectedChatSlice";
import { updateChatLastMessage, addChat } from "../../../store/chatSlice";
import { selectUser, selectToken } from "../../../store/userSlice";
import {
  BsFileEarmark,
  BsDownload,
  BsThreeDotsVertical,
  BsPencil,
  BsTrash,
} from "react-icons/bs";
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
  const [editingMessage, setEditingMessage] = useState<Message | null>(null);
  const [editText, setEditText] = useState("");
  const [menuMessageId, setMenuMessageId] = useState<string | null>(null);
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);
  const hasLoadedRef = useRef(false);
  const editInputRef = useRef<HTMLTextAreaElement>(null);

  const getToken = useCallback(() => {
    return token || localStorage.getItem("token");
  }, [token]);

  useEffect(() => {
    if (editingMessage && editInputRef.current) {
      editInputRef.current.focus();
      editInputRef.current.setSelectionRange(editText.length, editText.length);
    }
  }, [editingMessage]);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (menuMessageId) {
        const target = e.target as HTMLElement;
        if (!target.closest(`.${style.messageMenu}`)) {
          setMenuMessageId(null);
        }
      }
    };

    document.addEventListener("click", handleClickOutside);
    return () => document.removeEventListener("click", handleClickOutside);
  }, [menuMessageId]);

  const handleEditMessage = async (message: Message) => {
    if (!editText.trim() || editText === message.text) {
      setEditingMessage(null);
      setEditText("");
      return;
    }

    const authToken = getToken();
    if (!authToken) return;

    try {
      const response = await fetch(
        `http://localhost:8080/api/messages?message_id=${message.id}`,
        {
          method: "PUT",
          headers: {
            Authorization: `Bearer ${authToken}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ text: editText }),
        },
      );

      if (response.ok) {
        const updatedMessage = {
          ...message,
          text: editText,
          is_edited: true,
          updated_at: new Date().toISOString(),
        };
        dispatch(updateMessage(updatedMessage));

        if (messages[messages.length - 1]?.id === message.id && currentChat) {
          dispatch(
            updateCurrentChatLastMessage({ lastMessage: updatedMessage }),
          );
          dispatch(
            updateChatLastMessage({
              chatId: currentChat.id,
              lastMessage: updatedMessage,
            }),
          );
        }

        setEditingMessage(null);
        setEditText("");
      } else {
        console.error("Failed to edit message:", await response.text());
      }
    } catch (error) {
      console.error("Failed to edit message:", error);
    }
    setMenuMessageId(null);
  };

  const handleCancelEdit = () => {
    setEditingMessage(null);
    setEditText("");
  };

  const handleEditKeyDown = (e: React.KeyboardEvent, message: Message) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleEditMessage(message);
    } else if (e.key === "Escape") {
      handleCancelEdit();
    }
  };

  const handleDeleteMessage = async (messageId: string) => {
    if (!window.confirm("Удалить сообщение?")) return;

    const authToken = getToken();
    if (!authToken) return;

    try {
      const response = await fetch(
        `http://localhost:8080/api/messages?message_id=${messageId}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        },
      );

      if (response.ok) {
        const updatedMessages = messages.filter((msg) => msg.id !== messageId);
        dispatch(setMessages(updatedMessages));

        if (currentChat && updatedMessages.length > 0) {
          const lastMessage = updatedMessages[updatedMessages.length - 1];
          dispatch(updateCurrentChatLastMessage({ lastMessage }));
          dispatch(
            updateChatLastMessage({
              chatId: currentChat.id,
              lastMessage: lastMessage,
            }),
          );
        }

        setTimeout(() => loadMessages(), 500);
      }
    } catch (error) {
      console.error("Failed to delete message:", error);
    }
    setMenuMessageId(null);
  };

  const loadMessages = useCallback(async () => {
    if (!currentChat) return;
    if (currentChat.id.startsWith("temp_")) return;
    if (isLoadingMessages) return;

    setIsLoadingMessages(true);
    dispatch(setLoading(true));

    const authToken = getToken();
    if (!authToken) {
      setIsLoadingMessages(false);
      dispatch(setLoading(false));
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

        if (sortedMessages.length > 0 && currentChat) {
          const lastMessage = sortedMessages[sortedMessages.length - 1];
          dispatch(updateCurrentChatLastMessage({ lastMessage }));
          dispatch(
            updateChatLastMessage({
              chatId: currentChat.id,
              lastMessage: lastMessage,
            }),
          );
        }
      }
    } catch (error) {
      console.error("Error fetching messages:", error);
    } finally {
      setIsLoadingMessages(false);
      dispatch(setLoading(false));
    }
  }, [currentChat, getToken, dispatch]);

  useEffect(() => {
    if (
      currentChat &&
      !currentChat.id.startsWith("temp_") &&
      !hasLoadedRef.current
    ) {
      hasLoadedRef.current = true;
      loadMessages();
    }

    return () => {};
  }, [currentChat?.id, loadMessages]);

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
      if (!groups[date]) groups[date] = [];
      groups[date].push(message);
    });
    return groups;
  };

  const renderMessageContent = (message: Message) => {
    if (editingMessage?.id === message.id) {
      return (
        <div className={style.editMode}>
          <textarea
            ref={editInputRef}
            value={editText}
            onChange={(e) => setEditText(e.target.value)}
            onKeyDown={(e) => handleEditKeyDown(e, message)}
            className={style.editInput}
            rows={3}
          />
          <div className={style.editActions}>
            <button onClick={handleCancelEdit}>Отмена</button>
            <button onClick={() => handleEditMessage(message)}>
              Сохранить
            </button>
          </div>
        </div>
      );
    }

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
          {currentChat.id.startsWith("temp_") && (
            <div className={style.dateDivider}>
              <span>Новый чат</span>
            </div>
          )}
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
                    id={`message-${message.id}`}
                    className={`${style.message} ${
                      isOwn ? style.own : style.other
                    }`}
                  >
                    <div className={style.messageBubble}>
                      {renderMessageContent(message)}
                      <div
                        style={{
                          display: "flex",
                          justifyContent: "space-between",
                          alignItems: "center",
                          marginTop: "4px",
                        }}
                      >
                        <span className={style.timestamp}>
                          {formatTime(message.created_at)}
                          {message.is_edited && (
                            <span style={{ fontStyle: "italic", opacity: 0.7 }}>
                              {" "}
                              (ред.)
                            </span>
                          )}
                        </span>
                        {isOwn && message.type === "text" && (
                          <div className={style.messageMenu}>
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                setMenuMessageId(
                                  menuMessageId === message.id
                                    ? null
                                    : message.id,
                                );
                              }}
                              className={style.menuButton}
                            >
                              <BsThreeDotsVertical size={14} />
                            </button>
                            {menuMessageId === message.id && (
                              <div className={style.contextMenu}>
                                <button
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    setEditingMessage(message);
                                    setEditText(message.text);
                                    setMenuMessageId(null);
                                  }}
                                >
                                  <BsPencil size={14} /> Редактировать
                                </button>
                                <button
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handleDeleteMessage(message.id);
                                  }}
                                  className={style.deleteBtn}
                                >
                                  <BsTrash size={14} /> Удалить
                                </button>
                              </div>
                            )}
                          </div>
                        )}
                      </div>
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
