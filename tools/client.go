package tools

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"srun4-pay/configs"
)

var (
	Client *http.Client
)

func init() {
	if Client == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		Client = &http.Client{
			Transport: tr,
		}
	}
}

func Get(url string) ([]byte, error) {
	resp, err := Client.Get(url)
	if err != nil {
		configs.Log.WithField("HTTP GET ERROR", err).Error()
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		configs.Log.WithField("读取 GET Response 失败", err).Error()
		return nil, err
	}
	return res, nil
}

func Post(url string, params url.Values) ([]byte, error) {
	resp, err := Client.PostForm(url, params)
	if err != nil {
		configs.Log.WithField(fmt.Sprintf("HTTP POST[%s] ERROR", url), err).Error()
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		configs.Log.WithField("读取 POST Response 失败", err).Error()
		return nil, err
	}
	return res, nil
}
