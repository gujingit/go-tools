package net

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var log *logrus.Entry

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
		PrettyPrint:      false,
	})
	logger.SetOutput(os.Stdout)

	log = logger.WithField("controller", "algorithm")
}

func DoPost(url string, data interface{}) (string, error) {
	if data == nil {
		log.Errorf("post data is nil, skip")
		return "", nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		return "", fmt.Errorf("post %s, data: %s, error: %s", url, string(jsonBytes), err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	defer safeClose(resp)
	if err != nil {
		return "", fmt.Errorf("call %s api error: %s", url, err.Error())
	}
	if resp == nil {
		return "", fmt.Errorf("call %s api resp is nil", url)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("url %s response error: %+v", url, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read predict response error: %v", err.Error())
	}

	return string(body), nil
}

func safeClose(resp *http.Response) {
	if resp != nil && !resp.Close {
		err := resp.Body.Close()
		if err != nil {
			log.Errorf("error close response body:%v\n", err)
		}
	}
}
