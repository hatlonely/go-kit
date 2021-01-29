package config

import "context"

func init() {
	RegisterProvider("Memory", NewMemoryProviderWithOptions)
}

type MemoryProvider struct {
	buf []byte
}

func NewMemoryProviderWithOptions(options *MemoryProviderOptions) *MemoryProvider {
	return &MemoryProvider{
		buf: []byte(options.Buffer),
	}
}

type MemoryProviderOptions struct {
	Buffer string
}

func (p *MemoryProvider) Events() <-chan struct{} {
	panic("unimplemented method")
}

func (p *MemoryProvider) Errors() <-chan error {
	panic("unimplemented method")
}

func (p *MemoryProvider) Load() ([]byte, error) {
	return p.buf, nil
}

func (p *MemoryProvider) Dump(buf []byte) error {
	p.buf = buf
	return nil
}

func (p *MemoryProvider) EventLoop(ctx context.Context) error {
	panic("unimplemented method")
}
