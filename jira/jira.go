package jira

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/fatih/color"
	"github.com/go-logr/logr"
	"github.com/olekukonko/tablewriter"
)

const (
	// jiraBaseURL is the env variable containing the jira base URL
	jiraBaseURL = "JIRA_BASE_URL"
	// jiraProject is the name of the env variable with the project name
	jiraProject = "JIRA_PROJECT"
	// jiraBoardName is the name of the env variable with the board name
	jiraBoardName = "JIRA_BOARD"
	// username is the name of the env variable with the username
	username = "JIRA_USERNAME"
	// password is the name of the env variable with the password (base64 encoded)
	password = "JIRA_PASSWORD"
)

// VerifyEnvVariables verifies all needed environment variables are set
func VerifyEnvVariables(logger logr.Logger) {
	logger.Info("Verifying all needed environment variables are set")

	logger.V(5).Info(fmt.Sprintf("Verify %s (this must contain the jira base url)", jiraBaseURL))
	if _, ok := os.LookupEnv(jiraBaseURL); !ok {
		logger.Info(fmt.Sprintf("Env variable %s not found.", jiraBaseURL))
		panic(1)
	}

	logger.V(5).Info(fmt.Sprintf("Verify %s (this must contain the jira project name)", jiraProject))
	if _, ok := os.LookupEnv(jiraProject); !ok {
		logger.Info(fmt.Sprintf("Env variable %s not found.", jiraProject))
		panic(1)
	}

	logger.V(5).Info(fmt.Sprintf("Verify %s (this must contain the board name)", jiraBoardName))
	if _, ok := os.LookupEnv(jiraBoardName); !ok {
		logger.Info(fmt.Sprintf("Env variable %s not found.", jiraBoardName))
		panic(1)
	}

	logger.V(5).Info(fmt.Sprintf("Verify %s (this must contain the username)", username))
	if _, ok := os.LookupEnv(username); !ok {
		logger.Info(fmt.Sprintf("Env variable %s not found.", username))
		panic(1)
	}

	logger.V(5).Info(fmt.Sprintf("Verify %s (this must password, base64 encoded)", password))
	if _, ok := os.LookupEnv(password); !ok {
		logger.Info(fmt.Sprintf("Env variable %s not found.", password))
		panic(1)
	}
}

// GetJiraClient returns a new Jira API client.
func GetJiraClient(ctx context.Context, username, password string, logger logr.Logger) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	baseURL, ok := os.LookupEnv(jiraBaseURL)
	if !ok {
		msg := fmt.Sprintf("Env variable %s not found.", jiraBaseURL)
		logger.Info(msg)
		return nil, fmt.Errorf("%s", msg)
	}

	jiraClient, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to get jira client. Err: %v", err))
		return nil, err
	}

	return jiraClient, nil
}

// GetJiraProject returns the jira.Project with name projectName
func GetJiraProject(ctxt context.Context, jiraClient *jira.Client, projectName string, logger logr.Logger) (*jira.Project, error) {
	if projectName == "" {
		var ok bool
		projectName, ok = os.LookupEnv(jiraProject)
		if !ok || projectName == "" {
			msg := fmt.Sprintf("ProjectName was not passed and env variable %s is not set", jiraProject)
			logger.Info(msg)
			return nil, fmt.Errorf("%s", msg)
		}
	}

	url := fmt.Sprintf("rest/api/2/project/%s", projectName)
	req, _ := jiraClient.NewRequest("GET", url, nil)
	project := &jira.Project{}
	if _, err := jiraClient.Do(req, project); err != nil {
		logger.Info(fmt.Sprintf("Failed to get project with name: %s. Error: %v", projectName, err))
		return nil, err
	}

	return project, nil
}

