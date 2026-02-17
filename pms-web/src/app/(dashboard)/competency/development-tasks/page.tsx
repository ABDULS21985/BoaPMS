"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Loader2, ListTodo } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getDevelopmentPlans, saveDevelopmentPlan, getCompetencyReviewProfiles, getTrainingTypes } from "@/lib/api/competency";
import type { DevelopmentPlan, CompetencyReviewProfile, TrainingType } from "@/types/competency";

export default function DevelopmentTasksPage() {
  const { data: session } = useSession();
  const employeeNumber = session?.user?.id ?? "";
  const [plans, setPlans] = useState<DevelopmentPlan[]>([]);
  const [profiles, setProfiles] = useState<CompetencyReviewProfile[]>([]);
  const [trainingTypes, setTrainingTypes] = useState<TrainingType[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({ competencyReviewProfileId: "", activity: "", trainingTypeId: "", targetDate: "", learningResource: "" });

  const loadData = async () => {
    if (!employeeNumber) return;
    setLoading(true);
    try {
      const [planRes, profRes, ttRes] = await Promise.all([
        getDevelopmentPlans(),
        getCompetencyReviewProfiles(employeeNumber),
        getTrainingTypes(),
      ]);
      if (planRes?.data) setPlans(Array.isArray(planRes.data) ? planRes.data : []);
      if (profRes?.data) setProfiles(Array.isArray(profRes.data) ? profRes.data : []);
      if (ttRes?.data) setTrainingTypes(Array.isArray(ttRes.data) ? ttRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [employeeNumber]);

  const handleSave = async () => {
    if (!formData.activity || !formData.competencyReviewProfileId) { toast.error("Profile and activity are required."); return; }
    setSaving(true);
    try {
      const res = await saveDevelopmentPlan({
        competencyReviewProfileId: Number(formData.competencyReviewProfileId),
        activity: formData.activity,
        trainingTypeId: formData.trainingTypeId ? Number(formData.trainingTypeId) : undefined,
        targetDate: formData.targetDate || undefined,
        learningResource: formData.learningResource || undefined,
        employeeNumber,
      });
      if (res?.isSuccess) { toast.success("Development plan created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const statusColor = (status?: string) => {
    if (!status) return "secondary";
    const s = status.toLowerCase();
    if (s === "completed") return "default" as const;
    if (s === "in progress") return "secondary" as const;
    return "outline" as const;
  };

  const columns: ColumnDef<DevelopmentPlan>[] = [
    { accessorKey: "activity", header: "Activity" },
    { accessorKey: "competencyName", header: "Competency" },
    { accessorKey: "competencyCategoryName", header: "Category" },
    { accessorKey: "trainingTypeName", header: "Training Type", cell: ({ row }) => row.original.trainingTypeName ?? "—" },
    { accessorKey: "targetDate", header: "Target Date", cell: ({ row }) => row.original.targetDate?.split("T")[0] ?? "—" },
    { accessorKey: "completionDate", header: "Completed", cell: ({ row }) => row.original.completionDate?.split("T")[0] ?? "—" },
    { accessorKey: "taskStatus", header: "Status", cell: ({ row }) => <Badge variant={statusColor(row.original.taskStatus)}>{row.original.taskStatus ?? "Pending"}</Badge> },
    { accessorKey: "currentGap", header: "Gap", cell: ({ row }) => row.original.currentGap != null ? <Badge variant={row.original.currentGap > 0 ? "destructive" : "default"}>{row.original.currentGap}</Badge> : "—" },
  ];

  if (loading) return <div><PageHeader title="Development Tasks" breadcrumbs={[{ label: "Competency" }, { label: "Development Tasks" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Development Tasks" description="Manage your competency development plans"
        breadcrumbs={[{ label: "Competency" }, { label: "Development Tasks" }]}
        actions={<Button size="sm" onClick={() => { setFormData({ competencyReviewProfileId: "", activity: "", trainingTypeId: "", targetDate: "", learningResource: "" }); setOpen(true); }}><Plus className="mr-2 h-4 w-4" />Add Plan</Button>}
      />

      {plans.length > 0 ? (
        <DataTable columns={columns} data={plans} searchKey="activity" searchPlaceholder="Search development plans..." />
      ) : (
        <EmptyState icon={ListTodo} title="No Development Plans" description="You have no development plans yet." />
      )}

      <FormSheet open={open} onOpenChange={setOpen} title="Add Development Plan">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Competency Profile *</Label>
            <Select value={formData.competencyReviewProfileId} onValueChange={(v) => setFormData({ ...formData, competencyReviewProfileId: v })}>
              <SelectTrigger><SelectValue placeholder="Select profile" /></SelectTrigger>
              <SelectContent>{profiles.map((p) => <SelectItem key={p.competencyReviewProfileId} value={String(p.competencyReviewProfileId)}>{p.competencyName} ({p.reviewPeriodName})</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Activity *</Label><Input value={formData.activity} onChange={(e) => setFormData({ ...formData, activity: e.target.value })} /></div>
          <div className="space-y-2">
            <Label>Training Type</Label>
            <Select value={formData.trainingTypeId} onValueChange={(v) => setFormData({ ...formData, trainingTypeId: v })}>
              <SelectTrigger><SelectValue placeholder="Select training type" /></SelectTrigger>
              <SelectContent>{trainingTypes.map((t) => <SelectItem key={t.trainingTypeId} value={String(t.trainingTypeId)}>{t.trainingTypeName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Target Date</Label><Input type="date" value={formData.targetDate} onChange={(e) => setFormData({ ...formData, targetDate: e.target.value })} /></div>
          <div className="space-y-2"><Label>Learning Resource</Label><Input value={formData.learningResource} onChange={(e) => setFormData({ ...formData, learningResource: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
