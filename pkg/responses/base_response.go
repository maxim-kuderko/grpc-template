package responses

type BaseResponse struct {
	StatusCode int `json:"-"`
}

func (br BaseResponse) ResponseStatusCode() int {
	return br.StatusCode
}

type BaseResponser interface {
	ResponseStatusCode() int
}
