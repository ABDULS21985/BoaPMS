# PROMPT 9D: PMS Report Pages Part 2 (4 Pages)

## Context
You are converting a .NET Blazor PMS application to Next.js 15 (App Router) + TypeScript + Tailwind CSS + shadcn/ui. This prompt covers the remaining 4 PMS report pages under the `(dashboard)` route group.

**PREREQUISITE**: PROMPT 9A must be completed first (types, API functions, helpers).

## Target Directory
`/Users/enaira/Desktop/ENTERPRISE PMS/go-pms/pms-web/`

## Existing API Functions (DO NOT recreate)
```ts
// pms-engine.ts
getCommittees(params?) → Committee[]
updateCommittee(data) → ResponseVm
getProjects(params?) → Project[]
updateProject(data) → ResponseVm
getAllCompetencyReviewFeedbacks(staffId) → CompetencyReviewFeedback[]

// grievance.ts
getGrievanceReport() → Grievance[]
updateGrievance(data) → ResponseVm
submitGrievanceResolution(data) → ResponseVm
updateGrievanceResolution(data) → ResponseVm
escalateGrievance(data) → ResponseVm
closeGrievance(data) → ResponseVm
getGrievanceTypes() → { id: number; name: string }[]

// dashboard.ts
getEmployeeDetail(userId) → EmployeeErpDetails

// organogram.ts
getDepartments(), getDivisions(), getDivisionsByDepartment(deptId), getOffices(), getOfficesByDivision(divId)

// review-periods.ts
getReviewPeriods() → PerformanceReviewPeriod[]
```

## Existing Reusable Components
```ts
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { FormDialog } from "@/components/shared/form-dialog";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
```

## Standard Page Pattern
All pages follow: "use client" → useState/useEffect → loadData → columns → DataTable. See PROMPT 9C for detailed pattern.

---

## Page 1: `src/app/(dashboard)/feedback360-report/page.tsx`

**Route**: `/feedback360-report`
**Source**: `Feedback360Report.razor` + `.razor.cs`
**Access**: Admin, SuperAdmin, HrReportAdmin, HrAdmin
**Sidebar**: Reports → 360 Feedback Report

### Behavior (from .NET source)
1. Admin-only page.
2. On load:
   - Call `getReviewPeriods()` — populate review period filter
   - Call `getDepartments()` — populate department filter
3. **Cascading filters**: Same pattern as PMS Period Scores Report (PROMPT 9C, Page 5)
   - Department → Division → Office cascading dropdowns
   - Review Period dropdown
4. On "Search" click:
   - This report needs staff 360 feedback data. The existing API `getAllCompetencyReviewFeedbacks(staffId)` is per-staff.
   - For admin report view, use a batch approach: load period scores to get staff list, then show summary.
   - Alternative: Add a new API function if needed:
     ```ts
     export const get360FeedbackReport = (params?: { reviewPeriodId?: string; departmentId?: number; divisionId?: number; officeId?: number }) => {
       const q = new URLSearchParams();
       if (params?.reviewPeriodId) q.set("reviewPeriodId", params.reviewPeriodId);
       if (params?.departmentId) q.set("departmentId", String(params.departmentId));
       if (params?.divisionId) q.set("divisionId", String(params.divisionId));
       if (params?.officeId) q.set("officeId", String(params.officeId));
       return get<BaseAPIResponse<CompetencyReviewFeedback[]>>(`/pms-engine/competency-review/feedbacks/report?${q}`);
     };
     ```
   - If the endpoint doesn't exist, use `getPeriodScores(reviewPeriodId)` to get staff list and display basic info with a "View Details" action that calls `getAllCompetencyReviewFeedbacks(staffId)` per-row.

### UI Structure
- **PageHeader**: title="360 Feedback Report", breadcrumbs: Home → 360 Feedback Report
- **Filter Row**: Review Period, Department (cascading), Division (cascading), Office (cascading), Search button
- **DataTable**:
  - Columns: Staff ID, Staff Name, Department, Division, Office, Feedback Status, Score %, Actions
  - Actions: "View Details" button → opens FormDialog with competency ratings

### Detail Dialog Content
When "View Details" is clicked:
- Call `getAllCompetencyReviewFeedbacks(staffId)`
- Show table inside FormDialog:
  - Columns: Competency, Reviewer, Rating Score, Status
  - Summary row: Average score

---

## Page 2: `src/app/(dashboard)/committees-report/page.tsx`

**Route**: `/committees-report`
**Source**: `CommitteesReport.razor` + `.razor.cs`
**Access**: Admin, SuperAdmin, HrReportAdmin, HrAdmin
**Sidebar**: Reports → Committees Report

### Behavior (from .NET source)
1. Admin-only page.
2. On load, call `getCommittees()` — no params = get ALL committees (admin view).
3. Display as DataTable.
4. Action: "Re-assign Chair Person" — opens a FormDialog where admin can search for a new staff member by ID and reassign.

### UI Structure
- **PageHeader**: title="Committees Report", breadcrumbs: Home → Committees Report
- **DataTable**:
  - Columns: #, Name, Description, Objective/KPI, Department, Chair Person, Status (StatusBadge), Start Date, End Date, Action
  - Search by committee name
  - Action: "Re-assign" button (UserCog icon)

