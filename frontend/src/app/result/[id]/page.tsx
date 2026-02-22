/**
 * Result Page — Comprehensive assessment results with:
 * - "How This Works" explainer section
 * - Basic scoring results (works without AI)
 * - AI-powered deep analysis (when available)
 * - Interactive AI chatbot for follow-up questions
 */
"use client";

import { useEffect, useState, useRef } from "react";
import { useRouter, useParams } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { assessmentAPI, chatAPI } from "@/lib/api";
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

// Simple explanations for each career path
const careerExplanations: Record<string, string> = {
  "IT / Software Jobs": "This path is ideal if you enjoy problem-solving, coding, and technology. India's IT industry offers strong salaries even at entry level, and you can start earning within 4 years of graduation.",
  "MBA (India)": "MBA suits you if you're interested in business, management, and leadership. After 2 years of MBA from a good college, you can land roles in consulting, marketing, or finance with excellent pay.",
  "Government Exams": "Government jobs offer stability, pension, and social respect. If you're disciplined and good at studying for competitive exams like UPSC, SSC, or banking, this could be your path.",
  "Startup / Entrepreneurship": "If you have a business idea and high risk tolerance, entrepreneurship can be very rewarding. It requires patience, creativity, and the ability to handle uncertainty.",
  "Higher Studies (India)": "Pursuing M.Tech, PhD, or specialized courses in India is great if you love research and deep learning. GATE, NET, and other exams open doors to top institutes.",
  "MS Abroad": "Studying abroad (MS in USA/Europe) gives you global exposure, higher salaries, and cutting-edge research opportunities. It requires good GRE/IELTS scores and financial planning.",
};

