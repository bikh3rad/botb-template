import { NextResponse } from "next/server";

import {
  clearSessionCookies,
  gatewayBase,
  getRefreshToken,
} from "@/lib/admin/session";

/**
 * POST /admin/api/logout
 *
 * Revokes the refresh token at the backend (best-effort) and clears the
 * session cookies.
 */
export async function POST() {
  const refresh = await getRefreshToken();

  if (refresh) {
    try {
      await fetch(`${gatewayBase()}/apis/adminauth/v1/logout`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refresh }),
        cache: "no-store",
      });
    } catch {
      // Non-fatal: we clear local cookies regardless.
    }
  }

  await clearSessionCookies();
  return NextResponse.json({ ok: true });
}
