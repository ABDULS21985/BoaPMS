"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter, useParams } from "next/navigation";
import { Briefcase, Percent, Trophy, Award } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatCard } from "@/components/shared/charts/stat-card";
import { StaffProfileCard } from "@/components/shared/charts/staff-profile-card";
import { PerformanceDoughnut } from "@/components/shared/charts/performance-doughnut";
import { CompetencyBarChart } from "@/components/shared/charts/competency-bar-chart";
import { ScoreTrendChart } from "@/components/shared/charts/score-trend-chart";
import {
  getEmployeeDetail,
  getStaffIdMask,
  getStaffAnnualScoreCard,
} from "@/lib/api/dashboard";
import { getGradeFromScore, getGradeNumericValue, formatPercent } from "@/lib/scorecard-helpers";
import type { EmployeeErpDetails, StaffScoreCardDetails } from "@/types/dashboard";
import type { StaffIdMask } from "@/types/staff";

export default function StaffAnnualScoreCardPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const params = useParams<{ staffId: string; year: string }>();
  const staffId = params.staffId;
  const year = parseInt(params.year, 10);
  const loggedInStaffId = session?.user?.id ?? "";
  const userRoles = session?.user?.roles ?? [];

  const [loading, setLoading] = useState(true);
  const [employee, setEmployee] = useState<EmployeeErpDetails | null>(null);
  const [idMask, setIdMask] = useState<StaffIdMask | null>(null);
  const [scoreCards, setScoreCards] = useState<StaffScoreCardDetails[]>([]);

  useEffect(() => {
    if (!staffId || isNaN(year)) return;

    const load = async () => {
      setLoading(true);
      try {
        const [empRes, maskRes, annualRes] = await Promise.allSettled([
          getEmployeeDetail(staffId),
          getStaffIdMask(staffId),
          getStaffAnnualScoreCard(staffId, year),
        ]);

        if (empRes.status === "fulfilled" && empRes.value?.data) {
          const emp = empRes.value.data;
          setEmployee(emp);

          const canView =
            staffId === loggedInStaffId ||
            emp.supervisorId === loggedInStaffId ||
            emp.headOfOfficeId === loggedInStaffId ||
            emp.headOfDivisionId === loggedInStaffId ||
            emp.headOfDepartmentId === loggedInStaffId ||
            userRoles.includes("HrReportAdmin");

          if (!canView) {
            router.replace("/access-denied");
            return;
          }
        } else {
          router.replace("/404");
          return;
        }

        if (maskRes.status === "fulfilled" && maskRes.value?.data) {
          setIdMask(maskRes.value.data);
        }

        if (annualRes.status === "fulfilled" && annualRes.value?.scoreCards) {
          setScoreCards(annualRes.value.scoreCards);
        }
      } catch {
        /* ignore */
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [staffId, year, loggedInStaffId, userRoles, router]);

  if (loading) return <PageSkeleton />;
  if (!employee) return <PageSkeleton />;

  // Aggregate across all periods
  const totalWorkProducts = scoreCards.reduce((s, c) => s + c.totalWorkProducts, 0);
  const percentageGapsClosure = scoreCards.reduce((s, c) => s + c.percentageGapsClosure, 0);
  const actualPoints = scoreCards.reduce((s, c) => s + c.actualPoints, 0);
  const totalPoints = scoreCards.reduce((s, c) => s + c.maxPoints, 0);
  const percentageScore = totalPoints > 0 ? (100 * actualPoints) / totalPoints : 0;
  const grade = getGradeFromScore(percentageScore);

  // Aggregate competency categories: group by key, average
  const catMap = new Map<string, number[]>();
  for (const sc of scoreCards) {
    for (const [key, val] of Object.entries(sc.pmsCompetencyCategory || {})) {
      if (!catMap.has(key)) catMap.set(key, []);
      catMap.get(key)!.push(val);
    }
  }
  const categoryPoints = Array.from(catMap.entries()).map(([key, vals]) => ({
    name: key,
    value: vals.reduce((a, b) => a + b, 0) / vals.length,
  }));

  // Aggregate competencies: group by name, average ratingScores
  const compMap = new Map<string, number[]>();
  for (const sc of scoreCards) {
    for (const c of sc.pmsCompetencies || []) {
      if (!compMap.has(c.pmsCompetency)) compMap.set(c.pmsCompetency, []);
      compMap.get(c.pmsCompetency)!.push(c.ratingScore);
    }
  }
  const competencyData = Array.from(compMap.entries()).map(([name, vals]) => ({
    name,
    value: vals.reduce((a, b) => a + b, 0) / vals.length,
  }));

  // Per-period grade bar chart
  const gradeData = scoreCards.map((sc) => ({
    name: sc.reviewPeriodShortName || sc.reviewPeriod,
    value: getGradeNumericValue(sc.staffPerformanceGrade),
  }));

  // Score trend line chart
  const trendData = scoreCards.map((sc) => ({
    name: sc.reviewPeriodShortName || sc.reviewPeriod,
    score: Math.round(sc.percentageScore * 100) / 100,
  }));

  return (
    <div className="space-y-6">
      <PageHeader
        title={`${year} Performance - Score Card`}
        breadcrumbs={[
          { label: "Performance", href: "/performance-score-card" },
          { label: `${year} Annual Score Card` },
        ]}
      />

      <StaffProfileCard employee={employee} photoUrl={idMask?.currentStaffPhoto} />

      {/* 4 stat cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Total Work Products"
          value={totalWorkProducts}
          icon={Briefcase}
          iconColor="#3b82f6"
        />
        <StatCard
          title="Competency Gap Closure"
          value={formatPercent(percentageGapsClosure)}
          icon={Percent}
          iconColor="#22c55e"
        />
        <StatCard
          title="Points Earned"
          value={`${actualPoints.toFixed(1)} / ${totalPoints}`}
          icon={Trophy}
          iconColor="#f59e0b"
        />
        <StatCard
          title="Grade"
          value={grade}
          icon={Award}
          iconColor="#8b5cf6"
        />
      </div>

      {/* Charts Row 1: Grade Bar + Competency Bar */}
      <div className="grid gap-4 md:grid-cols-2">
        <CompetencyBarChart
          title="Performance Grade by Period"
          data={gradeData}
          height={250}
        />
        <CompetencyBarChart
          title="Competency Ratings (Aggregated)"
          data={competencyData}
          height={Math.max(250, competencyData.length * 40)}
        />
      </div>

      {/* Charts Row 2: Score Trend + Doughnut + Competency Points */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <div className="lg:col-span-2">
          <ScoreTrendChart title="Score Trend" data={trendData} />
        </div>

        <PerformanceDoughnut title="Overall Performance" percentage={percentageScore} />

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Competency Points</CardTitle>
          </CardHeader>
          <CardContent>
            {categoryPoints.length === 0 ? (
              <p className="text-sm text-muted-foreground">No competency data</p>
            ) : (
              <div className="space-y-3">
                {categoryPoints.map((cp) => (
                  <div key={cp.name} className="flex items-center justify-between">
                    <span className="text-sm">{cp.name}</span>
                    <span className="text-sm font-semibold">{cp.value.toFixed(2)}</span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
