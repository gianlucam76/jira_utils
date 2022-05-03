package show

import (
	"context"
	"fmt"
	"os"
	"strings"

	docopt "github.com/docopt/docopt-go"
	"github.com/olekukonko/tablewriter"
	"k8s.io/klog/v2/klogr"

	"github.com/gianlucam76/jira_utils/jira"
)

// Sprints displays information about issues in a given sprint
func Sprints(ctx context.Context, args []string) error {
	doc := `Usage:
	jira-utils show sprints [--project=<name>] [--board=<name>]
Options:
  -h --help       Show this screen.
     --project    Show Jira issues in current project (value in JIRA_PROJECT will be used by default)
     --board      Show Jira issues in current project/board (value in JIRA_BOARD will be used by default)

Description:
  The show sprints command shows information about jira issues.
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

	sprints, err := jira.GetJiraSprints(ctx, jiraClient, fmt.Sprintf("%d", board.ID), logger)
	if err != nil {
		return fmt.Errorf("failed to get jira sprints")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"SPRINT", "STATE"})
	table.SetReflowDuringAutoWrap(false)
	table.SetRowLine(true)

	for i := range sprints {
		table.Append([]string{sprints[i].Name, sprints[i].State})
	}

	table.Render()
	return nil
}
