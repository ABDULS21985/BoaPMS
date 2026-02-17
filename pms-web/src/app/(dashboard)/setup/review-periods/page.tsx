"use client";

import { useEffect, useState, useCallback } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2, Settings2, ChevronRight, ChevronLeft, Trash2, CheckCircle, RotateCcw, XCircle } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Separator } from "@/components/ui/separator";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import {
  getReviewPeriods, saveDraftReviewPeriod, updateReviewPeriod, submitDraftReviewPeriod,
  approveReviewPeriod, returnReviewPeriod, closeReviewPeriod, cancelReviewPeriod,
  enableObjectivePlanning, disableObjectivePlanning, enableWorkProductPlanning,
  disableWorkProductPlanning, enableWorkProductEvaluation, disableWorkProductEvaluation,
  getReviewPeriodObjectives, addReviewPeriodObjective,
  getReviewPeriodCategoryDefinitions, saveDraftCategoryDefinition,
} from "@/lib/api/review-periods";
import { getStrategies, getEnterpriseObjectives, getObjectiveCategories } from "@/lib/api/performance";
import { getJobGradeGroups } from "@/lib/api/competency";
import type { PerformanceReviewPeriod, CategoryDefinition, Strategy, EnterpriseObjective, ObjectiveCategory } from "@/types/performance";
import { ReviewPeriodRange, Status } from "@/types/enums";

const rangeLabels: Record<number, string> = { [ReviewPeriodRange.Quarterly]: "Quarterly", [ReviewPeriodRange.BiAnnual]: "Bi-Annual", [ReviewPeriodRange.Annual]: "Annual" };

interface GradeGroup { jobGradeGroupId: number; groupName: string; }

