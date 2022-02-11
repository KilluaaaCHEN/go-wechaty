package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var tryCount = 0

func Get(url string, header map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	for key, val := range header {
		req.Header.Set(key, val)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		tryCount++
		if tryCount > 3 {
			return nil, err
		}
		fmt.Printf("网络请求失败:%v\r", err)
		time.Sleep(time.Second * 3)
		return Get(url, header)
	}
	if resp.StatusCode != 200 {
		tryCount++
		if tryCount > 3 {
			fmt.Printf("%v ,尝试多次失败: %v\n", url, resp.Status)
			return nil, errors.New("重试多次失败" + resp.Status)
		}
		fmt.Printf("网络请求失败:%s\n3秒后开始重试第%d次...\r", resp.Status, tryCount)
		time.Sleep(time.Second * 3)
		return Get(url, header)
	}
	if tryCount > 0 {
		tryCount = 0
	}
	return resp, nil
}

const (
	contentTypeForm = "application/x-www-form-urlencoded"
	contentTypeJson = "application/json"
	methodGet       = "GET"
	methodPost      = "POST"
)

func PostJson(url string, params interface{}, headers map[string]interface{}) (string, int, error) {
	return Request(url, params, headers, methodPost, contentTypeJson)
}

func Request(url string, params interface{}, headers map[string]interface{}, method string, contentType string) (string, int, error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest(method, url, getReader(params, contentType))
	if err != nil {
		return "", 0, err
	}
	setHeader(req, headers, contentType)
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	return string(body), resp.StatusCode, nil
}

func getReader(params interface{}, contentType string) io.Reader {
	switch vv := params.(type) {
	case string:
		return strings.NewReader(vv)
	case map[string]interface{}:
		switch contentType {
		case contentTypeJson:
			bytesData, err := json.Marshal(vv)
			if err != nil {
				return nil
			}
			return bytes.NewReader(bytesData)
		case contentTypeForm:
			var paramStr string
			for k, v := range vv {
				paramStr += fmt.Sprintf("%s=%v&", k, v)
			}
			return strings.NewReader(strings.TrimRight(paramStr, "&"))
		}
	}
	return nil
}

func setHeader(req *http.Request, headers map[string]interface{}, contentType string) {
	if headers == nil {
		headers = map[string]interface{}{}
	}
	headers["Content-Type"] = contentType
	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
}
