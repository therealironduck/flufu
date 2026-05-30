package pet

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	renderInterval = 50 * time.Millisecond
	nextFrameTicks = 6
	paddingBottom  = -5
	pet            = "duck"
)

func Render(ctx context.Context) {
	ticker := time.NewTicker(renderInterval)
	defer ticker.Stop()

	frame := 0
	nextFrameTick := 0

	for {
		select {
		case <-ctx.Done():
			clearOverlay(pet, frame)
			return

		case <-ticker.C:
			clearOverlay(pet, frame)

			nextFrameTick = (nextFrameTick + 1) % nextFrameTicks
			if nextFrameTick == 0 {
				frame = (frame + 1) % len(pets[pet])
			}

			drawOverlay(pet, frame)
		}
	}
}

func getWidthAndHeight(pet string, frame int) (w, startCol, startRow int) {
	cols, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return
	}

	w = 0
	for _, l := range pets[pet][frame] {
		if utf8.RuneCountInString(l) > w {
			w = utf8.RuneCountInString(l)
		}
	}
	h := len(pets[pet][frame])

	startCol = cols - w - 1
	startRow = rows - h - 1 + paddingBottom

	return
}

func clearOverlay(pet string, frame int) {
	petW, startCol, startRow := getWidthAndHeight(pet, frame)

	var buf strings.Builder
	buf.WriteString("\0337") // save cursor
	blank := fmt.Sprintf("%*s", petW, "")

	for i := range pets[pet][frame] {
		fmt.Fprintf(&buf, "\033[%d;%dH", startRow+i, startCol)
		buf.WriteString(blank)
	}

	buf.WriteString("\0338") // restore cursor
	os.Stdout.WriteString(buf.String())
}

func drawOverlay(pet string, frame int) {
	_, startCol, startRow := getWidthAndHeight(pet, frame)

	var buf strings.Builder
	buf.WriteString("\0337") // save cursor

	for i, line := range pets[pet][frame] {
		fmt.Fprintf(&buf, "\033[%d;%dH", startRow+i, startCol)
		buf.WriteString(line)
	}

	buf.WriteString("\0338") // restore cursor
	os.Stdout.WriteString(buf.String())
}
