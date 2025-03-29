package jenkins_suggest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HTTPStatusError struct {
	URL    string
	Code   int
	Status string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("bad http status: %d: %s", e.Code, e.Status)
}

// SearchResult 表示整个 JSON 对象
type SearchResult struct {
	Class       string       `json:"_class"`      // 映射到 JSON 中的 "_class" 键
	Suggestions []Suggestion `json:"suggestions"` // 映射到 JSON 中的 "suggestions" 键
}

// Suggestion 表示 suggestions 数组中的每个对象
type Suggestion struct {
	Name string `json:"name"` // 映射到 JSON 中的 "name" 键
}
type Auth struct {
	Username string
	ApiToken string
}

type Jenkins struct {
	auth    *Auth
	baseUrl string
	client  *http.Client
}

func NewJenkins(auth *Auth, baseUrl string) *Jenkins {
	return &Jenkins{
		auth:    auth,
		baseUrl: baseUrl,
		client:  http.DefaultClient,
	}
}

func (jenkins *Jenkins) buildUrl(path string, params url.Values) (requestUrl string) {
	requestUrl = jenkins.baseUrl + path
	if params != nil {
		queryString := params.Encode()
		if queryString != "" {
			requestUrl = requestUrl + "?" + queryString
		}
	}
	return
}

func (jenkins *Jenkins) sendRequest(req *http.Request) (*http.Response, error) {
	if jenkins.auth != nil {
		req.SetBasicAuth(jenkins.auth.Username, jenkins.auth.ApiToken)
	}
	res, err := jenkins.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, &HTTPStatusError{
			URL:    req.URL.String(),
			Code:   res.StatusCode,
			Status: res.Status,
		}
	}
	return res, nil
}

func (jenkins *Jenkins) parseResponse(resp *http.Response, body interface{}) (err error) {
	defer resp.Body.Close()

	if body == nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(data, body)
}

func (jenkins *Jenkins) get(path string, params url.Values, body interface{}) (err error) {
	requestUrl := jenkins.buildUrl(path, params)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return
	}

	resp, err := jenkins.sendRequest(req)
	if err != nil {
		return
	}
	return jenkins.parseResponse(resp, body)
}

// Query 模糊查询
func (jenkins *Jenkins) Query(param string) (names []string, err error) {
	var searchResult SearchResult
	err = jenkins.get(fmt.Sprintf("/search/suggest?query=%s", param), nil, &searchResult)
	suggestions := searchResult.Suggestions
	//var names []string
	for _, suggestion := range suggestions {
		names = append(names, suggestion.Name)
	}
	return
}
