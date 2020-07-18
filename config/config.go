package config

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

func NewConfigWithBaseFile(filename string) (*Config, error) {
	bd, _ := NewDecoder("json")
	bp, err := NewLocalProvider(filename)
	if err != nil {
		return nil, err
	}
	bc, err := NewConfig(bd, bp, nil)
	if err != nil {
		return nil, err
	}

	cd, err := NewDecoderWithConfig(bc.Sub("Decoder"))
	if err != nil {
		return nil, err
	}
	cp, err := NewProviderWithConfig(bc.Sub("Provider"))
	if err != nil {
		return nil, err
	}
	cc, err := NewCipherWithConfig(bc.Sub("Cipher"))
	if err != nil {
		return nil, err
	}
	return NewConfig(cd, cp, cc)
}

func NewConfig(decoder Decoder, provider Provider, cipher Cipher) (*Config, error) {
	buf, err := provider.Load()
	if err != nil {
		return nil, err
	}
	storage, err := decoder.Decode(buf)
	if err != nil {
		return nil, err
	}
	if err := storage.Decrypt(cipher); err != nil {
		return nil, err
	}
	return &Config{
		provider:     provider,
		storage:      storage,
		decoder:      decoder,
		log:          logrus.New(), // todo
		cipher:       cipher,
		itemHandlers: map[string][]OnChangeHandler{},
	}, nil
}

type Config struct {
	provider Provider
	storage  *Storage
	decoder  Decoder
	cipher   Cipher

	itemHandlers map[string][]OnChangeHandler
	log          Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (c *Config) GetComponent() (Provider, *Storage, Decoder, Cipher) {
	return c.provider, c.storage, c.decoder, c.cipher
}

func (c *Config) Get(key string) (interface{}, error) {
	return c.storage.Get(key)
}

func (c *Config) UnsafeSet(key string, val interface{}) error {
	return c.storage.Set(key, val)
}

func (c *Config) Unmarshal(v interface{}) error {
	return c.storage.Unmarshal(v)
}

func (c *Config) Sub(key string) *Config {
	storage := c.storage.Sub(key)
	if storage == nil {
		return nil
	}

	return &Config{storage: storage}
}

func (c *Config) SubArr(key string) ([]*Config, error) {
	vs, err := c.storage.SubArr(key)
	if err != nil {
		return nil, err
	}
	var configs []*Config
	for _, storage := range vs {
		configs = append(configs, &Config{storage: storage})
	}
	return configs, nil
}

func (c *Config) SubMap(key string) (map[string]*Config, error) {
	kvs, err := c.storage.SubMap(key)
	if err != nil {
		return nil, err
	}
	configMap := map[string]*Config{}
	for k, v := range kvs {
		configMap[k] = &Config{storage: v}
	}
	return configMap, nil
}

func (c *Config) Stop() {
	c.cancel()
	c.wg.Wait()
}

func (c *Config) Watch() error {
	traveled := map[string]bool{}
	_ = c.storage.Travel(func(key string, val interface{}) error {
		for key != "" {
			if !traveled[key] {
				traveled[key] = true
				for _, handler := range c.itemHandlers[key] {
					handler(c.Sub(""))
				}
			}

			var err error
			_, key, err = getLastToken(key)
			if err != nil {
				c.log.Warnf("get token failed. key [%v]", key)
			}
		}

		for _, handler := range c.itemHandlers[""] {
			handler(c.Sub(""))
		}
		return nil
	})

	c.ctx, c.cancel = context.WithCancel(context.Background())
	if err := c.provider.EventLoop(c.ctx); err != nil {
		return err
	}

	c.wg.Add(1)
	go func() {
		ticker := time.NewTicker(300 * time.Second)
		defer ticker.Stop()
	out:
		for {
			select {
			case <-c.provider.Events():
				for len(c.provider.Events()) != 0 {
					<-c.provider.Events()
				}
				buf, err := c.provider.Load()
				if err != nil {
					c.log.Warnf("provider load failed. err: [%v]", err)
					continue
				}
				storage, err := c.decoder.Decode(buf)
				if err != nil {
					c.log.Warnf("decoder decode failed. err: [%v]", err)
					continue
				}
				if err := storage.Decrypt(c.cipher); err != nil {
					c.log.Warnf("storage decrypt failed. err: [%v]", err)
					continue
				}
				diffKeys, err := storage.Diff(*c.storage)
				if err != nil {
					c.log.Warnf("storage diff failed. err: [%v]", err)
					continue
				}
				c.storage = storage
				c.log.Infof("reload config success. storage: %v", c.storage)

				traveled := map[string]bool{}
				for _, key := range diffKeys {
					for key != "" {
						if !traveled[key] {
							traveled[key] = true
							for _, handler := range c.itemHandlers[key] {
								handler(c.Sub(""))
							}
						}
						_, key, err = getLastToken(key)
						if err != nil {
							c.log.Warnf("get token failed. key [%v]", key)
						}
					}
				}
				for _, handler := range c.itemHandlers[""] {
					handler(c.Sub(""))
				}
			case err := <-c.provider.Errors():
				c.log.Warnf("provider error [%v]", err)
			case <-ticker.C:
				c.log.Infof("tick")
			case <-c.ctx.Done():
				c.log.Infof("stop watch. exit")
				break out
			}
		}
		c.wg.Done()
	}()

	return nil
}

type OnChangeHandler func(conf *Config)

func (c *Config) AddOnChangeHandler(handler OnChangeHandler) {
	c.AddOnItemChangeHandler("", handler)
}

func (c *Config) AddOnItemChangeHandler(key string, handler OnChangeHandler) {
	c.itemHandlers[key] = append(c.itemHandlers[key], handler)
}
