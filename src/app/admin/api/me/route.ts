import { NextResponse } from "next/server";

import {
  gatewayBase,
  getAccessToken,
  refreshSession,
} from "@/lib/admin/session";

/**
 * GET /admin/api/me
 *
 * Returns the current admin's profile+role, transparently refreshing the
 * access token once if it has expired. 401 when there is no valid session.
 */
export async function GET() {
  let token = await getAccessToken();

  const fetchMe = (t: string) =>
    fetch(`${gatewayBase()}/apis/adminauth/v1/me`, {
      headers: { Authorization: `Bearer ${t}` },
      cache: "no-store",
    });

  let res = token
    ? await fetchMe(token)
    : new Response(null, { status: 401 });

  if (res.status === 401) {
    const refreshed = await refreshSession();
    if (!refreshed) {
      return NextResponse.json({ message: "unauthorized" }, { status: 401 });
    }
    token = refreshed;
    res = await fetchMe(token);
  }

  if (!res.ok) {
    return NextResponse.json({ message: "unauthorized" }, { status: 401 });
  }

  return NextResponse.json(await res.json());
}
