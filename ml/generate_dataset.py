"""
CareerManifest — Synthetic Training Data Generator
===================================================
Generates realistic student career profiles inspired by:
  - Kaggle: "Career path prediction for different fields" (Minister John)
  - Kaggle: "Computer Science Students Career Prediction" (Rugved Patil)  
  - Kaggle: "HR Analytics Job Prediction" (Faisal Qureshi)
  - ResearchGate: "Students career prediction" (Ramez Maged)

Maps to our 9-feature UserProfile:
  [0] AcademicStrength    — GPA, course rigor, competitive exam prep
  [1] FinancialPressure   — Loans, dependents, family income constraints
  [2] RiskTolerance       — Willingness to take career risks
  [3] LeadershipScore     — Leadership positions, team management, public speaking
  [4] TechAffinity        — Coding skills, tech interest, projects
  [5] GovtInterest        — Interest in government/public sector stability
  [6] AbroadInterest      — Desire for international exposure
  [7] IncomeUrgency       — How quickly income is needed
  [8] CareerInstability   — Tolerance for career instability/uncertainty

Target: 6 career labels matching CareerManifest engine:
  0 = IT / Software Jobs
  1 = MBA (India)
  2 = Government Exams
  3 = Startup / Entrepreneurship
  4 = Higher Studies (India)
  5 = MS Abroad

Generates 10,000 samples with realistic feature correlations and noise.
"""

import numpy as np
import pandas as pd
import os

np.random.seed(42)

NUM_SAMPLES = 10000
NUM_FEATURES = 9
NUM_CAREERS = 6

FEATURE_NAMES = [
    "academic_strength",
    "financial_pressure",
    "risk_tolerance",
    "leadership_score",
    "tech_affinity",
    "govt_interest",
    "abroad_interest",
    "income_urgency",
    "career_instability",
]

CAREER_LABELS = [
    "IT_Software",
    "MBA_India",
    "Government_Exams",
    "Startup_Entrepreneurship",
    "Higher_Studies_India",
    "MS_Abroad",
]

# ─── Career archetype distributions ─────────────────────────────────────────
# Each career has a "typical student" profile defined by (mean, std) per feature.
# These are informed by:
#   - Kaggle dataset feature distributions (GPA→Academic, CodingSkills→TechAffinity, etc.)
#   - Research paper findings on career predictors
#   - Indian student career landscape reality

