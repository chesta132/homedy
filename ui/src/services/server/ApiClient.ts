import axios, { isAxiosError, type AxiosInstance, type AxiosRequestConfig } from "axios";
import { ServerSuccess } from "./ServerSuccess";
import { ServerError } from "./ServerError";
import type { ApiResponse } from "@/types/server";
import { AUTH_ERROR_CODES } from "@/types/server";

export type ApiConfig<D = any> = AxiosRequestConfig<D>;

/**
 * Thin HTTP client bound to a shared axios instance + a URL prefix.
 *
 * Sub-clients (auth, samba, sambaConfig, convert) are plain instances of this class
 * with a different prefix — they do NOT spawn their own sub-clients, so there
 * is no recursive construction and no stack overflow.
 */
class ApiClient {
  private readonly instance: AxiosInstance;
  private readonly prefix: string;

  // These are only populated on the root client (prefix === "")
  readonly auth!: ApiClient;
  readonly samba!: ApiClient;
  readonly sambaConfig!: ApiClient;
  readonly convert!: ApiClient;

  constructor(instance: AxiosInstance, prefix = "") {
    this.instance = instance;
    this.prefix = prefix;

    // Only the root client builds sub-clients to avoid infinite recursion
    if (prefix === "") {
      (this as any).auth = new ApiClient(instance, "/auth");
      (this as any).samba = new ApiClient(instance, "/samba");
      (this as any).sambaConfig = new ApiClient(instance, "/samba/config");
      (this as any).convert = new ApiClient(instance, "/convert");
    }
  }

  private async request<T>(config: ApiConfig): Promise<ServerSuccess<T>> {
    try {
      const response = await this.instance.request<ApiResponse<T>>({
        ...config,
        url: `${this.prefix}${config.url ?? ""}`,
      });
      return new ServerSuccess(response);
    } catch (error) {
      if (isAxiosError(error) && error.response) {
        throw new ServerError(error as any);
      }
      throw error;
    }
  }

  get<T>(url: string, config?: ApiConfig) {
    return this.request<T>({ ...config, url, method: "GET" });
  }

  post<T>(url: string, data?: unknown, config?: ApiConfig) {
    return this.request<T>({ ...config, url, method: "POST", data });
  }

  put<T>(url: string, data?: unknown, config?: ApiConfig) {
    return this.request<T>({ ...config, url, method: "PUT", data });
  }

  patch<T>(url: string, data?: unknown, config?: ApiConfig) {
    return this.request<T>({ ...config, url, method: "PATCH", data });
  }

  delete<T>(url: string, config?: ApiConfig) {
    return this.request<T>({ ...config, url, method: "DELETE" });
  }

  /**
   * For file-download endpoints that return a Blob (not a JSON ApiResponse).
   * Extracts filename from Content-Disposition header automatically.
   */
  async postBlob(
    url: string,
    data?: unknown,
    config?: ApiConfig
  ): Promise<{ blob: Blob; filename: string }> {
    const res = await this.instance.request<Blob>({
      ...config,
      url: `${this.prefix}${url}`,
      method: "POST",
      data,
      responseType: "blob",
      // Let browser set Content-Type + boundary automatically for FormData.
      // Explicitly unsetting overrides the instance-level application/json default.
      headers: {
        ...config?.headers,
        ...(data instanceof FormData ? { "Content-Type": undefined } : {}),
      },
    });
    const disposition: string = res.headers["content-disposition"] ?? "";
    const match = disposition.match(/filename="?([^";\r\n]+)"?/);
    const filename = match?.[1]?.trim() ?? "download";
    return { blob: res.data, filename };
  }
}

/**
 * Creates the single shared axios instance with auth-error interceptor.
 * Base URL /api is proxied to the backend by Vite in dev.
 */
function createAxiosInstance(): AxiosInstance {
  const instance = axios.create({
    baseURL: "/api",
    withCredentials: true,
    headers: { "Content-Type": "application/json" },
  });

  instance.interceptors.response.use(
    (res) => res,
    (error) => {
      if (isAxiosError(error) && error.response) {
        const code = error.response.data?.data?.code as string;
        if (AUTH_ERROR_CODES.includes(code as any)) {
          if (
            window.location.pathname !== "/signin" &&
            window.location.pathname !== "/signup"
          ) {
            window.location.href = "/signin";
          }
        }
      }
      return Promise.reject(error);
    }
  );

  return instance;
}

const api = new ApiClient(createAxiosInstance());
export default api;
