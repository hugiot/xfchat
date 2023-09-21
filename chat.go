package xfchat

import "io"

type Options struct {
	AppID     string
	APISecret string
	APIKey    string
}

type Chat struct {
}

func New(o Options) *Chat {
	return &Chat{}
}

func (c *Chat) Ask(question string, answer io.Writer) (err error) {

	return nil
}
