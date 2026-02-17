"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { Calendar, Eye } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getReviewPeriods } from "@/lib/api/review-periods";
import { Roles } from "@/stores/auth-store";
import { Status } from "@/types/enums";
import type { PerformanceReviewPeriod } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

const excludedStatuses = [Status.Cancelled, Status.PendingApproval, Status.Rejected, Status.Returned];

export default function RequestReviewPeriodsPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const userRoles = session?.user?.roles ?? [];
  const adminRoles: string[] = [Roles.Admin, Roles.SuperAdmin, Roles.HrAdmin, Roles.HrReportAdmin, Roles.HrApprover];
  const isAdmin = userRoles.some((r: string) => adminRoles.includes(r));

  const [periods, setPeriods] = useState<PerformanceReviewPeriod[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAdmin) { router.replace("/"); return; }
    (async () => {
      setLoading(true);
      try {
        const res = await getReviewPeriods();
        if (res?.data) {
          const all = Array.isArray(res.data) ? res.data : [];
          setPeriods(all.filter((p) => !excludedStatuses.includes(p.recordStatus as Status)));
        }
      } catch { /* silent */ } finally { setLoading(false); }
    })();
  }, [isAdmin, router]);

  const columns: ColumnDef<PerformanceReviewPeriod>[] = [
    { accessorKey: "name", header: "Period Name", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "year", header: "Year" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    { accessorKey: "startDate", header: "Start Date", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End Date", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <Button size="sm" variant="outline" onClick={() => router.push(`/requests-review-period/${row.original.periodId}`)}>
          <Eye className="mr-1 h-3.5 w-3.5" />View Requests
        </Button>
      ),
    },
  ];

  if (!isAdmin) return null;
  if (loading) return <div><PageHeader title="All Requests" breadcrumbs={[{ label: "Requests" }, { label: "Review Periods" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="All Requests - Select Review Period"
        description="Select a review period to view all requests"
        breadcrumbs={[{ label: "Requests" }, { label: "Review Periods" }]}
      />
      {periods.length > 0 ? (
        <DataTable columns={columns} data={periods} searchKey="name" searchPlaceholder="Search by period name..." />
      ) : (
        <EmptyState icon={Calendar} title="No Review Periods" description="No review periods available." />
      )}
    </div>
  );
}
