import { z } from 'zod';

export const loginSchema = z.object({
  email: z.string().trim().email('Enter a valid email address.'),
  password: z.string().min(1, 'Enter your password.'),
});

export type LoginCredentials = z.infer<typeof loginSchema>;

export const userSchema = z
  .object({
    email: z.string().email(),
    first_name: z.string().optional().default(''),
    last_name: z.string().optional().default(''),
  })
  .transform((data) => ({
    email: data.email,
    firstName: data.first_name,
    lastName: data.last_name,
  }));

export type User = z.infer<typeof userSchema>;
