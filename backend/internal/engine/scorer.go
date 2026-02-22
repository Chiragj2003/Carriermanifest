package engine

// Dot computes the dot product of two equal-length vectors.
func Dot(a, b []float64) float64 {
	sum := 0.0
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		sum += a[i] * b[i]
	}
	return sum
}

// RawCareerScore holds the raw (un-normalized) dot product score for a career.
type RawCareerScore struct {
	Career Career
	Score  float64
}

// ComputeRawScores calculates the dot product of the user profile vector
// against each career weight vector, producing raw scores.
func ComputeRawScores(profile *UserProfile) []RawCareerScore {
	userVec := profile.Vector()
	scores := make([]RawCareerScore, int(NumCareers))

	for _, c := range AllCareers() {
		weights := GetCareerWeights(c)
		scores[int(c)] = RawCareerScore{
			Career: c,
			Score:  Dot(userVec, weights),
		}
	}

	return scores
}
