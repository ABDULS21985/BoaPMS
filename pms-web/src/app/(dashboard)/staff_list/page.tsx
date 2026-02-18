"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Users } from "lucide-react";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getStaffList } from "@/lib/api/staff";
import type { Staff } from "@/types/staff";
import { Badge } from "@/components/ui/badge";

export default function StaffListPage() {
  const [items, setItems] = useState<Staff[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const res = await getStaffList();
        if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      } catch { /* */ } finally { setLoading(false); }
    })();
  }, []);

  const columns: ColumnDef<Staff>[] = [
    { accessorKey: "employeeNumber", header: "Staff ID" },
    { accessorKey: "firstName", header: "First Name" },
    { accessorKey: "lastName", header: "Last Name" },
    { accessorKey: "email", header: "Email" },
    { accessorKey: "jobName", header: "Grade", cell: ({ row }) => row.original.gradeName ?? "N/A" },
    { accessorKey: "officeName", header: "Office", cell: ({ row }) => row.original.officeName ?? "N/A" },
    { accessorKey: "divisionName", header: "Division", cell: ({ row }) => row.original.divisionName ?? "N/A" },
    { accessorKey: "departmentName", header: "Department", cell: ({ row }) => row.original.departmentName ?? "N/A" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Disabled"}</Badge> },
  ];

  if (loading) return <div><PageHeader title="Staff List" breadcrumbs={[{ label: "Organogram" }, { label: "Staff List" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Staff List" description="View all staff members from the system" breadcrumbs={[{ label: "Organogram" }, { label: "Staff List" }]} />
      {items.length > 0 ? (
        <DataTable columns={columns} data={items} searchKey="firstName" searchPlaceholder="Search by name..." />
      ) : (
        <EmptyState icon={Users} title="No Staff" description="No staff records found." />
      )}
    </div>
  );
}
