"use client";

import { useState, useEffect } from "react";
import { type ColumnDef } from "@tanstack/react-table";
import { toast } from "sonner";
import { UserCog } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { FormDialog } from "@/components/shared/form-dialog";
import { getProjects, updateProject } from "@/lib/api/pms-engine";
import { getEmployeeDetail } from "@/lib/api/dashboard";
import type { Project } from "@/types/performance";
import type { EmployeeErpDetails } from "@/types/dashboard";

export default function ProjectsReportPage() {
  const [loading, setLoading] = useState(true);
  const [projects, setProjects] = useState<Project[]>([]);
  const [reassignOpen, setReassignOpen] = useState(false);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [newManagerId, setNewManagerId] = useState("");
  const [verifiedEmployee, setVerifiedEmployee] = useState<EmployeeErpDetails | null>(null);
  const [verifying, setVerifying] = useState(false);
  const [saving, setSaving] = useState(false);

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getProjects();
      setProjects(res?.data ?? []);
    } catch {
      toast.error("Failed to load projects.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { loadData(); }, []);

  const handleVerify = async () => {
    if (!newManagerId.trim()) {
      toast.error("Enter a staff ID.");
      return;
    }
    setVerifying(true);
    try {
      const res = await getEmployeeDetail(newManagerId.trim());
      if (res?.data?.employeeNumber) {
        setVerifiedEmployee(res.data);
        toast.success(`Found: ${res.data.firstName} ${res.data.lastName}`);
      } else {
        setVerifiedEmployee(null);
        toast.error("Staff not found.");
      }
    } catch {
      toast.error("Verification failed.");
    } finally {
      setVerifying(false);
    }
  };

  const handleReassign = async () => {
    if (!selectedProject || !verifiedEmployee) return;
    setSaving(true);
    try {
      const res = await updateProject({
        projectId: selectedProject.projectId,
        projectManager: verifiedEmployee.employeeNumber,
      });
      if (res?.isSuccess) {
        toast.success("Project manager re-assigned successfully.");
        setReassignOpen(false);
        loadData();
      } else {
        toast.error(res?.message || "Failed to re-assign.");
      }
    } catch {
      toast.error("Failed to re-assign project manager.");
    } finally {
      setSaving(false);
    }
  };

  const openReassign = (project: Project) => {
    setSelectedProject(project);
    setNewManagerId("");
    setVerifiedEmployee(null);
    setReassignOpen(true);
  };

  const columns: ColumnDef<Project>[] = [
    {
      id: "index",
      header: "#",
      cell: ({ row }) => row.index + 1,
    },
    { accessorKey: "name", header: "Name" },
    { accessorKey: "description", header: "Description" },
    { accessorKey: "deliverables", header: "Objective/KPI" },
    {
      accessorKey: "departmentId",
      header: "Department",
      cell: ({ row }) => row.original.departmentId ?? "-",
    },
    { accessorKey: "projectManager", header: "Project Manager" },
    {
      accessorKey: "recordStatus",
      header: "Status",
      cell: ({ row }) => <StatusBadge status={row.original.recordStatus ?? 0} />,
    },
    {
      accessorKey: "startDate",
      header: "Start Date",
      cell: ({ row }) => row.original.startDate ? new Date(row.original.startDate).toLocaleDateString() : "-",
    },
    {
      accessorKey: "endDate",
      header: "End Date",
      cell: ({ row }) => row.original.endDate ? new Date(row.original.endDate).toLocaleDateString() : "-",
    },
    {
      id: "actions",
      header: "Action",
      cell: ({ row }) => (
        <Button variant="ghost" size="sm" onClick={() => openReassign(row.original)}>
          <UserCog className="mr-1 h-4 w-4" /> Re-assign PM
        </Button>
      ),
    },
  ];

  if (loading) return <PageSkeleton />;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Projects Report"
        breadcrumbs={[{ label: "Projects Report" }]}
      />

      <DataTable
        columns={columns}
        data={projects}
        searchKey="name"
        searchPlaceholder="Search by project name..."
      />

      <FormDialog
        open={reassignOpen}
        onOpenChange={setReassignOpen}
        title="Re-assign Project Manager"
        description={selectedProject ? `Project: ${selectedProject.name}` : ""}
      >
        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label>Current Project Manager</Label>
            <Input value={selectedProject?.projectManager || "N/A"} readOnly className="bg-muted" />
          </div>
          <div className="space-y-1.5">
            <Label>New Project Manager Staff ID</Label>
            <div className="flex gap-2">
              <Input
                value={newManagerId}
                onChange={(e) => { setNewManagerId(e.target.value); setVerifiedEmployee(null); }}
                placeholder="Enter staff ID"
              />
              <Button variant="outline" onClick={handleVerify} disabled={verifying}>
                {verifying ? "Verifying..." : "Verify"}
              </Button>
            </div>
          </div>
          {verifiedEmployee && (
            <div className="rounded-md border p-3 text-sm">
              <p className="font-medium">{verifiedEmployee.firstName} {verifiedEmployee.lastName}</p>
              <p className="text-muted-foreground">{verifiedEmployee.jobTitle}</p>
              <p className="text-muted-foreground">{verifiedEmployee.departmentName}</p>
            </div>
          )}
          <div className="flex justify-end gap-2 pt-2">
            <Button variant="outline" onClick={() => setReassignOpen(false)}>Cancel</Button>
            <Button onClick={handleReassign} disabled={saving || !verifiedEmployee}>
              {saving ? "Re-assigning..." : "Re-assign"}
            </Button>
          </div>
        </div>
      </FormDialog>
    </div>
  );
}
