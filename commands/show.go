package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	docopt "github.com/docopt/docopt-go"

	"github.com/gianlucam76/jira_utils/commands/show"
)

// Show takes keyword then calls subcommand.
func Show(ctx context.Context, args []string) error {
	doc := `Usage:
	jira-utils show <command> [<args>...]

    issues           show jira issues.
    filed            show jira issues filed by user.
    sprints          show all sprints.
    e2e              show all open issues filed for e2e.

Options:
	-h --help      Show this screen.

Description:
	See 'jira-utils show <command> --help' to read about a specific subcommand.
  `
	parser := &docopt.Parser{
		HelpHandler:   docopt.PrintHelpAndExit,
		OptionsFirst:  true,
		SkipHelpFlags: false,
	}

	opts, err := parser.ParseArgs(doc, nil, "1.0")
	if err != nil {
		if _, ok := err.(*docopt.UserError); ok {
			fmt.Printf(
				"Invalid option: 'jira-util %s'. Use flag '--help' to read about a specific subcommand.\n",
				strings.Join(os.Args[1:], " "),
			)
		}
		os.Exit(1)
	}

	command := opts["<command>"].(string)
	arguments := append([]string{"show", command}, opts["<args>"].([]string)...)

	switch command {
	case "issues":
		return show.Issues(ctx, arguments)
	case "filed":
		return show.Filed(ctx, arguments)
	case "sprints":
		return show.Sprints(ctx, arguments)
	case "e2e":
		return show.E2EIssues(ctx, arguments)
	default:
		fmt.Println(doc)
	}

	return nil
}
