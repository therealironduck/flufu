package ai

import "github.com/hybridgroup/yzma/pkg/llama"

type Instance struct {
	llamaCtx llama.Context
	model    llama.Model
	vocab    llama.Vocab
	sampler  llama.Sampler
	ready    chan struct{}
}

func New() *Instance {
	return &Instance{ready: make(chan struct{})}
}

func (a *Instance) Ready() <-chan struct{} {
	return a.ready
}
