package show

import (
	"context"
	"fmt"
	"jira_utils/jira"
	"strings"

	docopt "github.com/docopt/docopt-go"
	"k8s.io/klog/v2/klogr"
)

// Assigned displays information about issues assigned to user
func Assigned(ctx context.Context, args []string) error {
	doc := `Usage:
	jira-utils show assigned [--sprint=<name>|--active] [--project=<name>] [--board=<name>]
Options:
  -h --help             Show this screen.
     --active           Show Jira issues in current active sprint.
     --sprint=<name>    Show Jira issues in current specified sprint.
     --project=<name>	Show Jira issues in current project (value in JIRA_PROJECT will be used by default)
     --board=<name>		Show Jira issues in current project/boardf (value in JIRA_BOARD will be used by default)

Description:
  The show assigned command shows information about jira issues assigned to user..
`
	parsedArgs, err := docopt.ParseArgs(doc, nil, "1.0")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf(
			"invalid option: 'jira-utils %s'. Use flag '--help' to read about a specific subcommand. Error: %v",
			strings.Join(args, " "),
			err,
		)
	}
	if len(parsedArgs) == 0 {
		return nil
	}

	logger := klogr.New()

	jiraClient, err := jira.GetJiraClient(ctx, jira.GetUsername(logger), jira.GetPassword(logger), logger)
	if err != nil {
		return err
	}

	projectName := ""
	if passedProject := parsedArgs["--project"]; passedProject != nil {
		projectName = passedProject.(string)
	}

	project, err := jira.GetJiraProject(ctx, jiraClient, projectName, logger)
	if err != nil || project == nil {
		return fmt.Errorf("failed to get jira project")
	}

	boardName := ""
	if passedBoard := parsedArgs["--board"]; passedBoard != nil {
		boardName = passedBoard.(string)
	}

	board, err := jira.GetJiraBoard(ctx, jiraClient, project.Key, boardName, logger)
	if err != nil || board == nil {
		return fmt.Errorf("failed to get jira board")
	}

	var jql string

	sprintName := ""
	if passedSprint := parsedArgs["--sprint"]; passedSprint != nil {
		sprintName = passedSprint.(string)
	}

	active := parsedArgs["--active"].(bool)
	if active {
		activeSprint, err := jira.GetJiraActiveSprint(ctx, jiraClient, fmt.Sprintf("%d", board.ID), logger)
		if err != nil || activeSprint == nil {
			return fmt.Errorf("failed to get jira active sprint")
		}
		jql = fmt.Sprintf("assignee = %s and Status NOT IN (Resolved,Closed) and sprint = %s",
			jira.GetUsername(logger), activeSprint.Name)
	} else if sprintName != "" {
		sprint, err := jira.GetJiraSprint(ctx, jiraClient, fmt.Sprintf("%d", board.ID), sprintName, logger)
		if err != nil || sprint == nil {
			return fmt.Errorf("%s", fmt.Sprintf("failed to get jira sprint %s", sprintName))
		}
		jql = fmt.Sprintf("assignee = %s and Status NOT IN (Resolved,Closed) and sprint = %s", jira.GetUsername(logger), sprintName)
	} else {
		jql = fmt.Sprintf("assignee = %s and Status NOT IN (Resolved,Closed)", jira.GetUsername(logger))
	}
	return jira.DisplayJiraIssues(ctx, jiraClient, jql, logger)
}
