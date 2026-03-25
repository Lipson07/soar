import { z } from "zod";

export const loginSchema = z.object({
  login: z
    .string()
    .min(3, "Логин должен содержать минимум 3 символа")
    .max(20, "Логин не должен превышать 20 символов")
    .regex(
      /^[a-zA-Z0-9_]+$/,
      "Логин может содержать только буквы, цифры и underscore",
    ),

  password: z
    .string()
    .min(1, "Введите пароль")
    .min(6, "Пароль должен содержать минимум 6 символов"),
});

export const registerSchema = z
  .object({
    login: z
      .string()
      .min(3, "Логин должен содержать минимум 3 символа")
      .max(20, "Логин не должен превышать 20 символов")
      .regex(
        /^[a-zA-Z0-9_]+$/,
        "Логин может содержать только буквы, цифры и underscore",
      ),

    email: z
      .string()
      .min(1, "Email обязателен")
      .email("Введите корректный email"),

    password: z
      .string()
      .min(8, "Пароль должен быть минимум 8 символов")
      .regex(/[A-Z]/, "Пароль должен содержать хотя бы одну заглавную букву")
      .regex(/[a-z]/, "Пароль должен содержать хотя бы одну строчную букву")
      .regex(/[0-9]/, "Пароль должен содержать хотя бы одну цифру")
      .regex(
        /[^a-zA-Z0-9]/,
        "Пароль должен содержать хотя бы один специальный символ",
      ),

    confirmPassword: z.string().min(1, "Подтвердите пароль"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Пароли не совпадают",
    path: ["confirmPassword"],
  });
