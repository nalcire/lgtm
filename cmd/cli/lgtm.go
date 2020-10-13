package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nalcire/lgtm/internal"
)

type Specs struct {
	Username string
	Password string
}

func main() {
	pr := os.Args[1]

	var s Specs
	envconfig.Process("lgtm", &s)
	internal.GitHubApprove(pr, s.Username, s.Password)
}
