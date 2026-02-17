"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { Eye, TrendingDown } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getCompetencyReviewProfiles, getCompetencyReviewPeriods } from "@/lib/api/competency";
import type { CompetencyReviewProfile, CompetencyReviewPeriod } from "@/types/competency";

export default function CompetencyGapsPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const employeeNumber = session?.user?.id ?? "";
  const [profiles, setProfiles] = useState<CompetencyReviewProfile[]>([]);
  const [periods, setPeriods] = useState<CompetencyReviewPeriod[]>([]);
  const [selectedPeriod, setSelectedPeriod] = useState<string>("all");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!employeeNumber) return;
    setLoading(true);
    Promise.all([getCompetencyReviewProfiles(employeeNumber), getCompetencyReviewPeriods()])
      .then(([profRes, periodRes]) => {
        if (profRes?.data) setProfiles(Array.isArray(profRes.data) ? profRes.data : []);
        if (periodRes?.data) setPeriods(Array.isArray(periodRes.data) ? periodRes.data : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [employeeNumber]);

  const filtered = (selectedPeriod === "all" ? profiles : profiles.filter((p) => p.reviewPeriodId === Number(selectedPeriod)))
    .filter((p) => p.averageRatingValue < p.expectedRatingValue);

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
    {
      id: "actions", header: "", cell: ({ row }) => (
        <Button size="sm" variant="outline" onClick={() => router.push(`/competency/gaps/${row.original.employeeNumber}`)}>
          <Eye className="mr-1 h-3.5 w-3.5" />View
        </Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Competency Gaps" breadcrumbs={[{ label: "Competency" }, { label: "Gaps" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Competency Gaps" description="View competencies where actual rating is below expected"
        breadcrumbs={[{ label: "Competency" }, { label: "Gaps" }]}
      />

      <div className="flex gap-3">
        <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
          <SelectTrigger className="w-[220px]"><SelectValue placeholder="All Periods" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Periods</SelectItem>
            {periods.map((p) => <SelectItem key={p.reviewPeriodId} value={String(p.reviewPeriodId)}>{p.name}</SelectItem>)}
          </SelectContent>
        </Select>
      </div>

      {filtered.length > 0 ? (
        <DataTable columns={columns} data={filtered} searchKey="competencyName" searchPlaceholder="Search competencies..." />
      ) : (
        <EmptyState icon={TrendingDown} title="No Gaps Found" description="No competency gaps for the selected period." />
      )}
    </div>
  );
}
