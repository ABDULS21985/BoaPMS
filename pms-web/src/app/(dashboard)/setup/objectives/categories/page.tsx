"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getObjectiveCategories, createObjectiveCategory, updateObjectiveCategory } from "@/lib/api/performance";
import type { ObjectiveCategory } from "@/types/performance";

export default function ObjectiveCategoriesPage() {
  const [items, setItems] = useState<ObjectiveCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<ObjectiveCategory | null>(null);
  const [formData, setFormData] = useState({ name: "", description: "", isActive: true });

  const loadData = async () => {
    setLoading(true);
    try { const res = await getObjectiveCategories(); if (res?.data) setItems(Array.isArray(res.data) ? res.data : []); }
    catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ name: "", description: "", isActive: true }); setOpen(true); };
  const openEdit = (item: ObjectiveCategory) => {
    setEditItem(item);
    setFormData({ name: item.name, description: item.description ?? "", isActive: item.isActive });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name) { toast.error("Name is required."); return; }
    setSaving(true);
    try {
      const res = editItem
        ? await updateObjectiveCategory({ ...formData, objectiveCategoryId: editItem.objectiveCategoryId })
        : await createObjectiveCategory(formData);
      if (res?.isSuccess) { toast.success(editItem ? "Category updated." : "Category created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<ObjectiveCategory>[] = [
    { accessorKey: "name", header: "Category Name" },
    { accessorKey: "description", header: "Description" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "secondary"}>{row.original.isActive ? "Active" : "Disabled"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Objective Categories" breadcrumbs={[{ label: "Setup" }, { label: "Objective Categories" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Objective Categories" description="Manage objective category definitions" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Objective Categories" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Category</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search categories..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Category" : "Add Category"} isEdit={!!editItem} editWarning="Update the selected category below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Category Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="flex items-center gap-3"><Switch checked={formData.isActive} onCheckedChange={(v) => setFormData({ ...formData, isActive: v })} /><Label>Active</Label></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
