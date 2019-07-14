package request

import (
	"errors"
	"github.com/swift9/gorequest"
	"reflect"
	"time"
)

type Request struct {
	proxy   *gorequest.SuperAgent
	timeout time.Duration
}

func New() *Request {
	request := &Request{
		proxy:   gorequest.New(),
		timeout: 3 * time.Minute,
	}
	request.proxy.Timeout(request.timeout)
	return request
}

func (request *Request) AddHeader(name string, value string) *Request {
	request.proxy.Set(name, value)
	return request
}

func (request *Request) RemoveHeader(name string) *Request {
	delete(request.proxy.Header, name)
	return request
}

func (request *Request) Timeout(t time.Duration) *Request {
	request.proxy.Timeout(t)
	return request
}

func (request *Request) RealRequest() *gorequest.SuperAgent {
	return request.proxy
}

func (request *Request) cloneRealRequest() *gorequest.SuperAgent {
	return request.proxy.Clone()
}

func (request *Request) Get(url string, result interface{}, retries ...int) []error {
	retryTimes := 0
	length := len(retries)
	if length > 0 {
		retryTimes = retries[0]
	}
	errs := []error{errors.New("unknown error")}
	res := ""
	req := request.cloneRealRequest()
	req.Get(url)
	for i := 0; i < retryTimes+1; i++ {
		switch result.(type) {
		case *string:
			_, res, errs = req.End()
			if errs == nil {
				reflect.ValueOf(result).Elem().SetString(res)
				return errs
			}
		default:
			_, _, errs = req.EndStruct(&result)
			if errs == nil {
				return errs
			}
		}
		if i != retryTimes {
			time.Sleep(3 * time.Second)
		}
	}
	return errs
}

func (request *Request) Post(contentType string, url string, data interface{}, result interface{}, retries ...int) []error {
	retryTimes := 0
	length := len(retries)
	if length > 0 {
		retryTimes = retries[0]
	}
	errs := []error{errors.New("unknown error")}
	res := ""
	req := request.cloneRealRequest()
	req.Type(contentType)
	req.Post(url).Send(data)

	for i := 0; i < retryTimes+1; i++ {
		switch result.(type) {
		case *string:
			_, res, errs = req.End()
			if errs == nil {
				reflect.ValueOf(result).Elem().SetString(res)
				return errs
			}
		default:
			_, _, errs = req.EndStruct(&result)
			if errs == nil {
				return errs
			}
		}
		if i != retryTimes {
			time.Sleep(3 * time.Second)
		}
	}
	return errs
}

func (request *Request) PostForm(url string, data interface{}, result *interface{}, retries ...int) []error {
	return request.Post("form", url, data, result, retries...)
}

func (request *Request) PostJson(url string, data interface{}, result *interface{}, retries ...int) []error {
	return request.Post("json", url, data, result, retries...)
}
