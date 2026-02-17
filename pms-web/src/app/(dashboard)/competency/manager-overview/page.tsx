"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import type { ColumnDef } from "@tanstack/react-table";
import { Users } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getOfficeCompetencyReviews, getCompetencyReviewPeriods } from "@/lib/api/competency";
import { getStaffDetails } from "@/lib/api/staff";
import type { OfficeCompetencyReview, CompetencyReviewPeriod } from "@/types/competency";

export default function ManagerOverviewPage() {
  const { data: session } = useSession();
  const [reviews, setReviews] = useState<OfficeCompetencyReview[]>([]);
  const [periods, setPeriods] = useState<CompetencyReviewPeriod[]>([]);
  const [selectedPeriod, setSelectedPeriod] = useState<string>("");
  const [officeId, setOfficeId] = useState<number>(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      try {
        const [periodRes, profileRes] = await Promise.all([
          getCompetencyReviewPeriods(),
          session?.user?.id ? getStaffDetails(session.user.id) : Promise.resolve(null),
        ]);
        if (periodRes?.data) {
          const ps = Array.isArray(periodRes.data) ? periodRes.data : [];
          setPeriods(ps);
          const approved = ps.find((p) => p.isApproved);
          if (approved) setSelectedPeriod(String(approved.reviewPeriodId));
        }
        if (profileRes?.data) {
          const profile = profileRes.data as { officeId?: number };
          if (profile.officeId) setOfficeId(profile.officeId);
        }
      } catch { /* */ } finally { setLoading(false); }
    };
    init();
  }, [session?.user?.id]);

  useEffect(() => {
    if (!officeId || !selectedPeriod) return;
    setLoading(true);
    getOfficeCompetencyReviews(officeId, Number(selectedPeriod))
      .then((res) => {
        if (res?.data) setReviews(Array.isArray(res.data) ? res.data : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [officeId, selectedPeriod]);

  const completed = reviews.filter((r) => r.isCompleted).length;
  const pending = reviews.length - completed;

  const columns: ColumnDef<OfficeCompetencyReview>[] = [
    { accessorKey: "employeeName", header: "Employee", cell: ({ row }) => row.original.employeeName ?? row.original.employeeNumber },
    { accessorKey: "gradeName", header: "Grade" },
    { accessorKey: "jobRoleName", header: "Job Role" },
    { accessorKey: "departmentName", header: "Department" },
    {
      accessorKey: "isCompleted", header: "Status",
      cell: ({ row }) => <Badge variant={row.original.isCompleted ? "default" : "secondary"}>{row.original.isCompleted ? "Completed" : "Pending"}</Badge>,
    },
  ];

  if (loading) return <div><PageHeader title="Manager Overview" breadcrumbs={[{ label: "Competency" }, { label: "Manager Overview" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Manager Overview" description="Competency review status for your team"
        breadcrumbs={[{ label: "Competency" }, { label: "Manager Overview" }]}
      />

      <div className="flex gap-3">
        <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
          <SelectTrigger className="w-[220px]"><SelectValue placeholder="Select Period" /></SelectTrigger>
          <SelectContent>{periods.map((p) => <SelectItem key={p.reviewPeriodId} value={String(p.reviewPeriodId)}>{p.name}</SelectItem>)}</SelectContent>
        </Select>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Total Staff</CardTitle></CardHeader><CardContent><p className="text-2xl font-bold">{reviews.length}</p></CardContent></Card>
        <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Completed</CardTitle></CardHeader><CardContent><p className="text-2xl font-bold text-green-600">{completed}</p></CardContent></Card>
        <Card><CardHeader className="pb-2"><CardTitle className="text-sm text-muted-foreground">Pending</CardTitle></CardHeader><CardContent><p className="text-2xl font-bold text-orange-600">{pending}</p></CardContent></Card>
      </div>

      {reviews.length > 0 ? (
        <DataTable columns={columns} data={reviews} searchKey="employeeName" searchPlaceholder="Search employees..." />
      ) : (
        <EmptyState icon={Users} title="No Reviews" description="No team competency reviews found for the selected period." />
      )}
    </div>
  );
}
