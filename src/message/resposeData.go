// Package msg
// @Description
// @Author root_wang
// @Date 2022/12/10 20:43
package msg

type GroupResp struct {
	GroupId    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}
