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
import { getFeedbackQuestionnaires, saveFeedbackQuestionnaires, getPmsCompetencies, type PmsCompetency } from "@/lib/api/performance";
import type { FeedbackQuestionnaire } from "@/types/performance";

export default function QuestionnairesPage() {
  const [items, setItems] = useState<FeedbackQuestionnaire[]>([]);
  const [competencies, setCompetencies] = useState<PmsCompetency[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<FeedbackQuestionnaire | null>(null);
  const [formData, setFormData] = useState({ question: "", description: "", pmsCompetencyId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, compRes] = await Promise.all([getFeedbackQuestionnaires(), getPmsCompetencies()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (compRes?.data) setCompetencies(Array.isArray(compRes.data) ? compRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ question: "", description: "", pmsCompetencyId: "" }); setOpen(true); };
  const openEdit = (item: FeedbackQuestionnaire) => {
    setEditItem(item);
    setFormData({ question: item.question, description: item.description ?? "", pmsCompetencyId: item.pmsCompetencyId ?? "" });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.question) { toast.error("Question is required."); return; }
    setSaving(true);
    try {
      const payload = { feedbackQuestionnaireId: editItem?.feedbackQuestionaireId, question: formData.question, description: formData.description, pmsCompetencyId: formData.pmsCompetencyId || undefined, recordStatus: 8 };
      const res = await saveFeedbackQuestionnaires([payload]);
      if (res?.isSuccess) { toast.success(editItem ? "Questionnaire updated." : "Questionnaire created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<FeedbackQuestionnaire>[] = [
    { accessorKey: "pmsCompetencyName", header: "Competency" },
    { accessorKey: "question", header: "Question" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Feedback Questionnaires" breadcrumbs={[{ label: "Setup" }, { label: "Questionnaires" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Feedback Questionnaires" description="Manage 360 feedback questionnaires" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Questionnaires" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Questionnaire</Button>} />
      <DataTable columns={columns} data={items} searchKey="question" searchPlaceholder="Search questionnaires..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Questionnaire" : "Add Questionnaire"} isEdit={!!editItem} editWarning="Update the selected questionnaire below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Competency</Label>
            <Select value={formData.pmsCompetencyId} onValueChange={(v) => setFormData({ ...formData, pmsCompetencyId: v })}>
              <SelectTrigger><SelectValue placeholder="Select competency" /></SelectTrigger>
              <SelectContent>{competencies.map((c) => <SelectItem key={c.pmsCompetencyId} value={c.pmsCompetencyId}>{c.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Question *</Label><Input value={formData.question} onChange={(e) => setFormData({ ...formData, question: e.target.value })} /></div>
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
