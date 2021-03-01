package responses

type BaseResponse struct {
	StatusCode int `json:"-"`
}

func (br BaseResponse) ResponseStatusCode() int {
	if br.StatusCode == 0 {
		return 200
	}
	return br.StatusCode
}

type BaseResponser interface {
	ResponseStatusCode() int
}
