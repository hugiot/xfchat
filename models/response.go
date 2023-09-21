package models

type Response struct {
	Header  ResponseHeader  `json:"header"`
	Payload ResponsePayload `json:"payload"`
}

type ResponseHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	SID     string `json:"sid"`
	Status  int    `json:"status"`
}

type ResponsePayload struct {
	Choices ResponsePayloadChoices `json:"choices"`
	Usage   ResponsePayloadUsage   `json:"usage"`
}

type ResponsePayloadChoices struct {
	Status int                          `json:"status"`
	Seq    int                          `json:"seq"`
	Text   []ResponsePayloadChoicesText `json:"text"`
}

type ResponsePayloadChoicesText struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Index   int    `json:"index"`
}

type ResponsePayloadUsage struct {
	Text ResponsePayloadUsageText `json:"text"`
}

type ResponsePayloadUsageText struct {
	QuestionTokens   int `json:"question_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
