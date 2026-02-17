export { auth as middleware } from "@/lib/auth";

export const config = {
  matcher: [
    /*
     * Match all request paths except for:
     * - /login
     * - /api/auth (NextAuth routes)
     * - /_next/static (static files)
     * - /_next/image (image optimization)
     * - /favicon.ico
     * - /public files
     */
    "/((?!login|api/auth|_next/static|_next/image|favicon.ico|loginasset|assets|avater.jpg).*)",
  ],
};
