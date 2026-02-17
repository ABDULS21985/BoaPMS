"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import {
  Send, Plus, Loader2, Briefcase, CheckCircle2, XCircle, Pause, Play,
  RotateCcw, ClipboardCheck, Eye, Pencil, Trash2, ListTodo,
} from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { FormSheet } from "@/components/shared/form-sheet";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { EmptyState } from "@/components/shared/empty-state";
import {
  getWorkProductDetails,
  getWorkProductTasks,
  getWorkProductEvaluation,
  submitDraftWorkProduct,
  approveWorkProduct,
  rejectWorkProduct,
  returnWorkProduct,
  reSubmitWorkProduct,
  completeWorkProduct,
  pauseWorkProduct,
  resumeWorkProduct,
  cancelWorkProduct,
  createWorkProductTask,
  updateWorkProductTask,
  completeWorkProductTask,
  cancelWorkProductTask,
  initiateReEvaluation,
} from "@/lib/api/pms-engine";
import type { WorkProduct, WorkProductTask, WorkProductEvaluation } from "@/types/performance";
import { Status, WorkProductType } from "@/types/enums";

const wpTypeLabels: Record<number, string> = {
  [WorkProductType.Operational]: "Operational",
  [WorkProductType.Project]: "Project",
  [WorkProductType.Committee]: "Committee",
};

