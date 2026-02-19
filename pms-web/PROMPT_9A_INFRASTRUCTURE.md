# PROMPT 9A: Infrastructure + Shared Chart Components

## Context
You are converting a .NET Blazor PMS application to Next.js 15 (App Router) + TypeScript + Tailwind CSS + shadcn/ui + Recharts. This prompt covers the foundational infrastructure needed before building the 17 scorecard/report pages.

## Target Directory
`/Users/enaira/Desktop/ENTERPRISE PMS/go-pms/pms-web/`

## Existing Architecture Patterns

### API Client Pattern
All API functions use thin wrappers around fetch:
```ts
import { get, post, put } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
```

### Chart Pattern (from existing dashboard page)
```ts
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from "recharts";
const PIE_COLORS = ["hsl(var(--chart-1))", "hsl(var(--chart-2))", "hsl(var(--chart-3))", "hsl(var(--chart-4))"];
```

### Card Pattern (shadcn)
```tsx
<Card>
  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
    <CardTitle className="text-sm font-medium">Title</CardTitle>
    <Icon className="h-4 w-4 text-muted-foreground" />
  </CardHeader>
  <CardContent>
    <div className="text-2xl font-bold">Value</div>
    <p className="text-xs text-muted-foreground">Description</p>
  </CardContent>
</Card>
```

---

## Task 1: Add Types to `src/types/dashboard.ts`

Add these interfaces AFTER the existing `EmployeeErpDetails` interface (line 166):

```ts
// --- Subordinates ScoreCard ---
export interface SubordinateScoreCard {
  staffId: string;
  staffName: string;
  totalWorkProducts: number;
  percentageWorkProductsCompletion: number;
  actualPoints: number;
  maxPoints: number;
  percentageScore: number;
  staffPerformanceGrade: string;
  totalWorkProductsCompletedOnSchedule: number;
  totalWorkProductsBehindSchedule: number;
}

export interface SubordinatesScoreCardResponse {
  isSuccess: boolean;
  message: string;
  managerId: string;
  reviewPeriodId: string;
  scoreCards: SubordinateScoreCard[];
}

// --- Period Score Data (for PMS Score Report) ---
export interface PeriodScoreData {
  periodScoreId: string;
  staffId: string;
  staffFullName: string;
  reviewPeriodId: string;
  reviewPeriod: string;
  year: number;
  startDate: string;
  endDate: string;
  finalScore: number;
  maxPoint: number;
  scorePercentage: number;
  finalGradeName: string;
  strategyName?: string;
  hrdDeductedPoints: number;
  departmentId?: number;
  departmentName?: string;
  divisionId?: number;
  divisionName?: string;
  officeId?: number;
  officeName?: string;
  staffGrade?: string;
  isUnderPerforming: boolean;
  minNoOfObjectives?: number;
  maxNoOfObjectives?: number;
  locationId?: string;
}

// --- Competency Group Report ---
export interface CompetencyRatingStat {
  ratingOrder: number;
  ratingName: string;
  numberOfStaff: number;
  staffPercentage: number;
}

export interface CategoryCompetencyDetailStat {
  categoryName: string;
  competencyRatingStat: CompetencyRatingStat[];
  averageRating: number;
  highestRating: number;
  lowestRating: number;
  mostCommonRating: number;
  groupCompetencyRatings: ChartDataVm[];
}

export interface CategoryCompetencyStat {
  categoryName: string;
  actual: number;
  expected: number;
}

export interface GroupedCompetencyReviewProfile {
  categoryCompetencyStats: CategoryCompetencyStat[];
  categoryCompetencyDetailStats: CategoryCompetencyDetailStat[];
}

export interface ChartDataVm {
  label: string;
  actual: number;
  expected: number;
}

// --- Competency Matrix ---
export interface CompetencyMatrixDetail {
  competencyName: string;
  averageScore: number;
  expectedRatingValue: number;
}

export interface CompetencyMatrixReviewProfile {
  employeeId: string;
  employeeName: string;
  position: string;
  grade: string;
  officeName: string;
  divisionName: string;
  departmentName: string;
  noOfCompetencies: number;
  noOfCompetent: number;
  gapCount: number;
  overallAverage: number;
  competencyMatrixDetails: CompetencyMatrixDetail[];
}

export interface CompetencyMatrixReviewOverview {
  competencyNames: string[];
  competencyMatrixReviewProfiles: CompetencyMatrixReviewProfile[];
}
```

