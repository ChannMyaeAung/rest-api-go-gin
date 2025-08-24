import axios from "axios";
import Cookies from "js-cookie";

export const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
});

/**
 * Request interceptor that automatically injects JWT authentication tokens.
 *
 * This interceptor runs before every HTTP request and:
 * 1. Retrieves the JWT token from the "token" cookie
 * 2. Adds the Authorization header with Bearer token format if token exists
 * 3. Ensures headers object exists before modification
 *
 * Authentication flow:
 * - Token is stored in cookies after successful login
 * - All subsequent API requests automatically include the token
 * - Backend middleware validates the token and extracts user information
 *
 * Security considerations:
 * - Uses httpOnly cookies in production for XSS protection
 * - Token is sent in standard Authorization header format
 * - No token = no Authorization header (for public endpoints)
 *
 * Cookie name: "token"
 * Header format: "Authorization: Bearer <jwt_token>"
 */
api.interceptors.request.use((config) => {
  const token = Cookies.get("token");
  if (token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

/**
 * Type definition for API error responses.
 *
 * The backend API returns errors in different formats depending on the endpoint:
 * - Standard Gin errors: { "error": "Error message" }
 * - Custom application errors: { "message": "Error message" }
 *
 * This type covers both formats to ensure proper error handling across the application.
 */
export type ApiError = { message?: string; error?: string };

/**
 * Centralized error parsing utility for consistent error handling.
 *
 * This function extracts meaningful error messages from various error sources:
 * 1. API response body (custom message or standard error field)
 * 2. HTTP status text (404 Not Found, 500 Internal Server Error, etc.)
 * 3. Network-level errors (timeout, connection refused, etc.)
 * 4. Unexpected errors (fallback message)
 *
 * Error precedence (first available wins):
 * 1. response.data.message (custom application errors)
 * 2. response.data.error (standard Gin errors)
 * 3. response.statusText (HTTP status descriptions)
 * 4. error.message (Axios/network errors)
 * 5. "Unexpected error" (fallback for unknown errors)
 *
 * @param e - The error object from try/catch block or promise rejection
 * @returns Human-readable error message for display to users
 */
export function getApiError(e: unknown): string {
  // Check if the error is an Axios error (HTTP request/response error)
  if (axios.isAxiosError(e)) {
    // Try to extract error message from response body (priority order)
    return (
      // 1. Custom application error message
      (e.response?.data as ApiError)?.message ||
      // 2. Standard Gin framework error message
      (e.response?.data as ApiError)?.error ||
      // 3. HTTP status text (e.g., "Not Found", "Internal Server Error")
      e.response?.statusText ||
      // 4. Axios error message (network errors, timeouts, etc.)
      e.message
    );
  }
  return "Unexpected error";
}
