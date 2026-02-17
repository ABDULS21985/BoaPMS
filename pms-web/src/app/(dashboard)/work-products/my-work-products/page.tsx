"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { Plus, Eye, Loader2, Briefcase } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getStaffWorkProducts, saveDraftWorkProduct } from "@/lib/api/pms-engine";
import { getStaffActiveReviewPeriod } from "@/lib/api/review-periods";
import type { WorkProduct, PerformanceReviewPeriod } from "@/types/performance";
import { WorkProductType } from "@/types/enums";
import type { ColumnDef } from "@tanstack/react-table";

const wpTypeLabels: Record<number, string> = {
  [WorkProductType.Operational]: "Operational",
  [WorkProductType.Project]: "Project",
  [WorkProductType.Committee]: "Committee",
};

export default function MyWorkProductsPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const staffId = session?.user?.id ?? "";
  const [workProducts, setWorkProducts] = useState<WorkProduct[]>([]);
  const [activePeriod, setActivePeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({ name: "", description: "", deliverables: "", startDate: "", endDate: "", maxPoint: "" });

  const loadData = async () => {
    if (!staffId) return;
    setLoading(true);
    try {
      const periodRes = await getStaffActiveReviewPeriod();
      const period = periodRes?.data ?? null;
      setActivePeriod(period);
      if (period) {
        const wpRes = await getStaffWorkProducts(staffId, period.periodId);
        if (wpRes?.data) setWorkProducts(Array.isArray(wpRes.data) ? wpRes.data : []);
      }
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [staffId]);

  const handleSave = async () => {
    if (!formData.name) { toast.error("Please enter a name."); return; }
    setSaving(true);
    try {
      const res = await saveDraftWorkProduct({
        name: formData.name,
        description: formData.description,
        deliverables: formData.deliverables,
        startDate: formData.startDate || undefined,
        endDate: formData.endDate || undefined,
        maxPoint: formData.maxPoint ? Number(formData.maxPoint) : 0,
        staffId,
        reviewPeriodId: activePeriod?.periodId,
        workProductType: WorkProductType.Operational,
      });
      if (res?.isSuccess) { toast.success("Work product saved as draft."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const openAdd = () => { setFormData({ name: "", description: "", deliverables: "", startDate: "", endDate: "", maxPoint: "" }); setOpen(true); };

  const columns: ColumnDef<WorkProduct>[] = [
    { accessorKey: "name", header: "Work Product", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "workProductType", header: "Type", cell: ({ row }) => <Badge variant="secondary">{wpTypeLabels[row.original.workProductType] ?? "—"}</Badge> },
    { accessorKey: "deliverables", header: "Deliverables", cell: ({ row }) => <span className="line-clamp-1">{row.original.deliverables}</span> },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "maxPoint", header: "Max Pts" },
    { accessorKey: "finalScore", header: "Score", cell: ({ row }) => <Badge variant={row.original.finalScore > 0 ? "default" : "secondary"}>{row.original.finalScore}</Badge> },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="outline" onClick={() => router.push(`/work-products/${row.original.workProductId}`)}><Eye className="mr-1 h-3.5 w-3.5" />View</Button> },
  ];

  if (loading) return <div><PageHeader title="My Work Products" breadcrumbs={[{ label: "Work Products" }, { label: "My Work Products" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title="My Work Products"
        description={activePeriod ? `Review period: ${activePeriod.name}` : "No active review period"}
        breadcrumbs={[{ label: "Work Products", href: "/work-products" }, { label: "My Work Products" }]}
        actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />New Work Product</Button>}
      />
      {workProducts.length > 0 ? (
        <DataTable columns={columns} data={workProducts} searchKey="name" searchPlaceholder="Search work products..." />
      ) : (
        <EmptyState icon={Briefcase} title="No Work Products" description="You have no work products for the current review period." />
      )}

      <FormSheet open={open} onOpenChange={setOpen} title="New Work Product" className="sm:max-w-lg overflow-y-auto">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Name *</Label><Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} /></div>
          <div className="space-y-2"><Label>Description</Label><Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} /></div>
          <div className="space-y-2"><Label>Deliverables</Label><Input value={formData.deliverables} onChange={(e) => setFormData({ ...formData, deliverables: e.target.value })} /></div>
          <div className="space-y-2"><Label>Max Points</Label><Input type="number" value={formData.maxPoint} onChange={(e) => setFormData({ ...formData, maxPoint: e.target.value })} /></div>
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-2"><Label>Start Date</Label><Input type="date" value={formData.startDate} onChange={(e) => setFormData({ ...formData, startDate: e.target.value })} /></div>
            <div className="space-y-2"><Label>End Date</Label><Input type="date" value={formData.endDate} onChange={(e) => setFormData({ ...formData, endDate: e.target.value })} /></div>
          </div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save Draft</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
