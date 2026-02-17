"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Briefcase } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { getProjectDetails, getProjectWorkProducts } from "@/lib/api/pms-engine";
import type { Project, WorkProduct } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

export default function ProjectWorkProductsPage() {
  const { projectId } = useParams<{ projectId: string }>();
  const [project, setProject] = useState<Project | null>(null);
  const [workProducts, setWorkProducts] = useState<WorkProduct[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!projectId) return;
    setLoading(true);
    Promise.all([getProjectDetails(projectId), getProjectWorkProducts(projectId)])
      .then(([projRes, wpRes]) => {
        if (projRes?.data) setProject(projRes.data);
        if (wpRes?.data) setWorkProducts(Array.isArray(wpRes.data) ? wpRes.data : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [projectId]);

  const columns: ColumnDef<WorkProduct>[] = [
    { accessorKey: "name", header: "Work Product", cell: ({ row }) => <span className="font-medium">{row.original.name}</span> },
    { accessorKey: "deliverables", header: "Deliverables", cell: ({ row }) => <span className="line-clamp-1">{row.original.deliverables}</span> },
    { accessorKey: "staffId", header: "Assigned To" },
    { accessorKey: "startDate", header: "Start", cell: ({ row }) => row.original.startDate?.split("T")[0] },
    { accessorKey: "endDate", header: "End", cell: ({ row }) => row.original.endDate?.split("T")[0] },
    { accessorKey: "maxPoint", header: "Max Pts" },
    { accessorKey: "finalScore", header: "Score", cell: ({ row }) => <Badge variant={row.original.finalScore > 0 ? "default" : "secondary"}>{row.original.finalScore}</Badge> },
    { accessorKey: "recordStatus", header: "Status", cell: ({ row }) => row.original.recordStatus != null ? <StatusBadge status={row.original.recordStatus} /> : "—" },
  ];

  if (loading) return <div><PageHeader title="Project Work Products" breadcrumbs={[{ label: "Projects" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={`Work Products: ${project?.name ?? ""}`} description={project?.description} breadcrumbs={[{ label: "Projects", href: "/projects" }, { label: project?.name ?? "Project", href: `/projects/${projectId}` }, { label: "Work Products" }]} />
      {workProducts.length > 0 ? (
        <DataTable columns={columns} data={workProducts} searchKey="name" searchPlaceholder="Search work products..." />
      ) : (
        <EmptyState icon={Briefcase} title="No Work Products" description="No work products assigned to this project." />
      )}
    </div>
  );
}
