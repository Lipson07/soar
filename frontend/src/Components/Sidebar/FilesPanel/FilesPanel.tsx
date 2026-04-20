import React, { useState, useEffect } from "react";
import style from "./FilesPanel.module.scss";

interface FileItem {
  id: string;
  name: string;
  url: string;
  size: number;
  type: string;
  mime_type: string;
  created_at: string;
}

interface FilesPanelProps {
  isOpen: boolean;
  onClose: () => void;
}

const API_BASE = "http://localhost:8080/api";

function FilesPanel({ isOpen, onClose }: FilesPanelProps) {
  const [files, setFiles] = useState<FileItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<"all" | "images" | "documents">("all");
  const [searchQuery, setSearchQuery] = useState("");

  useEffect(() => {
    if (isOpen) {
      fetchFiles();
    }
  }, [isOpen]);

  const fetchFiles = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${API_BASE}/files`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const data = await response.json();
        setFiles(Array.isArray(data) ? data : []);
      } else {
        setFiles([]);
      }
    } catch (error) {
      console.error("Failed to fetch files:", error);
      setFiles([]);
    } finally {
      setLoading(false);
    }
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
    return (bytes / (1024 * 1024)).toFixed(1) + " MB";
  };

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));

    if (days === 0) return "Сегодня";
    if (days === 1) return "Вчера";
    if (days < 7) return `${days} дн назад`;
    return date.toLocaleDateString("ru-RU");
  };

  const getFileIcon = (mimeType: string) => {
    if (mimeType.startsWith("image/")) {
      return (
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
          <rect
            x="2"
            y="2"
            width="20"
            height="20"
            rx="2"
            stroke="currentColor"
            strokeWidth="2"
          />
          <circle cx="8.5" cy="8.5" r="1.5" fill="currentColor" />
          <path
            d="M21 15L16 10L5 21"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
      );
    }
    return (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
        <path
          d="M14 2H6C5.46957 2 4.96086 2.21071 4.58579 2.58579C4.21071 2.96086 4 3.46957 4 4V20C4 20.5304 4.21071 21.0391 4.58579 21.4142C4.96086 21.7893 5.46957 22 6 22H18C18.5304 22 19.0391 21.7893 19.4142 21.4142C19.7893 21.0391 20 20.5304 20 20V8L14 2Z"
          stroke="currentColor"
          strokeWidth="2"
        />
        <path d="M14 2V8H20" stroke="currentColor" strokeWidth="2" />
      </svg>
    );
  };

  const filteredFiles = files.filter((file) => {
    const matchesFilter =
      filter === "all" ||
      (filter === "images" && file.type === "image") ||
      (filter === "documents" && file.type === "document");

    const matchesSearch = file.name
      .toLowerCase()
      .includes(searchQuery.toLowerCase());

    return matchesFilter && matchesSearch;
  });

  const handleDownload = (file: FileItem) => {
    window.open(`http://localhost:8080${file.url}`, "_blank");
  };

  const handleDelete = async (file: FileItem) => {
    if (!confirm(`Удалить файл "${file.name}"?`)) return;

    try {
      const token = localStorage.getItem("token");
      const urlPath = file.url.replace("/uploads/", "");
      const response = await fetch(`${API_BASE}/files/${urlPath}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        setFiles((prev) => prev.filter((f) => f.id !== file.id));
      }
    } catch (error) {
      console.error("Failed to delete file:", error);
    }
  };

  if (!isOpen) return null;

  return (
    <>
      <div className={style.overlay} onClick={onClose} />
      <div className={style.panel}>
        <div className={style.header}>
          <h2>Файлы</h2>
          <button className={style.closeButton} onClick={onClose}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path
                d="M18 6L6 18M6 6L18 18"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
          </button>
        </div>

        <div className={style.content}>
          <div className={style.searchWrapper}>
            <input
              type="text"
              placeholder="Поиск файлов..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className={style.searchInput}
            />
            <svg
              className={style.searchIcon}
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
            >
              <path
                d="M15.5 15.5L19 19M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
          </div>

          <div className={style.filters}>
            <button
              className={`${style.filterBtn} ${filter === "all" ? style.active : ""}`}
              onClick={() => setFilter("all")}
            >
              Все
            </button>
            <button
              className={`${style.filterBtn} ${filter === "images" ? style.active : ""}`}
              onClick={() => setFilter("images")}
            >
              Изображения
            </button>
            <button
              className={`${style.filterBtn} ${filter === "documents" ? style.active : ""}`}
              onClick={() => setFilter("documents")}
            >
              Документы
            </button>
          </div>

          <div className={style.filesList}>
            {loading ? (
              <div className={style.loadingState}>
                <div className={style.spinner}></div>
                <span>Загрузка файлов...</span>
              </div>
            ) : filteredFiles.length === 0 ? (
              <div className={style.emptyState}>
                <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M14 2H6C5.46957 2 4.96086 2.21071 4.58579 2.58579C4.21071 2.96086 4 3.46957 4 4V20C4 20.5304 4.21071 21.0391 4.58579 21.4142C4.96086 21.7893 5.46957 22 6 22H18C18.5304 22 19.0391 21.7893 19.4142 21.4142C19.7893 21.0391 20 20.5304 20 20V8L14 2Z"
                    stroke="currentColor"
                    strokeWidth="2"
                  />
                  <path d="M14 2V8H20" stroke="currentColor" strokeWidth="2" />
                </svg>
                <p>{searchQuery ? "Ничего не найдено" : "Нет файлов"}</p>
              </div>
            ) : (
              filteredFiles.map((file) => (
                <div key={file.id} className={style.fileItem}>
                  <div className={style.fileIcon}>
                    {getFileIcon(file.mime_type)}
                  </div>
                  <div className={style.fileInfo}>
                    <div className={style.fileName}>{file.name}</div>
                    <div className={style.fileMeta}>
                      {formatFileSize(file.size)} •{" "}
                      {formatDate(file.created_at)}
                    </div>
                  </div>
                  <div className={style.fileActions}>
                    <button
                      className={style.actionBtn}
                      onClick={() => handleDownload(file)}
                      title="Скачать"
                    >
                      <svg
                        width="18"
                        height="18"
                        viewBox="0 0 24 24"
                        fill="none"
                      >
                        <path
                          d="M12 3V16M12 16L8 12M12 16L16 12M5 21H19"
                          stroke="currentColor"
                          strokeWidth="2"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                        />
                      </svg>
                    </button>
                    <button
                      className={style.actionBtn}
                      onClick={() => handleDelete(file)}
                      title="Удалить"
                    >
                      <svg
                        width="18"
                        height="18"
                        viewBox="0 0 24 24"
                        fill="none"
                      >
                        <path
                          d="M4 7H20M10 11V17M14 11V17M5 7L6 19C6 19.5304 6.21071 20.0391 6.58579 20.4142C6.96086 20.7893 7.46957 21 8 21H16C16.5304 21 17.0391 20.7893 17.4142 20.4142C17.7893 20.0391 18 19.5304 18 19L19 7M9 7V4C9 3.73478 9.10536 3.48043 9.29289 3.29289C9.48043 3.10536 9.73478 3 10 3H14C14.2652 3 14.5196 3.10536 14.7071 3.29289C14.8946 3.48043 15 3.73478 15 4V7"
                          stroke="currentColor"
                          strokeWidth="2"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                        />
                      </svg>
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </>
  );
}

export default FilesPanel;
