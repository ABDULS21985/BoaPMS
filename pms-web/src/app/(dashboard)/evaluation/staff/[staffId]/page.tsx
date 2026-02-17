"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Loader2, Star } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getStaffIndividualObjectives, saveDraftEvaluation, submitDraftEvaluation } from "@/lib/api/pms-engine";
import { getStaffActiveReviewPeriod } from "@/lib/api/review-periods";
import { getStaffDetails } from "@/lib/api/staff";
import { getEvaluationOptions } from "@/lib/api/performance";
import type { IndividualPlannedObjective, PerformanceReviewPeriod, EvaluationOption } from "@/types/performance";
import type { Staff } from "@/types/staff";
import { Status } from "@/types/enums";

interface EvalForm {
  objectiveId: string;
  score: string;
  comment: string;
}

export default function StaffEvaluationPage() {
  const { staffId } = useParams<{ staffId: string }>();
  const [objectives, setObjectives] = useState<IndividualPlannedObjective[]>([]);
  const [period, setPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [staff, setStaff] = useState<Staff | null>(null);
  const [evalOptions, setEvalOptions] = useState<EvaluationOption[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [forms, setForms] = useState<Record<string, EvalForm>>({});

  const loadData = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const [staffRes, periodRes, optRes] = await Promise.all([getStaffDetails(staffId), getStaffActiveReviewPeriod(), getEvaluationOptions()]);
      if (staffRes?.data) setStaff(staffRes.data);
      if (optRes?.data) setEvalOptions(Array.isArray(optRes.data) ? optRes.data : []);
      if (periodRes?.data) {
        setPeriod(periodRes.data);
        const objRes = await getStaffIndividualObjectives(staffId, periodRes.data.periodId);
        if (objRes?.data) {
          const objs = Array.isArray(objRes.data) ? objRes.data : [];
          setObjectives(objs.filter((o) => o.recordStatus === Status.ApprovedAndActive || o.recordStatus === Status.Active));
          const initial: Record<string, EvalForm> = {};
          objs.forEach((o) => { initial[o.individualPlannedObjectiveId] = { objectiveId: o.individualPlannedObjectiveId, score: "", comment: "" }; });
          setForms(initial);
        }
      }
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [staffId]);

  const updateForm = (id: string, field: keyof EvalForm, value: string) => {
    setForms((prev) => ({ ...prev, [id]: { ...prev[id], [field]: value } }));
  };

  const handleSaveAll = async () => {
    if (!period) return;
    setSaving(true);
    try {
      for (const obj of objectives) {
        const form = forms[obj.individualPlannedObjectiveId];
        if (form?.score) {
          await saveDraftEvaluation({ reviewPeriodId: period.periodId, staffId, objectiveId: obj.individualPlannedObjectiveId, score: Number(form.score), comment: form.comment });
        }
      }
      toast.success("Evaluations saved as draft.");
      loadData();
    } catch { toast.error("Save failed."); } finally { setSaving(false); }
  };

  const handleSubmitAll = async () => {
    if (!period) return;
    setSaving(true);
    try {
      for (const obj of objectives) {
        const form = forms[obj.individualPlannedObjectiveId];
        if (form?.score) {
          await submitDraftEvaluation({ reviewPeriodId: period.periodId, staffId, objectiveId: obj.individualPlannedObjectiveId, score: Number(form.score), comment: form.comment });
        }
      }
      toast.success("Evaluations submitted.");
      loadData();
    } catch { toast.error("Submit failed."); } finally { setSaving(false); }
  };

  const staffName = staff ? `${staff.firstName} ${staff.lastName}` : staffId;
  const totalWeight = objectives.reduce((sum, o) => sum + (o.weight || 0), 0);

  if (loading) return <div><PageHeader title="Staff Evaluation" breadcrumbs={[{ label: "Evaluation" }, { label: "Staff" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={`Evaluate: ${staffName}`} description={period ? `${period.name} | Total Weight: ${totalWeight}%` : ""} breadcrumbs={[{ label: "Evaluation", href: "/evaluation/direct-reports" }, { label: staffName }]} actions={
        <div className="flex gap-2">
          <Button size="sm" variant="outline" onClick={handleSaveAll} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save Draft</Button>
          <Button size="sm" onClick={handleSubmitAll} disabled={saving}>Submit Evaluations</Button>
        </div>
      } />

      {objectives.length > 0 ? (
        <div className="space-y-4">
          {objectives.map((obj) => (
            <Card key={obj.individualPlannedObjectiveId}>
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <CardTitle className="text-base">{obj.objectiveName || obj.title}</CardTitle>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary">{obj.weight}%</Badge>
                    {obj.recordStatus != null && <StatusBadge status={obj.recordStatus} />}
                  </div>
                </div>
                {obj.keyPerformanceIndicator && <p className="text-sm text-muted-foreground">KPI: {obj.keyPerformanceIndicator}</p>}
              </CardHeader>
              <Separator />
              <CardContent className="pt-4 space-y-3">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label>Score *</Label>
                    <Input type="number" min="0" max={period?.maxPoints ?? 100} placeholder="Enter score" value={forms[obj.individualPlannedObjectiveId]?.score ?? ""} onChange={(e) => updateForm(obj.individualPlannedObjectiveId, "score", e.target.value)} />
                  </div>
                  <div className="space-y-2">
                    <Label>Comment</Label>
                    <Input placeholder="Evaluation comment" value={forms[obj.individualPlannedObjectiveId]?.comment ?? ""} onChange={(e) => updateForm(obj.individualPlannedObjectiveId, "comment", e.target.value)} />
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <EmptyState icon={Star} title="No Objectives to Evaluate" description="This staff member has no approved objectives for evaluation." />
      )}
    </div>
  );
}
