# PROMPT 9E: Competency Report Pages (3 Pages)

## Context
You are converting a .NET Blazor PMS application to Next.js 15 (App Router) + TypeScript + Tailwind CSS + shadcn/ui + Recharts. This prompt covers the 3 competency report pages under the `(dashboard)` route group.

**PREREQUISITE**: PROMPT 9A must be completed first (types including `CompetencyMatrixReviewOverview`, `GroupedCompetencyReviewProfile`, `ChartDataVm`, API functions, chart components).

## Target Directory
`/Users/enaira/Desktop/ENTERPRISE PMS/go-pms/pms-web/`

## Existing API Functions (DO NOT recreate)
```ts
// competency.ts (already has all needed functions)
getCompetencyMatrixReviewProfiles(params: { reviewPeriodId?: number; officeId?: number; divisionId?: number; departmentId?: number })
  → CompetencyReviewProfile[]
getTechnicalCompetencyMatrixReviewProfiles(reviewPeriodId: number, jobRoleId: number)
  → CompetencyReviewProfile[]
getGroupCompetencyReviewProfiles(params: { reviewPeriodId?: number; officeId?: number; divisionId?: number; departmentId?: number })
  → CompetencyReviewProfile[]
getCompetencyReviewPeriods() → CompetencyReviewPeriod[]
getJobRoles() → JobRole[]

// organogram.ts
getDepartments(), getDivisionsByDepartment(departmentId), getOfficesByDivision(divisionId)
```

## Key Types

### From `src/types/dashboard.ts` (added in PROMPT 9A)
```ts
interface CompetencyMatrixDetail {
  competencyName: string;
  averageScore: number;
  expectedRatingValue: number;
}
interface CompetencyMatrixReviewProfile {
  employeeId: string; employeeName: string; position: string; grade: string;
  officeName: string; divisionName: string; departmentName: string;
  noOfCompetencies: number; noOfCompetent: number; gapCount: number;
  overallAverage: number; competencyMatrixDetails: CompetencyMatrixDetail[];
}
interface CompetencyMatrixReviewOverview {
  competencyNames: string[];
  competencyMatrixReviewProfiles: CompetencyMatrixReviewProfile[];
}
interface CompetencyRatingStat { ratingOrder: number; ratingName: string; numberOfStaff: number; staffPercentage: number; }
interface CategoryCompetencyDetailStat {
  categoryName: string; competencyRatingStat: CompetencyRatingStat[];
  averageRating: number; highestRating: number; lowestRating: number; mostCommonRating: number;
  groupCompetencyRatings: ChartDataVm[];
}
interface CategoryCompetencyStat { categoryName: string; actual: number; expected: number; }
interface GroupedCompetencyReviewProfile {
  categoryCompetencyStats: CategoryCompetencyStat[];
  categoryCompetencyDetailStats: CategoryCompetencyDetailStat[];
}
interface ChartDataVm { label: string; actual: number; expected: number; }
```

### From `src/types/competency.ts` (already exists)
```ts
interface CompetencyReviewPeriod { ... reviewPeriodId: number; name: string; ... }
interface JobRole { ... jobRoleId: number; name: string; ... }
```

## Important Note on API Response Shape
The existing competency.ts API functions return `BaseAPIResponse<CompetencyReviewProfile[]>`, but the .NET source shows the matrix reports return a special response with `competencyNames` + `competencyMatrixReviewProfiles`. Check if the existing `getCompetencyMatrixReviewProfiles()` function needs its return type adjusted to `BaseAPIResponse<CompetencyMatrixReviewOverview>` instead:

```ts
// In competency.ts, may need to update these return types:
export const getCompetencyMatrixReviewProfiles = (params: { ... }) => {
  // ...
  return get<BaseAPIResponse<CompetencyMatrixReviewOverview>>(`/competency/review-profiles/matrix?${q}`);
};
export const getTechnicalCompetencyMatrixReviewProfiles = (reviewPeriodId: number, jobRoleId: number) =>
  get<BaseAPIResponse<CompetencyMatrixReviewOverview>>(
    `/competency/review-profiles/technical-matrix?reviewPeriodId=${reviewPeriodId}&jobRoleId=${jobRoleId}`
  );
export const getGroupCompetencyReviewProfiles = (params: { ... }) => {
  // ...
  return get<BaseAPIResponse<GroupedCompetencyReviewProfile>>(`/competency/review-profiles/group?${q}`);
};
```

Check the actual Go API response types and adjust accordingly. The matrix endpoints should return `competencyNames[]` alongside the profile data.

---

## Page 1: `src/app/(dashboard)/competency_matrix/page.tsx`

