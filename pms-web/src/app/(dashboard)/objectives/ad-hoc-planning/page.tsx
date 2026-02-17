"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { Plus, Loader2, Zap } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getStaffIndividualObjectives, saveDraftIndividualObjective } from "@/lib/api/pms-engine";
import { getStaffActiveReviewPeriod } from "@/lib/api/review-periods";
import { getConsolidatedObjectives } from "@/lib/api/performance";
import type { IndividualPlannedObjective, PerformanceReviewPeriod, ConsolidatedObjective } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function AdHocPlanningPage() {
  const { data: session } = useSession();
  const staffId = session?.user?.id ?? "";
  const [objectives, setObjectives] = useState<IndividualPlannedObjective[]>([]);
  const [period, setPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [availableObjs, setAvailableObjs] = useState<ConsolidatedObjective[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({ objectiveId: "", weight: "", targetDate: "", keyPerformanceIndicator: "", description: "" });

  const loadData = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const periodRes = await getStaffActiveReviewPeriod();
      if (periodRes?.data) {
        setPeriod(periodRes.data);
        const [objRes, consolRes] = await Promise.all([
          getStaffIndividualObjectives(staffId, periodRes.data.periodId),
          getConsolidatedObjectives(),
        ]);
        if (objRes?.data) setObjectives(Array.isArray(objRes.data) ? objRes.data : []);
        if (consolRes?.data?.objectives) setAvailableObjs(consolRes.data.objectives.filter((o) => o.objectiveLevel === "Office"));
      }
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [staffId]);

  const handleSave = async () => {
    if (!formData.objectiveId || !period) { toast.error("Please select an objective."); return; }
    setSaving(true);
    try {
      const res = await saveDraftIndividualObjective({ ...formData, reviewPeriodId: period.periodId, staffId, weight: Number(formData.weight) || 0 });
      if (res?.isSuccess) { toast.success("Ad-hoc objective saved."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<IndividualPlannedObjective>[] = [
    { accessorKey: "objectiveName", header: "Objective", cell: ({ row }) => <span className="font-medium">{row.original.objectiveName || row.original.title}</span> },
    { accessorKey: "objectiveLevel", header: "Level" },
    { accessorKey: "weight", header: "Weight", cell: ({ row }) => `${row.original.weight}%` },
    { accessorKey: "keyPerformanceIndicator", header: "KPI" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
  ];

  if (loading) return <div><PageHeader title="Ad-Hoc Planning" breadcrumbs={[{ label: "Objectives" }, { label: "Ad-Hoc Planning" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Ad-Hoc Planning" description="Plan ad-hoc objectives outside standard categories" breadcrumbs={[{ label: "Objectives", href: "/objectives/my-objectives" }, { label: "Ad-Hoc Planning" }]} actions={period?.allowObjectivePlanning ? <Button size="sm" onClick={() => { setFormData({ objectiveId: "", weight: "", targetDate: "", keyPerformanceIndicator: "", description: "" }); setOpen(true); }}><Plus className="mr-2 h-4 w-4" />Add Ad-Hoc Objective</Button> : undefined} />

      {objectives.length > 0 ? (
        <DataTable columns={columns} data={objectives} searchKey="objectiveName" searchPlaceholder="Search objectives..." />
      ) : (
        <EmptyState icon={Zap} title="No Ad-Hoc Objectives" description="No ad-hoc objectives have been planned for this period." />
      )}

      <FormSheet open={open} onOpenChange={setOpen} title="Add Ad-Hoc Objective" className="sm:max-w-lg overflow-y-auto">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Objective *</Label>
            <Select value={formData.objectiveId} onValueChange={(v) => setFormData({ ...formData, objectiveId: v })}>
              <SelectTrigger><SelectValue placeholder="Select objective" /></SelectTrigger>
              <SelectContent>{availableObjs.map((o) => <SelectItem key={o.objectiveId} value={o.objectiveId}>{o.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Weight (%)</Label><Input type="number" value={formData.weight} onChange={(e) => setFormData({ ...formData, weight: e.target.value })} /></div>
          <div className="space-y-2"><Label>Target Date</Label><Input type="date" value={formData.targetDate} onChange={(e) => setFormData({ ...formData, targetDate: e.target.value })} /></div>
          <div className="space-y-2"><Label>KPI</Label><Input value={formData.keyPerformanceIndicator} onChange={(e) => setFormData({ ...formData, keyPerformanceIndicator: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
