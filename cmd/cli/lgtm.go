package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nalcire/lgtm/internal/github"
)

type Specs struct {
	Username string
	Password string
}

func main() {
	pr := os.Args[1]

	var s Specs
	envconfig.Process("lgtm", &s)
	github.Approve(pr, s.Username, s.Password)
}
