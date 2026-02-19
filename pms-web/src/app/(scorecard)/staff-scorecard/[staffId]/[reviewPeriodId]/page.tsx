"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter, useParams } from "next/navigation";
import {
  CheckSquare,
  AlertTriangle,
  Percent,
  PlusCircle,
  MinusCircle,
  CheckCircle,
  Gauge,
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatCard } from "@/components/shared/charts/stat-card";
import { StaffProfileCard } from "@/components/shared/charts/staff-profile-card";
import { PerformanceDoughnut } from "@/components/shared/charts/performance-doughnut";
import { CompetencyBarChart } from "@/components/shared/charts/competency-bar-chart";
import { getEmployeeDetail, getStaffIdMask, getStaffScoreCard } from "@/lib/api/dashboard";
import { formatPercent } from "@/lib/scorecard-helpers";
import type { EmployeeErpDetails, StaffScoreCardDetails } from "@/types/dashboard";
import type { StaffIdMask } from "@/types/staff";

export default function StaffScoreCardPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const params = useParams<{ staffId: string; reviewPeriodId: string }>();
  const staffId = params.staffId;
  const reviewPeriodId = params.reviewPeriodId;
  const loggedInStaffId = session?.user?.id ?? "";
  const userRoles = session?.user?.roles ?? [];

  const [loading, setLoading] = useState(true);
  const [employee, setEmployee] = useState<EmployeeErpDetails | null>(null);
  const [idMask, setIdMask] = useState<StaffIdMask | null>(null);
  const [scoreCard, setScoreCard] = useState<StaffScoreCardDetails | null>(null);

  useEffect(() => {
    if (!staffId || !reviewPeriodId) return;

    const load = async () => {
      setLoading(true);
      try {
        const [empRes, maskRes, scRes] = await Promise.allSettled([
          getEmployeeDetail(staffId),
          getStaffIdMask(staffId),
          getStaffScoreCard(staffId, reviewPeriodId),
        ]);

        if (empRes.status === "fulfilled" && empRes.value?.data) {
          const emp = empRes.value.data;
          setEmployee(emp);

          // Access control: self, supervisor, head of unit, or HrReportAdmin
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

        if (scRes.status === "fulfilled" && scRes.value?.scoreCard) {
          setScoreCard(scRes.value.scoreCard);
        }
      } catch {
        /* ignore */
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [staffId, reviewPeriodId, loggedInStaffId, userRoles, router]);

  if (loading) return <PageSkeleton />;
  if (!employee || !scoreCard) return <PageSkeleton />;

  const percentageScore =
    scoreCard.maxPoints > 0
      ? Math.round((100 * scoreCard.actualPoints / scoreCard.maxPoints) * 100) / 100
      : 0;

  const competencyData = (scoreCard.pmsCompetencies || []).map((c) => ({
    name: c.pmsCompetency,
    value: c.ratingScore,
  }));

  const categoryPoints = Object.entries(scoreCard.pmsCompetencyCategory || {});

  return (
    <div className="space-y-6">
      <PageHeader
        title={`${scoreCard.year} ${scoreCard.reviewPeriod} - Score Card`}
        breadcrumbs={[
          { label: "Performance", href: "/performance-score-card" },
          { label: "Score Card" },
        ]}
      />

      <StaffProfileCard
        employee={employee}
        photoUrl={idMask?.currentStaffPhoto}
      />

      {/* Stat Cards - 2 rows of 4 */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Timely Work Products"
          value={scoreCard.totalWorkProductsCompletedOnSchedule}
          icon={CheckSquare}
          iconColor="#22c55e"
        />
        <StatCard
          title="Overdue Work Products"
          value={scoreCard.totalWorkProductsBehindSchedule}
          icon={AlertTriangle}
          iconColor="#f59e0b"
        />
        <StatCard
          title="WP Completion"
          value={formatPercent(scoreCard.percentageWorkProductsCompletion)}
          icon={Percent}
          iconColor="#22c55e"
        />
        <StatCard
          title="Competency Gap Closure"
          value={formatPercent(scoreCard.percentageGapsClosure)}
          icon={Percent}
          iconColor="#3b82f6"
        />
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Accumulated Points"
          value={scoreCard.accumulatedPoints.toFixed(2)}
          icon={PlusCircle}
          iconColor="#1e3a5f"
        />
        <StatCard
          title="Deducted Points"
          value={scoreCard.deductedPoints.toFixed(2)}
          icon={MinusCircle}
          iconColor="#ef4444"
        />
        <StatCard
          title="Actual Points"
          value={scoreCard.actualPoints.toFixed(2)}
          icon={CheckCircle}
          iconColor="#22c55e"
        />
        <StatCard
          title="Grade"
          value={scoreCard.staffPerformanceGrade || "N/A"}
          icon={Gauge}
          iconColor="#3b82f6"
        />
      </div>

      {/* Charts */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <div className="lg:col-span-2">
          <CompetencyBarChart
            title="Competency Ratings"
            data={competencyData}
            height={Math.max(250, competencyData.length * 40)}
          />
        </div>

        <PerformanceDoughnut
          title="Performance Score"
          percentage={percentageScore}
        />

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Competency Points</CardTitle>
          </CardHeader>
          <CardContent>
            {categoryPoints.length === 0 ? (
              <p className="text-sm text-muted-foreground">No competency data</p>
            ) : (
              <div className="space-y-3">
                {categoryPoints.map(([category, points]) => (
                  <div key={category} className="flex items-center justify-between">
                    <span className="text-sm">{category}</span>
                    <span className="text-sm font-semibold">{points.toFixed(2)}</span>
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
