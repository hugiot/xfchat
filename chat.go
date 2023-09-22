package xfchat

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hugiot/xfchat/models"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	v1 string = "v1.1"
	v2 string = "v2.1"
)

const (
	domainV1 string = "general"
	domainV2 string = "generalv2"
)

var (
	defaultTemperature           = 0.5
	defaultMaxTokens             = 2048
	defaultTopK                  = 4
	defaultOutput      io.Writer = os.Stderr
)

var urlTemplate = "wss://spark-api.xf-yun.com/%s/chat"

type Chat struct {
	appID     string
	apiKey    string
	apiSecret string
	// version v1 or v2
	version string
	// domain general or generalv2
	domain string
	// temperature sampling threshold [0, 1], default 0.5
	temperature float64
	// maxTokens answer tokens length
	maxTokens int
	// topK [1, 6], default 4
	topK int
	// answer output, default os.Stdout
	output    io.Writer
	answering bool
}

type Option interface {
	apply(*Chat)
}

type optionFunc func(c *Chat)

func (f optionFunc) apply(c *Chat) {
	f(c)
}

func New(appID string, apiKey string, apiSecret string, options ...Option) (*Chat, error) {
	chat := &Chat{
		appID:       appID,
		apiKey:      apiKey,
		apiSecret:   apiSecret,
		version:     v1,
		domain:      domainV1,
		temperature: defaultTemperature,
		maxTokens:   defaultMaxTokens,
		topK:        defaultTopK,
		output:      defaultOutput,
		answering:   false,
	}

	if len(options) > 0 {
		for i, _ := range options {
			options[i].apply(chat)
		}
	}

	return chat, nil
}

// Ask reconnect every time
func (c *Chat) Ask(question string) (err error) {
	if c.answering {
		_, _ = fmt.Fprintln(c.output, "please wait")
		return
	}

	question = strings.Trim(question, " ")
	if question == "" {
		_, _ = fmt.Fprintln(c.output, "please enter your question")
		return
	}

	c.answering = true
	conn, err := c.connect()
	if err != nil {
		return
	}
	defer conn.Close()
	if err = conn.WriteJSON(c.createRequest(question)); err != nil {
		return err
	}

	var answer []byte
	for {
		_, answer, err = conn.ReadMessage()
		if err != nil {
			return
		}
		var res models.Response
		if err = json.Unmarshal(answer, &res); err != nil {
			return
		}
		_, _ = fmt.Fprintf(c.output, "%s", res.Payload.Choices.Text[0].Content)
		// over
		if res.Header.Status == 2 {
			c.answering = false
			break
		}
	}
	return
}

func (c *Chat) createRequest(question string) models.Request {
	return models.Request{
		Header: models.RequestHeader{
			AppID: c.appID,
			UID:   "123456",
		},
		Parameter: models.RequestParameter{
			Chat: models.RequestParameterChat{
				Domain:      c.domain,
				Temperature: c.temperature,
				MaxTokens:   c.maxTokens,
				TopK:        c.topK,
				ChatID:      "",
			},
		},
		Payload: models.RequestPayload{
			Message: models.RequestPayloadMessage{
				Text: []models.RequestPayloadMessageText{
					{
						Role:    "user",
						Content: question,
					},
				},
			},
		},
	}
}

func (c *Chat) connect() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.getAuthUrl(), nil)
	return conn, err
}

// getAuthUrl url with authentication information
func (c *Chat) getAuthUrl() string {
	base, _ := url.Parse(c.getBaseUrl())
	// date Tue, 28 May 2019 09:10:42 MST
	date := time.Now().UTC().Format(time.RFC1123)
	// sign string
	signString := fmt.Sprintf("host: %s\ndate: %s\nGET %s HTTP/1.1", base.Host, date, base.Path)
	// result
	sha := base64.StdEncoding.EncodeToString(HmacSHA256(c.apiSecret, signString))
	// auth url
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", c.apiKey, "hmac-sha256", "host date request-line", sha)
	// base64
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))
	// url
	v := url.Values{}
	v.Add("host", base.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	return c.getBaseUrl() + "?" + v.Encode()
}

func (c *Chat) Close() error {
	return nil
}

// getBaseUrl get base url, for example: wss://spark-api.xf-yun.com/v1.1/chat
func (c *Chat) getBaseUrl() string {
	return fmt.Sprintf(urlTemplate, c.version)
}

// UseVersion1 use v1.1
func UseVersion1() Option {
	return optionFunc(func(c *Chat) {
		c.version = v1
		c.domain = domainV1
	})
}

// UseVersion2 use v2.1
func UseVersion2() Option {
	return optionFunc(func(c *Chat) {
		c.version = v2
		c.domain = domainV2
	})
}

// SetTemperature sampling threshold
// [0, 1], default 0.5
func SetTemperature(t float64) Option {
	return optionFunc(func(c *Chat) {
		if t >= 0 && t <= 1 {
			c.temperature = t
		}
	})
}

// SetMaxTokens set the maximum length of the model answer tokens
// v1 [1, 4096]
// v2 [1, 8192]
func SetMaxTokens(max int) Option {
	return optionFunc(func(c *Chat) {
		switch c.version {
		case v1:
			if max >= 1 && max <= 4096 {
				c.maxTokens = max
			}
		case v2:
			if max >= 1 && max <= 8192 {
				c.maxTokens = max
			}
		}
	})
}

// SetTopK randomly select one of the k candidates (equal probability)
// [1, 6], default 4
func SetTopK(v int) Option {
	return optionFunc(func(c *Chat) {
		if v >= 1 && v <= 6 {
			c.topK = v
		}
	})
}

// SetOutput set answer output
func SetOutput(w io.Writer) Option {
	return optionFunc(func(c *Chat) {
		c.output = w
	})
}
