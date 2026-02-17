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
import { getPmsCompetencies, createPmsCompetency, updatePmsCompetency, getObjectiveCategories, type PmsCompetency } from "@/lib/api/performance";
import type { ObjectiveCategory } from "@/types/performance";

export default function PmsCompetenciesPage() {
  const [items, setItems] = useState<PmsCompetency[]>([]);
  const [categories, setCategories] = useState<ObjectiveCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<PmsCompetency | null>(null);
  const [formData, setFormData] = useState({ name: "", description: "", objectCategoryId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, catRes] = await Promise.all([getPmsCompetencies(), getObjectiveCategories()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (catRes?.data) setCategories(Array.isArray(catRes.data) ? catRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ name: "", description: "", objectCategoryId: "" }); setOpen(true); };
  const openEdit = (item: PmsCompetency) => {
    setEditItem(item);
    setFormData({ name: item.name, description: item.description ?? "", objectCategoryId: item.objectCategoryId });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.objectCategoryId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const res = editItem
        ? await updatePmsCompetency({ ...formData, pmsCompetencyId: editItem.pmsCompetencyId, recordStatus: editItem.recordStatus })
        : await createPmsCompetency(formData);
      if (res?.isSuccess) { toast.success(editItem ? "Competency updated." : "Competency created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const getCategoryName = (id: string) => categories.find((c) => c.objectiveCategoryId === id)?.name ?? id;

  const columns: ColumnDef<PmsCompetency>[] = [
    { accessorKey: "name", header: "Competency Name" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "objectCategoryId", header: "Objective Category", cell: ({ row }) => getCategoryName(row.original.objectCategoryId) },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="PMS Competencies" breadcrumbs={[{ label: "Setup" }, { label: "PMS Competencies" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="PMS Competencies" description="Manage PMS competency dimensions" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "PMS Competencies" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Competency</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search competencies..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit PMS Competency" : "Add PMS Competency"} isEdit={!!editItem} editWarning="Update the selected competency below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Objective Category *</Label>
            <Select value={formData.objectCategoryId} onValueChange={(v) => setFormData({ ...formData, objectCategoryId: v })}>
              <SelectTrigger><SelectValue placeholder="Select category" /></SelectTrigger>
              <SelectContent>{categories.map((c) => <SelectItem key={c.objectiveCategoryId} value={c.objectiveCategoryId}>{c.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Competency Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
