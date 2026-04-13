import React, { useState, useRef } from "react";
import { IoSendOutline, IoClose } from "react-icons/io5";
import { IoMdAttach } from "react-icons/io";
import { BsEmojiSmile, BsFileEarmark, BsImage } from "react-icons/bs";
import { useDispatch, useSelector } from "react-redux";
import {
  addMessage,
  selectCurrentChat,
  setMessages,
} from "../../../store/selectedChatSlice";
import { selectUser, selectToken } from "../../../store/userSlice";
import style from "./MessageBar.module.scss";

interface UploadedFile {
  file_url: string;
  file_name: string;
  file_size: number;
  mime_type: string;
}

function MessageBar() {
  const [message, setMessage] = useState("");
  const [isSending, setIsSending] = useState(false);
  const [uploadedFile, setUploadedFile] = useState<UploadedFile | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const imageInputRef = useRef<HTMLInputElement>(null);

  const dispatch = useDispatch();
  const currentChat = useSelector(selectCurrentChat);
  const currentUser = useSelector(selectUser);
  const token = useSelector(selectToken);

  const getToken = () => {
    return token || localStorage.getItem("token");
  };

  const uploadFile = async (file: File): Promise<UploadedFile | null> => {
    const authToken = getToken();
    if (!authToken) return null;

    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch(
        "http://localhost:8080/api/messages/upload",
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
          body: formData,
        },
      );

      if (!response.ok) {
        throw new Error("Ошибка загрузки файла");
      }

      return await response.json();
    } catch (error) {
      console.error("Upload error:", error);
      return null;
    }
  };

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setIsUploading(true);
    const uploaded = await uploadFile(file);

    if (uploaded) {
      setUploadedFile(uploaded);
    }

    setIsUploading(false);
    if (fileInputRef.current) fileInputRef.current.value = "";
    if (imageInputRef.current) imageInputRef.current.value = "";
  };

  const clearUploadedFile = () => {
    setUploadedFile(null);
  };

  const sendMessageToServer = async (text: string, fileData?: UploadedFile) => {
    if (!currentChat) return null;

    const authToken = getToken();
    if (!authToken) {
      throw new Error("Токен авторизации отсутствует");
    }

    let messageType = "text";
    let body: any = { type: "text", text };

    if (fileData) {
      const isImage = fileData.mime_type.startsWith("image/");
      messageType = isImage ? "image" : "file";
      body = {
        type: messageType,
        text: text || fileData.file_name,
        file_url: fileData.file_url,
        file_name: fileData.file_name,
        file_size: fileData.file_size,
        mime_type: fileData.mime_type,
      };
    }

    const response = await fetch(
      `http://localhost:8080/api/messages?chat_id=${currentChat.id}`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
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
    if ((!message.trim() && !uploadedFile) || !currentChat || isSending) return;

    setIsSending(true);
    const messageText = message.trim();

    const tempMessage: any = {
      id: `temp-${Date.now()}`,
      chat_id: currentChat.id,
      user_id: currentUser?.id || "",
      type: uploadedFile
        ? uploadedFile.mime_type.startsWith("image/")
          ? "image"
          : "file"
        : "text",
      text: messageText || uploadedFile?.file_name || "",
      file_url: uploadedFile?.file_url,
      file_name: uploadedFile?.file_name,
      file_size: uploadedFile?.file_size,
      mime_type: uploadedFile?.mime_type,
      reply_to: null,
      is_edited: false,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      deleted_at: null,
    };

    dispatch(addMessage(tempMessage));
    setMessage("");
    setUploadedFile(null);

    try {
      await sendMessageToServer(messageText, uploadedFile || undefined);
      await fetchMessages();
    } catch (error) {
      console.error("Failed to send message:", error);
      setMessage(messageText);
      if (uploadedFile) setUploadedFile(uploadedFile);
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

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
    return (bytes / (1024 * 1024)).toFixed(1) + " MB";
  };

  return (
    <div className={style.messageBar}>
      {uploadedFile && (
        <div className={style.filePreview}>
          {uploadedFile.mime_type.startsWith("image/") ? (
            <div className={style.imagePreview}>
              <img
                src={`http://localhost:8080${uploadedFile.file_url}`}
                alt={uploadedFile.file_name}
              />
              <button className={style.removeFile} onClick={clearUploadedFile}>
                <IoClose size={16} />
              </button>
            </div>
          ) : (
            <div className={style.fileAttachment}>
              <BsFileEarmark size={24} />
              <div className={style.fileInfo}>
                <span className={style.fileName}>{uploadedFile.file_name}</span>
                <span className={style.fileSize}>
                  {formatFileSize(uploadedFile.file_size)}
                </span>
              </div>
              <button className={style.removeFile} onClick={clearUploadedFile}>
                <IoClose size={16} />
              </button>
            </div>
          )}
        </div>
      )}

      <div className={style.attachments}>
        <input
          ref={imageInputRef}
          type="file"
          accept="image/*"
          style={{ display: "none" }}
          onChange={handleFileSelect}
        />
        <input
          ref={fileInputRef}
          type="file"
          accept="*/*"
          style={{ display: "none" }}
          onChange={handleFileSelect}
        />

        <button
          className={style.attachButton}
          title="Прикрепить изображение"
          onClick={() => imageInputRef.current?.click()}
          disabled={isUploading}
        >
          <BsImage size={20} />
        </button>
        <button
          className={style.attachButton}
          title="Прикрепить файл"
          onClick={() => fileInputRef.current?.click()}
          disabled={isUploading}
        >
          <IoMdAttach size={20} />
        </button>
        <button className={style.attachButton} title="Эмодзи">
          <BsEmojiSmile size={20} />
        </button>
      </div>

      <textarea
        className={style.messageInput}
        placeholder={
          uploadedFile
            ? "Добавить подпись (необязательно)"
            : "Введите сообщение..."
        }
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        onKeyPress={handleKeyPress}
        rows={1}
        disabled={isSending || isUploading}
      />

      <button
        className={style.sendButton}
        onClick={handleSend}
        disabled={
          (!message.trim() && !uploadedFile) || isSending || isUploading
        }
        title="Отправить"
      >
        <IoSendOutline size={20} />
      </button>
    </div>
  );
}

export default MessageBar;
