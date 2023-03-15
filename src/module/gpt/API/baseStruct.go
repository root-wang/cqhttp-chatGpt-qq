// Package API
// @Description
// @Author root_wang
// @Date 2022/12/12 21:36
package API

import (
	"cqhttp-client/src/log"
	"cqhttp-client/src/message"
	"cqhttp-client/src/module"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	UserAgent       = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
	Authorization   = "Authorization"
)

type SaasName string

func (s SaasName) String() string {
	return string(s)
}

type Parser interface {
	module.Moduler
	Parse(interface{}) (string, error)
}

type API struct {
	apis   map[SaasName]Parser
	urls   map[SaasName]string
	ApiKey string
}

func InitApi() *API {
	a := &API{
		apis:   make(map[SaasName]Parser),
		urls:   make(map[SaasName]string),
		ApiKey: "",
	}
	a.apis[TextCompletion] = &TextApi{
		a,
	}
	a.urls[TextCompletion] = TextBaseURL

	a.apis[ImageGeneration] = &ImageApi{
		a,
	}
	a.urls[ImageGeneration] = ImageBaseURL

	return a
}

func (A *API) APIByName(n SaasName) module.Moduler {
	return A.apis[n]
}

func (A *API) MakeBody(api SaasName, msg string) interface{} {
	switch api {
	case TextCompletion:
		return &TextReq{
			Model:       Davinci,
			Prompt:      msg,
			MaxTokens:   TextMaxTokens,
			Temperature: TextTemperature,
			N:           TextN,
		}
	case ImageGeneration:
		return &ImageReq{
			Prompt: msg,
			N:      ImageN,
			Size:   ImageSize,
		}
	default:
		panic("can't make a api req body")
	}
}

func (A *API) MakeResp(api SaasName) interface{} {
	switch api {
	case TextCompletion:
		return &TextResp{}
	case ImageGeneration:
		return &ImageResp{}
	default:
		panic("can't make a api resp")
	}
}

func (A *API) ParseMessage(api SaasName, resp interface{}) (string, error) {
	parse, err := A.apis[api].Parse(resp)
	if err != nil {
		return "", fmt.Errorf("failed to parse message from api%s %v", api, err)
	}
	return parse, nil
}

func (A *API) HandlerMessage(s string, api SaasName) (string, error) {
	body, _ := json.Marshal(A.MakeBody(api, s))
	req, err := http.NewRequest("POST", A.urls[api], strings.NewReader(string(body)))
	if err != nil {
		return "", log.ErrorInsidef("failed to create api request: %v", err)
	}
	req.Header.Set(ContentType, ApplicationJson)
	req.Header.Set(Authorization, fmt.Sprintf("Bearer %s", A.ApiKey))
	req.Header.Set("User-Agent", UserAgent)
	h := &http.Client{}
	resp, err := h.Do(req)
	if err != nil {
		return "", log.ErrorInsidef("failed to connect to api: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", log.ErrorInsidef("failed to connect to api: %v", resp.Status)
	}

	respMsg := A.MakeResp(api)
	err = json.NewDecoder(resp.Body).Decode(respMsg)
	if err != nil {
		return "", log.ErrorInsidef("failed to parse message from api:%s :%v", api, err)
	}
	return A.ParseMessage(api, respMsg)
}

type TextApi struct {
	*API
}

func (t TextApi) Parse(resp interface{}) (string, error) {
	response := resp.(*TextResp)
	return fmt.Sprintf(
		"%s\n\n本次共消费%f美元", response.Choices[0].Text, float32(response.Usage.TotalTokens)*0.002/1000,
	), nil
}

func (t TextApi) HandlerMessage(s string) (string, error) {
	handlerMessage, err := t.API.HandlerMessage(s, TextCompletion)
	if err != nil {
		return "", err
	}
	return handlerMessage, nil
}

type ImageApi struct {
	*API
}

func (i ImageApi) HandlerMessage(s string) (string, error) {
	handlerMessage, err := i.API.HandlerMessage(s, ImageGeneration)
	if err != nil {
		return "", err
	}
	return handlerMessage, nil
}

func (i ImageApi) Parse(resp interface{}) (string, error) {
	response := resp.(*ImageResp)
	url := response.Data[0].Url
	// url := "https://tenfei02.cfp.cn/creative/vcg/veer/1600water/veer-153029426.jpg"
	cqCode := message.NewCQCode("", message.IMAGE)
	cqCode.SetKeyValue([]message.CQKEY{message.FILE}, url)

	return fmt.Sprintf("%s\n\n本次共消费0.016美元", cqCode), nil
}
