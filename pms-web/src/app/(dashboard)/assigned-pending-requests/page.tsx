"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { ClipboardCheck, Eye, Loader2, Check, CornerDownLeft, X } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { FormDialog } from "@/components/shared/form-dialog";
import {
  getPendingFeedbackActions,
  getRequestDetails,
  treatAssignedRequest,
} from "@/lib/api/pms-engine";
import { getFeedbackRequestTypeLabel } from "@/lib/feedback-helpers";
import { OperationType } from "@/types/enums";
import type { StaffPendingRequest, FeedbackRequestLog } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function AssignedPendingRequestsPage() {
  const { data: session } = useSession();
  const staffId = session?.user?.id ?? "";

  const [requests, setRequests] = useState<StaffPendingRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [detailsLoading, setDetailsLoading] = useState(false);
  const [selectedRequest, setSelectedRequest] = useState<FeedbackRequestLog | null>(null);
  const [selectedPending, setSelectedPending] = useState<StaffPendingRequest | null>(null);
  const [comment, setComment] = useState("");
  const [processing, setProcessing] = useState(false);

  const loadData = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const res = await getPendingFeedbackActions(staffId);
      if (res?.data) setRequests(Array.isArray(res.data) ? res.data : []);
    } catch { /* silent */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [staffId]);

  const openDetails = async (request: StaffPendingRequest) => {
    setSelectedPending(request);
    setComment("");
    setDetailsOpen(true);
    setDetailsLoading(true);
    try {
      const res = await getRequestDetails(request.feedbackRequestLogId);
      setSelectedRequest(res?.data ?? null);
    } catch { setSelectedRequest(null); } finally { setDetailsLoading(false); }
  };

  const handleAction = async (operationType: OperationType) => {
    if (!selectedPending) return;
    if ((operationType === OperationType.Return || operationType === OperationType.Reject) && !comment.trim()) {
      toast.error("Please provide a comment.");
      return;
    }
    setProcessing(true);
    try {
      const res = await treatAssignedRequest({
        requestId: selectedPending.feedbackRequestLogId,
        operationType,
        comment: comment.trim() || undefined,
      });
      if (res?.isSuccess) {
        toast.success(
          operationType === OperationType.Approve ? "Request approved." :
          operationType === OperationType.Return ? "Request returned." :
          "Request rejected."
        );
        setDetailsOpen(false);
        loadData();
      } else {
        toast.error(res?.message || "Action failed.");
      }
    } catch { toast.error("An error occurred."); } finally { setProcessing(false); }
  };

  const columns: ColumnDef<StaffPendingRequest>[] = [
    { accessorKey: "requestOwnerStaffId", header: "Requestor ID" },
    { accessorKey: "feedbackRequestType", header: "Request Type", cell: ({ row }) => <Badge variant="secondary">{getFeedbackRequestTypeLabel(row.original.feedbackRequestType)}</Badge> },
    { accessorKey: "referenceId", header: "Reference", cell: ({ row }) => <span className="truncate max-w-[100px] inline-block font-mono text-xs">{row.original.referenceId}</span> },
    { accessorKey: "timeInitiated", header: "Initiated", cell: ({ row }) => row.original.timeInitiated?.split("T")[0] ?? "—" },
    { accessorKey: "hasSla", header: "SLA", cell: ({ row }) => <Badge variant={row.original.hasSla ? "outline" : "secondary"}>{row.original.hasSla ? "Yes" : "No"}</Badge> },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <Button size="sm" variant="outline" onClick={() => openDetails(row.original)}>
          <Eye className="mr-1 h-3.5 w-3.5" />Review
        </Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Pending Requests" breadcrumbs={[{ label: "Requests" }, { label: "Pending" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Pending Requests"
        description="Review and take action on requests assigned to you"
        breadcrumbs={[{ label: "Requests" }, { label: "Pending Requests" }]}
      />

      {requests.length > 0 ? (
        <DataTable columns={columns} data={requests} searchKey="referenceId" searchPlaceholder="Search by reference..." />
      ) : (
        <EmptyState icon={ClipboardCheck} title="No Pending Requests" description="You have no pending requests to review." />
      )}

      <FormDialog open={detailsOpen} onOpenChange={setDetailsOpen} title="Review Request" className="sm:max-w-lg">
        {detailsLoading ? (
          <div className="flex items-center justify-center py-8"><Loader2 className="h-6 w-6 animate-spin" /></div>
        ) : selectedRequest ? (
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-3 text-sm">
              <div><span className="text-muted-foreground">Type:</span> <span className="font-medium">{getFeedbackRequestTypeLabel(selectedRequest.feedbackRequestType)}</span></div>
              <div><span className="text-muted-foreground">Reference:</span> <span className="font-medium font-mono text-xs">{selectedRequest.referenceId}</span></div>
              <div><span className="text-muted-foreground">Requestor:</span> <span className="font-medium">{selectedRequest.requestOwnerStaffName ?? selectedRequest.requestOwnerStaffId}</span></div>
              <div><span className="text-muted-foreground">Initiated:</span> <span className="font-medium">{selectedRequest.timeInitiated?.split("T")[0]}</span></div>
              <div><span className="text-muted-foreground">SLA Bound:</span> <span className="font-medium">{selectedRequest.hasSla ? "Yes" : "No"}</span></div>
            </div>

            {selectedRequest.requestOwnerComment && (
              <div className="text-sm">
                <span className="text-muted-foreground">Requestor Comment:</span>
                <p className="mt-1 rounded-md bg-muted p-2">{selectedRequest.requestOwnerComment}</p>
              </div>
            )}

            <div className="space-y-2">
              <Label>Comment (required for Return/Reject)</Label>
              <Textarea value={comment} onChange={(e) => setComment(e.target.value)} placeholder="Enter your comment..." rows={3} />
            </div>

            <div className="flex gap-2 pt-2">
              <Button className="flex-1" onClick={() => handleAction(OperationType.Approve)} disabled={processing}>
                {processing ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Check className="mr-2 h-4 w-4" />}
                Approve
              </Button>
              <Button variant="outline" className="flex-1" onClick={() => handleAction(OperationType.Return)} disabled={processing}>
                <CornerDownLeft className="mr-2 h-4 w-4" />Return
              </Button>
              <Button variant="destructive" className="flex-1" onClick={() => handleAction(OperationType.Reject)} disabled={processing}>
                <X className="mr-2 h-4 w-4" />Reject
              </Button>
            </div>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No details available.</p>
        )}
      </FormDialog>
    </div>
  );
}
