package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

var errNoUserMessage = errors.New("no user message found in transcript")

type transcriptEntry struct {
	Type    string         `json:"type"`
	Message *transcriptMsg `json:"message,omitempty"`
}

type transcriptMsg struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

const readChunkSize = 4096

func ReadNewestUserMessageFromTranscript(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error opening transcript: %w", err)
	}
	defer f.Close()

	size, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return "", fmt.Errorf("error seeking transcript: %w", err)
	}

	var tail []byte
	remaining := size

	for remaining > 0 {
		chunkSize := min(int64(readChunkSize), remaining)

		remaining -= chunkSize
		if _, err := f.Seek(remaining, io.SeekStart); err != nil {
			return "", fmt.Errorf("error seeking transcript: %w", err)
		}

		chunk := make([]byte, chunkSize, chunkSize+int64(len(tail)))
		if _, err := io.ReadFull(f, chunk); err != nil {
			return "", fmt.Errorf("error reading transcript: %w", err)
		}

		tail = append(chunk, tail...)

		lines := bytes.Split(tail, []byte("\n"))
		// Keep the first element as a partial line carried into the next chunk.
		tail = lines[0]

		for i := len(lines) - 1; i >= 1; i-- {
			if msg := extractUserMessage(lines[i]); msg != "" {
				return msg, nil
			}
		}
	}

	// Process any remaining data at the start of the file.
	for line := range bytes.SplitSeq(tail, []byte("\n")) {
		if msg := extractUserMessage(line); msg != "" {
			return msg, nil
		}
	}

	return "", errNoUserMessage
}

func extractUserMessage(line []byte) string {
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return ""
	}

	var entry transcriptEntry
	if err := json.Unmarshal(line, &entry); err != nil {
		return ""
	}

	if entry.Type != "user" || entry.Message == nil || entry.Message.Role != "user" {
		return ""
	}

	var text string
	if err := json.Unmarshal(entry.Message.Content, &text); err == nil {
		return text
	}

	var blocks []contentBlock
	if err := json.Unmarshal(entry.Message.Content, &blocks); err == nil {
		for _, b := range blocks {
			if b.Type == "text" {
				return b.Text
			}
		}
	}

	return ""
}
