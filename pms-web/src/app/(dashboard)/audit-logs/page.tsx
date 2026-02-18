"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Eye, Loader2, FileText } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormDialog } from "@/components/shared/form-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getAuditLogs } from "@/lib/api/pms-engine";
import type { AuditLog } from "@/types/staff";

const auditEventLabels: Record<number, string> = {
  1: "Added",
  2: "Deleted",
  3: "Modified",
};

export default function AuditLogsPage() {
  const [items, setItems] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<AuditLog | null>(null);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const res = await getAuditLogs();
        if (res?.auditLogs) setItems(Array.isArray(res.auditLogs) ? res.auditLogs : []);
      } catch { /* */ } finally { setLoading(false); }
    })();
  }, []);

  const formatDate = (dateStr: string) => {
    if (!dateStr) return "N/A";
    try {
      return new Date(dateStr).toLocaleString("en-GB", {
        day: "2-digit", month: "short", year: "numeric",
        hour: "2-digit", minute: "2-digit", second: "2-digit",
      });
    } catch { return dateStr; }
  };

  const columns: ColumnDef<AuditLog>[] = [
    { accessorKey: "userName", header: "User", cell: ({ row }) => row.original.userName || "SYSTEM" },
    { accessorKey: "auditEventDateUTC", header: "Event Date", cell: ({ row }) => formatDate(row.original.auditEventDateUTC) },
    { accessorKey: "auditEventType", header: "Action", cell: ({ row }) => <Badge variant="secondary">{auditEventLabels[row.original.auditEventType] ?? `Type ${row.original.auditEventType}`}</Badge> },
    { accessorKey: "tableName", header: "Table", cell: ({ row }) => row.original.tableName || "N/A" },
    { accessorKey: "fieldName", header: "Field", cell: ({ row }) => row.original.fieldName || "N/A" },
    {
      id: "actions", header: "", cell: ({ row }) => (
        <Button size="sm" variant="outline" onClick={() => { setSelectedItem(row.original); setDetailsOpen(true); }}><Eye className="h-3.5 w-3.5" /></Button>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Audit Logs" breadcrumbs={[{ label: "Reports" }, { label: "Audit Logs" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Audit Logs" description="View system audit trail" breadcrumbs={[{ label: "Reports" }, { label: "Audit Logs" }]} />

      {items.length > 0 ? (
        <DataTable columns={columns} data={items} searchKey="userName" searchPlaceholder="Search by user..." />
      ) : (
        <EmptyState icon={FileText} title="No Audit Logs" description="No audit log entries found." />
      )}

      {/* Details Dialog */}
      <FormDialog open={detailsOpen} onOpenChange={setDetailsOpen} title="Audit Log Details" className="sm:max-w-lg">
        {selectedItem ? (
          <div className="space-y-3 text-sm">
            <div className="grid grid-cols-2 gap-3">
              <div><span className="text-muted-foreground">User:</span> <span className="font-medium">{selectedItem.userName || "SYSTEM"}</span></div>
              <div><span className="text-muted-foreground">Event Date:</span> <span className="font-medium">{formatDate(selectedItem.auditEventDateUTC)}</span></div>
              <div><span className="text-muted-foreground">Action:</span> <span className="font-medium">{auditEventLabels[selectedItem.auditEventType] ?? `Type ${selectedItem.auditEventType}`}</span></div>
              <div><span className="text-muted-foreground">Table:</span> <span className="font-medium">{selectedItem.tableName || "N/A"}</span></div>
              <div><span className="text-muted-foreground">Record ID:</span> <span className="font-medium font-mono text-xs">{selectedItem.recordId || "N/A"}</span></div>
              <div><span className="text-muted-foreground">Field:</span> <span className="font-medium">{selectedItem.fieldName || "N/A"}</span></div>
            </div>
            <div>
              <span className="text-muted-foreground text-xs">Original Value:</span>
              <p className="text-sm mt-1 bg-muted/50 rounded-md p-2 break-all">{selectedItem.originalValue || "N/A"}</p>
            </div>
            <div>
              <span className="text-muted-foreground text-xs">New Value:</span>
              <p className="text-sm mt-1 bg-muted/50 rounded-md p-2 break-all">{selectedItem.newValue || "No new value provided"}</p>
            </div>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">No details available.</p>
        )}
      </FormDialog>
    </div>
  );
}
