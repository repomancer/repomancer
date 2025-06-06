# Repomancer
A tool for bulk updating repositories

![Repomancer main screen while jobs are running](docs/images/mainscreen.png)


## Features

- Manage pull requests across multiple repositories
- Execute commands on some or all repositories
- Track open/merged pull requests and statistics

See the [walkthrough](docs/walkthrough.md) for example usage.

## Installation

**Development version instructions**

Assuming you have `go` installed, install Repomancer by:

1. Install the GitHub CLI (`gh`): https://cli.github.com/
   ```shell
   # macOS: 
   brew install gh
   ```

2. Install Repomancer:
    ```shell
    GOPROXY=direct go install github.com/repomancer/repomancer@v0.0.1
    ```

3. Configure GitHub credentials:

   You **MUST** login to each GitHub host, even if you normally use an SSH key [^1]. **Connect with HTTPS.**

    ```shell
   # Repeat for each GitHub host you want to use:
   gh auth login
    ```
4. Run Repomancer:
   ```shell
   repomancer
    ```

[^1]: `gh` needs a personal access token to connect to the GitHub API, even if it uses the SSH key for cloning/pushing.
