// SERVER-ONLY admin session helpers.
//
// Token-storage design (deviation from "access token in memory"): both the
// access and refresh tokens are kept in HttpOnly, SameSite=Lax cookies on the
// FRONTEND origin, never exposed to client JavaScript. The browser talks only
// to same-origin Next.js route handlers (src/app/admin/api/*), which attach the
// bearer token to the backend gateway server-side. Rationale: the gateway
// (:8080) and the site (:3000) are different origins, so an HttpOnly refresh
// cookie scoped to the gateway would not be sent by the browser on same-origin
// admin XHRs; proxying through the frontend origin keeps every token out of JS
// (stricter than access-in-memory) and out of localStorage entirely.
import { cookies } from "next/headers";

/** Server-side gateway base URL (internal compose address preferred). */
export function gatewayBase(): string {
  return (
    process.env.API_BASE_URL_INTERNAL ||
    process.env.NEXT_PUBLIC_API_BASE_URL ||
    "http://localhost:8080"
  );
}

export const ACCESS_COOKIE = "botb_admin_access";
export const REFRESH_COOKIE = "botb_admin_refresh";

/** Shared cookie options. httpOnly + sameSite=lax keeps tokens out of JS. */
function cookieOptions(maxAge: number) {
  return {
    httpOnly: true,
    sameSite: "lax" as const,
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge,
  };
}

export interface AdminProfile {
  id: string;
  name: string;
  email: string;
  role: string;
  is_active?: boolean;
}

export interface TokenResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  refresh_token: string;
  admin: AdminProfile;
}

/** Persist a freshly issued token pair into HttpOnly cookies. */
export async function setSessionCookies(tokens: TokenResponse): Promise<void> {
  const store = await cookies();
  store.set(ACCESS_COOKIE, tokens.access_token, cookieOptions(tokens.expires_in));
  // Refresh cookie lives longer than the access token; 7 days matches the
  // backend default refresh TTL.
  store.set(REFRESH_COOKIE, tokens.refresh_token, cookieOptions(60 * 60 * 24 * 7));
}

/** Clear the session cookies (logout). */
export async function clearSessionCookies(): Promise<void> {
  const store = await cookies();
  store.delete(ACCESS_COOKIE);
  store.delete(REFRESH_COOKIE);
}

export async function getAccessToken(): Promise<string | undefined> {
  return (await cookies()).get(ACCESS_COOKIE)?.value;
}

export async function getRefreshToken(): Promise<string | undefined> {
  return (await cookies()).get(REFRESH_COOKIE)?.value;
}

/**
 * Exchange the stored refresh token for a new pair, rotating cookies. Returns
 * the new access token, or null when refresh is impossible (caller should treat
 * as logged-out). Refresh-token rotation is enforced by the backend.
 */
export async function refreshSession(): Promise<string | null> {
  const refresh = await getRefreshToken();
  if (!refresh) return null;

  try {
    const res = await fetch(`${gatewayBase()}/apis/adminauth/v1/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refresh }),
      cache: "no-store",
    });
    if (!res.ok) {
      await clearSessionCookies();
      return null;
    }
    const tokens = (await res.json()) as TokenResponse;
    await setSessionCookies(tokens);
    return tokens.access_token;
  } catch {
    return null;
  }
}
