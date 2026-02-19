# PROMPT 9C: PMS Report Pages Part 1 (5 Pages)

## Context
You are converting a .NET Blazor PMS application to Next.js 15 (App Router) + TypeScript + Tailwind CSS + shadcn/ui. This prompt covers the first 5 PMS report pages under the `(dashboard)` route group.

**PREREQUISITE**: PROMPT 9A must be completed first (types, API functions, helpers).

## Target Directory
`/Users/enaira/Desktop/ENTERPRISE PMS/go-pms/pms-web/`

## Route Group Layout
The `(dashboard)` route group uses `AppSidebar` + `TopNavbar`. Layout at `src/app/(dashboard)/layout.tsx`:
```tsx
<SidebarProvider>
  <AppSidebar />
  <SidebarInset>
    <TopNavbar />
    <main className="flex-1 overflow-auto p-4 md:p-6">{children}</main>
  </SidebarInset>
  <Toaster richColors position="top-right" />
</SidebarProvider>
```

## Existing API Functions (DO NOT recreate)
```ts
// dashboard.ts
getEmployeeDetail(userId), getPeriodScores(reviewPeriodId), getPeriodScoreDetails(reviewPeriodId, staffId)
getStaffReviewPeriods(staffId)

// review-periods.ts
getReviewPeriods(), getReviewPeriodObjectives(reviewPeriodId), getStaffPeriodScore(staffId, reviewPeriodId)

// organogram.ts
getDepartments(), getDivisions(), getDivisionsByDepartment(departmentId), getOffices(), getOfficesByDivision(divisionId)

// pms-engine.ts
getAllStaffWorkProducts(), getWorkProductTasks(workProductId), evaluateWorkProduct(data)
getStaffObjectives(staffId, reviewPeriodId), reInstateWorkProduct(data)
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

## Session & Role Check Pattern
```ts
import { useSession } from "@/lib/auth-context";
const { user } = useSession();
const isAdmin = user?.roles?.some(r => ["Admin", "SuperAdmin", "HrReportAdmin", "HrAdmin"].includes(r));
if (!isAdmin) router.push("/access-denied");
```

## Standard Page Pattern (from existing pages like /projects, /committees)
```tsx
"use client";
export default function ReportPage() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await apiCall();
      if (res?.data) setData(Array.isArray(res.data) ? res.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };
  useEffect(() => { loadData(); }, []);

  const columns: ColumnDef<T>[] = [...];

  if (loading) return (<div><PageHeader ... /><PageSkeleton /></div>);
  return (
    <div className="space-y-6">
      <PageHeader title="..." breadcrumbs={[...]} />
      <DataTable columns={columns} data={data} searchKey="..." />
    </div>
  );
}
```

---

## Page 1: `src/app/(dashboard)/staff-review-periods/page.tsx`

**Route**: `/staff-review-periods`
**Source**: `StaffReviewPeriodReport.razor` + `.razor.cs`
**Access**: Admin, SuperAdmin, HrReportAdmin, HrAdmin
**Sidebar**: Reports → Staff Review Periods

### Behavior (from .NET source)
1. Admin-only page. Check roles on mount — redirect to `/access-denied` if unauthorized.
2. Show a search input for Staff ID with a "Search" button.
3. When user clicks Search, call `getStaffReviewPeriods(staffId)` with the entered staffId.
4. Display matching review periods in a DataTable.
5. Each row has a "Details" action button that opens a FormDialog showing the period score breakdown.
6. On "Details" click: Call `getPeriodScoreDetails(reviewPeriodId, staffId)` to get the `PeriodScoreData`.

### UI Structure
- **PageHeader**: title="Staff Review Periods", breadcrumbs: Home → Staff Review Periods
- **Search area**: `<Input>` for Staff ID + `<Button>` "Search" (with Loader2 spinner during search)
- **DataTable** (only shown after search):
  - Columns: #, Name, Year, Start Date, End Date, Action
  - Action: "Details" button (Info icon)
- **FormDialog** for score details (when Details clicked):
  - Display fields (read-only, not form inputs):
    - Review Period: `{name} - {startDate} to {endDate}`
    - Name: `{staffFullName} ({staffId})`
    - Grade: `{staffGrade}`
    - Office: `{officeName}, {divisionName}, {departmentName}`
    - Performance Score: `{finalScore} / {maxPoint}`
    - Deducted Points: `{hrdDeductedPoints}`
    - Score Percentage: `{scorePercentage}%`
    - Performance Grade: `{finalGradeName}`
  - Close button

### Key State
```ts
const [staffIdSearch, setStaffIdSearch] = useState("");
const [periods, setPeriods] = useState<PerformanceReviewPeriod[]>([]);
const [searching, setSearching] = useState(false);
const [hasSearched, setHasSearched] = useState(false);
const [detailOpen, setDetailOpen] = useState(false);
const [scoreDetails, setScoreDetails] = useState<PeriodScoreData | null>(null);
```

---

## Page 2: `src/app/(dashboard)/review-period-report/page.tsx`

**Route**: `/review-period-report`
**Source**: `ReviewPeriodReport.razor` + `.razor.cs`
**Access**: Admin, SuperAdmin, HrReportAdmin, HrAdmin
**Sidebar**: Reports → Objectives & Work Products

### Behavior (from .NET source)
1. Admin-only page.
2. On load, call `getReviewPeriods()` to get all review periods.
3. Display as DataTable.
4. Each period has 2 action buttons: "View Objectives" and "View Work Products" that navigate to sub-pages.

### UI Structure
- **PageHeader**: title="Objectives & Work Products Report", breadcrumbs: Home → Objectives & Work Products
- **DataTable**:
  - Columns: #, Name, Year, Status (StatusBadge), Start Date, End Date, Actions
  - Actions (2 buttons):
    - "Objectives" (FileText icon) → navigates to `/review-period-report/[reviewPeriodId]/objectives`
    - "Work Products" (Package icon) → navigates to `/review-period-report/[reviewPeriodId]/work-products`
  - Search by period name

---

## Page 3: `src/app/(dashboard)/review-period-report/[reviewPeriodId]/objectives/page.tsx`

**Route**: `/review-period-report/[reviewPeriodId]/objectives`
**Source**: `StaffPlannedObjectives.razor` + `.razor.cs`
**Access**: Admin, SuperAdmin, HrReportAdmin, HrAdmin

### Behavior (from .NET source)
1. Extract `reviewPeriodId` from params.
2. Call `getReviewPeriodObjectives(reviewPeriodId)` — returns array of planned objectives for ALL staff in the period.
3. Display as DataTable.
4. Each rejected objective has a "Re-instate" action.

### UI Structure
- **PageHeader**: title="Staff Planned Objectives", breadcrumbs: Home → Reports → Objectives & Work Products → Objectives
- **DataTable**:
  - Columns: #, ID (truncated), Staff ID, Objective Name, Level/Category, Status (StatusBadge)
  - Actions: "Re-instate" button (only shown for rejected status, uses ConfirmationDialog)
  - Search by staff ID or objective name

### Re-instate Flow
```ts
// On confirm re-instate:
const handleReinstate = async (objectiveId: string) => {
  const res = await reInstateObjective({ objectiveId });
  if (res?.isSuccess) { toast.success("Objective re-instated."); loadData(); }
  else toast.error(res?.message || "Failed.");
};
```

Note: Check if `reInstateObjective` exists in review-periods.ts or pms-engine.ts. If not, add:
```ts
export const reInstateObjective = (data: unknown) =>
  post<ResponseVm>("/review-periods/objectives/reinstate", data);
