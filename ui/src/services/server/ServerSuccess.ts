import type { AxiosResponse } from "axios";
import type { ApiResponse } from "@/types/server";

/**
 * Wraps a successful axios response into a typed result object.
 */
export class ServerSuccess<T> {
  readonly data: T;
  readonly meta: ApiResponse<T>["meta"];

  constructor(response: AxiosResponse<ApiResponse<T>>) {
    this.data = response.data.data;
    this.meta = response.data.meta;
  }
}
