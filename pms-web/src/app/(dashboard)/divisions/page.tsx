"use client";

import { Suspense, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
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
import { getDepartments, getDivisions, saveDivision } from "@/lib/api/organogram";
import type { Division, Department } from "@/types/organogram";

export default function DivisionsPageWrapper() {
  return <Suspense fallback={<div><PageHeader title="Divisions" breadcrumbs={[{ label: "Organogram" }, { label: "Divisions" }]} /><PageSkeleton /></div>}><DivisionsPage /></Suspense>;
}

function DivisionsPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const filterDeptId = searchParams.get("departmentId");

  const [items, setItems] = useState<Division[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<Division | null>(null);
  const [formData, setFormData] = useState({ divisionName: "", divisionCode: "", departmentId: "", isActive: true });

  const loadData = async () => {
    setLoading(true);
    try {
      const [divRes, deptRes] = await Promise.all([getDivisions(), getDepartments()]);
      if (divRes?.data) {
        const all = Array.isArray(divRes.data) ? divRes.data : [];
        setItems(filterDeptId ? all.filter((d) => d.departmentId === Number(filterDeptId)) : all);
      }
      if (deptRes?.data) setDepartments(Array.isArray(deptRes.data) ? deptRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [filterDeptId]);

  const openAdd = () => {
    setEditItem(null);
    setFormData({ divisionName: "", divisionCode: "", departmentId: filterDeptId ?? "", isActive: true });
    setOpen(true);
  };

  const openEdit = (item: Division) => {
    setEditItem(item);
    setFormData({
      divisionName: item.divisionName,
      divisionCode: item.divisionCode ?? "",
      departmentId: String(item.departmentId),
      isActive: item.isActive,
    });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.divisionName || !formData.departmentId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { divisionId: editItem.divisionId } : {}),
        divisionName: formData.divisionName,
        divisionCode: formData.divisionCode,
        departmentId: Number(formData.departmentId),
        isActive: formData.isActive,
      };
      const res = await saveDivision(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Division updated." : "Division created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<Division>[] = [
    { accessorKey: "departmentName", header: "Department" },
    { accessorKey: "divisionCode", header: "Code" },
    { accessorKey: "divisionName", header: "Division Name" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Disabled"}</Badge> },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button>
          <Button size="sm" variant="ghost" onClick={() => router.push(`/offices?divisionId=${row.original.divisionId}`)} title="View Offices"><Eye className="h-3.5 w-3.5" /></Button>
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Divisions" breadcrumbs={[{ label: "Organogram" }, { label: "Divisions" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Divisions" description="Manage divisions under departments" breadcrumbs={[{ label: "Organogram" }, { label: "Departments", href: "/departments" }, { label: "Divisions" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Division</Button>} />
      <DataTable columns={columns} data={items} searchKey="divisionName" searchPlaceholder="Search divisions..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Division" : "Add Division"} isEdit={!!editItem} editWarning="Update the selected division below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Department *</Label>
            <Select value={formData.departmentId} onValueChange={(v) => setFormData({ ...formData, departmentId: v })}>
              <SelectTrigger><SelectValue placeholder="Select department" /></SelectTrigger>
              <SelectContent>{departments.map((d) => <SelectItem key={d.departmentId} value={String(d.departmentId)}>{d.departmentName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Division Code</Label><Input value={formData.divisionCode} onChange={(e) => setFormData({ ...formData, divisionCode: e.target.value })} /></div>
          <div className="space-y-2"><Label>Division Name *</Label><Input value={formData.divisionName} onChange={(e) => setFormData({ ...formData, divisionName: e.target.value })} /></div>
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
