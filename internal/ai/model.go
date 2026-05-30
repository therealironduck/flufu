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

func (a *Instance) Init(ctx context.Context) error {
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

	a.llamaCtx = llamaCtx
	a.model = model
	a.sampler = sampler
	a.vocab = vocab
	close(a.ready)

	<-ctx.Done()

	return nil
}

func initModel(path string) (ctx llama.Context, model llama.Model, vocab llama.Vocab, sampler llama.Sampler, err error) {
	model, err = llama.ModelLoadFromFile(path, llama.ModelDefaultParams())
	if err != nil || model == 0 {
		return
	}

	vocab = llama.ModelGetVocab(model)

	ctx, err = llama.InitFromModel(model, llama.ContextDefaultParams())
	if err != nil || ctx == 0 {
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
