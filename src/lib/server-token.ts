// SERVER-ONLY. Mints a short-lived HS256 admin bearer token so Server
// Components can call the backend's JWT-guarded admin endpoints (draws + users
// lists) needed to build the public winners feed.
//
// This module MUST NOT be imported by any Client Component. It reads the
// server-only `JWT_SECRET` (never a NEXT_PUBLIC_* var) and uses node:crypto —
// no third-party JWT dependency. The `server-only` package is not installed in
// this template, so isolation is enforced by convention: only api.ts (itself
// server-only) imports this file.
import crypto from "node:crypto";

/** Default dev secret — mirrors the backend gateway default. Override in prod. */
const DEFAULT_SECRET = "dev-insecure-change-me";

/** Token lifetime. The backend only validates signature + exp (no role claim). */
const TOKEN_TTL_SECONDS = 300;

/** base64url-encode a Buffer or string (RFC 7515, no padding). */
function base64url(input: Buffer | string): string {
  const buf = typeof input === "string" ? Buffer.from(input, "utf8") : input;
  return buf
    .toString("base64")
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
}

/**
 * Mint an HS256 JWT for the SSR service account. Claims: `sub` identifies the
 * caller, `exp` bounds its lifetime. Signed with HMAC-SHA256 over
 * `base64url(header).base64url(payload)` using the shared secret.
 */
export function mintAdminToken(): string {
  const secret = process.env.JWT_SECRET || DEFAULT_SECRET;

  const header = { alg: "HS256", typ: "JWT" };
  const now = Math.floor(Date.now() / 1000);
  const payload = { sub: "frontend-ssr", exp: now + TOKEN_TTL_SECONDS };

  const signingInput = `${base64url(JSON.stringify(header))}.${base64url(
    JSON.stringify(payload),
  )}`;
  const signature = crypto
    .createHmac("sha256", secret)
    .update(signingInput)
    .digest();

  return `${signingInput}.${base64url(signature)}`;
}
