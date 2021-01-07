package comm

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)


// SendRequest 发送request
func SendRequest(url string, body io.Reader, addHeaders map[string]string, method string) (resp []byte, err error) {
	// 1、创建req
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	//req.Header.Add("Content-Type", "application/json")

	// 2、设置headers
	if len(addHeaders) > 0 {
		for k, v := range addHeaders {
			req.Header.Add(k, v)
		}
	}

	// 3、发送http请求
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errors.New("http status err")
		fmt.Errorf("sendRequest failed, url=%v, response status code=%d", url, response.StatusCode)
		return
	}

	// 4、结果读取
	resp, err = ioutil.ReadAll(response.Body)
	return
}