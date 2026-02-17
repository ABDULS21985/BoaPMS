"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2, Lock, Unlock } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getPmsConfigurations, createPmsConfiguration, updatePmsConfiguration } from "@/lib/api/performance";
import type { PmsConfiguration } from "@/types/performance";

const configTypes = ["String", "Bool", "Int", "Decimal", "DateTime"];

export default function PmsConfigPage() {
  const [items, setItems] = useState<PmsConfiguration[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<PmsConfiguration | null>(null);
  const [formData, setFormData] = useState({ name: "", value: "", type: "", isEncrypted: false });

  const loadData = async () => {
    setLoading(true);
    try { const res = await getPmsConfigurations(); if (res?.data) setItems(Array.isArray(res.data) ? res.data : []); }
    catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ name: "", value: "", type: "", isEncrypted: false }); setOpen(true); };
  const openEdit = (item: PmsConfiguration) => {
    setEditItem(item);
    setFormData({ name: item.name, value: item.value ?? "", type: item.type ?? "", isEncrypted: item.isEncrypted });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.value || !formData.type) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const res = editItem
        ? await updatePmsConfiguration({ pmsConfigurationId: editItem.pmsConfigurationId, ...formData })
        : await createPmsConfiguration(formData);
      if (res?.isSuccess) { toast.success(editItem ? "Configuration updated." : "Configuration created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<PmsConfiguration>[] = [
    { accessorKey: "name", header: "Name" },
    { accessorKey: "value", header: "Value", cell: ({ row }) => row.original.isEncrypted ? <span className="text-muted-foreground italic">***encrypted***</span> : <span className="line-clamp-1">{row.original.value}</span> },
    { accessorKey: "type", header: "Type", cell: ({ row }) => <Badge variant="outline">{row.original.type}</Badge> },
    { accessorKey: "isEncrypted", header: "Encrypted", cell: ({ row }) => row.original.isEncrypted ? <Lock className="h-4 w-4 text-amber-500" /> : <Unlock className="h-4 w-4 text-muted-foreground" /> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="PMS Configurations" breadcrumbs={[{ label: "Setup" }, { label: "PMS Config" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="PMS Configurations" description="Manage PMS configuration entries" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "PMS Config" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Configuration</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search configurations..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Configuration" : "Add Configuration"} isEdit={!!editItem} editWarning="Update the selected configuration below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Value *</Label><Input type={formData.isEncrypted ? "password" : "text"} value={formData.value} onChange={(e) => setFormData({ ...formData, value: e.target.value })} /></div>
          <div className="space-y-2">
            <Label>Type *</Label>
            <Select value={formData.type} onValueChange={(v) => setFormData({ ...formData, type: v })}>
              <SelectTrigger><SelectValue placeholder="Select type" /></SelectTrigger>
              <SelectContent>{configTypes.map((t) => <SelectItem key={t} value={t}>{t}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="flex items-center gap-3">
            <Switch checked={formData.isEncrypted} onCheckedChange={(v) => setFormData({ ...formData, isEncrypted: v })} />
            <Label>Encrypted</Label>
          </div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
