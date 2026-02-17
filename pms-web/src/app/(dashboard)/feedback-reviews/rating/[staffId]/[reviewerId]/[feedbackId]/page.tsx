"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import { useSession } from "next-auth/react";
import { Loader2, CheckCircle, Circle, Send } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { Badge } from "@/components/ui/badge";
import {
  getQuestionnaire,
  getReviewerFeedbackDetails,
  getCompetencyReviewDetail,
  add360Rating,
  update360Rating,
  reviewerComplete360Review,
} from "@/lib/api/pms-engine";
import { getEmployeeDetail } from "@/lib/api/dashboard";
import type {
  PmsCompetency,
  CompetencyReviewer,
  CompetencyReviewerRating,
  FeedbackQuestionnaireWithOptions,
} from "@/types/performance";
import type { EmployeeErpDetails } from "@/types/dashboard";
import { cn } from "@/lib/utils";

export default function StaffRatingPage() {
  const params = useParams<{ staffId: string; reviewerId: string; feedbackId: string }>();
  const router = useRouter();
  const { data: session } = useSession();

  const [loading, setLoading] = useState(true);
  const [staffDetails, setStaffDetails] = useState<EmployeeErpDetails | null>(null);
  const [competencies, setCompetencies] = useState<PmsCompetency[]>([]);
  const [reviewerData, setReviewerData] = useState<CompetencyReviewer | null>(null);
  const [existingRatings, setExistingRatings] = useState<CompetencyReviewerRating[]>([]);
  const [selectedCompetencyId, setSelectedCompetencyId] = useState<string>("");
  const [selections, setSelections] = useState<Record<string, string>>({});
  const [saving, setSaving] = useState(false);
  const [finalizeOpen, setFinalizeOpen] = useState(false);
  const [isReviewComplete, setIsReviewComplete] = useState(false);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const [staffRes, questionnaireRes, reviewerRes] = await Promise.allSettled([
          getEmployeeDetail(params.staffId),
          getQuestionnaire(params.staffId),
          getReviewerFeedbackDetails(params.reviewerId),
        ]);

        if (staffRes.status === "fulfilled" && staffRes.value?.data) setStaffDetails(staffRes.value.data);
        if (questionnaireRes.status === "fulfilled" && questionnaireRes.value?.data) {
          const comps = Array.isArray(questionnaireRes.value.data) ? questionnaireRes.value.data : [];
          setCompetencies(comps);
          if (comps.length > 0) setSelectedCompetencyId(comps[0].pmsCompetencyId);
        }
        if (reviewerRes.status === "fulfilled" && reviewerRes.value?.data) {
          const rd = reviewerRes.value.data;
          setReviewerData(rd);
          const ratings = rd.competencyReviewerRatings ?? [];
          setExistingRatings(ratings);
          // Pre-populate selections from existing ratings
          const existing: Record<string, string> = {};
          for (const r of ratings) {
            if (r.pmsCompetencyId && r.feedbackQuestionaireOptionId) {
              existing[`${r.pmsCompetencyId}`] = r.feedbackQuestionaireOptionId;
            }
          }
          setSelections(existing);
          // Check if review already completed (status >= 10)
          if (rd.recordStatus != null && rd.recordStatus >= 10) {
            setIsReviewComplete(true);
          }
        }
      } catch { /* silent */ } finally {
        setLoading(false);
      }
    })();
  }, [params.staffId, params.reviewerId]);

  const selectedCompetency = competencies.find((c) => c.pmsCompetencyId === selectedCompetencyId);

  const isCompetencyRated = useCallback((compId: string) => {
    return existingRatings.some((r) => r.pmsCompetencyId === compId);
  }, [existingRatings]);

  const allCompetenciesRated = competencies.length > 0 && competencies.every((c) => isCompetencyRated(c.pmsCompetencyId));

  const handleSaveCompetencyRating = async () => {
    if (!selectedCompetency) return;
    const questions = selectedCompetency.feedbackQuestionaires ?? [];
    if (questions.length === 0) { toast.error("No questions for this competency."); return; }

    // For each question, find the selected option
    const unAnswered = questions.filter(
      (q) => !selections[`${selectedCompetency.pmsCompetencyId}_${q.feedbackQuestionaireId}`]
    );
    if (unAnswered.length > 0) {
      toast.error("Please answer all questions before saving.");
      return;
    }

    setSaving(true);
    try {
      let allOk = true;
      for (const q of questions) {
        const optionId = selections[`${selectedCompetency.pmsCompetencyId}_${q.feedbackQuestionaireId}`];
        const existingRating = existingRatings.find(
          (r) => r.pmsCompetencyId === selectedCompetency.pmsCompetencyId && r.feedbackQuestionaireOptionId
        );

        const payload = {
          pmsCompetencyId: selectedCompetency.pmsCompetencyId,
          feedbackQuestionaireOptionId: optionId,
          competencyReviewerId: params.reviewerId,
          competencyReviewerRatingId: existingRating?.competencyReviewerRatingId ?? "",
        };

        const res = existingRating
          ? await update360Rating(payload)
          : await add360Rating(payload);

        if (!res?.isSuccess) { allOk = false; break; }
      }

      if (allOk) {
        toast.success("Ratings saved successfully.");
        // Refresh reviewer data
        const refreshed = await getReviewerFeedbackDetails(params.reviewerId);
        if (refreshed?.data) {
          setReviewerData(refreshed.data);
          setExistingRatings(refreshed.data.competencyReviewerRatings ?? []);
        }
      } else {
        toast.error("Some ratings failed to save.");
      }
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const handleFinalize = async () => {
    try {
      const res = await reviewerComplete360Review({
        reviewStaffId: session?.user?.id ?? "",
        competencyReviewFeedbackId: params.feedbackId,
      });
      if (res?.isSuccess) {
        toast.success("360 Feedback finalized successfully.");
        setIsReviewComplete(true);
        router.push("/feedback-reviews");
      } else {
        toast.error(res?.message || "Finalization failed.");
      }
    } catch { toast.error("An error occurred."); }
  };

  if (loading) return <div><PageHeader title="Staff Rating" breadcrumbs={[{ label: "360 Feedback", href: "/feedback-reviews" }, { label: "Rating" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader
        title={`Rate: ${staffDetails?.firstName ?? ""} ${staffDetails?.lastName ?? params.staffId}`}
        description={staffDetails ? `${staffDetails.jobTitle ?? ""} - ${staffDetails.departmentName ?? ""}` : undefined}
        breadcrumbs={[
          { label: "360 Feedback", href: "/feedback-reviews" },
          { label: `Rating ${staffDetails?.firstName ?? params.staffId}` },
        ]}
        actions={
          allCompetenciesRated && !isReviewComplete ? (
            <Button onClick={() => setFinalizeOpen(true)}>
              <Send className="mr-2 h-4 w-4" />Finalize Feedback
            </Button>
          ) : undefined
        }
      />

      {isReviewComplete && (
        <div className="rounded-md border border-green-500/50 bg-green-50 p-4 dark:bg-green-950/20">
          <p className="text-sm font-medium text-green-700 dark:text-green-400">This feedback review has been completed and finalized.</p>
        </div>
      )}

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-4">
        {/* Competency Navigation */}
        <div className="space-y-2">
          <h3 className="text-sm font-semibold text-muted-foreground mb-3">Competencies</h3>
          {competencies.map((comp) => (
            <button
              key={comp.pmsCompetencyId}
              onClick={() => setSelectedCompetencyId(comp.pmsCompetencyId)}
              className={cn(
                "w-full flex items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition-colors",
                selectedCompetencyId === comp.pmsCompetencyId
                  ? "bg-primary text-primary-foreground"
                  : "hover:bg-muted"
              )}
            >
              {isCompetencyRated(comp.pmsCompetencyId) ? (
                <CheckCircle className="h-4 w-4 shrink-0 text-green-500" />
              ) : (
                <Circle className="h-4 w-4 shrink-0" />
              )}
              <span className="truncate">{comp.name}</span>
            </button>
          ))}
          <div className="mt-4 text-xs text-muted-foreground">
            {competencies.filter((c) => isCompetencyRated(c.pmsCompetencyId)).length} / {competencies.length} completed
          </div>
        </div>

        {/* Questions & Options */}
        <div className="lg:col-span-3">
          {selectedCompetency ? (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">{selectedCompetency.name}</CardTitle>
                {selectedCompetency.description && (
                  <p className="text-sm text-muted-foreground">{selectedCompetency.description}</p>
                )}
              </CardHeader>
              <CardContent className="space-y-6">
                {(selectedCompetency.feedbackQuestionaires ?? []).map((q: FeedbackQuestionnaireWithOptions, idx: number) => (
                  <div key={q.feedbackQuestionaireId} className="space-y-3">
                    <div>
                      <Label className="text-sm font-medium">
                        {idx + 1}. {q.question}
                      </Label>
                      {q.description && <p className="text-xs text-muted-foreground mt-0.5">{q.description}</p>}
                    </div>
                    <RadioGroup
                      value={selections[`${selectedCompetency.pmsCompetencyId}_${q.feedbackQuestionaireId}`] ?? ""}
                      onValueChange={(val) => {
                        setSelections((prev) => ({
                          ...prev,
                          [`${selectedCompetency.pmsCompetencyId}_${q.feedbackQuestionaireId}`]: val,
                        }));
                      }}
                      disabled={isReviewComplete}
                      className="space-y-2"
                    >
                      {(q.options ?? []).map((opt) => (
                        <div key={opt.feedbackQuestionaireOptionId} className="flex items-start gap-3 rounded-md border p-3">
                          <RadioGroupItem value={opt.feedbackQuestionaireOptionId} id={opt.feedbackQuestionaireOptionId} className="mt-0.5" />
                          <div className="flex-1">
                            <label htmlFor={opt.feedbackQuestionaireOptionId} className="text-sm cursor-pointer">
                              {opt.optionStatement}
                            </label>
                            {opt.description && <p className="text-xs text-muted-foreground">{opt.description}</p>}
                          </div>
                          <Badge variant="outline" className="shrink-0">{opt.score} pts</Badge>
                        </div>
                      ))}
                    </RadioGroup>
                  </div>
                ))}

                {!isReviewComplete && (selectedCompetency.feedbackQuestionaires ?? []).length > 0 && (
                  <div className="flex justify-end pt-4 border-t">
                    <Button onClick={handleSaveCompetencyRating} disabled={saving}>
                      {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                      {isCompetencyRated(selectedCompetency.pmsCompetencyId) ? "Update Rating" : "Save Rating"}
                    </Button>
                  </div>
                )}
              </CardContent>
            </Card>
          ) : (
            <p className="text-sm text-muted-foreground">Select a competency from the left panel.</p>
          )}
        </div>
      </div>

      <ConfirmationDialog
        open={finalizeOpen}
        onOpenChange={setFinalizeOpen}
        title="Finalize 360 Feedback"
        description="Are you sure you want to finalize your feedback? This action cannot be undone."
        confirmLabel="Finalize"
        onConfirm={handleFinalize}
      />
    </div>
  );
}
