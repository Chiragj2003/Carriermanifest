/**
 * Type Definitions — Shared TypeScript interfaces for the CareerManifest app.
 *
 * Covers: User, Question, Assessment, Results (scores, risk, salary, roadmap).
 */

// Type definitions for CareerManifest frontend

export interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface QuestionOption {
  label: string;
  value: number;
}

export interface QuestionWeight {
  option_index: number;
  scores: Record<string, number>;
  risk_factors?: Record<string, number>;
}

export interface Question {
  id: number;
  category: string;
  question_text: string;
  options: QuestionOption[];
  weights?: QuestionWeight[];
  display_order: number;
  is_active?: boolean;
}

export interface CareerScore {
  category: string;
  score: number;
  max_score: number;
  percentage: number;
}

export interface RiskAssessment {
  score: number;
  level: string;
  factors: Record<string, number>;
}

export interface SalaryProjection {
  year_1: string;
  year_2: string;
  year_3: string;
  year_4: string;
  year_5: string;
}

export interface RoadmapStep {
  step: number;
  title: string;
  description: string;
  duration: string;
}

export interface FeatureContribution {
  feature: string;
  user_value: number;
  career_weight: number;
  contribution: number;
  percentage: number;
}

export interface CareerExplanation {
  career: string;
  top_factors: FeatureContribution[];
  summary: string;
  penalties?: string[];
}

export interface UserProfile {
  academic_strength: number;
  financial_pressure: number;
  risk_tolerance: number;
  leadership_score: number;
  tech_affinity: number;
  govt_interest: number;
  abroad_interest: number;
  income_urgency: number;
  career_instability: number;
}

export interface VersionInfo {
  assessment: string;
  weight_matrix: string;
  feature_map: string;
  model_type: string;
  model_accuracy: number;
  model_f1_score: number;
}

export interface AssessmentResult {
  scores: CareerScore[];
  best_career_path: string;
  confidence: number;
  is_multi_fit: boolean;
  risk: RiskAssessment;
  profile: UserProfile;
  explanations: CareerExplanation[];
  salary_projection: SalaryProjection;
  roadmap: RoadmapStep[];
  required_skills: string[];
  suggested_exams: string[];
  suggested_colleges: string[];
  version: VersionInfo;
  ai_explanation?: string;
}

export interface AssessmentResponse {
  id: number;
  user_id: number;
  result: AssessmentResult;
  created_at: string;
}

export interface AssessmentListItem {
  id: number;
  best_career_path: string;
  risk_level: string;
  created_at: string;
}

export interface AdminStatsResponse {
  total_users: number;
  total_assessments: number;
  total_questions: number;
  career_distribution: Record<string, number>;
  risk_distribution: Record<string, number>;
}
