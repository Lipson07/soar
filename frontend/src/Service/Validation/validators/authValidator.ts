import { z } from "zod";
import { loginSchema, registerSchema } from "../schemas/authSchema";
import type {
  LoginFormData,
  RegisterFormData,
  ValidationResult,
} from "../types/validation.types";

export class AuthValidator {
  static validateLogin(data: unknown): ValidationResult<LoginFormData> {
    const result = loginSchema.safeParse(data);

    if (result.success) {
      return {
        success: true,
        data: result.data,
      };
    } else {
      const errors = result.error.issues.map((err) => ({
        field: err.path.join("."),
        message: err.message,
      }));
      return {
        success: false,
        errors,
      };
    }
  }

  static validateLoginField(
    field: keyof LoginFormData,
    value: unknown,
  ): string | null {
    try {
      const fieldSchema = loginSchema.shape[field];
      fieldSchema.parse(value);
      return null;
    } catch (error) {
      if (error instanceof z.ZodError) {
        return error.issues[0]?.message || "Ошибка валидации";
      }
      return "Ошибка валидации";
    }
  }

  static validateRegister(data: unknown): ValidationResult<RegisterFormData> {
    const result = registerSchema.safeParse(data);

    if (result.success) {
      return {
        success: true,
        data: result.data,
      };
    } else {
      const errors = result.error.issues.map((err) => ({
        field: err.path.join("."),
        message: err.message,
      }));
      return {
        success: false,
        errors,
      };
    }
  }

  static validateRegisterField(
    field: keyof RegisterFormData,
    value: unknown,
  ): string | null {
    try {
      const fieldSchema = registerSchema.shape[field];
      fieldSchema.parse(value);
      return null;
    } catch (error) {
      if (error instanceof z.ZodError) {
        return error.issues[0]?.message || "Ошибка валидации";
      }
      return "Ошибка валидации";
    }
  }
}
