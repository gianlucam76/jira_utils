package show

import (
	"context"
	"fmt"
	"jira_utils/jira"
	"strconv"
	"strings"

	docopt "github.com/docopt/docopt-go"
	"k8s.io/klog/v2/klogr"
)

// Filed displays information about issues filed by user (by default user defined in env variable JIRA_USERNAME)
func Filed(ctx context.Context, args []string) error {
	doc := `Usage:
	jira-utils show filed [--sprint=<name>|--active] [--project=<name>] [--board=<name>] [--username=<name>] [--warn-after=<days>]
Options:
  -h --help             Show this screen.
     --active           Show Jira issues in current active sprint.
     --username=<name>  Show Jira issues for specified user (by default user defined in env variable JIRA_USERNAME)
     --sprint=<name>    Show Jira issues in current specified sprint.
     --project=<name>	Show Jira issues in current project (value in JIRA_PROJECT will be used by default)
     --board=<name>     Show Jira issues in current project/board (value in JIRA_BOARD will be used by default)
     --warn-after=<days>  Highlights any issue ii progressing status for more than number of days specified.

Description:
  The show filed command shows information about jira issues filed by user (by default user defined in env variable JIRA_USERNAME)
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

	username := ""
	if passedUsername := parsedArgs["--username"]; passedUsername != nil {
		username = passedUsername.(string)
	} else {
		username = jira.GetUsername(logger)
	}

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

	warnAfter := 0
	if passedWarnAfter := parsedArgs["--warn-after"]; passedWarnAfter != nil {
		warnAfter, err = strconv.Atoi(passedWarnAfter.(string))
		if err != nil {
			return err
		}
	}

	active := parsedArgs["--active"].(bool)
	if active {
		activeSprint, err := jira.GetJiraActiveSprint(ctx, jiraClient, fmt.Sprintf("%d", board.ID), logger)
		if err != nil || activeSprint == nil {
			return fmt.Errorf("failed to get jira active sprint")
		}
		jql = fmt.Sprintf("Status NOT IN (Resolved,Closed) and sprint = %s and reporter = %s", activeSprint.Name, username)
	} else if sprintName != "" {
		sprint, err := jira.GetJiraSprint(ctx, jiraClient, fmt.Sprintf("%d", board.ID), sprintName, logger)
		if err != nil || sprint == nil {
			return fmt.Errorf("%s", fmt.Sprintf("failed to get jira sprint %s", sprintName))
		}
		jql = fmt.Sprintf("Status NOT IN (Resolved,Closed) and sprint = %s and reporter = %s", sprintName, username)
	} else {
		jql = fmt.Sprintf("Status NOT IN (Resolved,Closed) and reporter = %s", username)
	}

	return jira.DisplayJiraIssues(ctx, jiraClient, jql, warnAfter, logger)
}