CAREER_ARCHETYPES = {
    # IT/Software: High tech, good academics, moderate risk, low govt interest
    0: {
        "academic_strength":  (0.70, 0.15),
        "financial_pressure":  (0.35, 0.20),
        "risk_tolerance":      (0.50, 0.18),
        "leadership_score":    (0.40, 0.20),
        "tech_affinity":       (0.85, 0.10),
        "govt_interest":       (0.15, 0.12),
        "abroad_interest":     (0.50, 0.22),
        "income_urgency":      (0.40, 0.20),
        "career_instability":  (0.30, 0.15),
    },
    # MBA: Strong leadership, moderate academics, low tech, career-oriented
    1: {
        "academic_strength":  (0.60, 0.18),
        "financial_pressure":  (0.30, 0.18),
        "risk_tolerance":      (0.55, 0.18),
        "leadership_score":    (0.80, 0.12),
        "tech_affinity":       (0.30, 0.18),
        "govt_interest":       (0.15, 0.12),
        "abroad_interest":     (0.45, 0.22),
        "income_urgency":      (0.35, 0.18),
        "career_instability":  (0.30, 0.15),
    },
    # Government: High stability preference, moderate academics, low risk, high govt interest
    2: {
        "academic_strength":  (0.55, 0.18),
        "financial_pressure":  (0.55, 0.18),
        "risk_tolerance":      (0.20, 0.12),
        "leadership_score":    (0.40, 0.18),
        "tech_affinity":       (0.20, 0.15),
        "govt_interest":       (0.90, 0.08),
        "abroad_interest":     (0.10, 0.10),
        "income_urgency":      (0.45, 0.18),
        "career_instability":  (0.10, 0.08),
    },
    # Startup: Very high risk tolerance, leadership, moderate tech, low stability need
    3: {
        "academic_strength":  (0.50, 0.20),
        "financial_pressure":  (0.25, 0.18),
        "risk_tolerance":      (0.85, 0.10),
        "leadership_score":    (0.78, 0.14),
        "tech_affinity":       (0.60, 0.20),
        "govt_interest":       (0.08, 0.08),
        "abroad_interest":     (0.35, 0.22),
        "income_urgency":      (0.25, 0.18),
        "career_instability":  (0.65, 0.18),
    },
    # Higher Studies (India): Very high academics, research-oriented, moderate everything
    4: {
        "academic_strength":  (0.88, 0.08),
        "financial_pressure":  (0.35, 0.18),
        "risk_tolerance":      (0.35, 0.15),
        "leadership_score":    (0.30, 0.18),
        "tech_affinity":       (0.50, 0.22),
        "govt_interest":       (0.25, 0.18),
        "abroad_interest":     (0.30, 0.20),
        "income_urgency":      (0.20, 0.15),
        "career_instability":  (0.20, 0.12),
    },
    # MS Abroad: High academics, high abroad interest, good tech, low financial pressure
    5: {
        "academic_strength":  (0.80, 0.12),
        "financial_pressure":  (0.25, 0.16),
        "risk_tolerance":      (0.50, 0.18),
        "leadership_score":    (0.40, 0.20),
        "tech_affinity":       (0.60, 0.20),
        "govt_interest":       (0.08, 0.08),
        "abroad_interest":     (0.90, 0.08),
        "income_urgency":      (0.25, 0.15),
        "career_instability":  (0.25, 0.15),
    },
}

# Career distribution (reflects Indian student landscape)
# IT is most popular, followed by Govt, MBA, Higher Studies, MS Abroad, Startup
CAREER_DISTRIBUTION = [0.28, 0.18, 0.18, 0.10, 0.14, 0.12]


def add_feature_correlations(features: np.ndarray) -> np.ndarray:
    """
    Add realistic correlations between features:
    - High financial_pressure ↔ high income_urgency (r≈0.6)
    - High academic_strength ↔ low career_instability (r≈-0.3)
    - High risk_tolerance ↔ high career_instability (r≈0.4)
    - High govt_interest ↔ low risk_tolerance (r≈-0.5)
    - High abroad_interest ↔ low govt_interest (r≈-0.4)
    """
    # financial_pressure → income_urgency correlation
    features[:, 7] = 0.4 * features[:, 7] + 0.6 * features[:, 1] + np.random.normal(0, 0.08, len(features))
    
    # academic_strength → career_instability (inverse)
    features[:, 8] = 0.7 * features[:, 8] - 0.3 * features[:, 0] + np.random.normal(0, 0.08, len(features))
    
    # risk_tolerance → career_instability (positive)
    features[:, 8] = 0.6 * features[:, 8] + 0.4 * features[:, 2] + np.random.normal(0, 0.05, len(features))
    
    # govt_interest → risk_tolerance (inverse)
    mask_govt = features[:, 5] > 0.6
    features[mask_govt, 2] *= 0.6
    
    # abroad_interest → govt_interest (inverse)
    mask_abroad = features[:, 6] > 0.7
    features[mask_abroad, 5] *= 0.3
    
    return features


def add_noise_and_boundary_cases(features: np.ndarray, labels: np.ndarray) -> tuple:
    """
    Add realistic noise:
    - 5% of samples get random feature perturbation (real students are messy)
    - 3% of samples get "confused" labels (students who end up in unexpected careers)
    - Boundary cases: students with mixed profiles
    """
    n = len(features)
    
    # Random feature perturbation
    noise_mask = np.random.random(n) < 0.05
    noise_amount = np.random.normal(0, 0.15, (noise_mask.sum(), NUM_FEATURES))
    features[noise_mask] += noise_amount
    
    # Confused labels (3% label noise — simulates students who change career paths)
    confused_mask = np.random.random(n) < 0.03
    labels[confused_mask] = np.random.randint(0, NUM_CAREERS, confused_mask.sum())
    
    # Add boundary cases: 2% of samples with mixed high features
    boundary_mask = np.random.random(n) < 0.02
    for idx in np.where(boundary_mask)[0]:
        # Student with 2-3 high features from different careers
        high_features = np.random.choice(NUM_FEATURES, size=np.random.randint(2, 4), replace=False)
        features[idx, high_features] = np.random.uniform(0.7, 0.95, len(high_features))
    
    return features, labels


