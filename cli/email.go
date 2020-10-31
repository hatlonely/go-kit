package cli

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

func NewEmail(from string, password string, opts ...EmailOption) *EmailCli {
	options := defaultEmailOption
	for _, opt := range opts {
		opt(&options)
	}
	options.From = from
	options.Password = password

	return NewEmailWithOptions(&options)
}

func NewEmailWithConfig(cfg *config.Config, opts ...refx.Option) (*EmailCli, error) {
	options := defaultEmailOption
	if err := cfg.Unmarshal(&options, opts...); err != nil {
		return nil, err
	}

	return NewEmailWithOptions(&options), nil
}

func NewEmailWithOptions(options *EmailOptions) *EmailCli {
	return &EmailCli{
		Password: options.Password,
		From:     options.From,
		Server:   options.Server,
		Port:     options.Port,
	}
}

type EmailOptions struct {
	Password string
	From     string
	Server   string
	Port     int
}

var defaultEmailOption = EmailOptions{
	Server: "smtp.qq.com",
	Port:   25,
}

type EmailOption func(options *EmailOptions)

func WithEmailQQServer() EmailOption {
	return func(options *EmailOptions) {
		options.Server = "smtp.qq.com"
		options.Port = 25
	}
}

type EmailCli struct {
	Password string
	From     string
	Server   string
	Port     int
}

func (m *EmailCli) Send(to, subject, body string) error {
	content := fmt.Sprintf(`From: %v
To: %v
Subject: %v
Content-Type: text/html; charset=UTF-8;
%v
`, m.From, to, subject, body)

	return smtp.SendMail(
		fmt.Sprintf("%v:%v", m.Server, m.Port),
		smtp.PlainAuth("", m.From, m.Password, m.Server),
		m.From,
		strings.Split(to, ";"),
		[]byte(content),
	)
}
