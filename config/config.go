package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
)

type Options struct {
	Cipher   CipherOptions
	Decoder  DecoderOptions
	Provider ProviderOptions
}

func NewConfigWithOptions(options *Options) (*Config, error) {
	cipher, err := NewCipherWithOptions(&options.Cipher)
	if err != nil {
		return nil, err
	}
	decoder, err := NewDecoderWithOptions(&options.Decoder)
	if err != nil {
		return nil, err
	}
	provider, err := NewProviderWithOptions(&options.Provider)
	if err != nil {
		return nil, err
	}
	return NewConfig(decoder, provider, cipher)
}

func NewConfigWithBaseFile(filename string, opts ...refx.Option) (*Config, error) {
	cfg, err := NewConfigWithSimpleFile(filename)
	if err != nil {
		return nil, err
	}

	var options Options
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, errors.Wrap(err, "cfg.Unmarshal failed.")
	}
	return NewConfigWithOptions(&options)
}

type SimpleFileOptions struct {
	FileType  string
	DecodeKey string
}

var defaultSimpleFileOptions = SimpleFileOptions{
	FileType: "Json",
}

type SimpleFileOption func(options *SimpleFileOptions)

func WithSimpleFileType(fileType string) SimpleFileOption {
	return func(options *SimpleFileOptions) {
		options.FileType = fileType
	}
}

func WithSimpleFileKey(key string) SimpleFileOption {
	return func(options *SimpleFileOptions) {
		options.DecodeKey = key
	}
}

func NewConfigWithSimpleFile(filename string, opts ...SimpleFileOption) (*Config, error) {
	simpleFileOptions := defaultSimpleFileOptions
	for _, opt := range opts {
		opt(&simpleFileOptions)
	}

	options := &Options{
		Decoder: DecoderOptions{
			Type: "Json",
		},
		Provider: ProviderOptions{
			Type: "Local",
			LocalProvider: LocalProviderOptions{
				Filename: filename,
			},
		},
	}
	if simpleFileOptions.DecodeKey != "" {
		options.Cipher = CipherOptions{
			Type: "AES",
			AESCipher: AESCipherOptions{
				Key: simpleFileOptions.DecodeKey,
			},
		}
	}
	if simpleFileOptions.FileType != "" {
		options.Decoder = DecoderOptions{
			Type: simpleFileOptions.FileType,
		}
	}

	return NewConfigWithOptions(options)
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
		log:          StdoutLogger{},
		cipher:       cipher,
		itemHandlers: map[string][]OnChangeHandler{},
	}, nil
}

type Config struct {
	parent *Config
	prefix string

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
	if c.parent != nil {
		return c.parent.GetComponent()
	}

	return c.provider, c.storage, c.decoder, c.cipher
}

func (c *Config) Get(key string) (interface{}, bool) {
	if c.parent != nil {
		return c.parent.Get(prefixAppendKey(c.prefix, key))
	}

	val, err := c.GetE(key)
	if err != nil {
		return nil, false
	}
	return val, true
}

func (c *Config) GetE(key string) (interface{}, error) {
	if c.parent != nil {
		return c.parent.GetE(prefixAppendKey(c.prefix, key))
	}
	return c.storage.Get(key)
}

func (c *Config) UnsafeSet(key string, val interface{}) error {
	if c.parent != nil {
		return c.parent.UnsafeSet(prefixAppendKey(c.prefix, key), val)
	}
	return c.storage.Set(key, val)
}

func (c *Config) Unmarshal(v interface{}, opts ...refx.Option) error {
	if c.parent != nil {
		return c.parent.storage.Sub(c.prefix).Unmarshal(v, opts...)
	}
	return c.storage.Unmarshal(v, opts...)
}

func (c *Config) Sub(key string) *Config {
	return &Config{parent: c, prefix: key}
}

func (c *Config) SubArr(key string) ([]*Config, error) {
	vs, err := c.storage.SubArr(key)
	if err != nil {
		return nil, err
	}
	var configs []*Config
	for i := range vs {
		configs = append(configs, &Config{parent: c, prefix: prefixAppendIdx(key, i)})
	}
	return configs, nil
}

