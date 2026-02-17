"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2, CheckCircle } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getStrategies, createStrategy, updateStrategy, approveRecords } from "@/lib/api/performance";
import { getBankYears } from "@/lib/api/competency";
import type { Strategy } from "@/types/performance";
import type { CompetencyBankYear } from "@/types/competency";

export default function StrategiesPage() {
  const [items, setItems] = useState<Strategy[]>([]);
  const [bankYears, setBankYears] = useState<CompetencyBankYear[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<Strategy | null>(null);
  const [approveConfirm, setApproveConfirm] = useState<Strategy | null>(null);

  const [formData, setFormData] = useState({ name: "", description: "", bankYearId: "", startDate: "", endDate: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, yearsRes] = await Promise.all([getStrategies(), getBankYears()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (yearsRes?.data) setBankYears(Array.isArray(yearsRes.data) ? yearsRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => {
    setEditItem(null);
    setFormData({ name: "", description: "", bankYearId: "", startDate: "", endDate: "" });
    setOpen(true);
  };

  const openEdit = (item: Strategy) => {
    setEditItem(item);
    setFormData({
      name: item.name,
      description: item.description ?? "",
      bankYearId: String(item.bankYearId),
      startDate: item.startDate?.split("T")[0] ?? "",
      endDate: item.endDate?.split("T")[0] ?? "",
    });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.bankYearId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const payload = { ...formData, bankYearId: Number(formData.bankYearId) };
      const res = editItem
        ? await updateStrategy({ ...payload, strategyId: editItem.strategyId })
        : await createStrategy(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Strategy updated." : "Strategy created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const handleApprove = async () => {
    if (!approveConfirm) return;
    const res = await approveRecords({ entityType: "Strategy", recordIds: [approveConfirm.strategyId] });
    if (res?.isSuccess) { toast.success("Strategy approved."); loadData(); }
    else toast.error(res?.message || "Approval failed.");
  };

  const columns: ColumnDef<Strategy>[] = [
    { accessorKey: "name", header: "Strategy Name" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    {
      accessorKey: "bankYearId", header: "Year",
      cell: ({ row }) => bankYears.find((y) => y.bankYearId === row.original.bankYearId)?.yearName ?? row.original.bankYearId,
    },
    {
      accessorKey: "isActive", header: "Status",
      cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Approved" : "Pending"}</Badge>,
    },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button>
          {!row.original.isActive && (
            <Button size="sm" variant="ghost" onClick={() => setApproveConfirm(row.original)}><CheckCircle className="h-3.5 w-3.5 text-green-600" /></Button>
          )}
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Bank Strategies" breadcrumbs={[{ label: "Setup" }, { label: "Strategies" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Bank Strategies"
        description="Manage organizational strategies"
        breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Strategies" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Strategy</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search strategies..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Strategy" : "Add Strategy"} isEdit={!!editItem} editWarning="Update the selected strategy below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Strategy Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description *</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="space-y-2">
            <Label>Bank Year *</Label>
            <Select value={formData.bankYearId} onValueChange={(v) => setFormData({ ...formData, bankYearId: v })}>
              <SelectTrigger><SelectValue placeholder="Select year" /></SelectTrigger>
              <SelectContent>{bankYears.map((y) => <SelectItem key={y.bankYearId} value={String(y.bankYearId)}>{y.yearName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2"><Label>Start Date *</Label><Input type="date" value={formData.startDate} onChange={(e) => setFormData({ ...formData, startDate: e.target.value })} /></div>
            <div className="space-y-2"><Label>End Date *</Label><Input type="date" value={formData.endDate} onChange={(e) => setFormData({ ...formData, endDate: e.target.value })} /></div>
          </div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>
              {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}
            </Button>
          </div>
        </div>
      </FormSheet>

      <ConfirmationDialog open={!!approveConfirm} onOpenChange={() => setApproveConfirm(null)} title="Approve Strategy" description={`Approve "${approveConfirm?.name}"?`} confirmLabel="Approve" onConfirm={handleApprove} />
    </div>
  );
}
