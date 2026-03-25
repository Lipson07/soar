import { z } from "zod";

export const loginSchema = z.object({
  login: z.string().min(3, "логин должен содержать минимум 3 символа"),
});
