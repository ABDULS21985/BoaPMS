# PROMPT 9B: ScoreCard Pages (5 Pages)

## Context
You are converting a .NET Blazor PMS application to Next.js 15 (App Router) + TypeScript + Tailwind CSS + shadcn/ui + Recharts. This prompt covers the 5 scorecard pages under the `(scorecard)` route group.

**PREREQUISITE**: PROMPT 9A must be completed first (types, API functions, helpers, chart components).

## Target Directory
`/Users/enaira/Desktop/ENTERPRISE PMS/go-pms/pms-web/`

## Route Group Layout
The `(scorecard)` route group uses `TopNavbar` only (no sidebar). Layout already exists at `src/app/(scorecard)/layout.tsx`:
```tsx
<div className="min-h-screen bg-background">
  <TopNavbar />
  <main className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">{children}</main>
  <Toaster richColors position="top-right" />
</div>
```

## Existing API Functions (DO NOT recreate)
```ts
// src/lib/api/dashboard.ts
getStaffScoreCard(staffId, reviewPeriodId) → StaffScoreCardResponse
getStaffAnnualScoreCard(staffId, year) → StaffAnnualScoreCardResponse
getOrganogramPerformanceSummary(referenceId, reviewPeriodId, level?) → OrganogramPerformanceSummary
getSubordinatesScoreCard(managerId, reviewPeriodId) → SubordinatesScoreCardResponse
getOrganogramPerformanceSummaryList(headOfUnitId, reviewPeriodId, level?) → OrganogramPerformanceSummary[]
getHeadSubordinates(staffId) → unknown[]
getEmployeeDetail(userId) → EmployeeErpDetails
getStaffIdMask(userId) → StaffIdMask
getActiveReviewPeriod() → PerformanceReviewPeriod

// src/lib/api/review-periods.ts
getReviewPeriods() → PerformanceReviewPeriod[]

// src/lib/api/pms-engine.ts
getStaffObjectives(staffId, reviewPeriodId) → IndividualPlannedObjective[]
```

## Existing Reusable Components
```ts
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { StatusBadge } from "@/components/shared/status-badge";
import { EmptyState } from "@/components/shared/empty-state";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { StatCard } from "@/components/shared/charts/stat-card";
import { PerformanceDoughnut } from "@/components/shared/charts/performance-doughnut";
import { CompetencyBarChart } from "@/components/shared/charts/competency-bar-chart";
import { ScoreTrendChart } from "@/components/shared/charts/score-trend-chart";
import { StaffProfileCard } from "@/components/shared/charts/staff-profile-card";
```

## Existing Helpers
```ts
import { getGradeFromScore, getGradeInfo, formatPercent, formatPoints, CHART_COLORS } from "@/lib/scorecard-helpers";
```

## Session Access
```ts
import { useSession } from "@/lib/auth-context"; // or wherever the session hook lives
const { user } = useSession();
// user.id = staffId, user.roles = string[]
```

---

## Page 1: `src/app/(scorecard)/performance-score-card/page.tsx`

**Route**: `/performance-score-card`
**Source**: `Performance.razor` + `Performance.razor.cs`
**Access**: All authenticated users

### Behavior (from .NET source)
1. On load, get logged-in user's `staffId` from session
2. Call `getEmployeeDetail(staffId)` — if null, redirect to `/404`
3. Call `getHeadSubordinates(staffId)` — if array has items, user `IsHeadOfUnit = true`
4. Call `getReviewPeriods()` to get all review periods
5. Call `getStaffObjectives(staffId, ...)` to get staff's planned objectives
6. Filter review periods to only show periods where staff has planned objectives, excluding statuses: Cancelled, PendingApproval, Rejected, Returned
7. If no periods found via objectives, fallback to `getStaffReviewPeriods(staffId)` from dashboard API

### UI Structure
- **Two tabs**: "My Performance" and "Team Performance" (only show Team tab if `IsHeadOfUnit`)
- **Tab 1 - My Performance**: DataTable of filtered review periods
  - Columns: Name, Year, Status (StatusBadge), Start Date, End Date, Actions
  - Actions: "View ScoreCard" button → `/staff-scorecard/[staffId]/[reviewPeriodId]`, "View Annual" button → `/staff-annual-scorecard/[staffId]/[year]`
  - Search by period name
