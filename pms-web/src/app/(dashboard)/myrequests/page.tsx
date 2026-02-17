"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { ClipboardList, Eye, UserPlus, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { FormDialog } from "@/components/shared/form-dialog";
import { Badge } from "@/components/ui/badge";
import { getStaffRequests, getRequestDetails, reassignSelfRequest } from "@/lib/api/pms-engine";
import { getEmployeeDetail } from "@/lib/api/dashboard";
import { getFeedbackRequestTypeLabel } from "@/lib/feedback-helpers";
import type { FeedbackRequestLog } from "@/types/performance";
import type { EmployeeErpDetails } from "@/types/dashboard";
import type { ColumnDef } from "@tanstack/react-table";

export default function MyRequestsPage() {
  const { data: session } = useSession();
  const staffId = session?.user?.id ?? "";

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

  const loadData = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const res = await getStaffRequests(staffId);
      if (res?.data) setRequests(Array.isArray(res.data) ? res.data : []);
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
    if (!searchedStaff) { toast.error("Please search and select a staff member."); return; }
    setSubmitting(true);
    try {
      const res = await reassignSelfRequest({
        requestId: reassignRequestId,
        currentStaffId: staffId,
        newAssignedStaffId: searchedStaff.employeeNumber,
      });
      if (res?.isSuccess) { toast.success("Request reassigned."); setReassignOpen(false); loadData(); }
      else toast.error(res?.message || "Reassignment failed.");
    } catch { toast.error("An error occurred."); } finally { setSubmitting(false); }
  };

  const columns: ColumnDef<FeedbackRequestLog>[] = [
    { accessorKey: "feedbackRequestType", header: "Type", cell: ({ row }) => <Badge variant="secondary">{getFeedbackRequestTypeLabel(row.original.feedbackRequestType)}</Badge> },
    { accessorKey: "referenceId", header: "Reference", cell: ({ row }) => <span className="truncate max-w-[100px] inline-block font-mono text-xs">{row.original.referenceId}</span> },
    { accessorKey: "requestOwnerStaffName", header: "Requestor", cell: ({ row }) => row.original.requestOwnerStaffName ?? row.original.requestOwnerStaffId },
    { accessorKey: "assignedStaffName", header: "Assigned To", cell: ({ row }) => row.original.assignedStaffName ?? row.original.assignedStaffId },
    { accessorKey: "timeInitiated", header: "Initiated", cell: ({ row }) => row.original.timeInitiated?.split("T")[0] ?? "—" },
    { accessorKey: "timeCompleted", header: "Completed", cell: ({ row }) => row.original.timeCompleted?.split("T")[0] ?? "—" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    { accessorKey: "hasSla", header: "SLA", cell: ({ row }) => <Badge variant={row.original.hasSla ? "outline" : "secondary"}>{row.original.hasSla ? "Yes" : "No"}</Badge> },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="outline" onClick={() => openDetails(row.original.feedbackRequestLogId)}>
            <Eye className="mr-1 h-3.5 w-3.5" />View
          </Button>
          {row.original.recordStatus !== 10 && (
            <Button size="sm" variant="ghost" onClick={() => openReassign(row.original.feedbackRequestLogId)}>
              <UserPlus className="h-3.5 w-3.5" />
            </Button>
          )}
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="My Requests" breadcrumbs={[{ label: "Requests" }, { label: "My Requests" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="My Requests"
        description="View and manage your submitted requests"
        breadcrumbs={[{ label: "Requests" }, { label: "My Requests" }]}
      />

      {requests.length > 0 ? (
        <DataTable columns={columns} data={requests} searchKey="referenceId" searchPlaceholder="Search by reference..." />
      ) : (
        <EmptyState icon={ClipboardList} title="No Requests" description="You have no requests for the current review period." />
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
            {selectedRequest.requestOwnerComment && (
              <div><span className="text-muted-foreground">Owner Comment:</span><p className="mt-1 text-sm">{selectedRequest.requestOwnerComment}</p></div>
            )}
            {selectedRequest.assignedStaffComment && (
              <div><span className="text-muted-foreground">Assigned Staff Comment:</span><p className="mt-1 text-sm">{selectedRequest.assignedStaffComment}</p></div>
            )}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No details available.</p>
        )}
      </FormDialog>

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
    </div>
  );
}