// GetJiraBoard returns board with name boardName in project projectKey
// returns the board if only one is found or an error if any occurs.
// Returns nil if no board is found or more than one is found
func GetJiraBoard(ctx context.Context, jiraClient *jira.Client, projectKey, boardName string, logger logr.Logger) (*jira.Board, error) {
	if boardName == "" {
		var ok bool
		boardName, ok = os.LookupEnv(jiraBoardName)
		if !ok || boardName == "" {
			msg := fmt.Sprintf("boardName was not passed and env variable %s is not set", jiraBoardName)
			logger.Info(msg)
			return nil, fmt.Errorf("%s", msg)
		}
	}

	boardListOptions := &jira.BoardListOptions{ProjectKeyOrID: projectKey, Name: boardName}
	boardList, _, err := jiraClient.Board.GetAllBoards(boardListOptions)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to get board list. Error %v", err))
		return nil, err
	}

	if boardList.Values == nil {
		logger.Info(fmt.Sprintf("Got not result for GetAllBoards with projectKey: %s and boardName: %s ", projectKey, boardName))
		return nil, nil
	}

	if len(boardList.Values) != 1 {
		logger.Info(fmt.Sprintf("Got more than one result for GetAllBoards with projectKey: %s and boardName: %s ", projectKey, boardName))
		logger.Info(fmt.Sprintf("Result: %v", boardList.Values))
		return nil, nil
	}

	return &boardList.Values[0], nil
}

// GetJiraActiveSprint returns the active sprint for passed in board
// Returns active sprint if found or an error if any occurs.
// If no sprint is currently active, returns nil
func GetJiraActiveSprint(ctx context.Context, jiraClient *jira.Client, boardID string, logger logr.Logger) (*jira.Sprint, error) {
	if jiraClient == nil {
		msg := "jiraClient is nil"
		logger.Info(msg)
		return nil, fmt.Errorf(msg)
	}

	sprints, _, err := jiraClient.Board.GetAllSprintsWithContext(ctx, boardID)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to get board list. Error: %v", err))
		return nil, err
	}

	var activeSprint *jira.Sprint
	now := time.Now()
	for i := range sprints {
		if sprints[i].StartDate != nil && sprints[i].EndDate != nil {
			if sprints[i].StartDate.Before(now) && sprints[i].EndDate.After(now) {
				return &sprints[i], nil
			} else if sprints[i].StartDate.Before(now) {
				if activeSprint == nil {
					activeSprint = &sprints[i]
				} else if sprints[i].EndDate.After(*activeSprint.EndDate) {
					activeSprint = &sprints[i]
				}
			}
		}
	}

	return activeSprint, nil
}

// GetJiraSprints returns the specified sprint for passed in board
// Returns sprints or an error if any occurs.
func GetJiraSprints(ctx context.Context, jiraClient *jira.Client, boardID string, logger logr.Logger) ([]jira.Sprint, error) {
	if jiraClient == nil {
		msg := "jiraClient is nil"
		logger.Info(msg)
		return nil, fmt.Errorf(msg)
	}

	sprints, _, err := jiraClient.Board.GetAllSprintsWithContext(ctx, boardID)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to get board list. Error: %v", err))
		return nil, err
	}

	return sprints, nil
}

// GetJiraSprint returns all sprints for passed in board
// Returns sprint if found or an error if any occurs.
// If no matching sprint is found, returns nil
func GetJiraSprint(ctx context.Context, jiraClient *jira.Client, boardID, sprintName string, logger logr.Logger) (*jira.Sprint, error) {
	if jiraClient == nil {
		msg := "jiraClient is nil"
		logger.Info(msg)
		return nil, fmt.Errorf(msg)
	}

	sprints, _, err := jiraClient.Board.GetAllSprintsWithContext(ctx, boardID)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to get board list. Error: %v", err))
		return nil, err
	}

	for i := range sprints {
		if sprints[i].Name == sprintName {
			return &sprints[i], nil
		}
	}

	return nil, nil
}

// GetJiraIssues finds all issues matching passed jql
func GetJiraIssues(ctx context.Context, jiraClient *jira.Client, jql string, logger logr.Logger) ([]jira.Issue, error) {
	issues, _, err := jiraClient.Issue.SearchWithContext(ctx, jql, nil)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to get all issues matching jql:%s. Error: %v", jql, err))
		return nil, err
	}

	return issues, nil
}

// CreateIssue creates new issue of type bug which will be added to sprint
// - Comments will contain buildEnvironment (VCS vs UCS), run ID, failure message and full stack trace
// - Assignee is the user the bug will be assigned to
// - Reporter is the issue reporter
// Return the issue Key or empty an error occurred.
func CreateIssue(ctx context.Context, jiraClient *jira.Client, sprint *jira.Sprint, priority *jira.Priority,
	projectKey, componentName, assignee, testName, summary string,
	logger logr.Logger) (string, error) {
	component := jira.Component{Name: componentName}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: assignee,
			},
			Description: fmt.Sprintf("Test %s failed", testName),
			Type: jira.IssueType{
				Name: "Bug",
			},
			Project: jira.Project{
				Key: projectKey,
			},
			Components: []*jira.Component{
				&component,
			},
			Summary:  summary,
			Priority: priority,
		},
	}

	issue, resp, err := jiraClient.Issue.CreateWithContext(ctx, &i)
	if err != nil {
		body, _ := io.ReadAll(resp.Body)
		logger.Info(fmt.Sprintf("Failed to create issue. Error: %v. Resp %s", err, string(body)))
		return "", err
	}

	logger.Info(fmt.Sprintf("Created issue %s", issue.Key))

	return issue.Key, nil
}

