package hc

import (
	"encoding/json"
	"fmt"
	"github.com/feature-vector/harbor/base/utils"
	"io/ioutil"
	"net/http"
)

type HttpResponseHelper struct {
	Response *http.Response

	bodyBytes   []byte
	bodyReadErr error
	bodyRead    bool
}

func WrapHttpResponse(resp *http.Response) *HttpResponseHelper {
	return &HttpResponseHelper{
		Response: resp,
	}
}

func (r *HttpResponseHelper) Unmarshal(v interface{}) error {
	bytes, err := r.BodyBytes()
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return fmt.Errorf("http_%d %s", r.Response.StatusCode, string(bytes))
	}
	return nil
}

func (r *HttpResponseHelper) BodyBytes() ([]byte, error) {
	if r.bodyRead {
		return r.bodyBytes, r.bodyReadErr
	}
	defer utils.CloseSilent(r.Response.Body)
	r.bodyRead = true
	r.bodyBytes, r.bodyReadErr = ioutil.ReadAll(r.Response.Body)
	return r.bodyBytes, r.bodyReadErr
}

func (r *HttpResponseHelper) String() string {
	return fmt.Sprintf("%d %b", r.Response.StatusCode, r.bodyBytes)
}
