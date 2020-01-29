package scripting

import (
	"bytes"
	"fmt"
	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const HttpModule = "http.sky"

var client http.Client

func init() {
	client = http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			MaxIdleConns:        100,
			IdleConnTimeout:     30 * time.Second,
		},
	}
}

func LoadHttp() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(
		skylarkstruct.Default,
		skylark.StringDict{
			"request": skylark.NewBuiltin("request", request),
		},
	)
}

func makeUrl(uri skylark.String, query *skylark.Dict) string {
	u, err := url.Parse(AsString(uri))
	if err != nil {
		return ""
	}

	keys := query.Keys()
	if len(keys) > 0 {
		q := u.Query()
		for _, key := range keys {
			if strings.EqualFold(key.Type(), "string") {
				value, found, err := query.Get(key)
				if nil == err && found && strings.EqualFold(value.Type(), "string") {
					q.Set(AsString(key), AsString(value))
				}
			} else {
				return ""
			}
		}
		u.RawQuery = q.Encode()
	}

	return u.String()
}

func putHeaders(request *http.Request, header *skylark.Dict) {
	keys := header.Keys()
	if len(keys) > 0 {
		for _, key := range keys {
			if strings.EqualFold(key.Type(), "string") {
				value, found, err := header.Get(key)
				if nil == err && found && strings.EqualFold(value.Type(), "string") {
					request.Header.Add(AsString(key), AsString(value))
				}
			}
		}
	}
}

func request(
	thread *skylark.Thread,
	builtin *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var method, url, body skylark.String
	var query = &skylark.Dict{}
	var header = &skylark.Dict{}
	err := skylark.UnpackArgs("request", args, kwargs, "method", &method, "url", &url, "query", &query, "header", &header, "body", &body)
	if nil != err {
		return skylark.None, err
	}
	uri := makeUrl(url, query)
	if len(uri) == 0 {
		return skylark.None, fmt.Errorf("bad url or query")
	}
	request, err := http.NewRequest(strings.ToUpper(AsString(method)), uri, bytes.NewBufferString(AsString(body)))
	if nil != err {
		return skylark.None, err
	}
	putHeaders(request, header)
	response, err := client.Do(request)
	if nil != err {
		return skylark.None, err
	}
	defer response.Body.Close()
	code := response.StatusCode
	content, err := ioutil.ReadAll(response.Body)
	if nil != err {
		content = []byte("")
	}
	return skylarkstruct.FromStringDict(
		skylarkstruct.Default,
		skylark.StringDict{
			"status": skylark.MakeInt(code),
			"body":   skylark.String(content),
		},
	), nil
}
