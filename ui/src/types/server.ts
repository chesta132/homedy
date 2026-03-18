// Standard API error response shape
export interface ErrorResponse {
  code: string;
  message: string;
  title?: string;
  details?: string;
  field?: string;
}

// Standard API response envelope
export interface ApiResponse<T, Success extends boolean = true> {
  meta: {
    status: Success extends true ? "SUCCESS" : "ERROR";
    hasNext?: boolean;
    nextOffset?: number | null;
    information?: string;
  };
  data: T;
}

// Auth error codes that should redirect to signin
export const AUTH_ERROR_CODES = ["UNAUTHORIZED", "INVALID_AUTH", "INVALID_TOKEN"] as const;
export type AuthErrorCode = (typeof AUTH_ERROR_CODES)[number];

export type StateErrorServer = Omit<ErrorResponse, "field">;
