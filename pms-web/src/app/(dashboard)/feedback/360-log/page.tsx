"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { MessageSquare } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getMyReviewedCompetencies } from "@/lib/api/pms-engine";
import type { CompetencyReviewer } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function Feedback360LogPage() {
  const { data: session } = useSession();
  const staffId = session?.user?.id ?? "";
  const [logs, setLogs] = useState<CompetencyReviewer[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!staffId) return;
    (async () => {
      setLoading(true);
      try {
        const res = await getMyReviewedCompetencies(staffId);
        if (res?.data) setLogs(Array.isArray(res.data) ? res.data : []);
      } catch { /* silent */ } finally { setLoading(false); }
    })();
  }, [staffId]);

  const columns: ColumnDef<CompetencyReviewer>[] = [
    { accessorKey: "competencyReviewFeedback.staffName", header: "Staff Reviewed", cell: ({ row }) => row.original.competencyReviewFeedback?.staffName ?? "—" },
    { accessorKey: "competencyReviewFeedback.staffId", header: "Staff ID", cell: ({ row }) => row.original.competencyReviewFeedback?.staffId ?? "—" },
    { accessorKey: "initiatedDate", header: "Initiated", cell: ({ row }) => row.original.initiatedDate?.split("T")[0] ?? "—" },
    { accessorKey: "finalRating", header: "Final Rating", cell: ({ row }) => <Badge variant="outline">{row.original.finalRating?.toFixed(2) ?? "—"}</Badge> },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : row.original.recordStatusName ?? "—" },
  ];

  if (loading) return <div><PageHeader title="360 Feedback Log" breadcrumbs={[{ label: "Feedback" }, { label: "360 Log" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="360 Feedback Log"
        description="History of your completed 360-degree feedback reviews"
        breadcrumbs={[{ label: "Feedback" }, { label: "360 Log" }]}
      />
      {logs.length > 0 ? (
        <DataTable columns={columns} data={logs} searchKey="competencyReviewFeedback.staffName" searchPlaceholder="Search by staff name..." />
      ) : (
        <EmptyState icon={MessageSquare} title="No Feedback Logs" description="You have not completed any 360 feedback reviews yet." />
      )}
    </div>
  );
}
