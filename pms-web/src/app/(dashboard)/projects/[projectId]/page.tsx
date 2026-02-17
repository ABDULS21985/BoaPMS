"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Send, Loader2, UserPlus, FolderKanban, Briefcase, Users } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Separator } from "@/components/ui/separator";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { StatusBadge } from "@/components/shared/status-badge";
import { FormSheet } from "@/components/shared/form-sheet";
import { ConfirmationDialog } from "@/components/shared/confirmation-dialog";
import { EmptyState } from "@/components/shared/empty-state";
import { getProjectDetails, getProjectMembers, getProjectWorkProducts, addProjectMember, submitDraftProject, approveProject, rejectProject } from "@/lib/api/pms-engine";
import type { Project, WorkProduct, ProjectMember } from "@/types/performance";
import { Status } from "@/types/enums";

export default function ProjectDetailPage() {
  const { projectId } = useParams<{ projectId: string }>();
  const router = useRouter();
  const [project, setProject] = useState<Project | null>(null);
  const [members, setMembers] = useState<ProjectMember[]>([]);
  const [workProducts, setWorkProducts] = useState<WorkProduct[]>([]);
  const [loading, setLoading] = useState(true);
  const [memberOpen, setMemberOpen] = useState(false);
  const [memberStaffId, setMemberStaffId] = useState("");
  const [saving, setSaving] = useState(false);
  const [actionType, setActionType] = useState<"submit" | "approve" | "reject" | null>(null);

  const loadData = async () => {
    if (!projectId) return;
    setLoading(true);
    try {
      const [projRes, memRes, wpRes] = await Promise.all([
        getProjectDetails(projectId),
        getProjectMembers(projectId),
        getProjectWorkProducts(projectId),
      ]);
      if (projRes?.data) setProject(projRes.data);
      if (memRes?.data) setMembers(Array.isArray(memRes.data) ? memRes.data : []);
      if (wpRes?.data) setWorkProducts(Array.isArray(wpRes.data) ? wpRes.data : []);
    } catch { /* */ } finally { setLoading(false); }
  };

  useEffect(() => { loadData(); }, [projectId]);

  const handleAddMember = async () => {
    if (!memberStaffId) { toast.error("Enter a staff ID."); return; }
    setSaving(true);
    try {
      const res = await addProjectMember({ projectId, staffId: memberStaffId });
      if (res?.isSuccess) { toast.success("Member added."); setMemberOpen(false); setMemberStaffId(""); loadData(); }
      else toast.error(res?.message || "Failed to add member.");
    } catch { toast.error("An error occurred."); } finally { setSaving(false); }
  };

  const handleAction = async () => {
    if (!project || !actionType) return;
    try {
      const payload = { id: project.projectId };
      const fn = actionType === "submit" ? submitDraftProject : actionType === "approve" ? approveProject : rejectProject;
      const res = await fn(payload);
      if (res?.isSuccess) { toast.success(`Project ${actionType}ed.`); loadData(); }
      else toast.error(res?.message || "Action failed.");
    } catch { toast.error("An error occurred."); }
  };

  if (loading) return <div><PageHeader title="Project Details" breadcrumbs={[{ label: "Projects" }]} /><PageSkeleton /></div>;
  if (!project) return <div><PageHeader title="Project Not Found" breadcrumbs={[{ label: "Projects" }]} /><EmptyState icon={FolderKanban} title="Not Found" description="Project could not be loaded." /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title={project.name} description={project.description} breadcrumbs={[{ label: "Projects", href: "/projects" }, { label: project.name }]} actions={
        <div className="flex gap-2">
          {project.recordStatus === Status.Draft && <Button size="sm" onClick={() => setActionType("submit")}><Send className="mr-2 h-4 w-4" />Submit</Button>}
          {project.recordStatus === Status.PendingApproval && (
            <>
              <Button size="sm" onClick={() => setActionType("approve")}>Approve</Button>
              <Button size="sm" variant="destructive" onClick={() => setActionType("reject")}>Reject</Button>
            </>
          )}
        </div>
      } />

      {/* Overview Cards */}
      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        <Card><CardContent className="pt-4"><div className="flex items-center gap-2">{project.recordStatus != null && <StatusBadge status={project.recordStatus} />}</div><p className="text-xs text-muted-foreground mt-1">Status</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-2xl font-bold">{members.length}</div><p className="text-xs text-muted-foreground">Members</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-2xl font-bold">{workProducts.length}</div><p className="text-xs text-muted-foreground">Work Products</p></CardContent></Card>
        <Card><CardContent className="pt-4"><div className="text-sm">{project.startDate?.split("T")[0]} — {project.endDate?.split("T")[0]}</div><p className="text-xs text-muted-foreground">Duration</p></CardContent></Card>
      </div>

      <Tabs defaultValue="members">
        <TabsList><TabsTrigger value="members">Members</TabsTrigger><TabsTrigger value="workProducts">Work Products</TabsTrigger></TabsList>

        <TabsContent value="members" className="space-y-4">
          <div className="flex justify-end"><Button size="sm" variant="outline" onClick={() => setMemberOpen(true)}><UserPlus className="mr-2 h-4 w-4" />Add Member</Button></div>
          {members.length > 0 ? (
            <div className="space-y-2">
              {members.map((m) => (
                <div key={m.projectMemberId} className="flex items-center justify-between rounded-lg border p-3">
                  <div><p className="font-medium">{m.staffName ?? m.staffId}</p></div>
                  <div className="flex items-center gap-2">{m.recordStatus != null && <StatusBadge status={m.recordStatus} />}</div>
                </div>
              ))}
            </div>
          ) : <EmptyState icon={Users} title="No Members" description="No members assigned to this project yet." />}
        </TabsContent>

        <TabsContent value="workProducts" className="space-y-4">
          <div className="flex justify-end"><Button size="sm" variant="outline" onClick={() => router.push(`/projects/${projectId}/work-products`)}><Briefcase className="mr-2 h-4 w-4" />Manage WP</Button></div>
          {workProducts.length > 0 ? (
            <div className="space-y-2">
              {workProducts.map((wp) => (
                <div key={wp.workProductId} className="flex items-center justify-between rounded-lg border p-3">
                  <div className="space-y-1"><p className="font-medium">{wp.name}</p>{wp.deliverables && <p className="text-sm text-muted-foreground line-clamp-1">{wp.deliverables}</p>}</div>
                  <div className="flex items-center gap-2"><Badge variant="secondary">{wp.maxPoint} pts</Badge>{wp.recordStatus != null && <StatusBadge status={wp.recordStatus} />}</div>
                </div>
              ))}
            </div>
          ) : <EmptyState icon={Briefcase} title="No Work Products" description="No work products assigned to this project." />}
        </TabsContent>
      </Tabs>

      <FormSheet open={memberOpen} onOpenChange={setMemberOpen} title="Add Project Member">
        <div className="space-y-4">
          <div className="space-y-2"><Label>Staff ID *</Label><Input value={memberStaffId} onChange={(e) => setMemberStaffId(e.target.value)} placeholder="Enter staff ID" /></div>
          <div className="flex gap-3 pt-4">
            <Button variant="outline" className="flex-1" onClick={() => setMemberOpen(false)}>Cancel</Button>
            <Button className="flex-1" onClick={handleAddMember} disabled={saving}>{saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}Add</Button>
          </div>
        </div>
      </FormSheet>

      <ConfirmationDialog open={!!actionType} onOpenChange={(o) => { if (!o) setActionType(null); }} title={`${actionType === "submit" ? "Submit" : actionType === "approve" ? "Approve" : "Reject"} Project`} description={`Are you sure you want to ${actionType} "${project.name}"?`} variant={actionType === "reject" ? "destructive" : "default"} confirmLabel={actionType === "submit" ? "Submit" : actionType === "approve" ? "Approve" : "Reject"} showReasonInput={actionType === "reject"} reasonLabel="Reason" onConfirm={handleAction} />
    </div>
  );
}
