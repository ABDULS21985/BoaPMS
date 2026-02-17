"use client";

import type { ReactNode } from "react";
import { TopNavbar } from "@/components/layout/top-navbar";
import { Toaster } from "@/components/ui/sonner";

export default function ScoreCardLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-background">
      <TopNavbar />
      <main className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        {children}
      </main>
      <Toaster richColors position="top-right" />
    </div>
  );
}
