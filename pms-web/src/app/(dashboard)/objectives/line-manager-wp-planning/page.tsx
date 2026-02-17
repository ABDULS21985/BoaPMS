"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { Users, Eye, Briefcase } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getLineManagerEmployees } from "@/lib/api/pms-engine";
import type { ColumnDef } from "@tanstack/react-table";

interface StaffWpSummary {
  staffId: string;
  staffName?: string;
  employeeNumber?: string;
  departmentName?: string;
  workProductCount?: number;
  approvedCount?: number;
}

export default function LineManagerWpPlanningPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const staffId = session?.user?.id ?? "";
  const [reports, setReports] = useState<StaffWpSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!staffId) return;
    setLoading(true);
    getLineManagerEmployees(staffId, "WorkProductPlanning")
      .then((res) => { if (res?.data) setReports(Array.isArray(res.data) ? (res.data as StaffWpSummary[]) : []); })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [staffId]);

  const columns: ColumnDef<StaffWpSummary>[] = [
    { accessorKey: "staffName", header: "Staff Name", cell: ({ row }) => <span className="font-medium">{row.original.staffName ?? row.original.staffId}</span> },
    { accessorKey: "employeeNumber", header: "Employee No." },
    { accessorKey: "departmentName", header: "Department" },
    { accessorKey: "workProductCount", header: "Work Products", cell: ({ row }) => <Badge variant="secondary">{row.original.workProductCount ?? 0}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => (
      <Button size="sm" variant="outline" onClick={() => router.push(`/objectives/staff-planning/${row.original.staffId}`)}>
        <Eye className="mr-1 h-3.5 w-3.5" />Review
      </Button>
    )},
  ];

  if (loading) return <div><PageHeader title="Line Manager - WP Planning" breadcrumbs={[{ label: "Objectives" }, { label: "Manager WP Planning" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Line Manager - WP Planning" description="Review work product planning for your team members" breadcrumbs={[{ label: "Objectives", href: "/objectives/my-objectives" }, { label: "Manager WP Planning" }]} />
      {reports.length > 0 ? (
        <DataTable columns={columns} data={reports} searchKey="staffName" searchPlaceholder="Search staff..." />
      ) : (
        <EmptyState icon={Briefcase} title="No Staff" description="No staff members found for work product planning review." />
      )}
    </div>
  );
}
