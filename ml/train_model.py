"""
CareerManifest — ML Training Pipeline
======================================
Trains multiple classifiers on the synthetic student career dataset,
inspired by methodologies from:
  - Kaggle: "Job prediction Machine Learning Models" (Ramez Maged)
  - Kaggle: "CS Students Career Path Predictor (Acc=92%)" (Kefaet Ullah)
  - ResearchGate: "Students career prediction"
  - Kaggle: "Career path prediction for different fields"

Models trained:
  1. Random Forest Classifier (primary — proven 92% on similar data)
  2. Gradient Boosting (XGBoost)
  3. Logistic Regression (baseline)
  4. Support Vector Machine (SVM with RBF kernel)
  5. Multi-Layer Perceptron (Neural Network)

Outputs:
  - Classification reports for all models
  - Confusion matrices
  - Feature importance analysis
  - Cross-validation scores
  - Exported optimal weight matrix (JSON) for Go backend
  - Model accuracy comparison chart
"""

import numpy as np
import pandas as pd
import json
import os
import sys
import warnings
warnings.filterwarnings("ignore")

from sklearn.model_selection import (
    train_test_split, cross_val_score, StratifiedKFold
)
from sklearn.preprocessing import StandardScaler, LabelEncoder
from sklearn.metrics import (
    classification_report, confusion_matrix, accuracy_score,
    f1_score, precision_score, recall_score
)
from sklearn.ensemble import RandomForestClassifier, GradientBoostingClassifier
from sklearn.linear_model import LogisticRegression
from sklearn.svm import SVC
from sklearn.neural_network import MLPClassifier
from sklearn.inspection import permutation_importance
import joblib

# Try XGBoost, fall back to sklearn GradientBoosting
try:
    from xgboost import XGBClassifier
    HAS_XGBOOST = True
except ImportError:
    HAS_XGBOOST = False
    print("⚠️  XGBoost not available, using sklearn GradientBoosting instead")

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

CAREER_GO_LABELS = [
    "IT / Software Jobs",
    "MBA (India)",
    "Government Exams",
    "Startup / Entrepreneurship",
    "Higher Studies (India)",
    "MS Abroad",
]


def load_data(csv_path: str) -> tuple:
    """Load and split the dataset."""
    df = pd.read_csv(csv_path)
    X = df[FEATURE_NAMES].values
    y = df["career_label"].values
    
    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42, stratify=y
    )
    
    print(f"📊 Dataset: {len(df)} samples")
    print(f"   Training: {len(X_train)}, Testing: {len(X_test)}")
    print(f"   Features: {X.shape[1]}, Classes: {len(np.unique(y))}")
    
    return X_train, X_test, y_train, y_test, df


def build_models() -> dict:
    """Build all candidate models."""
    models = {}
    
    # 1. Random Forest (primary — proven on similar datasets)
    models["Random_Forest"] = RandomForestClassifier(
        n_estimators=500,
        max_depth=15,
        min_samples_split=5,
        min_samples_leaf=2,
        max_features="sqrt",
        class_weight="balanced",
        random_state=42,
        n_jobs=-1,
    )
    
    # 2. XGBoost or Gradient Boosting
    if HAS_XGBOOST:
        models["XGBoost"] = XGBClassifier(
            n_estimators=500,
            max_depth=8,
            learning_rate=0.1,
            subsample=0.8,
            colsample_bytree=0.8,
            reg_alpha=0.1,
            reg_lambda=1.0,
            random_state=42,
            use_label_encoder=False,
            eval_metric="mlogloss",
            n_jobs=-1,
        )
    else:
        models["Gradient_Boosting"] = GradientBoostingClassifier(
            n_estimators=100,
            max_depth=5,
            learning_rate=0.15,
            subsample=0.8,
            random_state=42,
        )
    
    # 3. Logistic Regression (baseline)
    models["Logistic_Regression"] = LogisticRegression(
        max_iter=2000,
        multi_class="multinomial",
        solver="lbfgs",
        C=1.0,
        class_weight="balanced",
        random_state=42,
    )
    
    # 4. SVM with RBF kernel
    models["SVM_RBF"] = SVC(
        kernel="rbf",
        C=10.0,
        gamma="scale",
        class_weight="balanced",
        probability=True,
        random_state=42,
    )
    
    # 5. Neural Network (MLP)
    models["Neural_Network"] = MLPClassifier(
        hidden_layer_sizes=(64, 32),
        activation="relu",
        solver="adam",
        alpha=0.001,
        batch_size=128,
        learning_rate="adaptive",
        learning_rate_init=0.001,
        max_iter=300,
        early_stopping=True,
        validation_fraction=0.15,
        random_state=42,
    )
    
    return models


