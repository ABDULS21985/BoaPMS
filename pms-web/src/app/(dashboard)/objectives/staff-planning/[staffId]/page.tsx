"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { CheckCircle, XCircle, RotateCcw, Target, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { EmptyState } from "@/components/shared/empty-state";
import { getStaffIndividualObjectives, approveIndividualObjective, rejectIndividualObjective, returnIndividualObjective } from "@/lib/api/pms-engine";
import { getStaffActiveReviewPeriod, getReviewPeriodCategoryDefinitions } from "@/lib/api/review-periods";
import { getStaffDetails } from "@/lib/api/staff";
import type { IndividualPlannedObjective, PerformanceReviewPeriod, CategoryDefinition } from "@/types/performance";
import type { Staff } from "@/types/staff";
import { Status } from "@/types/enums";

type ActionType = "approve" | "reject" | "return";

export default function StaffPlanningPage() {
  const { staffId } = useParams<{ staffId: string }>();
  const [objectives, setObjectives] = useState<IndividualPlannedObjective[]>([]);
  const [period, setPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [categories, setCategories] = useState<CategoryDefinition[]>([]);
  const [staff, setStaff] = useState<Staff | null>(null);
  const [loading, setLoading] = useState(true);
  const [actionType, setActionType] = useState<ActionType | null>(null);
  const [selectedObj, setSelectedObj] = useState<IndividualPlannedObjective | null>(null);
  const [processing, setProcessing] = useState(false);

  const loadData = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const [staffRes, periodRes] = await Promise.all([getStaffDetails(staffId), getStaffActiveReviewPeriod()]);
      if (staffRes?.data) setStaff(staffRes.data);
      if (periodRes?.data) {
        setPeriod(periodRes.data);
        const [objRes, catRes] = await Promise.all([
          getStaffIndividualObjectives(staffId, periodRes.data.periodId),
          getReviewPeriodCategoryDefinitions(periodRes.data.periodId),
        ]);
        if (objRes?.data) setObjectives(Array.isArray(objRes.data) ? objRes.data : []);
        if (catRes?.data) setCategories(Array.isArray(catRes.data) ? catRes.data : []);
      }
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [staffId]);

  const handleAction = async (reason?: string) => {
    if (!selectedObj || !period || !actionType) return;
    setProcessing(true);
    try {
      const payload = { id: selectedObj.individualPlannedObjectiveId, reviewPeriodId: period.periodId, staffId, remark: reason };
      const fn = actionType === "approve" ? approveIndividualObjective : actionType === "reject" ? rejectIndividualObjective : returnIndividualObjective;
      const res = await fn(payload);
      if (res?.isSuccess) { toast.success(`Objective ${actionType}d.`); loadData(); }
      else toast.error(res?.message || "Action failed.");
    } catch { toast.error("An error occurred."); } finally { setProcessing(false); }
  };

  const grouped = categories.map((cat) => ({
    category: cat,
    items: objectives.filter((o) => o.categoryName === cat.categoryName),
  }));

  const totalWeight = objectives.reduce((sum, o) => sum + (o.weight || 0), 0);
  const staffName = staff ? `${staff.firstName} ${staff.lastName}` : staffId;

  if (loading) return <div><PageHeader title="Staff Objectives" breadcrumbs={[{ label: "Objectives" }, { label: "Staff Planning" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={`${staffName}'s Objectives`} description={period ? `Review Period: ${period.name} | Total Weight: ${totalWeight}%` : ""} breadcrumbs={[{ label: "Objectives", href: "/objectives/my-objectives" }, { label: "Direct Reports", href: "/objectives/direct-reports" }, { label: staffName }]} />

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
                      {obj.description && <p className="text-sm text-muted-foreground">{obj.description}</p>}
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant="secondary">{obj.weight}%</Badge>
                      {obj.recordStatus != null && <StatusBadge status={obj.recordStatus} />}
                      {obj.recordStatus === Status.PendingApproval && (
                        <div className="flex gap-1">
                          <Button size="sm" variant="default" onClick={() => { setSelectedObj(obj); setActionType("approve"); }}><CheckCircle className="h-3.5 w-3.5" /></Button>
                          <Button size="sm" variant="outline" onClick={() => { setSelectedObj(obj); setActionType("return"); }}><RotateCcw className="h-3.5 w-3.5" /></Button>
                          <Button size="sm" variant="destructive" onClick={() => { setSelectedObj(obj); setActionType("reject"); }}><XCircle className="h-3.5 w-3.5" /></Button>
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">No objectives in this category.</p>
            )}
          </CardContent>
        </Card>
      )) : (
        <EmptyState icon={Target} title="No Objectives" description="This staff member has no planned objectives yet." />
      )}

      <ConfirmationDialog open={!!actionType} onOpenChange={(o) => { if (!o) setActionType(null); }} title={`${actionType === "approve" ? "Approve" : actionType === "reject" ? "Reject" : "Return"} Objective`} description={`Are you sure you want to ${actionType} "${selectedObj?.objectiveName || selectedObj?.title}"?`} variant={actionType === "reject" ? "destructive" : "default"} confirmLabel={actionType === "approve" ? "Approve" : actionType === "reject" ? "Reject" : "Return"} showReasonInput={actionType === "reject" || actionType === "return"} reasonLabel="Remark" onConfirm={handleAction} />
    </div>
  );
}
