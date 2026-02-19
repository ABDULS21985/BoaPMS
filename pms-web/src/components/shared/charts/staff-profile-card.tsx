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