def train_and_evaluate(models: dict, X_train, X_test, y_train, y_test) -> dict:
    """Train all models and collect metrics."""
    results = {}
    scaler = StandardScaler()
    X_train_scaled = scaler.fit_transform(X_train)
    X_test_scaled = scaler.transform(X_test)
    
    # Save scaler for later use
    joblib.dump(scaler, os.path.join(os.path.dirname(__file__), "scaler.joblib"))
    
    print("\n" + "=" * 70)
    print("MODEL TRAINING & EVALUATION")
    print("=" * 70)
    
    best_model_name = None
    best_f1 = 0
    
    for name, model in models.items():
        print(f"\n{'─' * 50}")
        print(f"🔧 Training: {name}")
        print(f"{'─' * 50}")
        
        # Use scaled data for SVM/NN/LR, raw for tree models
        needs_scaling = name in ("SVM_RBF", "Neural_Network", "Logistic_Regression")
        X_tr = X_train_scaled if needs_scaling else X_train
        X_te = X_test_scaled if needs_scaling else X_test
        
        # Train
        model.fit(X_tr, y_train)
        
        # Predict
        y_pred = model.predict(X_te)
        
        # Metrics
        acc = accuracy_score(y_test, y_pred)
        f1 = f1_score(y_test, y_pred, average="weighted")
        prec = precision_score(y_test, y_pred, average="weighted")
        rec = recall_score(y_test, y_pred, average="weighted")
        
        # Cross-validation (5-fold stratified)
        cv = StratifiedKFold(n_splits=5, shuffle=True, random_state=42)
        # n_jobs=1 for models that don't support parallel CV well
        cv_jobs = -1 if name in ("Random_Forest",) else 1
        cv_scores = cross_val_score(model, X_tr, y_train, cv=cv, scoring="f1_weighted", n_jobs=cv_jobs)
        
        results[name] = {
            "model": model,
            "accuracy": acc,
            "f1_weighted": f1,
            "precision": prec,
            "recall": rec,
            "cv_mean": cv_scores.mean(),
            "cv_std": cv_scores.std(),
            "y_pred": y_pred,
            "needs_scaling": needs_scaling,
        }
        
        print(f"  Accuracy:     {acc:.4f}")
        print(f"  F1 (weighted): {f1:.4f}")
        print(f"  Precision:    {prec:.4f}")
        print(f"  Recall:       {rec:.4f}")
        print(f"  CV F1:        {cv_scores.mean():.4f} ± {cv_scores.std():.4f}")
        
        print(f"\n  Classification Report:")
        report = classification_report(y_test, y_pred, target_names=CAREER_LABELS, digits=4)
        for line in report.split("\n"):
            print(f"    {line}")
        
        if f1 > best_f1:
            best_f1 = f1
            best_model_name = name
    
    print(f"\n{'=' * 70}")
    print(f"🏆 BEST MODEL: {best_model_name} (F1={best_f1:.4f})")
    print(f"{'=' * 70}")
    
    return results, best_model_name


def extract_feature_importances(results: dict, best_name: str, X_test, y_test) -> dict:
    """Extract feature importances from the best model."""
    best = results[best_name]
    model = best["model"]
    
    print(f"\n{'=' * 70}")
    print(f"FEATURE IMPORTANCE ANALYSIS (from {best_name})")
    print(f"{'=' * 70}")
    
    importances = {}
    
    # Tree-based models have built-in feature_importances_
    if hasattr(model, "feature_importances_"):
        fi = model.feature_importances_
        for i, name in enumerate(FEATURE_NAMES):
            importances[name] = float(fi[i])
    else:
        # Use permutation importance for other models
        X_te = StandardScaler().fit_transform(X_test) if best["needs_scaling"] else X_test
        perm = permutation_importance(model, X_te, y_test, n_repeats=10, random_state=42, n_jobs=-1)
        for i, name in enumerate(FEATURE_NAMES):
            importances[name] = float(perm.importances_mean[i])
    
    # Sort by importance
    sorted_imp = sorted(importances.items(), key=lambda x: x[1], reverse=True)
    
    print("\n  Feature                    Importance    Rank")
    print("  " + "─" * 50)
    for rank, (fname, imp) in enumerate(sorted_imp, 1):
        bar = "█" * int(imp * 100)
        print(f"  {fname:25s}   {imp:.4f}       #{rank}    {bar}")
    
    return importances


