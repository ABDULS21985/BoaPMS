"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useSession } from "next-auth/react";
import {
  Home,
  BarChart3,
  Users,
  Network,
  FolderKanban,
  MessageSquare,
  ListChecks,
  ClipboardCheck,
  BarChart,
  HelpCircle,
  ChevronRight,
  LineChart,
  Wrench,
  Settings2,
  Briefcase,
  UserCog,
  CheckCircle,
  Gauge,
  AlertTriangle,
  Settings,
  SquarePen,
  GitBranch,
} from "lucide-react";
import {
  Sidebar as ShadcnSidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarMenuSub,
  SidebarMenuSubItem,
  SidebarMenuSubButton,
  SidebarHeader,
  SidebarFooter,
} from "@/components/ui/sidebar";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { Button } from "@/components/ui/button";
import { Roles } from "@/stores/auth-store";
import { cn } from "@/lib/utils";

function hasAnyRole(userRoles: string[], requiredRoles: string[]): boolean {
  return userRoles.some((r) => requiredRoles.includes(r));
}

function hasHRRole(userRoles: string[]): boolean {
  return userRoles.some(
    (r) =>
      r === Roles.Admin ||
      r === Roles.SuperAdmin ||
      r === Roles.HrAdmin ||
      r === Roles.HrApprover ||
      r.toLowerCase().includes("hr")
  );
}

interface NavSubItem {
  label: string;
  href: string;
  roles?: string[];
}

interface NavItem {
  label: string;
  icon: React.ElementType;
  href?: string;
  roles?: string[];
  children?: NavSubItem[];
}

