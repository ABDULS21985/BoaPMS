import { get, post, del } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type { Staff, Role, Permission, StaffRole, RoleWithPermissions } from "@/types/staff";

// --- Staff ---
export const getStaffList = (searchString?: string) =>
  get<BaseAPIResponse<Staff[]>>(searchString ? `/staff?searchString=${searchString}` : "/staff");
export const addStaff = (data: unknown) =>
  post<ResponseVm>("/staff", data);
export const getStaffDetails = (staffId: string) =>
  get<BaseAPIResponse<Staff>>(`/staff?searchString=${staffId}`);

// --- Roles ---
export const getRoles = () =>
  get<BaseAPIResponse<Role[]>>("/staff/roles");
export const addRole = (data: unknown) =>
  post<ResponseVm>("/staff/roles", data);
export const deleteRole = (roleName: string) =>
  del<ResponseVm>(`/staff/roles?roleName=${encodeURIComponent(roleName)}`);

// --- Staff Roles ---
export const getStaffRoles = (staffId: string) =>
  get<BaseAPIResponse<StaffRole[]>>(`/staff/roles/by-staff?id=${staffId}`);
export const addStaffToRole = (data: unknown) =>
  post<ResponseVm>("/staff/roles/assign", data);
export const removeStaffFromRole = (userId: string, roleName: string) =>
  del<ResponseVm>(`/staff/roles/remove?userId=${userId}&roleName=${encodeURIComponent(roleName)}`);

// --- Role Permissions ---
export const getAllRolesWithPermissions = (roleId?: string) =>
  get<BaseAPIResponse<RoleWithPermissions>>(roleId ? `/rolemgmt/roles-with-permission?roleId=${roleId}` : "/rolemgmt/roles-with-permission");
export const getAllPermissions = () =>
  get<BaseAPIResponse<Permission[]>>("/rolemgmt/permissions");
export const getRolePermissions = (roleId: string) =>
  get<BaseAPIResponse<Permission[]>>(`/rolemgmt/permissions?roleId=${roleId}`);
export const addPermissionToRole = (data: unknown) =>
  post<ResponseVm>("/rolemgmt/permissions", data);
export const removePermissionFromRole = (roleId: string, permissionId: string) =>
  del<ResponseVm>(`/rolemgmt/permissions?roleId=${roleId}&permissionId=${permissionId}`);

// --- Employee info (ERP) ---
export const getEmployees = () =>
  get<BaseAPIResponse<unknown[]>>("/employees");
export const getEmployeesByDivision = (divisionId: number) =>
  get<BaseAPIResponse<unknown[]>>(`/employees/division/${divisionId}`);
export const getEmployeeDetails = (employeeNumber: string) =>
  get<BaseAPIResponse<unknown>>(`/employees/${employeeNumber}`);
