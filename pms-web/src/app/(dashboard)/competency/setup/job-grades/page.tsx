"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getJobGrades, saveJobGrade } from "@/lib/api/competency";
import type { JobGrade } from "@/types/competency";

export default function JobGradesPage() {
  const [items, setItems] = useState<JobGrade[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<JobGrade | null>(null);
  const [formData, setFormData] = useState({ gradeCode: "", gradeName: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getJobGrades();
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ gradeCode: "", gradeName: "" }); setOpen(true); };
  const openEdit = (item: JobGrade) => {
    setEditItem(item);
    setFormData({ gradeCode: item.gradeCode, gradeName: item.gradeName ?? "" });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.gradeCode) { toast.error("Grade code is required."); return; }
    setSaving(true);
    try {
      const payload = editItem
        ? { jobGradeId: editItem.jobGradeId, ...formData }
        : formData;
      const res = await saveJobGrade(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Job grade updated." : "Job grade created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<JobGrade>[] = [
    { accessorKey: "gradeCode", header: "Grade Code" },
    { accessorKey: "gradeName", header: "Grade Name" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Job Grades" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Job Grades" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Job Grades" description="Manage job grade levels"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Job Grades" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Job Grade</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="gradeCode" searchPlaceholder="Search grades..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Job Grade" : "Add Job Grade"} isEdit={!!editItem} editWarning="Update the job grade below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Grade Code *</Label><Input value={formData.gradeCode} onChange={(e) => setFormData({ ...formData, gradeCode: e.target.value })} /></div>
          <div className="space-y-2"><Label>Grade Name</Label><Input value={formData.gradeName} onChange={(e) => setFormData({ ...formData, gradeName: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
