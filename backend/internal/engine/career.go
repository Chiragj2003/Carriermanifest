// Package engine implements the vector-based career recommendation system.
// This is the core scoring engine of CareerManifest — a mathematically robust,
// explainable, and ML-ready career recommendation system for Indian students.
package engine

import "fmt"

// Career is a type-safe enum for career categories.
type Career int

const (
	// CareerIT represents IT / Software Jobs.
	CareerIT Career = iota
	// CareerMBA represents MBA (India).
	CareerMBA
	// CareerGovt represents Government Exams.
	CareerGovt
	// CareerStartup represents Startup / Entrepreneurship.
	CareerStartup
	// CareerHigherStudies represents Higher Studies (India).
	CareerHigherStudies
	// CareerMSAbroad represents MS Abroad.
	CareerMSAbroad
	// NumCareers is the total number of career categories.
	NumCareers
)

// careerLabels maps Career enum to human-readable labels.
var careerLabels = map[Career]string{
	CareerIT:            "IT / Software Jobs",
	CareerMBA:           "MBA (India)",
	CareerGovt:          "Government Exams",
	CareerStartup:       "Startup / Entrepreneurship",
	CareerHigherStudies: "Higher Studies (India)",
	CareerMSAbroad:      "MS Abroad",
}

// labelToCareer maps the legacy string labels to Career enum.
var labelToCareer = map[string]Career{
	"IT / Software Jobs":         CareerIT,
	"MBA (India)":                CareerMBA,
	"Government Exams":           CareerGovt,
	"Startup / Entrepreneurship": CareerStartup,
	"Higher Studies (India)":     CareerHigherStudies,
	"MS Abroad":                  CareerMSAbroad,
}

// String returns the human-readable label for a Career.
func (c Career) String() string {
	if label, ok := careerLabels[c]; ok {
		return label
	}
	return fmt.Sprintf("Unknown(%d)", int(c))
}

// AllCareers returns all valid Career values in order.
func AllCareers() []Career {
	return []Career{
		CareerIT,
		CareerMBA,
		CareerGovt,
		CareerStartup,
		CareerHigherStudies,
		CareerMSAbroad,
	}
}

// CareerFromLabel converts a legacy string label to a Career enum.
// Returns CareerIT and false if the label is unknown.
func CareerFromLabel(label string) (Career, bool) {
	c, ok := labelToCareer[label]
	return c, ok
}
