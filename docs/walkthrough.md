
# Walkthrough

An example of using Repomancer to add a file to multiple repositories, commit changes, and open pull requests.

1. Start Repomancer

   ![Start Screen](images/walkthrough-start.png)

2. Click New Project. Give your project a name _(note: the project name will also be used as the branch name)_ and click Create

   ![New Project](images/walkthrough-new_project.png)

3. Click Repository > Add Multiple. Add the repositories you want to work on, one per line. Click Add.

   ![New Project](images/walkthrough-add_repositories.png)

4. Now each of these repositories has been cloned and we're ready to work.

   ![New Project](images/walkthrough-cloned.png)

5. Enter a command in the command bar and press return (or click run) to run it in each repository. _(The same command is run in the root directory of each repository)_

   ![New Project](images/walkthrough-command.png)

6. Click "Logs" to view output from all commands run in a repository

   ![New Project](images/walkthrough-view_logs.png)

7. Click "Open" to view files in the repository

   ![New Project](images/walkthrough-open_files.png)

8. Commit changes - Click Git > Commit. Enter a commit message, then click Commit _(This is the same as running `git add . && git commit -m '<message>'`. Need more control? You can also run git commands in the command box!)_

   ![New Project](images/walkthrough-commit.png)

9. Push changes: Click Git > Push

   ![New Project](images/walkthrough-push.png)

10. Open Pull Requests: Click Git > Pull Request. Enter a title and description, then click Create Pull Request.

    ![New Project](images/walkthrough-pull_request_1.png)

    ![New Project](images/walkthrough-pull_request_2.png)

11. Check Pull Request status: Click Git > Refresh Status

    ![New Project](images/walkthrough-status.png)

