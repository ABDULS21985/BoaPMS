"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getAssignJobGradeGroups, saveAssignJobGradeGroup, getJobGrades, getJobGradeGroups } from "@/lib/api/competency";
import type { AssignJobGradeGroup, JobGrade, JobGradeGroup } from "@/types/competency";

export default function AssignGroupsPage() {
  const [items, setItems] = useState<AssignJobGradeGroup[]>([]);
  const [grades, setGrades] = useState<JobGrade[]>([]);
  const [groups, setGroups] = useState<JobGradeGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<AssignJobGradeGroup | null>(null);
  const [formData, setFormData] = useState({ jobGradeGroupId: "", jobGradeId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, gradeRes, groupRes] = await Promise.all([getAssignJobGradeGroups(), getJobGrades(), getJobGradeGroups()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (gradeRes?.data) setGrades(Array.isArray(gradeRes.data) ? gradeRes.data : []);
      if (groupRes?.data) setGroups(Array.isArray(groupRes.data) ? groupRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ jobGradeGroupId: "", jobGradeId: "" }); setOpen(true); };
  const openEdit = (item: AssignJobGradeGroup) => {
    setEditItem(item);
    setFormData({ jobGradeGroupId: String(item.jobGradeGroupId), jobGradeId: String(item.jobGradeId) });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.jobGradeGroupId || !formData.jobGradeId) { toast.error("Both fields are required."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { assignJobGradeGroupId: editItem.assignJobGradeGroupId } : {}),
        jobGradeGroupId: Number(formData.jobGradeGroupId),
        jobGradeId: Number(formData.jobGradeId),
      };
      const res = await saveAssignJobGradeGroup(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Updated." : "Assigned."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<AssignJobGradeGroup>[] = [
    { accessorKey: "jobGradeGroupName", header: "Grade Group" },
    { accessorKey: "jobGradeName", header: "Job Grade" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Assign Grade Groups" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Assign Groups" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Assign Grade Groups" description="Assign job grades to grade groups"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Assign Groups" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Assign Grade</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="jobGradeGroupName" searchPlaceholder="Search assignments..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Assignment" : "Assign Grade to Group"} isEdit={!!editItem}>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Grade Group *</Label>
            <Select value={formData.jobGradeGroupId} onValueChange={(v) => setFormData({ ...formData, jobGradeGroupId: v })}>
              <SelectTrigger><SelectValue placeholder="Select grade group" /></SelectTrigger>
              <SelectContent>{groups.map((g) => <SelectItem key={g.jobGradeGroupId} value={String(g.jobGradeGroupId)}>{g.groupName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Job Grade *</Label>
            <Select value={formData.jobGradeId} onValueChange={(v) => setFormData({ ...formData, jobGradeId: v })}>
              <SelectTrigger><SelectValue placeholder="Select job grade" /></SelectTrigger>
              <SelectContent>{grades.map((g) => <SelectItem key={g.jobGradeId} value={String(g.jobGradeId)}>{g.gradeCode}{g.gradeName ? ` - ${g.gradeName}` : ""}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Assign"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
