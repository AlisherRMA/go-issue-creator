package main

import (
	"github-issue-creator/flags"
	"github-issue-creator/surveys"
)

func main() {
	myFlags := flags.GetFlags()

	surveys.Start(myFlags)
}
