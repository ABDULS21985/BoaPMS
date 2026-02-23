import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";
import type { AuthenticateResponse } from "@/types/auth";

declare module "next-auth" {
  interface Session {
    accessToken: string;
    refreshToken: string;
    user: {
      id: string;
      name: string;
      email: string;
      firstName: string;
      lastName: string;
      roles: string[];
      permissions: string[];
      organizationalUnit?: string;
    };
  }

  interface User {
    id: string;
    name: string;
    email: string;
    firstName: string;
    lastName: string;
    roles: string[];
    permissions: string[];
    organizationalUnit?: string;
    accessToken: string;
    refreshToken: string;
    expiresAt: number;
  }
}

declare module "next-auth" {
  interface JWT {
    id: string;
    firstName: string;
    lastName: string;
    roles: string[];
    permissions: string[];
    organizationalUnit?: string;
    accessToken: string;
    refreshToken: string;
    expiresAt: number;
  }
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export const { handlers, signIn, signOut, auth } = NextAuth({
  secret: process.env.AUTH_SECRET,
  providers: [
    Credentials({
      name: "credentials",
      credentials: {
        username: { label: "Username", type: "text" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        if (!credentials?.username || !credentials?.password) return null;

        try {
          const res = await fetch(`${API_URL}/auth/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              username: credentials.username,
              password: credentials.password,
            }),
          });

          if (!res.ok) return null;

          const json = await res.json();
          // Go API wraps responses in { success, data }
          const data = json.data ?? json;

          return {
            id: data.user_id,
            name: `${data.first_name} ${data.last_name}`,
            email: data.email,
            firstName: data.first_name,
            lastName: data.last_name,
            roles: data.roles ?? [],
            permissions: data.permissions ?? [],
            organizationalUnit: data.organizational_unit,
            accessToken: data.access_token,
            refreshToken: data.refresh_token,
            expiresAt: data.expires_at,
          };
        } catch {
          return null;
        }
      },
    }),
  ],
  session: { strategy: "jwt", maxAge: 45 * 60 },
  pages: { signIn: "/login" },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.id = user.id as string;
        token.firstName = user.firstName;
        token.lastName = user.lastName;
        token.roles = user.roles;
        token.permissions = user.permissions;
        token.organizationalUnit = user.organizationalUnit;
        token.accessToken = user.accessToken;
        token.refreshToken = user.refreshToken;
        token.expiresAt = user.expiresAt;
      }

      // Token rotation: refresh if expired
      if (token.expiresAt && Date.now() / 1000 > (token.expiresAt as number)) {
        try {
          const res = await fetch(`${API_URL}/auth/refresh`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ refresh_token: token.refreshToken }),
          });
          if (res.ok) {
            const json = await res.json();
            const refreshed = json.data ?? json;
            token.accessToken = refreshed.access_token;
            token.refreshToken = refreshed.refresh_token;
            token.expiresAt = refreshed.expires_at;
          }
        } catch {
          // Refresh failed — user will be redirected to login
        }
      }

      return token;
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken as string;
      session.refreshToken = token.refreshToken as string;
      session.user = {
        ...session.user,
        id: token.id as string,
        firstName: token.firstName as string,
        lastName: token.lastName as string,
        roles: token.roles as string[],
        permissions: token.permissions as string[],
        organizationalUnit: token.organizationalUnit as string | undefined,
      };
      return session;
    },
  },
});
