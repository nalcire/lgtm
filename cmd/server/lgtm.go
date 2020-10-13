package main

import (
	"log"
	"net/http"

	"github.com/nalcire/lgtm/internal"
)

func main() {
	http.HandleFunc("/slack/lgtm", internal.SlackHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
