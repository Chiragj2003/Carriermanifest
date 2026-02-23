package engine

// Version constants for reproducibility and A/B testing.
const (
	// AssessmentVersion tracks the overall assessment logic version.
	AssessmentVersion = "3.0.0-ml"

	// WeightMatrixVersion tracks the career weight matrix version.
	// ML-optimized via Random Forest (accuracy: 88.45%, F1: 0.8844)
	// Trained on 10,000 synthetic student profiles inspired by Kaggle datasets
	WeightMatrixVersion = "3.0.0-ml"

	// FeatureMapVersion tracks the question-to-feature mapping version.
	FeatureMapVersion = "3.0.0-ml"

	// ModelType records the ML model that produced the weight matrix.
	ModelType = "Random_Forest"

	// ModelAccuracy records the test set accuracy of the trained model.
	ModelAccuracy = 0.8845

	// ModelF1Score records the weighted F1 score of the trained model.
	ModelF1Score = 0.8844
)