export default function ReviewPeriodsPage() {
  const [items, setItems] = useState<PerformanceReviewPeriod[]>([]);
  const [loading, setLoading] = useState(true);
  const [strategies, setStrategies] = useState<Strategy[]>([]);

  // Wizard state
  const [wizardOpen, setWizardOpen] = useState(false);
  const [step, setStep] = useState(1);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<PerformanceReviewPeriod | null>(null);
  const [formData, setFormData] = useState({ name: "", shortName: "", description: "", strategyId: "", maxPoints: "250", minNoOfObjectives: "1", maxNoOfObjectives: "10", range: "", rangeValue: "1", startDate: "", endDate: "" });

  // Step 2: Enterprise Objectives
  const [allObjectives, setAllObjectives] = useState<EnterpriseObjective[]>([]);
  const [periodObjectives, setPeriodObjectives] = useState<string[]>([]);
  const [selectedObjectives, setSelectedObjectives] = useState<Set<string>>(new Set());

  // Step 3: Category Definitions
  const [categories, setCategories] = useState<ObjectiveCategory[]>([]);
  const [gradeGroups, setGradeGroups] = useState<GradeGroup[]>([]);
  const [catDefs, setCatDefs] = useState<CategoryDefinition[]>([]);
  const [catForm, setCatForm] = useState({ objectiveCategoryId: "", gradeGroupId: "", weight: "", maxNoObjectives: "", maxNoWorkProduct: "", enforceWorkProductLimit: false });

  // Settings drawer
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [settingsItem, setSettingsItem] = useState<PerformanceReviewPeriod | null>(null);
  const [toggling, setToggling] = useState("");

  // Confirmation dialogs
  const [confirmAction, setConfirmAction] = useState<{ action: string; label: string; periodId: string } | null>(null);
  const [actionLoading, setActionLoading] = useState(false);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [res, stratRes] = await Promise.all([getReviewPeriods(), getStrategies()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (stratRes?.data) setStrategies(Array.isArray(stratRes.data) ? stratRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  }, []);

  useEffect(() => { loadData(); }, [loadData]);

  const openAdd = () => {
    setEditItem(null);
    setStep(1);
    setFormData({ name: "", shortName: "", description: "", strategyId: "", maxPoints: "250", minNoOfObjectives: "1", maxNoOfObjectives: "10", range: "", rangeValue: "1", startDate: "", endDate: "" });
    setSelectedObjectives(new Set());
    setPeriodObjectives([]);
    setCatDefs([]);
    setWizardOpen(true);
  };

  const openEdit = async (item: PerformanceReviewPeriod) => {
    setEditItem(item);
    setStep(1);
    setFormData({
      name: item.name, shortName: item.shortName ?? "", description: item.description ?? "",
      strategyId: item.strategyId ?? "", maxPoints: String(item.maxPoints), minNoOfObjectives: String(item.minNoOfObjectives),
      maxNoOfObjectives: String(item.maxNoOfObjectives), range: String(item.range), rangeValue: String(item.rangeValue),
      startDate: item.startDate?.split("T")[0] ?? "", endDate: item.endDate?.split("T")[0] ?? "",
    });
    setWizardOpen(true);
    // Load step 2 & 3 data
    try {
      const [objRes, catDefRes] = await Promise.all([
        getReviewPeriodObjectives(item.periodId),
        getReviewPeriodCategoryDefinitions(item.periodId),
      ]);
      const objIds = (objRes?.data as { enterpriseObjectiveId?: string }[] ?? []).map((o) => o.enterpriseObjectiveId).filter(Boolean) as string[];
      setPeriodObjectives(objIds);
      setSelectedObjectives(new Set(objIds));
      if (catDefRes?.data) setCatDefs(Array.isArray(catDefRes.data) ? catDefRes.data : []);
    } catch { /* */ }
  };

  const openSettings = (item: PerformanceReviewPeriod) => { setSettingsItem(item); setSettingsOpen(true); };

  // Load reference data for step 2 & 3
  const loadStepData = async () => {
    try {
      const [objRes, catRes, gradeRes] = await Promise.all([getEnterpriseObjectives(), getObjectiveCategories(), getJobGradeGroups()]);
      if (objRes?.data) setAllObjectives(Array.isArray(objRes.data) ? objRes.data : []);
      if (catRes?.data) setCategories(Array.isArray(catRes.data) ? catRes.data : []);
      if (gradeRes?.data) setGradeGroups(Array.isArray(gradeRes.data) ? gradeRes.data as GradeGroup[] : []);
    } catch { /* */ }
  };

  useEffect(() => { if (wizardOpen && step >= 2) loadStepData(); }, [wizardOpen, step]);

  // Step 1: Save draft
  const handleSaveStep1 = async () => {
    if (!formData.name || !formData.range || !formData.startDate || !formData.endDate) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const payload = {
        periodId: editItem?.periodId, name: formData.name, shortName: formData.shortName, description: formData.description,
        strategyId: formData.strategyId || undefined, maxPoints: Number(formData.maxPoints), minNoOfObjectives: Number(formData.minNoOfObjectives),
        maxNoOfObjectives: Number(formData.maxNoOfObjectives), range: Number(formData.range), rangeValue: Number(formData.rangeValue),
        startDate: formData.startDate, endDate: formData.endDate,
      };
      const res = editItem ? await updateReviewPeriod(payload) : await saveDraftReviewPeriod(payload);
      if (res?.isSuccess) {
        toast.success(editItem ? "Review period updated." : "Draft saved.");
        if (!editItem) {
          // Reload to get the new periodId
          const listRes = await getReviewPeriods();
          if (listRes?.data) {
            const newItems = Array.isArray(listRes.data) ? listRes.data : [];
            setItems(newItems);
            const newest = newItems.find((i) => i.name === formData.name);
            if (newest) setEditItem(newest);
          }
        }
        setStep(2);
      } else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  // Step 2: Add objective to period
  const toggleObjective = (id: string) => {
    setSelectedObjectives((prev) => { const n = new Set(prev); n.has(id) ? n.delete(id) : n.add(id); return n; });
  };

  const saveObjectives = async () => {
    if (!editItem) return;
    setSaving(true);
    try {
      const newIds = [...selectedObjectives].filter((id) => !periodObjectives.includes(id));
      for (const objId of newIds) {
        await addReviewPeriodObjective({ reviewPeriodId: editItem.periodId, enterpriseObjectiveId: objId });
      }
      if (newIds.length > 0) toast.success(`${newIds.length} objectives added.`);
      setPeriodObjectives([...selectedObjectives]);
      setStep(3);
    } catch { toast.error("Failed to save objectives."); } finally { setSaving(false); }
  };

  // Step 3: Add category definition
  const addCategoryDef = async () => {
    if (!editItem || !catForm.objectiveCategoryId || !catForm.weight) { toast.error("Category and weight are required."); return; }
    setSaving(true);
    try {
      const payload = {
        reviewPeriodId: editItem.periodId, objectiveCategoryId: catForm.objectiveCategoryId,
        gradeGroupId: catForm.gradeGroupId ? Number(catForm.gradeGroupId) : undefined,
        weight: Number(catForm.weight), maxNoObjectives: Number(catForm.maxNoObjectives) || 5,
        maxNoWorkProduct: Number(catForm.maxNoWorkProduct) || 10, enforceWorkProductLimit: catForm.enforceWorkProductLimit,
      };
      const res = await saveDraftCategoryDefinition(payload);
      if (res?.isSuccess) {
        toast.success("Category definition added.");
        setCatForm({ objectiveCategoryId: "", gradeGroupId: "", weight: "", maxNoObjectives: "", maxNoWorkProduct: "", enforceWorkProductLimit: false });
        const catDefRes = await getReviewPeriodCategoryDefinitions(editItem.periodId);
        if (catDefRes?.data) setCatDefs(Array.isArray(catDefRes.data) ? catDefRes.data : []);
      } else toast.error(res?.message || "Failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const finishWizard = () => { setWizardOpen(false); loadData(); };

  // Settings: Toggle functions
  const handleToggle = async (flag: string, current: boolean) => {
    if (!settingsItem) return;
    setToggling(flag);
    const payload = { periodId: settingsItem.periodId };
    try {
      let res;
      if (flag === "objectivePlanning") res = current ? await disableObjectivePlanning(payload) : await enableObjectivePlanning(payload);
      else if (flag === "workProductPlanning") res = current ? await disableWorkProductPlanning(payload) : await enableWorkProductPlanning(payload);
      else res = current ? await disableWorkProductEvaluation(payload) : await enableWorkProductEvaluation(payload);
      if (res?.isSuccess) {
        toast.success("Setting updated.");
        setSettingsItem({ ...settingsItem, [flag === "objectivePlanning" ? "allowObjectivePlanning" : flag === "workProductPlanning" ? "allowWorkProductPlanning" : "allowWorkProductEvaluation"]: !current });
        setItems((prev) => prev.map((i) => i.periodId === settingsItem.periodId ? { ...i, [flag === "objectivePlanning" ? "allowObjectivePlanning" : flag === "workProductPlanning" ? "allowWorkProductPlanning" : "allowWorkProductEvaluation"]: !current } : i));
      } else toast.error(res?.message || "Toggle failed.");
    } catch { toast.error("An error occurred."); } finally { setToggling(""); }
  };

  // Lifecycle actions
  const handleLifecycleAction = async () => {
    if (!confirmAction) return;
    setActionLoading(true);
    const payload = { periodId: confirmAction.periodId };
    try {
      let res;
      switch (confirmAction.action) {
        case "submit": res = await submitDraftReviewPeriod(payload); break;
        case "approve": res = await approveReviewPeriod(payload); break;
        case "return": res = await returnReviewPeriod(payload); break;
        case "close": res = await closeReviewPeriod(payload); break;
        case "cancel": res = await cancelReviewPeriod(payload); break;
        default: return;
      }
      if (res?.isSuccess) { toast.success(`Review period ${confirmAction.label.toLowerCase()}d.`); loadData(); setSettingsOpen(false); }
      else toast.error(res?.message || "Action failed.");
    } catch { toast.error("An error occurred."); } finally { setActionLoading(false); setConfirmAction(null); }
  };

  const columns: ColumnDef<PerformanceReviewPeriod>[] = [
    { accessorKey: "name", header: "Period Name" },
    { accessorKey: "range", header: "Range", cell: ({ row }) => <Badge variant="outline">{rangeLabels[row.original.range] ?? row.original.range}</Badge> },
    { accessorKey: "startDate", header: "Start Date", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End Date", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "maxPoints", header: "Max Points" },
    { accessorKey: "strategyName", header: "Strategy", cell: ({ row }) => row.original.strategyName ?? "-" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => <StatusBadge status={row.original.recordStatus ?? 1} /> },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button>
          <Button size="sm" variant="ghost" onClick={() => openSettings(row.original)}><Settings2 className="h-3.5 w-3.5" /></Button>
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Review Periods" breadcrumbs={[{ label: "Setup" }, { label: "Review Periods" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Review Periods" description="Manage performance review periods, objectives, and category definitions"
        breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Review Periods" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Review Period</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search periods..." />

      {/* Multi-step Wizard Sheet */}
      <Sheet open={wizardOpen} onOpenChange={setWizardOpen}>
        <SheetContent className="sm:max-w-xl overflow-y-auto">
          <SheetHeader>
            <SheetTitle>{editItem ? "Edit Review Period" : "New Review Period"}</SheetTitle>
          </SheetHeader>

          {/* Step indicators */}
          <div className="flex items-center gap-2 my-4">
            {[1, 2, 3].map((s) => (
              <div key={s} className="flex items-center gap-2">
                <div className={`flex h-8 w-8 items-center justify-center rounded-full text-sm font-medium ${step === s ? "bg-primary text-primary-foreground" : step > s ? "bg-primary/20 text-primary" : "bg-muted text-muted-foreground"}`}>{s}</div>
                <span className={`text-sm ${step === s ? "font-medium" : "text-muted-foreground"}`}>{s === 1 ? "Details" : s === 2 ? "Objectives" : "Categories"}</span>
                {s < 3 && <ChevronRight className="h-4 w-4 text-muted-foreground" />}
              </div>
            ))}
          </div>
          <Separator className="mb-4" />

          {/* Step 1: Period Details */}
          {step === 1 && (
            <div className="space-y-4">
              <div className="space-y-2"><Label>Period Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
              <div className="space-y-2"><Label>Short Name</Label><Input value={formData.shortName} onChange={(e) => setFormData({ ...formData, shortName: e.target.value })} /></div>
              <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
              <div className="space-y-2">
                <Label>Strategy</Label>
                <Select value={formData.strategyId} onValueChange={(v) => setFormData({ ...formData, strategyId: v })}>
                  <SelectTrigger><SelectValue placeholder="Select strategy" /></SelectTrigger>
                  <SelectContent>{strategies.map((s) => <SelectItem key={s.strategyId} value={s.strategyId}>{s.name}</SelectItem>)}</SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Period Range *</Label>
                <Select value={formData.range} onValueChange={(v) => setFormData({ ...formData, range: v })}>
                  <SelectTrigger><SelectValue placeholder="Select range" /></SelectTrigger>
                  <SelectContent>
                    {Object.entries(rangeLabels).map(([k, v]) => <SelectItem key={k} value={k}>{v}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-2"><Label>Start Date *</Label><Input type="date" value={formData.startDate} onChange={(e) => setFormData({ ...formData, startDate: e.target.value })} /></div>
                <div className="space-y-2"><Label>End Date *</Label><Input type="date" value={formData.endDate} onChange={(e) => setFormData({ ...formData, endDate: e.target.value })} /></div>
              </div>
              <div className="grid grid-cols-3 gap-3">
                <div className="space-y-2"><Label>Max Points</Label><Input type="number" value={formData.maxPoints} onChange={(e) => setFormData({ ...formData, maxPoints: e.target.value })} /></div>
                <div className="space-y-2"><Label>Min Objectives</Label><Input type="number" value={formData.minNoOfObjectives} onChange={(e) => setFormData({ ...formData, minNoOfObjectives: e.target.value })} /></div>
                <div className="space-y-2"><Label>Max Objectives</Label><Input type="number" value={formData.maxNoOfObjectives} onChange={(e) => setFormData({ ...formData, maxNoOfObjectives: e.target.value })} /></div>
              </div>
              <div className="flex gap-3 pt-4">
                <Button variant="outline" className="flex-1" onClick={() => setWizardOpen(false)}>Cancel</Button>
                <Button className="flex-1" onClick={handleSaveStep1} disabled={saving}>
                  {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  {editItem ? "Update & Next" : "Save Draft & Next"}
                  <ChevronRight className="ml-2 h-4 w-4" />
                </Button>
              </div>
            </div>
          )}

          {/* Step 2: Enterprise Objectives */}
          {step === 2 && (
            <div className="space-y-4">
              <p className="text-sm text-muted-foreground">Select enterprise objectives to associate with this review period.</p>
              <div className="max-h-[400px] overflow-y-auto border rounded-md">
                {allObjectives.length === 0 ? (
                  <p className="p-4 text-sm text-muted-foreground text-center">No enterprise objectives found.</p>
                ) : allObjectives.map((obj) => (
                  <div key={obj.enterpriseObjectiveId} className="flex items-start gap-3 p-3 border-b last:border-b-0 hover:bg-muted/50">
                    <Checkbox checked={selectedObjectives.has(obj.enterpriseObjectiveId)} onCheckedChange={() => toggleObjective(obj.enterpriseObjectiveId)} className="mt-0.5" />
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium">{obj.name}</p>
                      {obj.kpi && <p className="text-xs text-muted-foreground">KPI: {obj.kpi}</p>}
                      {obj.categoryName && <Badge variant="outline" className="mt-1 text-xs">{obj.categoryName}</Badge>}
                    </div>
                  </div>
                ))}
              </div>
              <p className="text-xs text-muted-foreground">{selectedObjectives.size} objective(s) selected</p>
              <div className="flex gap-3 pt-4">
                <Button variant="outline" onClick={() => setStep(1)}><ChevronLeft className="mr-2 h-4 w-4" />Back</Button>
                <Button className="flex-1" onClick={saveObjectives} disabled={saving}>
                  {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Save & Next<ChevronRight className="ml-2 h-4 w-4" />
                </Button>
              </div>
            </div>
          )}

          {/* Step 3: Category Definitions */}
          {step === 3 && (
            <div className="space-y-4">
              <p className="text-sm text-muted-foreground">Define objective categories, weights, and limits for this period.</p>

              {/* Existing definitions */}
              {catDefs.length > 0 && (
                <div className="space-y-2">
                  {catDefs.map((cd) => (
                    <Card key={cd.definitionId} className="p-3">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-sm font-medium">{cd.categoryName ?? cd.objectiveCategoryId}</p>
                          <p className="text-xs text-muted-foreground">Weight: {cd.weight}% | Max Obj: {cd.maxNoObjectives} | Max WP: {cd.maxNoWorkProduct}</p>
                        </div>
                        <Badge variant="outline">{cd.weight}%</Badge>
                      </div>
                    </Card>
                  ))}
                  <p className="text-xs text-muted-foreground">Total weight: {catDefs.reduce((sum, cd) => sum + cd.weight, 0)}%</p>
                </div>
              )}

              <Separator />

              {/* Add new definition */}
              <Card>
                <CardHeader className="pb-3"><CardTitle className="text-sm">Add Category Definition</CardTitle></CardHeader>
                <CardContent className="space-y-3">
                  <div className="space-y-2">
                    <Label>Objective Category *</Label>
                    <Select value={catForm.objectiveCategoryId} onValueChange={(v) => setCatForm({ ...catForm, objectiveCategoryId: v })}>
                      <SelectTrigger><SelectValue placeholder="Select category" /></SelectTrigger>
                      <SelectContent>{categories.map((c) => <SelectItem key={c.objectiveCategoryId} value={c.objectiveCategoryId}>{c.name}</SelectItem>)}</SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label>Grade Group</Label>
                    <Select value={catForm.gradeGroupId} onValueChange={(v) => setCatForm({ ...catForm, gradeGroupId: v })}>
                      <SelectTrigger><SelectValue placeholder="Select grade group" /></SelectTrigger>
                      <SelectContent>{gradeGroups.map((g) => <SelectItem key={g.jobGradeGroupId} value={String(g.jobGradeGroupId)}>{g.groupName}</SelectItem>)}</SelectContent>
                    </Select>
                  </div>
                  <div className="grid grid-cols-3 gap-3">
                    <div className="space-y-2"><Label>Weight % *</Label><Input type="number" value={catForm.weight} onChange={(e) => setCatForm({ ...catForm, weight: e.target.value })} /></div>
                    <div className="space-y-2"><Label>Max Obj.</Label><Input type="number" value={catForm.maxNoObjectives} onChange={(e) => setCatForm({ ...catForm, maxNoObjectives: e.target.value })} /></div>
                    <div className="space-y-2"><Label>Max WP</Label><Input type="number" value={catForm.maxNoWorkProduct} onChange={(e) => setCatForm({ ...catForm, maxNoWorkProduct: e.target.value })} /></div>
                  </div>
                  <div className="flex items-center gap-3">
                    <Switch checked={catForm.enforceWorkProductLimit} onCheckedChange={(v) => setCatForm({ ...catForm, enforceWorkProductLimit: v })} />
                    <Label>Enforce WP Limit</Label>
                  </div>
                  <Button size="sm" onClick={addCategoryDef} disabled={saving}>
                    {saving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Plus className="mr-2 h-4 w-4" />}Add Definition
                  </Button>
                </CardContent>
              </Card>

              <div className="flex gap-3 pt-4">
                <Button variant="outline" onClick={() => setStep(2)}><ChevronLeft className="mr-2 h-4 w-4" />Back</Button>
                <Button className="flex-1" onClick={finishWizard}><CheckCircle className="mr-2 h-4 w-4" />Finish</Button>
              </div>
            </div>
          )}
        </SheetContent>
      </Sheet>

      {/* Settings Drawer */}
      <Sheet open={settingsOpen} onOpenChange={setSettingsOpen}>
        <SheetContent className="sm:max-w-md">
          <SheetHeader>
            <SheetTitle>Period Settings</SheetTitle>
          </SheetHeader>
          {settingsItem && (
            <div className="space-y-6 mt-4">
              <div>
                <h3 className="font-medium">{settingsItem.name}</h3>
                <p className="text-sm text-muted-foreground">{settingsItem.startDate?.split("T")[0]} to {settingsItem.endDate?.split("T")[0]}</p>
                <div className="mt-2"><StatusBadge status={settingsItem.recordStatus ?? 1} /></div>
              </div>
              <Separator />

              {/* Lifecycle Actions */}
              <div className="space-y-3">
                <h4 className="text-sm font-medium">Lifecycle Actions</h4>
                <div className="flex flex-wrap gap-2">
                  {settingsItem.recordStatus === Status.Draft && (
                    <Button size="sm" onClick={() => setConfirmAction({ action: "submit", label: "Submit", periodId: settingsItem.periodId })}>
                      <CheckCircle className="mr-2 h-4 w-4" />Submit for Approval
                    </Button>
                  )}
                  {settingsItem.recordStatus === Status.PendingApproval && (
                    <>
                      <Button size="sm" onClick={() => setConfirmAction({ action: "approve", label: "Approve", periodId: settingsItem.periodId })}>
                        <CheckCircle className="mr-2 h-4 w-4" />Approve
                      </Button>
                      <Button size="sm" variant="outline" onClick={() => setConfirmAction({ action: "return", label: "Return", periodId: settingsItem.periodId })}>
                        <RotateCcw className="mr-2 h-4 w-4" />Return
                      </Button>
                    </>
                  )}
                  {(settingsItem.recordStatus === Status.Active || settingsItem.recordStatus === Status.ApprovedAndActive) && (
                    <Button size="sm" variant="outline" onClick={() => setConfirmAction({ action: "close", label: "Close", periodId: settingsItem.periodId })}>
                      <XCircle className="mr-2 h-4 w-4" />Close Period
                    </Button>
                  )}
                  {settingsItem.recordStatus !== Status.Cancelled && settingsItem.recordStatus !== Status.Closed && (
                    <Button size="sm" variant="destructive" onClick={() => setConfirmAction({ action: "cancel", label: "Cancel", periodId: settingsItem.periodId })}>
                      <Trash2 className="mr-2 h-4 w-4" />Cancel
                    </Button>
                  )}
                </div>
              </div>
              <Separator />

              {/* Feature Toggles */}
              <div className="space-y-4">
                <h4 className="text-sm font-medium">Feature Toggles</h4>
                <div className="flex items-center justify-between">
                  <div><Label>Objective Planning</Label><p className="text-xs text-muted-foreground">Allow staff to plan objectives</p></div>
                  <Switch checked={settingsItem.allowObjectivePlanning} disabled={toggling === "objectivePlanning"} onCheckedChange={() => handleToggle("objectivePlanning", settingsItem.allowObjectivePlanning)} />
                </div>
                <div className="flex items-center justify-between">
                  <div><Label>Work Product Planning</Label><p className="text-xs text-muted-foreground">Allow staff to plan work products</p></div>
                  <Switch checked={settingsItem.allowWorkProductPlanning} disabled={toggling === "workProductPlanning"} onCheckedChange={() => handleToggle("workProductPlanning", settingsItem.allowWorkProductPlanning)} />
                </div>
                <div className="flex items-center justify-between">
                  <div><Label>Work Product Evaluation</Label><p className="text-xs text-muted-foreground">Allow work product evaluation</p></div>
                  <Switch checked={settingsItem.allowWorkProductEvaluation} disabled={toggling === "workProductEvaluation"} onCheckedChange={() => handleToggle("workProductEvaluation", settingsItem.allowWorkProductEvaluation)} />
                </div>
              </div>
              <Separator />

              {/* Period Summary */}
              <div className="space-y-2">
                <h4 className="text-sm font-medium">Configuration</h4>
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div className="text-muted-foreground">Max Points</div><div>{settingsItem.maxPoints}</div>
                  <div className="text-muted-foreground">Min Objectives</div><div>{settingsItem.minNoOfObjectives}</div>
                  <div className="text-muted-foreground">Max Objectives</div><div>{settingsItem.maxNoOfObjectives}</div>
                  <div className="text-muted-foreground">Range</div><div>{rangeLabels[settingsItem.range] ?? settingsItem.range}</div>
                  <div className="text-muted-foreground">Strategy</div><div>{settingsItem.strategyName ?? "-"}</div>
                </div>
              </div>
            </div>
          )}
        </SheetContent>
      </Sheet>

      {/* Confirmation Dialog */}
      <ConfirmationDialog
        open={!!confirmAction}
        onOpenChange={(open) => { if (!open) setConfirmAction(null); }}
        title={`${confirmAction?.label} Review Period`}
        description={`Are you sure you want to ${confirmAction?.label.toLowerCase()} this review period?`}
        confirmLabel={confirmAction?.label ?? "Confirm"}
        onConfirm={handleLifecycleAction}
        variant={confirmAction?.action === "cancel" ? "destructive" : "default"}
      />
    </div>
  );
}
