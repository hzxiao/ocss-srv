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

func HttpPost(url, token, ctype string, reader io.Reader) (goutil.Map, error) {
	return httpRequest("POST", url, ctype, token, reader)
}

func HttpPut(url, token, ctype string, reader io.Reader) (goutil.Map, error) {
	return httpRequest("PUT", url, ctype, token, reader)
}

func HttpDelete(url, token, ctype string, reader io.Reader) (goutil.Map, error) {
	return httpRequest("DELETE", url, ctype, token, reader)
}

func httpRequest(method, url, ctype, token string, reader io.Reader) (goutil.Map, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", BearerToken(token))
	}
	if ctype != "" {
		req.Header.Set("Content-Type", ctype)
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
