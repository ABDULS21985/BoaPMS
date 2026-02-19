"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { format } from "date-fns";
import { Info, Loader2, Search } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormDialog } from "@/components/shared/form-dialog";
import { EmptyState } from "@/components/shared/empty-state";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getStaffReviewPeriods, getPeriodScoreDetails } from "@/lib/api/dashboard";
import type { PerformanceReviewPeriod } from "@/types/performance";
import type { PeriodScoreData } from "@/types/dashboard";
import { Roles } from "@/stores/auth-store";

const ADMIN_ROLES: string[] = [Roles.Admin, Roles.SuperAdmin, Roles.HrReportAdmin, Roles.HrAdmin];

export default function StaffReviewPeriodsPage() {
  const { data: session, status } = useSession();
  const router = useRouter();
  const userRoles = session?.user?.roles ?? [];

  const [staffIdSearch, setStaffIdSearch] = useState("");
  const [periods, setPeriods] = useState<PerformanceReviewPeriod[]>([]);
  const [searching, setSearching] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);
  const [detailOpen, setDetailOpen] = useState(false);
  const [scoreDetails, setScoreDetails] = useState<PeriodScoreData | null>(null);
  const [loadingDetails, setLoadingDetails] = useState(false);

  useEffect(() => {
    if (status === "authenticated" && !userRoles.some((r) => ADMIN_ROLES.includes(r))) {
      router.push("/access-denied");
    }
  }, [status, userRoles, router]);

  const handleSearch = async () => {
    if (!staffIdSearch.trim()) { toast.error("Please enter a Staff ID."); return; }
    setSearching(true);
    setHasSearched(true);
    try {
      const res = await getStaffReviewPeriods(staffIdSearch.trim());
      if (res?.data) setPeriods(Array.isArray(res.data) ? res.data : []);
      else setPeriods([]);
    } catch { toast.error("Failed to fetch review periods."); } finally { setSearching(false); }
  };

  const handleDetails = async (period: PerformanceReviewPeriod) => {
    setLoadingDetails(true);
    setDetailOpen(true);
    try {
      const res = await getPeriodScoreDetails(period.periodId, staffIdSearch.trim());
      if (res?.data) setScoreDetails(res.data);
      else setScoreDetails(null);
    } catch { toast.error("Failed to load score details."); } finally { setLoadingDetails(false); }
  };

  const columns: ColumnDef<PerformanceReviewPeriod>[] = [
    { id: "index", header: "#", cell: ({ row }) => row.index + 1 },
    { accessorKey: "name", header: "Name" },
    { accessorKey: "year", header: "Year" },
    { accessorKey: "startDate", header: "Start Date", cell: ({ row }) => row.original.startDate ? format(new Date(row.original.startDate), "dd MMM yyyy") : "-" },
    { accessorKey: "endDate", header: "End Date", cell: ({ row }) => row.original.endDate ? format(new Date(row.original.endDate), "dd MMM yyyy") : "-" },
    {
      id: "actions", header: "Action", cell: ({ row }) => (
        <Button size="sm" variant="ghost" onClick={() => handleDetails(row.original)} title="View Details"><Info className="h-3.5 w-3.5" /></Button>
      ),
    },
  ];

  if (status === "loading") return <div><PageHeader title="Staff Review Periods" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Staff Review Periods" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Staff Review Periods" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Staff Review Periods" }]} />

      <div className="flex items-end gap-3">
        <div className="space-y-2">
          <Label>Staff ID</Label>
          <Input placeholder="Enter Staff ID..." value={staffIdSearch} onChange={(e) => setStaffIdSearch(e.target.value)} onKeyDown={(e) => e.key === "Enter" && handleSearch()} className="w-64" />
        </div>
        <Button onClick={handleSearch} disabled={searching}>
          {searching ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Search className="mr-2 h-4 w-4" />}Search
        </Button>
      </div>

      {hasSearched && (
        periods.length > 0
          ? <DataTable columns={columns} data={periods} searchKey="name" searchPlaceholder="Filter by period name..." />
          : <EmptyState title="No Review Periods" description="No review periods found for the given Staff ID." />
      )}

      <FormDialog open={detailOpen} onOpenChange={setDetailOpen} title="Period Score Details" className="sm:max-w-lg">
        {loadingDetails ? (
          <div className="flex items-center justify-center py-8"><Loader2 className="h-6 w-6 animate-spin text-muted-foreground" /></div>
        ) : scoreDetails ? (
          <div className="space-y-3 text-sm">
            <div><span className="font-medium text-muted-foreground">Review Period:</span> <span>{scoreDetails.reviewPeriod} — {scoreDetails.startDate ? format(new Date(scoreDetails.startDate), "dd MMM yyyy") : ""} to {scoreDetails.endDate ? format(new Date(scoreDetails.endDate), "dd MMM yyyy") : ""}</span></div>
            <div><span className="font-medium text-muted-foreground">Name:</span> <span>{scoreDetails.staffFullName} ({scoreDetails.staffId})</span></div>
            <div><span className="font-medium text-muted-foreground">Grade:</span> <span>{scoreDetails.staffGrade || "-"}</span></div>
            <div><span className="font-medium text-muted-foreground">Office:</span> <span>{[scoreDetails.officeName, scoreDetails.divisionName, scoreDetails.departmentName].filter(Boolean).join(", ") || "-"}</span></div>
            <div><span className="font-medium text-muted-foreground">Performance Score:</span> <span>{scoreDetails.finalScore} / {scoreDetails.maxPoint}</span></div>
            <div><span className="font-medium text-muted-foreground">Deducted Points:</span> <span>{scoreDetails.hrdDeductedPoints}</span></div>
            <div><span className="font-medium text-muted-foreground">Score Percentage:</span> <span>{scoreDetails.scorePercentage?.toFixed(2)}%</span></div>
            <div><span className="font-medium text-muted-foreground">Performance Grade:</span> <span>{scoreDetails.finalGradeName || "-"}</span></div>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No score data available for this period.</p>
        )}
        <div className="flex justify-end pt-4">
          <Button variant="outline" onClick={() => setDetailOpen(false)}>Close</Button>
        </div>
      </FormDialog>
    </div>
  );
}