### Re-assign Chair Person Dialog (FormDialog)
- Fields:
  - Current Chair Person (read-only display)
  - New Chair Person Staff ID (`<Input>`)
  - "Verify" button → calls `getEmployeeDetail(staffId)` → shows employee name if found
  - Verified employee name display
- Footer: Cancel + "Re-assign" button
- On confirm: Call `updateCommittee({ committeeId, chairPersonId: newStaffId })` (or appropriate field name)
- Success: toast + reload data

### Key State
```ts
const [committees, setCommittees] = useState<Committee[]>([]);
const [reassignOpen, setReassignOpen] = useState(false);
const [selectedCommittee, setSelectedCommittee] = useState<Committee | null>(null);
const [newChairId, setNewChairId] = useState("");
const [verifiedEmployee, setVerifiedEmployee] = useState<EmployeeErpDetails | null>(null);
const [verifying, setVerifying] = useState(false);
const [saving, setSaving] = useState(false);
```

### Verify Staff Pattern
```ts
const handleVerify = async () => {
  if (!newChairId.trim()) { toast.error("Enter a staff ID."); return; }
  setVerifying(true);
  try {
    const res = await getEmployeeDetail(newChairId.trim());
    if (res?.data?.employeeNumber) {
      setVerifiedEmployee(res.data);
      toast.success(`Found: ${res.data.firstName} ${res.data.lastName}`);
    } else {
      setVerifiedEmployee(null);
      toast.error("Staff not found.");
    }
  } catch { toast.error("Verification failed."); } finally { setVerifying(false); }
};
```

---

## Page 3: `src/app/(dashboard)/grievance-report/page.tsx`

**Route**: `/grievance-report`
**Source**: `GreivanceReport.razor` + `.razor.cs`
**Access**: Admin, SuperAdmin, HrReportAdmin, HrAdmin
**Sidebar**: Reports → Grievance Report

### Behavior (from .NET source)
1. Admin-only page.
2. On load, call `getGrievanceReport()` — returns ALL grievances (admin view).
3. Display as DataTable.
4. Action: "View Details" → opens FormDialog showing full grievance info.

### UI Structure
- **PageHeader**: title="Grievance Report", breadcrumbs: Home → Grievance Report
- **DataTable**:
  - Columns: #, Grievance Type, Subject, Complainant (staff ID), Respondent (staff ID), Current Mediator, Status (StatusBadge), Date Created, Action
  - Search by subject or complainant
  - Action: "Details" button (Eye icon)

### Grievance Detail Dialog (FormDialog)
Display fields (read-only):
- Grievance Type
- Subject
- Description
- Complainant ID + Name (if available)
- Respondent ID + Name
- Current Mediator
- Status
- Date Created
- Resolution Status (if resolved)
- Resolution Notes
- Escalation History (if escalated)

### Grievance Type Interface
The `Grievance` type is defined in `src/types/performance.ts`. Check its fields. Key fields expected:
```ts
interface Grievance {
  grievanceId: string;
  grievanceType: number;
  grievanceTypeName?: string;
  subject: string;
  description?: string;
  complainantStaffId: string;
  respondentStaffId: string;
  currentMediatorStaffId?: string;
  status?: number;
  resolution?: string;
  resolutionDate?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}
```

---

## Page 4: `src/app/(dashboard)/projects-report/page.tsx`

**Route**: `/projects-report`
**Source**: `ProjectsReport.razor` + `.razor.cs`
**Access**: Admin, SuperAdmin, HrReportAdmin, HrAdmin, Smd
**Sidebar**: Reports → Projects Report

### Behavior (from .NET source)
1. Admin-only page (includes Smd role).
2. On load, call `getProjects()` — no params = get ALL projects (admin view).
3. Display as DataTable.
4. Action: "Re-assign Project Manager" — same pattern as Committees Report chair person reassign.

### UI Structure
- **PageHeader**: title="Projects Report", breadcrumbs: Home → Projects Report
- **DataTable**:
  - Columns: #, Name, Description, Objective/KPI (deliverables), Department, Project Manager, Status (StatusBadge), Start Date, End Date, Action
  - Search by project name
  - Action: "Re-assign PM" button (UserCog icon)

### Re-assign Project Manager Dialog (FormDialog)
Identical pattern to Committees Report re-assign:
- Fields: Current PM (read-only), New PM Staff ID input, Verify button, Verified name display
- On confirm: Call `updateProject({ projectId, projectManager: newStaffId })` (or appropriate field name)
- Success: toast + reload

### Key State (same pattern as committees)
```ts
const [projects, setProjects] = useState<Project[]>([]);
const [reassignOpen, setReassignOpen] = useState(false);
const [selectedProject, setSelectedProject] = useState<Project | null>(null);
const [newManagerId, setNewManagerId] = useState("");
const [verifiedEmployee, setVerifiedEmployee] = useState<EmployeeErpDetails | null>(null);
```

---

## Verification
```bash
cd /Users/enaira/Desktop/ENTERPRISE\ PMS/go-pms/pms-web && npm run build
```
Expect zero TypeScript/ESLint errors. All 4 routes should render.

## Files Created Summary
| File | Route |
|------|-------|
| `src/app/(dashboard)/feedback360-report/page.tsx` | `/feedback360-report` |
| `src/app/(dashboard)/committees-report/page.tsx` | `/committees-report` |
| `src/app/(dashboard)/grievance-report/page.tsx` | `/grievance-report` |
| `src/app/(dashboard)/projects-report/page.tsx` | `/projects-report` |
