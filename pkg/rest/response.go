package rest

import "encoding/json"

type Response struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	Result      any    `json:"result"`
}

func NewResponse(ok bool, description string, result any) *Response {
	return &Response{
		Ok:          ok,
		Description: description,
		Result:      result,
	}
}

func (r *Response) Marshal() []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return bytes
}
