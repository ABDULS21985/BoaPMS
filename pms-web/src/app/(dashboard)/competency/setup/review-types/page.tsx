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
import { getReviewTypes, saveReviewType } from "@/lib/api/competency";
import type { ReviewType } from "@/types/competency";

export default function ReviewTypesPage() {
  const [items, setItems] = useState<ReviewType[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<ReviewType | null>(null);
  const [formData, setFormData] = useState({ reviewTypeName: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getReviewTypes();
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ reviewTypeName: "" }); setOpen(true); };
  const openEdit = (item: ReviewType) => {
    setEditItem(item);
    setFormData({ reviewTypeName: item.reviewTypeName });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.reviewTypeName) { toast.error("Review type name is required."); return; }
    setSaving(true);
    try {
      const payload = editItem
        ? { reviewTypeId: editItem.reviewTypeId, ...formData }
        : formData;
      const res = await saveReviewType(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Review type updated." : "Review type created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<ReviewType>[] = [
    { accessorKey: "reviewTypeName", header: "Review Type" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Review Types" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Review Types" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Review Types" description="Manage competency review types"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Review Types" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Review Type</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="reviewTypeName" searchPlaceholder="Search review types..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Review Type" : "Add Review Type"} isEdit={!!editItem} editWarning="Update the review type below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Review Type Name *</Label><Input value={formData.reviewTypeName} onChange={(e) => setFormData({ ...formData, reviewTypeName: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
