package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kubex-ecosystem/ghbex/internal/metrics"
	"github.com/kubex-ecosystem/ghbex/internal/render"
	"github.com/spf13/cobra"
)

type Dora struct {
	DeploymentFrequency float64 `json:"deployment_frequency" yaml:"deployment_frequency"`
	PeriodUnit          string  `json:"period_unit" yaml:"period_unit"`
	LeadTimeP95         float64 `json:"lead_time_p95" yaml:"lead_time_p95"` // hours
	LeadTimeP50         float64 `json:"lead_time_p50" yaml:"lead_time_p50"`
	ChangeFailRate      float64 `json:"change_fail_rate" yaml:"change_fail_rate"`
	MTTR                float64 `json:"mttr" yaml:"mttr"`
}

type Code struct {
	MI             float64   `json:"mi" yaml:"mi"`
	DuplicationPct float64   `json:"duplication_pct" yaml:"duplication_pct"`
	CyclomaticAvg  float64   `json:"cyclomatic_avg" yaml:"cyclomatic_avg"`
	Trend          []float64 `json:"trend" yaml:"trend"`
}

type Community struct {
	FirstReviewP50 float64   `json:"first_review_p50" yaml:"first_review_p50"`
	BusFactor      int       `json:"bus_factor" yaml:"bus_factor"`
	TrendLeadTime  []float64 `json:"trend_lead_time" yaml:"trend_lead_time"`
}

type Input struct {
	SchemaVersion string    `json:"schema_version" yaml:"schema_version"`
	Owner         string    `json:"owner" yaml:"owner"`
	Repo          string    `json:"repo" yaml:"repo"`
	PeriodDays    int       `json:"period_days" yaml:"period_days"`
	Dora          Dora      `json:"dora" yaml:"dora"`
	Code          Code      `json:"code" yaml:"code"`
	Community     Community `json:"community" yaml:"community"`
}

type Files struct {
	SparkCHI  string `json:"spark_chi" yaml:"spark_chi"`
	SparkLead string `json:"spark_lead" yaml:"spark_lead"`
	BadgesMD  string `json:"badges_md_path" yaml:"badges_md_path"`
}

type Output struct {
	CHI    float64  `json:"chi" yaml:"chi"`
	Grade  string   `json:"grade" yaml:"grade"`
	Badges []string `json:"badges_md" yaml:"badges_md"`
	Files  Files    `json:"files" yaml:"files"`
}

func renderScorecard(inPath, outDir *string, width, height *int) {

	b, err := os.ReadFile(*inPath)
	must(err)
	var sc Input
	must(json.Unmarshal(b, &sc))

	chi := metrics.ComputeCHI(sc.Code.MI, sc.Code.DuplicationPct, sc.Code.CyclomaticAvg, metrics.DefaultCHI)
	grade := metrics.GradeFromCHI(chi)

	// gera badges markdown
	badges := render.BuildScorecardBadges(render.Scorecard{
		DLeadP95Hours: sc.Dora.LeadTimeP95,
		DeployFreq:    sc.Dora.DeploymentFrequency,
		DeployUnit:    sc.Dora.PeriodUnit,
		CHI:           chi,
		ReviewP50H:    sc.Community.FirstReviewP50,
		BusFactor:     sc.Community.BusFactor,
	})

	must(os.MkdirAll(*outDir,
		0o755))

	// escreve sparklines
	sparkCHI := filepath.Join(*outDir,
		"sparkline-chi.svg")
	sparkLead := filepath.Join(*outDir,
		"sparkline-leadtime.svg")
	if len(sc.Code.Trend) > 0 {
		must(render.WriteSparklineSVG(sparkCHI, sc.Code.Trend, *width, *height))
	}
	if len(sc.Community.TrendLeadTime) > 0 {
		must(render.WriteSparklineSVG(sparkLead, sc.Community.TrendLeadTime, *width, *height))
	}
	// badges.md
	badgesMD := filepath.Join(*outDir,
		"badges.md")
	f, err := os.Create(badgesMD)
	must(err)
	defer f.Close()
	for _, b := range badges {
		fmt.Fprintln(f, b)
	}

	out := Output{CHI: chi, Grade: grade, Badges: badges}
	out.Files.SparkCHI = sparkCHI
	out.Files.SparkLead = sparkLead
	out.Files.BadgesMD = badgesMD

	ob, _ := json.MarshalIndent(out,
		"",
		"  ")
	fmt.Println(string(ob))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func ScoreCardRootCmd() *cobra.Command {
	var inPath, outDir string
	var width, height int

	short := "Generate a scorecard report"
	long := "Generates a scorecard report based on the provided input JSON file, including sparklines and badges."

	cmd := &cobra.Command{
		Use:     "scorecard",
		Aliases: []string{"sc", "score"},
		Short:   short,
		Long:    long,
		Annotations: GetDescriptions([]string{
			short,
			long,
		}, os.Getenv("GHBEX_HIDE_BANNER") == "true"),
		Run: func(cmd *cobra.Command, args []string) {
			renderScorecard(&inPath, &outDir, &width, &height)
		},
	}

	cmd.Flags().StringVarP(&inPath, "file", "f", "scorecard.json", "input scorecard json")
	cmd.Flags().StringVarP(&outDir, "out", "o", "dist", "output directory")
	cmd.Flags().IntVarP(&width, "width", "w", 220, "sparkline width")
	cmd.Flags().IntVarP(&height, "height", "", 40, "sparkline height")

	return cmd
}
