package cost

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type languageStat struct {
	Name  string
	Lines int
	Files int
}

type estimate struct {
	Root           string
	Date           string
	RatesSource    string
	Country        string
	CountryCode    string
	Currency       string
	HourlyRateLow  float64
	HourlyRateAvg  float64
	HourlyRateHigh float64
	TotalLines     int
	CodeLines      int
	TestLines      int
	DocLines       int
	ConfigLines    int
	Complexity     float64
	Languages      []languageStat
	BaseHours      float64
	TotalHours     float64
}

type rateTable struct {
	Low  float64
	Avg  float64
	High float64
}

var ignoredDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	".next":        true,
	"coverage":     true,
	".cleo":        true,
}

var extensions = map[string]struct {
	Lang     string
	Category string
	Weight   float64
}{
	".go":    {"Go", "code", 28},
	".ts":    {"TypeScript", "code", 26},
	".tsx":   {"TypeScript", "code", 24},
	".js":    {"JavaScript", "code", 30},
	".jsx":   {"JavaScript", "code", 26},
	".py":    {"Python", "code", 28},
	".rb":    {"Ruby", "code", 30},
	".php":   {"PHP", "code", 30},
	".java":  {"Java", "code", 24},
	".kt":    {"Kotlin", "code", 24},
	".swift": {"Swift", "code", 22},
	".rs":    {"Rust", "code", 18},
	".cpp":   {"C++", "code", 16},
	".cc":    {"C++", "code", 16},
	".c":     {"C", "code", 16},
	".cs":    {"C#", "code", 24},
	".m":     {"Objective-C", "code", 16},
	".mm":    {"Objective-C++", "code", 14},
	".sql":   {"SQL", "code", 20},
	".sh":    {"Shell", "code", 25},
	".yml":   {"YAML", "config", 0},
	".yaml":  {"YAML", "config", 0},
	".json":  {"JSON", "config", 0},
	".toml":  {"TOML", "config", 0},
	".xml":   {"XML", "config", 0},
	".md":    {"Markdown", "docs", 0},
}

