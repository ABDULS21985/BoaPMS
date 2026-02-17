"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { Eye, FolderKanban } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getProjectsAssigned } from "@/lib/api/pms-engine";
import type { Project } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function MyProjectsPage() {
  const { data: session } = useSession();
  const router = useRouter();
  const staffId = session?.user?.id ?? "";
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!staffId) return;
    setLoading(true);
    getProjectsAssigned(staffId)
      .then((res) => { if (res?.data) setProjects(Array.isArray(res.data) ? res.data : []); })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [staffId]);

  const columns: ColumnDef<Project>[] = [
    { accessorKey: "name", header: "Project Name", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "projectManager", header: "Manager" },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
    { id: "actions", header: "Actions", cell: ({ row }) => <Button size="sm" variant="outline" onClick={() => router.push(`/projects/${row.original.projectId}`)}><Eye className="mr-1 h-3.5 w-3.5" />View</Button> },
  ];

  if (loading) return <div><PageHeader title="My Projects" breadcrumbs={[{ label: "Projects" }, { label: "My Projects" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="My Projects" description="Projects you are assigned to" breadcrumbs={[{ label: "Projects", href: "/projects" }, { label: "My Projects" }]} />
      {projects.length > 0 ? (
        <DataTable columns={columns} data={projects} searchKey="name" searchPlaceholder="Search projects..." />
      ) : (
        <EmptyState icon={FolderKanban} title="No Projects" description="You are not assigned to any projects." />
      )}
    </div>
  );
}
