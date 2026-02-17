import { get, post, del } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type { Staff, Role, Permission, StaffRole } from "@/types/staff";

// --- Staff ---
export const getStaffList = () =>
  get<BaseAPIResponse<Staff[]>>("/staff");
export const addStaff = (data: unknown) =>
  post<ResponseVm>("/staff", data);
export const getStaffDetails = (staffId: string) =>
  get<BaseAPIResponse<Staff>>(`/staff/${staffId}`);

// --- Roles ---
export const getRoles = () =>
  get<BaseAPIResponse<Role[]>>("/roles");
export const addRole = (data: unknown) =>
  post<ResponseVm>("/roles", data);
export const deleteRole = (roleId: number) =>
  del<ResponseVm>(`/roles/${roleId}`);

// --- Staff Roles ---
export const getStaffRoles = (staffId: string) =>
  get<BaseAPIResponse<StaffRole[]>>(`/staff/${staffId}/roles`);
export const addStaffToRole = (data: unknown) =>
  post<ResponseVm>("/staff/roles", data);
export const removeStaffFromRole = (data: unknown) =>
  post<ResponseVm>("/staff/roles/remove", data);

// --- Permissions ---
export const getRolePermissions = (roleId: number) =>
  get<BaseAPIResponse<Permission[]>>(`/roles/${roleId}/permissions`);
export const addPermissionToRole = (data: unknown) =>
  post<ResponseVm>("/roles/permissions", data);
export const removePermissionFromRole = (data: unknown) =>
  post<ResponseVm>("/roles/permissions/remove", data);

// --- Employee info (ERP) ---
export const getEmployees = () =>
  get<BaseAPIResponse<unknown[]>>("/employees");
export const getEmployeesByDivision = (divisionId: number) =>
  get<BaseAPIResponse<unknown[]>>(`/employees/division/${divisionId}`);
export const getEmployeeDetails = (employeeNumber: string) =>
  get<BaseAPIResponse<unknown>>(`/employees/${employeeNumber}`);
