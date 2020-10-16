package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/kelseyhightower/envconfig"
)

type EventType struct {
	Type  string
	Event struct {
		Type string
	}
}

type MessageCallback struct {
	Type  string
	Event struct {
		Type    string
		Text    string
		Channel string
		User    string
	}
}

type Server struct {
	GithubUser  string
	GithubToken string
	SlackToken  string
}

func NewServer(name string) *Server {
	s := Server{}
	envconfig.Process(name, &s)
	return &s
}

func (s *Server) SlackHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("can't read body: %s", err.Error())
		errorResponse(w, 500, fmt.Sprintf("can't read body: %s", err.Error()))
		return
	}

	var t EventType
	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Printf("can't understand body: %s", err.Error())
		errorResponse(w, 500, fmt.Sprintf("can't understand body: %s", err.Error()))
		return
	}

	switch t.Type {
	case "url_verification":
		urlVerification(w, body)
	case "event_callback":
		switch t.Event.Type {
		case "message":
			s.message(w, body)
		}
	default:
		return
	}
}

func errorResponse(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

func urlVerification(w http.ResponseWriter, event []byte) {
	payload := make(map[string]string)
	json.Unmarshal(event, &payload)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte(payload["challenge"]))
}

var rStamp *regexp.Regexp = regexp.MustCompile(`^stamp (.*)$`)

func (s *Server) message(w http.ResponseWriter, event []byte) {
	var payload MessageCallback
	err := json.Unmarshal(event, &payload)
	if err != nil {
		log.Printf("can't unmarshal event: %s", err.Error())
		errorResponse(w, 500, "don't understand body")
	}

	m := rStamp.FindStringSubmatch(payload.Event.Text)
	if len(m) == 0 {
		w.WriteHeader(200)
		return
	}

	// can get username from id and cache it
	err = GitHubApprove(payload.Event.User, m[1], s.GithubUser, s.GithubToken)
	text := "👍 LGTM"
	if err != nil {
		text = err.Error()
	}

	url := fmt.Sprintf("https://slack.com/api/chat.postMessage?token=%s&channel=%s&text=%s", s.SlackToken, payload.Event.Channel, url.QueryEscape(text))
	resp, err := http.Get(url)
	w.WriteHeader(200)

	respBody, _ := ioutil.ReadAll(resp.Body)
	log.Printf("slack response: %s %s", resp.Status, respBody)
}
