import React, { useRef, useEffect, useState } from "react";
import style from "./Register.module.scss";
import { Link, useNavigate } from "react-router-dom";
import { Input } from "../../UI";

const Register = () => {
  const birdRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const navigate = useNavigate();

  const [formData, setFormData] = useState({
    login: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  useEffect(() => {
    const bird = birdRef.current;
    const cont = containerRef.current;

    setTimeout(() => {
      if (bird && cont) {
        cont.style.transform = "translateY(0)";

        requestAnimationFrame(() => {
          requestAnimationFrame(() => {
            bird.style.transform = "translateY(-25vh)";
            cont.style.transform = "translateY(-25vh)";
          });
        });
      }
    }, 500);
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const response = await fetch("http://localhost:8080/api/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: formData.login,
          email: formData.email,
          password: formData.password,
          role: "user",
        }),
      });
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Ошибка регистрации");
      }

      navigate("/login");
    } catch (error) {
      console.error("Ошибка регистрации:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <main className={style.main}>
      <div className={style.backgroundObjects}>
        <div className={`${style.cloud} ${style.cloud1}`} />
        <div className={`${style.cloud} ${style.cloud2}`} />
        <div className={`${style.cloud} ${style.cloud3}`} />

        <div className={`${style.feather} ${style.feather1}`}>
          <svg viewBox="0 0 24 24" fill="white">
            <path d="M12 2C8 4 4 8 3 12c-1 6 3 10 7 8 2-1 3-4 2-7-1-2-3-3-5-2-1 1-1 3 0 4 1 1 2 0 2-1 0-1-1-1-1-1 0 0 2-1 3 1 1 2 1 5-1 7-3 2-7 0-6-5 1-4 5-8 9-10l-1-4z" />
          </svg>
        </div>
        <div className={`${style.feather} ${style.feather2}`}>
          <svg viewBox="0 0 24 24" fill="white">
            <path d="M12 2C8 4 4 8 3 12c-1 6 3 10 7 8 2-1 3-4 2-7-1-2-3-3-5-2-1 1-1 3 0 4 1 1 2 0 2-1 0-1-1-1-1-1 0 0 2-1 3 1 1 2 1 5-1 7-3 2-7 0-6-5 1-4 5-8 9-10l-1-4z" />
          </svg>
        </div>
        <div className={`${style.feather} ${style.feather3}`}>
          <svg viewBox="0 0 24 24" fill="white">
            <path d="M12 2C8 4 4 8 3 12c-1 6 3 10 7 8 2-1 3-4 2-7-1-2-3-3-5-2-1 1-1 3 0 4 1 1 2 0 2-1 0-1-1-1-1-1 0 0 2-1 3 1 1 2 1 5-1 7-3 2-7 0-6-5 1-4 5-8 9-10l-1-4z" />
          </svg>
        </div>
        <div className={`${style.feather} ${style.feather4}`}>
          <svg viewBox="0 0 24 24" fill="white">
            <path d="M12 2C8 4 4 8 3 12c-1 6 3 10 7 8 2-1 3-4 2-7-1-2-3-3-5-2-1 1-1 3 0 4 1 1 2 0 2-1 0-1-1-1-1-1 0 0 2-1 3 1 1 2 1 5-1 7-3 2-7 0-6-5 1-4 5-8 9-10l-1-4z" />
          </svg>
        </div>

        <div className={`${style.smallBird} ${style.smallBird1}`}>
          <svg viewBox="0 0 24 24" fill="white">
            <path d="M21 12c0 3-2 5-4 6-3 2-8 1-10-3-2-4 0-9 4-11 3-1 6 0 8 3 1 1 2 3 2 5z" />
          </svg>
        </div>
        <div className={`${style.smallBird} ${style.smallBird2}`}>
          <svg viewBox="0 0 24 24" fill="white">
            <path d="M21 12c0 3-2 5-4 6-3 2-8 1-10-3-2-4 0-9 4-11 3-1 6 0 8 3 1 1 2 3 2 5z" />
          </svg>
        </div>
        <div className={`${style.smallBird} ${style.smallBird3}`}>
          <svg viewBox="0 0 24 24" fill="white">
            <path d="M21 12c0 3-2 5-4 6-3 2-8 1-10-3-2-4 0-9 4-11 3-1 6 0 8 3 1 1 2 3 2 5z" />
          </svg>
        </div>
        <div className={`${style.smallBird} ${style.smallBird4}`}>
          <svg viewBox="0 0 24 24" fill="white">
            <path d="M21 12c0 3-2 5-4 6-3 2-8 1-10-3-2-4 0-9 4-11 3-1 6 0 8 3 1 1 2 3 2 5z" />
          </svg>
        </div>

        <div className={`${style.wind} ${style.wind1}`} />
        <div className={`${style.wind} ${style.wind2}`} />
        <div className={`${style.wind} ${style.wind3}`} />
      </div>

      <div className={style.birds} ref={birdRef}>
        <svg
          version="1.0"
          xmlns="http://www.w3.org/2000/svg"
          width="48.000000pt"
          height="48.000000pt"
          viewBox="0 0 48.000000 48.000000"
          preserveAspectRatio="xMidYMid meet"
        >
          <g
            transform="translate(0.000000,48.000000) scale(0.100000,-0.100000)"
            fill="#ffffff"
            stroke="none"
          >
            <path
              d="M31 446 c-16 -19 0 -90 22 -94 14 -3 16 2 10 30 l-6 33 50 -46 c96
-90 143 -100 185 -39 15 22 39 42 53 46 54 14 74 -49 23 -73 -23 -12 -25 -18
-21 -51 7 -52 -10 -89 -52 -114 -43 -25 -34 -48 11 -27 40 20 66 61 72 115 5
52 19 69 65 81 l27 8 -43 48 c-49 53 -78 59 -119 28 l-25 -20 -44 45 c-58 59
-70 59 -84 -3 -5 -26 -3 -33 8 -33 8 0 17 10 21 23 5 21 7 21 42 -15 39 -38
41 -68 5 -68 -20 0 -97 57 -145 109 -33 35 -39 37 -55 17z"
            />
            <path
              d="M20 316 c0 -7 8 -24 18 -37 11 -13 22 -37 26 -53 4 -18 17 -32 34
-38 l27 -10 -37 -40 c-63 -65 -46 -128 33 -128 31 0 47 7 70 29 31 32 40 81
14 81 -8 0 -15 -9 -15 -20 0 -33 -37 -62 -73 -58 -52 5 -47 36 13 92 55 52 59
76 12 76 -33 0 -52 16 -52 45 0 11 -9 26 -20 33 -11 7 -20 19 -20 27 0 8 -7
15 -15 15 -8 0 -15 -6 -15 -14z"
            />
          </g>
        </svg>
        <svg
          version="1.0"
          xmlns="http://www.w3.org/2000/svg"
          width="48.000000pt"
          height="48.000000pt"
          viewBox="0 0 48.000000 48.000000"
          preserveAspectRatio="xMidYMid meet"
          className={style.mirrored}
        >
          <g
            transform="translate(0.000000,48.000000) scale(0.100000,-0.100000)"
            fill="#ffffff"
            stroke="none"
          >
            <path
              d="M31 446 c-16 -19 0 -90 22 -94 14 -3 16 2 10 30 l-6 33 50 -46 c96
-90 143 -100 185 -39 15 22 39 42 53 46 54 14 74 -49 23 -73 -23 -12 -25 -18
-21 -51 7 -52 -10 -89 -52 -114 -43 -25 -34 -48 11 -27 40 20 66 61 72 115 5
52 19 69 65 81 l27 8 -43 48 c-49 53 -78 59 -119 28 l-25 -20 -44 45 c-58 59
-70 59 -84 -3 -5 -26 -3 -33 8 -33 8 0 17 10 21 23 5 21 7 21 42 -15 39 -38
41 -68 5 -68 -20 0 -97 57 -145 109 -33 35 -39 37 -55 17z"
            />
            <path
              d="M20 316 c0 -7 8 -24 18 -37 11 -13 22 -37 26 -53 4 -18 17 -32 34
-38 l27 -10 -37 -40 c-63 -65 -46 -128 33 -128 31 0 47 7 70 29 31 32 40 81
14 81 -8 0 -15 -9 -15 -20 0 -33 -37 -62 -73 -58 -52 5 -47 36 13 92 55 52 59
76 12 76 -33 0 -52 16 -52 45 0 11 -9 26 -20 33 -11 7 -20 19 -20 27 0 8 -7
15 -15 15 -8 0 -15 -6 -15 -14z"
            />
          </g>
        </svg>
      </div>

      <form
        className={style.container}
        ref={containerRef}
        onSubmit={handleSubmit}
      >
        <h1>Регистрация</h1>

        <Input
          background="#333333"
          color="white"
          width="320px"
          label="Имя пользователя"
          placeholder="Введите имя"
          type="text"
          name="login"
          onChange={handleChange}
          value={formData.login}
        />

        <Input
          background="#333333"
          color="white"
          width="320px"
          label="Email"
          placeholder="Введите email"
          type="email"
          name="email"
          onChange={handleChange}
          value={formData.email}
        />
        <Input
          background="#333333"
          color="white"
          width="320px"
          label="Пароль"
          placeholder="Введите пароль"
          type="password"
          name="password"
          onChange={handleChange}
          value={formData.password}
        />

        <Input
          background="#333333"
          color="white"
          width="320px"
          label="Подтвердите пароль"
          placeholder="Введите пароль повторно"
          type="password"
          name="confirmPassword"
          onChange={handleChange}
          value={formData.confirmPassword}
        />
        <div className={style.buttondiv}>
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Регистрация..." : "Зарегистрироваться"}
          </button>
        </div>

        <div className={style.txt}>
          <p>
            Уже есть аккаунт?{" "}
            <Link to="/login">
              <span>Войти</span>
            </Link>
          </p>
          <p>Политика конфиденциальности</p>
        </div>
      </form>
    </main>
  );
};

export default Register;
