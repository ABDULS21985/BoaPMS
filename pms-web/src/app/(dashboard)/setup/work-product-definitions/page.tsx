"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getAllWorkProductDefinitions, saveWorkProductDefinitions } from "@/lib/api/performance";
import type { WorkProductDefinition } from "@/types/performance";

const objectiveLevelLabels: Record<string, string> = { "1": "Department", "2": "Division", "3": "Office", "4": "Enterprise" };

export default function WorkProductDefinitionsPage() {
  const [items, setItems] = useState<WorkProductDefinition[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<WorkProductDefinition | null>(null);
  const [formData, setFormData] = useState({ name: "", description: "", deliverables: "", objectiveLevel: "" });

  const loadData = async () => {
    setLoading(true);
    try { const res = await getAllWorkProductDefinitions(); if (res?.data) setItems(Array.isArray(res.data) ? res.data : []); }
    catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ name: "", description: "", deliverables: "", objectiveLevel: "" }); setOpen(true); };
  const openEdit = (item: WorkProductDefinition) => {
    setEditItem(item);
    setFormData({ name: item.name, description: item.description ?? "", deliverables: item.deliverables ?? "", objectiveLevel: String(item.objectiveLevel) });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.name) { toast.error("Name is required."); return; }
    setSaving(true);
    try {
      const payload = { workProductDefinitionId: editItem?.workProductDefinitionId, name: formData.name, description: formData.description, deliverables: formData.deliverables, objectiveLevel: formData.objectiveLevel };
      const res = await saveWorkProductDefinitions([payload]);
      if (res?.isSuccess) { toast.success(editItem ? "Definition updated." : "Definition created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<WorkProductDefinition>[] = [
    { accessorKey: "referenceNo", header: "Ref #" },
    { accessorKey: "name", header: "Work Product Name" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1">{row.original.description}</span> },
    { accessorKey: "deliverables", header: "Deliverables", cell: ({ row }) => <span className="line-clamp-1">{row.original.deliverables}</span> },
    { accessorKey: "objectiveLevel", header: "Level", cell: ({ row }) => <Badge variant="outline">{objectiveLevelLabels[String(row.original.objectiveLevel)] ?? row.original.objectiveLevel}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Work Product Definitions" breadcrumbs={[{ label: "Setup" }, { label: "WP Definitions" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Work Product Definitions" description="Manage work product definition templates" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "WP Definitions" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Definition</Button>} />
      <DataTable columns={columns} data={items} searchKey="name" searchPlaceholder="Search definitions..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Work Product Definition" : "Add Work Product Definition"} isEdit={!!editItem} editWarning="Update the selected definition below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Objective Level</Label>
            <Select value={formData.objectiveLevel} onValueChange={(v) => setFormData({ ...formData, objectiveLevel: v })}>
              <SelectTrigger><SelectValue placeholder="Select level" /></SelectTrigger>
              <SelectContent>{Object.entries(objectiveLevelLabels).map(([k, v]) => <SelectItem key={k} value={k}>{v}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="space-y-2"><Label>Deliverables</Label><Input value={formData.deliverables} onChange={(e) => setFormData({ ...formData, deliverables: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