```

---

## Page 4: `src/app/(dashboard)/review-period-report/[reviewPeriodId]/work-products/page.tsx`

**Route**: `/review-period-report/[reviewPeriodId]/work-products`
**Source**: `StaffWorkProducts.razor` + `.razor.cs`
**Access**: Admin, SuperAdmin, HrReportAdmin, HrAdmin

### Behavior (from .NET source)
1. Extract `reviewPeriodId` from params.
2. Call `getAllStaffWorkProducts()` — returns all work products. Client-side filter by `reviewPeriodId` if the API doesn't support it as a parameter.
3. Display as DataTable with many columns.
4. Actions: "Update Evaluation" (FormDialog), "View Tasks" (FormDialog), "Re-instate" (ConfirmationDialog)

### UI Structure
- **PageHeader**: title="Staff Work Products", breadcrumbs: Home → Reports → Objectives & Work Products → Work Products
- **DataTable**:
  - Columns: #, ID (truncated), Staff ID, Objective, Work Product Name, Type, Max Point, Score, Status (StatusBadge)
  - Search by staff ID or work product name

### Actions

**Update Evaluation** (FormDialog):
- Opens when admin clicks "Evaluate" button on a work product
- Fields: Timeliness (Select: 1-5), Quality (Select: 1-5), Output (Select: 1-5)
- Submit button calls `evaluateWorkProduct(data)` with `{ workProductId, timeliness, quality, output }`

**View Tasks** (FormDialog):
- Opens when clicking "Tasks" button
- Calls `getWorkProductTasks(workProductId)` to load tasks
- Shows table inside dialog: Task Name, Status, Start Date, End Date

**Re-instate** (ConfirmationDialog):
- For cancelled/rejected work products
- Calls `reInstateWorkProduct({ workProductId })`

---

## Page 5: `src/app/(dashboard)/pms-period-scores-report/page.tsx`

**Route**: `/pms-period-scores-report`
**Source**: `PmsScoreReport.razor` + `.razor.cs`
**Access**: Admin, SuperAdmin, HrReportAdmin
**Sidebar**: Reports → Performance Scores Report

### Behavior (from .NET source — most complex report)
1. Admin-only page.
2. On load:
   - Call `getActiveReviewPeriod()` — get current active period
   - Call `getPeriodScores(activePeriodId)` — load initial data
   - Call `getReviewPeriods()` — populate review period dropdown
   - Call `getDepartments()` — populate department dropdown
3. **Cascading filters**:
   - Department change → call `getDivisionsByDepartment(deptId)` → populate Division dropdown
   - Division change → call `getOfficesByDivision(divId)` → populate Office dropdown
4. On "Search" click:
   - Call `getPeriodScores(selectedReviewPeriodId)` — get all scores for period
   - Client-side filter by department, division, office if selected
5. **Export to Excel**: Generate CSV/Excel download from filtered data

### UI Structure
- **PageHeader**: title="Performance Scores Report", breadcrumbs: Home → Performance Scores Report
- **Filter Row** (4 dropdowns + Search button):
  - Review Period (`<Select>`)
  - Department (`<Select>`, cascading)
  - Division (`<Select>`, cascading, depends on Department)
  - Office (`<Select>`, cascading, depends on Division)
  - "Search" button (with loading spinner)
- **Export Button**: "Export to Excel" button in the DataTable toolbar area

- **DataTable** (horizontal scroll for many columns):
  - Columns (20): Staff ID, Name, Department, Division, Office, Staff Grade, Score % (badge), Final Score (badge), Max Point (badge), Performance Grade (badge), Min Objectives (badge), Max Objectives (badge), Deducted Points (badge), Review Period, Year, Start Date, End Date, Strategy, Location ID, Is Under Performing
  - Search by staff ID or name
  - Pagination: 10/20/50/100/500/1000

### Badge Display Pattern
Use shadcn `Badge` component for numeric columns:
```tsx
<Badge variant="outline" className="font-bold">{Math.round(value * 100) / 100}</Badge>
```

### Excel Export Implementation
```ts
const exportToExcel = () => {
  const headers = ["StaffId", "FullName", "Department", "Division", "Office", "Staff Grade",
    "Score Percentage", "Final Score", "FinalGradeName", "MinNoOfObjectives",
    "MaxNoOfObjectives", "HRDDeductedPoints", "Review Period", "Year"];
  const csvContent = [
    headers.join(","),
    ...filteredData.map(row => [
      row.staffId, `"${row.staffFullName}"`, `"${row.departmentName}"`, `"${row.divisionName}"`,
      `"${row.officeName}"`, row.staffGrade, row.scorePercentage.toFixed(2), row.finalScore.toFixed(4),
      row.finalGradeName, row.minNoOfObjectives, row.maxNoOfObjectives, row.hrdDeductedPoints,
      `"${row.reviewPeriod}"`, row.year,
    ].join(","))
  ].join("\n");

  const blob = new Blob([csvContent], { type: "text/csv" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "Staff_Performance_Score.csv";
  a.click();
  URL.revokeObjectURL(url);
};
```

### Cascading Filter State
```ts
const [reviewPeriodId, setReviewPeriodId] = useState("");
const [departmentId, setDepartmentId] = useState<number>(0);
const [divisionId, setDivisionId] = useState<number>(0);
const [officeId, setOfficeId] = useState<number>(0);
const [departments, setDepartments] = useState<Department[]>([]);
const [divisions, setDivisions] = useState<Division[]>([]);
const [offices, setOffices] = useState<Office[]>([]);
const [periods, setPeriods] = useState<PerformanceReviewPeriod[]>([]);
const [scores, setScores] = useState<PeriodScoreData[]>([]);
const [filteredScores, setFilteredScores] = useState<PeriodScoreData[]>([]);

// On department change:
const handleDeptChange = async (deptId: number) => {
  setDepartmentId(deptId);
  setDivisionId(0);
  setOfficeId(0);
  setOffices([]);
  if (deptId > 0) {
    const res = await getDivisionsByDepartment(deptId);
    setDivisions(res?.data ? (Array.isArray(res.data) ? res.data : []) : []);
  } else {
    setDivisions([]);
  }
};

// On division change:
const handleDivChange = async (divId: number) => {
  setDivisionId(divId);
  setOfficeId(0);
  if (divId > 0) {
    const res = await getOfficesByDivision(divId);
    setOffices(res?.data ? (Array.isArray(res.data) ? res.data : []) : []);
  } else {
    setOffices([]);
  }
};
```

---

## Verification
```bash
cd /Users/enaira/Desktop/ENTERPRISE\ PMS/go-pms/pms-web && npm run build
```
Expect zero TypeScript/ESLint errors. All 5 routes should render.

## Files Created Summary
| File | Route |
|------|-------|
| `src/app/(dashboard)/staff-review-periods/page.tsx` | `/staff-review-periods` |
| `src/app/(dashboard)/review-period-report/page.tsx` | `/review-period-report` |
| `src/app/(dashboard)/review-period-report/[reviewPeriodId]/objectives/page.tsx` | `/review-period-report/:id/objectives` |
| `src/app/(dashboard)/review-period-report/[reviewPeriodId]/work-products/page.tsx` | `/review-period-report/:id/work-products` |
| `src/app/(dashboard)/pms-period-scores-report/page.tsx` | `/pms-period-scores-report` |
