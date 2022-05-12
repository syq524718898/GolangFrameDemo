package httputil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func HTTPGetRequest(rawUrl string, headers map[string]string, paramses map[string]string) (rsp map[string]interface{}, err error) {
	params := url.Values{}
	Url, err := url.Parse(rawUrl)
	if err != nil {
		return
	}
	for k,v := range paramses {
		params.Set(k, v)
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()

	client := &http.Client{}
	req,_ := http.NewRequest("GET",urlPath,nil)

	for k,v := range headers {
		req.Header.Add(k, v)
	}
	resp,_ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &rsp)
	return
}

func HTTPPostRequest(rawUrl string, headers map[string]string, body map[string]string) (rsp map[string]interface{}, err error)   {
	client := &http.Client{}
	bytesData, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", rawUrl, bytes.NewReader(bytesData))
	for k,v := range headers {
		req.Header.Add(k, v)
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyRsp, &rsp)
	return
}