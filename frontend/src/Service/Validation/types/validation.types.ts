import { z } from "zod";
import { loginSchema, registerSchema } from "../schemas/authSchema";

export type LoginFormData = z.infer<typeof loginSchema>;
export type RegisterFormData = z.infer<typeof registerSchema>;

export interface ValidationError {
  field: string;
  message: string;
}

export interface ValidationResult<T> {
  success: boolean;
  data?: T;
  errors?: ValidationError[];
}

export type FormField = keyof LoginFormData | keyof RegisterFormData;
