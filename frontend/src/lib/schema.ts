import { z } from "zod";

export const loginSchema = z.object({
  email: z.email(),
  password: z.string().min(6),
});

export const registerSchema = loginSchema.extend({
  name: z.string().min(2),
  email: z.email(),
  password: z.string().min(6),
});

export const eventSchema = z.object({
  name: z.string().min(2),
  description: z.string().min(10),
  location: z.string().min(2),
  date: z.date(),
});
