package ai

import "github.com/hybridgroup/yzma/pkg/llama"

func buildPrompt(model llama.Model, userMessage string) string {
	template := llama.ModelChatTemplate(model, "")
	if template == "" {
		template = "chatml"
	}

	messages := []llama.ChatMessage{
		llama.NewChatMessage("system", "You are a helpful assistant. Answer concisely and accurately."),
		llama.NewChatMessage("user", userMessage),
	}

	const bufSize = 4096
	buf := make([]byte, bufSize)
	n := llama.ChatApplyTemplate(template, messages, true, buf)
	return string(buf[:n])
}

func (ai *AiInstance) MakePrompt(message string) string {
	return buildPrompt(ai.model, message)
}
