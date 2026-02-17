"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { Clock, Eye, CheckCircle, Undo2, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { FormDialog } from "@/components/shared/form-dialog";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import {
  getStaffFeedbackRequests,
  getRequestDetails,
  treatAssignedRequest,
} from "@/lib/api/pms-engine";
import { getFeedbackRequestTypeLabel } from "@/lib/feedback-helpers";
import { Status, OperationType } from "@/types/enums";
import type { FeedbackRequestLog } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function PendingRequestsPage() {
  const { data: session } = useSession();
  const staffId = session?.user?.id ?? "";

  const [requests, setRequests] = useState<FeedbackRequestLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [detailsLoading, setDetailsLoading] = useState(false);
  const [selectedRequest, setSelectedRequest] = useState<FeedbackRequestLog | null>(null);
  const [approveItem, setApproveItem] = useState<string | null>(null);
  const [returnItem, setReturnItem] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState(false);

  const loadData = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const res = await getStaffFeedbackRequests(staffId);
      if (res?.data) {
        const all = Array.isArray(res.data) ? res.data : [];
        // Show only active/pending requests (not completed)
        setRequests(all.filter((r) => r.recordStatus != null && r.recordStatus !== Status.Completed && r.recordStatus !== Status.Closed && r.recordStatus !== Status.Cancelled));
      }
    } catch { /* silent */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [staffId]);

  const openDetails = async (requestId: string) => {
    setDetailsOpen(true);
    setDetailsLoading(true);
    try {
      const res = await getRequestDetails(requestId);
      setSelectedRequest(res?.data ?? null);
    } catch { setSelectedRequest(null); } finally { setDetailsLoading(false); }
  };

  const handleApprove = async (reason?: string) => {
    if (!approveItem) return;
    setActionLoading(true);
    try {
      const res = await treatAssignedRequest({
        requestId: approveItem,
        operationType: OperationType.Approve,
        comment: reason ?? "",
      });
      if (res?.isSuccess) { toast.success("Request approved."); loadData(); }
      else toast.error(res?.message || "Approval failed.");
    } catch { toast.error("An error occurred."); } finally { setActionLoading(false); setApproveItem(null); }
  };

  const handleReturn = async (reason?: string) => {
    if (!returnItem) return;
    setActionLoading(true);
    try {
      const res = await treatAssignedRequest({
        requestId: returnItem,
        operationType: OperationType.Return,
        comment: reason ?? "",
      });
      if (res?.isSuccess) { toast.success("Request returned."); loadData(); }
      else toast.error(res?.message || "Return failed.");
    } catch { toast.error("An error occurred."); } finally { setActionLoading(false); setReturnItem(null); }
  };

  const columns: ColumnDef<FeedbackRequestLog>[] = [
    { accessorKey: "feedbackRequestType", header: "Type", cell: ({ row }) => <Badge variant="secondary">{getFeedbackRequestTypeLabel(row.original.feedbackRequestType)}</Badge> },
    { accessorKey: "requestOwnerStaffName", header: "Requestor", cell: ({ row }) => row.original.requestOwnerStaffName ?? row.original.requestOwnerStaffId },
    { accessorKey: "referenceId", header: "Reference", cell: ({ row }) => <span className="truncate max-w-[100px] inline-block font-mono text-xs">{row.original.referenceId}</span> },
    { accessorKey: "timeInitiated", header: "Initiated", cell: ({ row }) => row.original.timeInitiated?.split("T")[0] ?? "—" },
    { accessorKey: "hasSla", header: "SLA", cell: ({ row }) => <Badge variant={row.original.hasSla ? "outline" : "secondary"}>{row.original.hasSla ? "Yes" : "No"}</Badge> },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="outline" onClick={() => openDetails(row.original.feedbackRequestLogId)} title="View Details">
            <Eye className="h-3.5 w-3.5" />
          </Button>
          <Button size="sm" variant="ghost" onClick={() => setApproveItem(row.original.feedbackRequestLogId)} title="Approve">
            <CheckCircle className="h-3.5 w-3.5 text-green-600" />
          </Button>
          <Button size="sm" variant="ghost" onClick={() => setReturnItem(row.original.feedbackRequestLogId)} title="Return">
            <Undo2 className="h-3.5 w-3.5 text-orange-600" />
          </Button>
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Pending Requests" breadcrumbs={[{ label: "Requests" }, { label: "Pending" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Pending Requests"
        description="Feedback requests awaiting your action"
        breadcrumbs={[{ label: "Requests" }, { label: "Pending" }]}
      />

      {requests.length > 0 ? (
        <DataTable columns={columns} data={requests} searchKey="referenceId" searchPlaceholder="Search by reference..." />
      ) : (
        <EmptyState icon={Clock} title="No Pending Requests" description="You have no pending feedback requests to action." />
      )}

      {/* Details Dialog */}
      <FormDialog open={detailsOpen} onOpenChange={setDetailsOpen} title="Request Details" className="sm:max-w-lg">
        {detailsLoading ? (
          <div className="flex items-center justify-center py-8"><Loader2 className="h-6 w-6 animate-spin" /></div>
        ) : selectedRequest ? (
          <div className="space-y-3 text-sm">
            <div className="grid grid-cols-2 gap-3">
              <div><span className="text-muted-foreground">Type:</span> <span className="font-medium">{getFeedbackRequestTypeLabel(selectedRequest.feedbackRequestType)}</span></div>
              <div><span className="text-muted-foreground">Reference:</span> <span className="font-medium font-mono text-xs">{selectedRequest.referenceId}</span></div>
              <div><span className="text-muted-foreground">Requestor:</span> <span className="font-medium">{selectedRequest.requestOwnerStaffName ?? selectedRequest.requestOwnerStaffId}</span></div>
              <div><span className="text-muted-foreground">Initiated:</span> <span className="font-medium">{selectedRequest.timeInitiated?.split("T")[0]}</span></div>
              <div><span className="text-muted-foreground">SLA:</span> <span className="font-medium">{selectedRequest.hasSla ? "Yes" : "No"}</span></div>
              <div><span className="text-muted-foreground">Status:</span> {selectedRequest.recordStatus != null ? <StatusBadge status={selectedRequest.recordStatus} /> : "—"}</div>
            </div>
            {selectedRequest.requestOwnerComment && (
              <div>
                <span className="text-muted-foreground text-xs">Requestor Comment:</span>
                <p className="text-sm mt-1 bg-muted/50 rounded-md p-2">{selectedRequest.requestOwnerComment}</p>
              </div>
            )}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No details available.</p>
        )}
      </FormDialog>

      {/* Approve Confirmation */}
      <ConfirmationDialog
        open={!!approveItem}
        onOpenChange={() => setApproveItem(null)}
        title="Approve Request"
        description="Are you sure you want to approve this request?"
        confirmLabel={actionLoading ? "Approving..." : "Approve"}
        showReasonInput
        reasonLabel="Comment (optional)"
        onConfirm={handleApprove}
      />

      {/* Return Confirmation */}
      <ConfirmationDialog
        open={!!returnItem}
        onOpenChange={() => setReturnItem(null)}
        title="Return Request"
        description="Are you sure you want to return this request?"
        confirmLabel={actionLoading ? "Returning..." : "Return"}
        variant="destructive"
        showReasonInput
        reasonLabel="Reason for returning"
        onConfirm={handleReturn}
      />
    </div>
  );
}
