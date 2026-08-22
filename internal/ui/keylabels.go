package ui

import "strings"

// prettyKey turns a raw bubbletea key string into a nicer display form.
func prettyKey(k string) string {
	switch k {
	case "up":
		return "\u2191"
	case "down":
		return "\u2193"
	case "left":
		return "\u2190"
	case "right":
		return "\u2192"
	case "shift+up":
		return "\u21e7\u2191"
	case "shift+down":
		return "\u21e7\u2193"
	case "shift+left":
		return "\u21e7\u2190"
	case "shift+right":
		return "\u21e7\u2192"
	case " ", "space":
		return "space"
	default:
		return k
	}
}

func firstKey(binds []string) string {
	if len(binds) == 0 {
		return "?"
	}
	return prettyKey(binds[0])
}

// combo joins up to the first two binds with a slash, e.g. "+/-" or "\u2191/k".
// Duplicate-looking binds (e.g. " " and "space" both rendering as "space")
// are collapsed into one.
func combo(binds []string) string {
	if len(binds) == 0 {
		return "?"
	}
	first := prettyKey(binds[0])
	if len(binds) == 1 {
		return first
	}
	second := prettyKey(binds[1])
	if second == first {
		return first
	}
	return first + "/" + second
}

// pair combines two bind lists into one label like "\u2191\u2193/jk".
// Used for movement (up+down), where showing both primary and secondary
// binds together reads naturally.
func pair(a, b []string) string {
	primary := firstKey(a) + firstKey(b)
	if len(a) > 1 && len(b) > 1 {
		return primary + "/" + prettyKey(a[1]) + prettyKey(b[1])
	}
	return primary
}

// altPair shows just the primary key of two separate actions side by side,
// e.g. "+/-" for hit/undo. Unlike pair(), it doesn't try to also combine
// secondary alt-keys, since mixing two actions' alt-binds together (e.g.
// "+-/=_") reads as noise rather than a natural combo.
func altPair(a, b []string) string {
	return firstKey(a) + "/" + firstKey(b)
}

func label(binds []string, desc string) string {
	return combo(binds) + " " + desc
}

func joinHelp(parts ...string) string {
	return strings.Join(parts, "  ")
}