- **Tab 2 - Team Performance** (conditional): DataTable of ALL review periods
  - Same columns but Action: "View Team" button → `/team-performance/[reviewPeriodId]/[managerId]`
  - `managerId` = logged-in user's staffId

### Key Details
- Use `useRouter()` for navigation
- Date formatting: `new Date(dateStr).toLocaleDateString()`
- Loading: Show `PageSkeleton` while data loads
- Empty: Show `EmptyState` if no periods

---

## Page 2: `src/app/(scorecard)/staff-scorecard/[staffId]/[reviewPeriodId]/page.tsx`

**Route**: `/staff-scorecard/[staffId]/[reviewPeriodId]`
**Source**: `StaffScoreCard.razor` + `StaffScoreCard.razor.cs`
**Access**: Self, supervisor, head of office/division/department, HrReportAdmin

### Behavior (from .NET source)
1. Extract `staffId` and `reviewPeriodId` from route params
2. Call `getEmployeeDetail(staffId)` — if null, redirect to `/404`
3. Call `getStaffIdMask(staffId)` for staff photo
4. **Access control**: Check if logged-in user can view (self, supervisor, head of unit, or HrReportAdmin role). If not, redirect to `/access-denied`
5. Call `getStaffScoreCard(staffId, reviewPeriodId)` — get scorecard data
6. Calculate: `percentageScore = 100 * actualPoints / maxPoints`

### UI Structure
- **StaffProfileCard** at top showing employee info
- **Page title**: `{Year} {PeriodName} - Score Card`
- **8 stat cards** in 2 rows of 4 (grid md:grid-cols-2 lg:grid-cols-4):
  1. Timely Work Products (`totalWorkProductsCompletedOnSchedule`) — CheckSquare icon, green
  2. Overdue Work Products (`totalWorkProductsBehindSchedule`) — AlertTriangle icon, amber
  3. Work Products Completion (`percentageWorkProductsCompletion`) — Percent icon, green
  4. Competency Gap Closure (`percentageGapsClosure`) — Percent icon, blue
  5. Accumulated Points (`accumulatedPoints`) — PlusCircle icon, navy
  6. Deducted Points (`deductedPoints`) — MinusCircle icon, red
  7. Actual Points (`actualPoints`) — CheckCircle icon, green
  8. Grade (`staffPerformanceGrade`) — Gauge icon, blue

- **Charts row** (grid md:grid-cols-2 lg:grid-cols-4):
  1. **Competency Bar Chart** (col-span-2): Horizontal bars from `pmsCompetencies` array — each competency name + ratingScore
  2. **Performance Doughnut** (col-span-1): Center shows `percentageScore%`
  3. **Competency Points** (col-span-1): Card with list of `pmsCompetencyCategory` entries (key=category name, value=points)

### Key Data Mapping
```ts
const scoreCard = response.scoreCard;
const percentageScore = scoreCard.maxPoints > 0
  ? Math.round((100 * scoreCard.actualPoints / scoreCard.maxPoints) * 100) / 100
  : 0;

// For competency bar chart:
const competencyData = (scoreCard.pmsCompetencies || []).map(c => ({
  name: c.pmsCompetency,
  value: c.ratingScore,
}));

// For competency category points:
const categoryPoints = Object.entries(scoreCard.pmsCompetencyCategory || {});
```

---

## Page 3: `src/app/(scorecard)/staff-annual-scorecard/[staffId]/[year]/page.tsx`

**Route**: `/staff-annual-scorecard/[staffId]/[year]`
**Source**: `StaffAnnualScoreCard.razor` + `StaffAnnualScoreCard.razor.cs`
**Access**: Same as staff scorecard (self, supervisor, head, HrReportAdmin)

