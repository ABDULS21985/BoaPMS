"use client";

import { useState, useEffect } from "react";
import { type ColumnDef } from "@tanstack/react-table";
import { toast } from "sonner";
import { Eye } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { FormDialog } from "@/components/shared/form-dialog";
import { getGrievanceReport, getGrievanceTypes } from "@/lib/api/grievance";
import type { Grievance } from "@/types/performance";

export default function GrievanceReportPage() {
  const [loading, setLoading] = useState(true);
  const [grievances, setGrievances] = useState<Grievance[]>([]);
  const [grievanceTypes, setGrievanceTypes] = useState<{ id: number; name: string }[]>([]);
  const [detailOpen, setDetailOpen] = useState(false);
  const [selectedGrievance, setSelectedGrievance] = useState<Grievance | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const [gRes, typesRes] = await Promise.all([getGrievanceReport(), getGrievanceTypes()]);
        setGrievances(gRes?.data ?? []);
        setGrievanceTypes(typesRes?.data ?? []);
      } catch {
        toast.error("Failed to load grievance data.");
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const getTypeName = (typeId: number) => {
    return grievanceTypes.find((t) => t.id === typeId)?.name ?? `Type ${typeId}`;
  };

  const columns: ColumnDef<Grievance>[] = [
    {
      id: "index",
      header: "#",
      cell: ({ row }) => row.index + 1,
    },
    {
      accessorKey: "grievanceType",
      header: "Grievance Type",
      cell: ({ row }) => getTypeName(row.original.grievanceType),
    },
    { accessorKey: "subject", header: "Subject" },
    { accessorKey: "complainantStaffId", header: "Complainant" },
    { accessorKey: "respondentStaffId", header: "Respondent" },
    {
      accessorKey: "currentResolutionLevel",
      header: "Resolution Level",
      cell: ({ row }) => `Level ${row.original.currentResolutionLevel}`,
    },
    {
      accessorKey: "recordStatus",
      header: "Status",
      cell: ({ row }) => <StatusBadge status={row.original.recordStatus ?? 0} />,
    },
    {
      accessorKey: "dateCreated",
      header: "Date Created",
      cell: ({ row }) => new Date(row.original.dateCreated).toLocaleDateString(),
    },
    {
      id: "actions",
      header: "Action",
      cell: ({ row }) => (
        <Button variant="ghost" size="sm" onClick={() => { setSelectedGrievance(row.original); setDetailOpen(true); }}>
          <Eye className="mr-1 h-4 w-4" /> Details
        </Button>
      ),
    },
  ];

  if (loading) return <PageSkeleton />;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Grievance Report"
        breadcrumbs={[{ label: "Grievance Report" }]}
      />

      <DataTable
        columns={columns}
        data={grievances}
        searchKey="subject"
        searchPlaceholder="Search by subject..."
      />

      <FormDialog
        open={detailOpen}
        onOpenChange={setDetailOpen}
        title="Grievance Details"
        className="max-w-lg"
      >
        {selectedGrievance && (
          <div className="space-y-3 text-sm">
            <DetailField label="Grievance Type" value={getTypeName(selectedGrievance.grievanceType)} />
            <DetailField label="Subject" value={selectedGrievance.subject ?? selectedGrievance.subjectId} />
            <DetailField label="Description" value={selectedGrievance.description} />
            <DetailField label="Complainant Staff ID" value={selectedGrievance.complainantStaffId} />
            <DetailField label="Respondent Staff ID" value={selectedGrievance.respondentStaffId} />
            <DetailField label="Current Resolution Level" value={`Level ${selectedGrievance.currentResolutionLevel}`} />
            <DetailField label="Status" value={<StatusBadge status={selectedGrievance.recordStatus ?? 0} />} />
            <DetailField label="Date Created" value={new Date(selectedGrievance.dateCreated).toLocaleDateString()} />
            {selectedGrievance.respondentComment && (
              <DetailField label="Respondent Comment" value={selectedGrievance.respondentComment} />
            )}
          </div>
        )}
      </FormDialog>
    </div>
  );
}

function DetailField({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div>
      <Label className="text-xs text-muted-foreground">{label}</Label>
      <div className="mt-0.5">{value || "-"}</div>
    </div>
  );
}
