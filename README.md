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

To build,

```
make build
```

Then

```
./bin/jira_utils --help
```

Few examples.
To list all jira items assigned to yourself in current sprint:

```
./bin/jira_utils show issues --active 
I0429 11:00:12.453300     273 jira.go:30]  "msg"="Verifying all needed environment variables are set"  
+-----------------+----------------------------------------------------+---------+-------------+----------+
|       KEY       |                      SUMMARY                       | STATUS  | LAST UPDATE | ASSIGNEE |
+-----------------+----------------------------------------------------+---------+-------------+----------+
| CLOUDSTACK-2355 | Review GlobalClusterConfig update PR               | Backlog | 26 days     | mgianluc |
+-----------------+----------------------------------------------------+---------+-------------+----------+
| CLOUDSTACK-2349 | e2e: add dex information to each workload cluster  | Backlog | 26 days     | mgianluc |
+-----------------+----------------------------------------------------+---------+-------------+----------+
| CLOUDSTACK-2182 | AuthN/AuthZ: Kyverno policies in for SRE/LCS Admin | Blocked | 26 days     | mgianluc |
+-----------------+----------------------------------------------------+---------+-------------+----------+
| CLOUDSTACK-2067 | AuthN/AuthZ: LCS RBACs for LCS Admin and SREs      | Blocked | 26 days     | mgianluc |
```


To list all jira items filed by yourself, currently worked on in active sprint and still not closed/resolved 

```
./bin/jira_utils show filed  --active 
I0429 11:01:04.121460     415 jira.go:30]  "msg"="Verifying all needed environment variables are set"  
+-----------------+-------------------------------------------------------------------------------------+-------------+-------------+----------+
|       KEY       |                                       SUMMARY                                       |   STATUS    | LAST UPDATE | ASSIGNEE |
+-----------------+-------------------------------------------------------------------------------------+-------------+-------------+----------+
| CLOUDSTACK-2349 | e2e: add dex information to each workload cluster                                   | Backlog     | 26 days     | mgianluc |
+-----------------+-------------------------------------------------------------------------------------+-------------+-------------+----------+
| CLOUDSTACK-2287 | UCS failure: node not active                                                        | Backlog     | 26 days     | vikasd   |
+-----------------+-------------------------------------------------------------------------------------+-------------+-------------+----------+
| CLOUDSTACK-2263 | list requirement for local registry in workload cluster to reach external registry  | In Progress | 26 days     | rchincha |
+-----------------+-------------------------------------------------------------------------------------+-------------+-------------+----------+
| CLOUDSTACK-2261 | design how customer will specify which registry to use for a given workload cluster | In Progress | 26 days     | rchincha |
+-----------------+-------------------------------------------------------------------------------------+-------------+-------------+----------+
| CLOUDSTACK-1795 | issues detected by kube-bench and kube-hunter                                       | Backlog     | 26 days     | srgoli   |
+-----------------+-------------------------------------------------------------------------------------+-------------+-------------+----------+
| CLOUDSTACK-1786 | kube hunter detected issues                                                         | In Progress | 26 days     | rchincha |
+-----------------+-------------------------------------------------------------------------------------+-------------+-------------+----------+
| CLOUDSTACK-1774 | kube-bench detected issues                                                          | In Progress | 26 days     | rchincha |
+-----------------+-------------------------------------------------------------------------------------+-------------+-------------+----------+
```

To list all jira issues automatically filed for e2e automatic tagging sanities and still open

```
./bin/jira_utils show e2e             
I0429 11:01:54.201978     511 jira.go:30]  "msg"="Verifying all needed environment variables are set"  
+-----------------+---------------------------+---------+-------------+----------+
|       KEY       |          SUMMARY          | STATUS  | LAST UPDATE | ASSIGNEE |
+-----------------+---------------------------+---------+-------------+----------+
| CLOUDSTACK-2330 | Test jira creating issues | Backlog | 22 days     | mgianluc |
+-----------------+---------------------------+---------+-------------+----------+
```
