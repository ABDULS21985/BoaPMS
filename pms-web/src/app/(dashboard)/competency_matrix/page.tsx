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
  getCompetencyMatrixReviewProfiles,
} from "@/lib/api/competency";
import {
  getDepartments,
  getDivisionsByDepartment,
  getOfficesByDivision,
} from "@/lib/api/organogram";
import type { CompetencyReviewPeriod } from "@/types/competency";
import type {
  CompetencyMatrixReviewOverview,
  CompetencyMatrixReviewProfile,
} from "@/types/dashboard";
import type { Department, Division, Office } from "@/types/organogram";

export default function CompetencyMatrixPage() {
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [matrixData, setMatrixData] =
    useState<CompetencyMatrixReviewOverview | null>(null);

  // Filter data
  const [reviewPeriods, setReviewPeriods] = useState<CompetencyReviewPeriod[]>(
    []
  );
  const [departments, setDepartments] = useState<Department[]>([]);
  const [divisions, setDivisions] = useState<Division[]>([]);
  const [offices, setOffices] = useState<Office[]>([]);

  // Filter selections
  const [reviewPeriodId, setReviewPeriodId] = useState("");
  const [departmentId, setDepartmentId] = useState("");
  const [divisionId, setDivisionId] = useState("");
  const [officeId, setOfficeId] = useState("");

  useEffect(() => {
    Promise.all([getCompetencyReviewPeriods(), getDepartments()])
      .then(([rpRes, deptRes]) => {
        if (rpRes?.data)
          setReviewPeriods(Array.isArray(rpRes.data) ? rpRes.data : []);
        if (deptRes?.data)
          setDepartments(Array.isArray(deptRes.data) ? deptRes.data : []);
      })
      .finally(() => setInitialLoading(false));
  }, []);

  const onDeptChange = async (val: string) => {
    setDepartmentId(val);
    setDivisionId("");
    setOfficeId("");
    setDivisions([]);
    setOffices([]);
    if (val) {
      try {
        const r = await getDivisionsByDepartment(Number(val));
        if (r?.data) setDivisions(Array.isArray(r.data) ? r.data : []);
      } catch {
        /* */
      }
    }
  };

  const onDivChange = async (val: string) => {
    setDivisionId(val);
    setOfficeId("");
    setOffices([]);
    if (val) {
      try {
        const r = await getOfficesByDivision(Number(val));
        if (r?.data) setOffices(Array.isArray(r.data) ? r.data : []);
      } catch {
        /* */
      }
    }
  };

  const handleSearch = async () => {
    setLoading(true);
    try {
      const res = await getCompetencyMatrixReviewProfiles({
        reviewPeriodId: reviewPeriodId ? Number(reviewPeriodId) : undefined,
        departmentId: departmentId ? Number(departmentId) : undefined,
        divisionId: divisionId ? Number(divisionId) : undefined,
        officeId: officeId ? Number(officeId) : undefined,
      });
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
          title="Behavioral Competency Matrix"
          breadcrumbs={[{ label: "Behavioral Competency Matrix" }]}
        />
        <PageSkeleton />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Behavioral Competency Matrix"
        breadcrumbs={[{ label: "Behavioral Competency Matrix" }]}
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

        <Select value={departmentId} onValueChange={onDeptChange}>
          <SelectTrigger className="w-48">
            <SelectValue placeholder="All Departments" />
          </SelectTrigger>
          <SelectContent>
            {departments.map((d) => (
              <SelectItem
                key={d.departmentId}
                value={String(d.departmentId)}
              >
                {d.departmentName}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        {divisions.length > 0 && (
          <Select value={divisionId} onValueChange={onDivChange}>
            <SelectTrigger className="w-48">
              <SelectValue placeholder="All Divisions" />
            </SelectTrigger>
            <SelectContent>
              {divisions.map((d) => (
                <SelectItem
                  key={d.divisionId}
                  value={String(d.divisionId)}
                >
                  {d.divisionName}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )}

        {offices.length > 0 && (
          <Select value={officeId} onValueChange={setOfficeId}>
            <SelectTrigger className="w-48">
              <SelectValue placeholder="All Offices" />
            </SelectTrigger>
            <SelectContent>
              {offices.map((o) => (
                <SelectItem key={o.officeId} value={String(o.officeId)}>
                  {o.officeName}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )}

        <Button onClick={handleSearch} disabled={loading}>
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
