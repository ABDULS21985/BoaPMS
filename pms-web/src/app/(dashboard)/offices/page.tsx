"use client";

import { Suspense, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getDivisions, getOffices, saveOffice } from "@/lib/api/organogram";
import type { Office, Division } from "@/types/organogram";

export default function OfficesPageWrapper() {
  return <Suspense fallback={<div><PageHeader title="Offices" breadcrumbs={[{ label: "Organogram" }, { label: "Offices" }]} /><PageSkeleton /></div>}><OfficesPage /></Suspense>;
}

function OfficesPage() {
  const searchParams = useSearchParams();
  const filterDivId = searchParams.get("divisionId");

  const [items, setItems] = useState<Office[]>([]);
  const [divisions, setDivisions] = useState<Division[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<Office | null>(null);
  const [formData, setFormData] = useState({ officeName: "", officeCode: "", divisionId: "", isActive: true });

  const loadData = async () => {
    setLoading(true);
    try {
      const [offRes, divRes] = await Promise.all([getOffices(), getDivisions()]);
      if (offRes?.data) {
        const all = Array.isArray(offRes.data) ? offRes.data : [];
        setItems(filterDivId ? all.filter((o) => o.divisionId === Number(filterDivId)) : all);
      }
      if (divRes?.data) setDivisions(Array.isArray(divRes.data) ? divRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [filterDivId]);

  const openAdd = () => {
    setEditItem(null);
    setFormData({ officeName: "", officeCode: "", divisionId: filterDivId ?? "", isActive: true });
    setOpen(true);
  };

  const openEdit = (item: Office) => {
    setEditItem(item);
    setFormData({
      officeName: item.officeName,
      officeCode: item.officeCode ?? "",
      divisionId: String(item.divisionId),
      isActive: item.isActive,
    });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.officeName || !formData.divisionId) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const payload = {
        ...(editItem ? { officeId: editItem.officeId } : {}),
        officeName: formData.officeName,
        officeCode: formData.officeCode,
        divisionId: Number(formData.divisionId),
        isActive: formData.isActive,
      };
      const res = await saveOffice(payload);
      if (res?.isSuccess) { toast.success(editItem ? "Office updated." : "Office created."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<Office>[] = [
    { accessorKey: "divisionName", header: "Division" },
    { accessorKey: "officeCode", header: "Code" },
    { accessorKey: "officeName", header: "Office Name" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Disabled"}</Badge> },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <Button size="sm" variant="ghost" onClick={() => openEdit(row.original)}><Pencil className="h-3.5 w-3.5" /></Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Offices" breadcrumbs={[{ label: "Organogram" }, { label: "Offices" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Offices" description="Manage offices under divisions" breadcrumbs={[{ label: "Organogram" }, { label: "Divisions", href: "/divisions" }, { label: "Offices" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Office</Button>} />
      <DataTable columns={columns} data={items} searchKey="officeName" searchPlaceholder="Search offices..." />

      <FormSheet open={open} onOpenChange={setOpen} title={editItem ? "Edit Office" : "Add Office"} isEdit={!!editItem} editWarning="Update the selected office below.">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Division *</Label>
            <Select value={formData.divisionId} onValueChange={(v) => setFormData({ ...formData, divisionId: v })}>
              <SelectTrigger><SelectValue placeholder="Select division" /></SelectTrigger>
              <SelectContent>{divisions.map((d) => <SelectItem key={d.divisionId} value={String(d.divisionId)}>{d.divisionName}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Office Code</Label><Input value={formData.officeCode} onChange={(e) => setFormData({ ...formData, officeCode: e.target.value })} /></div>
          <div className="space-y-2"><Label>Office Name *</Label><Input value={formData.officeName} onChange={(e) => setFormData({ ...formData, officeName: e.target.value })} /></div>
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
