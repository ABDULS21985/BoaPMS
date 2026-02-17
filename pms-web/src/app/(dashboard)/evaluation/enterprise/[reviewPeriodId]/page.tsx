"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Building2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { PageHeader } from "@/components/shared/page-header";
import { DataTable } from "@/components/shared/data-table";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { getPeriodObjectiveEvaluations, getReviewPeriodDetails } from "@/lib/api/review-periods";
import type { PerformanceReviewPeriod } from "@/types/performance";
import type { ColumnDef } from "@tanstack/react-table";

interface EnterpriseEval {
  objectiveId: string;
  objectiveName?: string;
  categoryName?: string;
  score?: number;
  weight?: number;
  weightedScore?: number;
  comment?: string;
}

export default function EnterpriseEvaluationPage() {
  const { reviewPeriodId } = useParams<{ reviewPeriodId: string }>();
  const [evaluations, setEvaluations] = useState<EnterpriseEval[]>([]);
  const [period, setPeriod] = useState<PerformanceReviewPeriod | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!reviewPeriodId) return;
    setLoading(true);
    Promise.all([getReviewPeriodDetails(reviewPeriodId), getPeriodObjectiveEvaluations(reviewPeriodId)])
      .then(([periodRes, evalRes]) => {
        if (periodRes?.data) setPeriod(periodRes.data);
        if (evalRes?.data) setEvaluations(Array.isArray(evalRes.data) ? (evalRes.data as EnterpriseEval[]) : []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [reviewPeriodId]);

  const columns: ColumnDef<EnterpriseEval>[] = [
    { accessorKey: "objectiveName", header: "Enterprise Objective", cell: ({ row }) => <span className="font-medium">{row.original.objectiveName ?? row.original.objectiveId}</span> },
    { accessorKey: "categoryName", header: "Category", cell: ({ row }) => <Badge variant="outline">{row.original.categoryName}</Badge> },
    { accessorKey: "weight", header: "Weight", cell: ({ row }) => row.original.weight != null ? `${row.original.weight}%` : "—" },
    { accessorKey: "score", header: "Score", cell: ({ row }) => row.original.score?.toFixed(1) ?? "—" },
    { accessorKey: "weightedScore", header: "Weighted", cell: ({ row }) => row.original.weightedScore?.toFixed(2) ?? "—" },
    { accessorKey: "comment", header: "Comment", cell: ({ row }) => <span className="line-clamp-1">{row.original.comment}</span> },
  ];

  if (loading) return <div><PageHeader title="Enterprise Objectives Evaluation" breadcrumbs={[{ label: "Evaluation" }, { label: "Enterprise" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Enterprise Objectives Evaluation" description={period ? `Review Period: ${period.name}` : ""} breadcrumbs={[{ label: "Evaluation", href: "/evaluation/direct-reports" }, { label: "Enterprise Outcome" }]} />
      {evaluations.length > 0 ? (
        <DataTable columns={columns} data={evaluations} searchKey="objectiveName" searchPlaceholder="Search objectives..." />
      ) : (
        <EmptyState icon={Building2} title="No Evaluations" description="No enterprise objective evaluations found for this period." />
      )}
    </div>
  );
}
