package pet

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/therealironduck/flufu/internal/ai"
	"golang.org/x/term"
)

const (
	renderInterval  = 50 * time.Millisecond
	nextFrameTicks  = 6
	paddingBottom   = -5
	pet             = "duck"
	jokeInterval    = 10 * time.Second
	jokeDuration    = 4 * time.Second
	jokeFetchBuffer = 30 * time.Second
)

func Render(ctx context.Context, aiInstance *ai.Instance) {
	ticker := time.NewTicker(renderInterval)
	defer ticker.Stop()

	frame := 0
	nextFrameTick := 0

	var (
		currentJoke   string
		lastBubblePos bubblePos
		jokeExpiry    time.Time
		nextJokeAt    time.Time
		jokeCh        = make(chan string, 1)
	)

	for {
		select {
		case <-ctx.Done():
			clearBubble(lastBubblePos)
			clearOverlay(pet, frame)
			return

		case j := <-jokeCh:
			currentJoke = j
			jokeExpiry = time.Now().Add(jokeDuration)
			nextJokeAt = time.Now().Add(jokeInterval)

		case <-ticker.C:
			clearBubble(lastBubblePos)
			clearOverlay(pet, frame)
			lastBubblePos = bubblePos{}

			nextFrameTick = (nextFrameTick + 1) % nextFrameTicks
			if nextFrameTick == 0 {
				frame = (frame + 1) % len(pets[pet])
			}

			if currentJoke != "" && time.Now().After(jokeExpiry) {
				currentJoke = ""
			}

			select {
			case <-aiInstance.Ready():
				if nextJokeAt.IsZero() || time.Now().After(nextJokeAt) {
					nextJokeAt = time.Now().Add(jokeDuration + jokeInterval + jokeFetchBuffer)
					go func() {
						joke, err := aiInstance.Joke()
						if err != nil {
							return
						}
						select {
						case jokeCh <- strings.TrimSpace(joke):
						default:
						}
					}()
				}
			default:
			}

			drawOverlay(pet, frame)
			if currentJoke != "" {
				lastBubblePos = drawBubble(pet, frame, currentJoke)
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
