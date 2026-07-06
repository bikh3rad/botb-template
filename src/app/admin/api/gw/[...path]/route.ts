import { NextResponse, type NextRequest } from "next/server";

import {
  gatewayBase,
  getAccessToken,
  refreshSession,
} from "@/lib/admin/session";

// Same-origin proxy for admin API calls. The browser never sees a token: it
// calls /admin/api/gw/<backend-path>, and this handler attaches the bearer from
// the HttpOnly cookie, refreshing once on a 401. Bodies (including multipart
// uploads) are streamed straight through.

type Ctx = { params: Promise<{ path: string[] }> };

async function proxy(request: NextRequest, ctx: Ctx): Promise<Response> {
  const { path } = await ctx.params;
  const search = request.nextUrl.search;
  const target = `${gatewayBase()}/${path.join("/")}${search}`;

  // Buffer the body once so we can replay it after a token refresh.
  const method = request.method.toUpperCase();
  const hasBody = method !== "GET" && method !== "HEAD";
  const bodyBuf = hasBody ? await request.arrayBuffer() : undefined;

  const baseHeaders = new Headers();
  const contentType = request.headers.get("content-type");
  if (contentType) baseHeaders.set("content-type", contentType);

  const send = async (token: string) => {
    const headers = new Headers(baseHeaders);
    headers.set("Authorization", `Bearer ${token}`);
    return fetch(target, {
      method,
      headers,
      body: bodyBuf ? Buffer.from(bodyBuf) : undefined,
      cache: "no-store",
      redirect: "manual",
    });
  };

  let token = await getAccessToken();
  if (!token) {
    const refreshed = await refreshSession();
    if (!refreshed) {
      return NextResponse.json({ message: "unauthorized" }, { status: 401 });
    }
    token = refreshed;
  }

  let res = await send(token);

  // Access token expired mid-session → refresh once and retry.
  if (res.status === 401) {
    const refreshed = await refreshSession();
    if (!refreshed) {
      return NextResponse.json({ message: "unauthorized" }, { status: 401 });
    }
    res = await send(refreshed);
  }

  // Relay the backend response verbatim (status + body + content-type).
  const respHeaders = new Headers();
  const respType = res.headers.get("content-type");
  if (respType) respHeaders.set("content-type", respType);

  return new NextResponse(res.body, {
    status: res.status,
    headers: respHeaders,
  });
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const PATCH = proxy;
export const DELETE = proxy;