**Route**: `/competency_matrix`
**Source**: `CompetencyMatrixReportPage.razor` + `.razor.cs`
**Access**: Admin, Head roles, HrReportAdmin
**Sidebar**: Competency Reports → Behavioral Matrix

### Behavior (from .NET source)
1. On load:
   - Call `getCompetencyReviewPeriods()` → populate Review Period dropdown
   - Call `getDepartments()` → populate Department dropdown
2. **Cascading filters**: Department → Division → Office (same pattern as PROMPT 9C)
3. On "Search" click:
   - Call `getCompetencyMatrixReviewProfiles({ reviewPeriodId, officeId, divisionId, departmentId })`
   - Response contains `competencyNames: string[]` (dynamic column headers) + `competencyMatrixReviewProfiles: CompetencyMatrixReviewProfile[]`
4. Render a wide table with fixed columns + dynamic competency columns

### UI Structure
- **PageHeader**: title="Behavioral Competency Matrix", breadcrumbs: Home → Behavioral Competency Matrix
- **Filter Row**: Review Period, Department (cascading), Division (cascading), Office (cascading), Search button

- **Legend sidebar** (col-md-3 on the right):
  - Green badge = Matches (actual == expected)
  - Yellow/Amber badge = Gap (actual < expected)
  - Blue badge = Exceeds (actual > expected)

- **DataTable** (col-md-9, horizontal scroll):
  - **Fixed Columns**: Employee ID, Name, Role, Grade, Office, Division, Department, Total Competencies (badge), Competent Count (badge), Gaps (badge), Overall Average
  - **Dynamic Columns**: One column per competency name from `competencyNames[]`
    - Cell value: `{averageScore}/{expectedRatingValue}` with color-coded badge
    - Color logic:
      - `expected - actual > 0` → Warning/Amber badge (gap)
      - `expected - actual < 0` → Primary/Blue badge (exceeds)
      - `expected === actual` → Success/Green badge (matches)
  - Sort by Employee Name by default
  - Search by employee ID or name

### Dynamic Column Generation
Since DataTable uses TanStack React Table ColumnDef, generate columns dynamically:
```tsx
const baseColumns: ColumnDef<CompetencyMatrixReviewProfile>[] = [
  { accessorKey: "employeeId", header: "Employee ID" },
  { accessorKey: "employeeName", header: "Name" },
  { accessorKey: "position", header: "Role" },
  { accessorKey: "grade", header: "Grade" },
  { accessorKey: "officeName", header: "Office" },
  { accessorKey: "divisionName", header: "Division" },
  { accessorKey: "departmentName", header: "Department" },
  {
    accessorKey: "noOfCompetencies",
    header: "Total",
    cell: ({ row }) => <Badge variant="outline">{row.original.noOfCompetencies}</Badge>,
  },
  {
    accessorKey: "noOfCompetent",
    header: "Competent",
    cell: ({ row }) => <Badge variant="outline" className="border-green-500">{row.original.noOfCompetent}</Badge>,
  },
  {
    accessorKey: "gapCount",
    header: "Gaps",
    cell: ({ row }) => <Badge variant="outline" className="border-amber-500">{row.original.gapCount}</Badge>,
  },
  { accessorKey: "overallAverage", header: "Avg" },
];

// Dynamic competency columns:
const competencyColumns: ColumnDef<CompetencyMatrixReviewProfile>[] = (matrixData?.competencyNames || []).map((name, idx) => ({
  id: `competency_${idx}`,
  header: name,
  cell: ({ row }) => {
    const detail = row.original.competencyMatrixDetails?.[idx];
    if (!detail) return "—";
    const diff = detail.expectedRatingValue - detail.averageScore;
    let badgeClass = "";
    if (diff > 0) badgeClass = "bg-amber-100 text-amber-800 border-amber-300"; // gap
    else if (diff < 0) badgeClass = "bg-blue-100 text-blue-800 border-blue-300"; // exceeds
    else badgeClass = "bg-green-100 text-green-800 border-green-300"; // matches
    return (
      <Badge variant="outline" className={badgeClass}>
        <b>{detail.averageScore}</b>/<b>{detail.expectedRatingValue}</b>
      </Badge>
    );
  },
}));

const allColumns = [...baseColumns, ...competencyColumns];
```

### Table Wrapper
Since this table can be very wide, wrap the DataTable in a horizontally scrollable container:
```tsx
<div className="overflow-x-auto">
  <DataTable columns={allColumns} data={matrixData?.competencyMatrixReviewProfiles || []} searchKey="employeeName" />
</div>
```

---

## Page 2: `src/app/(dashboard)/technical_competency_matrix/page.tsx`

