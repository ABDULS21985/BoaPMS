"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Eye, Users2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getCommitteesAssigned } from "@/lib/api/pms-engine";
import { getStaffDetails } from "@/lib/api/staff";
import type { Committee } from "@/types/performance";
import type { Staff } from "@/types/staff";
import type { ColumnDef } from "@tanstack/react-table";

export default function StaffCommitteesPage() {
  const { staffId } = useParams<{ staffId: string }>();
  const router = useRouter();
  const [committees, setCommittees] = useState<Committee[]>([]);
  const [staff, setStaff] = useState<Staff | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!staffId) return;
    setLoading(true);
    Promise.all([getStaffDetails(staffId), getCommitteesAssigned(staffId)])
      .then(([staffRes, commRes]) => {
        if (staffRes?.data) setStaff(staffRes.data);
        if (commRes?.data) setCommittees(Array.isArray(commRes.data) ? commRes.data : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [staffId]);

  const staffName = staff ? `${staff.firstName} ${staff.lastName}` : staffId;

  const columns: ColumnDef<Committee>[] = [
    { accessorKey: "name", header: "Committee Name", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "chairperson", header: "Chairperson" },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="outline" onClick={() => router.push(`/committees/${row.original.committeeId}`)}><Eye className="mr-1 h-3.5 w-3.5" />View</Button> },
  ];

  if (loading) return <div><PageHeader title="Staff Committees" breadcrumbs={[{ label: "Committees" }, { label: "Staff" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={`${staffName}'s Committees`} breadcrumbs={[{ label: "Committees", href: "/committees" }, { label: staffName }]} />
      {committees.length > 0 ? (
        <DataTable columns={columns} data={committees} searchKey="name" searchPlaceholder="Search committees..." />
      ) : (
        <EmptyState icon={Users2} title="No Committees" description="This staff member is not assigned to any committees." />
      )}
    </div>
  );
}
