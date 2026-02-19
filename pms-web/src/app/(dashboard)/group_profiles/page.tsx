"use client";

import { useEffect, useState } from "react";
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
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { Search } from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import {
  getCompetencyReviewPeriods,
  getGroupCompetencyReviewProfiles,
} from "@/lib/api/competency";
import {
  getDepartments,
  getDivisionsByDepartment,
  getOfficesByDivision,
} from "@/lib/api/organogram";
import type { CompetencyReviewPeriod } from "@/types/competency";
import type { GroupedCompetencyReviewProfile } from "@/types/dashboard";
import type { Department, Division, Office } from "@/types/organogram";

export default function GroupProfilesPage() {
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [profile, setProfile] =
    useState<GroupedCompetencyReviewProfile | null>(null);

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
      const res = await getGroupCompetencyReviewProfiles({
        reviewPeriodId: reviewPeriodId ? Number(reviewPeriodId) : undefined,
        departmentId: departmentId ? Number(departmentId) : undefined,
        divisionId: divisionId ? Number(divisionId) : undefined,
        officeId: officeId ? Number(officeId) : undefined,
      });
      if (res?.data) setProfile(res.data);
    } catch {
      /* */
    } finally {
      setLoading(false);
    }
  };

  const summaryChartData = (profile?.categoryCompetencyStats || []).map(
    (s) => ({
      name: s.categoryName,
      actual: s.actual,
      expected: s.expected,
    })
  );

  if (initialLoading) {
    return (
      <div>
        <PageHeader
          title="Group Strengths and Development Needs"
          breadcrumbs={[{ label: "Group Profiles" }]}
        />
        <PageSkeleton />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Group Strengths and Development Needs"
        breadcrumbs={[{ label: "Group Profiles" }]}
      />

      <div className="grid gap-6 md:grid-cols-12">
        {/* Main Content */}
        <div className="md:col-span-9 space-y-6">
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

          {profile && (
            <>
              {/* Summary Section */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Summary</CardTitle>
                </CardHeader>
                <CardContent>
                  {summaryChartData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <BarChart data={summaryChartData}>
                        <CartesianGrid
                          strokeDasharray="3 3"
                          className="stroke-border"
                        />
                        <XAxis
                          dataKey="name"
                          tick={{
                            fill: "hsl(var(--muted-foreground))",
                            fontSize: 11,
                          }}
                        />
                        <YAxis
                          tick={{
                            fill: "hsl(var(--muted-foreground))",
                            fontSize: 12,
                          }}
                        />
                        <Tooltip
                          contentStyle={{
                            backgroundColor: "hsl(var(--popover))",
                            border: "1px solid hsl(var(--border))",
                            borderRadius: "6px",
                          }}
                        />
                        <Legend />
                        <Bar
                          dataKey="actual"
                          name="Current Proficiency"
                          fill="hsl(var(--chart-1))"
                          radius={[4, 4, 0, 0]}
                        />
                        <Bar
                          dataKey="expected"
                          name="Expected Proficiency"
                          fill="hsl(var(--chart-3))"
                          radius={[4, 4, 0, 0]}
                        />
                      </BarChart>
                    </ResponsiveContainer>
                  ) : (
                    <p className="text-sm text-muted-foreground">
                      No summary data available.
                    </p>
                  )}
                </CardContent>
              </Card>

              {/* Detail Analysis Section */}
              {(profile.categoryCompetencyDetailStats?.length ?? 0) > 0 && (
                <Card>
                  <CardHeader>
                    <CardTitle className="text-base">
                      Competency Group Detail Analysis
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <Accordion
                      type="multiple"
                      defaultValue={profile.categoryCompetencyDetailStats.map(
                        (_, idx) => `cat-${idx}`
                      )}
                    >
                      {profile.categoryCompetencyDetailStats.map(
                        (cat, idx) => (
                          <AccordionItem key={idx} value={`cat-${idx}`}>
                            <AccordionTrigger className="font-semibold">
                              {cat.categoryName}
                            </AccordionTrigger>
                            <AccordionContent>
                              <div className="grid gap-4 md:grid-cols-12">
                                {/* Rating stats table */}
                                <div className="md:col-span-7">
                                  <table className="w-full text-sm">
                                    <thead>
                                      <tr>
                                        <th className="text-left pb-2">
                                          Rating
                                        </th>
                                        <th className="text-center pb-2">
                                          Number
                                        </th>
                                        <th className="text-center pb-2">
                                          Percentage
                                        </th>
                                      </tr>
                                    </thead>
                                    <tbody>
                                      {[...cat.competencyRatingStat]
                                        .sort(
                                          (a, b) =>
                                            b.ratingOrder - a.ratingOrder
                                        )
                                        .map((r, ri) => (
                                          <tr
                                            key={ri}
                                            className="border-t border-border"
                                          >
                                            <td className="py-1.5">
                                              {r.ratingOrder}-{r.ratingName}
                                            </td>
                                            <td className="text-center py-1.5">
                                              {r.numberOfStaff}
                                            </td>
                                            <td className="text-center py-1.5">
                                              {r.staffPercentage}%
                                            </td>
                                          </tr>
                                        ))}
                                    </tbody>
                                  </table>
                                </div>

                                {/* Divider */}
                                <div className="md:col-span-1 hidden md:flex items-center justify-center">
                                  <div className="h-full w-px bg-border" />
                                </div>

                                {/* Category metrics */}
                                <div className="md:col-span-4">
                                  <table className="w-full text-sm">
                                    <tbody>
                                      <tr className="border-t border-border">
                                        <td className="py-1.5 text-muted-foreground">
                                          Average Proficiency
                                        </td>
                                        <td className="py-1.5 font-medium text-right">
                                          {cat.averageRating.toFixed(2)}
                                        </td>
                                      </tr>
                                      <tr className="border-t border-border">
                                        <td className="py-1.5 text-muted-foreground">
                                          Highest Proficiency
                                        </td>
                                        <td className="py-1.5 font-medium text-right">
                                          {cat.highestRating.toFixed(2)}
                                        </td>
                                      </tr>
                                      <tr className="border-t border-border">
                                        <td className="py-1.5 text-muted-foreground">
                                          Lowest Proficiency
                                        </td>
                                        <td className="py-1.5 font-medium text-right">
                                          {cat.lowestRating.toFixed(2)}
                                        </td>
                                      </tr>
                                      <tr className="border-t border-border">
                                        <td className="py-1.5 text-muted-foreground">
                                          Common Proficiency
                                        </td>
                                        <td className="py-1.5 font-medium text-right">
                                          {cat.mostCommonRating.toFixed(2)}
                                        </td>
                                      </tr>
                                    </tbody>
                                  </table>
                                </div>
                              </div>

                              {/* Per-category bar chart */}
                              {cat.groupCompetencyRatings?.length > 0 && (
                                <div className="mt-4">
                                  <ResponsiveContainer
                                    width="100%"
                                    height={250}
                                  >
                                    <BarChart
                                      data={cat.groupCompetencyRatings.map(
                                        (r) => ({
                                          name: r.label,
                                          actual: r.actual,
                                          expected: r.expected,
                                        })
                                      )}
                                    >
                                      <CartesianGrid
                                        strokeDasharray="3 3"
                                        className="stroke-border"
                                      />
                                      <XAxis
                                        dataKey="name"
                                        tick={{ fontSize: 10 }}
                                      />
                                      <YAxis />
                                      <Tooltip
                                        contentStyle={{
                                          backgroundColor:
                                            "hsl(var(--popover))",
                                          border:
                                            "1px solid hsl(var(--border))",
                                          borderRadius: "6px",
                                        }}
                                      />
                                      <Legend />
                                      <Bar
                                        dataKey="actual"
                                        name="Current"
                                        fill="hsl(var(--chart-1))"
                                        radius={[4, 4, 0, 0]}
                                      />
                                      <Bar
                                        dataKey="expected"
                                        name="Expected"
                                        fill="hsl(var(--chart-3))"
                                        radius={[4, 4, 0, 0]}
                                      />
                                    </BarChart>
                                  </ResponsiveContainer>
                                </div>
                              )}
                            </AccordionContent>
                          </AccordionItem>
                        )
                      )}
                    </Accordion>
                  </CardContent>
                </Card>
              )}
            </>
          )}
        </div>

        {/* Legend Sidebar */}
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
                <span className="text-sm">Current Proficiency</span>
              </div>
              <div className="flex items-center gap-2">
                <Badge variant="secondary">3</Badge>
                <span className="text-sm">Expected Proficiency</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
