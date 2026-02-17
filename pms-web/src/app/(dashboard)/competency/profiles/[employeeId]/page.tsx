"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Loader2, User } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getCompetencyReviewProfiles, getDevelopmentPlans, saveDevelopmentPlan, getTrainingTypes } from "@/lib/api/competency";
import type { CompetencyReviewProfile, DevelopmentPlan, TrainingType } from "@/types/competency";

export default function EmployeeProfilePage() {
  const params = useParams();
  const employeeId = params.employeeId as string;
  const [profiles, setProfiles] = useState<CompetencyReviewProfile[]>([]);
  const [devPlans, setDevPlans] = useState<DevelopmentPlan[]>([]);
  const [trainingTypes, setTrainingTypes] = useState<TrainingType[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [selectedProfile, setSelectedProfile] = useState<CompetencyReviewProfile | null>(null);
  const [formData, setFormData] = useState({ activity: "", trainingTypeId: "", targetDate: "", learningResource: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [profRes, ttRes] = await Promise.all([getCompetencyReviewProfiles(employeeId), getTrainingTypes()]);
      if (profRes?.data) {
        const profs = Array.isArray(profRes.data) ? profRes.data : [];
        setProfiles(profs);
        if (profs.length > 0) setSelectedProfile(profs[0]);
      }
      if (ttRes?.data) setTrainingTypes(Array.isArray(ttRes.data) ? ttRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  const loadDevPlans = async (profileId: number) => {
    try {
      const res = await getDevelopmentPlans(profileId);
      if (res?.data) setDevPlans(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ }
  };

  useEffect(() => { loadData(); }, [employeeId]);
  useEffect(() => { if (selectedProfile) loadDevPlans(selectedProfile.competencyReviewProfileId); }, [selectedProfile]);

  const handleSavePlan = async () => {
    if (!formData.activity || !selectedProfile) { toast.error("Activity is required."); return; }
    setSaving(true);
    try {
      const res = await saveDevelopmentPlan({
        competencyReviewProfileId: selectedProfile.competencyReviewProfileId,
        activity: formData.activity,
        trainingTypeId: formData.trainingTypeId ? Number(formData.trainingTypeId) : undefined,
        targetDate: formData.targetDate || undefined,
        learningResource: formData.learningResource || undefined,
        employeeNumber: employeeId,
      });
      if (res?.isSuccess) { toast.success("Development plan added."); setOpen(false); loadDevPlans(selectedProfile.competencyReviewProfileId); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const employeeName = profiles[0]?.employeeFullName ?? employeeId;
  const gap = selectedProfile ? (selectedProfile.expectedRatingValue - selectedProfile.averageRatingValue).toFixed(1) : "—";

  const profileColumns: ColumnDef<CompetencyReviewProfile>[] = [
    { accessorKey: "competencyName", header: "Competency" },
    { accessorKey: "competencyCategoryName", header: "Category" },
    { accessorKey: "reviewPeriodName", header: "Period" },
    { accessorKey: "expectedRatingValue", header: "Expected", cell: ({ row }) => <Badge variant="outline">{row.original.expectedRatingName ?? row.original.expectedRatingValue}</Badge> },
    { accessorKey: "averageRatingValue", header: "Actual", cell: ({ row }) => <Badge variant={row.original.averageRatingValue >= row.original.expectedRatingValue ? "default" : "destructive"}>{row.original.averageRatingName ?? row.original.averageRatingValue}</Badge> },
    { accessorKey: "averageScore", header: "Score", cell: ({ row }) => row.original.averageScore.toFixed(1) },
  ];

  const planColumns: ColumnDef<DevelopmentPlan>[] = [
    { accessorKey: "activity", header: "Activity" },
    { accessorKey: "trainingTypeName", header: "Training Type", cell: ({ row }) => row.original.trainingTypeName ?? "—" },
    { accessorKey: "targetDate", header: "Target Date", cell: ({ row }) => row.original.targetDate?.split("T")[0] ?? "—" },
    { accessorKey: "taskStatus", header: "Status", cell: ({ row }) => <Badge variant="outline">{row.original.taskStatus ?? "Pending"}</Badge> },
    { accessorKey: "learningResource", header: "Resource", cell: ({ row }) => <span className="line-clamp-1">{row.original.learningResource}</span> },
  ];

  if (loading) return <div><PageHeader title="Employee Profile" breadcrumbs={[{ label: "Competency" }, { label: "Profiles" }, { label: "Details" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={employeeName} description={`Competency profile for ${employeeName}`}
        breadcrumbs={[{ label: "Competency" }, { label: "Profiles", href: "/competency/profiles" }, { label: employeeName }]}
      />

      {selectedProfile && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Competency</CardTitle></CardHeader><CardContent><p className="text-lg font-semibold">{selectedProfile.competencyName}</p></CardContent></Card>
          <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Expected Rating</CardTitle></CardHeader><CardContent><p className="text-lg font-semibold">{selectedProfile.expectedRatingName ?? selectedProfile.expectedRatingValue}</p></CardContent></Card>
          <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Actual Rating</CardTitle></CardHeader><CardContent><p className="text-lg font-semibold">{selectedProfile.averageRatingName ?? selectedProfile.averageRatingValue}</p></CardContent></Card>
          <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Gap</CardTitle></CardHeader><CardContent><p className={`text-lg font-semibold ${Number(gap) > 0 ? "text-red-600" : "text-green-600"}`}>{gap}</p></CardContent></Card>
        </div>
      )}

      <Tabs defaultValue="profiles">
        <TabsList>
          <TabsTrigger value="profiles">Competency Profiles</TabsTrigger>
          <TabsTrigger value="plans">Development Plans {devPlans.length > 0 && <Badge variant="secondary" className="ml-2">{devPlans.length}</Badge>}</TabsTrigger>
        </TabsList>

        <TabsContent value="profiles" className="mt-4">
          {profiles.length > 0 ? (
            <DataTable columns={profileColumns} data={profiles} searchKey="competencyName" searchPlaceholder="Search competencies..." />
          ) : (
            <EmptyState icon={User} title="No Profiles" description="No competency profiles found." />
          )}
        </TabsContent>

        <TabsContent value="plans" className="mt-4">
          <div className="flex justify-end mb-4">
            <Button size="sm" onClick={() => { setFormData({ activity: "", trainingTypeId: "", targetDate: "", learningResource: "" }); setOpen(true); }}>
              <Plus className="mr-2 h-4 w-4" />Add Development Plan
            </Button>
          </div>
          {devPlans.length > 0 ? (
            <DataTable columns={planColumns} data={devPlans} searchKey="activity" searchPlaceholder="Search plans..." />
          ) : (
            <EmptyState icon={User} title="No Development Plans" description="No development plans for this profile." />
          )}
        </TabsContent>
      </Tabs>

      <FormSheet open={open} onOpenChange={setOpen} title="Add Development Plan">
        <div className="space-y-4">
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
            <Button className="flex-1" onClick={handleSavePlan} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
