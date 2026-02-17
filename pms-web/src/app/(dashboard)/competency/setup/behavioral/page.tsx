"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getBehavioralCompetencies, saveBehavioralCompetency, getCompetencies, getJobGradeGroups, getRatings } from "@/lib/api/competency";
import type { BehavioralCompetency, Competency, JobGradeGroup, Rating } from "@/types/competency";

export default function BehavioralCompetenciesPage() {
  const [items, setItems] = useState<BehavioralCompetency[]>([]);
  const [competencies, setCompetencies] = useState<Competency[]>([]);
  const [gradeGroups, setGradeGroups] = useState<JobGradeGroup[]>([]);
  const [ratings, setRatings] = useState<Rating[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<BehavioralCompetency | null>(null);
  const [formData, setFormData] = useState({ competencyId: "", jobGradeGroupId: "", ratingId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, compRes, grpRes, ratRes] = await Promise.all([
        getBehavioralCompetencies(), getCompetencies(), getJobGradeGroups(), getRatings(),
      ]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (compRes?.data) setCompetencies(Array.isArray(compRes.data) ? compRes.data : []);
      if (grpRes?.data) setGradeGroups(Array.isArray(grpRes.data) ? grpRes.data : []);
      if (ratRes?.data) setRatings(Array.isArray(ratRes.data) ? ratRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ competencyId: "", jobGradeGroupId: "", ratingId: "" }); setOpen(true); };
  const openEdit = (item: BehavioralCompetency) => {
    setEditItem(item);
    setFormData({ competencyId: String(item.competencyId), jobGradeGroupId: String(item.jobGradeGroupId), ratingId: item.ratingId ? String(item.ratingId) : "" });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.competencyId || !formData.jobGradeGroupId) { toast.error("Competency and grade group are required."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { behavioralCompetencyId: editItem.behavioralCompetencyId } : {}),
        competencyId: Number(formData.competencyId),
        jobGradeGroupId: Number(formData.jobGradeGroupId),
        ratingId: formData.ratingId ? Number(formData.ratingId) : undefined,
      };
      const res = await saveBehavioralCompetency(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Updated." : "Created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<BehavioralCompetency>[] = [
    { accessorKey: "competencyName", header: "Competency" },
    { accessorKey: "jobGradeGroupName", header: "Grade Group" },
    { accessorKey: "ratingName", header: "Expected Rating", cell: ({ row }) => row.original.ratingName ?? "—" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Behavioral Competencies" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Behavioral" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Behavioral Competencies" description="Map behavioral competencies to grade groups with expected ratings"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Behavioral" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Mapping</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="competencyName" searchPlaceholder="Search competencies..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Behavioral Competency" : "Add Behavioral Competency"} isEdit={!!editItem}>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Competency *</Label>
            <Select value={formData.competencyId} onValueChange={(v) => setFormData({ ...formData, competencyId: v })}>
              <SelectTrigger><SelectValue placeholder="Select competency" /></SelectTrigger>
              <SelectContent>{competencies.filter((c) => c.isApproved).map((c) => <SelectItem key={c.competencyId} value={String(c.competencyId)}>{c.competencyName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Grade Group *</Label>
            <Select value={formData.jobGradeGroupId} onValueChange={(v) => setFormData({ ...formData, jobGradeGroupId: v })}>
              <SelectTrigger><SelectValue placeholder="Select grade group" /></SelectTrigger>
              <SelectContent>{gradeGroups.map((g) => <SelectItem key={g.jobGradeGroupId} value={String(g.jobGradeGroupId)}>{g.groupName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Expected Rating</Label>
            <Select value={formData.ratingId} onValueChange={(v) => setFormData({ ...formData, ratingId: v })}>
              <SelectTrigger><SelectValue placeholder="Select rating" /></SelectTrigger>
              <SelectContent>{ratings.map((r) => <SelectItem key={r.ratingId} value={String(r.ratingId)}>{r.name} ({r.value})</SelectItem>)}</SelectContent>
            </Select>
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