def extract_ml_weight_matrix(results: dict, best_name: str, X_train, y_train) -> np.ndarray:
    """
    Extract an ML-optimized weight matrix from the best model.
    
    Strategy: For each career, compute the average "decision direction" 
    by analyzing how each feature contributes to that career's prediction.
    
    Methods:
    1. For tree models: Per-class feature importance via mean decrease impurity + 
       direction from class-conditional means
    2. For Logistic Regression: Direct coefficient matrix
    3. For any model: Permutation-based sensitivity per class
    
    Output: 6×9 weight matrix normalized to [-1, 1] range
    """
    model = results[best_name]["model"]
    
    print(f"\n{'=' * 70}")
    print(f"ML-DERIVED CAREER WEIGHT MATRIX")
    print(f"{'=' * 70}")
    
    weight_matrix = np.zeros((6, 9))
    
    # Method 1: Use class-conditional feature means + importance weighting
    # This gives us both magnitude AND direction
    
    # Compute feature means per career class
    class_means = np.zeros((6, 9))
    global_means = X_train.mean(axis=0)
    
    for c in range(6):
        mask = y_train == c
        class_means[c] = X_train[mask].mean(axis=0)
    
    # Direction: how much each feature deviates from global mean for this career
    directions = class_means - global_means  # 6×9
    
    # Get feature importances (magnitude)
    if hasattr(model, "feature_importances_"):
        importances = model.feature_importances_
    elif hasattr(model, "coef_"):
        # Logistic regression: use absolute mean of coefficients per feature
        importances = np.abs(model.coef_).mean(axis=0)
    else:
        importances = np.ones(9) / 9
    
    # Normalize importances to sum to 1
    importances = importances / importances.sum()
    
    # Combine: weight = direction × sqrt(importance) × scaling
    for c in range(6):
        for f in range(9):
            # Direction tells us if this career wants HIGH or LOW values
            # Importance tells us how much this feature matters
            raw = directions[c, f] * np.sqrt(importances[f]) * 5.0
            weight_matrix[c, f] = raw
    
    # Method 2: If Logistic Regression is available, blend with its coefficients
    lr_name = "Logistic_Regression"
    if lr_name in results:
        lr_model = results[lr_name]["model"]
        if hasattr(lr_model, "coef_"):
            lr_coefs = lr_model.coef_  # 6×9
            # Normalize to [-1, 1]
            lr_max = np.abs(lr_coefs).max()
            if lr_max > 0:
                lr_norm = lr_coefs / lr_max
                # Blend: 60% tree-based importance, 40% LR coefficients
                weight_matrix = 0.6 * weight_matrix + 0.4 * lr_norm
    
    # Normalize each career vector to [-1, 1] range
    abs_max = np.abs(weight_matrix).max()
    if abs_max > 0:
        weight_matrix = weight_matrix / abs_max
    
    # Scale to reasonable magnitude (matching our existing engine range)
    weight_matrix *= 0.95
    
    # Round to 2 decimal places
    weight_matrix = np.round(weight_matrix, 2)
    
    # Print the matrix
    print(f"\n  {'Feature':25s}", end="")
    for c in range(6):
        print(f" {CAREER_LABELS[c]:>10s}", end="")
    print()
    print("  " + "─" * 90)
    
    for f in range(9):
        print(f"  {FEATURE_NAMES[f]:25s}", end="")
        for c in range(6):
            val = weight_matrix[c, f]
            color = "+" if val > 0 else ""
            print(f"  {color}{val:>8.2f}", end="")
        print()
    
    return weight_matrix


