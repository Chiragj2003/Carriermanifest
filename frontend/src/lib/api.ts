/**
 * API Client — Axios instance pre-configured for the CareerManifest backend.
 *
 * Includes automatic JWT token injection via request interceptor.
 * Exports typed API modules: authAPI, assessmentAPI, questionsAPI, adminAPI.
 */
import axios from "axios";

// API client configured for the CareerManifest backend.
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  headers: {
    "Content-Type": "application/json",
  },
});

// Attach JWT token to every request if available.
api.interceptors.request.use((config: any) => {
  if (typeof window !== "undefined") {
    const token = localStorage.getItem("cm_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

// Handle 401 responses globally.
api.interceptors.response.use(
  (response: any) => response,
  (error: any) => {
    if (error.response?.status === 401 && typeof window !== "undefined") {
      localStorage.removeItem("cm_token");
      localStorage.removeItem("cm_user");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

// ============================================================
// AUTH API
// ============================================================
export const authAPI = {
  register: (data: { name: string; email: string; password: string }) =>
    api.post("/api/auth/register", data),
  login: (data: { email: string; password: string }) =>
    api.post("/api/auth/login", data),
  getProfile: () => api.get("/api/auth/profile"),
};

// ============================================================
// QUESTIONS API
// ============================================================
export const questionsAPI = {
  getActive: () => api.get("/api/questions"),
};

// ============================================================
// ASSESSMENT API
// ============================================================
export const assessmentAPI = {
  submit: (data: { answers: { question_id: number; selected: number }[] }) =>
    api.post("/api/assessment", data),
  getById: (id: number) => api.get(`/api/assessment/${id}`),
  list: () => api.get("/api/assessment"),
};

// ============================================================
// CHAT API
// ============================================================
export const chatAPI = {
  send: (message: string, assessmentId: number) =>
    api.post("/api/chat", { message, assessment_id: assessmentId }),
};

// ============================================================
// ADMIN API
// ============================================================
export const adminAPI = {
  getStats: () => api.get("/api/admin/stats"),
  getQuestions: () => api.get("/api/admin/questions"),
  createQuestion: (data: any) => api.post("/api/admin/questions", data),
  updateQuestion: (id: number, data: any) =>
    api.put(`/api/admin/questions/${id}`, data),
};

export default api;
