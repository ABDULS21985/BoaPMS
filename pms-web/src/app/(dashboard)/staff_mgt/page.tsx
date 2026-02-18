"use client";

import { useEffect, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { Plus, Pencil, Shield, Loader2, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { FormSheet } from "@/components/shared/form-sheet";
import { FormDialog } from "@/components/shared/form-dialog";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getStaffList, addStaff, getRoles, getStaffRoles, addStaffToRole, removeStaffFromRole } from "@/lib/api/staff";
import type { Staff, Role, StaffRole } from "@/types/staff";

export default function StaffManagementPage() {
  const [items, setItems] = useState<Staff[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({ userName: "", email: "", firstName: "", lastName: "", password: "" });

  // Role assignment state
  const [roleDialogOpen, setRoleDialogOpen] = useState(false);
  const [selectedStaff, setSelectedStaff] = useState<Staff | null>(null);
  const [staffRoles, setStaffRoles] = useState<StaffRole[]>([]);
  const [rolesLoading, setRolesLoading] = useState(false);
  const [selectedRole, setSelectedRole] = useState("");
  const [assigningRole, setAssigningRole] = useState(false);
  const [removeConfirm, setRemoveConfirm] = useState<StaffRole | null>(null);

  const loadData = async () => {
    setLoading(true);
    try {
      const [staffRes, rolesRes] = await Promise.all([getStaffList(), getRoles()]);
      if (staffRes?.data) setItems(Array.isArray(staffRes.data) ? staffRes.data : []);
      if (rolesRes?.data) setRoles(Array.isArray(rolesRes.data) ? rolesRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, []);

  const openAdd = () => {
    setFormData({ userName: "", email: "", firstName: "", lastName: "", password: "" });
    setOpen(true);
  };

  const handleSubmit = async () => {
    if (!formData.userName || !formData.email || !formData.firstName || !formData.lastName) {
      toast.error("Please fill all required fields.");
      return;
    }
    setSaving(true);
    try {
      const res = await addStaff({
        user_name: formData.userName,
        email: formData.email,
        first_name: formData.firstName,
        last_name: formData.lastName,
        password: formData.password,
      });
      if (res?.isSuccess) { toast.success("Staff added."); setOpen(false); loadData(); }
      else toast.error(res?.message || "Operation failed.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const openRoleDialog = async (staff: Staff) => {
    setSelectedStaff(staff);
    setSelectedRole("");
    setRoleDialogOpen(true);
    setRolesLoading(true);
    try {
      const res = await getStaffRoles(staff.staffId);
      if (res?.data) setStaffRoles(Array.isArray(res.data) ? res.data : []);
      else setStaffRoles([]);
    } catch { setStaffRoles([]); } finally { setRolesLoading(false); }
  };

  const handleAssignRole = async () => {
    if (!selectedStaff || !selectedRole) return;
    setAssigningRole(true);
    try {
      const res = await addStaffToRole({ user_id: selectedStaff.staffId, role_names: [selectedRole] });
      if (res?.isSuccess) {
        toast.success("Role assigned.");
        const updated = await getStaffRoles(selectedStaff.staffId);
        if (updated?.data) setStaffRoles(Array.isArray(updated.data) ? updated.data : []);
        setSelectedRole("");
      } else toast.error(res?.message || "Failed to assign role.");
    } catch { toast.error("An error occurred."); } finally { setAssigningRole(false); }
  };

  const handleRemoveRole = async () => {
    if (!selectedStaff || !removeConfirm) return;
    try {
      const res = await removeStaffFromRole(selectedStaff.staffId, removeConfirm.roleName ?? "");
      if (res?.isSuccess) {
        toast.success("Role removed.");
        const updated = await getStaffRoles(selectedStaff.staffId);
        if (updated?.data) setStaffRoles(Array.isArray(updated.data) ? updated.data : []);
      } else toast.error(res?.message || "Failed to remove role.");
    } catch { toast.error("An error occurred."); }
  };

  const columns: ColumnDef<Staff>[] = [
    { accessorKey: "employeeNumber", header: "Staff ID" },
    { accessorKey: "firstName", header: "First Name" },
    { accessorKey: "lastName", header: "Last Name" },
    { accessorKey: "userName", header: "Username" },
    { accessorKey: "email", header: "Email" },
    { accessorKey: "isActive", header: "Status", cell: ({ row }) => <Badge variant={row.original.isActive ? "default" : "outline"}>{row.original.isActive ? "Active" : "Disabled"}</Badge> },
    {
      id: "actions", header: "Actions", cell: ({ row }) => (
        <div className="flex gap-1">
          <Button size="sm" variant="ghost" onClick={() => openRoleDialog(row.original)} title="Manage Roles">
            <Shield className="h-3.5 w-3.5" />
          </Button>
        </div>
      ),
    },
  ];

  if (loading) return <div><PageHeader title="Manage Users" breadcrumbs={[{ label: "Settings" }, { label: "Manage Users" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Manage Users" description="Add staff and manage role assignments" breadcrumbs={[{ label: "Settings" }, { label: "Manage Users" }]} actions={<Button size="sm" onClick={openAdd}><Plus className="mr-2 h-4 w-4" />Add Staff</Button>} />
      <DataTable columns={columns} data={items} searchKey="firstName" searchPlaceholder="Search staff..." />

      {/* Add Staff Sheet */}
      <FormSheet open={open} onOpenChange={setOpen} title="Add Staff">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Username *</Label><Input value={formData.userName} onChange={(e) => setFormData({ ...formData, userName: e.target.value })} /></div>
          <div className="space-y-2"><Label>Email *</Label><Input type="email" value={formData.email} onChange={(e) => setFormData({ ...formData, email: e.target.value })} /></div>
          <div className="space-y-2"><Label>First Name *</Label><Input value={formData.firstName} onChange={(e) => setFormData({ ...formData, firstName: e.target.value })} /></div>
          <div className="space-y-2"><Label>Last Name *</Label><Input value={formData.lastName} onChange={(e) => setFormData({ ...formData, lastName: e.target.value })} /></div>
          <div className="space-y-2"><Label>Password</Label><Input type="password" value={formData.password} onChange={(e) => setFormData({ ...formData, password: e.target.value })} /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleSubmit} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Save</Button>
          </div>
        </div>
      </FormSheet>

      {/* Role Assignment Dialog */}
      <FormDialog open={roleDialogOpen} onOpenChange={setRoleDialogOpen} title={`Assign Roles — ${selectedStaff?.firstName} ${selectedStaff?.lastName}`} className="sm:max-w-lg">
        <div className="space-y-4">
          <div className="flex gap-2">
            <Select value={selectedRole} onValueChange={setSelectedRole}>
              <SelectTrigger className="flex-1"><SelectValue placeholder="Select role to assign" /></SelectTrigger>
              <SelectContent>{roles.map((r) => <SelectItem key={r.roleName} value={r.roleName}>{r.roleName}</SelectItem>)}</SelectContent>
            </Select>
            <Button onClick={handleAssignRole} disabled={!selectedRole || assigningRole}>
              {assigningRole ? <Loader2 className="h-4 w-4 animate-spin" /> : "Assign"}
            </Button>
          </div>

          {rolesLoading ? (
            <div className="flex items-center justify-center py-4"><Loader2 className="h-5 w-5 animate-spin" /></div>
          ) : staffRoles.length > 0 ? (
            <div className="space-y-2">
              <Label className="text-muted-foreground">Assigned Roles</Label>
              <div className="space-y-1">
                {staffRoles.map((sr) => (
                  <div key={sr.staffRoleId} className="flex items-center justify-between rounded-md border px-3 py-2 text-sm">
                    <span>{sr.roleName}</span>
                    <Button size="sm" variant="ghost" onClick={() => setRemoveConfirm(sr)}><Trash2 className="h-3.5 w-3.5 text-destructive" /></Button>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">No roles assigned.</p>
          )}
        </div>
      </FormDialog>

      {/* Remove Role Confirmation */}
      <ConfirmationDialog
        open={!!removeConfirm}
        onOpenChange={() => setRemoveConfirm(null)}
        title="Remove Role"
        description={`Remove role "${removeConfirm?.roleName}" from this staff member?`}
        confirmLabel="Remove"
        variant="destructive"
        onConfirm={handleRemoveRole}
      />
    </div>
  );
}
