package http

// Response represent http response
type Response struct {
	statusLine string
	headers map[string]string
	responseBody string
}

// ToData traslate response to an []byte
func (r *Response) ToData() ([]byte) {
	var data = r.statusLine
	data += "\r\n"
	for key, value := range r.headers {
		data += key
		data += ":"
		data += value
		data += "\r\n"
	}
	data += "\r\n"
	data += r.responseBody
	return []byte(data)
}

// PutHeader is put header into response
func (r *Response) PutHeader(key string, value string) {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[key] = value
}

// SetBody is set response body content
func (r *Response) SetBody(responseBody string) {
	r.responseBody = responseBody
}

// GetContentLength is get response body content length
func (r *Response) GetContentLength() (int){
	return len(r.responseBody)
}

// SetStatusLine is set response body statusLine
func (r *Response) SetStatusLine(statusLine string) {
	r.statusLine = statusLine
}