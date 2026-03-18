import type { AxiosError } from "axios";
import type { ApiResponse, ErrorResponse } from "@/types/server";

/**
 * Wraps an axios error response into a typed, queryable error object.
 */
export class ServerError {
  readonly data: ErrorResponse;
  readonly status: number;

  constructor(error: AxiosError<ApiResponse<ErrorResponse, false>>) {
    if (!error.response?.data) {
      throw new Error("Invalid server error: missing response data");
    }
    this.data = error.response.data.data;
    this.status = error.response.status;
  }

  getCode() {
    return this.data.code;
  }

  getMessage() {
    return this.data.message;
  }

  getField() {
    return this.data.field;
  }

  getTitle() {
    return this.data.title ?? "Something went wrong";
  }
}
