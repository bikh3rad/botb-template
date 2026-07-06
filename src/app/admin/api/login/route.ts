import { NextResponse, type NextRequest } from "next/server";

import {
  gatewayBase,
  setSessionCookies,
  type TokenResponse,
} from "@/lib/admin/session";

/**
 * POST /admin/api/login
 *
 * Verifies admin credentials against the adminauth service and, on success,
 * stores the issued token pair in HttpOnly cookies. Returns only the admin
 * profile to the client — never the tokens.
 */
export async function POST(request: NextRequest) {
  let body: { email?: string; password?: string };
  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ message: "invalid request" }, { status: 400 });
  }

  const res = await fetch(`${gatewayBase()}/apis/adminauth/v1/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      // Forward the real client IP so the backend's per-IP rate limit is
      // keyed on the browser, not the SSR container.
      "X-Forwarded-For":
        request.headers.get("x-forwarded-for") ??
        request.headers.get("x-real-ip") ??
        "",
    },
    body: JSON.stringify({ email: body.email, password: body.password }),
    cache: "no-store",
  });

  if (!res.ok) {
    // Pass through the generic error + status (401 bad creds, 429 rate limit).
    const err = await res.json().catch(() => ({ message: "login failed" }));
    return NextResponse.json(err, { status: res.status });
  }

  const tokens = (await res.json()) as TokenResponse;
  await setSessionCookies(tokens);

  return NextResponse.json({ admin: tokens.admin });
}
