"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Plus, Eye, Loader2, Users2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getCommittees, saveDraftCommittee } from "@/lib/api/pms-engine";
import { getReviewPeriods } from "@/lib/api/review-periods";
import { getDepartments } from "@/lib/api/organogram";
import type { Committee, PerformanceReviewPeriod } from "@/types/performance";
import type { Department } from "@/types/organogram";
import type { ColumnDef } from "@tanstack/react-table";

export default function CommitteesPage() {
  const router = useRouter();
  const [committees, setCommittees] = useState<Committee[]>([]);
  const [periods, setPeriods] = useState<PerformanceReviewPeriod[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({ name: "", description: "", deliverables: "", reviewPeriodId: "", departmentId: "", chairperson: "", startDate: "", endDate: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [commRes, periodRes, deptRes] = await Promise.all([getCommittees(), getReviewPeriods(), getDepartments()]);
      if (commRes?.data) setCommittees(Array.isArray(commRes.data) ? commRes.data : []);
      if (periodRes?.data) setPeriods(Array.isArray(periodRes.data) ? periodRes.data : []);
      if (deptRes?.data) setDepartments(Array.isArray(deptRes.data) ? deptRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const handleSave = async () => {
    if (!formData.name || !formData.reviewPeriodId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const res = await saveDraftCommittee({ ...formData, departmentId: formData.departmentId ? Number(formData.departmentId) : undefined });
      if (res?.isSuccess) { toast.success("Committee created as draft."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const openAdd = () => { setFormData({ name: "", description: "", deliverables: "", reviewPeriodId: "", departmentId: "", chairperson: "", startDate: "", endDate: "" }); setOpen(true); };

  const columns: ColumnDef<Committee>[] = [
    { accessorKey: "name", header: "Committee Name", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "chairperson", header: "Chairperson" },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="outline" onClick={() => router.push(`/committees/${row.original.committeeId}`)}><Eye className="mr-1 h-3.5 w-3.5" />View</Button> },
  ];

  if (loading) return <div><PageHeader title="Committees" breadcrumbs={[{ label: "Committees" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Committees" description="Manage all committees across the organization" breadcrumbs={[{ label: "Committees" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />New Committee</Button>} />
      {committees.length > 0 ? (
        <DataTable columns={columns} data={committees} searchKey="name" searchPlaceholder="Search committees..." />
      ) : (
        <EmptyState icon={Users2} title="No Committees" description="No committees have been created yet." />
      )}

      <FormSheet open={open} onOpenChange={setOpen} title="New Committee" className="sm:max-w-lg overflow-y-auto">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2">
            <Label>Review Period *</Label>
            <Select value={formData.reviewPeriodId} onValueChange={(v) => setFormData({ ...formData, reviewPeriodId: v })}>
              <SelectTrigger><SelectValue placeholder="Select period" /></SelectTrigger>
              <SelectContent>{periods.map((p) => <SelectItem key={p.periodId} value={p.periodId}>{p.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Department</Label>
            <Select value={formData.departmentId} onValueChange={(v) => setFormData({ ...formData, departmentId: v })}>
              <SelectTrigger><SelectValue placeholder="Select department" /></SelectTrigger>
              <SelectContent>{departments.map((d) => <SelectItem key={d.departmentId} value={String(d.departmentId)}>{d.departmentName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Chairperson</Label><Input value={formData.chairperson} onChange={(e) => setFormData({ ...formData, chairperson: e.target.value })} placeholder="Staff ID" /></div>
          <div className="space-y-2"><Label>Deliverables</Label><Input value={formData.deliverables} onChange={(e) => setFormData({ ...formData, deliverables: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-2"><Label>Start Date</Label><Input type="date" value={formData.startDate} onChange={(e) => setFormData({ ...formData, startDate: e.target.value })} /></div>
            <div className="space-y-2"><Label>End Date</Label><Input type="date" value={formData.endDate} onChange={(e) => setFormData({ ...formData, endDate: e.target.value })} /></div>
          </div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save Draft</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