// AddCommentToIssue append comment to current open issue while also resetting sprint and priority.
func AddCommentToIssue(ctx context.Context, jiraClient *jira.Client, issueID string,
	commentMsg string, logger logr.Logger) error {
	comment := jira.Comment{
		Body: commentMsg,
	}

	if _, resp, err := jiraClient.Issue.AddCommentWithContext(ctx, issueID, &comment); err != nil {
		body, _ := io.ReadAll(resp.Body)
		logger.Info(fmt.Sprintf("Failed to update issue %s. Error: %v. Resp %s", issueID, err, string(body)))
		return err
	}

	return nil
}

func MoveIssueToSprint(ctx context.Context, jiraClient *jira.Client, sprintID int, issueID string, logger logr.Logger) error {
	if resp, err := jiraClient.Sprint.MoveIssuesToSprintWithContext(ctx, sprintID, []string{issueID}); err != nil {
		body, _ := io.ReadAll(resp.Body)
		logger.Info(fmt.Sprintf("Failed to update issue %s. Error: %v. Resp %s", issueID, err, string(body)))
		return err
	}

	return nil
}

func DisplayJiraIssues(ctx context.Context, jiraClient *jira.Client, jql string, warnAfter int, logger logr.Logger) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"KEY", "SUMMARY", "STATUS", "LAST UPDATE", "ASSIGNEE"})
	table.SetAutoWrapText(false)
	table.SetRowLine(true)

	issues, err := GetJiraIssues(ctx, jiraClient, jql, logger)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		logger.Info("No issue found")
	}

	for i := range issues {
		key := issues[i].Key
		warning := false
		if warnAfter != 0 && shouldWarn(jiraClient, &issues[i], warnAfter) {
			warning = true
		}
		var summary, status, username string
		if issues[i].Fields != nil {
			summary = issues[i].Fields.Summary
			if issues[i].Fields.Status != nil {
				status = issues[i].Fields.Status.Name
			} else {
				status = "N/A"
			}
			if issues[i].Fields.Assignee != nil {
				username = issues[i].Fields.Assignee.Name
			}
		}

		// lastUpdate := string([]byte(time.Time(issues[i].Fields.Updated).Format("\"2006-01-02T\"")))
		lastUpdateTime := time.Time(issues[i].Fields.Updated)
		lastUpdate := fmt.Sprintf("%d days", lastUpdateTime.Day())

		if warning {
			table.Append([]string{color.New(color.FgRed).Sprintf(key),
				color.New(color.FgRed).Sprintf(summary), color.New(color.FgRed).Sprintf(status),
				color.New(color.FgRed).Sprintf(lastUpdate), color.New(color.FgRed).Sprintf(username)})
		} else {
			table.Append([]string{key, summary, status, lastUpdate, username})
		}
	}

	table.Render()
	return nil
}

func shouldWarn(jiraClient *jira.Client, issue *jira.Issue, warnAfter int) bool {
	var cIssue *jira.Issue
	var err error
	if issue.Fields.Status != nil &&
		issue.Fields.Status.Name == "In Progress" {
		cIssue, _, err = jiraClient.Sprint.GetIssue(issue.ID, &jira.GetQueryOptions{Expand: "changelog"})
		if err != nil {
			return false
		}
	} else {
		return false
	}

	if cIssue.Changelog == nil {
		return false
	}

	hours := time.Duration(-24 * warnAfter)

	for i := range cIssue.Changelog.Histories {
		history := cIssue.Changelog.Histories[i]
		historyTime, err := history.CreatedTime()
		if err != nil {
			continue
		}
		for j := range history.Items {
			logItem := history.Items[j]
			if logItem.ToString == "In Progress" &&
				historyTime.Before(time.Now().Add(-hours*time.Hour)) {
				return true
			}
		}
	}
	return false
}
