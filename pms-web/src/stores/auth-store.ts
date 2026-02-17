import { create } from "zustand";

interface AuthState {
  userId: string;
  name: string;
  email: string;
  firstName: string;
  lastName: string;
  roles: string[];
  permissions: string[];
  organizationalUnit: string;
  setUser: (user: Partial<AuthState>) => void;
  clearUser: () => void;
  hasRole: (role: string) => boolean;
  hasAnyRole: (roles: string[]) => boolean;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  userId: "",
  name: "",
  email: "",
  firstName: "",
  lastName: "",
  roles: [],
  permissions: [],
  organizationalUnit: "",
  setUser: (user) => set((state) => ({ ...state, ...user })),
  clearUser: () =>
    set({
      userId: "",
      name: "",
      email: "",
      firstName: "",
      lastName: "",
      roles: [],
      permissions: [],
      organizationalUnit: "",
    }),
  hasRole: (role) => get().roles.includes(role),
  hasAnyRole: (roles) => get().roles.some((r) => roles.includes(r)),
}));

// Role constants matching Go auth/roles.go
export const Roles = {
  SuperAdmin: "SuperAdmin",
  Admin: "Admin",
  Staff: "Staff",
  HeadOfOffice: "HeadOfOffice",
  HeadOfDivision: "HeadOfDivision",
  HeadOfDepartment: "HeadOfDepartment",
  Supervisor: "Supervisor",
  HRD: "HRD",
  HrAdmin: "HrAdmin",
  HrApprover: "HrApprover",
  HrReportAdmin: "HrReportAdmin",
  SecurityAdmin: "SecurityAdmin",
  Smd: "Smd",
  SmdOutcomeEvaluator: "SmdOutcomeEvaluator",
  Reviewer: "Reviewer",
  Approver: "Approver",
} as const;

export type RoleName = (typeof Roles)[keyof typeof Roles];