export default function ResultPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const params = useParams();
  const [data, setData] = useState<AssessmentResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showHowItWorks, setShowHowItWorks] = useState(false);

  // Chatbot state
  const [chatOpen, setChatOpen] = useState(false);
  const [chatMessages, setChatMessages] = useState<{ role: "user" | "bot"; text: string }[]>([
    { role: "bot", text: "Hi! I'm your AI career assistant. Ask me anything about your results — exams to prepare, colleges to target, salary expectations, or next steps!" },
  ]);
  const [chatInput, setChatInput] = useState("");
  const [chatLoading, setChatLoading] = useState(false);
  const chatEndRef = useRef<HTMLDivElement>(null);

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

  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [chatMessages]);

  const handleChatSend = async () => {
    if (!chatInput.trim() || chatLoading || !data) return;
    const msg = chatInput.trim();
    setChatInput("");
    setChatMessages((prev) => [...prev, { role: "user", text: msg }]);
    setChatLoading(true);

    try {
      const res = await chatAPI.send(msg, data.id);
      setChatMessages((prev) => [...prev, { role: "bot", text: res.data.reply }]);
    } catch {
      setChatMessages((prev) => [...prev, { role: "bot", text: "Sorry, I couldn't process that. Please try again." }]);
    } finally {
      setChatLoading(false);
    }
  };

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
  const bestExplanation = careerExplanations[result.best_career_path] || "This career path aligns well with your skills, interests, and financial situation based on your assessment answers.";

  return (
    <div className="container mx-auto px-4 py-6 sm:py-10 max-w-5xl">
      {/* Header */}
      <div className="text-center mb-6 sm:mb-10">
        <h1 className="text-2xl sm:text-3xl md:text-4xl font-bold mb-2">
          Your CareerManifest Results
        </h1>
        <p className="text-muted-foreground">
          Assessment #{data.id} • {new Date(data.created_at).toLocaleDateString("en-IN", {
            day: "numeric", month: "long", year: "numeric"
          })}
        </p>
      </div>

      {/* ═══════════════════════════════════════════════════
          HOW THIS ASSESSMENT WORKS (expandable section)
          ═══════════════════════════════════════════════════ */}
      <Card className="mb-6 border-dashed">
        <button
          onClick={() => setShowHowItWorks(!showHowItWorks)}
          className="w-full text-left"
        >
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <CardTitle className="text-base sm:text-lg flex items-center gap-2">
                💡 How This Assessment Works
              </CardTitle>
              <span className="text-muted-foreground text-lg">{showHowItWorks ? "▲" : "▼"}</span>
            </div>
            {!showHowItWorks && (
              <CardDescription>Click to learn how we calculated your career path</CardDescription>
            )}
          </CardHeader>
        </button>
        {showHowItWorks && (
          <CardContent className="space-y-4 text-sm">
            <div className="grid sm:grid-cols-2 gap-4">
              <div className="p-4 rounded-lg bg-blue-500/5 border border-blue-500/20">
                <h4 className="font-semibold mb-2 text-blue-600 dark:text-blue-400">📋 Step 1: You Answered Questions</h4>
                <p className="text-muted-foreground">We asked about your education, finances, personality, and career interests — 30 questions designed by career experts.</p>
              </div>
              <div className="p-4 rounded-lg bg-green-500/5 border border-green-500/20">
                <h4 className="font-semibold mb-2 text-green-600 dark:text-green-400">⚙️ Step 2: Feature Aggregation</h4>
                <p className="text-muted-foreground">Your answers are converted into 9 measurable traits (Academic Strength, Tech Affinity, Risk Tolerance, etc.) — your personal &quot;Profile DNA&quot;.</p>
              </div>
              <div className="p-4 rounded-lg bg-purple-500/5 border border-purple-500/20">
                <h4 className="font-semibold mb-2 text-purple-600 dark:text-purple-400">📊 Step 3: Vector Scoring</h4>
                <p className="text-muted-foreground">Your profile vector is compared against career weight vectors using mathematical dot products. Each career has a unique signature of traits it values.</p>
              </div>
              <div className="p-4 rounded-lg bg-orange-500/5 border border-orange-500/20">
                <h4 className="font-semibold mb-2 text-orange-600 dark:text-orange-400">🤖 Step 4: AI Enhancement</h4>
                <p className="text-muted-foreground">An AI (Llama3-70B) reads your scores and generates personalized advice, exam plans, and monthly milestones.</p>
              </div>
            </div>
            <div className="p-4 rounded-lg bg-muted/50 border">
              <h4 className="font-semibold mb-2">📦 Where is your data stored?</h4>
              <p className="text-muted-foreground">Your answers and results are stored securely in a PostgreSQL database (hosted on Neon). Only you can access your results — nobody else can see them. The AI doesn&apos;t store your conversations.</p>
            </div>
          </CardContent>
        )}
      </Card>

      {/* ═══════════════════════════════════════════════════
          SECTION 1: BASIC RESULTS (without AI)
          ═══════════════════════════════════════════════════ */}
      <div className="mb-4">
        <div className="inline-block px-3 py-1 rounded-full bg-blue-500/10 text-blue-600 dark:text-blue-400 text-xs font-medium border border-blue-500/20 mb-4">
          📊 Basic Results — Scoring Engine Analysis
        </div>
      </div>

      {/* Best Career Path Hero */}
      <Card className="mb-6 sm:mb-8 bg-gradient-to-br from-primary/5 to-primary/10 border-primary/20">
        <CardContent className="p-4 sm:p-8 text-center">
          <p className="text-sm text-muted-foreground uppercase tracking-wider mb-2">Your Best Career Path</p>
          <h2 className="text-xl sm:text-3xl md:text-4xl font-bold text-primary mb-3">
            {result.best_career_path}
          </h2>
          {/* Simple explanation */}
          <p className="text-sm text-muted-foreground max-w-2xl mx-auto mb-4 leading-relaxed">
            {bestExplanation}
          </p>
          <div className="flex items-center justify-center gap-4 flex-wrap">
            <Badge variant={riskVariant} className="text-sm px-4 py-1.5">
              {result.risk.level} Risk (Score: {result.risk.score}/10)
            </Badge>
            <Badge variant="secondary" className="text-sm px-4 py-1.5">
              Match: {result.scores[0].percentage}%
            </Badge>
            {result.confidence !== undefined && (
              <Badge variant={result.is_multi_fit ? "outline" : "default"} className="text-sm px-4 py-1.5">
                {result.is_multi_fit
                  ? "🔀 Multi-Fit Profile"
                  : `Confidence: ${(result.confidence * 100).toFixed(0)}%`}
              </Badge>
            )}
          </div>
        </CardContent>
      </Card>

      {/* User Profile — Feature Breakdown */}
      {result.profile && (
        <Card className="mb-6 sm:mb-8">
          <CardHeader>
            <CardTitle>🧬 Your Profile DNA</CardTitle>
            <CardDescription>
              Your answers were converted into 9 measurable traits. This is the foundation of our recommendation.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-3 gap-3">
              {[
                { label: "Academic Strength", value: result.profile.academic_strength, color: "bg-blue-500" },
                { label: "Tech Affinity", value: result.profile.tech_affinity, color: "bg-cyan-500" },
                { label: "Leadership", value: result.profile.leadership_score, color: "bg-purple-500" },
                { label: "Risk Tolerance", value: result.profile.risk_tolerance, color: "bg-orange-500" },
                { label: "Abroad Interest", value: result.profile.abroad_interest, color: "bg-teal-500" },
                { label: "Govt Interest", value: result.profile.govt_interest, color: "bg-amber-500" },
                { label: "Financial Pressure", value: result.profile.financial_pressure, color: "bg-red-500" },
                { label: "Income Urgency", value: result.profile.income_urgency, color: "bg-rose-500" },
                { label: "Career Instability", value: result.profile.career_instability, color: "bg-gray-500" },
              ].map((feat) => (
                <div key={feat.label} className="text-center p-3 rounded-lg bg-muted/50">
                  <p className="text-xs text-muted-foreground mb-1">{feat.label}</p>
                  <div className="w-full h-2 bg-muted rounded-full mb-1">
                    <div
                      className={`h-2 rounded-full ${feat.color}`}
                      style={{ width: `${Math.round(feat.value * 100)}%` }}
                    />
                  </div>
                  <p className="text-sm font-semibold">{(feat.value * 100).toFixed(0)}%</p>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Why This Career — Explanations */}
      {result.explanations && result.explanations.length > 0 && (
        <Card className="mb-6 sm:mb-8 border-green-500/20">
          <CardHeader>
            <CardTitle>🔍 Why {result.best_career_path}?</CardTitle>
            <CardDescription>
              Data-driven breakdown of how each trait contributed to your top recommendation
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {result.explanations.slice(0, 3).map((exp, expIdx) => (
              <div key={exp.career} className={expIdx > 0 ? "pt-4 border-t" : ""}>
                <div className="flex items-center gap-2 mb-3">
                  {expIdx === 0 && <span className="text-yellow-500">🏆</span>}
                  {expIdx === 1 && <span className="text-gray-400">🥈</span>}
                  {expIdx === 2 && <span className="text-amber-600">🥉</span>}
                  <h4 className="font-semibold">{exp.career}</h4>
                </div>
                <div className="space-y-2">
                  {exp.top_factors
                    .filter((f) => f.contribution > 0)
                    .slice(0, 4)
                    .map((f) => (
                      <div key={f.feature} className="flex items-center gap-3">
                        <span className="text-xs text-muted-foreground w-32 sm:w-40 truncate">
                          {f.feature}
                        </span>
                        <div className="flex-1 h-2 bg-muted rounded-full">
                          <div
                            className="h-2 rounded-full bg-green-500"
                            style={{ width: `${Math.min(f.percentage, 100)}%` }}
                          />
                        </div>
                        <span className="text-xs font-medium w-12 text-right text-green-600 dark:text-green-400">
                          +{f.percentage.toFixed(0)}%
                        </span>
                      </div>
                    ))}
                  {exp.top_factors
                    .filter((f) => f.contribution < 0 && Math.abs(f.contribution) > 0.02)
                    .slice(0, 2)
                    .map((f) => (
                      <div key={f.feature} className="flex items-center gap-3 opacity-70">
                        <span className="text-xs text-muted-foreground w-32 sm:w-40 truncate">
                          {f.feature}
                        </span>
                        <div className="flex-1 h-2 bg-muted rounded-full">
                          <div
                            className="h-2 rounded-full bg-red-400"
                            style={{ width: `${Math.min(Math.abs(f.percentage || (f.contribution * 100)), 100)}%` }}
                          />
                        </div>
                        <span className="text-xs font-medium w-12 text-right text-red-500">
                          drag
                        </span>
                      </div>
                    ))}
                </div>
                {exp.penalties && exp.penalties.length > 0 && (
                  <div className="mt-2 p-2 rounded bg-red-500/5 border border-red-500/10">
                    <p className="text-xs font-medium text-red-600 dark:text-red-400 mb-1">Risk Adjustments:</p>
                    {exp.penalties.map((p, i) => (
                      <p key={i} className="text-xs text-muted-foreground">{p}</p>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </CardContent>
        </Card>
      )}

      {/* Score Breakdown */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>📊 Career Score Summary</CardTitle>
          <CardDescription>How you scored across all 6 career categories. Higher % = better fit for you.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {result.scores.map((score, idx) => (
            <div key={score.category} className="space-y-1.5">
              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-2">
                  {idx === 0 && <span className="text-yellow-500">🏆</span>}
                  {idx === 1 && <span className="text-gray-400">🥈</span>}
                  {idx === 2 && <span className="text-amber-600">🥉</span>}
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

      {/* Risk Analysis - Simplified */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>⚠️ Risk Assessment</CardTitle>
          <CardDescription>This tells you how realistic your top career path is, considering your financial and personal situation</CardDescription>
        </CardHeader>
        <CardContent>
          {/* Simple summary first */}
          <div className={cn(
            "p-4 rounded-lg mb-6 border",
            result.risk.level === "Low" ? "bg-green-500/5 border-green-500/20" :
            result.risk.level === "Medium" ? "bg-yellow-500/5 border-yellow-500/20" :
            "bg-red-500/5 border-red-500/20"
          )}>
            <p className="text-sm leading-relaxed">
              {result.risk.level === "Low" && "✅ Great news! Your risk is low. You have a stable foundation and can comfortably pursue this career path. Take your time to build skills properly."}
              {result.risk.level === "Medium" && "⚡ Your risk is moderate. You should balance ambition with practical planning. Have a backup plan while working towards your primary goal."}
              {result.risk.level === "High" && "⚠️ Your risk is high due to financial or family pressures. Focus on paths that offer quicker returns first, then pivot to your dream career over time."}
            </p>
          </div>

          <div className="grid grid-cols-2 gap-3 sm:gap-4 mb-6">
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
              <strong>How we calculated this:</strong> RiskScore = (IncomeUrgency × 35%) + (FinancialPressure × 25%) + (RiskTolerance × 20%) + (CareerInstability × 20%)
            </p>
            <p className="text-sm mt-2">
              <strong>Your Score:</strong> {result.risk.score}/10 →{" "}
              <Badge variant={riskVariant}>{result.risk.level}</Badge>
              <span className="text-muted-foreground ml-2">(0-3 = Low, 4-6 = Medium, 7-10 = High)</span>
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Salary Projection */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>💰 5-Year Salary Projection</CardTitle>
          <CardDescription>Expected salary growth for {result.best_career_path} (in Indian Rupees). These are average figures — your actual salary may vary based on skills and location.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 sm:grid-cols-5 gap-2 md:gap-4">
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
          <CardDescription>Follow these steps one by one to reach your career goals</CardDescription>
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
            <CardTitle className="text-lg">🛠️ Skills You Need</CardTitle>
            <CardDescription>Start building these skills now</CardDescription>
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
            <CardTitle className="text-lg">📝 Exams to Prepare For</CardTitle>
            <CardDescription>Important competitive exams for this path</CardDescription>
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
          <CardDescription>Top colleges/institutes for {result.best_career_path} in India</CardDescription>
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

      {/* ═══════════════════════════════════════════════════
          SECTION 2: AI-POWERED DEEP ANALYSIS
          ═══════════════════════════════════════════════════ */}
      <div className="mb-4 mt-12">
        <div className="inline-block px-3 py-1 rounded-full bg-purple-500/10 text-purple-600 dark:text-purple-400 text-xs font-medium border border-purple-500/20 mb-4">
          🤖 AI-Powered Analysis — Enhanced by Llama3-70B
        </div>
      </div>

      {result.ai_explanation ? (
        <Card className="mb-8 border-purple-500/20">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              🤖 Personalized AI Career Analysis
            </CardTitle>
            <CardDescription>
              Our AI (Llama3-70B) analyzed your scores and created a personalized career plan just for you.
              This includes a detailed preparation plan with monthly milestones.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="prose prose-sm max-w-none whitespace-pre-wrap text-sm leading-relaxed dark:prose-invert">
              {result.ai_explanation}
            </div>
          </CardContent>
        </Card>
      ) : (
        <Card className="mb-8 border-dashed">
          <CardContent className="p-6 text-center">
            <p className="text-2xl mb-2">🤖</p>
            <p className="font-medium mb-1">AI Analysis Not Available</p>
            <p className="text-sm text-muted-foreground max-w-md mx-auto">
              The basic scoring engine has provided your results above. AI-powered deep analysis 
              (monthly study plans, detailed exam syllabi, personalized tips) requires the AI module to be enabled.
            </p>
          </CardContent>
        </Card>
      )}

      {/* Version info */}
      {result.version && (
        <div className="text-center text-xs text-muted-foreground mb-4">
          Engine v{result.version.assessment} • Matrix v{result.version.weight_matrix} • Features v{result.version.feature_map}
        </div>
      )}

      {/* Actions */}
      <div className="flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-4 mb-8">
        <Link href="/dashboard" className="w-full sm:w-auto">
          <Button variant="outline" className="w-full">← Back to Dashboard</Button>
        </Link>
        <Link href="/assessment" className="w-full sm:w-auto">
          <Button className="w-full">Take Another Assessment</Button>
        </Link>
        <Button
          variant="secondary"
          className="w-full sm:w-auto"
          onClick={() => setChatOpen(true)}
        >
          💬 Ask AI Questions
        </Button>
      </div>

      {/* ═══════════════════════════════════════════════════
          AI CHATBOT WIDGET (floating)
          ═══════════════════════════════════════════════════ */}
      {/* Floating chat button */}
      {!chatOpen && (
        <button
          onClick={() => setChatOpen(true)}
          className="fixed bottom-6 right-6 h-14 w-14 rounded-full bg-primary text-white shadow-lg hover:shadow-xl hover:scale-110 transition-all flex items-center justify-center text-2xl z-50"
          title="Ask AI about your results"
        >
          💬
        </button>
      )}

      {/* Chat panel */}
      {chatOpen && (
        <div className="fixed bottom-6 right-6 w-[90vw] sm:w-[400px] h-[500px] bg-background border rounded-2xl shadow-2xl flex flex-col z-50">
          {/* Chat header */}
          <div className="flex items-center justify-between p-4 border-b bg-primary/5 rounded-t-2xl">
            <div className="flex items-center gap-2">
              <span className="text-xl">🤖</span>
              <div>
                <p className="text-sm font-semibold">Career AI Assistant</p>
                <p className="text-xs text-muted-foreground">Ask about your results</p>
              </div>
            </div>
            <button
              onClick={() => setChatOpen(false)}
              className="h-8 w-8 rounded-full hover:bg-muted flex items-center justify-center text-muted-foreground"
            >
              ✕
            </button>
          </div>

          {/* Chat messages */}
          <div className="flex-1 overflow-y-auto p-4 space-y-3">
            {chatMessages.map((msg, i) => (
              <div
                key={i}
                className={cn(
                  "max-w-[85%] p-3 rounded-2xl text-sm leading-relaxed",
                  msg.role === "user"
                    ? "ml-auto bg-primary text-primary-foreground rounded-br-md"
                    : "bg-muted rounded-bl-md"
                )}
              >
                {msg.text}
              </div>
            ))}
            {chatLoading && (
              <div className="max-w-[85%] p-3 rounded-2xl bg-muted rounded-bl-md text-sm">
                <span className="animate-pulse">Thinking...</span>
              </div>
            )}
            <div ref={chatEndRef} />
          </div>

          {/* Suggested questions */}
          {chatMessages.length <= 1 && (
            <div className="px-4 pb-2 flex flex-wrap gap-1.5">
              {[
                "What exams should I focus on?",
                "Expected salary after 3 years?",
                "Best colleges for my path?",
                "How to start preparing today?",
              ].map((q) => (
                <button
                  key={q}
                  onClick={() => { setChatInput(q); }}
                  className="text-xs px-2.5 py-1.5 rounded-full border hover:bg-accent transition-colors"
                >
                  {q}
                </button>
              ))}
            </div>
          )}

          {/* Chat input */}
          <div className="p-3 border-t">
            <div className="flex gap-2">
              <input
                type="text"
                value={chatInput}
                onChange={(e) => setChatInput(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleChatSend()}
                placeholder="Ask about your career path..."
                className="flex-1 h-10 px-4 rounded-full border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
                disabled={chatLoading}
              />
              <button
                onClick={handleChatSend}
                disabled={chatLoading || !chatInput.trim()}
                className="h-10 w-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center disabled:opacity-50 hover:opacity-90 transition-opacity"
              >
                ↑
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
