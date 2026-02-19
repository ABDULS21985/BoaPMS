"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter, useParams } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { Eye, BarChart3, ArrowDown } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { EmptyState } from "@/components/shared/empty-state";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { Progress } from "@/components/shared/progress-bar";
import {
  getEmployeeDetail,
  getSubordinatesScoreCard,
  getOrganogramPerformanceSummaryList,
} from "@/lib/api/dashboard";
import { getGradeInfo } from "@/lib/scorecard-helpers";
import type {
  EmployeeErpDetails,
  SubordinateScoreCard,
  OrganogramPerformanceSummary,
} from "@/types/dashboard";

function getProgressColor(score: number): string {
  if (score >= 90) return "bg-purple-500";
  if (score >= 80) return "bg-blue-500";
  if (score >= 66) return "bg-green-500";
  if (score >= 50) return "bg-orange-500";
  return "bg-red-500";
}

export default function TeamPerformancePage() {
  const { data: session } = useSession();
  const router = useRouter();
  const params = useParams<{ reviewPeriodId: string; managerId: string }>();
  const reviewPeriodId = params.reviewPeriodId;
  const managerId = params.managerId;
  const loggedInRoles = session?.user?.roles ?? [];

  const [loading, setLoading] = useState(true);
  const [manager, setManager] = useState<EmployeeErpDetails | null>(null);
  const [staffScoreCards, setStaffScoreCards] = useState<SubordinateScoreCard[]>([]);
  const [unitSummaries, setUnitSummaries] = useState<OrganogramPerformanceSummary[]>([]);

  // Role flags
  const [isHrReportAdmin, setIsHrReportAdmin] = useState(false);
  const [isHeadOfDepartment, setIsHeadOfDepartment] = useState(false);
  const [isHeadOfDivision, setIsHeadOfDivision] = useState(false);
  const [isHeadOfOffice, setIsHeadOfOffice] = useState(false);
  const [unitTabLabel, setUnitTabLabel] = useState("");

  useEffect(() => {
    if (!reviewPeriodId || !managerId) return;

    const load = async () => {
      setLoading(true);
      try {
        const [empRes, subsRes] = await Promise.allSettled([
          getEmployeeDetail(managerId),
          getSubordinatesScoreCard(managerId, reviewPeriodId),
        ]);

        if (empRes.status === "fulfilled" && empRes.value?.data) {
          const emp = empRes.value.data;
          setManager(emp);

          const hrAdmin = loggedInRoles.includes("HrReportAdmin");
          setIsHrReportAdmin(hrAdmin);

          let level: number | undefined;
          let tabLabel = "";

          if (hrAdmin) {
            level = 1; // Department level
            tabLabel = "Departments";
          } else if (emp.headOfDepartmentId === managerId) {
            setIsHeadOfDepartment(true);
            level = 2; // Division level
            tabLabel = "Divisions";
          } else if (emp.headOfDivisionId === managerId) {
            setIsHeadOfDivision(true);
            level = 3; // Office level
            tabLabel = "Offices";
          } else if (emp.headOfOfficeId === managerId) {
            setIsHeadOfOffice(true);
          }

          if (level && tabLabel) {
            setUnitTabLabel(tabLabel);
            try {
              const unitRes = await getOrganogramPerformanceSummaryList(
                managerId,
                reviewPeriodId,
                level
              );
              if (unitRes?.data && Array.isArray(unitRes.data)) {
                setUnitSummaries(unitRes.data);
              }
            } catch {
              /* ignore */
            }
          }
        }

        if (subsRes.status === "fulfilled" && subsRes.value?.scoreCards) {
          setStaffScoreCards(subsRes.value.scoreCards);
        }
      } catch {
        /* ignore */
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [reviewPeriodId, managerId, loggedInRoles]);

  const staffColumns: ColumnDef<SubordinateScoreCard>[] = [
    {
      id: "index",
      header: "#",
      cell: ({ row }) => row.index + 1,
    },
    { accessorKey: "staffName", header: "Staff Name" },
    { accessorKey: "totalWorkProducts", header: "Work Products" },
    {
      accessorKey: "percentageWorkProductsCompletion",
      header: "Completion Rate",
      cell: ({ row }) => {
        const val = row.original.percentageWorkProductsCompletion;
        return (
          <div className="flex items-center gap-2">
            <Progress
              value={val}
              className={`h-2 w-24 [&>div]:${getProgressColor(val)}`}
            />
            <span className="text-xs">{val.toFixed(1)}%</span>
          </div>
        );
      },
    },
    {
      accessorKey: "percentageScore",
      header: "Score",
      cell: ({ row }) => {
        const val = row.original.percentageScore;
        return (
          <div className="flex items-center gap-2">
            <Progress
              value={val}
              className={`h-2 w-24 [&>div]:${getProgressColor(val)}`}
            />
            <span className="text-xs">{val.toFixed(1)}%</span>
          </div>
        );
      },
    },
    {
      accessorKey: "staffPerformanceGrade",
      header: "Grade",
      cell: ({ row }) => {
        const info = getGradeInfo(row.original.staffPerformanceGrade);
        return <Badge className={info.bgClass}>{info.label}</Badge>;
      },
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => {
        const sc = row.original;
        return (
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() =>
                router.push(`/staff-scorecard/${sc.staffId}/${reviewPeriodId}`)
              }
            >
              <Eye className="mr-1 h-3.5 w-3.5" />
              ScoreCard
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                const y = new Date().getFullYear();
                router.push(`/staff-annual-scorecard/${sc.staffId}/${y}`);
              }}
            >
              <BarChart3 className="mr-1 h-3.5 w-3.5" />
              Annual
            </Button>
          </div>
        );
      },
    },
  ];

  const getOrgUnitLevel = (): number => {
    if (isHrReportAdmin) return 2;
    if (isHeadOfDepartment) return 3;
    if (isHeadOfDivision) return 4;
    return 0;
  };

  const unitColumns: ColumnDef<OrganogramPerformanceSummary>[] = [
    { accessorKey: "referenceName", header: unitTabLabel.slice(0, -1) },
    { accessorKey: "managerId", header: "Head" },
    { accessorKey: "totalWorkProducts", header: "WPs" },
    {
      accessorKey: "performanceScore",
      header: "Performance %",
      cell: ({ row }) => `${row.original.performanceScore.toFixed(1)}%`,
    },
    {
      accessorKey: "actualScore",
      header: "Score",
      cell: ({ row }) => `${row.original.actualScore.toFixed(1)} / ${row.original.maxPoint}`,
    },
    {
      accessorKey: "earnedPerformanceGrade",
      header: "Grade",
      cell: ({ row }) => {
        const info = getGradeInfo(row.original.earnedPerformanceGrade);
        return <Badge className={info.bgClass}>{info.label}</Badge>;
      },
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => {
        const item = row.original;
        const level = getOrgUnitLevel();
        return (
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() =>
                router.push(`/team-performance/${reviewPeriodId}/${item.managerId}`)
              }
            >
              <ArrowDown className="mr-1 h-3.5 w-3.5" />
              Drill Down
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() =>
                router.push(`/unit-scorecard/${item.referenceId}/${reviewPeriodId}/${level}`)
              }
            >
              <Eye className="mr-1 h-3.5 w-3.5" />
              Unit ScoreCard
            </Button>
          </div>
        );
      },
    },
  ];

  if (loading) return <PageSkeleton />;

  const showUnitTab =
    isHrReportAdmin || isHeadOfDepartment || isHeadOfDivision;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Team Performance"
        description={
          manager
            ? `${manager.firstName} ${manager.lastName} — ${manager.departmentName ?? ""}`
            : undefined
        }
        breadcrumbs={[
          { label: "Performance", href: "/performance-score-card" },
          { label: "Team Performance" },
        ]}
      />

      <Tabs defaultValue="staff">
        <TabsList>
          <TabsTrigger value="staff">Staff Performance</TabsTrigger>
          {showUnitTab && (
            <TabsTrigger value="units">{unitTabLabel}</TabsTrigger>
          )}
        </TabsList>

        <TabsContent value="staff" className="mt-4">
          {staffScoreCards.length === 0 ? (
            <EmptyState
              title="No staff data"
              description="No subordinate performance data found for this period."
            />
          ) : (
            <DataTable
              columns={staffColumns}
              data={staffScoreCards}
              searchKey="staffName"
              searchPlaceholder="Search staff..."
            />
          )}
        </TabsContent>

        {showUnitTab && (
          <TabsContent value="units" className="mt-4">
            {unitSummaries.length === 0 ? (
              <EmptyState
                title={`No ${unitTabLabel.toLowerCase()} data`}
                description={`No ${unitTabLabel.toLowerCase()} performance data found for this period.`}
              />
            ) : (
              <DataTable
                columns={unitColumns}
                data={unitSummaries}
                searchKey="referenceName"
                searchPlaceholder={`Search ${unitTabLabel.toLowerCase()}...`}
              />
            )}
          </TabsContent>
        )}
      </Tabs>
    </div>
  );
}