---

## Task 2: Add API Functions to `src/lib/api/dashboard.ts`

Add these AFTER the existing `getStaffIdMask` function (line 63):

```ts
// --- Subordinates ScoreCard ---
export const getSubordinatesScoreCard = (managerId: string, reviewPeriodId: string) =>
  get<SubordinatesScoreCardResponse>(
    `/pms-engine/scorecard/subordinates?managerId=${managerId}&reviewPeriodId=${reviewPeriodId}`
  );

// --- Organogram Performance List ---
export const getOrganogramPerformanceSummaryList = (
  headOfUnitId: string,
  reviewPeriodId: string,
  level?: number
) =>
  get<BaseAPIResponse<OrganogramPerformanceSummary[]>>(
    `/pms-engine/organogram-performance/list?headOfUnitId=${headOfUnitId}&reviewPeriodId=${reviewPeriodId}${level !== undefined ? `&level=${level}` : ""}`
  );

// --- Period Scores (for PMS Score Report) ---
export const getPeriodScores = (reviewPeriodId: string) =>
  get<BaseAPIResponse<PeriodScoreData[]>>(
    `/pms-engine/period-scores/all?reviewPeriodId=${reviewPeriodId}`
  );

export const getPeriodScoreDetails = (reviewPeriodId: string, staffId: string) =>
  get<BaseAPIResponse<PeriodScoreData>>(
    `/pms-engine/period-scores?reviewPeriodId=${reviewPeriodId}&staffId=${staffId}`
  );

// --- Head Subordinates ---
export const getHeadSubordinates = (staffId: string) =>
  get<BaseAPIResponse<unknown[]>>(
    `/pms-engine/line-manager-employees?staffId=${staffId}&category=direct`
  );

// --- Staff Review Periods ---
export const getStaffReviewPeriods = (staffId: string) =>
  get<BaseAPIResponse<import("@/types/performance").PerformanceReviewPeriod[]>>(
    `/review-periods/staff?staffId=${staffId}`
  );
```

Also update the imports at top to include the new types:
```ts
import type {
  FeedbackRequestDashboardStats,
  PerformancePointsStats,
  WorkProductDashboardStats,
  WorkProductDetailsDashboardStats,
  StaffScoreCardResponse,
  StaffAnnualScoreCardResponse,
  OrganogramPerformanceSummary,
  PendingAction,
  EmployeeErpDetails,
  SubordinatesScoreCardResponse,
  PeriodScoreData,
} from "@/types/dashboard";
```

---

## Task 3: Create `src/lib/scorecard-helpers.ts`

