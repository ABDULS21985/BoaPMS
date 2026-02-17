"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { CheckCircle, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getConsolidatedObjectives, approveRecords } from "@/lib/api/performance";

interface ConsolidatedObjective {
  objectiveId: string;
  objectiveLevel: number;
  name: string;
  description?: string;
  kpi?: string;
  target?: string;
  sbuName?: string;
  levelName?: string;
  recordStatus?: number;
}

const levelNames: Record<number, string> = { 1: "Department", 2: "Division", 3: "Office", 4: "Enterprise" };

export default function ObjectiveApprovalsPage() {
  const [items, setItems] = useState<ConsolidatedObjective[]>([]);
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [confirmApprove, setConfirmApprove] = useState(false);
  const [approving, setApproving] = useState(false);

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getConsolidatedObjectives();
      if (res?.data) {
        const all = Array.isArray(res.data) ? res.data as ConsolidatedObjective[] : [];
        setItems(all.filter((i) => i.recordStatus !== 3)); // exclude already approved
      }
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const toggleSelect = (id: string) => setSelected((prev) => { const n = new Set(prev); n.has(id) ? n.delete(id) : n.add(id); return n; });
  const toggleAll = () => setSelected(selected.size === items.length ? new Set() : new Set(items.map((i) => i.objectiveId)));

  const handleApprove = async () => {
    setApproving(true);
    try {
      const res = await approveRecords({ entityType: "ConsolidatedObjective", recordIds: Array.from(selected) });
      if (res?.isSuccess) { toast.success(`${selected.size} objectives approved.`); setSelected(new Set()); loadData(); }
      else toast.error(res?.message || "Approval failed.");
    } catch { toast.error("An error occurred."); } finally { setApproving(false); setConfirmApprove(false); }
  };

  const columns: ColumnDef<ConsolidatedObjective>[] = [
    { id: "select", header: () => <Checkbox checked={items.length > 0 && selected.size === items.length} onCheckedChange={toggleAll} />, cell: ({ row }) => <Checkbox checked={selected.has(row.original.objectiveId)} onCheckedChange={() => toggleSelect(row.original.objectiveId)} /> },
    { accessorKey: "sbuName", header: "SBU Name" },
    { accessorKey: "objectiveLevel", header: "Level", cell: ({ row }) => <Badge variant="outline">{levelNames[row.original.objectiveLevel] ?? row.original.objectiveLevel}</Badge> },
    { accessorKey: "name", header: "Objective Name" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "kpi", header: "KPI" },
    { accessorKey: "target", header: "Target" },
  ];

  if (loading) return <div><PageHeader title="Pending Objective Approvals" breadcrumbs={[{ label: "Setup" }, { label: "Approvals" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Pending Objective Approvals" description="Review and approve pending objectives" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Approvals" }]}
        actions={selected.size > 0 ? <Button size="sm" onClick={() => setConfirmApprove(true)} disabled={approving}>{approving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <CheckCircle className="mr-2 h-4 w-4" />}Approve ({selected.size})</Button> : undefined}
      />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search objectives..." />
      <ConfirmationDialog open={confirmApprove} onOpenChange={setConfirmApprove} title="Approve Objectives" description={`Approve ${selected.size} selected objectives?`} confirmLabel="Approve All" onConfirm={handleApprove} />
    </div>
  );
}
