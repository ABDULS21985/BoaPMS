"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter, useParams } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { format } from "date-fns";
import { ClipboardCheck, ListTodo, Loader2, RotateCcw } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { FormDialog } from "@/components/shared/form-dialog";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getAllStaffWorkProducts, getWorkProductTasks, evaluateWorkProduct, reInstateWorkProduct } from "@/lib/api/pms-engine";
import type { WorkProduct, WorkProductTask } from "@/types/performance";
import { Status } from "@/types/enums";
import { Roles } from "@/stores/auth-store";

const ADMIN_ROLES: string[] = [Roles.Admin, Roles.SuperAdmin, Roles.HrReportAdmin, Roles.HrAdmin];

export default function WorkProductsReportPage() {
  const { data: session, status } = useSession();
  const router = useRouter();
  const params = useParams();
  const reviewPeriodId = params.reviewPeriodId as string;
  const userRoles = session?.user?.roles ?? [];

  const [items, setItems] = useState<WorkProduct[]>([]);
  const [loading, setLoading] = useState(true);

  // Evaluate dialog
  const [evalOpen, setEvalOpen] = useState(false);
  const [evalWpId, setEvalWpId] = useState("");
  const [evalForm, setEvalForm] = useState({ timeliness: "", quality: "", output: "" });
  const [evalSaving, setEvalSaving] = useState(false);

  // Tasks dialog
  const [tasksOpen, setTasksOpen] = useState(false);
  const [tasks, setTasks] = useState<WorkProductTask[]>([]);
  const [tasksLoading, setTasksLoading] = useState(false);

  // Reinstate dialog
  const [reinstateOpen, setReinstateOpen] = useState(false);
  const [reinstateWpId, setReinstateWpId] = useState("");

  useEffect(() => {
    if (status === "authenticated" && !userRoles.some((r) => ADMIN_ROLES.includes(r))) {
      router.push("/access-denied");
    }
  }, [status, userRoles, router]);

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getAllStaffWorkProducts();
      if (res?.data) {
        const all = Array.isArray(res.data) ? res.data : [];
        setItems(reviewPeriodId ? all.filter((wp) => (wp as WorkProduct & { reviewPeriodId?: string }).reviewPeriodId === reviewPeriodId) : all);
      }
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [reviewPeriodId]);

  const openEval = (wp: WorkProduct) => {
    setEvalWpId(wp.workProductId);
    setEvalForm({ timeliness: "", quality: "", output: "" });
    setEvalOpen(true);
  };

  const handleEval = async () => {
    if (!evalForm.timeliness || !evalForm.quality || !evalForm.output) { toast.error("Please fill all evaluation fields."); return; }
    setEvalSaving(true);
    try {
      const res = await evaluateWorkProduct({ workProductId: evalWpId, timeliness: Number(evalForm.timeliness), quality: Number(evalForm.quality), output: Number(evalForm.output) });
      if (res?.isSuccess) { toast.success("Evaluation updated."); setEvalOpen(false); loadData(); }
      else toast.error(res?.message || "Evaluation failed.");
    } catch { toast.error("An error occurred."); } finally { setEvalSaving(false); }
  };

  const openTasks = async (wp: WorkProduct) => {
    setTasksOpen(true);
    setTasksLoading(true);
    try {
      const res = await getWorkProductTasks(wp.workProductId);
      if (res?.data) setTasks(Array.isArray(res.data) ? res.data : []);
      else setTasks([]);
    } catch { toast.error("Failed to load tasks."); } finally { setTasksLoading(false); }
  };

  const handleReinstate = async () => {
    const res = await reInstateWorkProduct({ workProductId: reinstateWpId });
    if (res?.isSuccess) { toast.success("Work product re-instated."); loadData(); }
    else toast.error(res?.message || "Failed to re-instate.");
  };

  const columns: ColumnDef<WorkProduct>[] = [
    { id: "index", header: "#", cell: ({ row }) => row.index + 1 },
    { accessorKey: "workProductId", header: "ID", cell: ({ row }) => <span title={row.original.workProductId}>{row.original.workProductId?.slice(0, 8)}...</span> },
    { accessorKey: "staffId", header: "Staff ID" },
    { accessorKey: "objectiveName", header: "Objective", cell: ({ row }) => (row.original as WorkProduct & { objectiveName?: string }).objectiveName || "-" },
    { accessorKey: "name", header: "Work Product" },
    { accessorKey: "workProductType", header: "Type" },
    { accessorKey: "maxPoint", header: "Max Point" },
    { accessorKey: "finalScore", header: "Score" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => <StatusBadge status={row.original.recordStatus ?? 0} /> },
    {
      id: "actions", header: "Actions", cell: ({ row }) => {
        const wp = row.original;
        const canReinstate = wp.recordStatus === Status.Rejected || wp.recordStatus === Status.Cancelled;
        return (
          <div className="flex gap-1">
            <Button size="sm" variant="ghost" onClick={() => openEval(wp)} title="Evaluate"><ClipboardCheck className="h-3.5 w-3.5" /></Button>
            <Button size="sm" variant="ghost" onClick={() => openTasks(wp)} title="View Tasks"><ListTodo className="h-3.5 w-3.5" /></Button>
            {canReinstate && <Button size="sm" variant="ghost" onClick={() => { setReinstateWpId(wp.workProductId); setReinstateOpen(true); }} title="Re-instate"><RotateCcw className="h-3.5 w-3.5" /></Button>}
          </div>
        );
      },
    },
  ];

  if (loading) return <div><PageHeader title="Staff Work Products" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Reports" }, { label: "Objectives & Work Products", href: "/review-period-report" }, { label: "Work Products" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Staff Work Products" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Reports" }, { label: "Objectives & Work Products", href: "/review-period-report" }, { label: "Work Products" }]} />
      <DataTable columns={columns} data={items} searchKey="staffId" searchPlaceholder="Search by staff ID or work product name..." />

      {/* Evaluate Dialog */}
      <FormDialog open={evalOpen} onOpenChange={setEvalOpen} title="Update Evaluation">
        <div className="space-y-4">
          {(["timeliness", "quality", "output"] as const).map((field) => (
            <div key={field} className="space-y-2">
              <Label className="capitalize">{field}</Label>
              <Select value={evalForm[field]} onValueChange={(v) => setEvalForm({ ...evalForm, [field]: v })}>
                <SelectTrigger><SelectValue placeholder={`Select ${field} score`} /></SelectTrigger>
                <SelectContent>{[1, 2, 3, 4, 5].map((n) => <SelectItem key={n} value={String(n)}>{n}</SelectItem>)}</SelectContent>
              </Select>
            </div>
          ))}
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setEvalOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleEval} disabled={evalSaving}>{evalSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Submit</Button>
          </div>
        </div>
      </FormDialog>

      {/* Tasks Dialog */}
      <FormDialog open={tasksOpen} onOpenChange={setTasksOpen} title="Work Product Tasks" className="sm:max-w-2xl">
        {tasksLoading ? (
          <div className="flex items-center justify-center py-8"><Loader2 className="h-6 w-6 animate-spin text-muted-foreground" /></div>
        ) : tasks.length > 0 ? (
          <div className="overflow-auto">
            <table className="w-full text-sm">
              <thead><tr className="border-b"><th className="py-2 px-3 text-left font-medium">Task Name</th><th className="py-2 px-3 text-left font-medium">Status</th><th className="py-2 px-3 text-left font-medium">Start Date</th><th className="py-2 px-3 text-left font-medium">End Date</th></tr></thead>
              <tbody>
                {tasks.map((t) => (
                  <tr key={t.workProductTaskId} className="border-b">
                    <td className="py-2 px-3">{t.name}</td>
                    <td className="py-2 px-3"><StatusBadge status={t.recordStatus ?? 0} /></td>
                    <td className="py-2 px-3">{t.startDate ? format(new Date(t.startDate), "dd MMM yyyy") : "-"}</td>
                    <td className="py-2 px-3">{t.endDate ? format(new Date(t.endDate), "dd MMM yyyy") : "-"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No tasks found for this work product.</p>
        )}
        <div className="flex justify-end pt-4">
          <Button variant="outline" onClick={() => setTasksOpen(false)}>Close</Button>
        </div>
      </FormDialog>

      {/* Reinstate Dialog */}
      <ConfirmationDialog open={reinstateOpen} onOpenChange={setReinstateOpen} title="Re-instate Work Product" description="Are you sure you want to re-instate this work product?" confirmLabel="Re-instate" onConfirm={handleReinstate} />
    </div>
  );
}
