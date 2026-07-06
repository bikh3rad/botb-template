import { NextResponse, type NextRequest } from "next/server";

import { ACCESS_COOKIE, REFRESH_COOKIE } from "@/lib/admin/session";

// Guards the admin panel: any /admin route without a session cookie is
// redirected to the login page. This is a coarse gate (presence of a cookie);
// the real authorization is enforced by the backend role guards on every API
// call. The login page and the admin API routes are excluded so login itself
// and the token-refresh flow can run.
export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Allow the login page and all /admin/api/* handlers through.
  if (pathname === "/admin/login" || pathname.startsWith("/admin/api/")) {
    return NextResponse.next();
  }

  const hasSession =
    request.cookies.has(ACCESS_COOKIE) || request.cookies.has(REFRESH_COOKIE);

  if (!hasSession) {
    const url = request.nextUrl.clone();
    url.pathname = "/admin/login";
    url.searchParams.set("next", pathname);
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/admin/:path*"],
};
