"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import type { ColumnDef } from "@tanstack/react-table";
import { Eye, Users } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getGroupCompetencyReviewProfiles, getCompetencyReviewPeriods } from "@/lib/api/competency";
import { getStaffDetails } from "@/lib/api/staff";
import type { CompetencyReviewProfile, CompetencyReviewPeriod } from "@/types/competency";

export default function SubordinatesPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const [profiles, setProfiles] = useState<CompetencyReviewProfile[]>([]);
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
    getGroupCompetencyReviewProfiles({ officeId, reviewPeriodId: Number(selectedPeriod) })
      .then((res) => {
        if (res?.data) setProfiles(Array.isArray(res.data) ? res.data : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [officeId, selectedPeriod]);

  // Group by employee
  const employeeMap = new Map<string, CompetencyReviewProfile>();
  profiles.forEach((p) => {
    if (!employeeMap.has(p.employeeNumber)) employeeMap.set(p.employeeNumber, p);
  });
  const employees = Array.from(employeeMap.values());

  const columns: ColumnDef<CompetencyReviewProfile>[] = [
    { accessorKey: "employeeFullName", header: "Employee", cell: ({ row }) => row.original.employeeFullName ?? row.original.employeeNumber },
    { accessorKey: "jobRoleName", header: "Job Role" },
    { accessorKey: "gradeName", header: "Grade" },
    { accessorKey: "departmentName", header: "Department" },
    { accessorKey: "averageScore", header: "Avg Score", cell: ({ row }) => <Badge variant="outline">{row.original.averageScore.toFixed(1)}</Badge> },
    {
      id: "actions", header: "",
      cell: ({ row }) => (
        <Button size="sm" variant="outline" onClick={() => router.push(`/competency/profiles/${row.original.employeeNumber}`)}>
          <Eye className="mr-1 h-3.5 w-3.5" />View Profile
        </Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Subordinates" breadcrumbs={[{ label: "Competency" }, { label: "Subordinates" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Subordinates" description="View subordinate competency profiles"
        breadcrumbs={[{ label: "Competency" }, { label: "Subordinates" }]}
      />

      <div className="flex gap-3">
        <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
          <SelectTrigger className="w-[220px]"><SelectValue placeholder="Select Period" /></SelectTrigger>
          <SelectContent>{periods.map((p) => <SelectItem key={p.reviewPeriodId} value={String(p.reviewPeriodId)}>{p.name}</SelectItem>)}</SelectContent>
        </Select>
      </div>

      {employees.length > 0 ? (
        <DataTable columns={columns} data={employees} searchKey="employeeFullName" searchPlaceholder="Search employees..." />
      ) : (
        <EmptyState icon={Users} title="No Subordinates" description="No subordinate profiles found for the selected period." />
      )}
    </div>
  );
}
