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
import { getStrategies, getAllStrategicThemes, createStrategicTheme, updateStrategicTheme, approveRecords } from "@/lib/api/performance";
import type { Strategy, StrategicTheme } from "@/types/performance";

export default function StrategicThemesPage() {
  const [items, setItems] = useState<StrategicTheme[]>([]);
  const [strategies, setStrategies] = useState<Strategy[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<StrategicTheme | null>(null);
  const [approveConfirm, setApproveConfirm] = useState<StrategicTheme | null>(null);
  const [formData, setFormData] = useState({ name: "", description: "", strategyId: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const [res, stratRes] = await Promise.all([getAllStrategicThemes(), getStrategies()]);
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      if (stratRes?.data) setStrategies(Array.isArray(stratRes.data) ? stratRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ name: "", description: "", strategyId: "" }); setOpen(true); };
  const openEdit = (item: StrategicTheme) => {
    setEditItem(item);
    setFormData({ name: item.name, description: item.description ?? "", strategyId: item.strategyId });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.strategyId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const res = editItem
        ? await updateStrategicTheme({ ...formData, strategicThemeId: editItem.strategicThemeId })
        : await createStrategicTheme(formData);
      if (res?.isSuccess) { toast.success(editItem ? "Theme updated." : "Theme created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const handleApprove = async () => {
    if (!approveConfirm) return;
    const res = await approveRecords({ entityType: "StrategicTheme", recordIds: [approveConfirm.strategicThemeId] });
    if (res?.isSuccess) { toast.success("Theme approved."); loadData(); } else toast.error(res?.message || "Approval failed.");
  };

  const columns: ColumnDef<StrategicTheme>[] = [
    { accessorKey: "name", header: "Theme Name" },
    { accessorKey: "strategyName", header: "Strategy" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Approved" : "Pending"}</Badge> },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button>
          {!row.original.isActive && <Button size="sm" variant="ghost" onClick={() => setApproveConfirm(row.original)}><CheckCircle className="h-3.5 w-3.5 text-green-600" /></Button>}
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Strategic Themes" breadcrumbs={[{ label: "Setup" }, { label: "Strategic Themes" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Strategic Themes" description="Manage strategic themes under strategies" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Strategic Themes" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Theme</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search themes..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Strategic Theme" : "Add Strategic Theme"} isEdit={!!editItem} editWarning="Update the selected theme below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Strategy *</Label>
            <Select value={formData.strategyId} onValueChange={(v) => setFormData({ ...formData, strategyId: v })}>
              <SelectTrigger><SelectValue placeholder="Select strategy" /></SelectTrigger>
              <SelectContent>{strategies.map((s) => <SelectItem key={s.strategyId} value={s.strategyId}>{s.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Theme Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description *</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>

      <ConfirmationDialog open={!!approveConfirm} onOpenChange={() => setApproveConfirm(null)} title="Approve Theme" description={`Approve "${approveConfirm?.name}"?`} confirmLabel="Approve" onConfirm={handleApprove} />
    </div>
  );
}
