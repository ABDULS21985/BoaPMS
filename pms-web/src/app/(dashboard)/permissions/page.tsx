"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Lock } from "lucide-react";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getAllPermissions } from "@/lib/api/staff";
import type { Permission } from "@/types/staff";

export default function PermissionsPage() {
  const [items, setItems] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const res = await getAllPermissions();
        if (res?.data) setItems(Array.isArray(res.data) ? res.data : []);
      } catch { /* */ } finally { setLoading(false); }
    })();
  }, []);

  const columns: ColumnDef<Permission>[] = [
    { accessorKey: "permissionName", header: "Permission Name" },
    { accessorKey: "description", header: "Description", cell: ({ row }) => row.original.description ?? "—" },
  ];

  if (loading) return <div><PageHeader title="Permissions" breadcrumbs={[{ label: "Settings" }, { label: "Permissions" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Permissions" description="View all system permissions" breadcrumbs={[{ label: "Settings" }, { label: "Permissions" }]} />
      {items.length > 0 ? (
        <DataTable columns={columns} data={items} searchKey="permissionName" searchPlaceholder="Search permissions..." />
      ) : (
        <EmptyState icon={Lock} title="No Permissions" description="No permissions found in the system." />
      )}
    </div>
  );
}
