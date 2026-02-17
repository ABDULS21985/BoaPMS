"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { MessageSquare, Eye, Star, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { FormDialog } from "@/components/shared/form-dialog";
import {
  getAllCompetencyReviewFeedbacks,
  getCompetenciesToReview,
  getCompetencyReviewFeedbackDetails,
} from "@/lib/api/pms-engine";
import { getActiveReviewPeriod } from "@/lib/api/dashboard";
import type {
  CompetencyReviewFeedback,
  CompetencyReviewer,
  CompetencyReviewFeedbackDetails,
} from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function FeedbackReviewsPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const staffId = session?.user?.id ?? "";

  const [myFeedbacks, setMyFeedbacks] = useState<CompetencyReviewFeedback[]>([]);
  const [toReview, setToReview] = useState<CompetencyReviewer[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [detailsLoading, setDetailsLoading] = useState(false);
  const [selectedDetails, setSelectedDetails] = useState<CompetencyReviewFeedbackDetails | null>(null);

  useEffect(() => {
    if (!staffId) return;
    (async () => {
      setLoading(true);
      try {
        const [feedbackRes, reviewRes] = await Promise.allSettled([
          getAllCompetencyReviewFeedbacks(staffId),
          getCompetenciesToReview(staffId),
        ]);
        if (feedbackRes.status === "fulfilled" && feedbackRes.value?.data) {
          setMyFeedbacks(Array.isArray(feedbackRes.value.data) ? feedbackRes.value.data : []);
        }
        if (reviewRes.status === "fulfilled" && reviewRes.value?.data) {
          setToReview(Array.isArray(reviewRes.value.data) ? reviewRes.value.data : []);
        }
      } catch { /* silent */ } finally {
        setLoading(false);
      }
    })();
  }, [staffId]);

  const openDetails = async (feedbackId: string) => {
    setDetailsOpen(true);
    setDetailsLoading(true);
    try {
      const res = await getCompetencyReviewFeedbackDetails(feedbackId);
      setSelectedDetails(res?.data ?? null);
    } catch { setSelectedDetails(null); } finally {
      setDetailsLoading(false);
    }
  };

  const myFeedbackColumns: ColumnDef<CompetencyReviewFeedback>[] = [
    { accessorKey: "staffName", header: "Staff", cell: ({ row }) => row.original.staffName ?? row.original.staffId },
    { accessorKey: "maxPoints", header: "Max Points" },
    { accessorKey: "finalScore", header: "Final Score" },
    { accessorKey: "reviewPeriodId", header: "Review Period", cell: ({ row }) => <span className="truncate max-w-[120px] inline-block">{row.original.reviewPeriodId}</span> },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : row.original.recordStatusName ?? "—" },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <Button size="sm" variant="outline" onClick={() => openDetails(row.original.competencyReviewFeedbackId)}>
          <Eye className="mr-1 h-3.5 w-3.5" />Details
        </Button>
      ),
    },
  ];

  const toReviewColumns: ColumnDef<CompetencyReviewer>[] = [
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
            `/feedback-reviews/rating/${row.original.competencyReviewFeedback?.staffId ?? "unknown"}/${row.original.competencyReviewerId}/${row.original.competencyReviewFeedbackId}`
          )}
        >
          <Star className="mr-1 h-3.5 w-3.5" />Rate
        </Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="360 Feedback Reviews" breadcrumbs={[{ label: "360 Feedback" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="360 Feedback Reviews"
        description="Manage your 360-degree feedback reviews and rate colleagues"
        breadcrumbs={[{ label: "360 Feedback" }]}
      />

      <Tabs defaultValue="my-reviews" className="w-full">
        <TabsList>
          <TabsTrigger value="my-reviews">My 360 Reviews ({myFeedbacks.length})</TabsTrigger>
          <TabsTrigger value="to-review">Feedbacks To Review ({toReview.length})</TabsTrigger>
        </TabsList>

        <TabsContent value="my-reviews" className="mt-4">
          {myFeedbacks.length > 0 ? (
            <DataTable columns={myFeedbackColumns} data={myFeedbacks} searchKey="staffName" searchPlaceholder="Search by staff name..." />
          ) : (
            <EmptyState icon={MessageSquare} title="No Feedback Reviews" description="You have no 360 feedback reviews for the current period." />
          )}
        </TabsContent>

        <TabsContent value="to-review" className="mt-4">
          {toReview.length > 0 ? (
            <DataTable columns={toReviewColumns} data={toReview} searchKey="competencyReviewFeedback.staffName" searchPlaceholder="Search by staff name..." />
          ) : (
            <EmptyState icon={Star} title="No Pending Reviews" description="You have no pending 360 feedbacks to review." />
          )}
        </TabsContent>
      </Tabs>

      <FormDialog open={detailsOpen} onOpenChange={setDetailsOpen} title="Feedback Review Details" className="sm:max-w-lg">
        {detailsLoading ? (
          <div className="flex items-center justify-center py-8"><Loader2 className="h-6 w-6 animate-spin" /></div>
        ) : selectedDetails ? (
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-3 text-sm">
              <div><span className="text-muted-foreground">Staff:</span> <span className="font-medium">{selectedDetails.staffName ?? selectedDetails.staffId}</span></div>
              <div><span className="text-muted-foreground">Department:</span> <span className="font-medium">{selectedDetails.departmentName ?? "—"}</span></div>
              <div><span className="text-muted-foreground">Division:</span> <span className="font-medium">{selectedDetails.divisionName ?? "—"}</span></div>
              <div><span className="text-muted-foreground">Office:</span> <span className="font-medium">{selectedDetails.officeName ?? "—"}</span></div>
              <div><span className="text-muted-foreground">Max Points:</span> <span className="font-medium">{selectedDetails.maxPoints}</span></div>
              <div><span className="text-muted-foreground">Final Score:</span> <span className="font-medium">{selectedDetails.finalScore}</span></div>
              <div><span className="text-muted-foreground">Score %:</span> <span className="font-medium">{selectedDetails.finalScorePercentage?.toFixed(1)}%</span></div>
              <div><span className="text-muted-foreground">Status:</span> <span className="font-medium">{selectedDetails.recordStatusName ?? "—"}</span></div>
            </div>
            {selectedDetails.ratings && selectedDetails.ratings.length > 0 && (
              <div>
                <h4 className="text-sm font-semibold mb-2">Competency Ratings</h4>
                <div className="space-y-1">
                  {selectedDetails.ratings.map((r) => (
                    <div key={r.pmsCompetencyId} className="flex justify-between text-sm border-b pb-1">
                      <span>{r.pmsCompetencyName}</span>
                      <span className="font-medium">{r.averageRating?.toFixed(2)} ({r.totalReviewers} reviewers)</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No details available.</p>
        )}
      </FormDialog>
    </div>
  );
}
