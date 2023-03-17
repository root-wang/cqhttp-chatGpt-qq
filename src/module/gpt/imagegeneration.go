// Package gpt
// @Description
// @Author root_wang
// @Date 2022/12/13 9:30
package gpt

type (
	ImageReq struct {
		Prompt string `json:"prompt"`
		N      int    `json:"n"`
		Size   string `json:"size"`
	}

	ImageResp struct {
		Created int64  `json:"created"`
		Data    []data `json:"data"`
	}
	data struct {
		Url string `json:"url"`
	}
)

const (
	ImageGeneration SaasName = "imagegeneration"
	ImageBaseURL             = "https://api.openai.com/v1/images/generations"
)

const (
	ImageN = 1
	// ImageSize = "1024x1024"
	ImageSize = "256x256"
)
