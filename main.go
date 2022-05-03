package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"

	docopt "github.com/docopt/docopt-go"

	"github.com/gianlucam76/jira_utils/commands"
	"github.com/gianlucam76/jira_utils/jira"
)

func main() {
	ctx := context.Background()

	klog.InitFlags(nil)
	logger := klogr.New()

	jira.VerifyEnvVariables(logger)

	doc := `Usage:
	jira_utils [options] <command> [<args>...]

	show          Display information on jira issues

Options:
  -h --help     Show this screen.
     --version  Show version.

Description:
  The jira-utils command line tool is used to manage/display jira issues.
  See 'jira-utils <command> --help' to read about a specific subcommand.
`

	parser := &docopt.Parser{
		HelpHandler:   docopt.PrintHelpOnly,
		OptionsFirst:  true,
		SkipHelpFlags: false,
	}

	opts, err := parser.ParseArgs(doc, nil, "")
	if err != nil {
		if _, ok := err.(*docopt.UserError); ok {
			fmt.Printf(
				"Invalid option: 'jira-util %s'. Use flag '--help' to read about a specific subcommand.\n",
				strings.Join(os.Args[1:], " "),
			)
		}
		os.Exit(1)
	}

	if opts["<command>"] != nil {
		command := opts["<command>"].(string)
		args := append([]string{command}, opts["<args>"].([]string)...)
		var err error

		switch command {
		case "show":
			err = commands.Show(ctx, args)
		default:
			err = fmt.Errorf("unknown command: %q\n%s", command, doc)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
}
