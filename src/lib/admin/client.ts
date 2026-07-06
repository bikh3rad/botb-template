"use client";

// Client-side admin API helper. All calls go to the SAME-ORIGIN proxy
// (/admin/api/gw/...), which attaches the bearer token server-side from the
// HttpOnly cookie. No token is ever handled in the browser.

const GW = "/admin/api/gw";

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
    this.name = "ApiError";
  }
}

async function parse(res: Response): Promise<unknown> {
  const text = await res.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

/** Typed request against a backend path (e.g. "/apis/competition/v1/admin/..."). */
export async function api<T = unknown>(
  path: string,
  init?: RequestInit,
): Promise<T> {
  const res = await fetch(`${GW}${path}`, {
    ...init,
    headers: {
      ...(init?.body && !(init.body instanceof FormData)
        ? { "Content-Type": "application/json" }
        : {}),
      ...init?.headers,
    },
  });

  const data = await parse(res);
  if (!res.ok) {
    const message =
      (data as { message?: string })?.message ?? `request failed (${res.status})`;
    throw new ApiError(res.status, message);
  }
  return data as T;
}

export const apiGet = <T>(path: string) => api<T>(path);
export const apiPost = <T>(path: string, body?: unknown) =>
  api<T>(path, { method: "POST", body: body ? JSON.stringify(body) : undefined });
export const apiPut = <T>(path: string, body?: unknown) =>
  api<T>(path, { method: "PUT", body: body ? JSON.stringify(body) : undefined });
export const apiDelete = <T>(path: string) => api<T>(path, { method: "DELETE" });

/** Upload multipart form data (used by the media manager). */
export const apiUpload = <T>(path: string, form: FormData) =>
  api<T>(path, { method: "POST", body: form });

/** Browser-facing MinIO base for rendering stored media directly. */
export function mediaBaseUrl(): string {
  return (
    process.env.NEXT_PUBLIC_MEDIA_BASE_URL || "http://localhost:9000"
  );
}

/** Build a public object URL from a media ref. */
export function mediaUrl(ref?: {
  bucket: string;
  object_key: string;
}): string | null {
  if (!ref?.bucket || !ref?.object_key) return null;
  return `${mediaBaseUrl()}/${ref.bucket}/${ref.object_key}`;
}
