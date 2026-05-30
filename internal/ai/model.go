package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
)

const (
	modelURL      = "https://huggingface.co/Qwen/Qwen2.5-1.5B-Instruct-GGUF/resolve/main/qwen2.5-1.5b-instruct-q4_k_m.gguf"
	modelFileName = "qwen2.5-1.5b-instruct-q4_k_m.gguf"
	maxTokens     = int32(1024)
)

func Init(ctx context.Context) error {
	if err := initLlama(); err != nil {
		return fmt.Errorf("failed to init llama: %w", err)
	}

	defer llama.Close()

	path, err := loadModel(ctx)
	if err != nil {
		return fmt.Errorf("failed to load model: %w", err)
	}

	llamaCtx, model, vocab, sampler, err := initModel(path)
	if err != nil {
		return fmt.Errorf("failed to init model: %w", err)
	}

	defer func() {
		if err := llama.ModelFree(model); err != nil {
			fmt.Fprintf(os.Stderr, "failed to free model: %v", err)
		}

		if err := llama.Free(llamaCtx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to free llma context: %v", err)
		}

		llama.SamplerFree(sampler)
	}()

	prompt := buildPrompt(model, "What are these yellow animals called that live in ponds and swim?")
	response, err := generate(llamaCtx, vocab, sampler, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}

	fmt.Printf("\n\n%s\n\n", response)

	return nil
}

func initModel(path string) (ctx llama.Context, model llama.Model, vocab llama.Vocab, sampler llama.Sampler, err error) {
	model, err = llama.ModelLoadFromFile(path, llama.ModelDefaultParams())
	if err != nil || model == 0 {
		fmt.Printf("\n\n-----\n\n%v\n\n------\n", path)
		return
	}

	vocab = llama.ModelGetVocab(model)

	ctx, err = llama.InitFromModel(model, llama.ContextDefaultParams())
	if err != nil || ctx == 0 {
		fmt.Printf("HELLLOOO2")
		return
	}

	sp := llama.DefaultSamplerParams()
	sp.PenaltyRepeat = 1.3
	samplerTypes := []llama.SamplerType{
		llama.SamplerTypePenalties,
		llama.SamplerTypeTopK,
		llama.SamplerTypeTopP,
		llama.SamplerTypeTemperature,
	}
	sampler = llama.NewSampler(model, samplerTypes, sp)

	return
}

/**
 * ai.Init() <-- this should run in a go func and keep the model alive until ctx cancel
 * ---
 * ai.Respond("my cool prompt") -- string / error
 *
 * I was thinking we are using
 */
