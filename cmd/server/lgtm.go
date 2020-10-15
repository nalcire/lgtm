package main

import (
	"log"
	"net/http"

	"github.com/nalcire/lgtm/internal"
)

func main() {
	s := internal.NewServer("lgtm")
	http.HandleFunc("/slack/lgtm", s.SlackHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