func (c *Config) SubMap(key string) (map[string]*Config, error) {
	kvs, err := c.storage.SubMap(key)
	if err != nil {
		return nil, err
	}
	configMap := map[string]*Config{}
	for k := range kvs {
		configMap[k] = &Config{parent: c, prefix: prefixAppendKey(key, k)}
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
					handler(c.Sub(key))
				}
			}

			var err error
			_, key, err = getLastToken(key)
			if err != nil {
				c.log.Warnf("get token failed. key [%v]", key)
			}
		}

		for _, handler := range c.itemHandlers[""] {
			handler(c)
		}
		return nil
	})

	c.ctx, c.cancel = context.WithCancel(context.Background())
	if err := c.provider.EventLoop(c.ctx); err != nil {
		return err
	}

	c.wg.Add(1)
	go func() {
		c.log.Infof("start watch")
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
				c.log.Infof("reload config success. storage: %v", strx.JsonMarshal(c.storage.Interface()))

				traveled := map[string]bool{}
				for _, key := range diffKeys {
					for key != "" {
						if !traveled[key] {
							traveled[key] = true
							for _, handler := range c.itemHandlers[key] {
								handler(c.Sub(key))
							}
						}
						_, key, err = getLastToken(key)
						if err != nil {
							c.log.Warnf("get token failed. key [%v]", key)
						}
					}
				}
				for _, handler := range c.itemHandlers[""] {
					handler(c)
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
	if c.parent != nil {
		c.parent.AddOnItemChangeHandler(c.prefix, handler)
		return
	}
	c.AddOnItemChangeHandler("", handler)
}

func (c *Config) AddOnItemChangeHandler(key string, handler OnChangeHandler) {
	if c.parent != nil {
		c.parent.AddOnItemChangeHandler(prefixAppendKey(c.prefix, key), handler)
		return
	}
	c.itemHandlers[key] = append(c.itemHandlers[key], handler)
}

type TransformOptions struct {
	CipherKeys []string
	NoCipher   bool
}

type TransformOption func(options *TransformOptions)

func WithTransformCipherKeys(keys ...string) TransformOption {
	return func(options *TransformOptions) {
		options.CipherKeys = append(options.CipherKeys, keys...)
	}
}

func WithTransformNoCipher() TransformOption {
	return func(options *TransformOptions) {
		options.NoCipher = true
	}
}

func (c *Config) TransformWithOptions(options *Options, transformOptions *TransformOptions) (*Config, error) {
	if c.parent != nil {
		return nil, errors.New("transform is not supported by children")
	}
	provider, err := NewProviderWithOptions(&options.Provider)
	if err != nil {
		return nil, err
	}
	cipher, err := NewCipherWithOptions(&options.Cipher)
	if err != nil {
		return nil, err
	}
	decoder, err := NewDecoderWithOptions(&options.Decoder)
	if err != nil {
		return nil, err
	}

	storage := c.deepCopyStorage()
	if len(transformOptions.CipherKeys) != 0 {
		storage.SetCipherKeys(transformOptions.CipherKeys)
	}
	if transformOptions.NoCipher {
		storage.SetCipherKeys([]string{})
	}

	return &Config{
		provider:     provider,
		storage:      storage,
		decoder:      decoder,
		log:          StdoutLogger{},
		cipher:       cipher,
		itemHandlers: map[string][]OnChangeHandler{},
	}, nil
}

func (c *Config) Transform(options *Options, opts ...TransformOption) (*Config, error) {
	var transformOptions TransformOptions
	for _, opt := range opts {
		opt(&transformOptions)
	}
	return c.TransformWithOptions(options, &transformOptions)
}

func (c *Config) Bytes() ([]byte, error) {
	s := c.deepCopyStorage()
	if c.parent != nil {
		if err := s.Encrypt(c.parent.cipher); err != nil {
			return nil, err
		}
		return c.parent.decoder.Encode(s.Sub(c.prefix))
	}
	if err := s.Encrypt(c.cipher); err != nil {
		return nil, err
	}
	return c.decoder.Encode(s.Sub(c.prefix))
}

func (c *Config) Save() error {
	if c.parent != nil {
		return errors.New("children are not allow to save")
	}
	buf, err := c.Bytes()
	if err != nil {
		return err
	}
	return c.provider.Dump(buf)
}

func (c *Config) Diff(o *Config) {
	text1 := strx.JsonMarshalIndent(c.storage.Interface())
	text2 := strx.JsonMarshalIndent(o.storage.Interface())
	fmt.Println(strx.Diff(text1, text2))
}

func (c *Config) deepCopyStorage() *Storage {
	if c.parent != nil {
		return c.parent.deepCopyStorage()
	}

	buf, err := c.decoder.Encode(c.storage)
	if err != nil {
		panic(err)
	}
	s, err := c.decoder.Decode(buf)
	if err != nil {
		panic(err)
	}
	var keys []string
	for key := range c.storage.cipherKeySet {
		keys = append(keys, key)
	}
	s.SetCipherKeys(keys)
	return s
}

func (c *Config) ToString() string {
	buf, err := c.Bytes()
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (c *Config) ToJsonString() string {
	return strx.JsonMarshal(c.storage.Interface())
}

func (c *Config) SetLogger(log Logger) {
	c.log = log
}

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

type StdoutLogger struct{}

func (l StdoutLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (l StdoutLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
