"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import type { ColumnDef } from "@tanstack/react-table";
import { CheckCircle, XCircle, ShieldCheck } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { syncJobRoleSoa } from "@/lib/api/competency";

interface StaffJobRoleUpdate {
  staffJobRoleId: number;
  employeeId: string;
  fullName?: string;
  jobRoleName?: string;
  soaStatus: boolean;
  soaResponse?: string;
  isApproved: boolean;
  rejectionReason?: string;
  status?: string;
}

export default function JobRoleApprovalsPage() {
  const { data: session } = useSession();
  const [items, setItems] = useState<StaffJobRoleUpdate[]>([]);
  const [loading, setLoading] = useState(true);
  const [approveItem, setApproveItem] = useState<StaffJobRoleUpdate | null>(null);
  const [rejectItem, setRejectItem] = useState<StaffJobRoleUpdate | null>(null);
  const [rejectReason, setRejectReason] = useState("");
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    // Job role approvals are loaded from the sync endpoint
    // For now, show the page structure — data loading depends on backend providing a list endpoint
    setLoading(false);
  }, []);

  const handleApprove = async () => {
    if (!approveItem) return;
    setActionLoading(true);
    try {
      const res = await syncJobRoleSoa({ staffJobRoleId: approveItem.staffJobRoleId, isApproved: true });
      if (res?.isSuccess) { toast.success("Job role update approved."); setItems((prev) => prev.filter((i) => i.staffJobRoleId !== approveItem.staffJobRoleId)); }
      else toast.error(res?.message || "Approval failed.");
    } catch { toast.error("An error occurred."); } finally { setActionLoading(false); setApproveItem(null); }
  };

  const handleReject = async (reason?: string) => {
    if (!rejectItem) return;
    setActionLoading(true);
    try {
      const res = await syncJobRoleSoa({ staffJobRoleId: rejectItem.staffJobRoleId, isApproved: false, rejectionReason: reason ?? rejectReason });
      if (res?.isSuccess) { toast.success("Job role update rejected."); setItems((prev) => prev.filter((i) => i.staffJobRoleId !== rejectItem.staffJobRoleId)); }
      else toast.error(res?.message || "Rejection failed.");
    } catch { toast.error("An error occurred."); } finally { setActionLoading(false); setRejectItem(null); setRejectReason(""); }
  };

  const columns: ColumnDef<StaffJobRoleUpdate>[] = [
    { accessorKey: "fullName", header: "Employee", cell: ({ row }) => row.original.fullName ?? row.original.employeeId },
    { accessorKey: "jobRoleName", header: "Job Role" },
    { accessorKey: "soaResponse", header: "SOA Response", cell: ({ row }) => <span className="line-clamp-1">{row.original.soaResponse ?? "—"}</span> },
    { accessorKey: "status", header: "Status", cell: ({ row }) => <Badge variant="secondary">{row.original.status ?? "Pending"}</Badge> },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => setApproveItem(row.original)}><CheckCircle className="h-3.5 w-3.5 text-green-600" /></Button>
          <Button size="sm" variant="ghost" onClick={() => setRejectItem(row.original)}><XCircle className="h-3.5 w-3.5 text-red-600" /></Button>
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Job Role Approvals" breadcrumbs={[{ label: "Competency" }, { label: "Job Role Approvals" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Job Role Update Approvals" description="Review and approve job role update requests"
        breadcrumbs={[{ label: "Competency" }, { label: "Job Role Approvals" }]}
      />

      {items.length > 0 ? (
        <DataTable columns={columns} data={items} searchKey="fullName" searchPlaceholder="Search employees..." />
      ) : (
        <EmptyState icon={ShieldCheck} title="No Pending Approvals" description="All job role updates have been reviewed." />
      )}

      <ConfirmationDialog open={!!approveItem} onOpenChange={() => setApproveItem(null)} title="Approve Job Role Update" description={`Approve job role update for ${approveItem?.fullName}?`} confirmLabel={actionLoading ? "Approving..." : "Approve"} onConfirm={handleApprove} />

      <ConfirmationDialog open={!!rejectItem} onOpenChange={() => { setRejectItem(null); setRejectReason(""); }} title="Reject Job Role Update" description={`Reject job role update for ${rejectItem?.fullName}?`} confirmLabel={actionLoading ? "Rejecting..." : "Reject"} variant="destructive" showReasonInput reasonLabel="Rejection Reason" onConfirm={handleReject} />
    </div>
  );
}
