package http

type httpResponse struct {
	body []byte
}

func (r *httpResponse) Body() []byte {
	return r.body
}