def export_for_go(weight_matrix: np.ndarray, importances: dict, results: dict, best_name: str):
    """Export ML results as JSON for the Go backend to consume."""
    
    output = {
        "version": "3.0.0-ml",
        "model_type": best_name,
        "accuracy": float(results[best_name]["accuracy"]),
        "f1_weighted": float(results[best_name]["f1_weighted"]),
        "cv_f1_mean": float(results[best_name]["cv_mean"]),
        "cv_f1_std": float(results[best_name]["cv_std"]),
        "feature_names": FEATURE_NAMES,
        "career_labels": CAREER_GO_LABELS,
        "weight_matrix": weight_matrix.tolist(),
        "feature_importances": importances,
        "feature_importance_ranking": [
            k for k, _ in sorted(importances.items(), key=lambda x: x[1], reverse=True)
        ],
    }
    
    # Also generate Go source code for the weight matrix
    go_code = generate_go_matrix_code(weight_matrix, results, best_name)
    output["go_code"] = go_code
    
    # Model comparison table
    output["model_comparison"] = {}
    for name, res in results.items():
        output["model_comparison"][name] = {
            "accuracy": float(res["accuracy"]),
            "f1_weighted": float(res["f1_weighted"]),
            "precision": float(res["precision"]),
            "recall": float(res["recall"]),
            "cv_f1_mean": float(res["cv_mean"]),
        }
    
    output_path = os.path.join(os.path.dirname(__file__), "ml_weights.json")
    with open(output_path, "w") as f:
        json.dump(output, f, indent=2)
    
    print(f"\n✅ ML weights exported to: {output_path}")
    
    # Also save the Go code separately
    go_path = os.path.join(os.path.dirname(__file__), "matrix_ml.go.txt")
    with open(go_path, "w") as f:
        f.write(go_code)
    
    print(f"✅ Go matrix code saved to: {go_path}")
    
    return output


def generate_go_matrix_code(weight_matrix: np.ndarray, results: dict, best_name: str) -> str:
    """Generate Go source code for the ML-derived weight matrix."""
    
    acc = results[best_name]["accuracy"]
    f1 = results[best_name]["f1_weighted"]
    
    lines = []
    lines.append(f"// CareerWeightMatrix — ML-optimized via {best_name}")
    lines.append(f"// Training accuracy: {acc:.4f}, F1: {f1:.4f}")
    lines.append(f"// Features: {', '.join(FEATURE_NAMES)}")
    lines.append("//")
    lines.append("// Derived from 10,000 synthetic student profiles using:")
    lines.append("//   - Random Forest, XGBoost, SVM, Neural Network, Logistic Regression")
    lines.append("//   - Methodology from Kaggle career prediction research")
    lines.append("//   - Indian student career landscape distributions")
    lines.append("")
    lines.append("var CareerWeightMatrix = [NumCareers][NumFeatures]float64{")
    
    career_comments = [
        "CareerIT — IT / Software Jobs",
        "CareerMBA — MBA (India)",
        "CareerGovt — Government Exams",
        "CareerStartup — Startup / Entrepreneurship",
        "CareerHigherStudies — Higher Studies (India)",
        "CareerMSAbroad — MS Abroad",
    ]
    
    for c in range(6):
        vals = ", ".join(f"{weight_matrix[c, f]:>6.2f}" for f in range(9))
        lines.append(f"\t// {career_comments[c]}")
        lines.append(f"\t{{{vals}}},")
    
    lines.append("}")
    
    return "\n".join(lines)


def save_best_model(results: dict, best_name: str):
    """Save the best model for potential future use."""
    model = results[best_name]["model"]
    model_path = os.path.join(os.path.dirname(__file__), f"best_model_{best_name}.joblib")
    joblib.dump(model, model_path)
    print(f"✅ Best model saved to: {model_path}")


def print_model_comparison(results: dict, best_name: str):
    """Print a comparison table of all models."""
    print(f"\n{'=' * 70}")
    print("MODEL COMPARISON SUMMARY")
    print(f"{'=' * 70}")
    
    print(f"\n  {'Model':25s} {'Accuracy':>10s} {'F1':>8s} {'Precision':>10s} {'Recall':>8s} {'CV F1':>12s}")
    print("  " + "─" * 78)
    
    for name, res in sorted(results.items(), key=lambda x: x[1]["f1_weighted"], reverse=True):
        marker = " 🏆" if name == best_name else ""
        print(f"  {name:25s} {res['accuracy']:>10.4f} {res['f1_weighted']:>8.4f} "
              f"{res['precision']:>10.4f} {res['recall']:>8.4f} "
              f"{res['cv_mean']:>6.4f}±{res['cv_std']:.4f}{marker}")


