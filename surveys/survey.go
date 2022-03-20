package surveys

import (
	"fmt"
	"github-issue-creator/flags"
	github_issues "github-issue-creator/github-issues"
	"github.com/AlecAivazis/survey/v2"
	"log"
	"strings"
)

type BaseSurveyAnswers struct {
	Action string `survey:"color"` // or you can tag fields to match a specific name
}

var qs = []*survey.Question{
	//{
	//	Name:   "githubToken",
	//	Prompt: &survey.Input{Message: "Please enter your github access token", Help: "You can find more information here: https://docs.github.com/en/enterprise-server@3.3/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token"},
	//},
	{
		Name: "action",
		Prompt: &survey.Select{
			Message: "Choose action:",
			Options: []string{"create_issue", "get_issues", "get_issue_by_id"},
			Default: "red",
		},
	},
}

func Start(flags flags.UserFlags) {
	// the answers will be written to this struct
	answers := BaseSurveyAnswers{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(flags)

	github := new(github_issues.GithubService)
	github.New(flags)
	handleSelection(&answers, github)
}

func handleSelection(answers *BaseSurveyAnswers, client *github_issues.GithubService) {
	if answers.Action == "get_issues" {
		fmt.Printf("loading...\n")
		issueId, err := getIssuesList(client)
		if err != nil {
			log.Fatal(err)
		}
		getIssueById(client, issueId)
	} else if answers.Action == "get_issue_by_id" {
		var issueId string
		err := survey.AskOne(&survey.Input{Message: "Please, type issue number"}, &issueId)
		if err != nil {
			log.Fatal(err)
		}
		getIssueById(client, issueId)
	} else {
		createIssue(client)
	}
}

func getIssuesList(client *github_issues.GithubService) (string, error) {
	res, err := client.SearchIssues()
	if err != nil {
		fmt.Println("Error occurred")
		log.Fatal(err)
	}
	issueLabels := make([]string, 0, len(res))
	for _, item := range res {
		issueLabels = append(issueLabels, fmt.Sprintf("%d: %s", item.Number, item.Title))
	}

	var selectedIssue string

	err = survey.AskOne(&survey.Select{
		Message: "Choose issue:",
		Options: issueLabels,
	}, &selectedIssue)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	fmt.Printf("issue: %s", selectedIssue)
	num := strings.Split(selectedIssue, ":")

	return num[0], nil
}

func getIssueById(client *github_issues.GithubService, issueId string) {
	issue, err := client.GetIssueById(issueId)
	if err != nil {
		fmt.Println("Error occurred")
		log.Fatal(err)
	}
	fmt.Printf("\n issue_id: %s \n title: %s \n created_by: %s \n created_at: %s \n body: %s",
		issueId, issue.Title, issue.User.Login, issue.CreatedAt, issue.Body)
}

func createIssue(client *github_issues.GithubService) {
	payload := []*survey.Question{
		{
			Name: "title",
			Prompt: &survey.Input{
				Message: "Type issue title",
			},
			Validate: survey.Required,
		},
		{
			Name: "body",
			Prompt: &survey.Input{
				Message: "Issue body",
			},
			Validate: survey.Required,
		},
	}

	payloadAnswers := struct {
		Title string
		Body  string
	}{}

	survey.Ask(payload, &payloadAnswers)

	err := client.CreateIssue(payloadAnswers.Title, payloadAnswers.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Issue created successfully")
}
