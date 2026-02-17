"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getCategoryGradings, saveCategoryGrading, getCompetencyCategories, getReviewTypes } from "@/lib/api/competency";
import type { CompetencyCategoryGrading, CompetencyCategory, ReviewType } from "@/types/competency";

export default function CompetencyGradingPage() {
  const [items, setItems] = useState<CompetencyCategoryGrading[]>([]);
  const [categories, setCategories] = useState<CompetencyCategory[]>([]);
  const [reviewTypes, setReviewTypes] = useState<ReviewType[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<CompetencyCategoryGrading | null>(null);
  const [formData, setFormData] = useState({ competencyCategoryId: "", reviewTypeId: "", weightPercentage: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, catRes, rtRes] = await Promise.all([getCategoryGradings(), getCompetencyCategories(), getReviewTypes()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (catRes?.data) setCategories(Array.isArray(catRes.data) ? catRes.data : []);
      if (rtRes?.data) setReviewTypes(Array.isArray(rtRes.data) ? rtRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ competencyCategoryId: "", reviewTypeId: "", weightPercentage: "" }); setOpen(true); };
  const openEdit = (item: CompetencyCategoryGrading) => {
    setEditItem(item);
    setFormData({ competencyCategoryId: String(item.competencyCategoryId), reviewTypeId: String(item.reviewTypeId), weightPercentage: String(item.weightPercentage) });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.competencyCategoryId || !formData.reviewTypeId || !formData.weightPercentage) { toast.error("All fields are required."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { competencyCategoryGradingId: editItem.competencyCategoryGradingId } : {}),
        competencyCategoryId: Number(formData.competencyCategoryId),
        reviewTypeId: Number(formData.reviewTypeId),
        weightPercentage: Number(formData.weightPercentage),
      };
      const res = await saveCategoryGrading(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Grading updated." : "Grading created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<CompetencyCategoryGrading>[] = [
    { accessorKey: "competencyCategoryName", header: "Category" },
    { accessorKey: "reviewTypeName", header: "Review Type" },
    { accessorKey: "weightPercentage", header: "Weight %", cell: ({ row }) => <Badge variant="outline">{row.original.weightPercentage}%</Badge> },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Competency Grading" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Grading" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Competency Grading" description="Define weight distribution across categories and review types"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Grading" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Grading</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="competencyCategoryName" searchPlaceholder="Search gradings..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Grading" : "Add Grading"} isEdit={!!editItem}>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Category *</Label>
            <Select value={formData.competencyCategoryId} onValueChange={(v) => setFormData({ ...formData, competencyCategoryId: v })}>
              <SelectTrigger><SelectValue placeholder="Select category" /></SelectTrigger>
              <SelectContent>{categories.map((c) => <SelectItem key={c.competencyCategoryId} value={String(c.competencyCategoryId)}>{c.categoryName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Review Type *</Label>
            <Select value={formData.reviewTypeId} onValueChange={(v) => setFormData({ ...formData, reviewTypeId: v })}>
              <SelectTrigger><SelectValue placeholder="Select review type" /></SelectTrigger>
              <SelectContent>{reviewTypes.map((r) => <SelectItem key={r.reviewTypeId} value={String(r.reviewTypeId)}>{r.reviewTypeName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Weight Percentage *</Label><Input type="number" value={formData.weightPercentage} onChange={(e) => setFormData({ ...formData, weightPercentage: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
