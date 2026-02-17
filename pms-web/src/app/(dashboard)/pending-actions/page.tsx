"use client";

import { useEffect, useMemo, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { format, formatDistanceToNow } from "date-fns";
import type { ColumnDef } from "@tanstack/react-table";
import {
  Clock,
  CheckCircle,
  AlertTriangle,
  XCircle,
  Loader2,
  Eye,
  RefreshCw,
} from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Checkbox } from "@/components/ui/checkbox";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { getStaffPendingRequests } from "@/lib/api/dashboard";
import { completeFeedbackRequest } from "@/lib/api/pms-engine";
import { FeedbackRequestType, statusLabels } from "@/types/enums";
import type { PendingAction } from "@/types/dashboard";
import { toast } from "sonner";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

const feedbackTypeLabels: Record<number, string> = {
  [FeedbackRequestType.WorkProductEvaluation]: "Work Product Evaluation",
  [FeedbackRequestType.ObjectivePlanning]: "Objective Planning",
  [FeedbackRequestType.WorkProductPlanning]: "Work Product Planning",
  [FeedbackRequestType.ProjectPlanning]: "Project Planning",
  [FeedbackRequestType.CommitteePlanning]: "Committee Planning",
  [FeedbackRequestType.WorkProductCancellation]: "WP Cancellation",
  [FeedbackRequestType.WorkProductSuspension]: "WP Suspension",
  [FeedbackRequestType.ObjectiveCancellation]: "Objective Cancellation",
  [FeedbackRequestType.ProjectMemberRemoval]: "Project Member Removal",
  [FeedbackRequestType.CommitteeMemberRemoval]: "Committee Member Removal",
  [FeedbackRequestType.ObjectiveSuspension]: "Objective Suspension",
  [FeedbackRequestType.ObjectiveReInstatement]: "Objective Re-Instatement",
  [FeedbackRequestType.WorkProductReInstatement]: "WP Re-Instatement",
  [FeedbackRequestType.WorkProductResumption]: "WP Resumption",
  [FeedbackRequestType.ObjectiveResumption]: "Objective Resumption",
  [FeedbackRequestType.WorkProductReEvaluation]: "WP Re-Evaluation",
  [FeedbackRequestType.ReviewPeriodExtension]: "Review Period Extension",
};

type TabValue = "all" | "workproducts" | "objectives" | "projects" | "committees" | "360reviews";

const TAB_FILTERS: Record<TabValue, number[]> = {
  all: [],
  workproducts: [
    FeedbackRequestType.WorkProductEvaluation,
    FeedbackRequestType.WorkProductPlanning,
    FeedbackRequestType.WorkProductCancellation,
    FeedbackRequestType.WorkProductSuspension,
    FeedbackRequestType.WorkProductReInstatement,
    FeedbackRequestType.WorkProductResumption,
    FeedbackRequestType.WorkProductReEvaluation,
  ],
  objectives: [
    FeedbackRequestType.ObjectivePlanning,
    FeedbackRequestType.ObjectiveCancellation,
    FeedbackRequestType.ObjectiveSuspension,
    FeedbackRequestType.ObjectiveReInstatement,
    FeedbackRequestType.ObjectiveResumption,
  ],
  projects: [
    FeedbackRequestType.ProjectPlanning,
    FeedbackRequestType.ProjectMemberRemoval,
  ],
  committees: [
    FeedbackRequestType.CommitteePlanning,
    FeedbackRequestType.CommitteeMemberRemoval,
  ],
  "360reviews": [],
};

