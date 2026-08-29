package overlay

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

// Split is one row of overlay data, mirroring the TUI's tracker row.
type Split struct {
	Name   string
	Hits   int
	PB     int
	Beaten bool
	Active bool // true if this is the split the cursor is currently on
}

// Data is a full snapshot of what the tracker currently shows.
type Data struct {
	Game      string
	Category  string
	Splits    []Split
	TotalHits int
	TotalPB   int
	ShowPB    bool
}

// Render builds a self-contained HTML document for OBS Browser Source
// (Local File), auto-refreshing every refreshSeconds to pick up changes.
// accent is a "#RRGGBB" hex color used for the panel border, titles, and
// total row, matching the TUI's accent color. When d.ShowPB is false, the
// PB column is omitted entirely (matching the TUI's own show_pb setting).
func Render(d Data, refreshSeconds float64, accent string) string {
	if accent == "" {
		accent = "#FFD400"
	}

	pbCell := func(v int) string {
		if !d.ShowPB {
			return ""
		}
		return fmt.Sprintf(`<span class="pb">%d</span>`, v)
	}

	var rows strings.Builder
	for _, s := range d.Splits {
		class := rowStyle(s)
		rows.WriteString(fmt.Sprintf(
			`<div class="row %s"><span class="name">%s</span><span class="hits">%d</span>%s</div>`+"\n",
			class, html.EscapeString(s.Name), s.Hits, pbCell(s.PB),
		))
	}

	columns := "auto auto"
	headerPB := ""
	totalPB := ""
	if d.ShowPB {
		columns = "auto auto auto"
		headerPB = `<span class="pb">PB</span>`
		totalPB = pbCell(d.TotalPB)
	}

	refresh := strconv.FormatFloat(refreshSeconds, 'f', -1, 64)

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta http-equiv="refresh" content="%s">
<title>peepoRun overlay</title>
<style>
  html, body {
    margin: 0;
    padding: 0;
    background: transparent;
    font-family: 'Consolas', 'Courier New', monospace;
    color: #e4e4e4;
  }
  .panel {
    display: inline-grid;
    grid-template-columns: %[8]s;
    row-gap: 2px;
    column-gap: 0;
    background: rgba(20, 20, 20, 0.72);
    border: 2px solid %[2]s;
    border-radius: 10px;
    padding: 14px 20px;
  }
  .game, .category {
    grid-column: 1 / -1;
  }
  .game {
    color: %[2]s;
    font-weight: bold;
    font-size: 20px;
  }
  .category {
    color: %[2]s;
    font-style: italic;
    font-size: 14px;
    margin-bottom: 8px;
  }
  /* header/row/total wrapper divs don't render a box themselves - their
     span children become direct grid items, so every row shares the same
     column widths and a highlighted row's background fills the full
     column instead of stopping right after its own text. */
  .header, .row, .total {
    display: contents;
  }
  .name, .hits, .pb {
    line-height: 20px;
  }
  .name { padding: 2px 12px 2px 0; }
  .hits, .pb { text-align: right; padding: 2px 4px; }
  .header > * {
    color: #9a9a9a;
    font-weight: bold;
    font-size: 13px;
    border-bottom: 1px solid #444;
    padding-bottom: 4px;
  }
  .row.clear > * { color: #5fd75f; }
  .row.hit > * { color: #ff5f5f; }
  .row.active-clear > * { background: #98c379; color: #1a1a1a; }
  .row.active-hit > * { background: #e08a90; color: #1a1a1a; }
  .row.active-clear > .name, .row.active-hit > .name { border-radius: 4px 0 0 4px; }
  /* :last-child (rather than a fixed .pb selector) so the right-side
     rounding/bleed lands on whichever column is actually last - .pb
     normally, or .hits when PB is hidden. */
  .row.active-clear > *:last-child, .row.active-hit > *:last-child { border-radius: 0 4px 4px 0; }
  /* tiny background bleed past the row's own edges - box-shadow is
     paint-only, so it doesn't shift the actual text position at all. */
  .row.active-clear > .name { box-shadow: -4px 0 0 0 #98c379; }
  .row.active-hit > .name { box-shadow: -4px 0 0 0 #e08a90; }
  .row.active-clear > *:last-child { box-shadow: 4px 0 0 0 #98c379; }
  .row.active-hit > *:last-child { box-shadow: 4px 0 0 0 #e08a90; }
  .total > * {
    color: %[2]s;
    font-weight: bold;
    border-top: 1px solid #444;
    padding-top: 4px;
  }
</style>
</head>
<body>
<div class="panel">
  <div class="game">%[3]s</div>
  <div class="category">%[4]s</div>
  <div class="header"><span class="name">Split</span><span class="hits">Hits</span>%[9]s</div>
%[5]s  <div class="total"><span class="name">Total</span><span class="hits">%[6]d</span>%[10]s</div>
</div>
</body>
</html>
`, refresh, accent, html.EscapeString(orDefault(d.Game, "(untitled)")), html.EscapeString(d.Category), rows.String(), d.TotalHits, d.TotalPB, columns, headerPB, totalPB)
}

func rowStyle(s Split) string {
	switch {
	case s.Beaten && s.Hits > 0:
		return "hit"
	case s.Beaten:
		return "clear"
	case s.Active && s.Hits > 0:
		return "active-hit"
	case s.Active:
		return "active-clear"
	default:
		return ""
	}
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
