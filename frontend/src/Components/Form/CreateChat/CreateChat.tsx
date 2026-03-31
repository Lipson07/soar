import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import ReactDOM from "react-dom";
import {
  closeCreateChat,
  selectModalChat,
} from "../../../store/modalChatSlice";
import style from "./CreateChat.module.scss";

function CreateChat() {
  const dispatch = useDispatch();
  const isOpen = useSelector(selectModalChat);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "unset";
    }
    return () => {
      document.body.style.overflow = "unset";
    };
  }, [isOpen]);

  const handleOverlayClick = (e: any) => {
    if (e.target === e.currentTarget) {
      dispatch(closeCreateChat());
    }
  };

  const handleClose = () => {
    dispatch(closeCreateChat());
  };

  if (!isOpen) return null;

  return (
    <div className={style.overlay} onClick={handleOverlayClick}>
      <div className={style.modal}>
        <button className={style.closeButton} onClick={handleClose}>
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
            <path
              d="M18 6L6 18M6 6L18 18"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
            />
          </svg>
        </button>

        <h2>Создать чат</h2>
      </div>
    </div>
  );
}

export default CreateChat;
