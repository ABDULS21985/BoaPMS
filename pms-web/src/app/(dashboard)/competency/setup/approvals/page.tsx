"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { CheckCircle, XCircle } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import {
  getCompetencies, approveCompetency, rejectCompetency,
  getCompetencyReviewPeriods, approveCompetencyReviewPeriod,
} from "@/lib/api/competency";
import type { Competency, CompetencyReviewPeriod } from "@/types/competency";

export default function CompetencyApprovalsPage() {
  const [competencies, setCompetencies] = useState<Competency[]>([]);
  const [periods, setPeriods] = useState<CompetencyReviewPeriod[]>([]);
  const [loading, setLoading] = useState(true);
  const [approveItem, setApproveItem] = useState<{ type: "competency" | "period"; id: number; name: string } | null>(null);
  const [rejectItem, setRejectItem] = useState<{ id: number; name: string } | null>(null);
  const [rejectReason, setRejectReason] = useState("");
  const [actionLoading, setActionLoading] = useState(false);

  const loadData = async () => {
    setLoading(true);
    try {
      const [compRes, periodRes] = await Promise.all([getCompetencies(), getCompetencyReviewPeriods()]);
      if (compRes?.data) {
        const all = Array.isArray(compRes.data) ? compRes.data : [];
        setCompetencies(all.filter((c) => !c.isApproved && !c.isRejected));
      }
      if (periodRes?.data) {
        const all = Array.isArray(periodRes.data) ? periodRes.data : [];
        setPeriods(all.filter((p) => !p.isApproved));
      }
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const handleApprove = async () => {
    if (!approveItem) return;
    setActionLoading(true);
    try {
      let res;
      if (approveItem.type === "competency") {
        res = await approveCompetency({ competencyId: approveItem.id });
      } else {
        res = await approveCompetencyReviewPeriod({ reviewPeriodId: approveItem.id });
      }
      if (res?.isSuccess) { toast.success(`${approveItem.name} approved.`); loadData(); }
      else toast.error(res?.message || "Approval failed.");
    } catch { toast.error("An error occurred."); } finally { setActionLoading(false); setApproveItem(null); }
  };

  const handleReject = async (reason?: string) => {
    if (!rejectItem) return;
    setActionLoading(true);
    try {
      const res = await rejectCompetency({ competencyId: rejectItem.id, rejectionReason: reason ?? rejectReason });
      if (res?.isSuccess) { toast.success("Competency rejected."); loadData(); }
      else toast.error(res?.message || "Rejection failed.");
    } catch { toast.error("An error occurred."); } finally { setActionLoading(false); setRejectItem(null); setRejectReason(""); }
  };

  const compColumns: ColumnDef<Competency>[] = [
    { accessorKey: "competencyName", header: "Competency" },
    { accessorKey: "competencyCategoryName", header: "Category" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "createdBy", header: "Created By" },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => setApproveItem({ type: "competency", id: row.original.competencyId, name: row.original.competencyName })}><CheckCircle className="h-3.5 w-3.5 text-green-600" /></Button>
          <Button size="sm" variant="ghost" onClick={() => setRejectItem({ id: row.original.competencyId, name: row.original.competencyName })}><XCircle className="h-3.5 w-3.5 text-red-600" /></Button>
        </div>
      ),
    },
  ];

  const periodColumns: ColumnDef<CompetencyReviewPeriod>[] = [
    { accessorKey: "name", header: "Period Name" },
    { accessorKey: "bankYearName", header: "Bank Year" },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "createdBy", header: "Created By" },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <Button size="sm" variant="ghost" onClick={() => setApproveItem({ type: "period", id: row.original.reviewPeriodId, name: row.original.name })}>
          <CheckCircle className="h-3.5 w-3.5 text-green-600" />
        </Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Pending Approvals" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Approvals" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Pending Approvals" description="Review and approve pending competencies and review periods"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Approvals" }]}
      />

      <Tabs defaultValue="competencies">
        <TabsList>
          <TabsTrigger value="competencies">
            Competencies {competencies.length > 0 && <Badge variant="secondary" className="ml-2">{competencies.length}</Badge>}
          </TabsTrigger>
          <TabsTrigger value="periods">
            Review Periods {periods.length > 0 && <Badge variant="secondary" className="ml-2">{periods.length}</Badge>}
          </TabsTrigger>
        </TabsList>

        <TabsContent value="competencies" className="mt-4">
          {competencies.length > 0 ? (
            <DataTable columns={compColumns} data={competencies} searchKey="competencyName" searchPlaceholder="Search competencies..." />
          ) : (
            <EmptyState icon={CheckCircle} title="No Pending Competencies" description="All competencies have been reviewed." />
          )}
        </TabsContent>

        <TabsContent value="periods" className="mt-4">
          {periods.length > 0 ? (
            <DataTable columns={periodColumns} data={periods} searchKey="name" searchPlaceholder="Search periods..." />
          ) : (
            <EmptyState icon={CheckCircle} title="No Pending Periods" description="All review periods have been approved." />
          )}
        </TabsContent>
      </Tabs>

      <ConfirmationDialog open={!!approveItem} onOpenChange={() => setApproveItem(null)} title="Approve" description={`Approve "${approveItem?.name}"?`} confirmLabel={actionLoading ? "Approving..." : "Approve"} onConfirm={handleApprove} />

      <ConfirmationDialog open={!!rejectItem} onOpenChange={() => { setRejectItem(null); setRejectReason(""); }} title="Reject Competency" description={`Reject "${rejectItem?.name}"?`} confirmLabel={actionLoading ? "Rejecting..." : "Reject"} variant="destructive" showReasonInput reasonLabel="Rejection Reason" onConfirm={handleReject} />
    </div>
  );
}
