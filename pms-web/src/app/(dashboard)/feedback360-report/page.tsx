"use client";

import { useState, useEffect, useCallback } from "react";
import { type ColumnDef } from "@tanstack/react-table";
import { toast } from "sonner";
import { Eye } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { FormDialog } from "@/components/shared/form-dialog";
import { getReviewPeriods } from "@/lib/api/review-periods";
import { getDepartments, getDivisionsByDepartment, getOfficesByDivision } from "@/lib/api/organogram";
import { getAllCompetencyReviewFeedbacks, getCompetencyReviewFeedbackDetails } from "@/lib/api/pms-engine";
import { getPeriodScores } from "@/lib/api/dashboard";
import type { PerformanceReviewPeriod, CompetencyReviewFeedback } from "@/types/performance";
import type { CompetencyReviewFeedbackDetails, CompetencyReviewerRatingSummary } from "@/types/performance";
import type { Department, Division, Office } from "@/types/organogram";
import type { PeriodScoreData } from "@/types/dashboard";
import { formatPercent } from "@/lib/scorecard-helpers";

interface FeedbackRow {
  staffId: string;
  staffName: string;
  departmentName?: string;
  divisionName?: string;
  officeName?: string;
  feedbackStatus: string;
  scorePercent: number;
  feedbackId?: string;
}

