"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2, CheckCircle, XCircle } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getCompetencies, saveCompetency, approveCompetency, rejectCompetency, getCompetencyCategories } from "@/lib/api/competency";
import type { Competency, CompetencyCategory } from "@/types/competency";

export default function CompetenciesPage() {
  const [items, setItems] = useState<Competency[]>([]);
  const [categories, setCategories] = useState<CompetencyCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<Competency | null>(null);
  const [formData, setFormData] = useState({ competencyName: "", description: "", competencyCategoryId: "" });
  const [approveItem, setApproveItem] = useState<Competency | null>(null);
  const [rejectItem, setRejectItem] = useState<Competency | null>(null);
  const [rejectReason, setRejectReason] = useState("");

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, catRes] = await Promise.all([getCompetencies(), getCompetencyCategories()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (catRes?.data) setCategories(Array.isArray(catRes.data) ? catRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ competencyName: "", description: "", competencyCategoryId: "" }); setOpen(true); };
  const openEdit = (item: Competency) => {
    setEditItem(item);
    setFormData({ competencyName: item.competencyName, description: item.description ?? "", competencyCategoryId: String(item.competencyCategoryId) });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.competencyName || !formData.competencyCategoryId) { toast.error("Name and category are required."); return; }
    setSaving(true);
    try {
      const payload = editItem
        ? { competencyId: editItem.competencyId, competencyName: formData.competencyName, description: formData.description, competencyCategoryId: Number(formData.competencyCategoryId) }
        : { competencyName: formData.competencyName, description: formData.description, competencyCategoryId: Number(formData.competencyCategoryId) };
      const res = await saveCompetency(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Competency updated." : "Competency created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const handleApprove = async () => {
    if (!approveItem) return;
    try {
      const res = await approveCompetency({ competencyId: approveItem.competencyId });
      if (res?.isSuccess) { toast.success("Competency approved."); loadData(); }
      else toast.error(res?.message || "Approval failed.");
    } catch { toast.error("An error occurred."); } finally { setApproveItem(null); }
  };

  const handleReject = async (reason?: string) => {
    if (!rejectItem) return;
    try {
      const res = await rejectCompetency({ competencyId: rejectItem.competencyId, rejectionReason: reason ?? rejectReason });
      if (res?.isSuccess) { toast.success("Competency rejected."); loadData(); }
      else toast.error(res?.message || "Rejection failed.");
    } catch { toast.error("An error occurred."); } finally { setRejectItem(null); setRejectReason(""); }
  };

  const columns: ColumnDef<Competency>[] = [
    { accessorKey: "competencyName", header: "Competency" },
    { accessorKey: "competencyCategoryName", header: "Category" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    {
      accessorKey: "isApproved", header: "Status",
      cell: ({ row }) => {
        if (row.original.isRejected) return <Badge variant="destructive">Rejected</Badge>;
        if (row.original.isApproved) return <Badge variant="default">Approved</Badge>;
        return <Badge variant="secondary">Pending</Badge>;
      },
    },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button>
          {!row.original.isApproved && !row.original.isRejected && (
            <>
              <Button size="sm" variant="ghost" onClick={() => setApproveItem(row.original)}><CheckCircle className="h-3.5 w-3.5 text-green-600" /></Button>
              <Button size="sm" variant="ghost" onClick={() => setRejectItem(row.original)}><XCircle className="h-3.5 w-3.5 text-red-600" /></Button>
            </>
          )}
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Competencies" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Competencies" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Competencies" description="Manage competency definitions with approval workflow"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Competencies" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Competency</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="competencyName" searchPlaceholder="Search competencies..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Competency" : "Add Competency"} isEdit={!!editItem} editWarning="Update the competency below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Competency Name *</Label><Input value={formData.competencyName} onChange={(e) => setFormData({ ...formData, competencyName: e.target.value })} /></div>
          <div className="space-y-2">
            <Label>Category *</Label>
            <Select value={formData.competencyCategoryId} onValueChange={(v) => setFormData({ ...formData, competencyCategoryId: v })}>
              <SelectTrigger><SelectValue placeholder="Select category" /></SelectTrigger>
              <SelectContent>{categories.map((c) => <SelectItem key={c.competencyCategoryId} value={String(c.competencyCategoryId)}>{c.categoryName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Description</Label><Textarea value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} rows={3} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>

      <ConfirmationDialog open={!!approveItem} onOpenChange={() => setApproveItem(null)} title="Approve Competency" description={`Approve "${approveItem?.competencyName}"?`} confirmLabel="Approve" onConfirm={handleApprove} />

      <ConfirmationDialog open={!!rejectItem} onOpenChange={() => { setRejectItem(null); setRejectReason(""); }} title="Reject Competency" description={`Reject "${rejectItem?.competencyName}"? Please provide a reason.`} confirmLabel="Reject" variant="destructive" showReasonInput reasonLabel="Rejection Reason" onConfirm={handleReject} />
    </div>
  );
}
