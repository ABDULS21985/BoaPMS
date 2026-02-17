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
import { getDepartmentObjectives, getDivisionObjectives, createDivisionObjective, updateDivisionObjective } from "@/lib/api/performance";
import { getDepartments, getDivisionsByDepartment } from "@/lib/api/organogram";
import type { DepartmentObjective, DivisionObjective } from "@/types/performance";
import type { Department, Division } from "@/types/organogram";

export default function DivisionObjectivesPage() {
  const [items, setItems] = useState<DivisionObjective[]>([]);
  const [deptObjs, setDeptObjs] = useState<DepartmentObjective[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [divisions, setDivisions] = useState<Division[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<DivisionObjective | null>(null);
  const [selectedDept, setSelectedDept] = useState("");
  const [formData, setFormData] = useState({ name: "", description: "", kpi: "", target: "", departmentObjectiveId: "", divisionId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, doRes, deptRes] = await Promise.all([getDivisionObjectives(), getDepartmentObjectives(), getDepartments()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (doRes?.data) setDeptObjs(Array.isArray(doRes.data) ? doRes.data : []);
      if (deptRes?.data) setDepartments(Array.isArray(deptRes.data) ? deptRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const onDeptChange = async (deptId: string) => {
    setSelectedDept(deptId);
    setFormData((f) => ({ ...f, departmentObjectiveId: "", divisionId: "" }));
    try { const res = await getDivisionsByDepartment(Number(deptId)); if (res?.data) setDivisions(Array.isArray(res.data) ? res.data : []); }
    catch { setDivisions([]); }
  };

  const filteredDeptObjs = selectedDept ? deptObjs.filter((d) => d.departmentId === Number(selectedDept)) : deptObjs;

  const openAdd = () => { setEditItem(null); setSelectedDept(""); setDivisions([]); setFormData({ name: "", description: "", kpi: "", target: "", departmentObjectiveId: "", divisionId: "" }); setOpen(true); };
  const openEdit = (item: DivisionObjective) => {
    setEditItem(item);
    const parentDeptObj = deptObjs.find((d) => d.departmentObjectiveId === item.departmentObjectiveId);
    const deptId = parentDeptObj ? String(parentDeptObj.departmentId) : "";
    setSelectedDept(deptId);
    if (deptId) getDivisionsByDepartment(Number(deptId)).then((r) => { if (r?.data) setDivisions(Array.isArray(r.data) ? r.data : []); });
    setFormData({ name: item.name, description: item.description ?? "", kpi: item.kpi ?? "", target: item.target ?? "", departmentObjectiveId: item.departmentObjectiveId, divisionId: String(item.divisionId) });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.divisionId || !formData.departmentObjectiveId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const payload = { ...formData, divisionId: Number(formData.divisionId), departmentId: Number(selectedDept) };
      const res = editItem ? await updateDivisionObjective({ ...payload, divisionObjectiveId: editItem.divisionObjectiveId }) : await createDivisionObjective(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Objective updated." : "Objective created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const getDivName = (id: number) => divisions.find((d) => d.divisionId === id)?.divisionName ?? String(id);

  const columns: ColumnDef<DivisionObjective>[] = [
    { accessorKey: "name", header: "Objective Name" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "divisionName", header: "Division" },
    { accessorKey: "kpi", header: "KPI" },
    { accessorKey: "target", header: "Target" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Division Objectives" breadcrumbs={[{ label: "Setup" }, { label: "Division Objectives" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Division Objectives" description="Manage division-level objectives" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Division Objectives" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Objective</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search objectives..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Division Objective" : "Add Division Objective"} isEdit={!!editItem} editWarning="Update the selected objective below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Department *</Label>
            <Select value={selectedDept} onValueChange={onDeptChange}>
              <SelectTrigger><SelectValue placeholder="Select department" /></SelectTrigger>
              <SelectContent>{departments.map((d) => <SelectItem key={d.departmentId} value={String(d.departmentId)}>{d.departmentName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Department Objective *</Label>
            <Select value={formData.departmentObjectiveId} onValueChange={(v) => setFormData({ ...formData, departmentObjectiveId: v })}>
              <SelectTrigger><SelectValue placeholder="Select dept objective" /></SelectTrigger>
              <SelectContent>{filteredDeptObjs.map((o) => <SelectItem key={o.departmentObjectiveId} value={o.departmentObjectiveId}>{o.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Division *</Label>
            <Select value={formData.divisionId} onValueChange={(v) => setFormData({ ...formData, divisionId: v })}>
              <SelectTrigger><SelectValue placeholder="Select division" /></SelectTrigger>
              <SelectContent>{divisions.map((d) => <SelectItem key={d.divisionId} value={String(d.divisionId)}>{d.divisionName}</SelectItem>)}</SelectContent>
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
