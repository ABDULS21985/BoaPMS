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
import { getJobRoles, saveJobRole } from "@/lib/api/competency";
import type { JobRole } from "@/types/competency";

export default function JobRolesPage() {
  const [items, setItems] = useState<JobRole[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<JobRole | null>(null);
  const [formData, setFormData] = useState({ jobRoleName: "", description: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getJobRoles();
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ jobRoleName: "", description: "" }); setOpen(true); };
  const openEdit = (item: JobRole) => {
    setEditItem(item);
    setFormData({ jobRoleName: item.jobRoleName, description: item.description ?? "" });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.jobRoleName) { toast.error("Job role name is required."); return; }
    setSaving(true);
    try {
      const payload = editItem
        ? { jobRoleId: editItem.jobRoleId, ...formData }
        : formData;
      const res = await saveJobRole(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Job role updated." : "Job role created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<JobRole>[] = [
    { accessorKey: "jobRoleName", header: "Job Role" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Job Roles" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Job Roles" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Job Roles" description="Manage job roles for competency assessment"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Job Roles" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Job Role</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="jobRoleName" searchPlaceholder="Search job roles..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Job Role" : "Add Job Role"} isEdit={!!editItem} editWarning="Update the job role below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Job Role Name *</Label><Input value={formData.jobRoleName} onChange={(e) => setFormData({ ...formData, jobRoleName: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
