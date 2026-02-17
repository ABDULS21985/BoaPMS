"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import Link from "next/link";
import { format } from "date-fns";
import {
  Target,
  FolderKanban,
  ClipboardCheck,
  TrendingUp,
  Calendar,
  ArrowRight,
  Briefcase,
  Users,
  MessageSquare,
  AlertTriangle,
  RefreshCw,
} from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  Legend,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/shared/progress-bar";
import { PageHeader } from "@/components/shared/page-header";
import { CardSkeleton } from "@/components/shared/loading-skeleton";
import {
  getRequestStatistics,
  getPerformanceStatistics,
  getWorkProductStatistics,
  getWorkProductDetailsStatistics,
  getActiveReviewPeriod,
} from "@/lib/api/dashboard";
import { getStaffObjectives } from "@/lib/api/pms-engine";
import type {
  FeedbackRequestDashboardStats,
  PerformancePointsStats,
  WorkProductDashboardStats,
  WorkProductDetailsDashboardStats,
} from "@/types/dashboard";
import type { PerformanceReviewPeriod } from "@/types/performance";

const PIE_COLORS = ["hsl(var(--chart-1))", "hsl(var(--chart-2))", "hsl(var(--chart-3))", "hsl(var(--chart-4))"];

export default function DashboardPage() {
  const { data: session } = useSession();
  const user = session?.user;
  const staffId = user?.id ?? "";

  const [loading, setLoading] = useState(true);
  const [reviewPeriod, setReviewPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [requestStats, setRequestStats] = useState<FeedbackRequestDashboardStats | null>(null);
  const [performanceStats, setPerformanceStats] = useState<PerformancePointsStats | null>(null);
  const [wpStats, setWpStats] = useState<WorkProductDashboardStats | null>(null);
  const [wpDetails, setWpDetails] = useState<WorkProductDetailsDashboardStats | null>(null);
  const [objectiveCount, setObjectiveCount] = useState(0);

  const loadDashboard = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const [rpRes, reqRes, perfRes, wpRes, wpDetRes] = await Promise.allSettled([
        getActiveReviewPeriod(),
        getRequestStatistics(staffId),
        getPerformanceStatistics(staffId),
        getWorkProductStatistics(staffId),
        getWorkProductDetailsStatistics(staffId),
      ]);

      if (rpRes.status === "fulfilled" && rpRes.value?.data) {
        setReviewPeriod(rpRes.value.data);
        // Load objectives count using review period
        try {
          const objRes = await getStaffObjectives(staffId, rpRes.value.data.periodId);
          if (objRes?.data) {
            setObjectiveCount(Array.isArray(objRes.data) ? objRes.data.length : 0);
          }
        } catch { /* ignore */ }
      }
      if (reqRes.status === "fulfilled" && reqRes.value?.data) setRequestStats(reqRes.value.data);
      if (perfRes.status === "fulfilled" && perfRes.value?.data) setPerformanceStats(perfRes.value.data);
      if (wpRes.status === "fulfilled" && wpRes.value?.data) setWpStats(wpRes.value.data);
      if (wpDetRes.status === "fulfilled" && wpDetRes.value?.data) setWpDetails(wpDetRes.value.data);
    } catch {
      // Dashboard load errors are non-fatal
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (staffId) loadDashboard();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [staffId]);

  // --- Chart data ---
  const wpChartData = wpStats
    ? [
        { name: "Active", value: wpStats.noActiveWorkProducts },
        { name: "Awaiting Eval", value: wpStats.noWorkProductsAwaitingEvaluation },
        { name: "Closed", value: wpStats.noWorkProductsClosed },
        { name: "Pending Approval", value: wpStats.noWorkProductsPendingApproval },
      ]
    : [];

  const pointsChartData = performanceStats
    ? [
        { name: "Accumulated", points: performanceStats.accumulatedPoints },
        { name: "Deducted", points: performanceStats.deductedPoints },
        { name: "Actual", points: performanceStats.actualPoints },
        { name: "Max", points: performanceStats.maxPoints },
      ]
    : [];

  const scorePercentage = performanceStats && performanceStats.maxPoints > 0
    ? Math.round((performanceStats.actualPoints / performanceStats.maxPoints) * 100)
    : 0;

  // --- Near-due work products ---
  const nearDueWPs = wpDetails?.workProducts
    ?.filter((wp) => wp.recordStatus === 8) // Active
    .sort((a, b) => new Date(a.endDate).getTime() - new Date(b.endDate).getTime())
    .slice(0, 3) ?? [];

  return (
    <div className="space-y-6">
      <PageHeader
        title={`Welcome, ${user?.firstName ?? "User"}`}
        description="Your performance management overview"
        actions={
          <Button variant="outline" size="sm" onClick={loadDashboard} disabled={loading}>
            <RefreshCw className={`mr-2 h-4 w-4 ${loading ? "animate-spin" : ""}`} />
            Refresh
          </Button>
        }
      />

      {/* Review Period Banner */}
      {reviewPeriod && (
        <Card className="border-primary/20 bg-primary/5">
          <CardContent className="flex items-center gap-4 py-4">
            <Calendar className="h-8 w-8 text-primary" />
            <div className="flex-1">
              <p className="text-sm font-medium text-muted-foreground">Active Review Period</p>
              <p className="text-lg font-semibold">{reviewPeriod.name}</p>
            </div>
            <div className="text-right text-sm text-muted-foreground">
              <p>{format(new Date(reviewPeriod.startDate), "dd MMM yyyy")}</p>
              <p>to {format(new Date(reviewPeriod.endDate), "dd MMM yyyy")}</p>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Statistics Cards */}
      {loading ? (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <CardSkeleton key={i} />
          ))}
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Work Products</CardTitle>
              <FolderKanban className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{wpStats?.noAllWorkProducts ?? 0}</div>
              <p className="text-xs text-muted-foreground">
                {wpStats?.noWorkProductsClosed ?? 0} closed of {wpStats?.noAllWorkProducts ?? 0} total
              </p>
              {wpStats && wpStats.noAllWorkProducts > 0 && (
                <Progress
                  value={Math.round((wpStats.noWorkProductsClosed / wpStats.noAllWorkProducts) * 100)}
                  className="mt-2 h-1.5"
                />
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Objectives Planned</CardTitle>
              <Target className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{objectiveCount}</div>
              <p className="text-xs text-muted-foreground">
                Active objectives this period
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Pending Requests</CardTitle>
              <ClipboardCheck className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{requestStats?.pendingRequests ?? 0}</div>
              <p className="text-xs text-muted-foreground">
                {requestStats?.breachedRequests ?? 0} breached SLA
              </p>
              <Link
                href="/assigned-pending-requests"
                className="mt-1 inline-flex items-center text-xs text-primary hover:underline"
              >
                View all <ArrowRight className="ml-1 h-3 w-3" />
              </Link>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Performance Score</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{scorePercentage}%</div>
              <p className="text-xs text-muted-foreground">
                {performanceStats?.actualPoints?.toFixed(1) ?? 0} of {performanceStats?.maxPoints ?? 0} points
              </p>
              <Progress value={scorePercentage} className="mt-2 h-1.5" />
            </CardContent>
          </Card>
        </div>
      )}

      {/* Charts Section */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* Performance Points Chart */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Performance Points</CardTitle>
            <CardDescription>Points breakdown for current period</CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="h-[250px] animate-pulse rounded bg-muted" />
            ) : (
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={pointsChartData}>
                  <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
                  <XAxis dataKey="name" className="text-xs" tick={{ fill: "hsl(var(--muted-foreground))" }} />
                  <YAxis className="text-xs" tick={{ fill: "hsl(var(--muted-foreground))" }} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "hsl(var(--popover))",
                      border: "1px solid hsl(var(--border))",
                      borderRadius: "var(--radius)",
                      color: "hsl(var(--popover-foreground))",
                    }}
                  />
                  <Bar dataKey="points" radius={[4, 4, 0, 0]}>
                    {pointsChartData.map((_, index) => (
                      <Cell key={index} fill={PIE_COLORS[index % PIE_COLORS.length]} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            )}
          </CardContent>
        </Card>

        {/* Work Products by Status */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Work Products by Status</CardTitle>
            <CardDescription>Distribution of work product statuses</CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="h-[250px] animate-pulse rounded bg-muted" />
            ) : (
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  <Pie
                    data={wpChartData}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={90}
                    dataKey="value"
                    paddingAngle={2}
                  >
                    {wpChartData.map((_, index) => (
                      <Cell key={index} fill={PIE_COLORS[index % PIE_COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "hsl(var(--popover))",
                      border: "1px solid hsl(var(--border))",
                      borderRadius: "var(--radius)",
                      color: "hsl(var(--popover-foreground))",
                    }}
                  />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Points Progress Section */}
      {performanceStats && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Points Progress</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <div className="mb-1 flex justify-between text-sm">
                <span>Accumulated Points</span>
                <span className="text-muted-foreground">
                  {performanceStats.accumulatedPoints.toFixed(1)} / {performanceStats.maxPoints}
                </span>
              </div>
              <Progress
                value={performanceStats.maxPoints > 0 ? (performanceStats.accumulatedPoints / performanceStats.maxPoints) * 100 : 0}
                className="h-2.5"
              />
            </div>
            <div>
              <div className="mb-1 flex justify-between text-sm">
                <span>Deducted Points</span>
                <span className="text-destructive">
                  -{performanceStats.deductedPoints.toFixed(1)}
                </span>
              </div>
              <Progress
                value={performanceStats.maxPoints > 0 ? (performanceStats.deductedPoints / performanceStats.maxPoints) * 100 : 0}
                className="h-2.5 [&>div]:bg-destructive"
              />
            </div>
            <div>
              <div className="mb-1 flex justify-between text-sm">
                <span className="font-medium">Actual Points</span>
                <span className="font-medium">
                  {performanceStats.actualPoints.toFixed(1)} / {performanceStats.maxPoints}
                </span>
              </div>
              <Progress
                value={performanceStats.maxPoints > 0 ? (performanceStats.actualPoints / performanceStats.maxPoints) * 100 : 0}
                className="h-2.5"
              />
            </div>
          </CardContent>
        </Card>
      )}

      {/* Near-Due Work Products */}
      {nearDueWPs.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Work Products Nearing Deadline</CardTitle>
            <CardDescription>Active work products closest to their due dates</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {nearDueWPs.map((wp) => (
                <div key={wp.workProductId} className="flex items-start gap-4 rounded-lg border p-3">
                  <div className="flex-1 space-y-1">
                    <p className="text-sm font-medium">{wp.name}</p>
                    {wp.description && (
                      <p className="text-xs text-muted-foreground line-clamp-1">{wp.description}</p>
                    )}
                    <div className="flex items-center gap-4 text-xs text-muted-foreground">
                      <span>Due: {format(new Date(wp.endDate), "dd MMM yyyy")}</span>
                      {wp.objectiveName && <span>Obj: {wp.objectiveName}</span>}
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-medium">{Math.round(wp.percentageTaskCompletion)}%</p>
                    <p className="text-xs text-muted-foreground">
                      {wp.tasksCompleted}/{wp.totalTasks} tasks
                    </p>
                    <Progress value={wp.percentageTaskCompletion} className="mt-1 h-1.5 w-20" />
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Request Statistics */}
      {requestStats && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Request Statistics</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-4">
              <div className="rounded-lg border p-3 text-center">
                <p className="text-2xl font-bold">{requestStats.completedRequests}</p>
                <p className="text-xs text-muted-foreground">Completed</p>
              </div>
              <div className="rounded-lg border p-3 text-center">
                <p className="text-2xl font-bold">{requestStats.pendingRequests}</p>
                <p className="text-xs text-muted-foreground">Pending</p>
              </div>
              <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-3 text-center">
                <p className="text-2xl font-bold text-destructive">{requestStats.breachedRequests}</p>
                <p className="text-xs text-muted-foreground">Breached SLA</p>
              </div>
              <div className="rounded-lg border p-3 text-center">
                <p className="text-2xl font-bold">{requestStats.pending360FeedbacksToTreat}</p>
                <p className="text-xs text-muted-foreground">Pending 360 Reviews</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Quick Links */}
      <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-5">
        {[
          { label: "My Objectives", href: "/my-objectives", icon: Target },
          { label: "My Work Products", href: "/myworkproducts", icon: FolderKanban },
          { label: "My Projects", href: "/my-projects", icon: Briefcase },
          { label: "My Committees", href: "/my-committees", icon: Users },
          { label: "My Requests", href: "/myrequests", icon: MessageSquare },
        ].map((link) => (
          <Link key={link.href} href={link.href}>
            <Card className="cursor-pointer transition-colors hover:bg-accent">
              <CardContent className="flex items-center gap-3 py-3">
                <link.icon className="h-5 w-5 text-primary" />
                <span className="text-sm font-medium">{link.label}</span>
                <ArrowRight className="ml-auto h-4 w-4 text-muted-foreground" />
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}
