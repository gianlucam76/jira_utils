# jira_utils

Use make build to build it.

Some environment variables need to be set. jira_utils checks for those and prints which one is missing, if any.

```
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
```

Then

```
./bin/jira_utils --help
```