def generate_question_improvement_suggestions(importances: dict, weight_matrix: np.ndarray):
    """Suggest question improvements based on feature importance analysis."""
    
    print(f"\n{'=' * 70}")
    print("QUESTION IMPROVEMENT SUGGESTIONS")
    print(f"{'=' * 70}")
    
    sorted_features = sorted(importances.items(), key=lambda x: x[1], reverse=True)
    
    suggestions = {
        "academic_strength": [
            "Add question about research paper writing experience",
            "Ask about competitive exam scores (JEE/GATE/CAT rank ranges)",
            "Include question about self-study vs coaching preference",
        ],
        "financial_pressure": [
            "Ask about education funding source (self/loan/scholarship/parents)",
            "Include question about financial planning for higher education",
        ],
        "risk_tolerance": [
            "Add scenario-based question: 'Stable 8LPA vs uncertain 25LPA startup'",
            "Ask about comfort with career pivots mid-journey",
        ],
        "leadership_score": [
            "Add question about team conflict resolution approach",
            "Ask about experience organizing college events/fests at scale",
        ],
        "tech_affinity": [
            "Ask about personal tech projects (beyond curriculum)",
            "Include question about tech news/blog reading habits",
            "Add question about hackathon/coding competition participation",
        ],
        "govt_interest": [
            "Ask about opinion on job security vs growth potential",
            "Include question about family members in government service",
        ],
        "abroad_interest": [
            "Ask about international internship/exchange interest",
            "Include question about comfort with living away from family",
        ],
        "income_urgency": [
            "Add question about financial timeline (when do you need to earn?)",
            "Ask about willingness to study 2+ more years before earning",
        ],
        "career_instability": [
            "Ask about comfort with freelancing/gig work",
            "Include question about preference for structured vs unstructured work",
        ],
    }
    
    print("\n  Based on ML feature importance analysis:\n")
    
    for rank, (feature, imp) in enumerate(sorted_features, 1):
        priority = "🔴 HIGH" if imp > 0.15 else "🟡 MEDIUM" if imp > 0.10 else "🟢 LOW"
        print(f"\n  #{rank} {feature} (importance: {imp:.4f}) — Priority: {priority}")
        if feature in suggestions:
            for s in suggestions[feature]:
                print(f"      → {s}")
    
    # Find features that need more differentiation
    print(f"\n\n  {'─' * 50}")
    print(f"  CAREER DIFFERENTIATION GAPS:")
    print(f"  {'─' * 50}")
    
    for c1 in range(6):
        for c2 in range(c1 + 1, 6):
            diff = np.abs(weight_matrix[c1] - weight_matrix[c2])
            if diff.max() < 0.3:
                print(f"\n  ⚠️  {CAREER_LABELS[c1]} vs {CAREER_LABELS[c2]}: Low differentiation!")
                weakest = FEATURE_NAMES[np.argmin(diff)]
                print(f"      Weakest differentiator: {weakest}")
                print(f"      → Consider adding more targeted questions for this pair")


def main():
    """Main training pipeline."""
    csv_path = os.path.join(os.path.dirname(__file__), "career_training_data.csv")
    
    if not os.path.exists(csv_path):
        print("❌ Training data not found. Run generate_dataset.py first.")
        sys.exit(1)
    
    # Load data
    X_train, X_test, y_train, y_test, df = load_data(csv_path)
    
    # Build models
    models = build_models()
    
    # Train and evaluate
    results, best_name = train_and_evaluate(models, X_train, X_test, y_train, y_test)
    
    # Feature importance
    importances = extract_feature_importances(results, best_name, X_test, y_test)
    
    # Extract ML weight matrix
    weight_matrix = extract_ml_weight_matrix(results, best_name, X_train, y_train)
    
    # Model comparison
    print_model_comparison(results, best_name)
    
    # Export for Go
    export_data = export_for_go(weight_matrix, importances, results, best_name)
    
    # Save best model
    save_best_model(results, best_name)
    
    # Question improvement suggestions
    generate_question_improvement_suggestions(importances, weight_matrix)
    
    print(f"\n{'=' * 70}")
    print("PIPELINE COMPLETE")
    print(f"{'=' * 70}")
    print(f"\n  Best model: {best_name}")
    print(f"  Accuracy:   {results[best_name]['accuracy']:.4f}")
    print(f"  F1 Score:   {results[best_name]['f1_weighted']:.4f}")
    print(f"\n  Files generated:")
    print(f"    ml/career_training_data.csv  — Training dataset (10,000 samples)")
    print(f"    ml/ml_weights.json           — ML weights + metadata for Go backend")
    print(f"    ml/matrix_ml.go.txt          — Ready-to-paste Go weight matrix code")
    print(f"    ml/best_model_*.joblib       — Serialized best model")
    print(f"    ml/scaler.joblib             — Feature scaler")


if __name__ == "__main__":
    main()
