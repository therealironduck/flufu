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
	msgDuration    = 4 * time.Second
)

func Render(ctx context.Context, msgCh <-chan string) {
	ticker := time.NewTicker(renderInterval)
	defer ticker.Stop()

	frame := 0
	nextFrameTick := 0

	var (
		currentMsg    string
		lastBubblePos bubblePos
		msgExpiry     = time.Now().Add(msgDuration)
	)

	currentMsg = "Ahoy!"

	for {
		select {
		case <-ctx.Done():
			clearBubble(lastBubblePos)
			clearOverlay(pet, frame)
			return

		case msg := <-msgCh:
			currentMsg = msg
			msgExpiry = time.Now().Add(msgDuration)

		case <-ticker.C:
			clearBubble(lastBubblePos)
			clearOverlay(pet, frame)
			lastBubblePos = bubblePos{}

			nextFrameTick = (nextFrameTick + 1) % nextFrameTicks
			if nextFrameTick == 0 {
				frame = (frame + 1) % len(pets[pet])
			}

			if currentMsg != "" && time.Now().After(msgExpiry) {
				currentMsg = ""
			}

			drawOverlay(pet, frame)
			if currentMsg != "" {
				lastBubblePos = drawBubble(pet, frame, currentMsg)
			}
		}
	}
}

func getWidthAndHeight(petKey string, frame int) (w, startCol, startRow int) {
	cols, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return
	}

	w = 0
	for _, l := range pets[petKey][frame] {
		if utf8.RuneCountInString(l) > w {
			w = utf8.RuneCountInString(l)
		}
	}
	h := len(pets[petKey][frame])

	startCol = cols - w - 1
	startRow = rows - h - 1 + paddingBottom

	return
}

func clearOverlay(petKey string, frame int) {
	petW, startCol, startRow := getWidthAndHeight(petKey, frame)

	var buf strings.Builder
	buf.WriteString("\0337")
	blank := fmt.Sprintf("%*s", petW, "")

	for i := range pets[petKey][frame] {
		fmt.Fprintf(&buf, "\033[%d;%dH", startRow+i, startCol)
		buf.WriteString(blank)
	}

	buf.WriteString("\0338")
	os.Stdout.WriteString(buf.String())
}

func drawOverlay(petKey string, frame int) {
	_, startCol, startRow := getWidthAndHeight(petKey, frame)

	var buf strings.Builder
	buf.WriteString("\0337")

	for i, line := range pets[petKey][frame] {
		fmt.Fprintf(&buf, "\033[%d;%dH", startRow+i, startCol)
		buf.WriteString(line)
	}

	buf.WriteString("\0338")
	os.Stdout.WriteString(buf.String())
}
