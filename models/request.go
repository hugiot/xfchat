package models

type Request struct {
	Header    RequestHeader    `json:"header"`
	Parameter RequestParameter `json:"parameter"`
	Payload   RequestPayload   `json:"payload"`
}

type RequestHeader struct {
	AppID string `json:"app_id"`
	UID   string `json:"uid"`
}

type RequestParameter struct {
	Chat RequestParameterChat `json:"chat"`
}

type RequestParameterChat struct {
	Domain      string  `json:"domain"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	TopK        int     `json:"top_k"`
	ChatID      string  `json:"chat_id"`
}

type RequestPayload struct {
	Message RequestPayloadMessage `json:"message"`
}

type RequestPayloadMessage struct {
	Text []RequestPayloadMessageText `json:"text"`
}

type RequestPayloadMessageText struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
