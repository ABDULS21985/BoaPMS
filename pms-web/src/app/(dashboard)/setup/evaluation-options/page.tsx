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
import { getEvaluationOptions, saveEvaluationOptions } from "@/lib/api/performance";
import { EvaluationType } from "@/types/enums";
import type { EvaluationOption } from "@/types/performance";

const evalTypeLabels: Record<number, string> = { [EvaluationType.Timeliness]: "Timeliness", [EvaluationType.Quality]: "Quality", [EvaluationType.Output]: "Output" };

export default function EvaluationOptionsPage() {
  const [items, setItems] = useState<EvaluationOption[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<EvaluationOption | null>(null);
  const [formData, setFormData] = useState({ name: "", description: "", score: "", evaluationType: "" });

  const loadData = async () => {
    setLoading(true);
    try { const res = await getEvaluationOptions(); if (res?.data) setItems(Array.isArray(res.data) ? res.data : []); }
    catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ name: "", description: "", score: "", evaluationType: "" }); setOpen(true); };
  const openEdit = (item: EvaluationOption) => {
    setEditItem(item);
    setFormData({ name: item.name, description: item.description ?? "", score: String(item.score), evaluationType: String(item.evaluationType) });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.evaluationType || !formData.score) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const payload = {
        evaluationOptionId: editItem?.evaluationOptionId ?? undefined,
        name: formData.name, description: formData.description,
        score: Number(formData.score), evaluationType: Number(formData.evaluationType),
        recordStatus: editItem?.recordStatus ?? 8,
      };
      const res = await saveEvaluationOptions([payload]);
      if (res?.isSuccess) { toast.success(editItem ? "Option updated." : "Option created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<EvaluationOption>[] = [
    { accessorKey: "evaluationType", header: "Type", cell: ({ row }) => <Badge variant="outline">{evalTypeLabels[row.original.evaluationType] ?? row.original.evaluationType}</Badge> },
    { accessorKey: "name", header: "Statement" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "score", header: "Score" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Evaluation Options" breadcrumbs={[{ label: "Setup" }, { label: "Evaluation Options" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Evaluation Options" description="Configure evaluation statements and scores" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Evaluation Options" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Option</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search options..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Evaluation Option" : "Add Evaluation Option"} isEdit={!!editItem} editWarning="Update the selected option below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Evaluation Type *</Label>
            <Select value={formData.evaluationType} onValueChange={(v) => setFormData({ ...formData, evaluationType: v })}>
              <SelectTrigger><SelectValue placeholder="Select type" /></SelectTrigger>
              <SelectContent>
                {Object.entries(evalTypeLabels).map(([k, v]) => <SelectItem key={k} value={k}>{v}</SelectItem>)}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Evaluation Statement *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="space-y-2"><Label>Score *</Label><Input type="number" value={formData.score} onChange={(e) => setFormData({ ...formData, score: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
