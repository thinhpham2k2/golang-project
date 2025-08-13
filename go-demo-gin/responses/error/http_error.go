package error

import "encoding/json"

type Error struct {
	Error map[string]string `json:"error"`
}

type HTTPError struct {
	StatusCode int
	Message    Error
}

// Đảm bảo implement interface `error`
func (e *HTTPError) Error() string {
	bytes, err := json.Marshal(e.Message.Error)
	if err != nil {
		return "Internal Error"
	}
	return string(bytes)
}
