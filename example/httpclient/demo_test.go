package httpclient

import "testing"

func TestHttpDemo(t *testing.T) {
	res, err := httpDemo()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(res)
}