func Estimate(args []string) (string, error) {
	root := flagValue(args, "--path")
	if root == "" {
		root = "."
	}
	rateSource := flagValue(args, "--rates-source")
	if rateSource == "" {
		rateSource = "cached"
	}
	format := flagValue(args, "--format")
	if format == "" {
		format = "markdown"
	}
	resolvedRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	statsByLang := map[string]*languageStat{}
	var totalLines, testLines, docLines, configLines int
	var weightedHours float64
	var weightedLines float64
	manifestCount := 0

	err = filepath.WalkDir(resolvedRoot, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			if ignoredDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		meta, ok := extensions[ext]
		if !ok {
			return nil
		}
		count, err := countLines(path)
		if err != nil {
			return nil
		}
		if count == 0 {
			return nil
		}
		totalLines += count
		if isTestFile(path) {
			testLines += count
		}
		switch meta.Category {
		case "docs":
			docLines += count
		case "config":
			configLines += count
		case "code":
			if !isTestFile(path) {
				weightedHours += float64(count) / meta.Weight
				weightedLines += float64(count)
			}
		}
		stat := statsByLang[meta.Lang]
		if stat == nil {
			stat = &languageStat{Name: meta.Lang}
			statsByLang[meta.Lang] = stat
		}
		stat.Lines += count
		stat.Files++
		if isManifest(path) {
			manifestCount++
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	codeLines := totalLines - testLines - docLines - configLines
	if codeLines < 0 {
		codeLines = 0
	}

	if weightedHours == 0 {
		weightedHours = float64(codeLines) / 30.0
	}

	complexity := 1.0
	if len(statsByLang) > 3 {
		complexity += 0.1
	}
	if len(statsByLang) > 6 {
		complexity += 0.1
	}
	if manifestCount > 3 {
		complexity += 0.1
	}
	if weightedLines > 150000 {
		complexity += 0.15
	} else if weightedLines > 50000 {
		complexity += 0.1
	}

	baseHours := weightedHours + (float64(testLines) / 32.0)
	overheadMultiplier := 1.85
	totalHours := baseHours * complexity * overheadMultiplier

	rates, country, err := resolveRates(rateSource, flagValue(args, "--hourly-rate"), flagValue(args, "--country"))
	if err != nil {
		return "", err
	}

	langs := make([]languageStat, 0, len(statsByLang))
	for _, s := range statsByLang {
		langs = append(langs, *s)
	}
	sort.Slice(langs, func(i, j int) bool {
		return langs[i].Lines > langs[j].Lines
	})

	e := estimate{
		Root:           resolvedRoot,
		Date:           time.Now().Format("2006-01-02"),
		RatesSource:    rateSource,
		Country:        country.Name,
		CountryCode:    country.Code2,
		Currency:       country.Currency,
		HourlyRateLow:  rates.Low,
		HourlyRateAvg:  rates.Avg,
		HourlyRateHigh: rates.High,
		TotalLines:     totalLines,
		CodeLines:      codeLines,
		TestLines:      testLines,
		DocLines:       docLines,
		ConfigLines:    configLines,
		Complexity:     complexity,
		Languages:      langs,
		BaseHours:      baseHours,
		TotalHours:     totalHours,
	}
	return render(e, format)
}

func resolveRates(source string, manual string, countryInput string) (rateTable, countryInfo, error) {
	if source == "manual" {
		v, err := strconv.ParseFloat(strings.TrimSpace(manual), 64)
		if err != nil || v <= 0 {
			return rateTable{}, countryInfo{}, fmt.Errorf("--hourly-rate must be a positive number")
		}
		country, err := resolveCountry(countryInput)
		if err != nil {
			return rateTable{}, countryInfo{}, err
		}
		return rateTable{Low: v * 0.85, Avg: v, High: v * 1.25}, country, nil
	}
	if source == "live" {
		rates, country, err := resolveLiveRates(countryInput)
		if err != nil {
			return rateTable{}, countryInfo{}, err
		}
		return rates, country, nil
	}
	country, err := resolveCountry(countryInput)
	if err != nil {
		return rateTable{}, countryInfo{}, err
	}
	// Default cached benchmark table (US, 2025 baseline).
	return rateTable{Low: 95, Avg: 145, High: 220}, country, nil
}

func countLines(path string) (int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	if len(b) == 0 {
		return 0, nil
	}
	lines := 0
	start := 0
	for i := 0; i < len(b); i++ {
		if b[i] == '\n' {
			if i > start {
				if strings.TrimSpace(string(b[start:i])) != "" {
					lines++
				}
			}
			start = i + 1
		}
	}
	if start < len(b) && strings.TrimSpace(string(b[start:])) != "" {
		lines++
	}
	return lines, nil
}

func isTestFile(path string) bool {
	lower := strings.ToLower(path)
	return strings.Contains(lower, "/test/") || strings.Contains(lower, "/tests/") ||
		strings.HasSuffix(lower, "_test.go") || strings.HasSuffix(lower, ".test.ts") ||
		strings.HasSuffix(lower, ".spec.ts") || strings.HasSuffix(lower, ".test.js") || strings.HasSuffix(lower, ".spec.js")
}

func isManifest(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	switch base {
	case "package.json", "go.mod", "cargo.toml", "pom.xml", "build.gradle", "build.gradle.kts", "pyproject.toml", "requirements.txt", "gemfile", "composer.json", "package-lock.json", "pnpm-lock.yaml", "yarn.lock":
		return true
	default:
		return false
	}
}

func render(e estimate, format string) (string, error) {
	switch format {
	case "plain":
		return renderPlain(e), nil
	case "json":
		return renderJSON(e)
	default:
		return renderMarkdown(e), nil
	}
}

func renderMarkdown(e estimate) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Cleo Cost Estimate\n\n")
	fmt.Fprintf(&b, "Analysis Date: %s\n", e.Date)
	fmt.Fprintf(&b, "Project Root: %s\n", e.Root)
	fmt.Fprintf(&b, "Rates Source: %s\n\n", e.RatesSource)
	fmt.Fprintf(&b, "Market Country: %s (%s)", e.Country, e.CountryCode)
	if e.Currency != "" {
		fmt.Fprintf(&b, " | Currency: %s", e.Currency)
	}
	fmt.Fprintf(&b, "\n\n")

	fmt.Fprintf(&b, "## Codebase Metrics\n")
	fmt.Fprintf(&b, "%s\n\n", markdownTable(
		[]string{"Metric", "Value"},
		[][]string{
			{"Total non-empty lines", fmt.Sprintf("%d", e.TotalLines)},
			{"Code lines", fmt.Sprintf("%d", e.CodeLines)},
			{"Test lines", fmt.Sprintf("%d", e.TestLines)},
			{"Docs lines", fmt.Sprintf("%d", e.DocLines)},
			{"Config lines", fmt.Sprintf("%d", e.ConfigLines)},
			{"Complexity multiplier", fmt.Sprintf("%.2fx", e.Complexity)},
		},
		[]bool{false, true},
	))

	fmt.Fprintf(&b, "### Language Mix\n")
	langRows := make([][]string, 0, len(e.Languages))
	for i, lang := range e.Languages {
		if i >= 10 {
			break
		}
		langRows = append(langRows, []string{lang.Name, fmt.Sprintf("%d", lang.Lines), fmt.Sprintf("%d", lang.Files)})
	}
	fmt.Fprintf(&b, "%s\n\n", markdownTable(
		[]string{"Language", "Lines", "Files"},
		langRows,
		[]bool{false, true, true},
	))

	fmt.Fprintf(&b, "## Development Time\n")
	fmt.Fprintf(&b, "%s\n\n", markdownTable(
		[]string{"Metric", "Hours"},
		[][]string{
			{"Base development hours", fmt.Sprintf("%.1f", e.BaseHours)},
			{"Total estimated hours (with overhead)", fmt.Sprintf("%.1f", e.TotalHours)},
		},
		[]bool{false, true},
	))

	low := e.TotalHours * e.HourlyRateLow
	avg := e.TotalHours * e.HourlyRateAvg
	high := e.TotalHours * e.HourlyRateHigh

	fmt.Fprintf(&b, "## Cost Estimate (Engineering Only)\n")
	fmt.Fprintf(&b, "%s\n\n", markdownTable(
		[]string{"Scenario", "Hourly Rate", "Total Cost"},
		[][]string{
			{"Low", fmt.Sprintf("$%.0f/hr", e.HourlyRateLow), fmt.Sprintf("$%.0f", low)},
			{"Average", fmt.Sprintf("$%.0f/hr", e.HourlyRateAvg), fmt.Sprintf("$%.0f", avg)},
			{"High", fmt.Sprintf("$%.0f/hr", e.HourlyRateHigh), fmt.Sprintf("$%.0f", high)},
		},
		[]bool{false, true, true},
	))

	fmt.Fprintf(&b, "## Team-Loaded Cost\n")
	fmt.Fprintf(&b, "%s\n", markdownTable(
		[]string{"Company Type", "Multiplier", "Total Cost"},
		[][]string{
			{"Lean startup", "1.45x", fmt.Sprintf("$%.0f", avg*1.45)},
			{"Growth company", "2.20x", fmt.Sprintf("$%.0f", avg*2.2)},
			{"Enterprise", "2.65x", fmt.Sprintf("$%.0f", avg*2.65)},
		},
		[]bool{false, true, true},
	))

	return b.String()
}

func markdownTable(headers []string, rows [][]string, rightAlign []bool) string {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i := 0; i < len(headers) && i < len(row); i++ {
			if len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}

	var b strings.Builder
	writeTableRow(&b, headers, widths, nil)
	b.WriteString("\n")
	writeTableSeparator(&b, widths, rightAlign)
	for _, row := range rows {
		b.WriteString("\n")
		writeTableRow(&b, row, widths, rightAlign)
	}
	return b.String()
}

func writeTableRow(b *strings.Builder, row []string, widths []int, rightAlign []bool) {
	b.WriteString("|")
	for i := range widths {
		val := ""
		if i < len(row) {
			val = row[i]
		}
		alignRight := i < len(rightAlign) && rightAlign[i]
		if alignRight {
			fmt.Fprintf(b, " %*s |", widths[i], val)
		} else {
			fmt.Fprintf(b, " %-*s |", widths[i], val)
		}
	}
}

func writeTableSeparator(b *strings.Builder, widths []int, rightAlign []bool) {
	b.WriteString("|")
	for i, w := range widths {
		dashes := strings.Repeat("-", w+2)
		alignRight := i < len(rightAlign) && rightAlign[i]
		if alignRight {
			dashes = dashes[:len(dashes)-1] + ":"
		}
		fmt.Fprintf(b, "%s|", dashes)
	}
}

func renderPlain(e estimate) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Cleo Cost Estimate\n")
	fmt.Fprintf(&b, "Analysis Date: %s\n", e.Date)
	fmt.Fprintf(&b, "Project Root: %s\n", e.Root)
	fmt.Fprintf(&b, "Rates Source: %s\n", e.RatesSource)
	fmt.Fprintf(&b, "Market Country: %s (%s)", e.Country, e.CountryCode)
	if e.Currency != "" {
		fmt.Fprintf(&b, " | Currency: %s", e.Currency)
	}
	fmt.Fprintf(&b, "\n\n")
	fmt.Fprintf(&b, "Codebase Metrics\n")
	fmt.Fprintf(&b, "  Total lines: %d\n", e.TotalLines)
	fmt.Fprintf(&b, "  Code lines: %d\n", e.CodeLines)
	fmt.Fprintf(&b, "  Test lines: %d\n", e.TestLines)
	fmt.Fprintf(&b, "  Docs lines: %d\n", e.DocLines)
	fmt.Fprintf(&b, "  Config lines: %d\n", e.ConfigLines)
	fmt.Fprintf(&b, "  Complexity: %.2fx\n\n", e.Complexity)

	fmt.Fprintf(&b, "Language Mix\n")
	for i, lang := range e.Languages {
		if i >= 10 {
			break
		}
		fmt.Fprintf(&b, "  %s: %d lines across %d files\n", lang.Name, lang.Lines, lang.Files)
	}
	fmt.Fprintf(&b, "\n")

	fmt.Fprintf(&b, "Development Time\n")
	fmt.Fprintf(&b, "  Base hours: %.1f\n", e.BaseHours)
	fmt.Fprintf(&b, "  Total hours: %.1f\n\n", e.TotalHours)

	low := e.TotalHours * e.HourlyRateLow
	avg := e.TotalHours * e.HourlyRateAvg
	high := e.TotalHours * e.HourlyRateHigh

	fmt.Fprintf(&b, "Cost Estimate (Engineering)\n")
	fmt.Fprintf(&b, "  Low (%.0f/hr): $%.0f\n", e.HourlyRateLow, low)
	fmt.Fprintf(&b, "  Average (%.0f/hr): $%.0f\n", e.HourlyRateAvg, avg)
	fmt.Fprintf(&b, "  High (%.0f/hr): $%.0f\n\n", e.HourlyRateHigh, high)

	fmt.Fprintf(&b, "Team-Loaded Cost\n")
	fmt.Fprintf(&b, "  Lean startup (1.45x): $%.0f\n", avg*1.45)
	fmt.Fprintf(&b, "  Growth company (2.2x): $%.0f\n", avg*2.2)
	fmt.Fprintf(&b, "  Enterprise (2.65x): $%.0f\n", avg*2.65)
	return b.String()
}