### Behavior (from .NET source)
1. Extract `staffId` and `year` from route params
2. Same access control as Page 2
3. Call `getStaffAnnualScoreCard(staffId, year)` — returns `scoreCards: StaffScoreCardDetails[]` (one per period in the year)
4. Aggregate across all periods:
   - `totalWorkProducts = sum(scoreCards.map(s => s.totalWorkProducts))`
   - `percentageGapsClosure = sum(scoreCards.map(s => s.percentageGapsClosure))`
   - `actualPoints = sum(scoreCards.map(s => s.actualPoints))`
   - `totalPoints = sum(scoreCards.map(s => s.maxPoints))`
   - `percentageScore = 100 * actualPoints / totalPoints`
   - `grade = getGradeFromScore(percentageScore)` (use helper)
5. Aggregate competency categories: group by key, average the values
6. Aggregate competencies: group by name, average the ratingScores

### UI Structure
- **StaffProfileCard** at top
- **Title**: `{Year} Performance - Score Card`
- **4 stat cards** in grid (lg:grid-cols-4):
  1. Total Work Products
  2. Competency Gap Closure %
  3. Points Earned (`actualPoints/totalPoints`)
  4. Grade (from `getGradeFromScore`)

- **Charts** (2 rows of 2):
  - Row 1: **Performance Grade Bar Chart** (per period) + **Competency Bar Chart** (aggregated)
  - Row 2: **Score Trend Line Chart** (per period) + **Performance Doughnut** (overall %) + **Competency Points** card

### Performance Grade Bar Chart
For each scoreCard in the array, map `staffPerformanceGrade` to numeric value (1-5) using `getGradeNumericValue()`:
```ts
const gradeData = scoreCards.map(sc => ({
  name: sc.reviewPeriodShortName || sc.reviewPeriod,
  value: getGradeNumericValue(sc.staffPerformanceGrade),
}));
```

### Score Trend Line Chart
```ts
const trendData = scoreCards.map(sc => ({
  name: sc.reviewPeriodShortName || sc.reviewPeriod,
  score: Math.round(sc.percentageScore * 100) / 100,
}));
```

---

## Page 4: `src/app/(scorecard)/team-performance/[reviewPeriodId]/[managerId]/page.tsx`

**Route**: `/team-performance/[reviewPeriodId]/[managerId]`
**Source**: `TeamPerformances.razor` + `TeamPerformances.razor.cs` (808 lines — most complex scorecard page)
**Access**: Supervisor, Head of Office/Division/Department, HrReportAdmin

### Behavior (from .NET source)
1. Extract `reviewPeriodId` and `managerId` from params
2. Call `getEmployeeDetail(managerId)` to get manager's org details
3. Determine manager's role level:
   - If HrReportAdmin → `canViewAll = true`, show all units
   - If `headOfDeptId === managerId` → `isHeadOfDepartment = true` (level 2)
   - If `headOfDivId === managerId` → `isHeadOfDivision = true` (level 3)
   - If `headOfOfficeId === managerId` → `isHeadOfOffice = true` (level 4)
   - Else → `isSupervisorOnly = true` (direct reports only)
4. Call `getSubordinatesScoreCard(managerId, reviewPeriodId)` — get staff performance list
5. Based on level, also call `getOrganogramPerformanceSummaryList(...)` for sub-unit summaries

### UI Structure (Tabs vary by role level)
**Tab 1 - Staff Performance** (always shown): DataTable of subordinate staff
- Columns: #, Staff Name, Work Products, Completion Rate (progress bar), Score (progress bar), Grade (badge)
- Progress bars: Use `<Progress>` component, color based on grade
- Actions: "View ScoreCard" → `/staff-scorecard/[staffId]/[reviewPeriodId]`, "View Annual" → `/staff-annual-scorecard/[staffId]/[year]`

**Tab 2 - Org Unit Tab** (conditional):
- **HeadOfDepartment**: "Divisions" tab → DataTable of division performance summaries
  - Columns: Division, Head, WPs, Performance %, Score, Grade
  - Actions: "Drill Down" → `/team-performance/[reviewPeriodId]/[divisionHeadId]`, "Unit ScoreCard" → `/unit-scorecard/[divisionId]/[reviewPeriodId]/3`
