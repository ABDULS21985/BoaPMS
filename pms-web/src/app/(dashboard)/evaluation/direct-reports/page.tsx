"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { Users, Eye, Star } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getLineManagerEmployees } from "@/lib/api/pms-engine";
import type { ColumnDef } from "@tanstack/react-table";

interface StaffEvalSummary {
  staffId: string;
  staffName?: string;
  employeeNumber?: string;
  departmentName?: string;
  objectiveCount?: number;
  evaluatedCount?: number;
}

export default function EvalDirectReportsPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const staffId = session?.user?.id ?? "";
  const [reports, setReports] = useState<StaffEvalSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!staffId) return;
    setLoading(true);
    getLineManagerEmployees(staffId, "Evaluation")
      .then((res) => { if (res?.data) setReports(Array.isArray(res.data) ? (res.data as StaffEvalSummary[]) : []); })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [staffId]);

  const columns: ColumnDef<StaffEvalSummary>[] = [
    { accessorKey: "staffName", header: "Staff Name", cell: ({ row }) => <span className="font-medium">{row.original.staffName ?? row.original.staffId}</span> },
    { accessorKey: "employeeNumber", header: "Employee No." },
    { accessorKey: "departmentName", header: "Department" },
    { accessorKey: "objectiveCount", header: "Objectives", cell: ({ row }) => <Badge variant="secondary">{row.original.objectiveCount ?? 0}</Badge> },
    { accessorKey: "evaluatedCount", header: "Evaluated", cell: ({ row }) => <Badge variant="default">{row.original.evaluatedCount ?? 0}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => (
      <Button size="sm" variant="outline" onClick={() => router.push(`/evaluation/staff/${row.original.staffId}`)}>
        <Eye className="mr-1 h-3.5 w-3.5" />Evaluate
      </Button>
    )},
  ];

  if (loading) return <div><PageHeader title="Direct Reports - Evaluation" breadcrumbs={[{ label: "Evaluation" }, { label: "Direct Reports" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Direct Reports - Evaluation" description="Evaluate objectives for your team members" breadcrumbs={[{ label: "Evaluation", href: "/evaluation/direct-reports" }, { label: "Direct Reports" }]} />

      <div className="grid grid-cols-2 gap-4 md:grid-cols-3">
        <Card><CardContent className="flex items-center gap-3 pt-4"><Users className="h-8 w-8 text-muted-foreground" /><div><div className="text-2xl font-bold">{reports.length}</div><p className="text-xs text-muted-foreground">Total Staff</p></div></CardContent></Card>
        <Card><CardContent className="flex items-center gap-3 pt-4"><Star className="h-8 w-8 text-amber-500" /><div><div className="text-2xl font-bold">{reports.filter((r) => (r.evaluatedCount ?? 0) > 0).length}</div><p className="text-xs text-muted-foreground">Evaluated</p></div></CardContent></Card>
        <Card><CardContent className="flex items-center gap-3 pt-4"><Eye className="h-8 w-8 text-blue-500" /><div><div className="text-2xl font-bold">{reports.filter((r) => (r.objectiveCount ?? 0) > (r.evaluatedCount ?? 0)).length}</div><p className="text-xs text-muted-foreground">Pending</p></div></CardContent></Card>
      </div>

      {reports.length > 0 ? (
        <DataTable columns={columns} data={reports} searchKey="staffName" searchPlaceholder="Search staff..." />
      ) : (
        <EmptyState icon={Users} title="No Direct Reports" description="You don't have any staff members to evaluate." />
      )}
    </div>
  );
}
