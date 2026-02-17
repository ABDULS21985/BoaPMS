"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Eye, Briefcase } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getAllStaffWorkProducts } from "@/lib/api/pms-engine";
import { getReviewPeriods } from "@/lib/api/review-periods";
import type { WorkProduct, PerformanceReviewPeriod } from "@/types/performance";
import { WorkProductType, Status } from "@/types/enums";
import type { ColumnDef } from "@tanstack/react-table";

const wpTypeLabels: Record<number, string> = {
  [WorkProductType.Operational]: "Operational",
  [WorkProductType.Project]: "Project",
  [WorkProductType.Committee]: "Committee",
};

export default function WorkProductsPage() {
  const router = useRouter();
  const [workProducts, setWorkProducts] = useState<WorkProduct[]>([]);
  const [periods, setPeriods] = useState<PerformanceReviewPeriod[]>([]);
  const [loading, setLoading] = useState(true);
  const [typeFilter, setTypeFilter] = useState<string>("all");
  const [statusFilter, setStatusFilter] = useState<string>("all");

  useEffect(() => {
    setLoading(true);
    Promise.all([getAllStaffWorkProducts(), getReviewPeriods()])
      .then(([wpRes, pRes]) => {
        if (wpRes?.data) setWorkProducts(Array.isArray(wpRes.data) ? wpRes.data : []);
        if (pRes?.data) setPeriods(Array.isArray(pRes.data) ? pRes.data : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const filtered = workProducts.filter((wp) => {
    if (typeFilter !== "all" && wp.workProductType !== Number(typeFilter)) return false;
    if (statusFilter !== "all" && wp.recordStatus !== Number(statusFilter)) return false;
    return true;
  });

  const columns: ColumnDef<WorkProduct>[] = [
    { accessorKey: "name", header: "Work Product", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "workProductType", header: "Type", cell: ({ row }) => wpTypeLabels[row.original.workProductType] ?? "—" },
    { accessorKey: "deliverables", header: "Deliverables", cell: ({ row }) => <span className="line-clamp-1">{row.original.deliverables}</span> },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "maxPoint", header: "Max Pts" },
    { accessorKey: "finalScore", header: "Score" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="outline" onClick={() => router.push(`/work-products/${row.original.workProductId}`)}><Eye className="mr-1 h-3.5 w-3.5" />View</Button> },
  ];

  if (loading) return <div><PageHeader title="Work Products" breadcrumbs={[{ label: "Work Products" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Work Products" description="View all work products across the organization" breadcrumbs={[{ label: "Work Products" }]} />

      <div className="flex flex-wrap gap-3">
        <Select value={typeFilter} onValueChange={setTypeFilter}>
          <SelectTrigger className="w-[180px]"><SelectValue placeholder="All Types" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Types</SelectItem>
            <SelectItem value={String(WorkProductType.Operational)}>Operational</SelectItem>
            <SelectItem value={String(WorkProductType.Project)}>Project</SelectItem>
            <SelectItem value={String(WorkProductType.Committee)}>Committee</SelectItem>
          </SelectContent>
        </Select>
        <Select value={statusFilter} onValueChange={setStatusFilter}>
          <SelectTrigger className="w-[180px]"><SelectValue placeholder="All Statuses" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Statuses</SelectItem>
            <SelectItem value={String(Status.Draft)}>Draft</SelectItem>
            <SelectItem value={String(Status.PendingApproval)}>Pending Approval</SelectItem>
            <SelectItem value={String(Status.Active)}>Active</SelectItem>
            <SelectItem value={String(Status.Completed)}>Completed</SelectItem>
            <SelectItem value={String(Status.Closed)}>Closed</SelectItem>
            <SelectItem value={String(Status.Paused)}>Paused</SelectItem>
            <SelectItem value={String(Status.Cancelled)}>Cancelled</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {filtered.length > 0 ? (
        <DataTable columns={columns} data={filtered} searchKey="name" searchPlaceholder="Search work products..." />
      ) : (
        <EmptyState icon={Briefcase} title="No Work Products" description="No work products match the selected filters." />
      )}
    </div>
  );
}
