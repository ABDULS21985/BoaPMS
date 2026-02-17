"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Loader2, ClipboardList } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import {
  getCompetencyReviewByReviewer, saveCompetencyReview, getCompetencyReviewPeriods,
  getCompetencies, getRatings,
} from "@/lib/api/competency";
import type { CompetencyReview, CompetencyReviewPeriod, Competency, Rating } from "@/types/competency";

export default function CompetencyReviewsPage() {
  const { data: session } = useSession();
  const reviewerId = session?.user?.id ?? "";
  const [reviews, setReviews] = useState<CompetencyReview[]>([]);
  const [periods, setPeriods] = useState<CompetencyReviewPeriod[]>([]);
  const [competencies, setCompetencies] = useState<Competency[]>([]);
  const [ratings, setRatings] = useState<Rating[]>([]);
  const [selectedPeriod, setSelectedPeriod] = useState<string>("all");
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editItem, setEditItem] = useState<CompetencyReview | null>(null);
  const [formData, setFormData] = useState({ actualRatingId: "" });

  const loadData = async () => {
    if (!reviewerId) return;
    setLoading(true);
    try {
      const [revRes, periodRes, compRes, ratRes] = await Promise.all([
        getCompetencyReviewByReviewer(reviewerId),
        getCompetencyReviewPeriods(),
        getCompetencies(),
        getRatings(),
      ]);
      if (revRes?.data) setReviews(Array.isArray(revRes.data) ? revRes.data : []);
      if (periodRes?.data) setPeriods(Array.isArray(periodRes.data) ? periodRes.data : []);
      if (compRes?.data) setCompetencies(Array.isArray(compRes.data) ? compRes.data : []);
      if (ratRes?.data) setRatings(Array.isArray(ratRes.data) ? ratRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [reviewerId]);

  const filtered = selectedPeriod === "all" ? reviews : reviews.filter((r) => r.reviewPeriodId === Number(selectedPeriod));

  const openRate = (item: CompetencyReview) => {
    setEditItem(item);
    setFormData({ actualRatingId: item.actualRatingId ? String(item.actualRatingId) : "" });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!editItem || !formData.actualRatingId) { toast.error("Please select a rating."); return; }
    setSaving(true);
    try {
      const res = await saveCompetencyReview({
        competencyReviewId: editItem.competencyReviewId,
        actualRatingId: Number(formData.actualRatingId),
        reviewerId,
      });
      if (res?.isSuccess) { toast.success("Review saved."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Save failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const columns: ColumnDef<CompetencyReview>[] = [
    { accessorKey: "employeeName", header: "Employee", cell: ({ row }) => row.original.employeeName ?? row.original.employeeNumber },
    { accessorKey: "competencyName", header: "Competency" },
    { accessorKey: "competencyCategoryName", header: "Category" },
    { accessorKey: "reviewTypeName", header: "Review Type" },
    { accessorKey: "expectedRatingValue", header: "Expected", cell: ({ row }) => <Badge variant="outline">{row.original.expectedRatingName ?? row.original.expectedRatingValue}</Badge> },
    { accessorKey: "actualRatingValue", header: "Actual", cell: ({ row }) => row.original.actualRatingId ? <Badge variant={row.original.actualRatingValue >= row.original.expectedRatingValue ? "default" : "destructive"}>{row.original.actualRatingName ?? row.original.actualRatingValue}</Badge> : <Badge variant="secondary">Not Rated</Badge> },
    { accessorKey: "reviewPeriodName", header: "Period" },
    {
      id: "actions", header: "Actions",
      cell: ({ row }) => (
        <Button size="sm" variant="outline" onClick={() => openRate(row.original)}>Rate</Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Competency Reviews" breadcrumbs={[{ label: "Competency" }, { label: "Reviews" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Competency Reviews" description="Review and rate employee competencies"
        breadcrumbs={[{ label: "Competency" }, { label: "Reviews" }]}
      />

      <div className="flex gap-3">
        <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
          <SelectTrigger className="w-[220px]"><SelectValue placeholder="All Periods" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Periods</SelectItem>
            {periods.map((p) => <SelectItem key={p.reviewPeriodId} value={String(p.reviewPeriodId)}>{p.name}</SelectItem>)}
          </SelectContent>
        </Select>
      </div>

      {filtered.length > 0 ? (
        <DataTable columns={columns} data={filtered} searchKey="employeeName" searchPlaceholder="Search employees..." />
      ) : (
        <EmptyState icon={ClipboardList} title="No Reviews" description="No competency reviews assigned to you." />
      )}

      <FormSheet open={open} onOpenChange={setOpen} title={`Rate: ${editItem?.employeeName ?? ""} - ${editItem?.competencyName ?? ""}`}>
        <div className="space-y-4">
          {editItem?.competencyDefinition && (
            <div className="rounded-md bg-muted p-3 text-sm">{editItem.competencyDefinition}</div>
          )}
          <div className="text-sm text-muted-foreground">Expected: <Badge variant="outline">{editItem?.expectedRatingName ?? editItem?.expectedRatingValue}</Badge></div>
          <div className="space-y-2">
            <Label>Actual Rating *</Label>
            <Select value={formData.actualRatingId} onValueChange={(v) => setFormData({ actualRatingId: v })}>
              <SelectTrigger><SelectValue placeholder="Select rating" /></SelectTrigger>
              <SelectContent>{ratings.map((r) => <SelectItem key={r.ratingId} value={String(r.ratingId)}>{r.name} ({r.value})</SelectItem>)}</SelectContent>
            </Select>
          </div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSave} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save Rating</Button>
          </div>
        </div>
      </FormSheet>
    </div>
  );
}
