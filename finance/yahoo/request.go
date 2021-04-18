package yahoo

import "net/url"

type request struct {
	method    string
	path      string
	urlValues url.Values
	reqBody   interface{}
	response  interface{}
}
