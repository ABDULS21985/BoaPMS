"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { TrendingDown } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getCompetencyReviewProfiles } from "@/lib/api/competency";
import type { CompetencyReviewProfile } from "@/types/competency";

export default function EmployeeGapsPage() {
  const params = useParams();
  const employeeId = params.employeeId as string;
  const [profiles, setProfiles] = useState<CompetencyReviewProfile[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getCompetencyReviewProfiles(employeeId)
      .then((res) => {
        if (res?.data) setProfiles(Array.isArray(res.data) ? res.data : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [employeeId]);

  const gaps = profiles.filter((p) => p.averageRatingValue < p.expectedRatingValue);
  const employeeName = profiles[0]?.employeeFullName ?? employeeId;
  const totalGap = gaps.reduce((sum, g) => sum + (g.expectedRatingValue - g.averageRatingValue), 0);

  const columns: ColumnDef<CompetencyReviewProfile>[] = [
    { accessorKey: "competencyName", header: "Competency" },
    { accessorKey: "competencyCategoryName", header: "Category" },
    { accessorKey: "reviewPeriodName", header: "Period" },
    { accessorKey: "expectedRatingValue", header: "Expected", cell: ({ row }) => <Badge variant="outline">{row.original.expectedRatingName ?? row.original.expectedRatingValue}</Badge> },
    { accessorKey: "averageRatingValue", header: "Actual", cell: ({ row }) => <Badge variant="destructive">{row.original.averageRatingName ?? row.original.averageRatingValue}</Badge> },
    {
      id: "gap", header: "Gap",
      cell: ({ row }) => {
        const gap = row.original.expectedRatingValue - row.original.averageRatingValue;
        return <Badge variant="destructive">{gap.toFixed(1)}</Badge>;
      },
    },
    { accessorKey: "numberOfDevelopmentPlans", header: "Dev Plans" },
  ];

  if (loading) return <div><PageHeader title="Employee Gaps" breadcrumbs={[{ label: "Competency" }, { label: "Gaps" }, { label: "Details" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={`${employeeName} - Competency Gaps`} description="Competencies below expected rating"
        breadcrumbs={[{ label: "Competency" }, { label: "Gaps", href: "/competency/gaps" }, { label: employeeName }]}
      />

      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Total Competencies</CardTitle></CardHeader><CardContent><p className="text-2xl font-bold">{profiles.length}</p></CardContent></Card>
        <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Gaps Identified</CardTitle></CardHeader><CardContent><p className="text-2xl font-bold text-red-600">{gaps.length}</p></CardContent></Card>
        <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Total Gap Score</CardTitle></CardHeader><CardContent><p className="text-2xl font-bold text-red-600">{totalGap.toFixed(1)}</p></CardContent></Card>
      </div>

      {gaps.length > 0 ? (
        <DataTable columns={columns} data={gaps} searchKey="competencyName" searchPlaceholder="Search competencies..." />
      ) : (
        <EmptyState icon={TrendingDown} title="No Gaps" description="All competencies meet or exceed expected ratings." />
      )}
    </div>
  );
}
