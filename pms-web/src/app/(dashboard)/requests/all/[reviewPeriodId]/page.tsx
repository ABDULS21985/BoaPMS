"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useParams, useRouter } from "next/navigation";
import { ClipboardList, Eye, UserPlus, XCircle, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { FormDialog } from "@/components/shared/form-dialog";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import {
  getAllRequests,
  getRequestDetails,
  reassignRequest,
  closeRequest,
} from "@/lib/api/pms-engine";
import { getReviewPeriodDetails } from "@/lib/api/review-periods";
import { getEmployeeDetail } from "@/lib/api/dashboard";
import { getFeedbackRequestTypeLabel } from "@/lib/feedback-helpers";
import { Roles } from "@/stores/auth-store";
import { Status } from "@/types/enums";
import type { FeedbackRequestLog, PerformanceReviewPeriod } from "@/types/performance";
import type { EmployeeErpDetails } from "@/types/dashboard";
import type { ColumnDef } from "@tanstack/react-table";

export default function AllRequestsForPeriodPage() {
  const params = useParams<{ reviewPeriodId: string }>();
  const { data: session } = useSession();
  const router = useRouter();
  const staffId = session?.user?.id ?? "";
  const userRoles = session?.user?.roles ?? [];
  const isAdmin = userRoles.some((r: string) => ([Roles.Admin, Roles.SuperAdmin, Roles.HrAdmin, Roles.HrReportAdmin, Roles.HrApprover] as string[]).includes(r));

  const [period, setPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [requests, setRequests] = useState<FeedbackRequestLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [detailsLoading, setDetailsLoading] = useState(false);
  const [selectedRequest, setSelectedRequest] = useState<FeedbackRequestLog | null>(null);

  const [reassignOpen, setReassignOpen] = useState(false);
  const [reassignRequestId, setReassignRequestId] = useState("");
  const [newStaffId, setNewStaffId] = useState("");
  const [searchedStaff, setSearchedStaff] = useState<EmployeeErpDetails | null>(null);
  const [searching, setSearching] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const [closeOpen, setCloseOpen] = useState(false);
  const [closeRequestId, setCloseRequestId] = useState("");

  const loadData = async () => {
    if (!isAdmin || !staffId) return;
    setLoading(true);
    try {
      const [periodRes, requestsRes] = await Promise.allSettled([
        getReviewPeriodDetails(params.reviewPeriodId),
        getAllRequests(staffId),
      ]);
      if (periodRes.status === "fulfilled" && periodRes.value?.data) setPeriod(periodRes.value.data);
      if (requestsRes.status === "fulfilled" && requestsRes.value?.data) {
        const all = Array.isArray(requestsRes.value.data) ? requestsRes.value.data : [];
        setRequests(all.filter((r) => r.reviewPeriodId === params.reviewPeriodId));
      }
    } catch { /* silent */ } finally { setLoading(false); }
  };

  useEffect(() => {
    if (!isAdmin) { router.replace("/"); return; }
    loadData();
  }, [isAdmin, staffId, params.reviewPeriodId, router]);

  const openDetails = async (requestId: string) => {
    setDetailsOpen(true);
    setDetailsLoading(true);
    try {
      const res = await getRequestDetails(requestId);
      setSelectedRequest(res?.data ?? null);
    } catch { setSelectedRequest(null); } finally { setDetailsLoading(false); }
  };

  const openReassign = (requestId: string) => {
    setReassignRequestId(requestId);
    setNewStaffId("");
    setSearchedStaff(null);
    setReassignOpen(true);
  };

  const searchStaff = async () => {
    if (newStaffId.length < 3) { toast.error("Enter at least 3 characters."); return; }
    setSearching(true);
    try {
      const res = await getEmployeeDetail(newStaffId);
      if (res?.data) setSearchedStaff(res.data);
      else { setSearchedStaff(null); toast.error("Staff not found."); }
    } catch { setSearchedStaff(null); toast.error("Staff lookup failed."); } finally { setSearching(false); }
  };

  const handleReassign = async () => {
    if (!searchedStaff) return;
    setSubmitting(true);
    try {
      const res = await reassignRequest({
        requestId: reassignRequestId,
        newAssignedStaffId: searchedStaff.employeeNumber,
      });
      if (res?.isSuccess) { toast.success("Request reassigned."); setReassignOpen(false); loadData(); }
      else toast.error(res?.message || "Reassignment failed.");
    } catch { toast.error("An error occurred."); } finally { setSubmitting(false); }
  };

  const handleClose = async () => {
    try {
      const res = await closeRequest({ requestId: closeRequestId });
      if (res?.isSuccess) { toast.success("Request closed."); loadData(); }
      else toast.error(res?.message || "Close failed.");
    } catch { toast.error("An error occurred."); }
  };

  const columns: ColumnDef<FeedbackRequestLog>[] = [
    { accessorKey: "feedbackRequestType", header: "Type", cell: ({ row }) => <Badge variant="secondary">{getFeedbackRequestTypeLabel(row.original.feedbackRequestType)}</Badge> },
    { accessorKey: "referenceId", header: "Reference", cell: ({ row }) => <span className="truncate max-w-[100px] inline-block font-mono text-xs">{row.original.referenceId}</span> },
    { accessorKey: "assignedStaffName", header: "Assignee", cell: ({ row }) => row.original.assignedStaffName ?? row.original.assignedStaffId },
    { accessorKey: "requestOwnerStaffName", header: "Requestor", cell: ({ row }) => row.original.requestOwnerStaffName ?? row.original.requestOwnerStaffId },
    { accessorKey: "timeInitiated", header: "Initiated", cell: ({ row }) => row.original.timeInitiated?.split("T")[0] ?? "—" },
    { accessorKey: "timeCompleted", header: "Completed", cell: ({ row }) => row.original.timeCompleted?.split("T")[0] ?? "—" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    { accessorKey: "hasSla", header: "SLA", cell: ({ row }) => <Badge variant={row.original.hasSla ? "outline" : "secondary"}>{row.original.hasSla ? "Yes" : "No"}</Badge> },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="outline" onClick={() => openDetails(row.original.feedbackRequestLogId)}>
            <Eye className="h-3.5 w-3.5" />
          </Button>
          {row.original.recordStatus !== Status.Completed && (
            <>
              <Button size="sm" variant="ghost" onClick={() => openReassign(row.original.feedbackRequestLogId)} title="Reassign">
                <UserPlus className="h-3.5 w-3.5" />
              </Button>
              <Button size="sm" variant="ghost" onClick={() => { setCloseRequestId(row.original.feedbackRequestLogId); setCloseOpen(true); }} title="Close">
                <XCircle className="h-3.5 w-3.5 text-destructive" />
              </Button>
            </>
          )}
        </div>
      ),
    },
  ];

  if (!isAdmin) return null;
  if (loading) return <div><PageHeader title="All Requests" breadcrumbs={[{ label: "Requests" }, { label: "Review Periods", href: "/requests/review-periods" }, { label: "All Requests" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title={`All Requests - ${period?.name ?? params.reviewPeriodId}`}
        description={period ? `${period.startDate?.split("T")[0]} to ${period.endDate?.split("T")[0]}` : undefined}
        breadcrumbs={[
          { label: "Requests" },
          { label: "Review Periods", href: "/requests/review-periods" },
          { label: period?.name ?? "All Requests" },
        ]}
      />

      {requests.length > 0 ? (
        <DataTable columns={columns} data={requests} searchKey="referenceId" searchPlaceholder="Search by reference..." />
      ) : (
        <EmptyState icon={ClipboardList} title="No Requests" description="No requests found for this review period." />
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
              <div><span className="text-muted-foreground">Assignee:</span> <span className="font-medium">{selectedRequest.assignedStaffName ?? selectedRequest.assignedStaffId}</span></div>
              <div><span className="text-muted-foreground">Requestor:</span> <span className="font-medium">{selectedRequest.requestOwnerStaffName ?? selectedRequest.requestOwnerStaffId}</span></div>
              <div><span className="text-muted-foreground">Initiated:</span> <span className="font-medium">{selectedRequest.timeInitiated?.split("T")[0]}</span></div>
              <div><span className="text-muted-foreground">Completed:</span> <span className="font-medium">{selectedRequest.timeCompleted?.split("T")[0] ?? "Pending"}</span></div>
              <div><span className="text-muted-foreground">SLA:</span> <span className="font-medium">{selectedRequest.hasSla ? "Yes" : "No"}</span></div>
              <div><span className="text-muted-foreground">Status:</span> {selectedRequest.recordStatus != null ? <StatusBadge status={selectedRequest.recordStatus} /> : "—"}</div>
            </div>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No details available.</p>
        )}
      </FormDialog>

      {/* Reassign Dialog */}
      <FormDialog open={reassignOpen} onOpenChange={setReassignOpen} title="Reassign Request">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Staff ID / Employee Number</Label>
            <div className="flex gap-2">
              <Input value={newStaffId} onChange={(e) => { setNewStaffId(e.target.value); setSearchedStaff(null); }} placeholder="Enter staff ID..." />
              <Button variant="outline" onClick={searchStaff} disabled={searching}>
                {searching ? <Loader2 className="h-4 w-4 animate-spin" /> : "Search"}
              </Button>
            </div>
          </div>
          {searchedStaff && (
            <div className="rounded-md border p-3 text-sm bg-muted/50">
              <p className="font-medium">{searchedStaff.firstName} {searchedStaff.lastName}</p>
              <p className="text-muted-foreground">{searchedStaff.jobTitle} - {searchedStaff.departmentName}</p>
            </div>
          )}
          <div className="flex gap-3 pt-2">
            <Button variant="outline" className="flex-1" onClick={() => setReassignOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleReassign} disabled={submitting || !searchedStaff}>
              {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Reassign
            </Button>
          </div>
        </div>
      </FormDialog>

      {/* Close Confirmation */}
      <ConfirmationDialog
        open={closeOpen}
        onOpenChange={setCloseOpen}
        title="Close Request"
        description="Are you sure you want to close this request? This action cannot be undone."
        confirmLabel="Close Request"
        variant="destructive"
        onConfirm={handleClose}
      />
    </div>
  );
}
