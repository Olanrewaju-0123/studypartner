export interface User {
  id: number;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Note {
  id: number;
  user_id: number;
  title: string;
  content: string;
  file_type: string;
  file_name: string;
  file_size: number;
  created_at: string;
  updated_at: string;
}

export interface Summary {
  id: number;
  note_id: number;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface Flashcard {
  id: number;
  note_id: number;
  question: string;
  answer: string;
  created_at: string;
}

export interface Quiz {
  id: number;
  note_id: number;
  question: string;
  options: string[];
  answer: number;
  created_at: string;
}

export interface StudySession {
  id: number;
  user_id: number;
  note_id: number;
  type: string;
  score?: number;
  completed: boolean;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface UploadRequest {
  file: string; // Base64 encoded
  name: string;
}

export interface SearchRequest {
  query: string;
}

export interface SearchResult extends Note {
  similarity: number;
}
