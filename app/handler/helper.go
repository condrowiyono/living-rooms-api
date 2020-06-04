package handler

import (
	"io/ioutil"
	"net/http"
)

func getHTTPRequestGetBody(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func getHTTPRequestQuery(r *http.Request, query string) string {
	vars := r.URL.Query()
	value := string(vars.Get(query))

	return value
}
