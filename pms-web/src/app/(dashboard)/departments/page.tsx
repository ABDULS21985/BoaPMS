"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Eye, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getDirectorates, getDepartments, saveDepartment } from "@/lib/api/organogram";
import type { Department, Directorate } from "@/types/organogram";

export default function DepartmentsPage() {
  const router = useRouter();
  const [items, setItems] = useState<Department[]>([]);
  const [directorates, setDirectorates] = useState<Directorate[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<Department | null>(null);
  const [formData, setFormData] = useState({ departmentName: "", departmentCode: "", directorateId: "", isActive: true });

  const loadData = async () => {
    setLoading(true);
    try {
      const [deptRes, dirRes] = await Promise.all([getDepartments(), getDirectorates()]);
      if (deptRes?.data) setItems(Array.isArray(deptRes.data) ? deptRes.data : []);
      if (dirRes?.data) setDirectorates(Array.isArray(dirRes.data) ? dirRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => {
    setEditItem(null);
    setFormData({ departmentName: "", departmentCode: "", directorateId: "", isActive: true });
    setOpen(true);
  };

  const openEdit = (item: Department) => {
    setEditItem(item);
    setFormData({
      departmentName: item.departmentName,
      departmentCode: item.departmentCode ?? "",
      directorateId: String(item.directorateId ?? ""),
      isActive: item.isActive,
    });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.departmentName || !formData.directorateId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { departmentId: editItem.departmentId } : {}),
        departmentName: formData.departmentName,
        departmentCode: formData.departmentCode,
        directorateId: Number(formData.directorateId),
        isActive: formData.isActive,
      };
      const res = await saveDepartment(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Department updated." : "Department created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<Department>[] = [
    { accessorKey: "directorateName", header: "Directorate" },
    { accessorKey: "departmentCode", header: "Code" },
    { accessorKey: "departmentName", header: "Department Name" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Disabled"}</Badge> },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button>
          <Button size="sm" variant="ghost" onClick={() => router.push(`/divisions?departmentId=${row.original.departmentId}`)} title="View Divisions"><Eye className="h-3.5 w-3.5" /></Button>
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Departments" breadcrumbs={[{ label: "Organogram" }, { label: "Departments" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Departments" description="Manage departments under directorates" breadcrumbs={[{ label: "Organogram" }, { label: "Departments" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Department</Button>} />
      <DataTable columns={columns} data={items} searchKey="departmentName" searchPlaceholder="Search departments..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Department" : "Add Department"} isEdit={!!editItem} editWarning="Update the selected department below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Directorate *</Label>
            <Select value={formData.directorateId} onValueChange={(v) => setFormData({ ...formData, directorateId: v })}>
              <SelectTrigger><SelectValue placeholder="Select directorate" /></SelectTrigger>
              <SelectContent>{directorates.map((d) => <SelectItem key={d.directorateId} value={String(d.directorateId)}>{d.directorateName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Department Code</Label><Input value={formData.departmentCode} onChange={(e) => setFormData({ ...formData, departmentCode: e.target.value })} /></div>
          <div className="space-y-2"><Label>Department Name *</Label><Input value={formData.departmentName} onChange={(e) => setFormData({ ...formData, departmentName: e.target.value })} /></div>
          <div className="flex items-center gap-2"><Switch checked={formData.isActive} onCheckedChange={(v) => setFormData({ ...formData, isActive: v })} /><Label>Active</Label></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
