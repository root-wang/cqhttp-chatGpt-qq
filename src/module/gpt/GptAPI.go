// Package gpt
// @Description
// @Author root_wang
// @Date 2022/12/12 21:36
package gpt

import (
	"cqhttp-client/src/config"
	"cqhttp-client/src/constant"
	"cqhttp-client/src/log"
	"cqhttp-client/src/message"
	"cqhttp-client/src/module"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type SaasName string

func (s SaasName) String() string {
	return string(s)
}

type Parser interface {
	module.Moduler
	// Parse todo:解析错误未实现
	Parse(interface{}) (string, error)
}

type API struct {
	apis   map[SaasName]Parser
	urls   map[SaasName]string
	ApiKey string
}

func InitApi() *API {
	tmp := cfg.GetYAMLFiled("api-key")
	a := &API{
		apis:   make(map[SaasName]Parser),
		urls:   make(map[SaasName]string),
		ApiKey: tmp["api-key"].(string),
	}
	a.apis[TextCompletion] = &TextApi{
		a,
	}
	a.urls[TextCompletion] = TextBaseURL

	a.apis[ImageGeneration] = &ImageApi{
		a,
	}
	a.urls[ImageGeneration] = ImageBaseURL

	a.apis[ChatCompletion] = &ChatApi{
		a,
	}
	a.urls[ChatCompletion] = ChatBseURL
	return a
}

// APIByName 通过模块名称请求模块
func (A *API) APIByName(n SaasName) module.Moduler {
	return A.apis[n]
}

// MakeRequestStruct 为不同的请求构建不同的结构体
func (A *API) MakeRequestStruct(api SaasName, msg interface{}) interface{} {
	switch api {
	case TextCompletion:
		return &TextReq{
			Model:       Davinci,
			Prompt:      msg.(string),
			MaxTokens:   TextMaxTokens,
			Temperature: TextTemperature,
			N:           TextN,
		}
	case ImageGeneration:
		return &ImageReq{
			Prompt: msg.(string),
			N:      ImageN,
			Size:   ImageSize,
		}
	case ChatCompletion:
		return &ChatReq{
			Model:            Gpt35Turbo,
			Messages:         msg.([]*ChatMessage),
			Temperature:      ChatTemperature,
			TopP:             ChatTopP,
			N:                ChatN,
			MaxTokens:        ChatMaxTokens,
			PresencePenalty:  ChatPresencePenalty,
			FrequencyPenalty: ChatFrequencyPenalty,
		}
	default:
		panic("can't make a api req body")
	}
}

// MakeRespStruct 为不同的请求响应构建不同的结构体
func (A *API) MakeRespStruct(api SaasName) interface{} {
	switch api {
	case TextCompletion:
		return &TextResp{}
	case ImageGeneration:
		return &ImageResp{}
	case ChatCompletion:
		return &ChatResp{}
	default:
		panic("can't make a api resp")
	}
}

// ParseMessage 调用不同的模块的响应解析方法来解析请求返回的响应结构体中的消息
func (A *API) ParseMessage(api SaasName, resp interface{}) (string, error) {
	parseMsg, err := A.apis[api].Parse(resp)
	if err != nil {
		return cst.ErrorOpenAIParseResponse, fmt.Errorf("failed to parse message from api%s %v", api, err)
	}
	return parseMsg, nil
}

// HandlerMessage 集中处理不同模块的请求 同时在获得请求响应时调用对应的响应解析方法 ParseMessage
func (A *API) HandlerMessage(s interface{}, api SaasName) (string, error) {
	body, _ := json.Marshal(A.MakeRequestStruct(api, s))
	req, err := http.NewRequest("POST", A.urls[api], strings.NewReader(string(body)))
	if err != nil {
		return cst.ErrorOpenAIResponse, log.ErrorInsidef("failed to create api request: %v", err)
	}

	req.Header.Set(cst.ContentType, cst.ContentTypeValue)
	req.Header.Set(cst.Authorization, fmt.Sprintf("Bearer %s", A.ApiKey))
	req.Header.Set(cst.UserAgent, cst.UserAgentValue)

	h := &http.Client{}
	resp, err := h.Do(req)
	if err != nil {
		return cst.ErrorOpenAIResponse, log.ErrorInsidef("failed to connect to api: %v", err)
	}

	if resp.StatusCode != 200 {
		return cst.ErrorOpenAIResponse, log.ErrorInsidef("failed to connect to api: %v", resp.Status)
	}

	respStruct := A.MakeRespStruct(api)
	err = json.NewDecoder(resp.Body).Decode(respStruct)
	if err != nil {
		return cst.ErrorOpenAIResponse, log.ErrorInsidef("failed to parse message from api:%s :%v", api, err)
	}

	return A.ParseMessage(api, respStruct)
}

// TextApi 文本补全模块
type TextApi struct {
	*API
}

func (t TextApi) Matcher(s string) (bool, string) {
	compile := regexp.MustCompile(`^[pc]&`)
	isMatch := compile.MatchString(s)
	return !isMatch, s
}

// Parse 解析文本补全的响应结构体
func (t TextApi) Parse(resp interface{}) (string, error) {
	response := resp.(*TextResp)
	return fmt.Sprintf(
		"%s\n", response.Choices[0].Text,
	), nil
}

func (t TextApi) HandlerMessage(s interface{}) (string, error) {
	return t.API.HandlerMessage(s, TextCompletion)
}

// ImageApi 图片生成模块
type ImageApi struct {
	*API
}

func (i ImageApi) Matcher(s string) (ok bool, prompt string) {
	ok = false
	// use p& prefix
	imageReg := regexp.MustCompile(`p&amp;(.+)`)
	matches := imageReg.FindStringSubmatch(s)
	if matches != nil {
		ok = true
		prompt = matches[1]
		return
	}
	return
}

// Parse 解析图片生成的响应结构体
func (i ImageApi) Parse(resp interface{}) (string, error) {
	response := resp.(*ImageResp)
	url := response.Data[0].Url
	// url := "https://tenfei02.cfp.cn/creative/vcg/veer/1600water/veer-153029426.jpg"
	cqCode := msg.NewCQCode("", msg.IMAGE)
	cqCode.SetKeyValue([]msg.CQKEY{msg.FILE}, url)

	return fmt.Sprintf("%s\n", cqCode), nil
}

func (i ImageApi) HandlerMessage(s interface{}) (string, error) {
	return i.API.HandlerMessage(s, ImageGeneration)
}

// ChatApi 聊天对话模块
type ChatApi struct {
	*API
}

func (c ChatApi) Matcher(s string) (ok bool, prompt string) {
	ok = false
	// use c& prefix
	chatReg := regexp.MustCompile(`c&amp;(.+)`)
	matches := chatReg.FindStringSubmatch(s)
	if matches != nil {
		ok = true
		prompt = matches[1]
		return
	}
	return
}

// Parse 解析聊天对话的响应结构体
func (c ChatApi) Parse(resp interface{}) (string, error) {
	response := resp.(*ChatResp)
	return fmt.Sprintf(
		"%s", response.Choices[0].Message.Content,
	), nil
}

func (c ChatApi) HandlerMessage(s interface{}) (string, error) {
	return c.API.HandlerMessage(s, ChatCompletion)
}
