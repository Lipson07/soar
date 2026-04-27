import React, { useState, useRef, useEffect, useCallback } from "react";
import { useSelector } from "react-redux";
import { selectUser, selectToken } from "../../../store/userSlice";
import style from "./StoriesSection.module.scss";

interface Story {
  id: string;
  user_id: string;
  user_name: string;
  user_avatar: string | null;
  file_url: string;
  type: "image" | "video";
  created_at: string;
  expires_at: string;
  viewed: boolean;
}

function StoriesSection() {
  const [stories, setStories] = useState<Story[]>([]);
  const [selectedStory, setSelectedStory] = useState<Story | null>(null);
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const videoRef = useRef<HTMLVideoElement>(null);
  const progressIntervalRef = useRef<NodeJS.Timeout>();

  const currentUser = useSelector(selectUser);
  const token = useSelector(selectToken);

  const API_URL = "http://localhost:8080";

  const fetchStories = useCallback(async () => {
    try {
      if (!token) return;

      const response = await fetch(`${API_URL}/api/stories`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        console.log("Fetched stories:", data);
        setStories(data);
      }
    } catch (error) {
      console.error("Failed to fetch stories:", error);
    }
  }, [token]);

  useEffect(() => {
    fetchStories();
    const interval = setInterval(fetchStories, 30000);
    return () => clearInterval(interval);
  }, [fetchStories]);

  const uploadStory = async (file: File) => {
    if (!token || !currentUser) return;

    if (file.size > 50 * 1024 * 1024) {
      alert("Файл слишком большой. Максимальный размер 50MB");
      return;
    }

    const allowedTypes = [
      "image/jpeg",
      "image/png",
      "image/gif",
      "image/webp",
      "video/mp4",
      "video/webm",
    ];
    if (!allowedTypes.includes(file.type)) {
      alert("Неподдерживаемый формат файла");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    setUploading(true);
    try {
      const response = await fetch(`${API_URL}/api/stories/upload`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });

      if (response.ok) {
        await fetchStories();
      } else {
        const errorData = await response.json();
        alert(`Ошибка загрузки: ${errorData.error}`);
      }
    } catch (error) {
      console.error("Failed to upload story:", error);
      alert("Ошибка при загрузке сторис");
    } finally {
      setUploading(false);
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      uploadStory(file);
    }
  };

  const markAsViewed = async (storyId: string) => {
    try {
      if (!token) return;

      await fetch(`${API_URL}/api/stories/${storyId}/view`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    } catch (error) {
      console.error("Failed to mark story as viewed:", error);
    }
  };

  const handleStoryClick = (story: Story) => {
    setSelectedStory(story);
    markAsViewed(story.id);

    setStories((prev) =>
      prev.map((s) => (s.id === story.id ? { ...s, viewed: true } : s)),
    );

    if (story.type === "image") {
      setProgress(0);
      if (progressIntervalRef.current)
        clearInterval(progressIntervalRef.current);

      const duration = 5000;
      const interval = 50;
      const steps = duration / interval;
      let currentStep = 0;

      progressIntervalRef.current = setInterval(() => {
        currentStep++;
        const newProgress = (currentStep / steps) * 100;

        if (newProgress >= 100) {
          clearInterval(progressIntervalRef.current);
          setSelectedStory(null);
          setProgress(0);
        } else {
          setProgress(newProgress);
        }
      }, interval);
    }
  };

  const closeStory = () => {
    if (progressIntervalRef.current) {
      clearInterval(progressIntervalRef.current);
    }
    if (videoRef.current) {
      videoRef.current.pause();
    }
    setSelectedStory(null);
    setProgress(0);
  };

  const handleVideoEnded = () => {
    closeStory();
  };

  const formatTimeAgo = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diff = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (diff < 60) return `${diff} сек`;
    if (diff < 3600) return `${Math.floor(diff / 60)} мин`;
    if (diff < 86400) return `${Math.floor(diff / 3600)} ч`;
    if (diff < 604800) return `${Math.floor(diff / 86400)} дн`;
    return date.toLocaleDateString();
  };

  const getFullUrl = (path: string) => {
    if (!path) return "";
    if (path.startsWith("http")) return path;
    return `${API_URL}${path}`;
  };

  // Группируем сторис по пользователям
  const groupedStories = stories.reduce(
    (acc, story) => {
      if (!acc[story.user_id]) {
        acc[story.user_id] = [];
      }
      acc[story.user_id].push(story);
      return acc;
    },
    {} as Record<string, Story[]>,
  );

  // Берем только одну сторис от каждого пользователя
  const uniqueUserStories = Object.values(groupedStories).map(
    (userStories) => userStories[0],
  );

  return (
    <>
      <div className={style.storiesSection}>
        <div className={style.storiesContainer}>
          {/* Кнопка создания сторис */}
          <div className={style.storyItem}>
            <div
              className={style.createStory}
              onClick={() => !uploading && fileInputRef.current?.click()}
            >
              <div className={style.plusIcon}>
                {uploading ? (
                  <div className={style.spinner}></div>
                ) : (
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                    <path
                      d="M12 5V19M5 12H19"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                    />
                  </svg>
                )}
              </div>
              <span className={style.storyLabel}>
                {uploading ? "Загрузка..." : "Сторис"}
              </span>
            </div>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp,video/mp4,video/webm"
              style={{ display: "none" }}
              onChange={handleFileSelect}
              disabled={uploading}
            />
          </div>

          {/* Сторис пользователей */}
          {uniqueUserStories.map((story) => {
            const isMyStory = story.user_id === currentUser?.id;

            // Для своей сторис берем аватарку из currentUser
            const avatarSrc = isMyStory
              ? currentUser?.avatar_url
              : story.user_avatar;

            // Имя пользователя
            const displayName = isMyStory
              ? "Вы"
              : story.user_name || "Пользователь";

            return (
              <div
                key={story.id}
                className={style.storyItem}
                onClick={() => handleStoryClick(story)}
              >
                <div
                  className={`${style.storyAvatar} ${!story.viewed && !isMyStory ? style.notViewed : ""}`}
                >
                  {avatarSrc ? (
                    <img src={avatarSrc} alt={displayName} />
                  ) : (
                    <div className={style.defaultAvatar}>
                      {displayName.charAt(0).toUpperCase()}
                    </div>
                  )}
                </div>
                <span className={style.storyLabel}>
                  {displayName}
                  <br />
                  <small>{formatTimeAgo(story.created_at)}</small>
                </span>
              </div>
            );
          })}
        </div>
      </div>

      {/* Просмотрщик сторис */}
      {selectedStory && (
        <div className={style.storyViewerOverlay} onClick={closeStory}>
          <div
            className={style.storyViewer}
            onClick={(e) => e.stopPropagation()}
          >
            {selectedStory.type === "image" && (
              <div className={style.progressBar}>
                <div
                  className={style.progressFill}
                  style={{ width: `${progress}%` }}
                />
              </div>
            )}

            <div className={style.storyContent}>
              {selectedStory.type === "image" ? (
                <img
                  src={getFullUrl(selectedStory.file_url)}
                  alt="Story"
                  className={style.storyMedia}
                />
              ) : (
                <video
                  ref={videoRef}
                  src={getFullUrl(selectedStory.file_url)}
                  className={style.storyMedia}
                  autoPlay
                  playsInline
                  onEnded={handleVideoEnded}
                />
              )}
            </div>

            <div className={style.storyInfo}>
              <div className={style.storyUserInfo}>
                {/* Аватарка в просмотрщике */}
                {selectedStory.user_id === currentUser?.id ? (
                  currentUser?.avatar_url ? (
                    <img src={currentUser.avatar_url} alt="" />
                  ) : (
                    <div className={style.defaultAvatarSmall}>В</div>
                  )
                ) : selectedStory.user_avatar ? (
                  <img src={selectedStory.user_avatar} alt="" />
                ) : (
                  <div className={style.defaultAvatarSmall}>
                    {(selectedStory.user_name || "П")[0].toUpperCase()}
                  </div>
                )}
                <div>
                  <div className={style.userName}>
                    {selectedStory.user_id === currentUser?.id
                      ? "Вы"
                      : selectedStory.user_name || "Пользователь"}
                  </div>
                  <div className={style.storyTime}>
                    {formatTimeAgo(selectedStory.created_at)}
                  </div>
                </div>
              </div>
              <button className={style.closeButton} onClick={closeStory}>
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M18 6L6 18M6 6l12 12"
                    stroke="white"
                    strokeWidth="2"
                    strokeLinecap="round"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

export default StoriesSection;
