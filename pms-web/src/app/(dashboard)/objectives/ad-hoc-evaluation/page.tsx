"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { ClipboardCheck } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getStaffIndividualObjectives } from "@/lib/api/pms-engine";
import { getStaffActiveReviewPeriod } from "@/lib/api/review-periods";
import type { IndividualPlannedObjective, PerformanceReviewPeriod } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function AdHocEvaluationPage() {
  const { data: session } = useSession();
  const staffId = session?.user?.id ?? "";
  const [objectives, setObjectives] = useState<IndividualPlannedObjective[]>([]);
  const [period, setPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!staffId) return;
    setLoading(true);
    getStaffActiveReviewPeriod()
      .then(async (periodRes) => {
        if (periodRes?.data) {
          setPeriod(periodRes.data);
          const objRes = await getStaffIndividualObjectives(staffId, periodRes.data.periodId);
          if (objRes?.data) setObjectives(Array.isArray(objRes.data) ? objRes.data : []);
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [staffId]);

  const columns: ColumnDef<IndividualPlannedObjective>[] = [
    { accessorKey: "objectiveName", header: "Objective", cell: ({ row }) => <span className="font-medium">{row.original.objectiveName || row.original.title}</span> },
    { accessorKey: "categoryName", header: "Category" },
    { accessorKey: "weight", header: "Weight", cell: ({ row }) => <Badge variant="secondary">{row.original.weight}%</Badge> },
    { accessorKey: "keyPerformanceIndicator", header: "KPI" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
  ];

  if (loading) return <div><PageHeader title="Ad-Hoc Evaluation" breadcrumbs={[{ label: "Objectives" }, { label: "Ad-Hoc Evaluation" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Ad-Hoc Evaluation" description={period ? `Review Period: ${period.name}` : "No active review period"} breadcrumbs={[{ label: "Objectives", href: "/objectives/my-objectives" }, { label: "Ad-Hoc Evaluation" }]} />

      {objectives.length > 0 ? (
        <DataTable columns={columns} data={objectives} searchKey="objectiveName" searchPlaceholder="Search objectives..." />
      ) : (
        <EmptyState icon={ClipboardCheck} title="No Objectives for Evaluation" description="There are no objectives available for ad-hoc evaluation in this period." />
      )}
    </div>
  );
}
