"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getJobGradeGroups, saveJobGradeGroup } from "@/lib/api/competency";
import type { JobGradeGroup } from "@/types/competency";

export default function GradeGroupsPage() {
  const [items, setItems] = useState<JobGradeGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<JobGradeGroup | null>(null);
  const [formData, setFormData] = useState({ groupName: "", order: "" });

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getJobGradeGroups();
      if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => { setEditItem(null); setFormData({ groupName: "", order: "" }); setOpen(true); };
  const openEdit = (item: JobGradeGroup) => {
    setEditItem(item);
    setFormData({ groupName: item.groupName, order: String(item.order) });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!formData.groupName) { toast.error("Group name is required."); return; }
    setSaving(true);
    try {
      const payload = editItem
        ? { jobGradeGroupId: editItem.jobGradeGroupId, groupName: formData.groupName, order: Number(formData.order) || 0 }
        : { groupName: formData.groupName, order: Number(formData.order) || 0 };
      const res = await saveJobGradeGroup(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Grade group updated." : "Grade group created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<JobGradeGroup>[] = [
    { accessorKey: "groupName", header: "Group Name" },
    { accessorKey: "order", header: "Order", cell: ({ row }) => <Badge variant="outline">{row.original.order}</Badge> },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Inactive"}</Badge> },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button> },
  ];

  if (loading) return <div><PageHeader title="Grade Groups" breadcrumbs={[{ label: "Competency" }, { label: "Setup" }, { label: "Grade Groups" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Grade Groups" description="Manage job grade groups for behavioral competency mapping"
        breadcrumbs={[{ label: "Competency", href: "/competency/profiles" }, { label: "Setup" }, { label: "Grade Groups" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Grade Group</Button>}
      />
      <DataTable columns={columns} data={items} searchKey="groupName" searchPlaceholder="Search groups..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Grade Group" : "Add Grade Group"} isEdit={!!editItem} editWarning="Update the grade group below.">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Group Name *</Label><Input value={formData.groupName} onChange={(e) => setFormData({ ...formData, groupName: e.target.value })} /></div>
          <div className="space-y-2"><Label>Display Order</Label><Input type="number" value={formData.order} onChange={(e) => setFormData({ ...formData, order: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{editItem ? "Update" : "Save"}</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
