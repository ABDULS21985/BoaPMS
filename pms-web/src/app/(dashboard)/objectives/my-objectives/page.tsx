"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { Plus, Send, Loader2, Target } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { PageHeader } from "@/components/shared/page-header";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { EmptyState } from "@/components/shared/empty-state";
import { getStaffIndividualObjectives, saveDraftIndividualObjective, submitDraftIndividualObjective, cancelIndividualObjective } from "@/lib/api/pms-engine";
import { getStaffActiveReviewPeriod, getReviewPeriodCategoryDefinitions } from "@/lib/api/review-periods";
import { getConsolidatedObjectives } from "@/lib/api/performance";
import type { IndividualPlannedObjective, PerformanceReviewPeriod, CategoryDefinition, ConsolidatedObjective } from "@/types/performance";
import { Status } from "@/types/enums";

export default function MyObjectivesPage() {
  const { data: session } = useSession();
  const staffId = session?.user?.id ?? "";
  const [objectives, setObjectives] = useState<IndividualPlannedObjective[]>([]);
  const [period, setPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [categories, setCategories] = useState<CategoryDefinition[]>([]);
  const [availableObjs, setAvailableObjs] = useState<ConsolidatedObjective[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [submitOpen, setSubmitOpen] = useState(false);
  const [cancelOpen, setCancelOpen] = useState(false);
  const [selectedObj, setSelectedObj] = useState<IndividualPlannedObjective | null>(null);
  const [formData, setFormData] = useState({ objectiveId: "", weight: "", targetDate: "", keyPerformanceIndicator: "", description: "" });

  const loadData = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const periodRes = await getStaffActiveReviewPeriod();
      if (periodRes?.data) {
        setPeriod(periodRes.data);
        const [objRes, catRes, consolRes] = await Promise.all([
          getStaffIndividualObjectives(staffId, periodRes.data.periodId),
          getReviewPeriodCategoryDefinitions(periodRes.data.periodId),
          getConsolidatedObjectives(),
        ]);
        if (objRes?.data) setObjectives(Array.isArray(objRes.data) ? objRes.data : []);
        if (catRes?.data) setCategories(Array.isArray(catRes.data) ? catRes.data : []);
        if (consolRes?.data?.objectives) setAvailableObjs(consolRes.data.objectives);
      }
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [staffId]);

  const openAdd = () => {
    setFormData({ objectiveId: "", weight: "", targetDate: "", keyPerformanceIndicator: "", description: "" });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.objectiveId || !formData.weight || !period) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const res = await saveDraftIndividualObjective({ ...formData, reviewPeriodId: period.periodId, staffId, weight: Number(formData.weight) });
      if (res?.isSuccess) { toast.success("Objective saved as draft."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const handleSubmitAll = async () => {
    if (!period) return;
    try {
      const draftObjs = objectives.filter((o) => o.recordStatus === Status.Draft);
      for (const obj of draftObjs) {
        await submitDraftIndividualObjective({ id: obj.individualPlannedObjectiveId, reviewPeriodId: period.periodId, staffId });
      }
      toast.success("Objectives submitted for approval.");
      loadData();
    } catch { toast.error("Submit failed."); }
  };

  const handleCancel = async () => {
    if (!selectedObj || !period) return;
    try {
      const res = await cancelIndividualObjective({ id: selectedObj.individualPlannedObjectiveId, reviewPeriodId: period.periodId, staffId });
      if (res?.isSuccess) { toast.success("Objective cancelled."); loadData(); }
      else toast.error(res?.message || "Cancel failed.");
    } catch { toast.error("An error occurred."); }
  };

  const grouped = categories.map((cat) => ({
    category: cat,
    items: objectives.filter((o) => o.categoryName === cat.categoryName),
  }));

  const totalWeight = objectives.reduce((sum, o) => sum + (o.weight || 0), 0);
  const draftCount = objectives.filter((o) => o.recordStatus === Status.Draft).length;

  if (loading) return <div><PageHeader title="My Objectives" breadcrumbs={[{ label: "Objectives" }, { label: "My Objectives" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="My Objectives" description={period ? `Review Period: ${period.name}` : "No active review period"} breadcrumbs={[{ label: "Objectives", href: "/objectives/my-objectives" }, { label: "My Objectives" }]} actions={
        <div className="flex gap-2">
          {draftCount > 0 && <Button size="sm" variant="outline" onClick={() => setSubmitOpen(true)}><Send className="mr-2 h-4 w-4" />Submit All Drafts</Button>}
          {period?.allowObjectivePlanning && <Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Objective</Button>}
        </div>
      } />

      {/* Summary Cards */}
      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        <Card><CardContent className="pt-4"><div className="text-2xl font-bold">{objectives.length}</div><p className="text-xs text-muted-foreground">Total Objectives</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-2xl font-bold">{totalWeight}%</div><p className="text-xs text-muted-foreground">Total Weight</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-2xl font-bold">{objectives.filter((o) => o.recordStatus === Status.ApprovedAndActive).length}</div><p className="text-xs text-muted-foreground">Approved</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-2xl font-bold">{draftCount}</div><p className="text-xs text-muted-foreground">Drafts</p></CardContent></Card>
      </div>

      {/* Grouped by Category */}
      {grouped.length > 0 ? grouped.map(({ category, items }) => (
        <Card key={category.definitionId}>
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-base">{category.categoryName ?? "Uncategorized"}</CardTitle>
              <Badge variant="outline">Weight: {category.weight}%</Badge>
            </div>
          </CardHeader>
          <Separator />
          <CardContent className="pt-4">
            {items.length > 0 ? (
              <div className="space-y-3">
                {items.map((obj) => (
                  <div key={obj.individualPlannedObjectiveId} className="flex items-center justify-between rounded-lg border p-3">
                    <div className="space-y-1">
                      <p className="font-medium">{obj.objectiveName || obj.title}</p>
                      {obj.keyPerformanceIndicator && <p className="text-sm text-muted-foreground">KPI: {obj.keyPerformanceIndicator}</p>}
                    </div>
                    <div className="flex items-center gap-3">
                      <Badge variant="secondary">{obj.weight}%</Badge>
                      {obj.recordStatus != null && <StatusBadge status={obj.recordStatus} />}
                      {obj.recordStatus === Status.Draft && (
                        <Button size="sm" variant="ghost" className="text-destructive" onClick={() => { setSelectedObj(obj); setCancelOpen(true); }}>Cancel</Button>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">No objectives in this category yet.</p>
            )}
          </CardContent>
        </Card>
      )) : (
        <EmptyState icon={Target} title="No Objectives" description={period?.allowObjectivePlanning ? "Start by adding your planned objectives for this review period." : "Objective planning is not currently enabled for this period."} />
      )}

      {/* Add Objective Sheet */}
      <FormSheet open={open} onOpenChange={setOpen} title="Add Planned Objective" className="sm:max-w-lg overflow-y-auto">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Objective *</Label>
            <Select value={formData.objectiveId} onValueChange={(v) => setFormData({ ...formData, objectiveId: v })}>
              <SelectTrigger><SelectValue placeholder="Select objective" /></SelectTrigger>
              <SelectContent>{availableObjs.map((o) => <SelectItem key={o.objectiveId} value={o.objectiveId}>{o.name} ({o.objectiveLevel})</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Weight (%) *</Label><Input type="number" value={formData.weight} onChange={(e) => setFormData({ ...formData, weight: e.target.value })} /></div>
          <div className="space-y-2"><Label>Target Date</Label><Input type="date" value={formData.targetDate} onChange={(e) => setFormData({ ...formData, targetDate: e.target.value })} /></div>
          <div className="space-y-2"><Label>KPI</Label><Input value={formData.keyPerformanceIndicator} onChange={(e) => setFormData({ ...formData, keyPerformanceIndicator: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save Draft</Button>
          </div>
        </div>
      </FormSheet>

      <ConfirmationDialog open={submitOpen} onOpenChange={setSubmitOpen} title="Submit Objectives" description={`Submit ${draftCount} draft objective(s) for approval?`} confirmLabel="Submit" onConfirm={handleSubmitAll} />
      <ConfirmationDialog open={cancelOpen} onOpenChange={setCancelOpen} title="Cancel Objective" description="Are you sure you want to cancel this objective?" variant="destructive" confirmLabel="Cancel Objective" onConfirm={handleCancel} />
    </div>
  );
}
