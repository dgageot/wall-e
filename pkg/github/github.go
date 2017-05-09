package github

import (
	"fmt"
	"net/http"
)

func Get(token, uri string) (*http.Response, error) {
	var url = fmt.Sprintf("https://api.github.com%s", uri)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "token "+token)

	return http.DefaultClient.Do(request)
}
