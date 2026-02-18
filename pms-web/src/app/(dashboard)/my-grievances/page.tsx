"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Eye, Loader2, AlertTriangle } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { FormDialog } from "@/components/shared/form-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { StatusBadge } from "@/components/shared/status-badge";
import {
  getStaffGrievances,
  submitGrievance,
  submitGrievanceResolution,
  getGrievanceTypes,
} from "@/lib/api/grievance";
import type { Grievance } from "@/types/performance";

const resolutionLevels: Record<number, string> = {
  0: "None",
  1: "SBU Level",
  2: "Department Level",
  3: "HRD Level",
};

export default function MyGrievancesPage() {
  const { data: session } = useSession();
  const staffId = session?.user?.id ?? "";

  const [items, setItems] = useState<Grievance[]>([]);
  const [grievanceTypes, setGrievanceTypes] = useState<{ id: number; name: string }[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<Grievance | null>(null);

  const [formData, setFormData] = useState({ grievanceType: "", subjectId: "", description: "" });

  // Resolution form
  const [resolveOpen, setResolveOpen] = useState(false);
  const [resolveItem, setResolveItem] = useState<Grievance | null>(null);
  const [resolutionComment, setResolutionComment] = useState("");
  const [resolveSaving, setResolveSaving] = useState(false);

  const loadData = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const [gRes, typesRes] = await Promise.all([getStaffGrievances(staffId), getGrievanceTypes()]);
      if (gRes?.data) setItems(Array.isArray(gRes.data) ? gRes.data : []);
      if (typesRes?.data) setGrievanceTypes(Array.isArray(typesRes.data) ? typesRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [staffId]);

  const openAdd = () => {
    setFormData({ grievanceType: "", subjectId: "", description: "" });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.grievanceType || !formData.description) { toast.error("Please fill required fields."); return; }
    setSaving(true);
    try {
      const res = await submitGrievance({
        grievanceType: Number(formData.grievanceType),
        subjectId: formData.subjectId,
        description: formData.description,
        complainantStaffId: staffId,
      });
      if (res?.isSuccess) { toast.success("Grievance submitted."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Submission failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const openResolve = (item: Grievance) => {
    setResolveItem(item);
    setResolutionComment("");
    setResolveOpen(true);
  };

  const handleResolve = async () => {
    if (!resolveItem || !resolutionComment) { toast.error("Please enter a resolution comment."); return; }
    setResolveSaving(true);
    try {
      const res = await submitGrievanceResolution({
        grievanceId: resolveItem.grievanceId,
        resolutionComment,
        mediatorStaffId: staffId,
      });
      if (res?.isSuccess) { toast.success("Resolution submitted."); setResolveOpen(false); loadData(); }
      else toast.error(res?.message || "Resolution failed.");
    } catch { toast.error("An error occurred."); } finally { setResolveSaving(false); }
  };

  const getGrievanceTypeName = (type: number) => grievanceTypes.find((t) => t.id === type)?.name ?? `Type ${type}`;

  const columns: ColumnDef<Grievance>[] = [
    { accessorKey: "grievanceType", header: "Type", cell: ({ row }) => <Badge variant="secondary">{getGrievanceTypeName(row.original.grievanceType)}</Badge> },
    { accessorKey: "subject", header: "Subject", cell: ({ row }) => row.original.subject ?? row.original.subjectId },
    { accessorKey: "description", header: "Description", cell: ({ row }) => <span className="line-clamp-1 max-w-[200px] inline-block">{row.original.description}</span> },
    { accessorKey: "currentResolutionLevel", header: "Resolution Level", cell: ({ row }) => resolutionLevels[row.original.currentResolutionLevel] ?? "N/A" },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="outline" onClick={() => { setSelectedItem(row.original); setDetailsOpen(true); }}><Eye className="h-3.5 w-3.5" /></Button>
          {row.original.respondentStaffId === staffId && !row.original.respondentComment && (
            <Button size="sm" variant="ghost" onClick={() => openResolve(row.original)} title="Respond">Respond</Button>
          )}
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Grievances" breadcrumbs={[{ label: "Performance" }, { label: "Grievances" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Grievances" description="Manage and track grievances" breadcrumbs={[{ label: "Performance" }, { label: "Grievances" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Raise Grievance</Button>} />

      {items.length > 0 ? (
        <DataTable columns={columns} data={items} searchKey="description" searchPlaceholder="Search grievances..." />
      ) : (
        <EmptyState icon={AlertTriangle} title="No Grievances" description="You have no grievance records." />
      )}

      {/* Add Grievance Sheet */}
      <FormSheet open={open} onOpenChange={setOpen} title="Raise Grievance">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Grievance Type *</Label>
            <Select value={formData.grievanceType} onValueChange={(v) => setFormData({ ...formData, grievanceType: v })}>
              <SelectTrigger><SelectValue placeholder="Select type" /></SelectTrigger>
              <SelectContent>{grievanceTypes.map((t) => <SelectItem key={t.id} value={String(t.id)}>{t.name}</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="space-y-2"><Label>Subject / Target</Label><Input value={formData.subjectId} onChange={(e) => setFormData({ ...formData, subjectId: e.target.value })} placeholder="Enter subject or staff ID" /></div>
          <div className="space-y-2"><Label>Description *</Label><Textarea value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} rows={4} placeholder="Describe the grievance..." /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Submit</Button>
          </div>
        </div>
      </FormSheet>

      {/* Details Dialog */}
      <FormDialog open={detailsOpen} onOpenChange={setDetailsOpen} title="Grievance Details" className="sm:max-w-lg">
        {selectedItem ? (
          <div className="space-y-3 text-sm">
            <div className="grid grid-cols-2 gap-3">
              <div><span className="text-muted-foreground">Type:</span> <span className="font-medium">{getGrievanceTypeName(selectedItem.grievanceType)}</span></div>
              <div><span className="text-muted-foreground">Subject:</span> <span className="font-medium">{selectedItem.subject ?? selectedItem.subjectId}</span></div>
              <div><span className="text-muted-foreground">Complainant:</span> <span className="font-medium">{selectedItem.complainantStaffId}</span></div>
              <div><span className="text-muted-foreground">Respondent:</span> <span className="font-medium">{selectedItem.respondentStaffId}</span></div>
              <div><span className="text-muted-foreground">Resolution Level:</span> <span className="font-medium">{resolutionLevels[selectedItem.currentResolutionLevel] ?? "N/A"}</span></div>
              <div><span className="text-muted-foreground">Status:</span> {selectedItem.recordStatus != null ? <StatusBadge status={selectedItem.recordStatus} /> : "—"}</div>
            </div>
            <div><span className="text-muted-foreground text-xs">Description:</span><p className="text-sm mt-1 bg-muted/50 rounded-md p-2">{selectedItem.description}</p></div>
            {selectedItem.respondentComment && (
              <div><span className="text-muted-foreground text-xs">Respondent Comment:</span><p className="text-sm mt-1 bg-muted/50 rounded-md p-2">{selectedItem.respondentComment}</p></div>
            )}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No details available.</p>
        )}
      </FormDialog>

      {/* Resolution Sheet */}
      <FormSheet open={resolveOpen} onOpenChange={setResolveOpen} title="Respond to Grievance">
        <div className="space-y-4">
          {resolveItem && (
            <div className="rounded-md border p-3 text-sm bg-muted/50">
              <p className="font-medium">{getGrievanceTypeName(resolveItem.grievanceType)}</p>
              <p className="text-muted-foreground line-clamp-2">{resolveItem.description}</p>
            </div>
          )}
          <div className="space-y-2"><Label>Resolution Comment *</Label><Textarea value={resolutionComment} onChange={(e) => setResolutionComment(e.target.value)} rows={4} placeholder="Enter your response..." /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setResolveOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleResolve} disabled={resolveSaving}>{resolveSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Submit Response</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
