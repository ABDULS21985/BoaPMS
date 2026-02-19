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
