"use client";

import { useEffect, useState } from "react";
import { Briefcase } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getAllStaffWorkProducts } from "@/lib/api/pms-engine";
import type { WorkProduct } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

const wpTypeLabels: Record<number, string> = { 1: "Operational", 2: "Project", 3: "Committee" };

export default function AllWpPlanningPage() {
  const [items, setItems] = useState<WorkProduct[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getAllStaffWorkProducts()
      .then((res) => { if (res?.data) setItems(Array.isArray(res.data) ? res.data : []); })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const columns: ColumnDef<WorkProduct>[] = [
    { accessorKey: "name", header: "Work Product", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "workProductType", header: "Type", cell: ({ row }) => <Badge variant="outline">{wpTypeLabels[row.original.workProductType] ?? "Unknown"}</Badge> },
    { accessorKey: "staffId", header: "Staff ID" },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "maxPoint", header: "Max Points" },
    { accessorKey: "finalScore", header: "Score" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
  ];

  if (loading) return <div><PageHeader title="All Work Product Planning" breadcrumbs={[{ label: "Objectives" }, { label: "All WP Planning" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="All Work Product Planning" description="Organization-wide view of work product planning across all staff" breadcrumbs={[{ label: "Objectives", href: "/objectives/my-objectives" }, { label: "All WP Planning" }]} />
      {items.length > 0 ? (
        <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search work products..." />
      ) : (
        <EmptyState icon={Briefcase} title="No Work Products" description="No work products found across the organization." />
      )}
    </div>
  );
}
