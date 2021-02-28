package responses

import "context"

type Get struct {
	context.Context `json:"-"`
	BaseResponse    `json:"-"`
	Value           string `json:"value,omitempty"`
}
