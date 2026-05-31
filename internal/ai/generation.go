package ai

import (
	"fmt"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func generate(ctx llama.Context, vocab llama.Vocab, sampler llama.Sampler, prompt string) (string, error) {
	mem, err := llama.GetMemory(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get memory: %w", err)
	}

	if err := llama.MemoryClear(mem, true); err != nil {
		return "", fmt.Errorf("failed to clear memory: %w", err)
	}

	tokens := llama.Tokenize(vocab, prompt, true, true)
	batch := llama.BatchGetOne(tokens)

	var result []byte

	for pos := int32(0); pos < maxTokens; pos += batch.NTokens {
		if _, err := llama.Decode(ctx, batch); err != nil {
			return string(result), err
		}

		token := llama.SamplerSample(sampler, ctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			break
		}

		const bufSize = 36
		buf := make([]byte, bufSize)
		n := llama.TokenToPiece(vocab, token, buf, 0, true)
		result = append(result, buf[:n]...)

		batch = llama.BatchGetOne([]llama.Token{token})
	}

	return string(result), nil
}

func (a *Instance) Generate(prompt string) (string, error) {
	return generate(a.llamaCtx, a.vocab, a.sampler, prompt)
}