**Route**: `/technical_competency_matrix`
**Source**: `TechnicalCompetencyMatrixPage.razor` + `.razor.cs`
**Access**: Admin, Head roles, HrReportAdmin
**Sidebar**: Competency Reports → Technical Matrix

### Behavior (from .NET source)
1. On load:
   - Call `getCompetencyReviewPeriods()` → populate Review Period dropdown
   - Call `getJobRoles()` → populate Job Role dropdown
2. **No cascading** — only 2 dropdowns (Review Period + Job Role)
3. On "Search" click:
   - Call `getTechnicalCompetencyMatrixReviewProfiles(reviewPeriodId, jobRoleId)`
   - Response: same shape as behavioral matrix (`competencyNames[]` + `competencyMatrixReviewProfiles[]`)
4. Render identical matrix table as Page 1

### UI Structure
- **PageHeader**: title="Technical Competency Matrix", breadcrumbs: Home → Technical Competency Matrix
- **Filter Row**: Review Period (`<Select>`), Job Role (`<Select>`), Search button
- **Legend** (same as behavioral matrix)
- **DataTable**: Identical column structure to behavioral matrix (same dynamic columns, same color-coded badges)

### Differences from Behavioral Matrix
- Filters: Job Role instead of Department/Division/Office
- API call: `getTechnicalCompetencyMatrixReviewProfiles(reviewPeriodId, jobRoleId)` instead of `getCompetencyMatrixReviewProfiles(params)`
- Everything else (table columns, badge logic, dynamic columns) is identical

### Code Reuse
Since the matrix table is identical, consider extracting the matrix table into a shared component if desired. Otherwise, duplicate the column generation logic.

---

## Page 3: `src/app/(dashboard)/group_profiles/page.tsx`

**Route**: `/group_profiles`
**Source**: `GroupReportPage.razor` + `.razor.cs`
**Access**: Admin, Head roles, HrReportAdmin
**Sidebar**: Competency Reports → Group Report

### Behavior (from .NET source)
1. On load:
   - Get logged-in user's employee details → use their `officeId` for initial data
   - Call `getGroupCompetencyReviewProfiles({ officeId: employee.officeId })` → initial load
   - Call `getCompetencyReviewPeriods()` → populate Review Period dropdown
   - Call `getDepartments()` → populate Department dropdown
2. **Cascading filters**: Department → Division → Office
3. On "Search" click:
   - Call `getGroupCompetencyReviewProfiles({ reviewPeriodId, officeId, divisionId, departmentId })`
4. Response shape: `GroupedCompetencyReviewProfile` with:
   - `categoryCompetencyStats[]` → for summary bar chart
   - `categoryCompetencyDetailStats[]` → for expandable detail sections

### UI Structure
- **PageHeader**: title="Group Strengths and Development Needs", breadcrumbs: Home → Group Profiles
- **Layout**: 2 columns — main content (col-9) + legend sidebar (col-3)

**Main Content (left)**:
- **Filter Row**: Review Period, Department (cascading), Division (cascading), Office (cascading), Search button

- **Summary Section** (Collapsible/Accordion, initially expanded):
  - Title: "Summary"
  - **Bar Chart** (Recharts BarChart) showing overview of competency categories:
    - X-axis: Category names
    - Two bars per category: "Actual" (colored) + "Expected" (gray)
    - Data from `categoryCompetencyStats`:
      ```ts
      const summaryChartData = profile.categoryCompetencyStats.map(s => ({
        name: s.categoryName,
        actual: s.actual,
        expected: s.expected,
      }));
      ```

- **Detail Analysis Section** (Collapsible, initially expanded):
  - Title: "Competency Group Detail Analysis"
  - **Per-category expandable panels**: For each `CategoryCompetencyDetailStat`:

    **Left side (col-7)**: Rating statistics table
    | Rating | Number | Percentage |
    |--------|--------|------------|
    | 5-Expert | 3 | 15% |
    | 4-Advanced | 8 | 40% |
    | ... | ... | ... |

    **Right side (col-4)**: Category metrics table
    | Metric | Value |
    |--------|-------|
    | Average Proficiency | 3.45 |
    | Highest Proficiency | 5.00 |
    | Lowest Proficiency | 1.00 |
    | Common Proficiency | 3.00 |

    **Below both**: Per-category bar chart showing competencies from highest to lowest
    - Data from `categoryCharts.groupCompetencyRatings`:
      ```ts
      const categoryChartData = detailStat.groupCompetencyRatings.map(r => ({
        name: r.label,
        actual: r.actual,
        expected: r.expected,
      }));
      ```

**Legend Sidebar (right)**:
- Card with title "Legend"
- Items:
  - Green badge: "Matches"
  - Amber/Warning badge: "Current Proficiency" (gap)
  - Gray/Secondary badge: "Expected Proficiency"

