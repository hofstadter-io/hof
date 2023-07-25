import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// This function can be marked `async` if using `await` inside
export function middleware(request: NextRequest) {
  // console.log("middleware-PRE:", request.url)
  const response = NextResponse.next();
  // console.log("middleware-POST:", request.url)
  return response;
}

// See "Matching Paths" below to learn more
export const config = {
  matcher: ["/", "/api/(.*)", "/auth/(.*)", "/essays/(.*)"],
};
