"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter, useParams } from "next/navigation";
import {
  Trophy,
  Users,
  Briefcase,
  AlertTriangle,
  CheckSquare,
  Star,
} from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatCard } from "@/components/shared/charts/stat-card";
import { StaffProfileCard } from "@/components/shared/charts/staff-profile-card";
import { PerformanceDoughnut } from "@/components/shared/charts/performance-doughnut";
import {
  getOrganogramPerformanceSummary,
  getEmployeeDetail,
  getStaffIdMask,
} from "@/lib/api/dashboard";
import { getGradeInfo, CHART_COLORS } from "@/lib/scorecard-helpers";
import type { EmployeeErpDetails, OrganogramPerformanceSummary } from "@/types/dashboard";
import type { StaffIdMask } from "@/types/staff";

const LEVEL_LABELS: Record<number, string> = {
  2: "Department",
  3: "Division",
  4: "Office",
};

export default function UnitScoreCardPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const params = useParams<{ unitId: string; reviewPeriodId: string; level: string }>();
  const unitId = params.unitId;
  const reviewPeriodId = params.reviewPeriodId;
  const level = parseInt(params.level, 10);
  const loggedInStaffId = session?.user?.id ?? "";
  const userRoles = session?.user?.roles ?? [];

  const [loading, setLoading] = useState(true);
  const [summary, setSummary] = useState<OrganogramPerformanceSummary | null>(null);
  const [manager, setManager] = useState<EmployeeErpDetails | null>(null);
  const [idMask, setIdMask] = useState<StaffIdMask | null>(null);

  useEffect(() => {
    if (!unitId || !reviewPeriodId || isNaN(level)) return;

    const load = async () => {
      setLoading(true);
      try {
        const sumRes = await getOrganogramPerformanceSummary(unitId, reviewPeriodId, level);
        if (sumRes?.data) {
          const s = sumRes.data;
          setSummary(s);

          const [empRes, maskRes] = await Promise.allSettled([
            getEmployeeDetail(s.managerId),
            getStaffIdMask(s.managerId),
          ]);

          if (empRes.status === "fulfilled" && empRes.value?.data) {
            const emp = empRes.value.data;
            setManager(emp);

            // Access control
            const canView =
              s.managerId === loggedInStaffId ||
              emp.headOfOfficeId === loggedInStaffId ||
              emp.headOfDivisionId === loggedInStaffId ||
              emp.headOfDepartmentId === loggedInStaffId ||
              userRoles.includes("HrReportAdmin");

            if (!canView) {
              router.replace("/access-denied");
              return;
            }
          }

          if (maskRes.status === "fulfilled" && maskRes.value?.data) {
            setIdMask(maskRes.value.data);
          }
        }
      } catch {
        /* ignore */
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [unitId, reviewPeriodId, level, loggedInStaffId, userRoles, router]);

  if (loading) return <PageSkeleton />;
  if (!summary) return <PageSkeleton />;

  const gradeInfo = getGradeInfo(summary.earnedPerformanceGrade);
  const levelLabel = LEVEL_LABELS[level] || "Unit";

  const feedbackData = [
    { name: "Total", value: summary.total360Feedbacks, color: CHART_COLORS.info },
    { name: "Completed", value: summary.completed360FeedbacksToTreat, color: CHART_COLORS.success },
    { name: "Pending", value: summary.pending360FeedbacksToTreat, color: CHART_COLORS.warning },
  ];

  return (
    <div className="space-y-6">
      <PageHeader
        title={`${summary.referenceName} ${summary.year} ${summary.reviewPeriod} - Score Card`}
        breadcrumbs={[
          { label: "Performance", href: "/performance-score-card" },
          { label: `${levelLabel} Score Card` },
        ]}
      />

      {manager && (
        <StaffProfileCard employee={manager} photoUrl={idMask?.currentStaffPhoto} />
      )}

      {/* Highlighted Performance Score Card */}
      <Card className="border-green-200 bg-gradient-to-r from-green-50 to-emerald-50">
        <CardContent className="flex flex-col items-center gap-2 py-6 sm:flex-row sm:justify-center sm:gap-8">
          <div className="text-center">
            <p className="text-sm font-medium text-green-700">Performance Score</p>
            <p className="text-4xl font-bold text-green-800">
              {summary.performanceScore.toFixed(1)}%
            </p>
          </div>
          <div className="text-center">
            <p className="text-sm font-medium text-green-700">Performance Grade</p>
            <p className={`text-3xl font-bold ${gradeInfo.bgClass} inline-block rounded-md px-3 py-1`}>
              {gradeInfo.label}
            </p>
          </div>
        </CardContent>
      </Card>

      {/* 6 Stat Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <StatCard
          title="Points Earned"
          value={`${summary.actualScore.toFixed(1)} / ${summary.maxPoint}`}
          icon={Trophy}
          iconColor="#f59e0b"
        />
        <StatCard
          title="Total Staff"
          value={summary.totalStaff}
          icon={Users}
          iconColor="#3b82f6"
        />
        <StatCard
          title="Total Work Products"
          value={summary.totalWorkProducts}
          icon={Briefcase}
          iconColor="#22c55e"
        />
        <StatCard
          title="Overdue Work Products"
          value={summary.totalWorkProductsBehindSchedule}
          icon={AlertTriangle}
          iconColor="#ef4444"
        />
        <StatCard
          title="Timely Work Products"
          value={summary.totalWorkProductsCompletedOnSchedule}
          icon={CheckSquare}
          iconColor="#22c55e"
        />
        <StatCard
          title="Living the Values"
          value={summary.percentageGapsClosure > 0 ? `${summary.percentageGapsClosure.toFixed(1)}%` : "N/A"}
          icon={Star}
          iconColor="#8b5cf6"
        />
      </div>

      {/* Charts */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* 360 Review Feedbacks Bar Chart */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">360 Review Feedbacks</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={feedbackData}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
                <XAxis
                  dataKey="name"
                  tick={{ fill: "hsl(var(--muted-foreground))", fontSize: 12 }}
                />
                <YAxis tick={{ fill: "hsl(var(--muted-foreground))", fontSize: 12 }} />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "hsl(var(--popover))",
                    border: "1px solid hsl(var(--border))",
                    borderRadius: "6px",
                  }}
                />
                <Bar dataKey="value" radius={[4, 4, 0, 0]}>
                  {feedbackData.map((entry, index) => (
                    <Cell key={index} fill={entry.color} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Gap Closure Doughnut */}
        <PerformanceDoughnut
          title="Gap Closure"
          percentage={summary.percentageGapsClosure}
        />
      </div>
    </div>
  );
}
