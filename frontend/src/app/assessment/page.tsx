/**
 * Assessment Page — Interactive quiz with 30 career-oriented questions.
 *
 * Features:
 * - Category-grouped questions with progress bar
 * - Animated radio-button options with keyboard & touch support
 * - Navigation dots for quick question jumping
 * - Auto-scrolls to first unanswered on submit attempt
 */
"use client";

import { useEffect, useState } from "react";
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

export default function AssessmentPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const [questions, setQuestions] = useState<Question[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [answers, setAnswers] = useState<Record<number, number>>({});
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [loadingQuestions, setLoadingQuestions] = useState(true);

  useEffect(() => {
    if (!isLoading && !user) {
      router.push("/login");
      return;
    }

    if (user) {
      questionsAPI.getActive()
        .then((res) => {
          setQuestions(res.data);
        })
        .catch((err) => setError("Failed to load questions. Please try again."))
        .finally(() => setLoadingQuestions(false));
    }
  }, [user, isLoading, router]);

  if (isLoading || !user || loadingQuestions) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="animate-pulse text-muted-foreground">Loading assessment...</div>
      </div>
    );
  }

  if (questions.length === 0) {
    return (
      <div className="container mx-auto px-4 py-20 text-center">
        <h2 className="text-2xl font-bold mb-4">No Questions Available</h2>
        <p className="text-muted-foreground">Please contact the administrator to set up assessment questions.</p>
      </div>
    );
  }

  const currentQuestion = questions[currentIndex];
  const meta = categoryMeta[currentQuestion.category] || { label: currentQuestion.category, emoji: "❓", color: "bg-gray-500" };
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
    // Check all questions answered
    const unanswered = questions.filter((q) => answers[q.id] === undefined);
    if (unanswered.length > 0) {
      setError(`Please answer all questions. ${unanswered.length} remaining.`);
      // Navigate to first unanswered
      const firstIdx = questions.findIndex((q) => answers[q.id] === undefined);
      setCurrentIndex(firstIdx);
      return;
    }

    setSubmitting(true);
    setError("");

    try {
      const payload = {
        answers: Object.entries(answers).map(([qid, selected]) => ({
          question_id: parseInt(qid),
          selected,
        })),
      };

      const res = await assessmentAPI.submit(payload);
      router.push(`/result/${res.data.id}`);
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to submit assessment. Please try again.");
    } finally {
      setSubmitting(false);
    }
  };

  const answeredCount = Object.keys(answers).length;

  return (
    <div className="container mx-auto px-4 py-6 sm:py-10 max-w-3xl">
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
        <div className="p-3 rounded-md bg-destructive/10 text-destructive text-sm mb-4">
          {error}
        </div>
      )}

      {/* Navigation */}
      <div className="flex items-center justify-between">
        <Button
          variant="outline"
          onClick={handlePrev}
          disabled={currentIndex === 0}
        >
          ← Previous
        </Button>

        <div className="flex gap-2">
          {/* Question dots for navigation */}
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
          <Button
            onClick={handleSubmit}
            disabled={submitting || answeredCount < questions.length}
          >
            {submitting ? "Analyzing..." : "Submit Assessment ✨"}
          </Button>
        ) : (
          <Button onClick={handleNext} disabled={!isAnswered}>
            Next →
          </Button>
        )}
      </div>
    </div>
  );
}
