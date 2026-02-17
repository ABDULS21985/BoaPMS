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
import { getRatings, saveRating } from "@/lib/api/competency";
import type { Rating } from "@/types/competency";

export default function RatingsPage() {
  const [items, setItems] = useState<Rating[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<Rating | null>(null);
  const [formData, setFormData] = useState({ name: "", value: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getRatings();
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ name: "", value: "" }); setOpen(true); };
  const openEdit = (item: Rating) => {
    setEditItem(item);
    setFormData({ name: item.name, value: String(item.value) });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.name || !formData.value) { toast.error("Name and value are required."); return; }
    setSaving(true);
    try {
      const payload = editItem
        ? { ratingId: editItem.ratingId, name: formData.name, value: Number(formData.value) }
        : { name: formData.name, value: Number(formData.value) };
      const res = await saveRating(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Rating updated." : "Rating created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<Rating>[] = [
    { accessorKey: "name", header: "Rating Name" },
    { accessorKey: "value", header: "Value", cell: ({ row }) => <Badge variant="outline">{row.original.value}</Badge> },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Ratings" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Ratings" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Ratings" description="Manage competency rating levels"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Ratings" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Rating</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search ratings..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Rating" : "Add Rating"} isEdit={!!editItem} editWarning="Update the rating below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Rating Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Value *</Label><Input type="number" value={formData.value} onChange={(e) => setFormData({ ...formData, value: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
