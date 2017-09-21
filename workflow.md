# Workflows

## Git

At Kolide, we use GitHub for source control.

* Projects live in the [GOPATH](https://github.com/golang/go/wiki/GOPATH) at the original `$GOPATH/src/github.com/kolide/$repo` path.
* `github.com/kolide/$repo` is used as the git origin, with your fork being added as a remote. The workflow for a new feature branch becomes:

    ```
    # First you would clone a repo
    git clone git@github.com:kolide/kit.git $GOPATH/src/github.com/kolide/kit
    cd $GOPATH/src/github.com/kolide/kit

    # Add your fork as a git remote
    $username = "groob" # this should be whatever your GitHub username is
    git remote add $username git@github.com:$username/kit.git

    # Pull from origin
    git pull origin master --rebase

    # Create your feature
    git checkout -b feature-branch

    # Push to your fork
    git push -u $username feature-branch

    # Open a pull request on GitHub.

    # Continue to push to your fork as you iterate
    git add .
    git commit
    git push $username feature-branch
    ```

* Prefer small, self contained feature branches.
* Request code reviews from at least one person on your team.
* You can commit to your branch however many times you like, but we have found that using the "Squash and Merge" feature on GitHub works well for us. Once a Pull Request goes through code review and receives approval, the original author should squash and merge the pull request, adding a final commit message which will show up in the master branch's commit history.

## Go dependencies

Historically we've used [`glide`](https://github.com/Masterminds/glide#glide-vendor-package-management-for-golang) to manage Go dependencies, but we've started to adopt [`dep`](https://github.com/golang/dep) for newer projects.

Using dep requires that you edit the `Gopkg.toml` file with constraints and overrides for the project. See the [oficial docs](https://github.com/golang/dep/blob/master/docs/Gopkg.toml.md) for an up to date guide on the Gopkg file format.
You can run `dep ensure -examples` to see a list of commonly used `dep` commands.
