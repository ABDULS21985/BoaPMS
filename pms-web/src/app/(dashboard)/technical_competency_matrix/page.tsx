"use client";

import { useEffect, useState, useMemo } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { DataTable } from "@/components/shared/data-table";
import { Search } from "lucide-react";
import {
  getCompetencyReviewPeriods,
  getTechnicalCompetencyMatrixReviewProfiles,
  getJobRoles,
} from "@/lib/api/competency";
import type { CompetencyReviewPeriod, JobRole } from "@/types/competency";
import type {
  CompetencyMatrixReviewOverview,
  CompetencyMatrixReviewProfile,
} from "@/types/dashboard";

export default function TechnicalCompetencyMatrixPage() {
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [matrixData, setMatrixData] =
    useState<CompetencyMatrixReviewOverview | null>(null);

  // Filter data
  const [reviewPeriods, setReviewPeriods] = useState<CompetencyReviewPeriod[]>(
    []
  );
  const [jobRoles, setJobRoles] = useState<JobRole[]>([]);

  // Filter selections
  const [reviewPeriodId, setReviewPeriodId] = useState("");
  const [jobRoleId, setJobRoleId] = useState("");

  useEffect(() => {
    Promise.all([getCompetencyReviewPeriods(), getJobRoles()])
      .then(([rpRes, jrRes]) => {
        if (rpRes?.data)
          setReviewPeriods(Array.isArray(rpRes.data) ? rpRes.data : []);
        if (jrRes?.data)
          setJobRoles(Array.isArray(jrRes.data) ? jrRes.data : []);
      })
      .finally(() => setInitialLoading(false));
  }, []);

  const handleSearch = async () => {
    if (!reviewPeriodId || !jobRoleId) return;
    setLoading(true);
    try {
      const res = await getTechnicalCompetencyMatrixReviewProfiles(
        Number(reviewPeriodId),
        Number(jobRoleId)
      );
      if (res?.data) setMatrixData(res.data);
    } catch {
      /* */
    } finally {
      setLoading(false);
    }
  };

  const baseColumns: ColumnDef<CompetencyMatrixReviewProfile>[] = [
    { accessorKey: "employeeId", header: "Employee ID" },
    { accessorKey: "employeeName", header: "Name" },
    { accessorKey: "position", header: "Role" },
    { accessorKey: "grade", header: "Grade" },
    { accessorKey: "officeName", header: "Office" },
    { accessorKey: "divisionName", header: "Division" },
    { accessorKey: "departmentName", header: "Department" },
    {
      accessorKey: "noOfCompetencies",
      header: "Total",
      cell: ({ row }) => (
        <Badge variant="outline">{row.original.noOfCompetencies}</Badge>
      ),
    },
    {
      accessorKey: "noOfCompetent",
      header: "Competent",
      cell: ({ row }) => (
        <Badge variant="outline" className="border-green-500">
          {row.original.noOfCompetent}
        </Badge>
      ),
    },
    {
      accessorKey: "gapCount",
      header: "Gaps",
      cell: ({ row }) => (
        <Badge variant="outline" className="border-amber-500">
          {row.original.gapCount}
        </Badge>
      ),
    },
    { accessorKey: "overallAverage", header: "Avg" },
  ];

  const competencyColumns: ColumnDef<CompetencyMatrixReviewProfile>[] = useMemo(
    () =>
      (matrixData?.competencyNames || []).map((name, idx) => ({
        id: `competency_${idx}`,
        header: name,
        cell: ({
          row,
        }: {
          row: { original: CompetencyMatrixReviewProfile };
        }) => {
          const detail = row.original.competencyMatrixDetails?.[idx];
          if (!detail) return "—";
          const diff = detail.expectedRatingValue - detail.averageScore;
          let badgeClass = "";
          if (diff > 0)
            badgeClass = "bg-amber-100 text-amber-800 border-amber-300";
          else if (diff < 0)
            badgeClass = "bg-blue-100 text-blue-800 border-blue-300";
          else badgeClass = "bg-green-100 text-green-800 border-green-300";
          return (
            <Badge variant="outline" className={badgeClass}>
              <b>{detail.averageScore}</b>/<b>{detail.expectedRatingValue}</b>
            </Badge>
          );
        },
      })),
    [matrixData?.competencyNames]
  );

  const allColumns = useMemo(
    () => [...baseColumns, ...competencyColumns],
    [competencyColumns]
  );

  if (initialLoading) {
    return (
      <div>
        <PageHeader
          title="Technical Competency Matrix"
          breadcrumbs={[{ label: "Technical Competency Matrix" }]}
        />
        <PageSkeleton />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Technical Competency Matrix"
        breadcrumbs={[{ label: "Technical Competency Matrix" }]}
      />

      {/* Filters */}
      <div className="flex flex-wrap items-end gap-3">
        <Select value={reviewPeriodId} onValueChange={setReviewPeriodId}>
          <SelectTrigger className="w-52">
            <SelectValue placeholder="Review Period" />
          </SelectTrigger>
          <SelectContent>
            {reviewPeriods.map((rp) => (
              <SelectItem
                key={rp.reviewPeriodId}
                value={String(rp.reviewPeriodId)}
              >
                {rp.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Select value={jobRoleId} onValueChange={setJobRoleId}>
          <SelectTrigger className="w-52">
            <SelectValue placeholder="Job Role" />
          </SelectTrigger>
          <SelectContent>
            {jobRoles.map((jr) => (
              <SelectItem key={jr.jobRoleId} value={String(jr.jobRoleId)}>
                {jr.jobRoleName}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Button
          onClick={handleSearch}
          disabled={loading || !reviewPeriodId || !jobRoleId}
        >
          <Search className="mr-2 h-4 w-4" />
          {loading ? "Searching..." : "Search"}
        </Button>
      </div>

      {/* Main content + Legend */}
      <div className="grid gap-6 md:grid-cols-12">
        <div className="md:col-span-9">
          <div className="overflow-x-auto">
            <DataTable
              columns={allColumns}
              data={matrixData?.competencyMatrixReviewProfiles || []}
              searchKey="employeeName"
              searchPlaceholder="Search by name..."
            />
          </div>
        </div>

        <div className="md:col-span-3">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Legend</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center gap-2">
                <Badge
                  variant="outline"
                  className="bg-green-100 text-green-800 border-green-300"
                >
                  3/3
                </Badge>
                <span className="text-sm">Matches</span>
              </div>
              <div className="flex items-center gap-2">
                <Badge
                  variant="outline"
                  className="bg-amber-100 text-amber-800 border-amber-300"
                >
                  2/3
                </Badge>
                <span className="text-sm">Gap</span>
              </div>
              <div className="flex items-center gap-2">
                <Badge
                  variant="outline"
                  className="bg-blue-100 text-blue-800 border-blue-300"
                >
                  4/3
                </Badge>
                <span className="text-sm">Exceeds</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
