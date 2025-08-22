// Package metrics fornece funções para calcular o Índice de Saúde do Código (CHI).
package metrics

import "math"

// CHI (Code Health Index) compõe MI, duplicação e complexidade em 0..100.
// Pesos padrão podem ser ajustados conforme o time evolui.

// CHIWeights define os pesos para o cálculo do índice de saúde do código.
type CHIWeights struct {
	WMI      float64 // peso do Maintainability Index (0..1)
	WDup     float64 // peso de (100 - duplication%)
	WComplex float64 // peso do inverso da complexidade
	// Range esperado para complexidade média por função.
	// Ex.: 1..10 (Go idiomático costuma ficar ~2..4)
	ComplexityMin float64
	ComplexityMax float64
}

var DefaultCHI = CHIWeights{
	WMI: 0.5, WDup: 0.3, WComplex: 0.2,
	ComplexityMin: 1, ComplexityMax: 10,
}

// NormalizeComplexity converte complexidade média em 0..100 (maior é melhor).
func NormalizeComplexity(avg, min, max float64) float64 {
	if max <= min {
		return 100
	}
	if avg < min {
		avg = min
	}
	if avg > max {
		avg = max
	}
	// quanto menor a complexidade, maior a nota
	frac := (avg - min) / (max - min) // 0..1
	score := 100 * (1 - frac)
	return score
}

// ComputeCHI calcula o índice composto, clampando o resultado em 0..100.
func ComputeCHI(mi float64, duplicationPct float64, cyclomaticAvg float64, w CHIWeights) float64 {
	if w.WMI+w.WDup+w.WComplex <= 0 {
		w = DefaultCHI
	}
	mi = clamp(mi,
		0,
		100)
	dupScore := clamp(100-duplicationPct,
		0,
		100)
	cx := NormalizeComplexity(cyclomaticAvg, w.ComplexityMin, w.ComplexityMax)
	raw := w.WMI*mi + w.WDup*dupScore + w.WComplex*cx
	// normaliza pelos pesos para manter 0..100
	totalW := w.WMI + w.WDup + w.WComplex
	score := raw / totalW
	return clamp(score,
		0,
		100)
}

func clamp(x, lo, hi float64) float64 {
	return math.Max(lo, math.Min(hi, x))
}

// GradeFromCHI mapeia score para uma nota.
func GradeFromCHI(chi float64) string {
	switch {
	case chi >= 95:
		return "A+"
	case chi >= 90:
		return "A"
	case chi >= 85:
		return "A-"
	case chi >= 80:
		return "B+"
	case chi >= 75:
		return "B"
	case chi >= 70:
		return "B-"
	case chi >= 65:
		return "C+"
	case chi >= 55:
		return "C"
	case chi >= 45:
		return "C-"
	case chi >= 35:
		return "D"
	default:
		return "E"
	}
}

// DORA Metrics types and calculations

// DoraMetrics representa as métricas DORA fundamentais
type DoraMetrics struct {
	DeploymentFrequency float64 `json:"deployment_frequency"` // deploys per period
	PeriodUnit          string  `json:"period_unit"`          // "week", "month"
	LeadTimeP50         float64 `json:"lead_time_p50"`        // hours
	LeadTimeP95         float64 `json:"lead_time_p95"`        // hours
	ChangeFailRate      float64 `json:"change_fail_rate"`     // percentage
	MTTR                float64 `json:"mttr"`                 // hours
}

// DoraGrade calcula o grade DORA baseado nas métricas
func DoraGrade(metrics DoraMetrics) string {
	score := calculateDoraScore(metrics)
	return GradeFromCHI(score) // Reutiliza a mesma escala
}

func calculateDoraScore(metrics DoraMetrics) float64 {
	// Normalize each metric to 0-100 scale
	deployScore := normalizeDeployFreq(metrics.DeploymentFrequency, metrics.PeriodUnit)
	leadScore := normalizeLeadTime(metrics.LeadTimeP95)
	failScore := normalizeFailRate(metrics.ChangeFailRate)
	mttrScore := normalizeMTTR(metrics.MTTR)

	// Weighted average
	score := 0.3*deployScore + 0.3*leadScore + 0.2*failScore + 0.2*mttrScore
	return clamp(score, 0, 100)
}

func normalizeDeployFreq(freq float64, unit string) float64 {
	// Convert to weekly frequency for normalization
	weeklyFreq := freq
	if unit == "month" {
		weeklyFreq = freq / 4.33 // ~4.33 weeks per month
	}

	switch {
	case weeklyFreq >= 5: // Multiple deploys per day (daily+)
		return 100
	case weeklyFreq >= 1: // Weekly deploys
		return 80
	case weeklyFreq >= 0.25: // Monthly deploys
		return 60
	case weeklyFreq >= 0.08: // Quarterly deploys
		return 40
	default: // Less than quarterly
		return 20
	}
}

func normalizeLeadTime(hours float64) float64 {
	switch {
	case hours <= 24: // Less than 1 day - Elite
		return 100
	case hours <= 168: // Less than 1 week - High
		return 80
	case hours <= 720: // Less than 1 month - Medium
		return 60
	case hours <= 4320: // Less than 6 months - Low
		return 40
	default: // More than 6 months
		return 20
	}
}

func normalizeFailRate(percentage float64) float64 {
	// Lower is better for fail rate
	switch {
	case percentage <= 5: // 0-5% - Elite
		return 100
	case percentage <= 10: // 5-10% - High
		return 80
	case percentage <= 15: // 10-15% - Medium
		return 60
	case percentage <= 30: // 15-30% - Low
		return 40
	default: // >30%
		return 20
	}
}

func normalizeMTTR(hours float64) float64 {
	// Lower is better for MTTR
	switch {
	case hours <= 1: // Less than 1 hour - Elite
		return 100
	case hours <= 24: // Less than 1 day - High
		return 80
	case hours <= 168: // Less than 1 week - Medium
		return 60
	case hours <= 720: // Less than 1 month - Low
		return 40
	default: // More than 1 month
		return 20
	}
}
