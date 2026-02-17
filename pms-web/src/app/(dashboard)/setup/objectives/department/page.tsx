"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getEnterpriseObjectives, getDepartmentObjectives, createDepartmentObjective, updateDepartmentObjective } from "@/lib/api/performance";
import { getDepartments } from "@/lib/api/organogram";
import type { EnterpriseObjective, DepartmentObjective } from "@/types/performance";
import type { Department } from "@/types/organogram";

export default function DepartmentObjectivesPage() {
  const [items, setItems] = useState<DepartmentObjective[]>([]);
  const [enterpriseObjs, setEnterpriseObjs] = useState<EnterpriseObjective[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<DepartmentObjective | null>(null);
  const [formData, setFormData] = useState({ name: "", description: "", kpi: "", target: "", enterpriseObjectiveId: "", departmentId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, eoRes, deptRes] = await Promise.all([getDepartmentObjectives(), getEnterpriseObjectives(), getDepartments()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (eoRes?.data) setEnterpriseObjs(Array.isArray(eoRes.data) ? eoRes.data : []);
      if (deptRes?.data) setDepartments(Array.isArray(deptRes.data) ? deptRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ name: "", description: "", kpi: "", target: "", enterpriseObjectiveId: "", departmentId: "" }); setOpen(true); };
  const openEdit = (item: DepartmentObjective) => {
    setEditItem(item);
    setFormData({ name: item.name, description: item.description ?? "", kpi: item.kpi ?? "", target: item.target ?? "", enterpriseObjectiveId: item.enterpriseObjectiveId, departmentId: String(item.departmentId) });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.departmentId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const payload = { ...formData, departmentId: Number(formData.departmentId) };
      const res = editItem ? await updateDepartmentObjective({ ...payload, departmentObjectiveId: editItem.departmentObjectiveId }) : await createDepartmentObjective(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Objective updated." : "Objective created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const getDeptName = (id: number) => departments.find((d) => d.departmentId === id)?.departmentName ?? String(id);

  const columns: ColumnDef<DepartmentObjective>[] = [
    { accessorKey: "name", header: "Objective Name" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "departmentId", header: "Department", cell: ({ row }) => getDeptName(row.original.departmentId) },
    { accessorKey: "kpi", header: "KPI" },
    { accessorKey: "target", header: "Target" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Department Objectives" breadcrumbs={[{ label: "Setup" }, { label: "Department Objectives" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Department Objectives" description="Manage department-level objectives" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Department Objectives" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Objective</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search objectives..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Department Objective" : "Add Department Objective"} isEdit={!!editItem} editWarning="Update the selected objective below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Enterprise Objective</Label>
            <Select value={formData.enterpriseObjectiveId} onValueChange={(v) => setFormData({ ...formData, enterpriseObjectiveId: v })}>
              <SelectTrigger><SelectValue placeholder="Select enterprise objective" /></SelectTrigger>
              <SelectContent>{enterpriseObjs.map((o) => <SelectItem key={o.enterpriseObjectiveId} value={o.enterpriseObjectiveId}>{o.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Department *</Label>
            <Select value={formData.departmentId} onValueChange={(v) => setFormData({ ...formData, departmentId: v })}>
              <SelectTrigger><SelectValue placeholder="Select department" /></SelectTrigger>
              <SelectContent>{departments.map((d) => <SelectItem key={d.departmentId} value={String(d.departmentId)}>{d.departmentName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Objective Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="space-y-2"><Label>KPI</Label><Input value={formData.kpi} onChange={(e) => setFormData({ ...formData, kpi: e.target.value })} /></div>
          <div className="space-y-2"><Label>Target</Label><Input value={formData.target} onChange={(e) => setFormData({ ...formData, target: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
