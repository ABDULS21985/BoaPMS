"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Eye, FolderKanban } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getProjectsAssigned } from "@/lib/api/pms-engine";
import { getStaffDetails } from "@/lib/api/staff";
import type { Project } from "@/types/performance";
import type { Staff } from "@/types/staff";
import type { ColumnDef } from "@tanstack/react-table";

export default function StaffProjectsPage() {
  const { staffId } = useParams<{ staffId: string }>();
  const router = useRouter();
  const [projects, setProjects] = useState<Project[]>([]);
  const [staff, setStaff] = useState<Staff | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!staffId) return;
    setLoading(true);
    Promise.all([getStaffDetails(staffId), getProjectsAssigned(staffId)])
      .then(([staffRes, projRes]) => {
        if (staffRes?.data) setStaff(staffRes.data);
        if (projRes?.data) setProjects(Array.isArray(projRes.data) ? projRes.data : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [staffId]);

  const staffName = staff ? `${staff.firstName} ${staff.lastName}` : staffId;

  const columns: ColumnDef<Project>[] = [
    { accessorKey: "name", header: "Project Name", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "projectManager", header: "Manager" },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="outline" onClick={() => router.push(`/projects/${row.original.projectId}`)}><Eye className="mr-1 h-3.5 w-3.5" />View</Button> },
  ];

  if (loading) return <div><PageHeader title="Staff Projects" breadcrumbs={[{ label: "Projects" }, { label: "Staff" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={`${staffName}'s Projects`} breadcrumbs={[{ label: "Projects", href: "/projects" }, { label: staffName }]} />
      {projects.length > 0 ? (
        <DataTable columns={columns} data={projects} searchKey="name" searchPlaceholder="Search projects..." />
      ) : (
        <EmptyState icon={FolderKanban} title="No Projects" description="This staff member is not assigned to any projects." />
      )}
    </div>
  );
}