- **HeadOfDivision**: "Offices" tab → DataTable of office performance summaries
  - Actions: "Drill Down" → `/team-performance/[reviewPeriodId]/[officeHeadId]`, "Unit ScoreCard" → `/unit-scorecard/[officeId]/[reviewPeriodId]/4`
- **HrReportAdmin**: "Departments" tab showing all department summaries
  - Actions: "Drill Down" → `/team-performance/[reviewPeriodId]/[deptHeadId]`, "Unit ScoreCard" → `/unit-scorecard/[deptId]/[reviewPeriodId]/2`

### Progress Bar Color Logic
```ts
function getProgressColor(score: number): string {
  if (score >= 90) return "bg-purple-500";
  if (score >= 80) return "bg-blue-500";
  if (score >= 66) return "bg-green-500";
  if (score >= 50) return "bg-orange-500";
  return "bg-red-500";
}
```

---

## Page 5: `src/app/(scorecard)/unit-scorecard/[unitId]/[reviewPeriodId]/[level]/page.tsx`

**Route**: `/unit-scorecard/[unitId]/[reviewPeriodId]/[level]`
**Source**: `UnitScoreCard.razor` + `UnitScoreCard.razor.cs`
**Access**: Head of unit, HrReportAdmin
**Level mapping**: 2=Department, 3=Division, 4=Office

### Behavior (from .NET source)
1. Extract `unitId`, `reviewPeriodId`, `level` from params
2. Call `getOrganogramPerformanceSummary(unitId, reviewPeriodId, level)` — returns single `OrganogramPerformanceSummary`
3. Get `managerId` from response, call `getEmployeeDetail(managerId)` for profile
4. Access control: Check logged-in user can view (same head/HrReportAdmin check)

### UI Structure
- **StaffProfileCard** of unit head
- **Title**: `{referenceName} {year} {reviewPeriod} - Score Card`

- **Highlighted stat** (full-width green gradient card):
  - Performance Score: `{performanceScore}%` + Performance Grade: `{earnedPerformanceGrade}`

- **6 stat cards** in grid:
  1. Points Earned (`actualScore/maxPoint`)
  2. Total Staff (`totalStaff`)
  3. Total Work Products (`totalWorkProducts`)
  4. Overdue Work Products (`totalWorkProductsBehindSchedule`)
  5. Timely Work Products (`totalWorkProductsCompletedOnSchedule`)
  6. Living the Values Rating (if available)

- **Charts row** (grid md:grid-cols-2):
  1. **360 Review Feedbacks** bar chart — 3 bars: Total (`total360Feedbacks`), Completed (`completed360FeedbacksToTreat`), Pending (`pending360FeedbacksToTreat`)
  2. **Gap Closure Doughnut** — `percentageGapsClosure` vs `100 - percentageGapsClosure`

### 360 Feedback Bar Chart Data
```ts
const feedbackData = [
  { name: "Total", value: summary.total360Feedbacks },
  { name: "Completed", value: summary.completed360FeedbacksToTreat },
  { name: "Pending", value: summary.pending360FeedbacksToTreat },
];
```

---

## Verification
```bash
cd /Users/enaira/Desktop/ENTERPRISE\ PMS/go-pms/pms-web && npm run build
```
Expect zero TypeScript/ESLint errors. All 5 scorecard routes should render.

## Files Created Summary
| File | Route |
|------|-------|
| `src/app/(scorecard)/performance-score-card/page.tsx` | `/performance-score-card` |
| `src/app/(scorecard)/staff-scorecard/[staffId]/[reviewPeriodId]/page.tsx` | `/staff-scorecard/:staffId/:reviewPeriodId` |
| `src/app/(scorecard)/staff-annual-scorecard/[staffId]/[year]/page.tsx` | `/staff-annual-scorecard/:staffId/:year` |
| `src/app/(scorecard)/team-performance/[reviewPeriodId]/[managerId]/page.tsx` | `/team-performance/:reviewPeriodId/:managerId` |
| `src/app/(scorecard)/unit-scorecard/[unitId]/[reviewPeriodId]/[level]/page.tsx` | `/unit-scorecard/:unitId/:reviewPeriodId/:level` |
