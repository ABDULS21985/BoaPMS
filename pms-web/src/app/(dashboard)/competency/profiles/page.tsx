"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { Eye, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getCompetencyReviewProfiles, getCompetencyReviewPeriods } from "@/lib/api/competency";
import type { CompetencyReviewProfile, CompetencyReviewPeriod } from "@/types/competency";

export default function CompetencyProfilesPage() {
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

  const filtered = selectedPeriod === "all" ? profiles : profiles.filter((p) => p.reviewPeriodId === Number(selectedPeriod));

  const columns: ColumnDef<CompetencyReviewProfile>[] = [
    { accessorKey: "competencyName", header: "Competency" },
    { accessorKey: "competencyCategoryName", header: "Category" },
    { accessorKey: "reviewPeriodName", header: "Review Period" },
    { accessorKey: "expectedRatingName", header: "Expected", cell: ({ row }) => <Badge variant="outline">{row.original.expectedRatingName ?? row.original.expectedRatingValue}</Badge> },
    { accessorKey: "averageRatingName", header: "Actual", cell: ({ row }) => <Badge variant={row.original.averageRatingValue >= row.original.expectedRatingValue ? "default" : "destructive"}>{row.original.averageRatingName ?? row.original.averageRatingValue}</Badge> },
    { accessorKey: "averageScore", header: "Score", cell: ({ row }) => row.original.averageScore.toFixed(1) },
    { accessorKey: "numberOfDevelopmentPlans", header: "Dev Plans" },
    {
      id: "actions", header: "", cell: ({ row }) => (
        <Button size="sm" variant="outline" onClick={() => router.push(`/competency/profiles/${row.original.employeeNumber}?profileId=${row.original.competencyReviewProfileId}`)}>
          <Eye className="mr-1 h-3.5 w-3.5" />Details
        </Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="My Competency Profiles" breadcrumbs={[{ label: "Competency" }, { label: "Profiles" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="My Competency Profiles" description="View your competency review profiles and scores"
        breadcrumbs={[{ label: "Competency" }, { label: "Profiles" }]}
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
        <EmptyState icon={User} title="No Profiles" description="No competency profiles found for the selected period." />
      )}
    </div>
  );
}
