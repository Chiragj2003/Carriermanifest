/**
 * Assessment Page — Interactive quiz with smart conditional question flow.
 *
 * Features:
 * - Education Level gate: first asks your current education status
 * - Dynamically skips irrelevant questions (e.g. no CGPA question for 10th students)
 * - Category-grouped questions with progress bar
 * - Animated radio-button options with keyboard & touch support
 * - Navigation dots for quick question jumping
 * - Auto-scrolls to first unanswered on submit attempt
 */
"use client";

import { useEffect, useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { questionsAPI, assessmentAPI } from "@/lib/api";
import { Question } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { cn } from "@/lib/utils";

// Category display names and emojis
const categoryMeta: Record<string, { label: string; emoji: string; color: string }> = {
  academic: { label: "Academic Background", emoji: "📚", color: "bg-blue-500" },
  financial: { label: "Financial Situation", emoji: "💰", color: "bg-green-500" },
  personality: { label: "Personality & Risk", emoji: "🧠", color: "bg-purple-500" },
  career_interest: { label: "Career Interest", emoji: "🎯", color: "bg-orange-500" },
};

// Education levels for the gate question
const educationLevels = [
  { label: "Currently in 10th or below", value: "10th", emoji: "🏫" },
  { label: "Currently in 11th / 12th", value: "12th", emoji: "📖" },
  { label: "Doing Graduation (B.Tech / B.Sc / B.Com / BA etc.)", value: "graduation", emoji: "🎓" },
  { label: "Graduated (Completed degree)", value: "graduated", emoji: "👨‍🎓" },
  { label: "Post-Graduation / Working Professional", value: "postgrad", emoji: "💼" },
];

/**
 * Given the education level, returns keyword patterns to skip irrelevant questions.
 * E.g., a 10th grader shouldn't be asked about B.Tech CGPA or MBA fees.
 */
function getSkipPatterns(educationLevel: string): string[] {
  switch (educationLevel) {
    case "10th":
      return [
        "12th standard",
        "graduation branch",
        "current cgpa",
        "competitive exams",
        "work experience",
        "mba fees",
        "education loan",
      ];
    case "12th":
      return [
        "graduation branch",
        "current cgpa",
        "competitive exams",
        "work experience",
      ];
    default:
      return [];
  }
}

export default function AssessmentPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const [allQuestions, setAllQuestions] = useState<Question[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [answers, setAnswers] = useState<Record<number, number>>({});
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [loadingQuestions, setLoadingQuestions] = useState(true);

  // Gate: education level (null = not selected yet)
  const [educationLevel, setEducationLevel] = useState<string | null>(null);

  useEffect(() => {
    if (!isLoading && !user) {
      router.push("/login");
      return;
    }

    if (user) {
      questionsAPI.getActive()
        .then((res) => setAllQuestions(res.data))
        .catch(() => setError("Failed to load questions. Please try again."))
        .finally(() => setLoadingQuestions(false));
    }
  }, [user, isLoading, router]);

  // Filter questions based on education level
  const questions = useMemo(() => {
    if (!educationLevel) return allQuestions;
    const skipPatterns = getSkipPatterns(educationLevel);
    if (skipPatterns.length === 0) return allQuestions;
    return allQuestions.filter((q) => {
      const text = q.question_text.toLowerCase();
      return !skipPatterns.some((pattern) => text.includes(pattern));
    });
  }, [allQuestions, educationLevel]);

  if (isLoading || !user || loadingQuestions) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="animate-pulse text-muted-foreground">Loading assessment...</div>
      </div>
    );
  }

  if (allQuestions.length === 0) {
    return (
      <div className="container mx-auto px-4 py-20 text-center">
        <h2 className="text-2xl font-bold mb-4">No Questions Available</h2>
        <p className="text-muted-foreground">Please contact the administrator to set up assessment questions.</p>
      </div>
    );
  }

  // ── Gate Screen: Education Level Selection ──
  if (!educationLevel) {
    return (
      <div className="container mx-auto px-4 py-6 sm:py-10 max-w-3xl">
        <div className="text-center mb-8">
          <div className="inline-block px-4 py-1.5 mb-4 rounded-full bg-primary/10 text-primary text-sm font-medium border border-primary/20">
            Step 1 of 2
          </div>
          <h1 className="text-2xl sm:text-3xl font-bold mb-2">Before We Begin</h1>
          <p className="text-muted-foreground max-w-lg mx-auto">
            Tell us your current education level so we can ask you the most relevant questions.
            We won&apos;t ask about things that don&apos;t apply to your stage.
          </p>
        </div>

        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="text-lg">What is your current education status?</CardTitle>
            <CardDescription>This helps us personalize the assessment for you</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {educationLevels.map((level) => (
              <button
                key={level.value}
                onClick={() => {
                  setEducationLevel(level.value);
                  setCurrentIndex(0);
                  setAnswers({});
                }}
                className="w-full text-left p-4 rounded-lg border-2 border-border hover:border-primary/50 hover:bg-accent/50 transition-all duration-200 active:scale-[0.98] group"
              >
                <div className="flex items-center gap-3">
                  <span className="text-2xl group-hover:scale-110 transition-transform">{level.emoji}</span>
                  <span className="text-sm font-medium">{level.label}</span>
                </div>
              </button>
            ))}
          </CardContent>
        </Card>

        <div className="text-center">
          <p className="text-xs text-muted-foreground">
            💡 Based on your selection, we&apos;ll show only the most relevant questions out of {allQuestions.length}.
          </p>
        </div>
      </div>
    );
  }

  // ── Main Assessment Flow ──
  const currentQuestion = questions[currentIndex];
  if (!currentQuestion) return null;

  const meta = categoryMeta[currentQuestion.category] || { label: "General", emoji: "❓", color: "bg-gray-500" };
  const progress = ((currentIndex + 1) / questions.length) * 100;
  const isAnswered = answers[currentQuestion.id] !== undefined;

  const handleSelect = (optionIndex: number) => {
    setAnswers((prev) => ({ ...prev, [currentQuestion.id]: optionIndex }));
  };

  const handleNext = () => {
    if (currentIndex < questions.length - 1) {
      setCurrentIndex(currentIndex + 1);
    }
  };

  const handlePrev = () => {
    if (currentIndex > 0) {
      setCurrentIndex(currentIndex - 1);
    }
  };

  const handleSubmit = async () => {
    const unanswered = questions.filter((q) => answers[q.id] === undefined);
    if (unanswered.length > 0) {
      setError(`Please answer all questions. ${unanswered.length} remaining.`);
      const firstIdx = questions.findIndex((q) => answers[q.id] === undefined);
      setCurrentIndex(firstIdx);
      return;
    }

    setSubmitting(true);
    setError("");

    try {
      const payload = {
        answers: questions
          .filter((q) => answers[q.id] !== undefined)
          .map((q) => ({ question_id: q.id, selected: answers[q.id] })),
      };

      const res = await assessmentAPI.submit(payload);
      router.push(`/result/${res.data.id}`);
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to submit assessment. Please try again.");
    } finally {
      setSubmitting(false);
    }
  };

  const answeredCount = questions.filter((q) => answers[q.id] !== undefined).length;

  return (
    <div className="container mx-auto px-4 py-6 sm:py-10 max-w-3xl">
      {/* Education level badge + change link */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <span>{educationLevels.find((l) => l.value === educationLevel)?.emoji}</span>
          <span>{educationLevels.find((l) => l.value === educationLevel)?.label}</span>
        </div>
        <button
          onClick={() => { setEducationLevel(null); setCurrentIndex(0); setAnswers({}); }}
          className="text-xs text-primary hover:underline"
        >
          Change level
        </button>
      </div>

      {/* Progress Header */}
      <div className="mb-6 sm:mb-8">
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs sm:text-sm text-muted-foreground">
            Question {currentIndex + 1} of {questions.length}
          </span>
          <span className="text-xs sm:text-sm font-medium">
            {answeredCount}/{questions.length} answered
          </span>
        </div>
        <Progress value={progress} />
      </div>

      {/* Category Badge */}
      <div className="flex items-center gap-2 mb-4">
        <span className="text-2xl">{meta.emoji}</span>
        <span className={cn("text-xs px-2 py-1 rounded-full text-white", meta.color)}>
          {meta.label}
        </span>
      </div>

      {/* Question Card */}
      <Card className="mb-4 sm:mb-6">
        <CardHeader className="px-4 sm:px-6">
          <CardTitle className="text-lg sm:text-xl leading-relaxed">
            {currentQuestion.question_text}
          </CardTitle>
          <CardDescription>Select one option</CardDescription>
        </CardHeader>
        <CardContent className="space-y-2 sm:space-y-3 px-4 sm:px-6">
          {currentQuestion.options.map((option, index) => (
            <button
              key={index}
              onClick={() => handleSelect(index)}
              className={cn(
                "w-full text-left p-3 sm:p-4 rounded-lg border-2 transition-all duration-200 active:scale-[0.98]",
                answers[currentQuestion.id] === index
                  ? "border-primary bg-primary/5 ring-2 ring-primary/20"
                  : "border-border hover:border-primary/50 hover:bg-accent/50"
              )}
            >
              <div className="flex items-center gap-3">
                <div
                  className={cn(
                    "h-6 w-6 rounded-full border-2 flex items-center justify-center flex-shrink-0",
                    answers[currentQuestion.id] === index
                      ? "border-primary bg-primary"
                      : "border-muted-foreground/30"
                  )}
                >
                  {answers[currentQuestion.id] === index && (
                    <div className="h-2 w-2 rounded-full bg-white" />
                  )}
                </div>
                <span className="text-sm font-medium">{option.label}</span>
              </div>
            </button>
          ))}
        </CardContent>
      </Card>

      {/* Error */}
      {error && (
        <div className="p-3 rounded-lg bg-destructive/10 text-destructive text-sm mb-4 border border-destructive/20">
          {error}
        </div>
      )}

      {/* Navigation */}
      <div className="flex items-center justify-between gap-2">
        <Button variant="outline" onClick={handlePrev} disabled={currentIndex === 0} className="min-h-[44px] px-3 sm:px-4">
          ← <span className="hidden sm:inline ml-1">Previous</span>
        </Button>

        <div className="flex gap-2">
          {/* Mobile: show question counter */}
          <div className="flex md:hidden items-center gap-1">
            <span className="text-xs text-muted-foreground font-medium">
              {currentIndex + 1}/{questions.length}
            </span>
          </div>
          {/* Desktop: show navigation dots */}
          <div className="hidden md:flex items-center gap-1">
            {questions.map((q, i) => (
              <button
                key={q.id}
                onClick={() => setCurrentIndex(i)}
                className={cn(
                  "h-2.5 w-2.5 rounded-full transition-all",
                  i === currentIndex
                    ? "bg-primary w-6"
                    : answers[q.id] !== undefined
                      ? "bg-primary/40"
                      : "bg-muted-foreground/20"
                )}
                title={`Question ${i + 1}`}
              />
            ))}
          </div>
        </div>

        {currentIndex === questions.length - 1 ? (
          <Button onClick={handleSubmit} disabled={submitting || answeredCount < questions.length} className="min-h-[44px] px-3 sm:px-4">
            {submitting ? "Analyzing..." : <><span className="hidden sm:inline">Submit Assessment</span><span className="sm:hidden">Submit</span> ✨</>}
          </Button>
        ) : (
          <Button onClick={handleNext} disabled={!isAnswered} className="min-h-[44px] px-3 sm:px-4">
            Next →
          </Button>
        )}
      </div>

      {/* Mobile swipe hint */}
      <div className="md:hidden text-center mt-4">
        <p className="text-xs text-muted-foreground">Tap an option to select, then tap Next</p>
      </div>
    </div>
  );
}
