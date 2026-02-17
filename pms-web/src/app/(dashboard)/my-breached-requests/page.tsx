"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { AlertTriangle, Eye, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { FormDialog } from "@/components/shared/form-dialog";
import { getBreachedRequests, getRequestDetails } from "@/lib/api/pms-engine";
import { getActiveReviewPeriod } from "@/lib/api/dashboard";
import { getFeedbackRequestTypeLabel, formatOverdueDuration } from "@/lib/feedback-helpers";
import type { BreachedFeedbackRequestLog, FeedbackRequestLog } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function MyBreachedRequestsPage() {
  const { data: session } = useSession();
  const staffId = session?.user?.id ?? "";

  const [requests, setRequests] = useState<BreachedFeedbackRequestLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [detailsLoading, setDetailsLoading] = useState(false);
  const [selectedRequest, setSelectedRequest] = useState<FeedbackRequestLog | null>(null);

  useEffect(() => {
    if (!staffId) return;
    (async () => {
      setLoading(true);
      try {
        const periodRes = await getActiveReviewPeriod();
        const period = periodRes?.data ?? null;
        if (period) {
          const res = await getBreachedRequests(staffId, period.periodId);
          if (res?.data) setRequests(Array.isArray(res.data) ? res.data : []);
        }
      } catch { /* silent */ } finally { setLoading(false); }
    })();
  }, [staffId]);

  const openDetails = async (requestId: string) => {
    setDetailsOpen(true);
    setDetailsLoading(true);
    try {
      const res = await getRequestDetails(requestId);
      setSelectedRequest(res?.data ?? null);
    } catch { setSelectedRequest(null); } finally { setDetailsLoading(false); }
  };

  const columns: ColumnDef<BreachedFeedbackRequestLog>[] = [
    { accessorKey: "feedbackRequestType", header: "Type", cell: ({ row }) => <Badge variant="secondary">{getFeedbackRequestTypeLabel(row.original.feedbackRequestType)}</Badge> },
    { accessorKey: "referenceId", header: "Reference", cell: ({ row }) => <span className="truncate max-w-[100px] inline-block font-mono text-xs">{row.original.referenceId}</span> },
    { accessorKey: "requestOwnerStaffName", header: "Requestor", cell: ({ row }) => row.original.requestOwnerStaffName ?? row.original.requestOwnerStaffId },
    { accessorKey: "assignedStaffName", header: "Assigned To", cell: ({ row }) => row.original.assignedStaffName ?? row.original.assignedStaffId },
    { accessorKey: "timeInitiated", header: "Initiated", cell: ({ row }) => row.original.timeInitiated?.split("T")[0] ?? "—" },
    {
      accessorKey: "isBreached", header: "Overdue",
      cell: ({ row }) => (
        <Badge variant="destructive">
          {formatOverdueDuration(row.original.timeInitiated)}
        </Badge>
      ),
    },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <Button size="sm" variant="outline" onClick={() => openDetails(row.original.feedbackRequestLogId)}>
          <Eye className="mr-1 h-3.5 w-3.5" />View
        </Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Breached Requests" breadcrumbs={[{ label: "Requests" }, { label: "Breached" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Breached Requests"
        description="Requests that have exceeded their SLA deadline"
        breadcrumbs={[{ label: "Requests" }, { label: "Breached Requests" }]}
      />

      {requests.length > 0 ? (
        <DataTable columns={columns} data={requests} searchKey="referenceId" searchPlaceholder="Search by reference..." />
      ) : (
        <EmptyState icon={AlertTriangle} title="No Breached Requests" description="You have no breached or overdue requests." />
      )}

      <FormDialog open={detailsOpen} onOpenChange={setDetailsOpen} title="Request Details" className="sm:max-w-lg">
        {detailsLoading ? (
          <div className="flex items-center justify-center py-8"><Loader2 className="h-6 w-6 animate-spin" /></div>
        ) : selectedRequest ? (
          <div className="space-y-3 text-sm">
            <div className="grid grid-cols-2 gap-3">
              <div><span className="text-muted-foreground">Type:</span> <span className="font-medium">{getFeedbackRequestTypeLabel(selectedRequest.feedbackRequestType)}</span></div>
              <div><span className="text-muted-foreground">Reference:</span> <span className="font-medium font-mono text-xs">{selectedRequest.referenceId}</span></div>
              <div><span className="text-muted-foreground">Requestor:</span> <span className="font-medium">{selectedRequest.requestOwnerStaffName ?? selectedRequest.requestOwnerStaffId}</span></div>
              <div><span className="text-muted-foreground">Assigned To:</span> <span className="font-medium">{selectedRequest.assignedStaffName ?? selectedRequest.assignedStaffId}</span></div>
              <div><span className="text-muted-foreground">Initiated:</span> <span className="font-medium">{selectedRequest.timeInitiated?.split("T")[0]}</span></div>
              <div><span className="text-muted-foreground">Completed:</span> <span className="font-medium">{selectedRequest.timeCompleted?.split("T")[0] ?? "Pending"}</span></div>
              <div><span className="text-muted-foreground">SLA Bound:</span> <span className="font-medium">{selectedRequest.hasSla ? "Yes" : "No"}</span></div>
              <div><span className="text-muted-foreground">Status:</span> {selectedRequest.recordStatus != null ? <StatusBadge status={selectedRequest.recordStatus} /> : "—"}</div>
            </div>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No details available.</p>
        )}
      </FormDialog>
    </div>
  );
}
