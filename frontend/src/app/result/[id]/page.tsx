"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { assessmentAPI } from "@/lib/api";
import { AssessmentResponse } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { cn } from "@/lib/utils";

// Colors for each career category
const careerColors: Record<string, string> = {
  "IT / Software Jobs": "bg-blue-500",
  "MBA (India)": "bg-purple-500",
  "Government Exams": "bg-amber-500",
  "Startup / Entrepreneurship": "bg-red-500",
  "Higher Studies (India)": "bg-green-500",
  "MS Abroad": "bg-cyan-500",
};

export default function ResultPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const params = useParams();
  const [data, setData] = useState<AssessmentResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!isLoading && !user) {
      router.push("/login");
      return;
    }

    if (user && params.id) {
      assessmentAPI.getById(Number(params.id))
        .then((res) => setData(res.data))
        .catch((err) => setError(err.response?.data?.error || "Failed to load result"))
        .finally(() => setLoading(false));
    }
  }, [user, isLoading, router, params.id]);

  if (isLoading || loading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="animate-pulse text-muted-foreground">Loading your results...</div>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="container mx-auto px-4 py-20 text-center">
        <h2 className="text-2xl font-bold mb-4">Result Not Found</h2>
        <p className="text-muted-foreground mb-6">{error || "This assessment result could not be found."}</p>
        <Link href="/dashboard"><Button>Back to Dashboard</Button></Link>
      </div>
    );
  }

  const result = data.result;
  const riskVariant = result.risk.level === "Low" ? "success" : result.risk.level === "Medium" ? "warning" : "destructive";

  return (
    <div className="container mx-auto px-4 py-10 max-w-5xl">
      {/* Header */}
      <div className="text-center mb-10">
        <h1 className="text-3xl md:text-4xl font-bold mb-2">
          Your CareerManifest Results
        </h1>
        <p className="text-muted-foreground">
          Assessment #{data.id} • {new Date(data.created_at).toLocaleDateString("en-IN", {
            day: "numeric", month: "long", year: "numeric"
          })}
        </p>
      </div>

      {/* Best Career Path Hero */}
      <Card className="mb-8 bg-gradient-to-br from-primary/5 to-primary/10 border-primary/20">
        <CardContent className="p-8 text-center">
          <p className="text-sm text-muted-foreground uppercase tracking-wider mb-2">Your Best Career Path</p>
          <h2 className="text-3xl md:text-4xl font-bold text-primary mb-4">
            {result.best_career_path}
          </h2>
          <div className="flex items-center justify-center gap-4 flex-wrap">
            <Badge variant={riskVariant} className="text-sm px-4 py-1.5">
              {result.risk.level} Risk (Score: {result.risk.score}/10)
            </Badge>
            <Badge variant="secondary" className="text-sm px-4 py-1.5">
              Match: {result.scores[0].percentage}%
            </Badge>
          </div>
        </CardContent>
      </Card>

      {/* Score Breakdown */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>📊 Career Score Summary</CardTitle>
          <CardDescription>How you scored across all career categories</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {result.scores.map((score, idx) => (
            <div key={score.category} className="space-y-1.5">
              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-2">
                  {idx === 0 && <span className="text-yellow-500">🏆</span>}
                  <span className={cn("font-medium", idx === 0 && "text-primary")}>
                    {score.category}
                  </span>
                </div>
                <span className="text-muted-foreground">
                  {score.score}/{score.max_score} ({score.percentage}%)
                </span>
              </div>
              <Progress
                value={score.percentage}
                color={careerColors[score.category] || "bg-primary"}
              />
            </div>
          ))}
        </CardContent>
      </Card>

      {/* Risk Analysis */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>⚠️ Risk Assessment</CardTitle>
          <CardDescription>Based on your financial and personal situation</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            {Object.entries(result.risk.factors).map(([factor, value]) => (
              <div key={factor} className="text-center p-4 rounded-lg bg-muted/50">
                <p className="text-xs text-muted-foreground capitalize mb-1">
                  {factor.replace(/_/g, " ")}
                </p>
                <p className="text-2xl font-bold">{value.toFixed(1)}</p>
                <p className="text-xs text-muted-foreground">/10</p>
              </div>
            ))}
          </div>
          <div className="p-4 rounded-lg bg-muted/30">
            <p className="text-sm">
              <strong>Formula:</strong> RiskScore = (IncomeUrgency × 0.35) + (FamilyDependency × 0.25) + (RiskTolerance × 0.20) + (CareerInstability × 0.20)
            </p>
            <p className="text-sm mt-2">
              <strong>Your Score:</strong> {result.risk.score}/10 →{" "}
              <Badge variant={riskVariant}>{result.risk.level}</Badge>
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Salary Projection */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>💰 5-Year Salary Projection</CardTitle>
          <CardDescription>Expected salary growth for {result.best_career_path}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-5 gap-2 md:gap-4">
            {[
              { year: "Year 1", salary: result.salary_projection.year_1 },
              { year: "Year 2", salary: result.salary_projection.year_2 },
              { year: "Year 3", salary: result.salary_projection.year_3 },
              { year: "Year 4", salary: result.salary_projection.year_4 },
              { year: "Year 5", salary: result.salary_projection.year_5 },
            ].map((item, idx) => (
              <div
                key={item.year}
                className="text-center p-3 md:p-4 rounded-lg bg-gradient-to-t from-primary/5 to-primary/10 border"
              >
                <div
                  className="w-full bg-primary/20 rounded-full mb-3 mx-auto"
                  style={{ height: `${20 + idx * 20}px` }}
                >
                  <div
                    className="bg-primary rounded-full w-full"
                    style={{ height: `${20 + idx * 20}px` }}
                  />
                </div>
                <p className="text-xs text-muted-foreground">{item.year}</p>
                <p className="text-xs md:text-sm font-semibold mt-1">{item.salary}</p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Roadmap */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>🗺️ Preparation Roadmap</CardTitle>
          <CardDescription>Step-by-step plan to achieve your career goals</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-6">
            {result.roadmap.map((step, idx) => (
              <div key={step.step} className="flex gap-4">
                <div className="flex flex-col items-center">
                  <div className="h-10 w-10 rounded-full bg-primary flex items-center justify-center text-white font-bold text-sm flex-shrink-0">
                    {step.step}
                  </div>
                  {idx < result.roadmap.length - 1 && (
                    <div className="w-0.5 h-full bg-primary/20 mt-2" />
                  )}
                </div>
                <div className="pb-6">
                  <div className="flex items-center gap-2 mb-1">
                    <h4 className="font-semibold">{step.title}</h4>
                    <Badge variant="outline" className="text-xs">{step.duration}</Badge>
                  </div>
                  <p className="text-sm text-muted-foreground">{step.description}</p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Skills & Exams */}
      <div className="grid md:grid-cols-2 gap-8 mb-8">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">🛠️ Required Skills</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-2">
              {result.required_skills.map((skill) => (
                <Badge key={skill} variant="secondary" className="py-1.5">
                  {skill}
                </Badge>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">📝 Suggested Exams</CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="space-y-2">
              {result.suggested_exams.map((exam) => (
                <li key={exam} className="flex items-center gap-2 text-sm">
                  <span className="h-1.5 w-1.5 rounded-full bg-primary flex-shrink-0" />
                  {exam}
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
      </div>

      {/* Suggested Colleges */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>🏫 Suggested Institutions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid md:grid-cols-2 gap-3">
            {result.suggested_colleges.map((college) => (
              <div key={college} className="flex items-center gap-2 p-3 rounded-lg bg-muted/50 text-sm">
                <span>🎓</span>
                {college}
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* AI Explanation */}
      {result.ai_explanation && (
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>🤖 AI-Powered Analysis</CardTitle>
            <CardDescription>Personalized explanation from our AI engine</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="prose prose-sm max-w-none whitespace-pre-wrap text-sm leading-relaxed">
              {result.ai_explanation}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Actions */}
      <div className="flex items-center justify-center gap-4">
        <Link href="/dashboard">
          <Button variant="outline">← Back to Dashboard</Button>
        </Link>
        <Link href="/assessment">
          <Button>Take Another Assessment</Button>
        </Link>
      </div>
    </div>
  );
}
