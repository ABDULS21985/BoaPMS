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
import { getJobRoleCompetencies, saveJobRoleCompetency, getJobRoles, getCompetencies, getRatings } from "@/lib/api/competency";
import type { JobRoleCompetency, JobRole, Competency, Rating } from "@/types/competency";

export default function JobRoleCompetenciesPage() {
  const [items, setItems] = useState<JobRoleCompetency[]>([]);
  const [jobRoles, setJobRoles] = useState<JobRole[]>([]);
  const [competencies, setCompetencies] = useState<Competency[]>([]);
  const [ratings, setRatings] = useState<Rating[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<JobRoleCompetency | null>(null);
  const [formData, setFormData] = useState({ jobRoleId: "", competencyId: "", ratingId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, jrRes, compRes, ratRes] = await Promise.all([
        getJobRoleCompetencies(), getJobRoles(), getCompetencies(), getRatings(),
      ]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (jrRes?.data) setJobRoles(Array.isArray(jrRes.data) ? jrRes.data : []);
      if (compRes?.data) setCompetencies(Array.isArray(compRes.data) ? compRes.data : []);
      if (ratRes?.data) setRatings(Array.isArray(ratRes.data) ? ratRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ jobRoleId: "", competencyId: "", ratingId: "" }); setOpen(true); };
  const openEdit = (item: JobRoleCompetency) => {
    setEditItem(item);
    setFormData({ jobRoleId: String(item.jobRoleId), competencyId: String(item.competencyId), ratingId: String(item.ratingId) });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.jobRoleId || !formData.competencyId || !formData.ratingId) { toast.error("All fields are required."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { jobRoleCompetencyId: editItem.jobRoleCompetencyId } : {}),
        jobRoleId: Number(formData.jobRoleId),
        competencyId: Number(formData.competencyId),
        ratingId: Number(formData.ratingId),
      };
      const res = await saveJobRoleCompetency(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Updated." : "Created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<JobRoleCompetency>[] = [
    { accessorKey: "jobRoleName", header: "Job Role" },
    { accessorKey: "competencyName", header: "Competency" },
    { accessorKey: "ratingName", header: "Expected Rating" },
    { accessorKey: "officeName", header: "Office", cell: ({ row }) => row.original.officeName ?? "—" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Job Role Competencies" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Job Role Competencies" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Job Role Competencies" description="Map competencies to job roles with expected ratings"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Job Role Competencies" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Mapping</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="jobRoleName" searchPlaceholder="Search by job role..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Mapping" : "Add Job Role Competency"} isEdit={!!editItem}>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Job Role *</Label>
            <Select value={formData.jobRoleId} onValueChange={(v) => setFormData({ ...formData, jobRoleId: v })}>
              <SelectTrigger><SelectValue placeholder="Select job role" /></SelectTrigger>
              <SelectContent>{jobRoles.map((r) => <SelectItem key={r.jobRoleId} value={String(r.jobRoleId)}>{r.jobRoleName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Competency *</Label>
            <Select value={formData.competencyId} onValueChange={(v) => setFormData({ ...formData, competencyId: v })}>
              <SelectTrigger><SelectValue placeholder="Select competency" /></SelectTrigger>
              <SelectContent>{competencies.filter((c) => c.isApproved).map((c) => <SelectItem key={c.competencyId} value={String(c.competencyId)}>{c.competencyName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Expected Rating *</Label>
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
