"use client";

import { Construction } from "lucide-react";
import { PageHeader } from "@/components/shared/page-header";
import { EmptyState } from "@/components/shared/empty-state";

export default function ObjectiveAssignmentsPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Objective Assignments"
        description="Assign objectives to staff members"
        breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Objective Assignments" }]}
      />
      <EmptyState
        icon={Construction}
        title="Coming Soon"
        description="Objective assignment functionality is under development. This feature will allow you to assign objectives to individual staff members."
      />
    </div>
  );
}
