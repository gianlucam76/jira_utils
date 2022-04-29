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

// E2EIssues displays information about issues filed for e2e automatic tagging sanities
func E2EIssues(ctx context.Context, args []string) error {
	doc := `Usage:
	jira-utils show e2e [--warn-after=<days>]
Options:
  -h --help               Show this screen.
     --warn-after=<days>  Highlights any issue ii progressing status for more than number of days specified.

Description:
  The show e2e command shows information about jira issues filed for e2e automatic tagging sanities
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

	username := "atom-ci.gen"

	jiraClient, err := jira.GetJiraClient(ctx, jira.GetUsername(logger), jira.GetPassword(logger), logger)
	if err != nil {
		return err
	}

	projectName := ""
	project, err := jira.GetJiraProject(ctx, jiraClient, projectName, logger)
	if err != nil || project == nil {
		return fmt.Errorf("failed to get jira project")
	}

	boardName := ""
	board, err := jira.GetJiraBoard(ctx, jiraClient, project.Key, boardName, logger)
	if err != nil || board == nil {
		return fmt.Errorf("failed to get jira board")
	}

	var jql string

	warnAfter := 0
	if passedWarnAfter := parsedArgs["--warn-after"]; passedWarnAfter != nil {
		warnAfter, err = strconv.Atoi(passedWarnAfter.(string))
		if err != nil {
			return err
		}
	}

	jql = fmt.Sprintf("Status NOT IN (Resolved,Closed) and reporter = %s", username)

	return jira.DisplayJiraIssues(ctx, jiraClient, jql, warnAfter, logger)
}
