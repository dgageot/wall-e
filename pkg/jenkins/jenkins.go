package jenkins

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetJenkinsCrumb(user string, token string, server string) (string, error) {
	var url = fmt.Sprintf("https://%v:%v@%v/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)", user, token, server)

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("http response code (%v)", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Get(user, token, server, uri string) (*http.Response, error) {
	var url = fmt.Sprintf("https://%s:%s@%s%s", user, token, server, uri)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	crumb, err := GetJenkinsCrumb(user, token, server)
	if err != nil {
		return nil, err
	}

	crumbKeyValue := strings.Split(crumb, ":")
	request.Header.Set(crumbKeyValue[0], crumbKeyValue[1])

	return http.DefaultClient.Do(request)
}
