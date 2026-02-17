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
import { getStrategies, getObjectiveCategories, getEnterpriseObjectives, createEnterpriseObjective, updateEnterpriseObjective } from "@/lib/api/performance";
import type { Strategy, ObjectiveCategory, EnterpriseObjective } from "@/types/performance";

export default function EnterpriseObjectivesPage() {
  const [items, setItems] = useState<EnterpriseObjective[]>([]);
  const [strategies, setStrategies] = useState<Strategy[]>([]);
  const [categories, setCategories] = useState<ObjectiveCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<EnterpriseObjective | null>(null);
  const [formData, setFormData] = useState({ name: "", description: "", kpi: "", target: "", strategyId: "", enterpriseObjectivesCategoryId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, stratRes, catRes] = await Promise.all([getEnterpriseObjectives(), getStrategies(), getObjectiveCategories()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (stratRes?.data) setStrategies(Array.isArray(stratRes.data) ? stratRes.data : []);
      if (catRes?.data) setCategories(Array.isArray(catRes.data) ? catRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ name: "", description: "", kpi: "", target: "", strategyId: "", enterpriseObjectivesCategoryId: "" }); setOpen(true); };
  const openEdit = (item: EnterpriseObjective) => {
    setEditItem(item);
    setFormData({ name: item.name, description: item.description ?? "", kpi: item.kpi ?? "", target: item.target ?? "", strategyId: item.strategyId, enterpriseObjectivesCategoryId: item.enterpriseObjectivesCategoryId });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.strategyId || !formData.enterpriseObjectivesCategoryId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const res = editItem
        ? await updateEnterpriseObjective({ ...formData, enterpriseObjectiveId: editItem.enterpriseObjectiveId })
        : await createEnterpriseObjective(formData);
      if (res?.isSuccess) { toast.success(editItem ? "Objective updated." : "Objective created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<EnterpriseObjective>[] = [
    { accessorKey: "name", header: "Objective Name" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "kpi", header: "KPI" },
    { accessorKey: "target", header: "Target" },
    { accessorKey: "categoryName", header: "Category" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Enterprise Objectives" breadcrumbs={[{ label: "Setup" }, { label: "Enterprise Objectives" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Enterprise Objectives" description="Manage enterprise-level objectives" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Enterprise Objectives" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Objective</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search objectives..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Enterprise Objective" : "Add Enterprise Objective"} isEdit={!!editItem} editWarning="Update the selected objective below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Strategy *</Label>
            <Select value={formData.strategyId} onValueChange={(v) => setFormData({ ...formData, strategyId: v })}>
              <SelectTrigger><SelectValue placeholder="Select strategy" /></SelectTrigger>
              <SelectContent>{strategies.map((s) => <SelectItem key={s.strategyId} value={s.strategyId}>{s.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Objective Category *</Label>
            <Select value={formData.enterpriseObjectivesCategoryId} onValueChange={(v) => setFormData({ ...formData, enterpriseObjectivesCategoryId: v })}>
              <SelectTrigger><SelectValue placeholder="Select category" /></SelectTrigger>
              <SelectContent>{categories.map((c) => <SelectItem key={c.objectiveCategoryId} value={c.objectiveCategoryId}>{c.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Objective Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description *</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="space-y-2"><Label>KPI *</Label><Input value={formData.kpi} onChange={(e) => setFormData({ ...formData, kpi: e.target.value })} /></div>
          <div className="space-y-2"><Label>Target *</Label><Input value={formData.target} onChange={(e) => setFormData({ ...formData, target: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
