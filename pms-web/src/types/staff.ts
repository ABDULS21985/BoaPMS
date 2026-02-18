export interface Staff {
  staffId: string;
  employeeNumber: string;
  firstName: string;
  lastName: string;
  email: string;
  userName: string;
  phone?: string;
  jobName?: string;
  gradeName?: string;
  departmentName?: string;
  divisionName?: string;
  officeName?: string;
  supervisorId?: string;
  supervisorName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface Role {
  roleId: number;
  roleName: string;
  description?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface Permission {
  permissionId: number;
  permissionName: string;
  description?: string;
  isActive: boolean;
}

export interface StaffRole {
  staffRoleId: number;
  staffId: string;
  roleId: number;
  roleName?: string;
  staffName?: string;
  isActive: boolean;
}

export interface StaffIdMask {
  employeeNumber: string;
  currentStaffPhoto?: string;
  firstName?: string;
  lastName?: string;
}

export interface RoleWithPermissions {
  allPermissions: Permission[];
  rolesAndPermissions: {
    roleId: string;
    roleName: string;
    permissions: Permission[];
  };
}

export interface AuditLog {
  userName: string;
  auditEventDateUTC: string;
  auditEventType: number;
  tableName: string;
  recordId: string;
  fieldName: string;
  originalValue: string;
  newValue: string;
}
