// Package API
// @Description
// @Author root_wang
// @Date 2022/12/12 21:36
package API

import (
	"cqhttp-client/src/log"
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
	ApiKey string
}

func InitApi() *API {
	a := &API{
		apis:   make(map[SaasName]Parser),
		ApiKey: "",
	}
	a.apis[TextCompletion] = &TextApi{
		a,
	}

	return a
}

func (A *API) APIByName(n SaasName) module.Moduler {
	return A.apis[n]
}

func (A *API) MakeBody(api SaasName, msg string) interface{} {
	switch api {
	case TextCompletion:
		return &Request{
			Model:       Curie,
			Prompt:      msg,
			MaxTokens:   MaxTokens,
			Temperature: Temperature,
			N:           N,
		}
	default:
		panic("can't make a api req body")
	}
}

func (A *API) MakeResp(api SaasName) interface{} {
	switch api {
	case TextCompletion:
		return &Response{}
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
	req, err := http.NewRequest("POST", TextBaseURL, strings.NewReader(string(body)))
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
	response := resp.(*Response)
	return response.Choices[0].Text, nil
}

func (t TextApi) HandlerMessage(s string) (string, error) {
	handlerMessage, err := t.API.HandlerMessage(s, TextCompletion)
	if err != nil {
		return "", err
	}
	return handlerMessage, nil
}
