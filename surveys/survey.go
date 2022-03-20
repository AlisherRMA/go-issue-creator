package surveys

import (
	"fmt"
	"github-issue-creator/colors"
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
	if answers.Action == "create_issue" {
		createIssue(client)
	} else if answers.Action == "get_issue_by_id" {
		var issueId string
		err := survey.AskOne(&survey.Input{Message: "Please, type issue number"}, &issueId)
		if err != nil {
			log.Fatal(err)
		}
		getIssueById(client, issueId)
	} else {
		fmt.Printf("loading...\n")
		issueId, err := getIssuesList(client)
		if err != nil {
			log.Fatal(err)
		}
		getIssueById(client, issueId)
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
	fmt.Printf("\n\n")
	colors.ColorizeOutput("issue_id  ", issueId)
	colors.ColorizeOutput("title     ", issue.Title)
	colors.ColorizeOutput("created_by", issue.User.Login)
	colors.ColorizeOutput("created_at", issue.CreatedAt.String())
	colors.ColorizeOutput("body      ", issue.Body)
	fmt.Printf("\n\n")

	afterIssueSelected(client, issue)
}

func afterIssueSelected(client *github_issues.GithubService, issue *github_issues.Issue) {
	var nextAxtion string
	survey.AskOne(&survey.Select{
		Message: "What do you want to do next?",
		Options: []string{"update_issue", "read_another_issue", "exit"},
	}, &nextAxtion)

	switch nextAxtion {
	case "update_issue":
		updateIssue(client, issue)
	case "read_another_issue":
		handleSelection(&BaseSurveyAnswers{"get_issues"}, client)
	}
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
}

func updateIssue(client *github_issues.GithubService, issue *github_issues.Issue) {
	questions := []*survey.Question{
		{
			Name: "title",
			Prompt: &survey.Input{
				Message: "Enter a title or press Enter to leave as is",
				Default: issue.Title,
			},
		},
		{
			Name: "body",
			Prompt: &survey.Input{
				Message: "Enter a body or press Enter to leave as is",
				Default: issue.Body,
			},
		},
	}

	answers := struct {
		Title string
		Body  string
	}{}

	survey.Ask(questions, &answers)
	err := client.UpdateIssue(answers.Title, answers.Body, issue.Number)
	if err != nil {
		log.Fatal(err)
	}
}
