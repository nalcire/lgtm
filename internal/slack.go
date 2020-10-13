package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SlackHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, 500, fmt.Sprintf("can't read body: %s", err.Error()))
		return
	}

	payload := make(map[string]string)
	err = json.Unmarshal(body, &payload)
	if err != nil {
		errorResponse(w, 500, fmt.Sprintf("can't understand body: %s", err.Error()))
		return
	}

	switch payload["type"] {
	case "url_verification":
		urlVerification(w, payload)
	default:
		return
	}
}

func errorResponse(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

func urlVerification(w http.ResponseWriter, payload map[string]string) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte(payload["challenge"]))
}
