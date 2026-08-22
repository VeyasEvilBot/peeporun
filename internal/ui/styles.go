package ui

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func lipglossPlace(s string, w, h int) string {
	if w <= 0 || h <= 0 {
		return s
	}
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, s)
}

var (
	ColorAccent = lipgloss.Color("#FFD400") // yellow, overridable via SetAccentColor
	ColorMuted  = lipgloss.Color("#6c6c6c")
	ColorGood   = lipgloss.Color("#5fd75f")
	ColorBad    = lipgloss.Color("#ff5f5f")
	ColorFg     = lipgloss.Color("#e4e4e4")
	ColorBorder = lipgloss.Color("#3a3a3a")
)

var (
	StyleAppTitle    lipgloss.Style
	StyleGameTitle   lipgloss.Style
	StyleCategory    lipgloss.Style
	StyleHeader      lipgloss.Style
	StyleActiveRow   lipgloss.Style
	StyleActiveNoHit lipgloss.Style
	StyleActiveHit   lipgloss.Style
	StyleBeatenRow   lipgloss.Style
	StyleHitRow      lipgloss.Style
	StyleNormalRow   lipgloss.Style
	StyleTotal       lipgloss.Style
	StyleHelp        lipgloss.Style
	StyleBox         lipgloss.Style
	StyleModal       lipgloss.Style
)

func init() {
	rebuildStyles()
}

// SetAccentColor changes the shared accent color and rebuilds every style
// that depends on it. Call once at startup after loading theme.yaml.
func SetAccentColor(hex string) {
	if hex == "" {
		return
	}
	ColorAccent = lipgloss.Color(hex)
	rebuildStyles()
}

func rebuildStyles() {
	accentText := contrastText(string(ColorAccent))

	StyleAppTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	StyleGameTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	StyleCategory = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Italic(true)

	StyleHeader = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorMuted)

	StyleActiveRow = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(accentText)).
		Background(ColorAccent)

	StyleActiveNoHit = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1a1a1a")).
		Background(lipgloss.Color("#98c379")) // pastel green

	StyleActiveHit = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1a1a1a")).
		Background(lipgloss.Color("#e08a90")) // pastel red

	StyleBeatenRow = lipgloss.NewStyle().
		Foreground(ColorGood)

	StyleHitRow = lipgloss.NewStyle().
		Foreground(ColorBad)

	StyleNormalRow = lipgloss.NewStyle().
		Foreground(ColorFg)

	StyleTotal = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	StyleHelp = lipgloss.NewStyle().
		Foreground(ColorMuted)

	StyleBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAccent).
		Padding(1, 3)

	StyleModal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAccent).
		Padding(1, 2).
		Bold(true)
}

// contrastText picks black or white text depending on the perceived
// brightness of the given "#RRGGBB" hex color, so a highlighted row stays
// readable no matter how dark or bright the chosen accent color is.
func contrastText(hex string) string {
	if len(hex) != 7 || hex[0] != '#' {
		return "#1a1a1a"
	}
	r, err1 := strconv.ParseInt(hex[1:3], 16, 32)
	g, err2 := strconv.ParseInt(hex[3:5], 16, 32)
	b, err3 := strconv.ParseInt(hex[5:7], 16, 32)
	if err1 != nil || err2 != nil || err3 != nil {
		return "#1a1a1a"
	}
	luminance := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	if luminance > 140 {
		return "#1a1a1a"
	}
	return "#f5f5f5"
}

// widestLine returns the visual width (ANSI-aware) of the widest line among
// the given strings. Used to pad highlighted rows out to the full content
// width, instead of the background stopping right after the row's own text.
func widestLine(lines ...string) int {
	max := 0
	for _, l := range lines {
		if w := lipgloss.Width(l); w > max {
			max = w
		}
	}
	return max
}