def generate_career_samples(career_id: int, n_samples: int) -> tuple:
    """Generate n_samples for a specific career archetype."""
    archetype = CAREER_ARCHETYPES[career_id]
    features = np.zeros((n_samples, NUM_FEATURES))
    
    for i, fname in enumerate(FEATURE_NAMES):
        mean, std = archetype[fname]
        features[:, i] = np.random.normal(mean, std, n_samples)
    
    labels = np.full(n_samples, career_id)
    return features, labels


def generate_full_dataset() -> pd.DataFrame:
    """Generate the complete training dataset."""
    all_features = []
    all_labels = []
    
    # Generate samples per career based on distribution
    for career_id in range(NUM_CAREERS):
        n = int(NUM_SAMPLES * CAREER_DISTRIBUTION[career_id])
        features, labels = generate_career_samples(career_id, n)
        all_features.append(features)
        all_labels.append(labels)
    
    features = np.vstack(all_features)
    labels = np.concatenate(all_labels)
    
    # Add correlations
    features = add_feature_correlations(features)
    
    # Add noise and boundary cases
    features, labels = add_noise_and_boundary_cases(features, labels)
    
    # Clamp to [0, 1]
    features = np.clip(features, 0.0, 1.0)
    
    # Round to 4 decimal places
    features = np.round(features, 4)
    
    # Build DataFrame
    df = pd.DataFrame(features, columns=FEATURE_NAMES)
    df["career_label"] = labels.astype(int)
    df["career_name"] = df["career_label"].map(lambda x: CAREER_LABELS[x])
    
    # Shuffle
    df = df.sample(frac=1, random_state=42).reset_index(drop=True)
    
    return df


def print_summary(df: pd.DataFrame):
    """Print dataset summary statistics."""
    print("=" * 70)
    print("CareerManifest ML Training Dataset — Summary")
    print("=" * 70)
    print(f"\nTotal samples: {len(df)}")
    print(f"Features: {NUM_FEATURES}")
    print(f"Career classes: {NUM_CAREERS}")
    
    print("\n--- Class Distribution ---")
    for i, name in enumerate(CAREER_LABELS):
        count = (df["career_label"] == i).sum()
        pct = count / len(df) * 100
        print(f"  {name:30s}: {count:5d} ({pct:.1f}%)")
    
    print("\n--- Feature Statistics ---")
    for col in FEATURE_NAMES:
        print(f"  {col:25s}: mean={df[col].mean():.3f}  std={df[col].std():.3f}  "
              f"min={df[col].min():.3f}  max={df[col].max():.3f}")
    
    print("\n--- Feature Correlations (top) ---")
    corr = df[FEATURE_NAMES].corr()
    pairs = []
    for i in range(len(FEATURE_NAMES)):
        for j in range(i + 1, len(FEATURE_NAMES)):
            pairs.append((FEATURE_NAMES[i], FEATURE_NAMES[j], abs(corr.iloc[i, j])))
    pairs.sort(key=lambda x: x[2], reverse=True)
    for f1, f2, c in pairs[:8]:
        sign = "+" if corr.loc[f1, f2] > 0 else "-"
        print(f"  {f1:25s} ↔ {f2:25s}: {sign}{c:.3f}")


if __name__ == "__main__":
    df = generate_full_dataset()
    
    output_path = os.path.join(os.path.dirname(__file__), "career_training_data.csv")
    df.to_csv(output_path, index=False)
    print(f"\n✅ Dataset saved to: {output_path}")
    
    print_summary(df)
