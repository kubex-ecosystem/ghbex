// Package render fornece funções para gerar gráficos e visualizações simples em SVG.
package render

import (
	"fmt"
	"math"
	"os"
)

// WriteSparklineSVG gera um sparkline simples (sem deps) de uma série 0..infinito.
// Escala Y automaticamente pelo min/max da série.
func WriteSparklineSVG(path string, values []float64, width, height int) error {
	if len(values) == 0 {
		return os.WriteFile(path,
			[]byte(emptySVG(width, height)),
			0644)
	}
	minV, maxV := values[0], values[0]
	for _, v := range values {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	if math.Abs(maxV-minV) < 1e-9 {
		maxV = minV + 1
	} // evita div/0

	// padding visual
	pad := 4.0
	w := float64(width)
	h := float64(height)
	step := (w - 2*pad) / float64(len(values)-1)

	// monta polyline
	d := make([]byte,
		0,
		1024)
	for i, v := range values {
		x := pad + float64(i)*step
		y := pad + (h-2*pad)*(1-(v-minV)/(maxV-minV))
		d = append(d,
			[]byte(fmt.Sprintf("%.2f,%.2f ", x, y))...)
	}

	svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="sparkline">
  <polyline fill="none" stroke="currentColor" stroke-width="2" points="%s"/>
</svg>`, width, height, width, height, string(d))
	return os.WriteFile(path,
		[]byte(svg),
		0644)
}

func emptySVG(w, h int) string {
	return fmt.Sprintf(`<?xml version="1.0"?><svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg"></svg>`, w, h)
}
