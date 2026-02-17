"use client";

import { useEffect, useState } from "react";
import { BarChart3 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getPeriodObjectiveEvaluations, getActiveReviewPeriod } from "@/lib/api/review-periods";
import type { PerformanceReviewPeriod } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

interface OutcomeEval {
  objectiveId: string;
  objectiveName?: string;
  objectiveLevel?: string;
  score?: number;
  weight?: number;
  weightedScore?: number;
  evaluatorName?: string;
  comment?: string;
}

export default function OutcomeEvaluationPage() {
  const [evaluations, setEvaluations] = useState<OutcomeEval[]>([]);
  const [period, setPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getActiveReviewPeriod()
      .then(async (periodRes) => {
        if (periodRes?.data) {
          setPeriod(periodRes.data);
          const evalRes = await getPeriodObjectiveEvaluations(periodRes.data.periodId);
          if (evalRes?.data) setEvaluations(Array.isArray(evalRes.data) ? (evalRes.data as OutcomeEval[]) : []);
        }
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const columns: ColumnDef<OutcomeEval>[] = [
    { accessorKey: "objectiveName", header: "Objective", cell: ({ row }) => <span className="font-medium">{row.original.objectiveName ?? row.original.objectiveId}</span> },
    { accessorKey: "objectiveLevel", header: "Level", cell: ({ row }) => <Badge variant="outline">{row.original.objectiveLevel}</Badge> },
    { accessorKey: "weight", header: "Weight", cell: ({ row }) => row.original.weight != null ? `${row.original.weight}%` : "—" },
    { accessorKey: "score", header: "Score", cell: ({ row }) => row.original.score?.toFixed(1) ?? "—" },
    { accessorKey: "weightedScore", header: "Weighted Score", cell: ({ row }) => row.original.weightedScore?.toFixed(2) ?? "—" },
    { accessorKey: "evaluatorName", header: "Evaluator" },
    { accessorKey: "comment", header: "Comment", cell: ({ row }) => <span className="line-clamp-1">{row.original.comment}</span> },
  ];

  if (loading) return <div><PageHeader title="Outcome Evaluation" breadcrumbs={[{ label: "Evaluation" }, { label: "Outcome" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Outcome Evaluation" description={period ? `Review Period: ${period.name}` : "No active review period"} breadcrumbs={[{ label: "Evaluation", href: "/evaluation/direct-reports" }, { label: "Outcome" }]} />
      {evaluations.length > 0 ? (
        <DataTable columns={columns} data={evaluations} searchKey="objectiveName" searchPlaceholder="Search evaluations..." />
      ) : (
        <EmptyState icon={BarChart3} title="No Evaluations" description="No outcome evaluations available for this review period." />
      )}
    </div>
  );
}
