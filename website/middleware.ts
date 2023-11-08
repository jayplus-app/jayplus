import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  return NextResponse.redirect(new URL("/admin/booking", request.url));
}

export const config = {
  matcher: "/admin/dashboard",
};