### Summary Bar Chart Component
```tsx
<ResponsiveContainer width="100%" height={300}>
  <BarChart data={summaryChartData}>
    <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
    <XAxis dataKey="name" tick={{ fill: "hsl(var(--muted-foreground))", fontSize: 11 }} />
    <YAxis tick={{ fill: "hsl(var(--muted-foreground))", fontSize: 12 }} />
    <Tooltip contentStyle={{ backgroundColor: "hsl(var(--popover))", border: "1px solid hsl(var(--border))", borderRadius: "6px" }} />
    <Legend />
    <Bar dataKey="actual" name="Current Proficiency" fill="hsl(var(--chart-1))" radius={[4, 4, 0, 0]} />
    <Bar dataKey="expected" name="Expected Proficiency" fill="hsl(var(--chart-3))" radius={[4, 4, 0, 0]} />
  </BarChart>
</ResponsiveContainer>
```

### Expandable Panels
Use shadcn `Collapsible` or `Accordion` components:
```tsx
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
// OR
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
```

### Detail Section Per Category
```tsx
{profile.categoryCompetencyDetailStats?.map((cat, idx) => (
  <AccordionItem key={idx} value={`cat-${idx}`}>
    <AccordionTrigger className="font-semibold">{cat.categoryName}</AccordionTrigger>
    <AccordionContent>
      <div className="grid gap-4 md:grid-cols-12">
        {/* Rating stats table - col-7 */}
        <div className="md:col-span-7">
          <table className="w-full text-sm">
            <thead>
              <tr><th></th><th className="text-center">Number</th><th className="text-center">Percentage</th></tr>
            </thead>
            <tbody>
              {cat.competencyRatingStat
                .sort((a, b) => b.ratingOrder - a.ratingOrder)
                .map((r, ri) => (
                  <tr key={ri}>
                    <td>{r.ratingOrder}-{r.ratingName}</td>
                    <td className="text-center">{r.numberOfStaff}</td>
                    <td className="text-center">{r.staffPercentage}%</td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
        {/* Category metrics - col-4 (with col-1 divider) */}
        <div className="md:col-span-1 hidden md:flex items-center justify-center">
          <div className="h-full w-px bg-border" />
        </div>
        <div className="md:col-span-4">
          <table className="w-full text-sm">
            <tbody>
              <tr><td>Average Proficiency</td><td>{cat.averageRating.toFixed(2)}</td></tr>
              <tr><td>Highest Proficiency</td><td>{cat.highestRating.toFixed(2)}</td></tr>
              <tr><td>Lowest Proficiency</td><td>{cat.lowestRating.toFixed(2)}</td></tr>
              <tr><td>Common Proficiency</td><td>{cat.mostCommonRating.toFixed(2)}</td></tr>
            </tbody>
          </table>
        </div>
      </div>
      {/* Per-category bar chart */}
      <div className="mt-4">
        <ResponsiveContainer width="100%" height={250}>
          <BarChart data={cat.groupCompetencyRatings.map(r => ({ name: r.label, actual: r.actual, expected: r.expected }))}>
            <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
            <XAxis dataKey="name" tick={{ fontSize: 10 }} />
            <YAxis />
            <Tooltip />
            <Legend />
            <Bar dataKey="actual" name="Current" fill="hsl(var(--chart-1))" radius={[4, 4, 0, 0]} />
            <Bar dataKey="expected" name="Expected" fill="hsl(var(--chart-3))" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </AccordionContent>
  </AccordionItem>
))}
```

---

## Verification
```bash
cd /Users/enaira/Desktop/ENTERPRISE\ PMS/go-pms/pms-web && npm run build
```
Expect zero TypeScript/ESLint errors. All 3 routes should render.

## Files Created/Modified Summary
| Action | File | Route |
|--------|------|-------|
| MODIFY | `src/lib/api/competency.ts` | (update return types for matrix/group endpoints if needed) |
| CREATE | `src/app/(dashboard)/competency_matrix/page.tsx` | `/competency_matrix` |
| CREATE | `src/app/(dashboard)/technical_competency_matrix/page.tsx` | `/technical_competency_matrix` |
| CREATE | `src/app/(dashboard)/group_profiles/page.tsx` | `/group_profiles` |

## Final Verification (After All PROMPT 9A-9E Complete)
```bash
cd /Users/enaira/Desktop/ENTERPRISE\ PMS/go-pms/pms-web && npm run build
```
All 17 new pages should compile with zero errors. The sidebar navigation already has links for all report pages. Scorecard pages are accessed via action buttons on the performance-score-card entry page.
