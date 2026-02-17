"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import {
  User,
  Mail,
  Building2,
  MapPin,
  Briefcase,
  Users,
  Phone,
  BadgeCheck,
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { getEmployeeDetail, getStaffIdMask } from "@/lib/api/dashboard";
import type { EmployeeErpDetails } from "@/types/dashboard";
import type { StaffIdMask } from "@/types/staff";

interface InfoRowProps {
  icon: React.ElementType;
  label: string;
  value?: string | null;
}

function InfoRow({ icon: Icon, label, value }: InfoRowProps) {
  if (!value) return null;
  return (
    <div className="flex items-center gap-3 py-2">
      <Icon className="h-4 w-4 shrink-0 text-muted-foreground" />
      <div>
        <p className="text-xs text-muted-foreground">{label}</p>
        <p className="text-sm font-medium">{value}</p>
      </div>
    </div>
  );
}

export default function ProfilePage() {
  const { data: session } = useSession();
  const user = session?.user;
  const userId = user?.id ?? "";

  const [loading, setLoading] = useState(true);
  const [staffDetail, setStaffDetail] = useState<EmployeeErpDetails | null>(null);
  const [staffMask, setStaffMask] = useState<StaffIdMask | null>(null);
  const [supervisor, setSupervisor] = useState<EmployeeErpDetails | null>(null);
  const [headOfOffice, setHeadOfOffice] = useState<EmployeeErpDetails | null>(null);

  useEffect(() => {
    if (!userId) return;

    async function load() {
      setLoading(true);
      try {
        const [detailRes, maskRes] = await Promise.allSettled([
          getEmployeeDetail(userId),
          getStaffIdMask(userId),
        ]);

        if (detailRes.status === "fulfilled" && detailRes.value?.data) {
          const detail = detailRes.value.data;
          setStaffDetail(detail);

          // Load supervisor and head of office
          const [supRes, hofRes] = await Promise.allSettled([
            detail.supervisorId ? getEmployeeDetail(detail.supervisorId) : Promise.resolve(null),
            detail.headOfOfficeId ? getEmployeeDetail(detail.headOfOfficeId) : Promise.resolve(null),
          ]);

          if (supRes.status === "fulfilled" && supRes.value?.data) {
            setSupervisor(supRes.value.data);
          }
          if (hofRes.status === "fulfilled" && hofRes.value?.data) {
            setHeadOfOffice(hofRes.value.data);
          }
        }

        if (maskRes.status === "fulfilled" && maskRes.value?.data) {
          setStaffMask(maskRes.value.data);
        }
      } catch {
        // Non-fatal
      } finally {
        setLoading(false);
      }
    }

    load();
  }, [userId]);

  if (loading) {
    return (
      <div>
        <PageHeader title="My Profile" breadcrumbs={[{ label: "Profile" }]} />
        <PageSkeleton />
      </div>
    );
  }

  const initials = user
    ? `${(user.firstName ?? "")[0] ?? ""}${(user.lastName ?? "")[0] ?? ""}`.toUpperCase()
    : "?";

  return (
    <div className="space-y-6">
      <PageHeader
        title="My Profile"
        description="Your personal and organizational information"
        breadcrumbs={[{ label: "Profile" }]}
      />

      <div className="grid gap-6 md:grid-cols-3">
        {/* Profile Card */}
        <Card className="md:col-span-1">
          <CardContent className="flex flex-col items-center pt-6">
            <Avatar className="h-24 w-24">
              {staffMask?.currentStaffPhoto && (
                <AvatarImage src={staffMask.currentStaffPhoto} alt="Profile" />
              )}
              <AvatarFallback className="bg-primary/10 text-2xl font-bold text-primary">
                {initials}
              </AvatarFallback>
            </Avatar>
            <h2 className="mt-4 text-xl font-semibold">
              {user?.firstName} {user?.lastName}
            </h2>
            <p className="text-sm text-muted-foreground">{user?.email}</p>

            {/* Roles */}
            {user?.roles && user.roles.length > 0 && (
              <div className="mt-4 flex flex-wrap justify-center gap-1.5">
                {user.roles.map((role) => (
                  <Badge key={role} variant="secondary" className="text-xs">
                    <BadgeCheck className="mr-1 h-3 w-3" />
                    {role}
                  </Badge>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Details */}
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle>Employee Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1">
            <InfoRow icon={User} label="Username" value={user?.name} />
            <InfoRow icon={Mail} label="Email" value={user?.email} />
            <InfoRow icon={Phone} label="Phone" value={staffDetail?.phone} />
            <Separator className="my-2" />
            <InfoRow icon={Briefcase} label="Job Title" value={staffDetail?.jobTitle} />
            <InfoRow icon={BadgeCheck} label="Grade" value={staffDetail?.gradeName} />
            <InfoRow icon={Building2} label="Department" value={staffDetail?.departmentName} />
            <InfoRow icon={Building2} label="Division" value={staffDetail?.divisionName} />
            <InfoRow icon={MapPin} label="Office" value={staffDetail?.officeName} />
            <Separator className="my-2" />
            <InfoRow
              icon={Users}
              label="Supervisor"
              value={supervisor ? `${supervisor.firstName} ${supervisor.lastName}` : undefined}
            />
            <InfoRow
              icon={Users}
              label="Head of Office"
              value={headOfOffice ? `${headOfOffice.firstName} ${headOfOffice.lastName}` : undefined}
            />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
