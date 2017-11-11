# github webhook

Super simple single purpose app to handle incoming webhook request from github
and to auth with secret!

After authorising the request, we check the branch and repository url in the payload and if it matches we refresh the repository on disk.
We only handle a single repository!


# Prerequisites

- Generate an ssh key pair (with passphrase) and put the public key as a new Deploy key for the repo we want to update (on github)
- Generate a webhook for the repository with a secret, path and set the application type to json
- Ensure the user the service is running as can write to the repository directory (to clone / pull)

# Configuration

## Order of precedence:

1. Command line options
2. Environment variables
3. Configuration file
4. Default values

| Type   | CLI Flag                             | Environment (all prefixed with `GHW_`)  | File          | Default Value    | Notes              |
| ------ | :----------------------------------- |:--------------------------------------- |:------------- |:-----------------|:-------------------|
| int    | -port 1010                           | PORT=1010                               | port 1010     | 4567             | TCP port to listen |
| string | -path /payload                       | PATH="/payload"                         | path /payload | /payload | URI path, e.g. https://domain.com/payload |
| string | -secret gitHubWebHookSecret12345     | SECRET="gitHubWebHookSecret12345"       | secret gitHubWebHookSecret12345     |       | webhook secret set on github |
| string | -repo_ssh_key /path/to/private/key   | REPO_SSH_KEY="/path/to/private/key"     | repo_ssh_key /path/to/private/key   |       | path to ssh private key (deploy key) |
| string | -repo_ssh_pass sshPassphrase123      | REPO_SSH_PASS="sshPassphrase123"        | repo_ssh_pass sshPassphrase123      |       | passphrase to ssh key |
| string | -repo_url git@github.com:ns/repo.git | REPO_URL="git@github.com:ns/repo.git"   | repo_url git@github.com:ns/repo.git |       | git url (ssh, not http) |
| string | -repo_branch master                  | REPO_BRANCH="master"                    | repo_branch master                  | master| branch to clone / update |
| string | -repo_dir /path/to/clone/to          | REPO_DIR="/path/to/clone/to"            | repo_dir /path/to/clone/to          |       | local directory to clone repository to |
| int    | -health_check_port 9091              | HEALTH_CHECK_PORT=9091                  | health_check_port 9091              | 9091  | port to handle health check |
| string | -health_check_path /ping             | HEALTH_CHECK_PATH="/ping"               | health_check_path /ping             | /ping | path to handle health check |

If we set a configuration file, pass the path to `-config` on the cli

