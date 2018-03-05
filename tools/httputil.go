package tools

import (
	"encoding/json"
	"github.com/hzxiao/goutil"
	"io"
	"io/ioutil"
	"net/http"
)

var httpClient http.Client

func HttpGet(url, token string) (goutil.Map, error) {
	return httpRequest("GET", url, "", token, nil)
}

func HttpPost(url, token, cType string, reader io.Reader) (goutil.Map, error) {
	return httpRequest("POST", url, cType, token, reader)
}

func HttpPut(url, token, cType string, reader io.Reader) (goutil.Map, error) {
	return httpRequest("PUT", url, cType, token, reader)
}

func HttpDelete(url, token, cType string, reader io.Reader) (goutil.Map, error) {
	return httpRequest("DELETE", url, cType, token, reader)
}

func httpRequest(method, url, cType, token string, reader io.Reader) (goutil.Map, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", BearerToken(token))
	}
	if cType != "" {
		req.Header.Set("Content-Type", cType)
	}

	return doRequest(req)
}

func BearerToken(token string) string {
	return "Bearer " + token
}

func doRequest(req *http.Request) (goutil.Map, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data goutil.Map
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return data, err
}
