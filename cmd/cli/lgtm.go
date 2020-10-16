package main

import (
	"fmt"
	"os"
	"os/user"

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
	username := "unknown"
	u, err := user.Current()
	if err == nil {
		username = u.Name
	}

	err = internal.GitHubApprove(username, pr, s.Username, s.Token)
	if err != nil {
		fmt.Print(err.Error())
	}
}
