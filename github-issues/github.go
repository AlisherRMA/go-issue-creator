package github_issues

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github-issue-creator/flags"
	"io"
	"net/http"
)

type GithubService struct {
	token    string
	repo     string
	username string
	url      string
}

func (g *GithubService) New(flags flags.UserFlags) {
	g.token = flags.Token
	g.username = flags.Username
	g.repo = flags.Repo
	g.url = fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", g.username, g.repo)
}

func (g *GithubService) SearchIssues() ([]*Issue, error) {
	resp, err := http.Get(g.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сбой запроса: %s", resp.Status)
	}

	var result []*Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (g *GithubService) GetIssueById(id string) (*Issue, error) {
	resp, err := http.Get(g.url + "/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сбой запроса: %s", resp.Status)
	}

	var result *Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (g *GithubService) CreateIssue(title string, body string) error {
	values := map[string]string{"title": title, "body": body}
	jsonValue, _ := json.Marshal(values)

	req, err := g.MyRequest("POST", g.url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("сбой запроса: %s", resp.Status)
	}

	var result *Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	fmt.Println("Issue created successfully")
	return nil
}

func (g *GithubService) UpdateIssue(title, body string, id int) error {
	values := map[string]string{"title": title, "body": body}
	jsonValue, _ := json.Marshal(values)
	fmt.Printf("\n %+v \n", values)
	req, err := g.MyRequest("PATCH", fmt.Sprintf("%s/%d", g.url, id), bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сбой запроса: %s", resp.Status)
	}

	var result *Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	fmt.Println("Issue updated successfully")
	return nil
}

func (g *GithubService) MyRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "Bearer "+g.token)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}
