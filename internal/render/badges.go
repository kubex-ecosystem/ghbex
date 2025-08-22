package render

import (
	"fmt"
	"net/url"
)

// MarkdownShield retorna um badge do shields.io (estático) para embutir no README.
func MarkdownShield(label, value, color string) string {
	if color == "" {
		color = "blue"
	}
	u := fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s",
		url.PathEscape(label), url.PathEscape(value), url.PathEscape(color))
	return fmt.Sprintf("![%s](%s)", label, u)
}

// BuildScorecardBadges gera 5 badges canônicos do scorecard.

// Scorecard representa as métricas principais de um projeto.
type Scorecard struct {
	DLeadP95Hours float64
	DeployFreq    float64
	DeployUnit    string // day|week|month
	CHI           float64
	ReviewP50H    float64
	BusFactor     int
}

func BuildScorecardBadges(s Scorecard) []string {
	unit := s.DeployUnit
	if unit == "" {
		unit = "week"
	}
	b1 := MarkdownShield("Lead Time p95", fmt.Sprintf("%.1fh", s.DLeadP95Hours), colorForLead(s.DLeadP95Hours))
	b2 := MarkdownShield("Deploy Freq", fmt.Sprintf("%.1f/%s", s.DeployFreq, unit),
		"informational")
	b3 := MarkdownShield("Code Health (CHI)", fmt.Sprintf("%.0f/100", s.CHI), colorForCHI(s.CHI))
	b4 := MarkdownShield("First Review p50", fmt.Sprintf("%.1fh", s.ReviewP50H), colorForReview(s.ReviewP50H))
	b5 := MarkdownShield("Bus Factor", fmt.Sprintf("%d", s.BusFactor), colorForBus(s.BusFactor))
	return []string{b1, b2, b3, b4, b5}
}

func colorForCHI(x float64) string {
	switch {
	case x >= 85:
		return "brightgreen"
	case x >= 70:
		return "green"
	case x >= 60:
		return "yellowgreen"
	case x >= 50:
		return "yellow"
	case x >= 40:
		return "orange"
	default:
		return "red"
	}
}

func colorForLead(h float64) string {
	switch {
	case h <= 24:
		return "brightgreen"
	case h <= 48:
		return "green"
	case h <= 72:
		return "yellowgreen"
	case h <= 96:
		return "yellow"
	case h <= 120:
		return "orange"
	default:
		return "red"
	}
}

func colorForReview(h float64) string {
	switch {
	case h <= 8:
		return "brightgreen"
	case h <= 24:
		return "green"
	case h <= 36:
		return "yellowgreen"
	case h <= 48:
		return "yellow"
	case h <= 72:
		return "orange"
	default:
		return "red"
	}
}

func colorForBus(n int) string {
	switch {
	case n >= 5:
		return "brightgreen"
	case n >= 3:
		return "green"
	case n >= 2:
		return "yellow"
	default:
		return "red"
	}
}

// GHbexBadges gera badges específicos para o formato atual do GHbex
func GHbexBadges(score float64, grade string, language string, activity int) []string {
	var badges []string

	// Health Score Badge
	badges = append(badges, MarkdownShield("Health Score",
		fmt.Sprintf("%.1f/100", score), colorForCHI(score)))

	// Grade Badge
	badges = append(badges, MarkdownShield("Grade", grade, colorForGrade(grade)))

	// Primary Language Badge
	if language != "" {
		badges = append(badges, MarkdownShield("Language", language, "blue"))
	}

	// Activity Badge
	badges = append(badges, MarkdownShield("Activity",
		fmt.Sprintf("%d%%", activity), colorForActivity(activity)))

	return badges
}

func colorForGrade(grade string) string {
	switch grade {
	case "A+", "A":
		return "brightgreen"
	case "A-", "B+":
		return "green"
	case "B", "B-":
		return "yellowgreen"
	case "C+", "C":
		return "yellow"
	case "C-", "D":
		return "orange"
	default:
		return "red"
	}
}

func colorForActivity(activity int) string {
	switch {
	case activity >= 80:
		return "brightgreen"
	case activity >= 60:
		return "green"
	case activity >= 40:
		return "yellowgreen"
	case activity >= 20:
		return "yellow"
	case activity >= 10:
		return "orange"
	default:
		return "red"
	}
}
