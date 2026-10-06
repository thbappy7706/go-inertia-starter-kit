export interface User {
  id: number;
  name: string;
  email: string;
  email_verified_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface Product {
  id: number;
  name: string;
  slug: string;
  description: string | null;
  price: number;
  status: boolean;
  created_at: string;
  updated_at: string;
}

export interface PaginationMeta<T> {
  data: T[];
  current_page: number;
  last_page: number;
  per_page: number;
  total: number;
  from: number;
  to: number;
}

export interface Auth {
  user: User | null;
}

export interface FlashMessages {
  success?: string;
  error?: string;
}

export interface PageProps {
  auth: Auth;
  flash: FlashMessages;
  errors: Record<string, string>;
  [key: string]: unknown;
}