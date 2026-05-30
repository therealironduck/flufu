package pet

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	bubbleMaxWidth  = 28
	bubbleGap       = 1 // blank columns between bubble right edge and pet left edge
	bubbleSidePad   = 2 // total horizontal padding inside borders ("| " and " |")
	bubbleBorders   = 2 // left and right border characters
	bubbleVertBords = 2 // top and bottom border rows
)

// bubblePos records where a bubble was drawn so clearBubble uses the same coordinates
// regardless of subsequent terminal resizes.
type bubblePos struct {
	startCol   int
	topRow     int
	totalWidth int
	numLines   int
	tailCol    int
	tailRow    int
	valid      bool
}

func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var current strings.Builder

	for _, word := range words {
		wordLen := utf8.RuneCountInString(word)
		switch {
		case current.Len() == 0:
			current.WriteString(word)
		case current.Len()+1+wordLen <= width:
			current.WriteString(" ")
			current.WriteString(word)
		default:
			lines = append(lines, current.String())
			current.Reset()
			current.WriteString(word)
		}
	}

	if current.Len() > 0 {
		lines = append(lines, current.String())
	}

	return lines
}

func bubbleDims(joke string) (lines []string, innerWidth int) {
	lines = wrapText(joke, bubbleMaxWidth)
	for _, l := range lines {
		if n := utf8.RuneCountInString(l); n > innerWidth {
			innerWidth = n
		}
	}

	return
}

func drawBubble(petKey string, frame int, joke string) bubblePos {
	_, startCol, startRow := getWidthAndHeight(petKey, frame)
	lines, innerWidth := bubbleDims(joke)

	contentWidth := innerWidth + bubbleSidePad
	totalWidth := contentWidth + bubbleBorders

	bubbleStartCol := startCol - totalWidth - bubbleGap
	bubbleTopRow := startRow - len(lines) - bubbleVertBords

	pos := bubblePos{
		startCol:   bubbleStartCol,
		topRow:     bubbleTopRow,
		totalWidth: totalWidth,
		numLines:   len(lines),
		tailCol:    bubbleStartCol + totalWidth - 1,
		tailRow:    startRow,
	}

	if bubbleStartCol < 1 || bubbleTopRow < 1 {
		return pos
	}

	pos.valid = true

	var buf strings.Builder
	buf.WriteString("\0337")

	fmt.Fprintf(&buf, "\033[%d;%dH.%s.", bubbleTopRow, bubbleStartCol, strings.Repeat("-", contentWidth))

	for i, line := range lines {
		padding := innerWidth - utf8.RuneCountInString(line)
		fmt.Fprintf(&buf, "\033[%d;%dH| %s%s |", bubbleTopRow+1+i, bubbleStartCol, line, strings.Repeat(" ", padding))
	}

	fmt.Fprintf(&buf, "\033[%d;%dH'%s'", bubbleTopRow+1+len(lines), bubbleStartCol, strings.Repeat("-", contentWidth))

	// Tail character pointing right toward the pet
	fmt.Fprintf(&buf, "\033[%d;%dH\\", startRow, pos.tailCol)

	buf.WriteString("\0338")
	os.Stdout.WriteString(buf.String())

	return pos
}

func clearBubble(pos bubblePos) {
	if !pos.valid {
		return
	}

	blank := strings.Repeat(" ", pos.totalWidth)

	var buf strings.Builder
	buf.WriteString("\0337")

	for i := range pos.numLines + bubbleVertBords {
		fmt.Fprintf(&buf, "\033[%d;%dH%s", pos.topRow+i, pos.startCol, blank)
	}

	// Clear tail
	fmt.Fprintf(&buf, "\033[%d;%dH ", pos.tailRow, pos.tailCol)

	buf.WriteString("\0338")
	os.Stdout.WriteString(buf.String())
}
