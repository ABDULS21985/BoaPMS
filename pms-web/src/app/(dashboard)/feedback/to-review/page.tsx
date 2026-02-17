"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { Star, Users } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getCompetenciesToReview } from "@/lib/api/pms-engine";
import type { CompetencyReviewer } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function FeedbackToReviewPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const staffId = session?.user?.id ?? "";

  const [toReview, setToReview] = useState<CompetencyReviewer[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!staffId) return;
    (async () => {
      setLoading(true);
      try {
        const res = await getCompetenciesToReview(staffId);
        if (res?.data) setToReview(Array.isArray(res.data) ? res.data : []);
      } catch { /* silent */ } finally { setLoading(false); }
    })();
  }, [staffId]);

  const columns: ColumnDef<CompetencyReviewer>[] = [
    { accessorKey: "competencyReviewFeedback.staffName", header: "Staff Name", cell: ({ row }) => row.original.competencyReviewFeedback?.staffName ?? "—" },
    { accessorKey: "reviewStaffId", header: "Staff ID", cell: ({ row }) => row.original.competencyReviewFeedback?.staffId ?? row.original.reviewStaffId },
    { accessorKey: "initiatedDate", header: "Initiated", cell: ({ row }) => row.original.initiatedDate?.split("T")[0] ?? "—" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : row.original.recordStatusName ?? "—" },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <Button
          size="sm"
          onClick={() => router.push(
            `/feedback/staff-rating/${row.original.competencyReviewFeedback?.staffId ?? "unknown"}/${row.original.competencyReviewerId}/${row.original.competencyReviewFeedbackId}`
          )}
        >
          <Star className="mr-1 h-3.5 w-3.5" />Rate
        </Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Feedbacks To Review" breadcrumbs={[{ label: "Feedback" }, { label: "To Review" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Feedbacks To Review"
        description="Staff members assigned to you for 360-degree feedback review"
        breadcrumbs={[{ label: "Feedback" }, { label: "To Review" }]}
      />

      {toReview.length > 0 ? (
        <DataTable columns={columns} data={toReview} searchKey="competencyReviewFeedback.staffName" searchPlaceholder="Search by staff name..." />
      ) : (
        <EmptyState icon={Users} title="No Pending Reviews" description="You have no pending 360 feedbacks to review." />
      )}
    </div>
  );
}
