"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { Eye, BarChart3 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getEmployeeDetail, getHeadSubordinates, getStaffReviewPeriods } from "@/lib/api/dashboard";
import { getReviewPeriods } from "@/lib/api/review-periods";
import { getStaffObjectives } from "@/lib/api/pms-engine";
import type { PerformanceReviewPeriod } from "@/types/performance";
import { Status } from "@/types/enums";

export default function PerformanceScoreCardPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const staffId = session?.user?.id ?? "";

  const [loading, setLoading] = useState(true);
  const [isHeadOfUnit, setIsHeadOfUnit] = useState(false);
  const [myPeriods, setMyPeriods] = useState<PerformanceReviewPeriod[]>([]);
  const [allPeriods, setAllPeriods] = useState<PerformanceReviewPeriod[]>([]);

  useEffect(() => {
    if (!staffId) return;

    const load = async () => {
      setLoading(true);
      try {
        const [empRes, subsRes, periodsRes] = await Promise.allSettled([
          getEmployeeDetail(staffId),
          getHeadSubordinates(staffId),
          getReviewPeriods(),
        ]);

        if (empRes.status === "fulfilled" && !empRes.value?.data) {
          router.replace("/404");
          return;
        }

        if (subsRes.status === "fulfilled" && subsRes.value?.data) {
          setIsHeadOfUnit(Array.isArray(subsRes.value.data) && subsRes.value.data.length > 0);
        }

        const periods =
          periodsRes.status === "fulfilled" && periodsRes.value?.data
            ? Array.isArray(periodsRes.value.data) ? periodsRes.value.data : []
            : [];
        setAllPeriods(periods);

        // Filter periods where staff has planned objectives (exclude certain statuses)
        const excludedStatuses = new Set([
          Status.Cancelled,
          Status.PendingApproval,
          Status.Rejected,
          Status.Returned,
        ]);

        const periodObjectiveResults = await Promise.allSettled(
          periods.map((p) => getStaffObjectives(staffId, p.periodId))
        );

        const filteredPeriods = periods.filter((p, idx) => {
          const res = periodObjectiveResults[idx];
          if (res.status !== "fulfilled") return false;
          const objectives = res.value?.data;
          if (!Array.isArray(objectives) || objectives.length === 0) return false;
          return objectives.some(
            (o) => o.recordStatus != null && !excludedStatuses.has(o.recordStatus)
          );
        });

        if (filteredPeriods.length > 0) {
          setMyPeriods(filteredPeriods);
        } else {
          // Fallback: get staff review periods
          try {
            const fallback = await getStaffReviewPeriods(staffId);
            if (fallback?.data && Array.isArray(fallback.data)) {
              setMyPeriods(fallback.data);
            }
          } catch {
            /* ignore */
          }
        }
      } catch {
        /* ignore */
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [staffId, router]);

  const periodColumns: ColumnDef<PerformanceReviewPeriod>[] = [
    { accessorKey: "name", header: "Name" },
    { accessorKey: "year", header: "Year" },
    {
      accessorKey: "recordStatus",
      header: "Status",
      cell: ({ row }) => <StatusBadge status={row.original.recordStatus ?? 0} />,
    },
    {
      accessorKey: "startDate",
      header: "Start Date",
      cell: ({ row }) => new Date(row.original.startDate).toLocaleDateString(),
    },
    {
      accessorKey: "endDate",
      header: "End Date",
      cell: ({ row }) => new Date(row.original.endDate).toLocaleDateString(),
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => (
        <div className="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() =>
              router.push(`/staff-scorecard/${staffId}/${row.original.periodId}`)
            }
          >
            <Eye className="mr-1 h-3.5 w-3.5" />
            ScoreCard
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() =>
              router.push(`/staff-annual-scorecard/${staffId}/${row.original.year}`)
            }
          >
            <BarChart3 className="mr-1 h-3.5 w-3.5" />
            Annual
          </Button>
        </div>
      ),
    },
  ];

  const teamColumns: ColumnDef<PerformanceReviewPeriod>[] = [
    { accessorKey: "name", header: "Name" },
    { accessorKey: "year", header: "Year" },
    {
      accessorKey: "recordStatus",
      header: "Status",
      cell: ({ row }) => <StatusBadge status={row.original.recordStatus ?? 0} />,
    },
    {
      accessorKey: "startDate",
      header: "Start Date",
      cell: ({ row }) => new Date(row.original.startDate).toLocaleDateString(),
    },
    {
      accessorKey: "endDate",
      header: "End Date",
      cell: ({ row }) => new Date(row.original.endDate).toLocaleDateString(),
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => (
        <Button
          variant="outline"
          size="sm"
          onClick={() =>
            router.push(`/team-performance/${row.original.periodId}/${staffId}`)
          }
        >
          <Eye className="mr-1 h-3.5 w-3.5" />
          View Team
        </Button>
      ),
    },
  ];

  if (loading) return <PageSkeleton />;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Performance Score Card"
        description="View your performance score cards and team performance"
        breadcrumbs={[{ label: "Performance Score Card" }]}
      />

      <Tabs defaultValue="my-performance">
        <TabsList>
          <TabsTrigger value="my-performance">My Performance</TabsTrigger>
          {isHeadOfUnit && (
            <TabsTrigger value="team-performance">Team Performance</TabsTrigger>
          )}
        </TabsList>

        <TabsContent value="my-performance" className="mt-4">
          {myPeriods.length === 0 ? (
            <EmptyState
              title="No review periods found"
              description="You have no review periods with planned objectives."
            />
          ) : (
            <DataTable
              columns={periodColumns}
              data={myPeriods}
              searchKey="name"
              searchPlaceholder="Search review periods..."
            />
          )}
        </TabsContent>

        {isHeadOfUnit && (
          <TabsContent value="team-performance" className="mt-4">
            {allPeriods.length === 0 ? (
              <EmptyState
                title="No review periods found"
                description="No review periods are available."
              />
            ) : (
              <DataTable
                columns={teamColumns}
                data={allPeriods}
                searchKey="name"
                searchPlaceholder="Search review periods..."
              />
            )}
          </TabsContent>
        )}
      </Tabs>
    </div>
  );
}
