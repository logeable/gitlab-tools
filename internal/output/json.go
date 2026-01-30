package output

import (
	"encoding/json"
	"io"
)

// WriteJSON 将 v 序列化为 JSON 写入 w，字段使用 snake_case（由结构体 json tag 控制）
func WriteJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

// JSONError 用于 --json 时的结构化错误输出（stderr）
type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"`
}

// WriteJSONError 将错误信息以 JSON 写入 w（通常为 stderr），code 为退出码 1 或 2
func WriteJSONError(w io.Writer, message string, code int) error {
	return WriteJSON(w, JSONError{Error: message, Code: code})
}
