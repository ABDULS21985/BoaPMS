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
import { getDivisionObjectives, getOfficeObjectives, createOfficeObjective, updateOfficeObjective } from "@/lib/api/performance";
import { getDepartments, getDivisionsByDepartment, getOfficesByDivision } from "@/lib/api/organogram";
import { getJobGradeGroups } from "@/lib/api/competency";
import type { DivisionObjective, OfficeObjective } from "@/types/performance";
import type { Department, Division, Office } from "@/types/organogram";
import type { JobGradeGroup } from "@/types/competency";

export default function OfficeObjectivesPage() {
  const [items, setItems] = useState<OfficeObjective[]>([]);
  const [divObjs, setDivObjs] = useState<DivisionObjective[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [divisions, setDivisions] = useState<Division[]>([]);
  const [offices, setOffices] = useState<Office[]>([]);
  const [gradeGroups, setGradeGroups] = useState<JobGradeGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<OfficeObjective | null>(null);
  const [selectedDept, setSelectedDept] = useState("");
  const [selectedDiv, setSelectedDiv] = useState("");
  const [formData, setFormData] = useState({ name: "", description: "", kpi: "", target: "", divisionObjectiveId: "", officeId: "", jobGradeGroupId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, doRes, deptRes, ggRes] = await Promise.all([getOfficeObjectives(), getDivisionObjectives(), getDepartments(), getJobGradeGroups()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (doRes?.data) setDivObjs(Array.isArray(doRes.data) ? doRes.data : []);
      if (deptRes?.data) setDepartments(Array.isArray(deptRes.data) ? deptRes.data : []);
      if (ggRes?.data) setGradeGroups(Array.isArray(ggRes.data) ? ggRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const onDeptChange = async (deptId: string) => {
    setSelectedDept(deptId); setSelectedDiv(""); setOffices([]);
    setFormData((f) => ({ ...f, divisionObjectiveId: "", officeId: "" }));
    try { const res = await getDivisionsByDepartment(Number(deptId)); if (res?.data) setDivisions(Array.isArray(res.data) ? res.data : []); }
    catch { setDivisions([]); }
  };

  const onDivChange = async (divId: string) => {
    setSelectedDiv(divId);
    setFormData((f) => ({ ...f, officeId: "" }));
    try { const res = await getOfficesByDivision(Number(divId)); if (res?.data) setOffices(Array.isArray(res.data) ? res.data : []); }
    catch { setOffices([]); }
  };

  const filteredDivObjs = selectedDiv ? divObjs.filter((d) => d.divisionId === Number(selectedDiv)) : divObjs;

  const openAdd = () => { setEditItem(null); setSelectedDept(""); setSelectedDiv(""); setDivisions([]); setOffices([]); setFormData({ name: "", description: "", kpi: "", target: "", divisionObjectiveId: "", officeId: "", jobGradeGroupId: "" }); setOpen(true); };
  const openEdit = (item: OfficeObjective) => {
    setEditItem(item);
    setFormData({ name: item.name, description: item.description ?? "", kpi: item.kpi ?? "", target: item.target ?? "", divisionObjectiveId: item.divisionObjectiveId, officeId: String(item.officeId), jobGradeGroupId: String(item.jobGradeGroupId) });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.officeId || !formData.jobGradeGroupId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const payload = { ...formData, officeId: Number(formData.officeId), jobGradeGroupId: Number(formData.jobGradeGroupId) };
      const res = editItem ? await updateOfficeObjective({ ...payload, officeObjectiveId: editItem.officeObjectiveId }) : await createOfficeObjective(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Objective updated." : "Objective created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<OfficeObjective>[] = [
    { accessorKey: "name", header: "Objective Name" },
    { accessorKey: "officeName", header: "Office" },
    { accessorKey: "kpi", header: "KPI" },
    { accessorKey: "target", header: "Target" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Office Objectives" breadcrumbs={[{ label: "Setup" }, { label: "Office Objectives" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Office Objectives" description="Manage office-level objectives" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Office Objectives" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Objective</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search objectives..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Office Objective" : "Add Office Objective"} isEdit={!!editItem} editWarning="Update the selected objective below." className="sm:max-w-xl overflow-y-auto">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Department</Label><Select value={selectedDept} onValueChange={onDeptChange}><SelectTrigger><SelectValue placeholder="Select department" /></SelectTrigger><SelectContent>{departments.map((d) => <SelectItem key={d.departmentId} value={String(d.departmentId)}>{d.departmentName}</SelectItem>)}</SelectContent></Select></div>
          <div className="space-y-2"><Label>Division</Label><Select value={selectedDiv} onValueChange={onDivChange}><SelectTrigger><SelectValue placeholder="Select division" /></SelectTrigger><SelectContent>{divisions.map((d) => <SelectItem key={d.divisionId} value={String(d.divisionId)}>{d.divisionName}</SelectItem>)}</SelectContent></Select></div>
          <div className="space-y-2"><Label>Division Objective</Label><Select value={formData.divisionObjectiveId} onValueChange={(v) => setFormData({ ...formData, divisionObjectiveId: v })}><SelectTrigger><SelectValue placeholder="Select division objective" /></SelectTrigger><SelectContent>{filteredDivObjs.map((o) => <SelectItem key={o.divisionObjectiveId} value={o.divisionObjectiveId}>{o.name}</SelectItem>)}</SelectContent></Select></div>
          <div className="space-y-2"><Label>Office *</Label><Select value={formData.officeId} onValueChange={(v) => setFormData({ ...formData, officeId: v })}><SelectTrigger><SelectValue placeholder="Select office" /></SelectTrigger><SelectContent>{offices.map((o) => <SelectItem key={o.officeId} value={String(o.officeId)}>{o.officeName}</SelectItem>)}</SelectContent></Select></div>
          <div className="space-y-2"><Label>Grade Group *</Label><Select value={formData.jobGradeGroupId} onValueChange={(v) => setFormData({ ...formData, jobGradeGroupId: v })}><SelectTrigger><SelectValue placeholder="Select grade group" /></SelectTrigger><SelectContent>{gradeGroups.map((g) => <SelectItem key={g.jobGradeGroupId} value={String(g.jobGradeGroupId)}>{g.groupName}</SelectItem>)}</SelectContent></Select></div>
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
