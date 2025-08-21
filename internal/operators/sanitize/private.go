package sanitize

// calculateReleaseHealth computes realistic release health score based on actual actions
func calculateReleaseHealth(action *SanitizationAction) float64 {
	if action == nil {
		return 65.0 // No releases processed
	}

	// Base score for having organized releases
	baseScore := 70.0

	// Improvement based on cleanup actions
	if action.ItemsCount > 0 {
		// Each cleaned item improves score
		improvement := float64(action.ItemsCount) * 3.0
		return min(baseScore+improvement, 95.0)
	}

	return baseScore
}

// calculateHealthScore computes realistic overall health score
func calculateHealthScore(runs, artifacts, security int) float64 {
	baseScore := 75.0

	// Each category contributes to score
	runImpact := min(float64(runs)*2.0, 10.0)
	artifactImpact := min(float64(artifacts)*1.5, 8.0)
	securityImpact := min(float64(security)*5.0, 12.0)

	return min(baseScore+runImpact+artifactImpact+securityImpact, 98.0)
}
