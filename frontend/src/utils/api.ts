import {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
  Note,
  UploadRequest,
  Summary,
  Flashcard,
  Quiz,
  StudySession,
  SearchRequest,
  SearchResult,
} from "@/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
    if (typeof window !== "undefined") {
      this.token = localStorage.getItem("token");
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(options.headers as Record<string, string>),
    };

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ error: "Network error" }));
      throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
  }

  setToken(token: string) {
    this.token = token;
    if (typeof window !== "undefined") {
      localStorage.setItem("token", token);
    }
  }

  clearToken() {
    this.token = null;
    if (typeof window !== "undefined") {
      localStorage.removeItem("token");
    }
  }

  // Auth endpoints
  async register(data: RegisterRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>("/api/auth/register", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getCurrentUser(): Promise<User> {
    return this.request<User>("/api/auth/me");
  }

  // Notes endpoints
  async uploadNote(data: UploadRequest): Promise<Note> {
    return this.request<Note>("/api/notes/upload", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getNotes(): Promise<Note[]> {
    return this.request<Note[]>("/api/notes");
  }

  async getNote(id: number): Promise<Note> {
    return this.request<Note>(`/api/notes/${id}`);
  }

  async deleteNote(id: number): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/notes/${id}`, {
      method: "DELETE",
    });
  }

  async searchNotes(data: SearchRequest): Promise<SearchResult[]> {
    return this.request<SearchResult[]>("/api/notes/search", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  // Study endpoints
  async getSummary(noteId: number): Promise<Summary> {
    return this.request<Summary>(`/api/study/notes/${noteId}/summary`);
  }

  async generateSummary(noteId: number): Promise<Summary> {
    return this.request<Summary>(`/api/study/notes/${noteId}/summary`, {
      method: "POST",
    });
  }

  async getFlashcards(noteId: number): Promise<Flashcard[]> {
    return this.request<Flashcard[]>(`/api/study/notes/${noteId}/flashcards`);
  }

  async generateFlashcards(noteId: number): Promise<Flashcard[]> {
    return this.request<Flashcard[]>(`/api/study/notes/${noteId}/flashcards`, {
      method: "POST",
    });
  }

  async getQuiz(noteId: number): Promise<Quiz[]> {
    return this.request<Quiz[]>(`/api/study/notes/${noteId}/quiz`);
  }

  async generateQuiz(noteId: number): Promise<Quiz[]> {
    return this.request<Quiz[]>(`/api/study/notes/${noteId}/quiz`, {
      method: "POST",
    });
  }

  async createStudySession(data: {
    note_id: number;
    type: string;
  }): Promise<StudySession> {
    return this.request<StudySession>("/api/study/sessions", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateStudySession(
    id: number,
    data: { score?: number; completed: boolean }
  ): Promise<StudySession> {
    return this.request<StudySession>(`/api/study/sessions/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
