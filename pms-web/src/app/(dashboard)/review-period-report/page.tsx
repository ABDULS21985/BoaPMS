"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { format } from "date-fns";
import { FileText, Package } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getReviewPeriods } from "@/lib/api/review-periods";
import type { PerformanceReviewPeriod } from "@/types/performance";
import { Roles } from "@/stores/auth-store";

const ADMIN_ROLES: string[] = [Roles.Admin, Roles.SuperAdmin, Roles.HrReportAdmin, Roles.HrAdmin];

export default function ReviewPeriodReportPage() {
  const { data: session, status } = useSession();
  const router = useRouter();
  const userRoles = session?.user?.roles ?? [];

  const [periods, setPeriods] = useState<PerformanceReviewPeriod[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (status === "authenticated" && !userRoles.some((r) => ADMIN_ROLES.includes(r))) {
      router.push("/access-denied");
    }
  }, [status, userRoles, router]);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      try {
        const res = await getReviewPeriods();
        if (res?.data) setPeriods(Array.isArray(res.data) ? res.data : []);
      } catch { /* */ } finally { setLoading(false); }
    };
    load();
  }, []);

  const columns: ColumnDef<PerformanceReviewPeriod>[] = [
    { id: "index", header: "#", cell: ({ row }) => row.index + 1 },
    { accessorKey: "name", header: "Name" },
    { accessorKey: "year", header: "Year" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => <StatusBadge status={row.original.recordStatus ?? 0} /> },
    { accessorKey: "startDate", header: "Start Date", cell: ({ row }) => row.original.startDate ? format(new Date(row.original.startDate), "dd MMM yyyy") : "-" },
    { accessorKey: "endDate", header: "End Date", cell: ({ row }) => row.original.endDate ? format(new Date(row.original.endDate), "dd MMM yyyy") : "-" },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => router.push(`/review-period-report/${row.original.periodId}/objectives`)} title="View Objectives"><FileText className="h-3.5 w-3.5" /></Button>
          <Button size="sm" variant="ghost" onClick={() => router.push(`/review-period-report/${row.original.periodId}/work-products`)} title="View Work Products"><Package className="h-3.5 w-3.5" /></Button>
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Objectives & Work Products Report" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Objectives & Work Products" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Objectives & Work Products Report" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Objectives & Work Products" }]} />
      <DataTable columns={columns} data={periods} searchKey="name" searchPlaceholder="Search by period name..." />
    </div>
  );
}