```ts
// Performance grade display + color mapping
// Grade boundaries from .NET: <50=Developing, 50-65=Progressive, 66-79=Competent, 80-89=Accomplished, 90-100=Exemplary

export type PerformanceGrade = "Developing" | "Progressive" | "Competent" | "Accomplished" | "Exemplary";

export function getGradeFromScore(score: number): PerformanceGrade {
  if (score < 50) return "Developing";
  if (score < 66) return "Progressive";
  if (score < 80) return "Competent";
  if (score < 90) return "Accomplished";
  return "Exemplary";
}

export function getGradeInfo(grade: string): { label: string; color: string; bgClass: string } {
  switch (grade?.toLowerCase()) {
    case "developing":
      return { label: "Developing", color: "#ef4444", bgClass: "bg-red-100 text-red-800" };
    case "progressive":
      return { label: "Progressive", color: "#f97316", bgClass: "bg-orange-100 text-orange-800" };
    case "competent":
      return { label: "Competent", color: "#22c55e", bgClass: "bg-green-100 text-green-800" };
    case "accomplished":
      return { label: "Accomplished", color: "#3b82f6", bgClass: "bg-blue-100 text-blue-800" };
    case "exemplary":
      return { label: "Exemplary", color: "#8b5cf6", bgClass: "bg-purple-100 text-purple-800" };
    default:
      return { label: grade || "N/A", color: "#6b7280", bgClass: "bg-gray-100 text-gray-800" };
  }
}

export function getGradeNumericValue(grade: string): number {
  switch (grade?.toLowerCase()) {
    case "developing": return 1;
    case "progressive": return 2;
    case "competent": return 3;
    case "accomplished": return 4;
    case "exemplary": return 5;
    default: return 0;
  }
}

export const CHART_COLORS = {
  primary: "hsl(var(--chart-1))",
  secondary: "hsl(var(--chart-2))",
  tertiary: "hsl(var(--chart-3))",
  quaternary: "hsl(var(--chart-4))",
  earned: "#36A2EB",
  remaining: "#1b171842",
  gap: "#a12b2b",
  closed: "#0e293b",
  success: "#22c55e",
  warning: "#f59e0b",
  danger: "#ef4444",
  info: "#3b82f6",
};

export function formatPercent(value: number, decimals = 2): string {
  return `${(Math.round(value * Math.pow(10, decimals)) / Math.pow(10, decimals)).toFixed(decimals)}%`;
}

export function formatPoints(value: number, decimals = 4): string {
  return (Math.round(value * Math.pow(10, decimals)) / Math.pow(10, decimals)).toFixed(decimals);
}
```

---

## Task 4: Create Shared Chart Components

### 4a. `src/components/shared/charts/stat-card.tsx`
A reusable metric card for scorecard stat displays:
```tsx
"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { LucideIcon } from "lucide-react";

interface StatCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  iconColor?: string;
  description?: string;
  className?: string;
}

export function StatCard({ title, value, icon: Icon, iconColor, description, className }: StatCardProps) {
  return (
    <Card className={className}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="h-4 w-4" style={iconColor ? { color: iconColor } : undefined} />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {description && <p className="text-xs text-muted-foreground">{description}</p>}
      </CardContent>
    </Card>
  );
}
```

### 4b. `src/components/shared/charts/performance-doughnut.tsx`
A doughnut/ring chart showing a percentage with the value displayed in the center:
```tsx
"use client";

import { PieChart, Pie, Cell, ResponsiveContainer } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CHART_COLORS } from "@/lib/scorecard-helpers";

interface PerformanceDoughnutProps {
  title: string;
  percentage: number;
  earnedColor?: string;
  remainingColor?: string;
  className?: string;
}

export function PerformanceDoughnut({
  title,
  percentage,
  earnedColor = CHART_COLORS.earned,
  remainingColor = CHART_COLORS.remaining,
  className,
}: PerformanceDoughnutProps) {
  const data = [
    { name: "Earned", value: percentage },
    { name: "Remaining", value: Math.max(0, 100 - percentage) },
  ];

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle className="text-base">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="relative flex items-center justify-center" style={{ height: 250 }}>
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie data={data} cx="50%" cy="50%" innerRadius={60} outerRadius={90} dataKey="value" startAngle={90} endAngle={-270}>
                <Cell fill={earnedColor} />
                <Cell fill={remainingColor} />
              </Pie>
            </PieChart>
          </ResponsiveContainer>
          <div className="absolute text-3xl font-bold">{percentage.toFixed(1)}%</div>
        </div>
      </CardContent>
    </Card>
  );
}
```

### 4c. `src/components/shared/charts/competency-bar-chart.tsx`
A horizontal/vertical bar chart for competency ratings:
```tsx
"use client";

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface CompetencyBarChartProps {
  title: string;
  data: { name: string; value: number; color?: string }[];
  className?: string;
  height?: number;
}

const DEFAULT_COLORS = [
  "hsl(var(--chart-1))", "hsl(var(--chart-2))", "hsl(var(--chart-3))",
  "hsl(var(--chart-4))", "#8b5cf6", "#ec4899", "#14b8a6", "#f97316",
];

export function CompetencyBarChart({ title, data, className, height = 250 }: CompetencyBarChartProps) {
  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle className="text-base">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={height}>
          <BarChart data={data} layout="vertical" margin={{ left: 20 }}>
            <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
            <XAxis type="number" tick={{ fill: "hsl(var(--muted-foreground))", fontSize: 12 }} />
            <YAxis type="category" dataKey="name" width={120} tick={{ fill: "hsl(var(--muted-foreground))", fontSize: 11 }} />
            <Tooltip contentStyle={{ backgroundColor: "hsl(var(--popover))", border: "1px solid hsl(var(--border))", borderRadius: "6px" }} />
            <Bar dataKey="value" radius={[0, 4, 4, 0]}>
              {data.map((entry, index) => (
                <Cell key={index} fill={entry.color || DEFAULT_COLORS[index % DEFAULT_COLORS.length]} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
```

