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
import { getJobRoleGrades, saveJobRoleGrade, getJobRoles, getJobGrades } from "@/lib/api/competency";
import type { JobRoleGrade, JobRole, JobGrade } from "@/types/competency";

export default function RoleGradesPage() {
  const [items, setItems] = useState<JobRoleGrade[]>([]);
  const [jobRoles, setJobRoles] = useState<JobRole[]>([]);
  const [grades, setGrades] = useState<JobGrade[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<JobRoleGrade | null>(null);
  const [formData, setFormData] = useState({ jobRoleId: "", gradeId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, jrRes, gradeRes] = await Promise.all([getJobRoleGrades(), getJobRoles(), getJobGrades()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (jrRes?.data) setJobRoles(Array.isArray(jrRes.data) ? jrRes.data : []);
      if (gradeRes?.data) setGrades(Array.isArray(gradeRes.data) ? gradeRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ jobRoleId: "", gradeId: "" }); setOpen(true); };
  const openEdit = (item: JobRoleGrade) => {
    setEditItem(item);
    setFormData({ jobRoleId: String(item.jobRoleId), gradeId: item.gradeId ?? "" });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.jobRoleId || !formData.gradeId) { toast.error("Both fields are required."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { jobRoleGradeId: editItem.jobRoleGradeId } : {}),
        jobRoleId: Number(formData.jobRoleId),
        gradeId: formData.gradeId,
      };
      const res = await saveJobRoleGrade(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Updated." : "Created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<JobRoleGrade>[] = [
    { accessorKey: "jobRoleName", header: "Job Role" },
    { accessorKey: "gradeName", header: "Grade", cell: ({ row }) => row.original.gradeName ?? row.original.gradeId ?? "—" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Role Grades" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Role Grades" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Job Role Grades" description="Map job roles to applicable grade levels"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Role Grades" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Role Grade</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="jobRoleName" searchPlaceholder="Search role grades..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Role Grade" : "Add Role Grade"} isEdit={!!editItem}>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Job Role *</Label>
            <Select value={formData.jobRoleId} onValueChange={(v) => setFormData({ ...formData, jobRoleId: v })}>
              <SelectTrigger><SelectValue placeholder="Select job role" /></SelectTrigger>
              <SelectContent>{jobRoles.map((r) => <SelectItem key={r.jobRoleId} value={String(r.jobRoleId)}>{r.jobRoleName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Grade *</Label>
            <Select value={formData.gradeId} onValueChange={(v) => setFormData({ ...formData, gradeId: v })}>
              <SelectTrigger><SelectValue placeholder="Select grade" /></SelectTrigger>
              <SelectContent>{grades.map((g) => <SelectItem key={g.jobGradeId} value={String(g.jobGradeId)}>{g.gradeCode}{g.gradeName ? ` - ${g.gradeName}` : ""}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
