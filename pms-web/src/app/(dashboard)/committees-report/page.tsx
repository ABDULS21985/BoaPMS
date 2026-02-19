"use client";

import { useState, useEffect } from "react";
import { type ColumnDef } from "@tanstack/react-table";
import { toast } from "sonner";
import { UserCog } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { FormDialog } from "@/components/shared/form-dialog";
import { getCommittees, updateCommittee } from "@/lib/api/pms-engine";
import { getEmployeeDetail } from "@/lib/api/dashboard";
import type { Committee } from "@/types/performance";
import type { EmployeeErpDetails } from "@/types/dashboard";

export default function CommitteesReportPage() {
  const [loading, setLoading] = useState(true);
  const [committees, setCommittees] = useState<Committee[]>([]);
  const [reassignOpen, setReassignOpen] = useState(false);
  const [selectedCommittee, setSelectedCommittee] = useState<Committee | null>(null);
  const [newChairId, setNewChairId] = useState("");
  const [verifiedEmployee, setVerifiedEmployee] = useState<EmployeeErpDetails | null>(null);
  const [verifying, setVerifying] = useState(false);
  const [saving, setSaving] = useState(false);

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getCommittees();
      setCommittees(res?.data ?? []);
    } catch {
      toast.error("Failed to load committees.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { loadData(); }, []);

  const handleVerify = async () => {
    if (!newChairId.trim()) {
      toast.error("Enter a staff ID.");
      return;
    }
    setVerifying(true);
    try {
      const res = await getEmployeeDetail(newChairId.trim());
      if (res?.data?.employeeNumber) {
        setVerifiedEmployee(res.data);
        toast.success(`Found: ${res.data.firstName} ${res.data.lastName}`);
      } else {
        setVerifiedEmployee(null);
        toast.error("Staff not found.");
      }
    } catch {
      toast.error("Verification failed.");
    } finally {
      setVerifying(false);
    }
  };

  const handleReassign = async () => {
    if (!selectedCommittee || !verifiedEmployee) return;
    setSaving(true);
    try {
      const res = await updateCommittee({
        committeeId: selectedCommittee.committeeId,
        chairperson: verifiedEmployee.employeeNumber,
      });
      if (res?.isSuccess) {
        toast.success("Chair person re-assigned successfully.");
        setReassignOpen(false);
        loadData();
      } else {
        toast.error(res?.message || "Failed to re-assign.");
      }
    } catch {
      toast.error("Failed to re-assign chair person.");
    } finally {
      setSaving(false);
    }
  };

  const openReassign = (committee: Committee) => {
    setSelectedCommittee(committee);
    setNewChairId("");
    setVerifiedEmployee(null);
    setReassignOpen(true);
  };

  const columns: ColumnDef<Committee>[] = [
    {
      id: "index",
      header: "#",
      cell: ({ row }) => row.index + 1,
    },
    { accessorKey: "name", header: "Name" },
    { accessorKey: "description", header: "Description" },
    { accessorKey: "deliverables", header: "Objective/KPI" },
    {
      accessorKey: "departmentId",
      header: "Department",
      cell: ({ row }) => row.original.departmentId ?? "-",
    },
    { accessorKey: "chairperson", header: "Chair Person" },
    {
      accessorKey: "recordStatus",
      header: "Status",
      cell: ({ row }) => <StatusBadge status={row.original.recordStatus ?? 0} />,
    },
    {
      accessorKey: "startDate",
      header: "Start Date",
      cell: ({ row }) => row.original.startDate ? new Date(row.original.startDate).toLocaleDateString() : "-",
    },
    {
      accessorKey: "endDate",
      header: "End Date",
      cell: ({ row }) => row.original.endDate ? new Date(row.original.endDate).toLocaleDateString() : "-",
    },
    {
      id: "actions",
      header: "Action",
      cell: ({ row }) => (
        <Button variant="ghost" size="sm" onClick={() => openReassign(row.original)}>
          <UserCog className="mr-1 h-4 w-4" /> Re-assign
        </Button>
      ),
    },
  ];

  if (loading) return <PageSkeleton />;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Committees Report"
        breadcrumbs={[{ label: "Committees Report" }]}
      />

      <DataTable
        columns={columns}
        data={committees}
        searchKey="name"
        searchPlaceholder="Search by committee name..."
      />

      <FormDialog
        open={reassignOpen}
        onOpenChange={setReassignOpen}
        title="Re-assign Chair Person"
        description={selectedCommittee ? `Committee: ${selectedCommittee.name}` : ""}
      >
        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label>Current Chair Person</Label>
            <Input value={selectedCommittee?.chairperson || "N/A"} readOnly className="bg-muted" />
          </div>
          <div className="space-y-1.5">
            <Label>New Chair Person Staff ID</Label>
            <div className="flex gap-2">
              <Input
                value={newChairId}
                onChange={(e) => { setNewChairId(e.target.value); setVerifiedEmployee(null); }}
                placeholder="Enter staff ID"
              />
              <Button variant="outline" onClick={handleVerify} disabled={verifying}>
                {verifying ? "Verifying..." : "Verify"}
              </Button>
            </div>
          </div>
          {verifiedEmployee && (
            <div className="rounded-md border p-3 text-sm">
              <p className="font-medium">{verifiedEmployee.firstName} {verifiedEmployee.lastName}</p>
              <p className="text-muted-foreground">{verifiedEmployee.jobTitle}</p>
              <p className="text-muted-foreground">{verifiedEmployee.departmentName}</p>
            </div>
          )}
          <div className="flex justify-end gap-2 pt-2">
            <Button variant="outline" onClick={() => setReassignOpen(false)}>Cancel</Button>
            <Button onClick={handleReassign} disabled={saving || !verifiedEmployee}>
              {saving ? "Re-assigning..." : "Re-assign"}
            </Button>
          </div>
        </div>
      </FormDialog>
    </div>
  );
}
