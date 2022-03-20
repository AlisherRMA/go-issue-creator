package flags

import (
	"flag"
)

type UserFlags struct {
	Username string
	Repo     string
	Token    string
}

func GetFlags() UserFlags {
	username := *flag.String("username", "AlisherRMA", "Repo owner's username")
	repo := *flag.String("repo", "go-issue-creator", "Repo name")
	token := *flag.String("token", "", "Github access token")
	flag.Parse()
	return UserFlags{
		Username: username,
		Repo:     repo,
		Token:    token,
	}
}
