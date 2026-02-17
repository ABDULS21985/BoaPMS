import { Badge } from "@/components/ui/badge";
import { Status, statusLabels, getStatusVariant } from "@/types/enums";

interface StatusBadgeProps {
  status: number;
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const label = statusLabels[status] ?? `Unknown (${status})`;
  const variant = getStatusVariant(status as Status);

  return (
    <Badge variant={variant} className={className}>
      {label}
    </Badge>
  );
}
