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
import { getTrainingTypes, saveTrainingType } from "@/lib/api/competency";
import type { TrainingType } from "@/types/competency";

export default function TrainingTypesPage() {
  const [items, setItems] = useState<TrainingType[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<TrainingType | null>(null);
  const [formData, setFormData] = useState({ trainingTypeName: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getTrainingTypes();
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ trainingTypeName: "" }); setOpen(true); };
  const openEdit = (item: TrainingType) => {
    setEditItem(item);
    setFormData({ trainingTypeName: item.trainingTypeName });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.trainingTypeName) { toast.error("Training type name is required."); return; }
    setSaving(true);
    try {
      const payload = editItem
        ? { trainingTypeId: editItem.trainingTypeId, ...formData }
        : formData;
      const res = await saveTrainingType(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Training type updated." : "Training type created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<TrainingType>[] = [
    { accessorKey: "trainingTypeName", header: "Training Type" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Training Types" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Training Types" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Training Types" description="Manage development intervention / training types"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Training Types" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Training Type</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="trainingTypeName" searchPlaceholder="Search training types..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Training Type" : "Add Training Type"} isEdit={!!editItem} editWarning="Update the training type below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Training Type Name *</Label><Input value={formData.trainingTypeName} onChange={(e) => setFormData({ ...formData, trainingTypeName: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
