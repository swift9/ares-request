package request

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Test struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func TestRequest_Get(t *testing.T) {
	request := New()
	r := Test{
		Code: 20,
	}
	b, _ := json.Marshal(r)
	fmt.Println(string(b[:]))
	request.Get("https://api.apiopen.top/EmailSearch?number=1012002", &r)
	fmt.Println(r.Code)
}