### 4d. `src/components/shared/charts/score-trend-chart.tsx`
A line chart for quarterly score trends:
```tsx
"use client";

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface ScoreTrendChartProps {
  title: string;
  data: { name: string; score: number }[];
  className?: string;
  height?: number;
}

export function ScoreTrendChart({ title, data, className, height = 250 }: ScoreTrendChartProps) {
  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle className="text-base">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={height}>
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
            <XAxis dataKey="name" tick={{ fill: "hsl(var(--muted-foreground))", fontSize: 12 }} />
            <YAxis domain={[0, 100]} tick={{ fill: "hsl(var(--muted-foreground))", fontSize: 12 }} />
            <Tooltip contentStyle={{ backgroundColor: "hsl(var(--popover))", border: "1px solid hsl(var(--border))", borderRadius: "6px" }} />
            <Line type="monotone" dataKey="score" stroke="hsl(var(--chart-1))" strokeWidth={2} dot={{ r: 4 }} activeDot={{ r: 6 }} />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
```

### 4e. `src/components/shared/charts/staff-profile-card.tsx`
A staff information card with photo, name, title, org hierarchy:
```tsx
"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { EmployeeErpDetails } from "@/types/dashboard";

interface StaffProfileCardProps {
  employee: EmployeeErpDetails;
  photoUrl?: string;
  className?: string;
}

export function StaffProfileCard({ employee, photoUrl, className }: StaffProfileCardProps) {
  const initials = `${employee.firstName?.[0] || ""}${employee.lastName?.[0] || ""}`.toUpperCase();

  return (
    <Card className={className}>
      <CardContent className="flex items-center gap-4 p-4">
        <Avatar className="h-16 w-16">
          <AvatarImage src={photoUrl || employee.photoUrl} alt={`${employee.firstName} ${employee.lastName}`} />
          <AvatarFallback className="text-lg">{initials}</AvatarFallback>
        </Avatar>
        <div className="flex-1 space-y-1">
          <h3 className="text-lg font-semibold">{employee.firstName} {employee.lastName}</h3>
          <p className="text-sm text-muted-foreground">{employee.jobTitle}</p>
          <p className="text-xs text-muted-foreground">
            {employee.officeName && <span>{employee.officeName}</span>}
            {employee.divisionName && <span> · {employee.divisionName}</span>}
            {employee.departmentName && <span> · {employee.departmentName}</span>}
          </p>
          {employee.gradeName && (
            <p className="text-xs text-muted-foreground">Grade: {employee.gradeName}</p>
          )}
          {employee.supervisorName && (
            <p className="text-xs text-muted-foreground">Supervisor: {employee.supervisorName}</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
```

---

## Verification
After implementing all files, run:
```bash
cd /Users/enaira/Desktop/ENTERPRISE\ PMS/go-pms/pms-web && npm run build
```
Expect zero TypeScript/ESLint errors.

## Files Created/Modified Summary
| Action | File |
|--------|------|
| MODIFY | `src/types/dashboard.ts` |
| MODIFY | `src/lib/api/dashboard.ts` |
| CREATE | `src/lib/scorecard-helpers.ts` |
| CREATE | `src/components/shared/charts/stat-card.tsx` |
| CREATE | `src/components/shared/charts/performance-doughnut.tsx` |
| CREATE | `src/components/shared/charts/competency-bar-chart.tsx` |
| CREATE | `src/components/shared/charts/score-trend-chart.tsx` |
| CREATE | `src/components/shared/charts/staff-profile-card.tsx` |
