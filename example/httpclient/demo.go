package httpclient

import (
	"github.com/gowins/dionysus/httpclient"
	"io"
	"net/http"
)

func httpDemo() (string, error) {
	client := httpclient.New()

	resp, err := client.Get("https://skyao.io/learning-serverless/docs/spec/cloudevents/core.html", http.Header{})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
