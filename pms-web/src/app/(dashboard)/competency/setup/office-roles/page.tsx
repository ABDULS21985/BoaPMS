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
import { getOfficeJobRoles, saveOfficeJobRole, getJobRoles } from "@/lib/api/competency";
import { getOffices } from "@/lib/api/organogram";
import type { OfficeJobRole, JobRole } from "@/types/competency";

interface Office { officeId: number; officeName: string; }

export default function OfficeJobRolesPage() {
  const [items, setItems] = useState<OfficeJobRole[]>([]);
  const [jobRoles, setJobRoles] = useState<JobRole[]>([]);
  const [offices, setOffices] = useState<Office[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<OfficeJobRole | null>(null);
  const [formData, setFormData] = useState({ officeId: "", jobRoleId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, jrRes, offRes] = await Promise.all([getOfficeJobRoles(), getJobRoles(), getOffices()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (jrRes?.data) setJobRoles(Array.isArray(jrRes.data) ? jrRes.data : []);
      if (offRes?.data) setOffices(Array.isArray(offRes.data) ? offRes.data as Office[] : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ officeId: "", jobRoleId: "" }); setOpen(true); };
  const openEdit = (item: OfficeJobRole) => {
    setEditItem(item);
    setFormData({ officeId: String(item.officeId), jobRoleId: String(item.jobRoleId) });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.officeId || !formData.jobRoleId) { toast.error("Both fields are required."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { officeJobRoleId: editItem.officeJobRoleId } : {}),
        officeId: Number(formData.officeId),
        jobRoleId: Number(formData.jobRoleId),
      };
      const res = await saveOfficeJobRole(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Updated." : "Created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<OfficeJobRole>[] = [
    { accessorKey: "officeName", header: "Office" },
    { accessorKey: "jobRoleName", header: "Job Role" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Office Job Roles" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Office Roles" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Office Job Roles" description="Map job roles to offices"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Office Roles" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Office Role</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="officeName" searchPlaceholder="Search office roles..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Office Role" : "Add Office Role"} isEdit={!!editItem}>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Office *</Label>
            <Select value={formData.officeId} onValueChange={(v) => setFormData({ ...formData, officeId: v })}>
              <SelectTrigger><SelectValue placeholder="Select office" /></SelectTrigger>
              <SelectContent>{offices.map((o) => <SelectItem key={o.officeId} value={String(o.officeId)}>{o.officeName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Job Role *</Label>
            <Select value={formData.jobRoleId} onValueChange={(v) => setFormData({ ...formData, jobRoleId: v })}>
              <SelectTrigger><SelectValue placeholder="Select job role" /></SelectTrigger>
              <SelectContent>{jobRoles.map((r) => <SelectItem key={r.jobRoleId} value={String(r.jobRoleId)}>{r.jobRoleName}</SelectItem>)}</SelectContent>
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
