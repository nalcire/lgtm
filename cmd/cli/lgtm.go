package main

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nalcire/lgtm/internal"
)

type Specs struct {
	Username string
	Token    string
}

func main() {
	pr := os.Args[1]

	var s Specs
	envconfig.Process("lgtm", &s)
	err := internal.GitHubApprove(pr, s.Username, s.Token)
	if err != nil {
		fmt.Print(err.Error())
	}
}
