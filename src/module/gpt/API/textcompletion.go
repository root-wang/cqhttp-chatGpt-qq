// Package API
// @Description  use for chatting same as GptChat
// @Author root_wang
// @Date 2022/12/12 21:00
package API

type (
	TextReq struct {
		// need to get model list
		Model string `json:"model"`
		// send msg
		Prompt string `json:"prompt"`
		// max resp token
		MaxTokens int `json:"max_tokens"`
		// default 1
		Temperature int `json:"temperature"`
		// default 0
		TopP int `json:"top_p"`
		// default 1 how many resp of msg
		N int `json:"n"`
		// default false stream to show tokens
		Stream bool `json:"stream"`
		// default nil
		Logprobs int `json:"logprobs"`
		// default nil
		Stop string `json:"stop"`
	}

	TextResp struct {
		Id      string        `json:"id"`
		Object  string        `json:"object"`
		Created int64         `json:"created"`
		Model   string        `json:"model"`
		Choices []textChoices `json:"choices"`
		Usage   usage         `json:"usage"`
	}

	textChoices struct {
		Text         string   `json:"text"`
		Index        int      `json:"index"`
		Logprobs     logprobs `json:"logprobs"`
		FinishReason string   `json:"finish_reason"`
	}
	logprobs struct {
		Tokens        []string    `json:"tokens"`
		TokenLogprobs []float32   `json:"token_logprobs"`
		TopLogprobs   interface{} `json:"top_logprobs"`
		TextOffset    []int       `json:"text_offset"`
	}

	usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
)

const (
	TextCompletion SaasName = "textcompletion"
	TextBaseURL             = "https://api.openai.com/v1/completions"
)

const (
	TextMaxTokens   = 200
	TextTemperature = 1
	TextTopP        = 0
	TextN           = 1
	TextStream      = false
	TextLogprobs    = 0
	TextStop        = ""
)

// Model
const (
	// Davinci Good at: Complex intent, cause and effect, summarization for audience
	Davinci = "text-davinci-003"
	// Curie Good at: Language translation, complex classification, text sentiment, summarization `default`
	Curie = "text-curie-001"
	// Babbage Good at: Moderate classification, semantic search classification
	Babbage = "text-babbage-001"
	// Ada Good at: Parsing text, simple classification, address correction, keywords
	Ada = "text-ada-001"
)