export default function PendingActionsPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const staffId = session?.user?.id ?? "";

  const [loading, setLoading] = useState(true);
  const [actions, setActions] = useState<PendingAction[]>([]);
  const [activeTab, setActiveTab] = useState<TabValue>("all");
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [bulkProcessing, setBulkProcessing] = useState(false);
  const [confirmAction, setConfirmAction] = useState<{ id: string; name: string } | null>(null);

  const loadActions = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const res = await getStaffPendingRequests(staffId);
      if (res?.data) {
        setActions(res.data);
      }
    } catch {
      // Non-fatal
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (staffId) loadActions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [staffId]);

  const filteredActions = useMemo(() => {
    const filterTypes = TAB_FILTERS[activeTab];
    if (!filterTypes || filterTypes.length === 0) return actions;
    return actions.filter((a) => a.feedbackRequestType && filterTypes.includes(a.feedbackRequestType));
  }, [actions, activeTab]);

  const toggleSelect = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const toggleSelectAll = () => {
    if (selected.size === filteredActions.length) {
      setSelected(new Set());
    } else {
      setSelected(new Set(filteredActions.map((a) => a.id)));
    }
  };

  const handleBulkComplete = async () => {
    if (selected.size === 0) return;
    setBulkProcessing(true);
    try {
      const results = await Promise.allSettled(
        Array.from(selected).map((id) =>
          completeFeedbackRequest({ feedbackRequestLogId: id, comment: "Bulk approved" })
        )
      );
      const succeeded = results.filter((r) => r.status === "fulfilled").length;
      toast.success(`${succeeded} of ${selected.size} actions completed successfully.`);
      setSelected(new Set());
      loadActions();
    } catch {
      toast.error("Bulk action failed.");
    } finally {
      setBulkProcessing(false);
    }
  };

  const handleSingleComplete = async (id: string) => {
    try {
      const res = await completeFeedbackRequest({
        feedbackRequestLogId: id,
        comment: "Approved",
      });
      if (res?.isSuccess) {
        toast.success("Action completed successfully.");
        loadActions();
      } else {
        toast.error(res?.message || "Action failed.");
      }
    } catch {
      toast.error("An error occurred.");
    }
    setConfirmAction(null);
  };

  const columns: ColumnDef<PendingAction>[] = [
    {
      id: "select",
      header: () => (
        <Checkbox
          checked={filteredActions.length > 0 && selected.size === filteredActions.length}
          onCheckedChange={toggleSelectAll}
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={selected.has(row.original.id)}
          onCheckedChange={() => toggleSelect(row.original.id)}
        />
      ),
      enableSorting: false,
    },
    {
      accessorKey: "type",
      header: "Request Type",
      cell: ({ row }) => {
        const fType = row.original.feedbackRequestType;
        return (
          <Badge variant="outline" className="text-xs font-normal">
            {fType ? (feedbackTypeLabels[fType] ?? row.original.type) : row.original.type}
          </Badge>
        );
      },
    },
    {
      accessorKey: "name",
      header: "Reference",
      cell: ({ row }) => (
        <div>
          <p className="text-sm font-medium">{row.original.name}</p>
          {row.original.description && (
            <p className="text-xs text-muted-foreground line-clamp-1">{row.original.description}</p>
          )}
        </div>
      ),
    },
    {
      accessorKey: "assignedDate",
      header: "Assigned",
      cell: ({ row }) => {
        const date = row.original.assignedDate;
        if (!date) return "—";
        return (
          <div>
            <p className="text-sm">{format(new Date(date), "dd MMM yyyy")}</p>
            <p className="text-xs text-muted-foreground">
              {formatDistanceToNow(new Date(date), { addSuffix: true })}
            </p>
          </div>
        );
      },
    },
    {
      accessorKey: "dueDate",
      header: "SLA",
      cell: ({ row }) => {
        const due = row.original.dueDate;
        if (!due) return <span className="text-xs text-muted-foreground">No SLA</span>;
        const isOverdue = new Date(due) < new Date();
        return (
          <div className="flex items-center gap-1.5">
            {isOverdue ? (
              <AlertTriangle className="h-3.5 w-3.5 text-destructive" />
            ) : (
              <Clock className="h-3.5 w-3.5 text-muted-foreground" />
            )}
            <span className={`text-xs ${isOverdue ? "font-medium text-destructive" : "text-muted-foreground"}`}>
              {formatDistanceToNow(new Date(due), { addSuffix: true })}
            </span>
          </div>
        );
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => <StatusBadge status={row.original.status} />,
    },
    {
      id: "actions",
      header: "Action",
      cell: ({ row }) => (
        <Button
          size="sm"
          variant="outline"
          onClick={() =>
            setConfirmAction({ id: row.original.id, name: row.original.name })
          }
        >
          <Eye className="mr-1 h-3.5 w-3.5" />
          View
        </Button>
      ),
    },
  ];

  if (loading) {
    return (
      <div>
        <PageHeader title="Pending Actions" breadcrumbs={[{ label: "Pending Actions" }]} />
        <PageSkeleton />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Pending Actions"
        description="Review and act on pending feedback requests, approvals, and evaluations"
        breadcrumbs={[{ label: "Pending Actions" }]}
        actions={
          <div className="flex items-center gap-2">
            {selected.size > 0 && (
              <Button
                size="sm"
                onClick={handleBulkComplete}
                disabled={bulkProcessing}
              >
                {bulkProcessing ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <CheckCircle className="mr-2 h-4 w-4" />
                )}
                Complete ({selected.size})
              </Button>
            )}
            <Button variant="outline" size="sm" onClick={loadActions}>
              <RefreshCw className="mr-2 h-4 w-4" />
              Refresh
            </Button>
          </div>
        }
      />

      {/* Summary Cards */}
      <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-4">
        <Card>
          <CardContent className="flex items-center gap-3 py-3">
            <Clock className="h-5 w-5 text-yellow-500" />
            <div>
              <p className="text-lg font-bold">{actions.length}</p>
              <p className="text-xs text-muted-foreground">Total Pending</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="flex items-center gap-3 py-3">
            <AlertTriangle className="h-5 w-5 text-destructive" />
            <div>
              <p className="text-lg font-bold">
                {actions.filter((a) => a.dueDate && new Date(a.dueDate) < new Date()).length}
              </p>
              <p className="text-xs text-muted-foreground">Overdue</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="flex items-center gap-3 py-3">
            <CheckCircle className="h-5 w-5 text-green-500" />
            <div>
              <p className="text-lg font-bold">
                {actions.filter((a) => a.feedbackRequestType && TAB_FILTERS.workproducts.includes(a.feedbackRequestType)).length}
              </p>
              <p className="text-xs text-muted-foreground">Work Products</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="flex items-center gap-3 py-3">
            <XCircle className="h-5 w-5 text-blue-500" />
            <div>
              <p className="text-lg font-bold">
                {actions.filter((a) => a.feedbackRequestType && TAB_FILTERS.objectives.includes(a.feedbackRequestType)).length}
              </p>
              <p className="text-xs text-muted-foreground">Objectives</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Tabs + Table */}
      <Tabs value={activeTab} onValueChange={(v) => { setActiveTab(v as TabValue); setSelected(new Set()); }}>
        <TabsList>
          <TabsTrigger value="all">All Pending</TabsTrigger>
          <TabsTrigger value="workproducts">Work Products</TabsTrigger>
          <TabsTrigger value="objectives">Objectives</TabsTrigger>
          <TabsTrigger value="projects">Projects</TabsTrigger>
          <TabsTrigger value="committees">Committees</TabsTrigger>
          <TabsTrigger value="360reviews">360 Reviews</TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="mt-4">
          <DataTable
            columns={columns}
            data={filteredActions}
            searchKey="name"
            searchPlaceholder="Search pending actions..."
          />
        </TabsContent>
      </Tabs>

      {/* Action Confirmation Dialog */}
      <AlertDialog open={!!confirmAction} onOpenChange={() => setConfirmAction(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Complete Action</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to complete the action for &ldquo;{confirmAction?.name}&rdquo;?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => confirmAction && handleSingleComplete(confirmAction.id)}
            >
              Complete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