export function AppSidebar() {
  const pathname = usePathname();
  const { data: session } = useSession();
  const userRoles = session?.user?.roles ?? [];

  const isActive = (href: string) => pathname === href;
  const isGroupActive = (items: NavSubItem[]) =>
    items.some((i) => pathname === i.href);

  // ─── Performance Section ───
  const performanceItems: NavItem[] = [
    { label: "Overview", icon: BarChart3, href: "/overview" },
    {
      label: "My Performance",
      icon: Users,
      children: [
        { label: "Objectives", href: "/my-objectives" },
        { label: "Work Products", href: "/myworkproducts" },
      ],
    },
    {
      label: "Line Manager",
      icon: Network,
      roles: [
        Roles.HeadOfOffice,
        Roles.Supervisor,
        Roles.HeadOfDivision,
        Roles.HeadOfDepartment,
        Roles.SuperAdmin,
      ],
      children: [
        {
          label: "Manage Staff Performance",
          href: "/direct-reports-planning",
        },
      ],
    },
    {
      label: "Ad-Hoc",
      icon: FolderKanban,
      children: [
        { label: "Planning", href: "/ad-hoc-planning" },
        { label: "Projects", href: "/my-projects" },
        { label: "Committees", href: "/my-committees" },
        { label: "Evaluation", href: "/ad-hoc-evaluation" },
      ],
    },
    { label: "360 Feedback", icon: MessageSquare, href: "/feedback-reviews" },
    { label: "Grievances", icon: AlertTriangle, href: "/my-grievances" },
    {
      label: "Outcome Evaluation",
      icon: CheckCircle,
      href: "/outcome-evaluation",
      roles: [Roles.SmdOutcomeEvaluator],
    },
    {
      label: "Score Card",
      icon: Gauge,
      href: "/performance-score-card",
    },
    {
      label: "Assigned Requests",
      icon: ClipboardCheck,
      children: [
        { label: "My Requests", href: "/myrequests" },
        { label: "Breached Requests", href: "/my-breached-requests" },
        { label: "Pending Requests", href: "/assigned-pending-requests" },
        {
          label: "All Requests",
          href: "/requests-review-period",
          roles: [
            Roles.Admin,
            Roles.SuperAdmin,
            Roles.HrReportAdmin,
            Roles.HrAdmin,
            Roles.HrApprover,
          ],
        },
      ],
    },
    {
      label: "Reports",
      icon: BarChart,
      roles: [
        Roles.Admin,
        Roles.SuperAdmin,
        Roles.HrReportAdmin,
        Roles.HrAdmin,
        Roles.SecurityAdmin,
        Roles.Smd,
      ],
      children: [
        {
          label: "Audit Logs",
          href: "/audit-logs",
          roles: [Roles.SecurityAdmin],
        },
        {
          label: "Staff Review Periods",
          href: "/staff-review-periods",
          roles: [
            Roles.Admin,
            Roles.SuperAdmin,
            Roles.HrReportAdmin,
            Roles.HrAdmin,
          ],
        },
        {
          label: "Objectives & Work Products",
          href: "/review-period-report",
          roles: [
            Roles.Admin,
            Roles.SuperAdmin,
            Roles.HrReportAdmin,
            Roles.HrAdmin,
          ],
        },
        {
          label: "Performance Scores Report",
          href: "/pms-period-scores-report",
          roles: [
            Roles.Admin,
            Roles.SuperAdmin,
            Roles.HrReportAdmin,
            Roles.HrAdmin,
          ],
        },
        {
          label: "360 Feedback Report",
          href: "/feedback360-report",
          roles: [
            Roles.Admin,
            Roles.SuperAdmin,
            Roles.HrReportAdmin,
            Roles.HrAdmin,
          ],
        },
        {
          label: "Committees Report",
          href: "/committees-report",
          roles: [
            Roles.Admin,
            Roles.SuperAdmin,
            Roles.HrReportAdmin,
            Roles.HrAdmin,
          ],
        },
        {
          label: "Grievance Report",
          href: "/grievance-report",
          roles: [
            Roles.Admin,
            Roles.SuperAdmin,
            Roles.HrReportAdmin,
            Roles.HrAdmin,
          ],
        },
        {
          label: "Projects Report",
          href: "/projects-report",
          roles: [
            Roles.Admin,
            Roles.SuperAdmin,
            Roles.HrReportAdmin,
            Roles.HrAdmin,
            Roles.Smd,
          ],
        },
      ],
    },
    { label: "Help & Support", icon: HelpCircle, href: "/help-support" },
  ];

  // ─── Competency Section ───
  const competencyItems: NavItem[] = [
    {
      label: "Review Competencies",
      icon: SquarePen,
      href: "/reviews",
    },
    {
      label: "Gap Analysis",
      icon: GitBranch,
      href: "/competency_profiles",
    },
    {
      label: "Development Plan",
      icon: ListChecks,
      href: "/competency_gaps",
    },
    {
      label: "Manage Staff",
      icon: Users,
      roles: [
        Roles.SuperAdmin,
        Roles.Supervisor,
        Roles.HeadOfOffice,
        Roles.HeadOfDivision,
        Roles.HeadOfDepartment,
      ],
      children: [
        { label: "Development Task", href: "/manager_overview" },
        { label: "Approve Job Role", href: "/approve_jobrole_update" },
      ],
    },
    {
      label: "Reports",
      icon: BarChart,
      roles: [
        Roles.HeadOfOffice,
        Roles.HeadOfDivision,
        Roles.HeadOfDepartment,
        Roles.SuperAdmin,
      ],
      children: [
        { label: "Group Report", href: "/group_profiles" },
        { label: "Behavioral Matrix", href: "/competency_matrix" },
        { label: "Technical Matrix", href: "/technical_competency_matrix" },
      ],
    },
  ];

  // ─── Admin: Performance Setup ───
  const performanceSetupItems: NavSubItem[] = [
    { label: "Bank Strategies", href: "/strategysetup" },
    { label: "Strategic Themes", href: "/strategic-theme-setup" },
    { label: "Objective Categories", href: "/objective-categories" },
    { label: "Objectives Bulk Upload", href: "/objective-setup" },
    { label: "Enterprise Objectives", href: "/SetupEnterpriseObjective" },
    { label: "Departmental Objectives", href: "/DepartmentObjectiveSetup" },
    { label: "Divisional Objectives", href: "/DivisionObjectiveSetup" },
    { label: "Office Objectives", href: "/OfficeObjectiveSetup" },
    { label: "PMS Review Periods", href: "/review-periods-list" },
    { label: "Work Product Upload", href: "/work-product-setup" },
    { label: "PMS Competencies", href: "/pms-competencies-setup" },
    { label: "Feedback Questionnaires", href: "/feedback-questionaires" },
    { label: "Evaluation Statements", href: "/evaluation-options-setup" },
    {
      label: "Configurations",
      href: "/pms-configs",
      roles: [Roles.Admin, Roles.HrAdmin],
    },
    {
      label: "System Settings",
      href: "/settings",
      roles: [Roles.Admin, Roles.HrAdmin],
    },
  ];

  // ─── Admin: Manage Competencies ───
  const manageCompetenciesItems: NavSubItem[] = [
    {
      label: "Technical Competencies",
      href: "/jobrolecompetencies",
      roles: [Roles.Admin, Roles.SuperAdmin, Roles.HrAdmin],
    },
    {
      label: "Behavioral Competencies",
      href: "/behavioral_competencies",
      roles: [Roles.Admin, Roles.SuperAdmin, Roles.HrAdmin],
    },
    {
      label: "Competencies",
      href: "/competencies",
      roles: [Roles.Admin, Roles.SuperAdmin, Roles.HrAdmin],
    },
    {
      label: "Approve Competencies",
      href: "/pending_competencies",
      roles: [Roles.SuperAdmin, Roles.HrApprover],
    },
  ];

  // ─── Admin: Competency Setup ───
  const competencySetupItems: NavSubItem[] = [
    { label: "Review Types", href: "/reviewtypes" },
    { label: "Review Period", href: "/reviewperiod" },
    { label: "Categories", href: "/competency_group" },
    { label: "Ratings", href: "/ratings" },
    { label: "Grading", href: "/competency_category_grade" },
    { label: "Intervention Types", href: "/development_interventions" },
    {
      label: "Approve Review Period",
      href: "/approve_reviewperiod",
      roles: [Roles.Admin, Roles.HrApprover],
    },
  ];

  // ─── Admin: Portfolio Setup ───
  const portfolioItems: NavSubItem[] = [
    { label: "Office Job Roles", href: "/office_jobRoles" },
    { label: "Job Roles", href: "/jobroles" },
    { label: "Assign Grade Groups", href: "/assign_jobGrades" },
    { label: "Grade Groups", href: "/grade_groups" },
    { label: "Job Grades", href: "/job_grades" },
  ];

  // ─── Admin: System Organogram ───
  const organogramItems: NavSubItem[] = [
    { label: "Offices", href: "/offices" },
    { label: "Divisions", href: "/divisions" },
    { label: "Departments", href: "/departments" },
    { label: "Bank Year", href: "/bankyear" },
    { label: "Staff", href: "/staff_list" },
  ];

  // ─── Settings: User Management ───
  const userMgtItems: NavSubItem[] = [
    { label: "Manage Users", href: "/staff_mgt" },
    { label: "Manage Roles", href: "/manage_roles" },
    { label: "Permissions", href: "/permissions" },
  ];

  const showAdminSection = hasAnyRole(userRoles, [
    Roles.Admin,
    Roles.SuperAdmin,
    Roles.HrAdmin,
    Roles.HrApprover,
  ]);
  const showCompetencySetup =
    hasAnyRole(userRoles, [Roles.Admin]) || hasHRRole(userRoles);
  const showPortfolioOrg = hasAnyRole(userRoles, [
    Roles.Admin,
    Roles.HrAdmin,
    Roles.SuperAdmin,
  ]);
  const showSettings = hasAnyRole(userRoles, [
    Roles.Admin,
    Roles.SuperAdmin,
    Roles.SecurityAdmin,
  ]);

  function renderNavItem(item: NavItem) {
    if (item.roles && !hasAnyRole(userRoles, item.roles)) return null;

    if (item.href && !item.children) {
      return (
        <SidebarMenuItem key={item.label}>
          <SidebarMenuButton asChild isActive={isActive(item.href)}>
            <Link href={item.href}>
              <item.icon className="h-4 w-4" />
              <span>{item.label}</span>
            </Link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      );
    }

    if (item.children) {
      const visibleChildren = item.children.filter(
        (c) => !c.roles || hasAnyRole(userRoles, c.roles)
      );
      if (visibleChildren.length === 0) return null;

      return (
        <Collapsible
          key={item.label}
          defaultOpen={isGroupActive(visibleChildren)}
          className="group/collapsible"
        >
          <SidebarMenuItem>
            <CollapsibleTrigger asChild>
              <SidebarMenuButton>
                <item.icon className="h-4 w-4" />
                <span>{item.label}</span>
                <ChevronRight className="ml-auto h-4 w-4 transition-transform group-data-[state=open]/collapsible:rotate-90" />
              </SidebarMenuButton>
            </CollapsibleTrigger>
            <CollapsibleContent>
              <SidebarMenuSub>
                {visibleChildren.map((child) => (
                  <SidebarMenuSubItem key={child.href}>
                    <SidebarMenuSubButton
                      asChild
                      isActive={isActive(child.href)}
                    >
                      <Link href={child.href}>{child.label}</Link>
                    </SidebarMenuSubButton>
                  </SidebarMenuSubItem>
                ))}
              </SidebarMenuSub>
            </CollapsibleContent>
          </SidebarMenuItem>
        </Collapsible>
      );
    }

    return null;
  }

  function renderAdminGroup(
    label: string,
    icon: React.ElementType,
    items: NavSubItem[]
  ) {
    const visibleItems = items.filter(
      (i) => !i.roles || hasAnyRole(userRoles, i.roles)
    );
    if (visibleItems.length === 0) return null;
    const Icon = icon;

    return (
      <Collapsible
        key={label}
        defaultOpen={isGroupActive(visibleItems)}
        className="group/collapsible"
      >
        <SidebarMenuItem>
          <CollapsibleTrigger asChild>
            <SidebarMenuButton>
              <Icon className="h-4 w-4" />
              <span>{label}</span>
              <ChevronRight className="ml-auto h-4 w-4 transition-transform group-data-[state=open]/collapsible:rotate-90" />
            </SidebarMenuButton>
          </CollapsibleTrigger>
          <CollapsibleContent>
            <SidebarMenuSub>
              {visibleItems.map((item) => (
                <SidebarMenuSubItem key={item.href}>
                  <SidebarMenuSubButton
                    asChild
                    isActive={isActive(item.href)}
                  >
                    <Link href={item.href}>{item.label}</Link>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              ))}
            </SidebarMenuSub>
          </CollapsibleContent>
        </SidebarMenuItem>
      </Collapsible>
    );
  }

  return (
    <ShadcnSidebar collapsible="icon" variant="sidebar">
      <SidebarHeader className="border-b px-4 py-3">
        <Link href="/" className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <span className="text-sm font-bold">P</span>
          </div>
          <span className="font-semibold group-data-[collapsible=icon]:hidden">
            PMS
          </span>
        </Link>
      </SidebarHeader>

      <SidebarContent>
        {/* Dashboard */}
        <SidebarGroup>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton asChild isActive={isActive("/")}>
                <Link href="/">
                  <Home className="h-4 w-4" />
                  <span>Dashboard</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>

        {/* Performance */}
        <SidebarGroup>
          <SidebarGroupLabel>Performance</SidebarGroupLabel>
          <SidebarMenu>
            {performanceItems.map((item) => renderNavItem(item))}
          </SidebarMenu>
        </SidebarGroup>

        {/* Competency */}
        <SidebarGroup>
          <SidebarGroupLabel>Competency</SidebarGroupLabel>
          <SidebarMenu>
            {competencyItems.map((item) => renderNavItem(item))}
          </SidebarMenu>
        </SidebarGroup>

        {/* Administrator Section */}
        {(showAdminSection || hasHRRole(userRoles)) && (
          <SidebarGroup>
            <SidebarGroupLabel>Administrator</SidebarGroupLabel>
            <SidebarMenu>
              {showAdminSection &&
                renderAdminGroup(
                  "Performance Setup",
                  LineChart,
                  performanceSetupItems
                )}
              {showAdminSection &&
                renderAdminGroup(
                  "Manage Competencies",
                  Wrench,
                  manageCompetenciesItems
                )}
              {showCompetencySetup &&
                renderAdminGroup(
                  "Competency Setup",
                  Settings2,
                  competencySetupItems
                )}
              {showPortfolioOrg &&
                renderAdminGroup(
                  "Portfolio Setup",
                  Briefcase,
                  portfolioItems
                )}
              {showPortfolioOrg &&
                renderAdminGroup(
                  "System Organogram",
                  Network,
                  organogramItems
                )}
            </SidebarMenu>
          </SidebarGroup>
        )}

        {/* Settings */}
        {showSettings && (
          <SidebarGroup>
            <SidebarGroupLabel>Settings</SidebarGroupLabel>
            <SidebarMenu>
              {renderAdminGroup(
                "User Management",
                UserCog,
                userMgtItems
              )}
            </SidebarMenu>
          </SidebarGroup>
        )}
      </SidebarContent>

      <SidebarFooter className="border-t p-4">
        <Button variant="outline" className="w-full" asChild>
          <Link href="/api/auth/signout">Log Out</Link>
        </Button>
      </SidebarFooter>
    </ShadcnSidebar>
  );
}
