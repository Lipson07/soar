import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import ReactDOM from "react-dom";
import {
  closeCreateChat,
  selectModalChat,
} from "../../../store/modalChatSlice";
import style from "./CreateChat.module.scss";
import { Input } from "../../UI";

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

        <form>
          <Input
            background="#333333"
            color="white"
            width="500px"
            label="Название чата"
            placeholder="Введите название"
            type="text"
            name="name"
          />

          <div className={style.formGroup}>
            <label className={style.label}>Описание (необязательно)</label>
            <textarea
              className={style.textarea}
              name="description"
              placeholder="Введите описание (максимум 500 символов)"
              rows={3}
              maxLength={500}
            />
          </div>

          <div className={style.typeGroup}>
            <label className={style.typeLabel}>Тип чата</label>
            <div className={style.typeOptions}>
              <label className={style.typeOption}>
                <input
                  type="radio"
                  name="type"
                  value="private"
                  defaultChecked
                />
                <span className={style.typeRadio}></span>
                <div className={style.typeContent}>
                  <span className={style.typeTitle}>Приватный</span>
                  <span className={style.typeDesc}>
                    Только для приглашенных
                  </span>
                </div>
              </label>

              <label className={style.typeOption}>
                <input type="radio" name="type" value="group" />
                <span className={style.typeRadio}></span>
                <div className={style.typeContent}>
                  <span className={style.typeTitle}>Группа</span>
                  <span className={style.typeDesc}>
                    Общайтесь с несколькими людьми
                  </span>
                </div>
              </label>

              <label className={style.typeOption}>
                <input type="radio" name="type" value="channel" />
                <span className={style.typeRadio}></span>
                <div className={style.typeContent}>
                  <span className={style.typeTitle}>Канал</span>
                  <span className={style.typeDesc}>
                    Для трансляций и широкой аудитории
                  </span>
                </div>
              </label>
            </div>
          </div>

          <div className={style.actions}>
            <button
              type="button"
              className={style.cancelButton}
              onClick={handleClose}
            >
              Отмена
            </button>
            <button type="submit" className={style.submitButton}>
              Создать чат
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default CreateChat;
