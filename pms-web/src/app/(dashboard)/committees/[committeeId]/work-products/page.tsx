"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Briefcase } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getCommitteeDetails, getCommitteeWorkProducts } from "@/lib/api/pms-engine";
import type { Committee, WorkProduct } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function CommitteeWorkProductsPage() {
  const { committeeId } = useParams<{ committeeId: string }>();
  const [committee, setCommittee] = useState<Committee | null>(null);
  const [workProducts, setWorkProducts] = useState<WorkProduct[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!committeeId) return;
    setLoading(true);
    Promise.all([getCommitteeDetails(committeeId), getCommitteeWorkProducts(committeeId)])
      .then(([commRes, wpRes]) => {
        if (commRes?.data) setCommittee(commRes.data);
        if (wpRes?.data) setWorkProducts(Array.isArray(wpRes.data) ? wpRes.data : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [committeeId]);

  const columns: ColumnDef<WorkProduct>[] = [
    { accessorKey: "name", header: "Work Product", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "deliverables", header: "Deliverables", cell: ({ row }) => <span className="line-clamp-1">{row.original.deliverables}</span> },
    { accessorKey: "staffId", header: "Assigned To" },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "maxPoint", header: "Max Pts" },
    { accessorKey: "finalScore", header: "Score", cell: ({ row }) => <Badge variant={row.original.finalScore > 0 ? "default" : "secondary"}>{row.original.finalScore}</Badge> },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
  ];

  if (loading) return <div><PageHeader title="Committee Work Products" breadcrumbs={[{ label: "Committees" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={`Work Products: ${committee?.name ?? ""}`} description={committee?.description} breadcrumbs={[{ label: "Committees", href: "/committees" }, { label: committee?.name ?? "Committee", href: `/committees/${committeeId}` }, { label: "Work Products" }]} />
      {workProducts.length > 0 ? (
        <DataTable columns={columns} data={workProducts} searchKey="name" searchPlaceholder="Search work products..." />
      ) : (
        <EmptyState icon={Briefcase} title="No Work Products" description="No work products assigned to this committee." />
      )}
    </div>
  );
}
