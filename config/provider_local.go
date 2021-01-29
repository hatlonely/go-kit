package config

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

func NewLocalProvider(filename string) (*LocalProvider, error) {
	return NewLocalProviderWithOptions(&LocalProviderOptions{Filename: filename})
}

func NewLocalProviderWithOptions(options *LocalProviderOptions) (*LocalProvider, error) {
	abs, err := filepath.Abs(options.Filename)
	if err != nil {
		return nil, err
	}

	return &LocalProvider{
		abs:    abs,
		events: make(chan struct{}, 10),
		errors: make(chan error, 10),
	}, nil
}

func NewLocalProviderWithConfig(cfg *Config, opts ...refx.Option) (*LocalProvider, error) {
	var options LocalProviderOptions
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
	}
	return NewLocalProviderWithOptions(&options)
}

type LocalProvider struct {
	abs    string
	events chan struct{}
	errors chan error
}

func (p *LocalProvider) Events() <-chan struct{} {
	return p.events
}

func (p *LocalProvider) Errors() <-chan error {
	return p.errors
}

func (p *LocalProvider) Load() ([]byte, error) {
	fp, err := os.Open(p.abs)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return ioutil.ReadAll(fp)
}

func (p *LocalProvider) Dump(buf []byte) error {
	return ioutil.WriteFile(p.abs, buf, 0644)
}

func (p *LocalProvider) EventLoop(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := watcher.Add(filepath.Dir(p.abs)); err != nil {
		return err
	}

	go func() {
	out:
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
					continue
				}
				if event.Name != p.abs {
					continue
				}
				for len(watcher.Events) != 0 {
					<-watcher.Events
				}
				p.events <- struct{}{}
			case err := <-watcher.Errors:
				p.errors <- err
			case <-ctx.Done():
				break out
			}
		}
	}()

	return nil
}

type LocalProviderOptions struct {
	Filename string
}
