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
import { getCompetencyCategories, saveCompetencyCategory } from "@/lib/api/competency";
import type { CompetencyCategory } from "@/types/competency";

export default function CompetencyCategoriesPage() {
  const [items, setItems] = useState<CompetencyCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<CompetencyCategory | null>(null);
  const [formData, setFormData] = useState({ categoryName: "", isTechnical: false });

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getCompetencyCategories();
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ categoryName: "", isTechnical: false }); setOpen(true); };
  const openEdit = (item: CompetencyCategory) => {
    setEditItem(item);
    setFormData({ categoryName: item.categoryName, isTechnical: item.isTechnical });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.categoryName) { toast.error("Category name is required."); return; }
    setSaving(true);
    try {
      const payload = editItem
        ? { competencyCategoryId: editItem.competencyCategoryId, ...formData }
        : formData;
      const res = await saveCompetencyCategory(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Category updated." : "Category created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<CompetencyCategory>[] = [
    { accessorKey: "categoryName", header: "Category Name" },
    { accessorKey: "isTechnical", header: "Type", cell: ({ row }) => <Badge variant={row.original.isTechnical ? "default" : "secondary"}>{row.original.isTechnical ? "Technical" : "Behavioral"}</Badge> },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Competency Categories" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Categories" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Competency Categories" description="Manage competency categories (Technical / Behavioral)"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Categories" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Category</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="categoryName" searchPlaceholder="Search categories..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Category" : "Add Category"} isEdit={!!editItem} editWarning="Update the category details below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Category Name *</Label><Input value={formData.categoryName} onChange={(e) => setFormData({ ...formData, categoryName: e.target.value })} /></div>
          <div className="flex items-center gap-3">
            <Switch checked={formData.isTechnical} onCheckedChange={(v) => setFormData({ ...formData, isTechnical: v })} />
            <Label>Technical Category</Label>
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