func renderJSON(e estimate) (string, error) {
	low := e.TotalHours * e.HourlyRateLow
	avg := e.TotalHours * e.HourlyRateAvg
	high := e.TotalHours * e.HourlyRateHigh
	payload := map[string]any{
		"title":         "Cleo Cost Estimate",
		"analysis_date": e.Date,
		"project_root":  e.Root,
		"rates": map[string]any{
			"source": e.RatesSource,
			"country": map[string]any{
				"name":     e.Country,
				"code":     e.CountryCode,
				"currency": e.Currency,
			},
			"hourly": map[string]float64{
				"low":     e.HourlyRateLow,
				"average": e.HourlyRateAvg,
				"high":    e.HourlyRateHigh,
			},
		},
		"lines": map[string]int{
			"total":  e.TotalLines,
			"code":   e.CodeLines,
			"tests":  e.TestLines,
			"docs":   e.DocLines,
			"config": e.ConfigLines,
		},
		"complexity_multiplier": e.Complexity,
		"languages":             e.Languages,
		"hours": map[string]float64{
			"base":  e.BaseHours,
			"total": e.TotalHours,
		},
		"cost_engineering": map[string]float64{
			"low":     low,
			"average": avg,
			"high":    high,
		},
		"cost_team_loaded": map[string]float64{
			"lean_startup":   avg * 1.45,
			"growth_company": avg * 2.2,
			"enterprise":     avg * 2.65,
		},
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
