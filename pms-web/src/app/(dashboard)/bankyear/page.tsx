"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Switch } from "@/components/ui/switch";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getBankYears, saveBankYear } from "@/lib/api/competency";
import type { CompetencyBankYear } from "@/types/competency";

export default function BankYearPage() {
  const [items, setItems] = useState<CompetencyBankYear[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<CompetencyBankYear | null>(null);
  const [formData, setFormData] = useState({ yearName: "", isActive: true });

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getBankYears();
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ yearName: "", isActive: true }); setOpen(true); };
  const openEdit = (item: CompetencyBankYear) => {
    setEditItem(item);
    setFormData({ yearName: item.yearName, isActive: item.isActive });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.yearName) { toast.error("Please enter a year name."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { bankYearId: editItem.bankYearId } : {}),
        yearName: formData.yearName,
        isActive: formData.isActive,
      };
      const res = await saveBankYear(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Bank year updated." : "Bank year created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<CompetencyBankYear>[] = [
    { accessorKey: "yearName", header: "Bank Year" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Disabled"}</Badge> },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Bank Year" breadcrumbs={[{ label: "Organogram" }, { label: "Bank Year" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Bank Year" description="Manage bank year configuration" breadcrumbs={[{ label: "Organogram" }, { label: "Bank Year" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Bank Year</Button>} />
      <DataTable columns={columns} data={items} searchKey="yearName" searchPlaceholder="Search bank years..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Bank Year" : "Add Bank Year"} isEdit={!!editItem} editWarning="Update the selected bank year below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Year Name *</Label><Input value={formData.yearName} onChange={(e) => setFormData({ ...formData, yearName: e.target.value })} placeholder="e.g. 2026" /></div>
          <div className="flex items-center gap-2"><Switch checked={formData.isActive} onCheckedChange={(v) => setFormData({ ...formData, isActive: v })} /><Label>Active</Label></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