export default function Feedback360ReportPage() {
  const [loading, setLoading] = useState(true);
  const [searching, setSearching] = useState(false);
  const [reviewPeriods, setReviewPeriods] = useState<PerformanceReviewPeriod[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [divisions, setDivisions] = useState<Division[]>([]);
  const [offices, setOffices] = useState<Office[]>([]);

  const [selectedPeriod, setSelectedPeriod] = useState("");
  const [selectedDept, setSelectedDept] = useState("");
  const [selectedDiv, setSelectedDiv] = useState("");
  const [selectedOffice, setSelectedOffice] = useState("");

  const [rows, setRows] = useState<FeedbackRow[]>([]);
  const [detailOpen, setDetailOpen] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailData, setDetailData] = useState<CompetencyReviewFeedbackDetails | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const [rpRes, deptRes] = await Promise.all([getReviewPeriods(), getDepartments()]);
        setReviewPeriods(rpRes?.data ?? []);
        setDepartments(deptRes?.data ?? []);
      } catch {
        toast.error("Failed to load filter data.");
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const handleDeptChange = useCallback(async (deptId: string) => {
    setSelectedDept(deptId);
    setSelectedDiv("");
    setSelectedOffice("");
    setDivisions([]);
    setOffices([]);
    if (deptId) {
      try {
        const res = await getDivisionsByDepartment(Number(deptId));
        setDivisions(res?.data ?? []);
      } catch { /* ignore */ }
    }
  }, []);

  const handleDivChange = useCallback(async (divId: string) => {
    setSelectedDiv(divId);
    setSelectedOffice("");
    setOffices([]);
    if (divId) {
      try {
        const res = await getOfficesByDivision(Number(divId));
        setOffices(res?.data ?? []);
      } catch { /* ignore */ }
    }
  }, []);

  const handleSearch = async () => {
    if (!selectedPeriod) {
      toast.error("Please select a review period.");
      return;
    }
    setSearching(true);
    try {
      const scoresRes = await getPeriodScores(selectedPeriod);
      let scores: PeriodScoreData[] = scoresRes?.data ?? [];

      if (selectedDept) scores = scores.filter((s) => s.departmentId === Number(selectedDept));
      if (selectedDiv) scores = scores.filter((s) => s.divisionId === Number(selectedDiv));
      if (selectedOffice) scores = scores.filter((s) => s.officeId === Number(selectedOffice));

      const feedbackRows: FeedbackRow[] = scores.map((s) => ({
        staffId: s.staffId,
        staffName: s.staffFullName,
        departmentName: s.departmentName,
        divisionName: s.divisionName,
        officeName: s.officeName,
        feedbackStatus: "Loaded",
        scorePercent: s.scorePercentage,
      }));

      setRows(feedbackRows);
      if (feedbackRows.length === 0) toast.info("No records found for the selected criteria.");
    } catch {
      toast.error("Failed to search feedback data.");
    } finally {
      setSearching(false);
    }
  };

  const handleViewDetails = async (row: FeedbackRow) => {
    setDetailOpen(true);
    setDetailLoading(true);
    setDetailData(null);
    try {
      const fbRes = await getAllCompetencyReviewFeedbacks(row.staffId);
      const feedbacks: CompetencyReviewFeedback[] = fbRes?.data ?? [];
      if (feedbacks.length > 0) {
        const detailRes = await getCompetencyReviewFeedbackDetails(feedbacks[0].competencyReviewFeedbackId);
        setDetailData(detailRes?.data ?? null);
      }
    } catch {
      toast.error("Failed to load feedback details.");
    } finally {
      setDetailLoading(false);
    }
  };

  const columns: ColumnDef<FeedbackRow>[] = [
    { accessorKey: "staffId", header: "Staff ID" },
    { accessorKey: "staffName", header: "Staff Name" },
    { accessorKey: "departmentName", header: "Department" },
    { accessorKey: "divisionName", header: "Division" },
    { accessorKey: "officeName", header: "Office" },
    { accessorKey: "feedbackStatus", header: "Feedback Status" },
    {
      accessorKey: "scorePercent",
      header: "Score %",
      cell: ({ row }) => formatPercent(row.original.scorePercent),
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => (
        <Button variant="ghost" size="sm" onClick={() => handleViewDetails(row.original)}>
          <Eye className="mr-1 h-4 w-4" /> Details
        </Button>
      ),
    },
  ];

  if (loading) return <PageSkeleton />;

  return (
    <div className="space-y-6">
      <PageHeader
        title="360 Feedback Report"
        breadcrumbs={[{ label: "360 Feedback Report" }]}
      />

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5">
        <div className="space-y-1.5">
          <Label>Review Period</Label>
          <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
            <SelectTrigger><SelectValue placeholder="Select period" /></SelectTrigger>
            <SelectContent>
              {reviewPeriods.map((rp) => (
                <SelectItem key={rp.periodId} value={rp.periodId}>{rp.name}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-1.5">
          <Label>Department</Label>
          <Select value={selectedDept} onValueChange={handleDeptChange}>
            <SelectTrigger><SelectValue placeholder="All departments" /></SelectTrigger>
            <SelectContent>
              {departments.map((d) => (
                <SelectItem key={d.departmentId} value={String(d.departmentId)}>{d.departmentName}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-1.5">
          <Label>Division</Label>
          <Select value={selectedDiv} onValueChange={handleDivChange} disabled={!selectedDept}>
            <SelectTrigger><SelectValue placeholder="All divisions" /></SelectTrigger>
            <SelectContent>
              {divisions.map((d) => (
                <SelectItem key={d.divisionId} value={String(d.divisionId)}>{d.divisionName}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-1.5">
          <Label>Office</Label>
          <Select value={selectedOffice} onValueChange={setSelectedOffice} disabled={!selectedDiv}>
            <SelectTrigger><SelectValue placeholder="All offices" /></SelectTrigger>
            <SelectContent>
              {offices.map((o) => (
                <SelectItem key={o.officeId} value={String(o.officeId)}>{o.officeName}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex items-end">
          <Button onClick={handleSearch} disabled={searching} className="w-full">
            {searching ? "Searching..." : "Search"}
          </Button>
        </div>
      </div>

      {rows.length > 0 ? (
        <DataTable columns={columns} data={rows} searchKey="staffName" searchPlaceholder="Search by staff name..." />
      ) : (
        !searching && <EmptyState title="No feedback data" description="Select filters and click Search to view 360 feedback report." />
      )}

      <FormDialog open={detailOpen} onOpenChange={setDetailOpen} title="360 Feedback Details" className="max-w-2xl">
        {detailLoading ? (
          <div className="py-8 text-center text-sm text-muted-foreground">Loading details...</div>
        ) : detailData ? (
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-2 text-sm">
              <div><span className="font-medium">Staff:</span> {detailData.staffName}</div>
              <div><span className="font-medium">Department:</span> {detailData.departmentName}</div>
              <div><span className="font-medium">Division:</span> {detailData.divisionName}</div>
              <div><span className="font-medium">Office:</span> {detailData.officeName}</div>
              <div><span className="font-medium">Final Score:</span> {detailData.finalScore} / {detailData.maxPoints}</div>
              <div><span className="font-medium">Score %:</span> {formatPercent(detailData.finalScorePercentage)}</div>
            </div>
            {detailData.ratings && detailData.ratings.length > 0 && (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Competency</TableHead>
                    <TableHead>Avg Rating</TableHead>
                    <TableHead>Total Reviewers</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {detailData.ratings.map((r: CompetencyReviewerRatingSummary) => (
                    <TableRow key={r.pmsCompetencyId}>
                      <TableCell>{r.pmsCompetencyName}</TableCell>
                      <TableCell>{r.averageRating.toFixed(2)}</TableCell>
                      <TableCell>{r.totalReviewers}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </div>
        ) : (
          <div className="py-8 text-center text-sm text-muted-foreground">No feedback details available.</div>
        )}
      </FormDialog>
    </div>
  );
}
