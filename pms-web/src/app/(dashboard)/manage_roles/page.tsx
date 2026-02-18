"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Trash2, Shield, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { FormDialog } from "@/components/shared/form-dialog";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import {
  getRoles,
  addRole,
  deleteRole,
  getAllRolesWithPermissions,
  addPermissionToRole,
  removePermissionFromRole,
} from "@/lib/api/staff";
import type { Role, Permission } from "@/types/staff";

interface RoleDisplay {
  roleName: string;
}

export default function ManageRolesPage() {
  const [items, setItems] = useState<RoleDisplay[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [newRoleName, setNewRoleName] = useState("");
  const [deleteConfirm, setDeleteConfirm] = useState<RoleDisplay | null>(null);

  // Permission assignment state
  const [permDialogOpen, setPermDialogOpen] = useState(false);
  const [selectedRole, setSelectedRole] = useState<RoleDisplay | null>(null);
  const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
  const [rolePermissions, setRolePermissions] = useState<Permission[]>([]);
  const [permsLoading, setPermsLoading] = useState(false);
  const [selectedPerm, setSelectedPerm] = useState("");
  const [assigningPerm, setAssigningPerm] = useState(false);
  const [removePermConfirm, setRemovePermConfirm] = useState<Permission | null>(null);

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getRoles();
      if (res?.data) {
        const all = Array.isArray(res.data) ? res.data : [];
        setItems(all.map((r) => ({ roleName: typeof r === "string" ? r : (r as Role).roleName ?? (r as unknown as string) })));
      }
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const handleAddRole = async () => {
    if (!newRoleName.trim()) { toast.error("Enter a role name."); return; }
    setSaving(true);
    try {
      const res = await addRole({ role_name: newRoleName.trim() });
      if (res?.isSuccess) { toast.success("Role created."); setOpen(false); setNewRoleName(""); loadData(); }
      else toast.error(res?.message || "Failed to create role.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const handleDeleteRole = async () => {
    if (!deleteConfirm) return;
    try {
      const res = await deleteRole(deleteConfirm.roleName);
      if (res?.isSuccess) { toast.success("Role deleted."); loadData(); }
      else toast.error(res?.message || "Failed to delete role.");
    } catch { toast.error("An error occurred."); }
  };

  const openPermDialog = async (role: RoleDisplay) => {
    setSelectedRole(role);
    setSelectedPerm("");
    setPermDialogOpen(true);
    setPermsLoading(true);
    try {
      const res = await getAllRolesWithPermissions(role.roleName);
      if (res?.data) {
        setAllPermissions(res.data.allPermissions ?? []);
        setRolePermissions(res.data.rolesAndPermissions?.permissions ?? []);
      }
    } catch { setAllPermissions([]); setRolePermissions([]); } finally { setPermsLoading(false); }
  };

  const handleAssignPerm = async () => {
    if (!selectedRole || !selectedPerm) return;
    setAssigningPerm(true);
    try {
      const res = await addPermissionToRole({ role_id: selectedRole.roleName, permission_id: selectedPerm });
      if (res?.isSuccess) {
        toast.success("Permission assigned.");
        const updated = await getAllRolesWithPermissions(selectedRole.roleName);
        if (updated?.data) {
          setAllPermissions(updated.data.allPermissions ?? []);
          setRolePermissions(updated.data.rolesAndPermissions?.permissions ?? []);
        }
        setSelectedPerm("");
      } else toast.error(res?.message || "Failed to assign permission.");
    } catch { toast.error("An error occurred."); } finally { setAssigningPerm(false); }
  };

  const handleRemovePerm = async () => {
    if (!selectedRole || !removePermConfirm) return;
    try {
      const res = await removePermissionFromRole(selectedRole.roleName, removePermConfirm.permissionId?.toString() ?? removePermConfirm.permissionName);
      if (res?.isSuccess) {
        toast.success("Permission removed.");
        const updated = await getAllRolesWithPermissions(selectedRole.roleName);
        if (updated?.data) {
          setAllPermissions(updated.data.allPermissions ?? []);
          setRolePermissions(updated.data.rolesAndPermissions?.permissions ?? []);
        }
      } else toast.error(res?.message || "Failed to remove permission.");
    } catch { toast.error("An error occurred."); }
  };

  const columns: ColumnDef<RoleDisplay>[] = [
    { accessorKey: "roleName", header: "Role Name" },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => openPermDialog(row.original)} title="Manage Permissions"><Shield className="h-3.5 w-3.5" /></Button>
          <Button size="sm" variant="ghost" onClick={() => setDeleteConfirm(row.original)} title="Delete"><Trash2 className="h-3.5 w-3.5 text-destructive" /></Button>
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Manage Roles" breadcrumbs={[{ label: "Settings" }, { label: "Manage Roles" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Manage Roles" description="Create and manage system roles" breadcrumbs={[{ label: "Settings" }, { label: "Manage Roles" }]} actions={<Button size="sm" onClick={() => setOpen(true)}><Plus className="mr-2 h-4 w-4" />Add Role</Button>} />
      <DataTable columns={columns} data={items} searchKey="roleName" searchPlaceholder="Search roles..." />

      {/* Add Role Sheet */}
      <FormSheet open={open} onOpenChange={setOpen} title="Add Role">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Role Name *</Label><Input value={newRoleName} onChange={(e) => setNewRoleName(e.target.value)} placeholder="Enter role name" /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleAddRole} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save</Button>
          </div>
        </div>
      </FormSheet>

      {/* Delete Confirmation */}
      <ConfirmationDialog open={!!deleteConfirm} onOpenChange={() => setDeleteConfirm(null)} title="Delete Role" description={`Delete role "${deleteConfirm?.roleName}"? This cannot be undone.`} confirmLabel="Delete" variant="destructive" onConfirm={handleDeleteRole} />

      {/* Permission Assignment Dialog */}
      <FormDialog open={permDialogOpen} onOpenChange={setPermDialogOpen} title={`Permissions — ${selectedRole?.roleName}`} className="sm:max-w-lg">
        <div className="space-y-4">
          <div className="flex gap-2">
            <Select value={selectedPerm} onValueChange={setSelectedPerm}>
              <SelectTrigger className="flex-1"><SelectValue placeholder="Select permission" /></SelectTrigger>
              <SelectContent>{allPermissions.filter((p) => !rolePermissions.some((rp) => rp.permissionName === p.permissionName)).map((p) => <SelectItem key={p.permissionName} value={p.permissionName}>{p.permissionName}</SelectItem>)}</SelectContent>
            </Select>
            <Button onClick={handleAssignPerm} disabled={!selectedPerm || assigningPerm}>
              {assigningPerm ? <Loader2 className="h-4 w-4 animate-spin" /> : "Add"}
            </Button>
          </div>

          {permsLoading ? (
            <div className="flex items-center justify-center py-4"><Loader2 className="h-5 w-5 animate-spin" /></div>
          ) : rolePermissions.length > 0 ? (
            <div className="space-y-2">
              <Label className="text-muted-foreground">Assigned Permissions</Label>
              <div className="space-y-1 max-h-[300px] overflow-y-auto">
                {rolePermissions.map((p) => (
                  <div key={p.permissionName} className="flex items-center justify-between rounded-md border px-3 py-2 text-sm">
                    <span>{p.permissionName}</span>
                    <Button size="sm" variant="ghost" onClick={() => setRemovePermConfirm(p)}><Trash2 className="h-3.5 w-3.5 text-destructive" /></Button>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">No permissions assigned.</p>
          )}
        </div>
      </FormDialog>

      {/* Remove Permission Confirmation */}
      <ConfirmationDialog open={!!removePermConfirm} onOpenChange={() => setRemovePermConfirm(null)} title="Remove Permission" description={`Remove permission "${removePermConfirm?.permissionName}" from this role?`} confirmLabel="Remove" variant="destructive" onConfirm={handleRemovePerm} />
    </div>
  );
}
