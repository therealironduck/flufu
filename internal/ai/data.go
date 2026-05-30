package ai

import "github.com/hybridgroup/yzma/pkg/llama"

type AiInstance struct {
	llamaCtx llama.Context
	model    llama.Model
	vocab    llama.Vocab
	sampler  llama.Sampler
	ready    chan struct{}
}

func New() *AiInstance {
	return &AiInstance{ready: make(chan struct{})}
}

func (a *AiInstance) Ready() <-chan struct{} {
	return a.ready
}
