package flags

import (
	"flag"
	"fmt"
)

type UserFlags struct {
	Username string
	Repo     string
	Token    string
}

func GetFlags() UserFlags {
	flags := UserFlags{}

	flag.StringVar(&flags.Username, "username", "AlisherRMA", "Repo owner's username")
	flag.StringVar(&flags.Repo, "repo", "go-issue-creator", "Repo name")
	flag.StringVar(&flags.Token, "token", "", "Github access token")
	flag.Parse()
	fmt.Println(flags)
	return flags
}
