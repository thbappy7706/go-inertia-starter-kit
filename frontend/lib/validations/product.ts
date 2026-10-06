import { z } from "zod"

export const productSchema = z.object({
  name: z
    .string()
    .min(2, "Name must be at least 2 characters")
    .max(255, "Name cannot exceed 255 characters"),
  slug: z
    .string()
    .min(1, "Slug is required")
    .regex(/^[a-z0-9]+(?:-[a-z0-9]+)*$/, "Slug must only contain lowercase letters, numbers, and hyphens"),
  description: z.string().optional().default(""),
  price: z.coerce.number().positive("Price must be a positive number"),
  status: z.boolean().default(true),
})

export type ProductFormValues = z.infer<typeof productSchema>