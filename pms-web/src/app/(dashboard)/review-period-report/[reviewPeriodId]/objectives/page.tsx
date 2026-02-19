"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter, useParams } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { RotateCcw } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getReviewPeriodObjectives, reInstateObjective } from "@/lib/api/review-periods";
import type { IndividualPlannedObjective } from "@/types/performance";
import { Status } from "@/types/enums";
import { Roles } from "@/stores/auth-store";

const ADMIN_ROLES: string[] = [Roles.Admin, Roles.SuperAdmin, Roles.HrReportAdmin, Roles.HrAdmin];

export default function ObjectivesReportPage() {
  const { data: session, status } = useSession();
  const router = useRouter();
  const params = useParams();
  const reviewPeriodId = params.reviewPeriodId as string;
  const userRoles = session?.user?.roles ?? [];

  const [items, setItems] = useState<IndividualPlannedObjective[]>([]);
  const [loading, setLoading] = useState(true);
  const [reinstateOpen, setReinstateOpen] = useState(false);
  const [selectedId, setSelectedId] = useState<string | null>(null);

  useEffect(() => {
    if (status === "authenticated" && !userRoles.some((r) => ADMIN_ROLES.includes(r))) {
      router.push("/access-denied");
    }
  }, [status, userRoles, router]);

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getReviewPeriodObjectives(reviewPeriodId);
      if (res?.data) setItems(Array.isArray(res.data) ? (res.data as IndividualPlannedObjective[]) : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { if (reviewPeriodId) loadData(); }, [reviewPeriodId]);

  const handleReinstate = async () => {
    if (!selectedId) return;
    const res = await reInstateObjective({ objectiveId: selectedId });
    if (res?.isSuccess) { toast.success("Objective re-instated."); loadData(); }
    else toast.error(res?.message || "Failed to re-instate objective.");
  };

  const columns: ColumnDef<IndividualPlannedObjective>[] = [
    { id: "index", header: "#", cell: ({ row }) => row.index + 1 },
    { accessorKey: "individualPlannedObjectiveId", header: "ID", cell: ({ row }) => <span title={row.original.individualPlannedObjectiveId}>{row.original.individualPlannedObjectiveId?.slice(0, 8)}...</span> },
    { accessorKey: "staffId", header: "Staff ID" },
    { accessorKey: "objectiveName", header: "Objective Name", cell: ({ row }) => row.original.objectiveName || row.original.title || "-" },
    { accessorKey: "objectiveLevel", header: "Level/Category", cell: ({ row }) => row.original.categoryName || row.original.objectiveLevel || "-" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => <StatusBadge status={row.original.recordStatus ?? 0} /> },
    {
      id: "actions", header: "Actions", cell: ({ row }) => {
        const isRejected = row.original.recordStatus === Status.Rejected;
        if (!isRejected) return null;
        return (
          <Button size="sm" variant="ghost" onClick={() => { setSelectedId(row.original.individualPlannedObjectiveId); setReinstateOpen(true); }} title="Re-instate">
            <RotateCcw className="h-3.5 w-3.5" />
          </Button>
        );
      },
    },
  ];

  if (loading) return <div><PageHeader title="Staff Planned Objectives" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Reports" }, { label: "Objectives & Work Products", href: "/review-period-report" }, { label: "Objectives" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Staff Planned Objectives" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Reports" }, { label: "Objectives & Work Products", href: "/review-period-report" }, { label: "Objectives" }]} />
      <DataTable columns={columns} data={items} searchKey="staffId" searchPlaceholder="Search by staff ID or objective name..." />

      <ConfirmationDialog open={reinstateOpen} onOpenChange={setReinstateOpen} title="Re-instate Objective" description="Are you sure you want to re-instate this objective? This will restore it to its previous active status." confirmLabel="Re-instate" onConfirm={handleReinstate} />
    </div>
  );
}
