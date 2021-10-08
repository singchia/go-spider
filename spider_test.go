package spider

import (
	"encoding/json"
	"net/http"
	"testing"
)

const (
	UserAgent = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36`
)

//go test -v -run=Test_Spider
func Test_Spider(t *testing.T) {
	//request, err := http.NewRequest(http.MethodGet, "http://testphp.vulnweb.com/", nil)
	request, err := http.NewRequest(http.MethodGet, "http://swj.zjtz.gov.cn", nil)
	if err != nil {
		t.Error(err)
		return
	}

	optionHeader := OptionSpiderRequestHeader("User-Agent", UserAgent)
	result := NewSpider(optionHeader).AddRequest(request).Run().Result()
	data, _ := json.Marshal(result)
	t.Log(string(data))
}
