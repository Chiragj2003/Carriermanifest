/**
 * Dashboard Page — Shows user's assessment history, stats, and quick actions.
 *
 * Redirects unauthenticated users to /login.
 * Fetches assessment list on mount and displays stats cards + history list.
 */
"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { assessmentAPI } from "@/lib/api";
import { AssessmentListItem } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

export default function DashboardPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const [assessments, setAssessments] = useState<AssessmentListItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isLoading && !user) {
      router.push("/login");
      return;
    }

    if (user) {
      assessmentAPI.list()
        .then((res) => setAssessments(res.data))
        .catch(console.error)
        .finally(() => setLoading(false));
    }
  }, [user, isLoading, router]);

  if (isLoading || !user) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="animate-pulse text-muted-foreground">Loading...</div>
      </div>
    );
  }

  const riskBadgeVariant = (level: string) => {
    switch (level) {
      case "Low": return "success" as const;
      case "Medium": return "warning" as const;
      case "High": return "destructive" as const;
      default: return "secondary" as const;
    }
  };

  return (
    <div className="container mx-auto px-4 sm:px-6 py-6 sm:py-10">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6 sm:mb-8">
        <div>
          <h1 className="text-xl sm:text-2xl md:text-3xl font-bold">Welcome back, {user.name}!</h1>
          <p className="text-muted-foreground text-sm sm:text-base mt-1">Your career discovery dashboard</p>
        </div>
        <Link href="/assessment" className="w-full sm:w-auto">
          <Button size="lg" className="w-full sm:w-auto min-h-[44px] active:scale-[0.98] transition-transform">Take New Assessment →</Button>
        </Link>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 sm:gap-4 md:gap-6 mb-6 sm:mb-10">
        <Card className="hover:shadow-md transition-shadow">
          <CardHeader className="pb-2 p-4 sm:p-6">
            <CardDescription className="text-xs sm:text-sm">Total Assessments</CardDescription>
            <CardTitle className="text-3xl sm:text-4xl">{assessments.length}</CardTitle>
          </CardHeader>
        </Card>
        <Card className="hover:shadow-md transition-shadow">
          <CardHeader className="pb-2 p-4 sm:p-6">
            <CardDescription className="text-xs sm:text-sm">Latest Career Match</CardDescription>
            <CardTitle className="text-base sm:text-lg">
              {assessments.length > 0 ? assessments[0].best_career_path : "Take your first assessment"}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card className="hover:shadow-md transition-shadow">
          <CardHeader className="pb-2 p-4 sm:p-6">
            <CardDescription className="text-xs sm:text-sm">Risk Level</CardDescription>
            <CardTitle>
              {assessments.length > 0 ? (
                <Badge variant={riskBadgeVariant(assessments[0].risk_level)} className="text-base px-3 py-1">
                  {assessments[0].risk_level} Risk
                </Badge>
              ) : (
                <span className="text-muted-foreground text-sm">Not assessed yet</span>
              )}
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Assessment History */}
      <Card>
        <CardHeader>
          <CardTitle>Assessment History</CardTitle>
          <CardDescription>View all your past career assessments</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="text-center py-8 text-muted-foreground">Loading assessments...</div>
          ) : assessments.length === 0 ? (
            <div className="text-center py-12">
              <div className="text-4xl mb-4">🎯</div>
              <h3 className="text-lg font-semibold mb-2">No assessments yet</h3>
              <p className="text-muted-foreground mb-6">
                Take your first career assessment to discover your ideal career path.
              </p>
              <Link href="/assessment">
                <Button>Start Assessment</Button>
              </Link>
            </div>
          ) : (
            <div className="space-y-3">
              {assessments.map((a) => (
                <Link
                  key={a.id}
                  href={`/result/${a.id}`}
                  className="flex items-center justify-between p-4 rounded-lg border hover:bg-accent/50 transition-colors"
                >
                  <div className="flex items-center gap-3 sm:gap-4 min-w-0">
                    <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center text-primary font-semibold flex-shrink-0">
                      #{a.id}
                    </div>
                    <div className="min-w-0">
                      <p className="font-medium truncate">{a.best_career_path}</p>
                      <p className="text-sm text-muted-foreground">
                        {new Date(a.created_at).toLocaleDateString("en-IN", {
                          day: "numeric", month: "long", year: "numeric"
                        })}
                      </p>
                    </div>
                  </div>
                  <Badge variant={riskBadgeVariant(a.risk_level)}>
                    {a.risk_level} Risk
                  </Badge>
                </Link>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
