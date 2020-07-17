package scrapper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Meta    interface{} `json:"meta"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
}

type Meta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Page   int   `json:"page"`
	Total  int64 `json:"total"`
}

var err error

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

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, meta interface{}, payload interface{}) {
	response := Response{
		Meta:    meta,
		Data:    payload,
		Message: "success",
		Error:   "null"}

	printResponse, err := json.Marshal(response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(printResponse)
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	response := Response{
		Meta:    nil,
		Data:    nil,
		Message: "failed",
		Error:   message}

	printResponse, err := json.Marshal(response)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(printResponse)
}
