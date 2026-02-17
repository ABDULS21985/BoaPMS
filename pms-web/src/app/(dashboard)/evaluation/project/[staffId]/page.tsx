"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { FolderKanban } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getStaffProjectWorkProducts } from "@/lib/api/pms-engine";
import { getStaffActiveReviewPeriod } from "@/lib/api/review-periods";
import { getStaffDetails } from "@/lib/api/staff";
import type { WorkProduct, PerformanceReviewPeriod } from "@/types/performance";
import type { Staff } from "@/types/staff";
import type { ColumnDef } from "@tanstack/react-table";

export default function ProjectWpEvaluationPage() {
  const { staffId } = useParams<{ staffId: string }>();
  const [workProducts, setWorkProducts] = useState<WorkProduct[]>([]);
  const [period, setPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [staff, setStaff] = useState<Staff | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!staffId) return;
    setLoading(true);
    Promise.all([getStaffDetails(staffId), getStaffActiveReviewPeriod()])
      .then(async ([staffRes, periodRes]) => {
        if (staffRes?.data) setStaff(staffRes.data);
        if (periodRes?.data) {
          setPeriod(periodRes.data);
          const wpRes = await getStaffProjectWorkProducts(staffId, periodRes.data.periodId);
          if (wpRes?.data) setWorkProducts(Array.isArray(wpRes.data) ? wpRes.data : []);
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [staffId]);

  const staffName = staff ? `${staff.firstName} ${staff.lastName}` : staffId;

  const columns: ColumnDef<WorkProduct>[] = [
    { accessorKey: "name", header: "Work Product", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "deliverables", header: "Deliverables", cell: ({ row }) => <span className="line-clamp-1">{row.original.deliverables}</span> },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "maxPoint", header: "Max Points" },
    { accessorKey: "finalScore", header: "Score", cell: ({ row }) => <Badge variant={row.original.finalScore > 0 ? "default" : "secondary"}>{row.original.finalScore}</Badge> },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
  ];

  if (loading) return <div><PageHeader title="Project WP Evaluation" breadcrumbs={[{ label: "Evaluation" }, { label: "Project" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={`Project WP Evaluation: ${staffName}`} description={period ? `Review Period: ${period.name}` : ""} breadcrumbs={[{ label: "Evaluation", href: "/evaluation/direct-reports" }, { label: "Project WP", href: "/evaluation/direct-reports" }, { label: staffName }]} />
      {workProducts.length > 0 ? (
        <DataTable columns={columns} data={workProducts} searchKey="name" searchPlaceholder="Search work products..." />
      ) : (
        <EmptyState icon={FolderKanban} title="No Project Work Products" description="No project work products found for evaluation." />
      )}
    </div>
  );
}