export default function WorkProductDetailPage() {
  const { workProductId } = useParams<{ workProductId: string }>();
  const [wp, setWp] = useState<WorkProduct | null>(null);
  const [tasks, setTasks] = useState<WorkProductTask[]>([]);
  const [evaluation, setEvaluation] = useState<WorkProductEvaluation | null>(null);
  const [loading, setLoading] = useState(true);

  // Task form
  const [taskOpen, setTaskOpen] = useState(false);
  const [taskSaving, setTaskSaving] = useState(false);
  const [editingTask, setEditingTask] = useState<WorkProductTask | null>(null);
  const [taskForm, setTaskForm] = useState({ name: "", description: "", startDate: "", endDate: "" });

  // Lifecycle action
  const [actionType, setActionType] = useState<string | null>(null);

  const loadData = async () => {
    if (!workProductId) return;
    setLoading(true);
    try {
      const [wpRes, taskRes, evalRes] = await Promise.all([
        getWorkProductDetails(workProductId),
        getWorkProductTasks(workProductId),
        getWorkProductEvaluation(workProductId),
      ]);
      if (wpRes?.data) setWp(wpRes.data);
      if (taskRes?.data) setTasks(Array.isArray(taskRes.data) ? taskRes.data : []);
      if (evalRes?.data) setEvaluation(evalRes.data as WorkProductEvaluation);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [workProductId]);

  // Task handlers
  const openAddTask = () => {
    setEditingTask(null);
    setTaskForm({ name: "", description: "", startDate: "", endDate: "" });
    setTaskOpen(true);
  };

  const openEditTask = (t: WorkProductTask) => {
    setEditingTask(t);
    setTaskForm({ name: t.name, description: t.description ?? "", startDate: t.startDate?.split("T")[0] ?? "", endDate: t.endDate?.split("T")[0] ?? "" });
    setTaskOpen(true);
  };

  const handleSaveTask = async () => {
    if (!taskForm.name) { toast.error("Enter a task name."); return; }
    setTaskSaving(true);
    try {
      const payload = { ...taskForm, workProductId };
      const res = editingTask
        ? await updateWorkProductTask({ ...payload, workProductTaskId: editingTask.workProductTaskId })
        : await createWorkProductTask(payload);
      if (res?.isSuccess) { toast.success(editingTask ? "Task updated." : "Task added."); setTaskOpen(false); loadData(); }
      else toast.error(res?.message || "Failed.");
    } catch { toast.error("An error occurred."); } finally { setTaskSaving(false); }
  };

  const handleCompleteTask = async (taskId: string) => {
    try {
      const res = await completeWorkProductTask({ workProductTaskId: taskId });
      if (res?.isSuccess) { toast.success("Task completed."); loadData(); }
      else toast.error(res?.message || "Failed.");
    } catch { toast.error("An error occurred."); }
  };

  const handleCancelTask = async (taskId: string) => {
    try {
      const res = await cancelWorkProductTask({ workProductTaskId: taskId });
      if (res?.isSuccess) { toast.success("Task cancelled."); loadData(); }
      else toast.error(res?.message || "Failed.");
    } catch { toast.error("An error occurred."); }
  };

  // Lifecycle actions
  const handleAction = async () => {
    if (!wp || !actionType) return;
    const payload = { workProductId: wp.workProductId };
    const fnMap: Record<string, (d: unknown) => Promise<{ isSuccess?: boolean; message?: string } | null>> = {
      submit: submitDraftWorkProduct,
      approve: approveWorkProduct,
      reject: rejectWorkProduct,
      return: returnWorkProduct,
      resubmit: reSubmitWorkProduct,
      complete: completeWorkProduct,
      pause: pauseWorkProduct,
      resume: resumeWorkProduct,
      cancel: cancelWorkProduct,
    };
    const fn = fnMap[actionType];
    if (!fn) return;
    try {
      const res = await fn(payload);
      if (res?.isSuccess) { toast.success(`Work product ${actionType}ed.`); loadData(); }
      else toast.error(res?.message || "Action failed.");
    } catch { toast.error("An error occurred."); }
  };

  const handleReEvaluation = async () => {
    if (!workProductId) return;
    try {
      const res = await initiateReEvaluation(workProductId);
      if (res?.isSuccess) { toast.success("Re-evaluation initiated."); loadData(); }
      else toast.error(res?.message || "Failed.");
    } catch { toast.error("An error occurred."); }
  };

  const status = wp?.recordStatus;
  const activeTasks = tasks.filter((t) => t.recordStatus !== Status.Cancelled);

  if (loading) return <div><PageHeader title="Work Product Details" breadcrumbs={[{ label: "Work Products" }]} /><PageSkeleton /></div>;
  if (!wp) return <div><PageHeader title="Not Found" breadcrumbs={[{ label: "Work Products" }]} /><EmptyState icon={Briefcase} title="Not Found" description="Work product could not be loaded." /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title={wp.name}
        description={wp.description}
        breadcrumbs={[{ label: "Work Products", href: "/work-products" }, { label: wp.name }]}
        actions={
          <div className="flex flex-wrap gap-2">
            {status === Status.Draft && <Button size="sm" onClick={() => setActionType("submit")}><Send className="mr-2 h-4 w-4" />Submit</Button>}
            {status === Status.Returned && <Button size="sm" onClick={() => setActionType("resubmit")}><RotateCcw className="mr-2 h-4 w-4" />Resubmit</Button>}
            {status === Status.PendingApproval && (
              <>
                <Button size="sm" onClick={() => setActionType("approve")}>Approve</Button>
                <Button size="sm" variant="destructive" onClick={() => setActionType("reject")}>Reject</Button>
                <Button size="sm" variant="outline" onClick={() => setActionType("return")}>Return</Button>
              </>
            )}
            {status === Status.Active && (
              <>
                <Button size="sm" onClick={() => setActionType("complete")}><CheckCircle2 className="mr-2 h-4 w-4" />Complete</Button>
                <Button size="sm" variant="outline" onClick={() => setActionType("pause")}><Pause className="mr-2 h-4 w-4" />Pause</Button>
                <Button size="sm" variant="destructive" onClick={() => setActionType("cancel")}><XCircle className="mr-2 h-4 w-4" />Cancel</Button>
              </>
            )}
            {status === Status.Paused && <Button size="sm" onClick={() => setActionType("resume")}><Play className="mr-2 h-4 w-4" />Resume</Button>}
          </div>
        }
      />

      {/* Overview Cards */}
      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        <Card><CardContent className="pt-4"><div className="flex items-center gap-2">{status != null && <StatusBadge status={status} />}</div><p className="text-xs text-muted-foreground mt-1">Status</p></CardContent></Card>
        <Card><CardContent className="pt-4"><Badge variant="secondary">{wpTypeLabels[wp.workProductType] ?? "—"}</Badge><p className="text-xs text-muted-foreground mt-1">Type</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-2xl font-bold">{wp.maxPoint}</div><p className="text-xs text-muted-foreground">Max Points</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-2xl font-bold">{wp.finalScore}</div><p className="text-xs text-muted-foreground">Final Score</p></CardContent></Card>
      </div>
      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        <Card><CardContent className="pt-4"><div className="text-sm">{wp.startDate?.split("T")[0]}</div><p className="text-xs text-muted-foreground">Start Date</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-sm">{wp.endDate?.split("T")[0]}</div><p className="text-xs text-muted-foreground">End Date</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-sm">{wp.completionDate?.split("T")[0] ?? "—"}</div><p className="text-xs text-muted-foreground">Completion Date</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-2xl font-bold">{activeTasks.length}</div><p className="text-xs text-muted-foreground">Tasks</p></CardContent></Card>
      </div>
      {wp.deliverables && (
        <Card><CardContent className="pt-4"><p className="text-sm text-muted-foreground mb-1">Deliverables</p><p className="text-sm">{wp.deliverables}</p></CardContent></Card>
      )}

      {/* Tabs */}
      <Tabs defaultValue="tasks">
        <TabsList><TabsTrigger value="tasks">Tasks</TabsTrigger><TabsTrigger value="evaluation">Evaluation</TabsTrigger></TabsList>

        {/* Tasks Tab */}
        <TabsContent value="tasks" className="space-y-4">
          {status === Status.Active && (
            <div className="flex justify-end">
              <Button size="sm" variant="outline" onClick={openAddTask}><Plus className="mr-2 h-4 w-4" />Add Task</Button>
            </div>
          )}
          {activeTasks.length > 0 ? (
            <div className="space-y-2">
              {activeTasks.map((t) => (
                <div key={t.workProductTaskId} className="flex items-center justify-between rounded-lg border p-3">
                  <div className="space-y-1">
                    <p className="font-medium">{t.name}</p>
                    {t.description && <p className="text-sm text-muted-foreground line-clamp-1">{t.description}</p>}
                    <div className="flex gap-3 text-xs text-muted-foreground">
                      <span>{t.startDate?.split("T")[0]} — {t.endDate?.split("T")[0]}</span>
                      {t.completionDate && <span>Completed: {t.completionDate.split("T")[0]}</span>}
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {t.recordStatus != null && <StatusBadge status={t.recordStatus} />}
                    {status === Status.Active && t.recordStatus !== Status.Completed && (
                      <>
                        <Button size="sm" variant="ghost" onClick={() => openEditTask(t)}><Pencil className="h-3.5 w-3.5" /></Button>
                        <Button size="sm" variant="ghost" onClick={() => handleCompleteTask(t.workProductTaskId)}><CheckCircle2 className="h-3.5 w-3.5" /></Button>
                        <Button size="sm" variant="ghost" onClick={() => handleCancelTask(t.workProductTaskId)}><Trash2 className="h-3.5 w-3.5" /></Button>
                      </>
                    )}
                  </div>
                </div>
              ))}
            </div>
          ) : <EmptyState icon={ListTodo} title="No Tasks" description="No tasks have been added to this work product." />}
        </TabsContent>

        {/* Evaluation Tab */}
        <TabsContent value="evaluation" className="space-y-4">
          {evaluation ? (
            <Card>
              <CardHeader><CardTitle className="text-base">Evaluation Scores</CardTitle></CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
                  <div><p className="text-xs text-muted-foreground">Timeliness</p><p className="text-2xl font-bold">{evaluation.timeliness}</p></div>
                  <div><p className="text-xs text-muted-foreground">Quality</p><p className="text-2xl font-bold">{evaluation.quality}</p></div>
                  <div><p className="text-xs text-muted-foreground">Output</p><p className="text-2xl font-bold">{evaluation.output}</p></div>
                  <div><p className="text-xs text-muted-foreground">Outcome</p><p className="text-2xl font-bold">{evaluation.outcome}</p></div>
                </div>
                <Separator />
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-muted-foreground">Re-evaluated: {evaluation.isReEvaluated ? "Yes" : "No"}</p>
                    {evaluation.evaluatorStaffId && <p className="text-sm text-muted-foreground">Evaluator: {evaluation.evaluatorStaffId}</p>}
                  </div>
                  {(status === Status.Closed || status === Status.Completed) && (
                    <Button size="sm" variant="outline" onClick={handleReEvaluation}><RotateCcw className="mr-2 h-4 w-4" />Re-evaluate</Button>
                  )}
                </div>
              </CardContent>
            </Card>
          ) : (
            <EmptyState icon={ClipboardCheck} title="No Evaluation" description="This work product has not been evaluated yet." />
          )}
        </TabsContent>
      </Tabs>

      {/* Task Form Sheet */}
      <FormSheet open={taskOpen} onOpenChange={setTaskOpen} title={editingTask ? "Edit Task" : "Add Task"}>
        <div className="space-y-4">
          <div className="space-y-2"><Label>Task Name *</Label><Input value={taskForm.name} onChange={(e) => setTaskForm({ ...taskForm, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={taskForm.description} onChange={(e) => setTaskForm({ ...taskForm, description: e.target.value })} /></div>
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-2"><Label>Start Date</Label><Input type="date" value={taskForm.startDate} onChange={(e) => setTaskForm({ ...taskForm, startDate: e.target.value })} /></div>
            <div className="space-y-2"><Label>End Date</Label><Input type="date" value={taskForm.endDate} onChange={(e) => setTaskForm({ ...taskForm, endDate: e.target.value })} /></div>
          </div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setTaskOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSaveTask} disabled={taskSaving}>{taskSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editingTask ? "Update" : "Add"}</Button>
          </div>
        </div>
      </FormSheet>

      {/* Lifecycle Confirmation Dialog */}
      <ConfirmationDialog
        open={!!actionType}
        onOpenChange={(o) => { if (!o) setActionType(null); }}
        title={`${actionType?.charAt(0).toUpperCase()}${actionType?.slice(1) ?? ""} Work Product`}
        description={`Are you sure you want to ${actionType} "${wp.name}"?`}
        variant={actionType === "reject" || actionType === "cancel" ? "destructive" : "default"}
        confirmLabel={actionType?.charAt(0).toUpperCase() + (actionType?.slice(1) ?? "")}
        showReasonInput={actionType === "reject" || actionType === "return"}
        reasonLabel="Reason"
        onConfirm={handleAction}
      />
    </div>
  );
}
