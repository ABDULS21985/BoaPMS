"use client";

import Link from "next/link";
import { useSession, signOut } from "next-auth/react";
import {
  ClipboardList,
  Clock,
  LogOut,
  Moon,
  Sun,
  User,
} from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Separator } from "@/components/ui/separator";
import { useUIStore } from "@/stores/ui-store";

export function TopNavbar() {
  const { data: session } = useSession();
  const { theme, toggleTheme } = useUIStore();
  const user = session?.user;

  const initials = user
    ? `${(user.firstName ?? "")[0] ?? ""}${(user.lastName ?? "")[0] ?? ""}`.toUpperCase()
    : "?";

  const displayName = user
    ? `${user.firstName} ${user.lastName}`
    : "Loading...";

  return (
    <header className="sticky top-0 z-30 flex h-14 items-center gap-4 border-b bg-background px-4">
      <SidebarTrigger className="-ml-1" />
      <Separator orientation="vertical" className="h-6" />

      <div className="flex flex-1 items-center gap-4">
        <h2 className="hidden text-sm font-semibold text-muted-foreground sm:block">
          Performance Management System
        </h2>
      </div>

      <div className="flex items-center gap-3">
        {/* Pending Requests */}
        <Link href="/assigned-pending-requests">
          <Button variant="ghost" size="icon" className="relative">
            <ClipboardList className="h-5 w-5 text-muted-foreground" />
            <Clock className="absolute bottom-1 right-1 h-3 w-3 text-muted-foreground" />
          </Button>
        </Link>

        {/* Theme Toggle */}
        <Button variant="ghost" size="icon" onClick={toggleTheme}>
          {theme === "dark" ? (
            <Sun className="h-5 w-5" />
          ) : (
            <Moon className="h-5 w-5" />
          )}
        </Button>

        {/* User Dropdown */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              className="flex items-center gap-2 px-2"
            >
              <div className="hidden text-right sm:block">
                <p className="text-sm font-medium leading-none">
                  {displayName}
                </p>
                <p className="text-xs text-muted-foreground">
                  {user?.email ?? ""}
                </p>
              </div>
              <Avatar className="h-8 w-8">
                <AvatarFallback className="bg-primary/10 text-sm font-medium text-primary">
                  {initials}
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56">
            <DropdownMenuLabel>
              <div className="flex flex-col gap-1">
                <p className="font-medium">{displayName}</p>
                <p className="text-xs text-muted-foreground">
                  {user?.email}
                </p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem asChild>
              <Link href="/my-profile" className="cursor-pointer">
                <User className="mr-2 h-4 w-4" />
                My Profile
              </Link>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={() => signOut({ callbackUrl: "/login" })}
              className="cursor-pointer text-destructive focus:text-destructive"
            >
              <LogOut className="mr-2 h-4 w-4" />
              Log Out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
