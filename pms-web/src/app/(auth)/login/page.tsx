"use client";

import { useState } from "react";
import { signIn } from "next-auth/react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2, Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { toast } from "sonner";

const loginSchema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
});

type LoginValues = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const form = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { username: "", password: "" },
  });

  async function onSubmit(values: LoginValues) {
    setLoading(true);
    try {
      const result = await signIn("credentials", {
        username: values.username,
        password: values.password,
        redirect: false,
      });

      if (result?.error) {
        toast.error("Invalid username or password");
      } else {
        router.push("/");
        router.refresh();
      }
    } catch {
      toast.error("An error occurred. Please try again.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen">
      {/* Left — Login Form */}
      <div className="flex w-full flex-col items-center justify-center px-6 lg:w-5/12">
        <div className="w-full max-w-md space-y-8">
          <div className="flex flex-col items-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-xl bg-primary/10">
              <span className="text-2xl font-bold text-primary">PMS</span>
            </div>
            <div className="text-center">
              <h1 className="text-2xl font-bold tracking-tight">
                Performance Management System
              </h1>
              <p className="mt-1 text-sm text-muted-foreground">
                Sign in to your account
              </p>
            </div>
          </div>

          <Card className="border-0 shadow-lg">
            <CardHeader className="pb-4">
              <CardTitle className="text-lg">Sign In</CardTitle>
            </CardHeader>
            <CardContent>
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSubmit)}
                  className="space-y-4"
                >
                  <FormField
                    control={form.control}
                    name="username"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Username</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="Enter your username"
                            autoComplete="username"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Password</FormLabel>
                        <FormControl>
                          <div className="relative">
                            <Input
                              type={showPassword ? "text" : "password"}
                              placeholder="Enter your password"
                              autoComplete="current-password"
                              {...field}
                            />
                            <button
                              type="button"
                              className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                              onClick={() => setShowPassword(!showPassword)}
                            >
                              {showPassword ? (
                                <EyeOff className="h-4 w-4" />
                              ) : (
                                <Eye className="h-4 w-4" />
                              )}
                            </button>
                          </div>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <Button
                    type="submit"
                    className="w-full"
                    disabled={loading}
                  >
                    {loading && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    {loading ? "Signing in..." : "Sign In"}
                  </Button>
                </form>
              </Form>
            </CardContent>
          </Card>

          <p className="text-center text-xs text-muted-foreground">
            Having trouble? Contact{" "}
            <a
              href="mailto:ithelpdesk@cbn.gov.ng"
              className="text-primary hover:underline"
            >
              IT Help Desk
            </a>
          </p>
        </div>
      </div>

      {/* Right — Decorative Panel */}
      <div className="hidden w-7/12 bg-gradient-to-br from-primary/90 to-primary lg:flex lg:items-center lg:justify-center">
        <div className="max-w-lg space-y-6 px-12 text-white">
          <div className="inline-flex rounded-full bg-white/20 px-4 py-2 text-sm font-medium backdrop-blur-sm">
            Point Based
          </div>
          <h2 className="text-4xl font-bold leading-tight">
            Discover the best of Performance Management
          </h2>
          <p className="text-lg text-white/80">
            Earn Points for completed activities. Track objectives, manage work
            products, and evaluate performance across your organization.
          </p>
          <div className="flex gap-3">
            <div className="h-2 w-12 rounded-full bg-white/60" />
            <div className="h-2 w-12 rounded-full bg-white/30" />
            <div className="h-2 w-12 rounded-full bg-white/30" />
          </div>
        </div>
      </div>
    </div>
  );
}
