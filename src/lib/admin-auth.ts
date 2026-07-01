import { timingSafeEqual } from "node:crypto"
import type { NextRequest } from "next/server"

/**
 * Shared-secret guard for admin-only routes.
 *
 * A request is authorized when it presents the value of the `ADMIN_TOKEN`
 * environment variable, either as a bearer token
 * (`Authorization: Bearer <token>`) or as the `admin_token` cookie.
 *
 * Fails closed: if `ADMIN_TOKEN` is unset, every request is rejected.
 */

export const ADMIN_COOKIE = "admin_token"

function safeEqual(provided: string, expected: string): boolean {
  const providedBuf = Buffer.from(provided)
  const expectedBuf = Buffer.from(expected)
  // timingSafeEqual throws on length mismatch; compare lengths first.
  if (providedBuf.length !== expectedBuf.length) return false
  return timingSafeEqual(providedBuf, expectedBuf)
}

export function isAuthorizedAdmin(request: NextRequest): boolean {
  const expected = process.env.ADMIN_TOKEN
  if (!expected) return false

  const authHeader = request.headers.get("authorization")
  const bearer = authHeader?.startsWith("Bearer ")
    ? authHeader.slice("Bearer ".length).trim()
    : null

  const cookie = request.cookies.get(ADMIN_COOKIE)?.value ?? null
  const provided = bearer ?? cookie
  if (!provided) return false

  return safeEqual(provided, expected)
}
