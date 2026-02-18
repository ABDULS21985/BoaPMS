import Link from "next/link";
import { ShieldAlert } from "lucide-react";

export default function AccessDeniedPage() {
  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4 text-center">
      <ShieldAlert className="h-16 w-16 text-destructive" />
      <h1 className="text-3xl font-bold">Access Denied</h1>
      <p className="text-muted-foreground max-w-md">
        You do not have permission to access this page. Please contact your administrator if you believe this is an error.
      </p>
      <Link
        href="/"
        className="mt-4 inline-flex items-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
      >
        Go to Dashboard
      </Link>
    </div>
  );
}
