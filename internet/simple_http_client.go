package internet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type simpleHTTPClient struct {
	url     string
	method  string
	body    io.Reader
	errs    []error
	headers map[string]string
	timeout time.Duration
}

// NewSimpleHTTPClient is the best practise way to use the http client for calling any API by default it already have header
//
//	"Content-Type": "application/json"
//
// and the timeout is 3 second
// example how to use
//
//	req := YourRequestStruct{Message: "Helloworld"}
//	var res YourResponseStruct
//
//	url := "https://yourapi.com"
//
//	err := internet.NewSimpleHTTPClient(http.MethodPost, url, req).
//		Header("Authorization", "Bearer YourAuthorization").
//		Call(&res)
//
//	if err != nil {
//		panic(err)
//	}
func NewSimpleHTTPClient(method, url string, requestData ...any) *simpleHTTPClient {

	errs := make([]error, 0)

	var err error
	var body io.Reader

	if len(requestData) > 0 {
		body, err = constructBody(requestData[0])
		if err != nil {
			errs = append(errs, err)
		}
	}

	return &simpleHTTPClient{
		url:    url,
		method: method,
		body:   body,
		errs:   errs,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		timeout: time.Second * 3,
	}
}

func (r *simpleHTTPClient) Method(method string) *simpleHTTPClient {
	if r.method != "" {
		return r
	}
	r.method = method
	return r
}

func (r *simpleHTTPClient) URL(url string) *simpleHTTPClient {
	if r.url != "" {
		return r
	}
	if url == "" {
		r.errs = append(r.errs, fmt.Errorf("url must not empty"))
	}
	r.url = url
	return r
}

func constructBody(requestData any) (io.Reader, error) {
	jsonInBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonInBytes), nil
}

func (r *simpleHTTPClient) Body(requestData any) *simpleHTTPClient {
	if r.body != nil {
		return r
	}

	body, err := constructBody(requestData)
	if err != nil {
		r.errs = append(r.errs, err)
		return r
	}

	r.body = body

	return r
}

func (r *simpleHTTPClient) Header(key string, value string) *simpleHTTPClient {
	r.headers[key] = value
	return r
}

func (r *simpleHTTPClient) Timeout(duration time.Duration) *simpleHTTPClient {
	r.timeout = duration
	return r
}

func (r *simpleHTTPClient) Call(responseData any) error {

	if len(r.errs) > 0 {
		errMessage := ""
		for i, s := range r.errs {
			if i == 0 {
				errMessage += s.Error()
				continue
			}
			errMessage += fmt.Sprintf(", %s", s.Error())
		}
		return fmt.Errorf(errMessage)
	}

	request, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return err
	}

	for k, v := range r.headers {
		request.Header.Set(k, v)
	}

	var client = &http.Client{Timeout: r.timeout}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			return
		}
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, responseData)
	if err != nil {
		return err
	}

	return nil
}

func (r *simpleHTTPClient) CallAndPrint(responseData any) error {
	err := r.Call(responseData)
	if err != nil {
		return err
	}

	arrBytes, err := json.MarshalIndent(responseData, "", " ")
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", string(arrBytes))

	return nil
}
