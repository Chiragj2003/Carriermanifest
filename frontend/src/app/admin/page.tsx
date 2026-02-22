/**
 * Admin Page — Manages question bank and displays platform analytics.
 *
 * Features:
 * - Stats dashboard (users, assessments, questions, career paths)
 * - Career & risk distribution charts
 * - Question CRUD with JSON editor modal (bottom-sheet on mobile)
 * - Category filter for question bank
 *
 * Access: Admin-role users only; redirects others to /dashboard.
 */
"use client";

import { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { adminAPI } from "@/lib/api";
import { AdminStatsResponse, Question } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

// ──────────────────────────────────────────
// Modal for Add/Edit Question
// ──────────────────────────────────────────
function QuestionModal({
  question,
  onSave,
  onClose,
}: {
  question: Question | null;
  onSave: (q: Partial<Question>) => void;
  onClose: () => void;
}) {
  const isNew = !question;
  const [questionText, setQuestionText] = useState(question?.question_text || "");
  const [category, setCategory] = useState(question?.category || "academic");
  const [optionsRaw, setOptionsRaw] = useState(
    question?.options ? JSON.stringify(question.options, null, 2) : "[]"
  );
  const [weightsRaw, setWeightsRaw] = useState(
    question?.weights ? JSON.stringify(question.weights, null, 2) : "[]"
  );
  const [isActive, setIsActive] = useState(question?.is_active ?? true);
  const [parseError, setParseError] = useState("");

  const handleSubmit = () => {
    try {
      const options = JSON.parse(optionsRaw);
      const weights = JSON.parse(weightsRaw);
      setParseError("");
      onSave({
        question_text: questionText,
        category,
        options,
        weights,
        is_active: isActive,
        ...(question ? { id: question.id } : {}),
      });
    } catch {
      setParseError("Invalid JSON in options or weights field.");
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50">
      <Card className="w-full sm:max-w-2xl max-h-[85vh] sm:max-h-[90vh] overflow-auto sm:m-4 rounded-b-none sm:rounded-b-lg">
        <CardHeader className="px-4 sm:px-6">
          <CardTitle>{isNew ? "Add Question" : "Edit Question"}</CardTitle>
          <CardDescription>
            {isNew ? "Create a new assessment question" : `Editing question #${question?.id}`}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <label className="text-sm font-medium mb-1 block">Question Text</label>
            <Input value={questionText} onChange={(e) => setQuestionText(e.target.value)} placeholder="Enter question..." />
          </div>
          <div>
            <label className="text-sm font-medium mb-1 block">Category</label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="academic">Academic</option>
              <option value="financial">Financial</option>
              <option value="personality">Personality</option>
              <option value="career_interest">Career Interest</option>
            </select>
          </div>
          <div>
            <label className="text-sm font-medium mb-1 block">Options (JSON array of strings)</label>
            <textarea
              value={optionsRaw}
              onChange={(e) => setOptionsRaw(e.target.value)}
              className="w-full h-28 rounded-md border border-input bg-background px-3 py-2 text-sm font-mono"
            />
          </div>
          <div>
            <label className="text-sm font-medium mb-1 block">Weights (JSON object)</label>
            <textarea
              value={weightsRaw}
              onChange={(e) => setWeightsRaw(e.target.value)}
              className="w-full h-40 rounded-md border border-input bg-background px-3 py-2 text-sm font-mono"
            />
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="active"
              checked={isActive}
              onChange={(e) => setIsActive(e.target.checked)}
              className="rounded"
            />
            <label htmlFor="active" className="text-sm">Active</label>
          </div>
          {parseError && <p className="text-sm text-red-500">{parseError}</p>}
          <div className="flex justify-end gap-2 pt-2">
            <Button variant="outline" onClick={onClose}>Cancel</Button>
            <Button onClick={handleSubmit}>{isNew ? "Create" : "Update"}</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

// ──────────────────────────────────────────
// Admin Page
// ──────────────────────────────────────────
export default function AdminPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const [stats, setStats] = useState<AdminStatsResponse | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [loadingStats, setLoadingStats] = useState(true);
  const [loadingQuestions, setLoadingQuestions] = useState(true);
  const [error, setError] = useState("");
  const [modalQuestion, setModalQuestion] = useState<Question | null | "new">(null);
  const [filterSection, setFilterSection] = useState("all");

  const fetchData = useCallback(() => {
    adminAPI.getStats()
      .then((res) => setStats(res.data))
      .catch(() => setError("Failed to load stats"))
      .finally(() => setLoadingStats(false));

    adminAPI.getQuestions()
      .then((res) => setQuestions(res.data))
      .catch(() => {})
      .finally(() => setLoadingQuestions(false));
  }, []);

  useEffect(() => {
    if (!isLoading && (!user || user.role !== "admin")) {
      router.push("/dashboard");
      return;
    }
    if (user?.role === "admin") {
      fetchData();
    }
  }, [user, isLoading, router, fetchData]);

  const handleSave = async (q: Partial<Question>) => {
    try {
      if (q.id) {
        await adminAPI.updateQuestion(q.id, q);
      } else {
        await adminAPI.createQuestion(q);
      }
      setModalQuestion(null);
      setLoadingQuestions(true);
      adminAPI.getQuestions()
        .then((res) => setQuestions(res.data))
        .finally(() => setLoadingQuestions(false));
    } catch {
      alert("Failed to save question");
    }
  };

  if (isLoading || loadingStats) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="animate-pulse text-muted-foreground">Loading admin panel...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-20 text-center">
        <p className="text-red-500">{error}</p>
      </div>
    );
  }

  const filteredQuestions = filterSection === "all"
    ? questions
    : questions.filter((q) => q.category === filterSection);

  return (
    <div className="container mx-auto px-4 py-10 max-w-6xl">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-6 sm:mb-8">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold">Admin Panel</h1>
          <p className="text-muted-foreground">Manage questions & view analytics</p>
        </div>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-10">
          <Card>
            <CardContent className="p-6 text-center">
              <p className="text-3xl font-bold text-primary">{stats.total_users}</p>
              <p className="text-sm text-muted-foreground mt-1">Total Users</p>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-6 text-center">
              <p className="text-3xl font-bold text-primary">{stats.total_assessments}</p>
              <p className="text-sm text-muted-foreground mt-1">Total Assessments</p>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-6 text-center">
              <p className="text-3xl font-bold text-primary">{stats.total_questions}</p>
              <p className="text-sm text-muted-foreground mt-1">Total Questions</p>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-6 text-center">
              <p className="text-3xl font-bold text-primary">
                {stats.career_distribution ? Object.keys(stats.career_distribution).length : 0}
              </p>
              <p className="text-sm text-muted-foreground mt-1">Career Paths</p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Distribution Charts (simple) */}
      {stats && (
        <div className="grid md:grid-cols-2 gap-8 mb-10">
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Career Distribution</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {stats.career_distribution && Object.entries(stats.career_distribution).length > 0 ? (
                Object.entries(stats.career_distribution).map(([career, count]: [string, number]) => (
                  <div key={career} className="flex items-center justify-between text-sm">
                    <span>{career}</span>
                    <Badge variant="secondary">{count}</Badge>
                  </div>
                ))
              ) : (
                <p className="text-sm text-muted-foreground">No assessments yet</p>
              )}
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Risk Distribution</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {stats.risk_distribution && Object.entries(stats.risk_distribution).length > 0 ? (
                Object.entries(stats.risk_distribution).map(([level, count]: [string, number]) => (
                  <div key={level} className="flex items-center justify-between text-sm">
                    <span className="flex items-center gap-2">
                      <Badge
                        variant={level === "Low" ? "success" : level === "Medium" ? "warning" : "destructive"}
                      >
                        {level}
                      </Badge>
                    </span>
                    <span className="font-medium">{count}</span>
                  </div>
                ))
              ) : (
                <p className="text-sm text-muted-foreground">No assessments yet</p>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* Questions Management */}
      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 sm:gap-4">
            <div>
              <CardTitle>Question Bank</CardTitle>
              <CardDescription>{questions.length} questions total</CardDescription>
            </div>
            <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-2 w-full sm:w-auto">
              <select
                value={filterSection}
                onChange={(e) => setFilterSection(e.target.value)}
                className="rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                <option value="all">All Categories</option>
                <option value="academic">Academic</option>
                <option value="financial">Financial</option>
                <option value="personality">Personality</option>
                <option value="career_interest">Career Interest</option>
              </select>
              <Button onClick={() => setModalQuestion("new")}>+ Add Question</Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {loadingQuestions ? (
            <div className="text-center py-8 text-muted-foreground">Loading questions...</div>
          ) : filteredQuestions.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">No questions found</div>
          ) : (
            <div className="space-y-3">
              {filteredQuestions.map((q) => (
                <div
                  key={q.id}
                  className="flex items-start justify-between gap-4 p-4 rounded-lg border hover:bg-muted/50 transition-colors"
                >
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <Badge variant="outline" className="text-xs capitalize">{q.category}</Badge>
                      {q.is_active === false && <Badge variant="destructive" className="text-xs">Inactive</Badge>}
                    </div>
                    <p className="text-sm font-medium">{q.question_text}</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {q.options.length} options{q.weights ? ` • ${q.weights.length} weight groups` : ""}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setModalQuestion(q)}
                  >
                    Edit
                  </Button>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Modal */}
      {modalQuestion !== null && (
        <QuestionModal
          question={modalQuestion === "new" ? null : modalQuestion}
          onSave={handleSave}
          onClose={() => setModalQuestion(null)}
        />
      )}
    </div>
  );
}
