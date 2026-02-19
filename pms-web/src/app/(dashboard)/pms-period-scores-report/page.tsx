"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { format } from "date-fns";
import { Download, Loader2, Search } from "lucide-react";
import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getActiveReviewPeriod, getPeriodScores } from "@/lib/api/dashboard";
import { getReviewPeriods } from "@/lib/api/review-periods";
import { getDepartments, getDivisionsByDepartment, getOfficesByDivision } from "@/lib/api/organogram";
import type { PerformanceReviewPeriod } from "@/types/performance";
import type { PeriodScoreData } from "@/types/dashboard";
import type { Department, Division, Office } from "@/types/organogram";
import { Roles } from "@/stores/auth-store";

const ADMIN_ROLES: string[] = [Roles.Admin, Roles.SuperAdmin, Roles.HrReportAdmin];

export default function PmsPeriodScoresReportPage() {
  const { data: session, status } = useSession();
  const router = useRouter();
  const userRoles = session?.user?.roles ?? [];

  const [loading, setLoading] = useState(true);
  const [searching, setSearching] = useState(false);

  // Filter state
  const [reviewPeriodId, setReviewPeriodId] = useState("");
  const [departmentId, setDepartmentId] = useState<number>(0);
  const [divisionId, setDivisionId] = useState<number>(0);
  const [officeId, setOfficeId] = useState<number>(0);

  // Dropdown data
  const [periods, setPeriods] = useState<PerformanceReviewPeriod[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [divisions, setDivisions] = useState<Division[]>([]);
  const [offices, setOffices] = useState<Office[]>([]);

  // Score data
  const [scores, setScores] = useState<PeriodScoreData[]>([]);
  const [filteredScores, setFilteredScores] = useState<PeriodScoreData[]>([]);

  useEffect(() => {
    if (status === "authenticated" && !userRoles.some((r) => ADMIN_ROLES.includes(r))) {
      router.push("/access-denied");
    }
  }, [status, userRoles, router]);

  // Initial load
  useEffect(() => {
    const init = async () => {
      setLoading(true);
      try {
        const [activeRes, periodsRes, deptsRes] = await Promise.all([
          getActiveReviewPeriod(),
          getReviewPeriods(),
          getDepartments(),
        ]);
        if (periodsRes?.data) setPeriods(Array.isArray(periodsRes.data) ? periodsRes.data : []);
        if (deptsRes?.data) setDepartments(Array.isArray(deptsRes.data) ? deptsRes.data : []);

        const activePeriod = activeRes?.data;
        if (activePeriod?.periodId) {
          setReviewPeriodId(activePeriod.periodId);
          const scoresRes = await getPeriodScores(activePeriod.periodId);
          if (scoresRes?.data) {
            const all = Array.isArray(scoresRes.data) ? scoresRes.data : [];
            setScores(all);
            setFilteredScores(all);
          }
        }
      } catch { /* */ } finally { setLoading(false); }
    };
    init();
  }, []);

  // Cascading filters
  const handleDeptChange = async (deptId: number) => {
    setDepartmentId(deptId);
    setDivisionId(0);
    setOfficeId(0);
    setOffices([]);
    if (deptId > 0) {
      const res = await getDivisionsByDepartment(deptId);
      setDivisions(res?.data ? (Array.isArray(res.data) ? res.data : []) : []);
    } else {
      setDivisions([]);
    }
  };

  const handleDivChange = async (divId: number) => {
    setDivisionId(divId);
    setOfficeId(0);
    if (divId > 0) {
      const res = await getOfficesByDivision(divId);
      setOffices(res?.data ? (Array.isArray(res.data) ? res.data : []) : []);
    } else {
      setOffices([]);
    }
  };

  // Search / filter
  const handleSearch = async () => {
    if (!reviewPeriodId) { toast.error("Please select a review period."); return; }
    setSearching(true);
    try {
      const res = await getPeriodScores(reviewPeriodId);
      if (res?.data) {
        const all = Array.isArray(res.data) ? res.data : [];
        setScores(all);
        applyFilters(all, departmentId, divisionId, officeId);
      }
    } catch { toast.error("Failed to fetch scores."); } finally { setSearching(false); }
  };

  const applyFilters = (data: PeriodScoreData[], deptId: number, divId: number, offId: number) => {
    let filtered = data;
    if (deptId > 0) filtered = filtered.filter((s) => s.departmentId === deptId);
    if (divId > 0) filtered = filtered.filter((s) => s.divisionId === divId);
    if (offId > 0) filtered = filtered.filter((s) => s.officeId === offId);
    setFilteredScores(filtered);
  };

  // Re-apply client-side filters when department/division/office changes
  useEffect(() => {
    applyFilters(scores, departmentId, divisionId, officeId);
  }, [departmentId, divisionId, officeId, scores]);

  // CSV Export
  const exportToExcel = () => {
    const headers = [
      "StaffId", "FullName", "Department", "Division", "Office", "Staff Grade",
      "Score Percentage", "Final Score", "FinalGradeName", "MinNoOfObjectives",
      "MaxNoOfObjectives", "HRDDeductedPoints", "Review Period", "Year",
    ];
    const csvContent = [
      headers.join(","),
      ...filteredScores.map((row) => [
        row.staffId, `"${row.staffFullName ?? ""}"`, `"${row.departmentName ?? ""}"`, `"${row.divisionName ?? ""}"`,
        `"${row.officeName ?? ""}"`, row.staffGrade ?? "", (row.scorePercentage ?? 0).toFixed(2), (row.finalScore ?? 0).toFixed(4),
        row.finalGradeName ?? "", row.minNoOfObjectives ?? "", row.maxNoOfObjectives ?? "", row.hrdDeductedPoints ?? 0,
        `"${row.reviewPeriod ?? ""}"`, row.year ?? "",
      ].join(","))
    ].join("\n");

    const blob = new Blob([csvContent], { type: "text/csv" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "Staff_Performance_Score.csv";
    a.click();
    URL.revokeObjectURL(url);
  };

  const numBadge = (value: number | undefined) => (
    <Badge variant="outline" className="font-bold">{value != null ? Math.round(value * 100) / 100 : "-"}</Badge>
  );

  const columns: ColumnDef<PeriodScoreData>[] = [
    { accessorKey: "staffId", header: "Staff ID" },
    { accessorKey: "staffFullName", header: "Name" },
    { accessorKey: "departmentName", header: "Department" },
    { accessorKey: "divisionName", header: "Division" },
    { accessorKey: "officeName", header: "Office" },
    { accessorKey: "staffGrade", header: "Grade" },
    { accessorKey: "scorePercentage", header: "Score %", cell: ({ row }) => numBadge(row.original.scorePercentage) },
    { accessorKey: "finalScore", header: "Final Score", cell: ({ row }) => numBadge(row.original.finalScore) },
    { accessorKey: "maxPoint", header: "Max Point", cell: ({ row }) => numBadge(row.original.maxPoint) },
    { accessorKey: "finalGradeName", header: "Grade", cell: ({ row }) => <Badge variant="outline" className="font-bold">{row.original.finalGradeName || "-"}</Badge> },
    { accessorKey: "minNoOfObjectives", header: "Min Obj", cell: ({ row }) => numBadge(row.original.minNoOfObjectives) },
    { accessorKey: "maxNoOfObjectives", header: "Max Obj", cell: ({ row }) => numBadge(row.original.maxNoOfObjectives) },
    { accessorKey: "hrdDeductedPoints", header: "Deducted", cell: ({ row }) => numBadge(row.original.hrdDeductedPoints) },
    { accessorKey: "reviewPeriod", header: "Review Period" },
    { accessorKey: "year", header: "Year" },
    { accessorKey: "startDate", header: "Start Date", cell: ({ row }) => row.original.startDate ? format(new Date(row.original.startDate), "dd MMM yyyy") : "-" },
    { accessorKey: "endDate", header: "End Date", cell: ({ row }) => row.original.endDate ? format(new Date(row.original.endDate), "dd MMM yyyy") : "-" },
    { accessorKey: "strategyName", header: "Strategy" },
    { accessorKey: "locationId", header: "Location ID" },
    { accessorKey: "isUnderPerforming", header: "Under Performing", cell: ({ row }) => row.original.isUnderPerforming ? "Yes" : "No" },
  ];

  if (loading) return <div><PageHeader title="Performance Scores Report" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Performance Scores Report" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Performance Scores Report" breadcrumbs={[{ label: "Home", href: "/" }, { label: "Performance Scores Report" }]} />

      {/* Filter Row */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5 items-end">
        <div className="space-y-2">
          <Label>Review Period</Label>
          <Select value={reviewPeriodId} onValueChange={setReviewPeriodId}>
            <SelectTrigger><SelectValue placeholder="Select period" /></SelectTrigger>
            <SelectContent>{periods.map((p) => <SelectItem key={p.periodId} value={p.periodId}>{p.name}</SelectItem>)}</SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label>Department</Label>
          <Select value={departmentId > 0 ? String(departmentId) : "all"} onValueChange={(v) => handleDeptChange(v === "all" ? 0 : Number(v))}>
            <SelectTrigger><SelectValue placeholder="All departments" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Departments</SelectItem>
              {departments.map((d) => <SelectItem key={d.departmentId} value={String(d.departmentId)}>{d.departmentName}</SelectItem>)}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label>Division</Label>
          <Select value={divisionId > 0 ? String(divisionId) : "all"} onValueChange={(v) => handleDivChange(v === "all" ? 0 : Number(v))} disabled={departmentId === 0}>
            <SelectTrigger><SelectValue placeholder="All divisions" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Divisions</SelectItem>
              {divisions.map((d) => <SelectItem key={d.divisionId} value={String(d.divisionId)}>{d.divisionName}</SelectItem>)}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label>Office</Label>
          <Select value={officeId > 0 ? String(officeId) : "all"} onValueChange={(v) => setOfficeId(v === "all" ? 0 : Number(v))} disabled={divisionId === 0}>
            <SelectTrigger><SelectValue placeholder="All offices" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Offices</SelectItem>
              {offices.map((o) => <SelectItem key={o.officeId} value={String(o.officeId)}>{o.officeName}</SelectItem>)}
            </SelectContent>
          </Select>
        </div>
        <Button onClick={handleSearch} disabled={searching}>
          {searching ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Search className="mr-2 h-4 w-4" />}Search
        </Button>
      </div>

      {/* Export + Table */}
      <div className="flex justify-end">
        <Button variant="outline" size="sm" onClick={exportToExcel} disabled={filteredScores.length === 0}>
          <Download className="mr-2 h-4 w-4" />Export to Excel
        </Button>
      </div>

      <div className="overflow-x-auto">
        <DataTable columns={columns} data={filteredScores} searchKey="staffId" searchPlaceholder="Search by staff ID or name..." />
      </div>
    </div>
  );
}
