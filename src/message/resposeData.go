// Package message
// @Description
// @Author root_wang
// @Date 2022/12/10 20:43
package message

type groupResp struct {
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}
