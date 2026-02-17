"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { format } from "date-fns";
import {
  Calendar,
  Users,
  Settings,
  Shield,
  ArrowRight,
  BarChart3,
  FileText,
  UserCog,
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { PageHeader } from "@/components/shared/page-header";
import { CardSkeleton } from "@/components/shared/loading-skeleton";
import { Roles } from "@/stores/auth-store";
import { getActiveReviewPeriod } from "@/lib/api/dashboard";
import { getStaffList } from "@/lib/api/staff";
import type { PerformanceReviewPeriod } from "@/types/performance";

export default function AdminDashboardPage() {
  const { data: session, status } = useSession();
  const router = useRouter();
  const userRoles = session?.user?.roles ?? [];

  const isAdmin = userRoles.includes(Roles.Admin) || userRoles.includes(Roles.SuperAdmin);

  const [loading, setLoading] = useState(true);
  const [reviewPeriod, setReviewPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [totalStaff, setTotalStaff] = useState(0);

  useEffect(() => {
    if (status === "loading") return;
    if (!isAdmin) {
      router.replace("/");
      return;
    }

    async function load() {
      setLoading(true);
      try {
        const [rpRes, staffRes] = await Promise.allSettled([
          getActiveReviewPeriod(),
          getStaffList(),
        ]);
        if (rpRes.status === "fulfilled" && rpRes.value?.data) {
          setReviewPeriod(rpRes.value.data);
        }
        if (staffRes.status === "fulfilled" && staffRes.value?.data) {
          setTotalStaff(Array.isArray(staffRes.value.data) ? staffRes.value.data.length : 0);
        }
      } catch {
        // Non-fatal
      } finally {
        setLoading(false);
      }
    }

    load();
  }, [status, isAdmin, router]);

  if (!isAdmin) return null;

  const quickLinks = [
    { label: "Performance Setup", href: "/review-periods-list", icon: Settings, description: "Manage review periods and categories" },
    { label: "Staff Management", href: "/staff_list", icon: Users, description: "View and manage all staff members" },
    { label: "Manage Roles", href: "/manage_roles", icon: Shield, description: "Configure roles and permissions" },
    { label: "Audit Logs", href: "/audit-logs", icon: FileText, description: "View system audit trail" },
    { label: "Reports", href: "/pms-period-scores-report", icon: BarChart3, description: "Performance scores and analytics" },
    { label: "User Management", href: "/staff_mgt", icon: UserCog, description: "Manage user accounts and access" },
  ];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Administrator's Dashboard"
        description="System overview and administrative controls"
        breadcrumbs={[{ label: "Admin Dashboard" }]}
      />

      {/* Stats Cards */}
      {loading ? (
        <div className="grid gap-4 md:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => <CardSkeleton key={i} />)}
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-3">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Staff</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{totalStaff}</div>
              <p className="text-xs text-muted-foreground">Registered system users</p>
            </CardContent>
          </Card>

          <Card className="border-primary/20 bg-primary/5">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active Review Period</CardTitle>
              <Calendar className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              {reviewPeriod ? (
                <>
                  <div className="text-lg font-bold">{reviewPeriod.name}</div>
                  <p className="text-xs text-muted-foreground">
                    {format(new Date(reviewPeriod.startDate), "dd MMM yyyy")} –{" "}
                    {format(new Date(reviewPeriod.endDate), "dd MMM yyyy")}
                  </p>
                </>
              ) : (
                <>
                  <div className="text-lg font-bold text-muted-foreground">None</div>
                  <p className="text-xs text-muted-foreground">No active review period</p>
                </>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Max Points</CardTitle>
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{reviewPeriod?.maxPoints ?? "—"}</div>
              <p className="text-xs text-muted-foreground">Per staff this period</p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Quick Links */}
      <div>
        <h2 className="mb-4 text-lg font-semibold">Quick Actions</h2>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {quickLinks.map((link) => (
            <Link key={link.href} href={link.href}>
              <Card className="cursor-pointer transition-colors hover:bg-accent">
                <CardContent className="flex items-start gap-4 py-4">
                  <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                    <link.icon className="h-5 w-5 text-primary" />
                  </div>
                  <div className="flex-1">
                    <p className="text-sm font-medium">{link.label}</p>
                    <p className="text-xs text-muted-foreground">{link.description}</p>
                  </div>
                  <ArrowRight className="mt-1 h-4 w-4 text-muted-foreground" />
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
