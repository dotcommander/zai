package cmd

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"

	"github.com/dotcommander/zai/internal/app"
)

type usageDisplay struct {
	Badge, Heading               string
	Types                        map[string]string
	Labels                       map[string]string
	Filled, Empty, PercentMark   string
	GroupSeparator, DurationHour string
	DurationMinute, ResetFormat  string
	Indent, Gap                  string
	BarWidth                     int
}

func currentUsageDisplay() usageDisplay {
	get := func(key string) string { return viper.GetString("quota.display." + key) }
	return usageDisplay{
		Badge: get("badge"), Heading: viper.GetString("quota.heading"),
		Types: map[string]string{
			"TIME_LIMIT":   get("types.time_limit"),
			"TOKENS_LIMIT": get("types.tokens_limit"),
			"CREDIT_LIMIT": get("types.credit_limit"),
		},
		Labels: map[string]string{
			"used":      viper.GetString("quota.labels.used"),
			"limit":     viper.GetString("quota.labels.limit"),
			"remaining": viper.GetString("quota.labels.remaining"),
			"duration":  viper.GetString("quota.labels.duration"),
			"reset":     viper.GetString("quota.labels.reset"),
		},
		Filled: get("filled"), Empty: get("empty"), PercentMark: get("percent_mark"),
		GroupSeparator: get("group_separator"), DurationHour: get("duration_hour"),
		DurationMinute: get("duration_minute"), ResetFormat: get("reset_format"),
		Indent: get("indent"), Gap: get("gap"), BarWidth: viper.GetInt("quota.display.bar_width"),
	}
}

// writeUsagePretty renders only the terminal presentation; no terminal state is changed.
func writeUsagePretty(out io.Writer, windows []app.QuotaWindow, display usageDisplay, width int, colored bool) error {
	renderer := lipgloss.NewRenderer(out)
	style := func(base lipgloss.Style) lipgloss.Style {
		if !colored {
			return renderer.NewStyle()
		}
		return base.Renderer(renderer)
	}
	title := style(theme.Title)
	section := style(theme.Section)
	percentage := style(theme.Info)
	barUsed := style(theme.Command)
	barLeft := style(theme.Divider)
	dim := style(theme.Dim)

	var b strings.Builder
	badge := title.Render(display.Badge)
	heading := dim.Render(display.Heading)
	if width > 0 && lipgloss.Width(badge)+lipgloss.Width(display.Gap)+lipgloss.Width(heading) > width {
		fmt.Fprintln(&b, badge)
		fmt.Fprintln(&b, heading)
	} else {
		fmt.Fprintln(&b, badge+display.Gap+heading)
	}
	for _, window := range windows {
		name := display.Types[window.Type]
		if name == "" {
			name = window.Type
		}
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, section.Render(name))

		label := fmt.Sprintf("%g%s %s", window.Percentage, display.PercentMark, display.Labels["used"])
		line := display.Indent + percentage.Render(label)
		available := display.BarWidth
		if width > 0 {
			remaining := width - lipgloss.Width(display.Indent) - lipgloss.Width(label) - lipgloss.Width(display.Gap)
			if remaining < available {
				available = remaining
			}
		}
		if available > 0 && display.Filled != "" && display.Empty != "" {
			percent := math.Max(0, math.Min(window.Percentage, 100))
			filled := int(percent * float64(available) / 100)
			if percent > 0 && filled == 0 {
				filled = 1
			}
			line += display.Gap + barUsed.Render(strings.Repeat(display.Filled, filled)) + barLeft.Render(strings.Repeat(display.Empty, available-filled))
		}
		fmt.Fprintln(&b, line)

		details := usageWindowDetails(window, display)
		// Keep counts together when they fit; wrap complete values, never truncate them.
		line = display.Indent
		for _, detail := range details {
			if line != display.Indent && width > 0 && lipgloss.Width(line)+lipgloss.Width(display.Gap)+lipgloss.Width(detail) > width {
				fmt.Fprintln(&b, dim.Render(line))
				line = display.Indent
			}
			if line != display.Indent {
				line += display.Gap
			}
			line += detail
		}
		if line != display.Indent {
			fmt.Fprintln(&b, dim.Render(line))
		}
	}
	_, err := io.WriteString(out, b.String())
	return err
}

func usageWindowDetails(window app.QuotaWindow, display usageDisplay) []string {
	var details []string
	for _, item := range []struct {
		key   string
		value *float64
	}{
		{"used", window.Used}, {"limit", window.Limit}, {"remaining", window.Remaining},
	} {
		if item.value != nil {
			details = append(details, display.Labels[item.key]+" "+groupUsageNumber(*item.value, display.GroupSeparator))
		}
	}
	if window.DurationMinutes != nil {
		minutes := *window.DurationMinutes
		value := fmt.Sprintf("%d%s", minutes, display.DurationMinute)
		if minutes != 0 && minutes%60 == 0 {
			value = fmt.Sprintf("%d%s", minutes/60, display.DurationHour)
		}
		details = append(details, display.Labels["duration"]+" "+value)
	}
	if window.ResetsAt != nil {
		details = append(details, display.Labels["reset"]+" "+window.ResetsAt.UTC().Format(display.ResetFormat))
	}
	return details
}

func groupUsageNumber(value float64, separator string) string {
	s := strconv.FormatFloat(value, 'f', -1, 64)
	integer := s
	fraction := ""
	if i := strings.IndexAny(s, ".eE"); i >= 0 {
		integer, fraction = s[:i], s[i:]
	}
	start := 0
	if strings.HasPrefix(integer, "-") {
		start = 1
	}
	if len(integer)-start <= 3 {
		return s
	}
	var b strings.Builder
	b.WriteString(integer[:start])
	for i := start; i < len(integer); i++ {
		if i > start && (len(integer)-i)%3 == 0 {
			b.WriteString(separator)
		}
		b.WriteByte(integer[i])
	}
	b.WriteString(fraction)
	return b.String()
}
