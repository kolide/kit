## Kolide Kit

Kolide Kit is a collection of Go libraries used in projects at Kolide. This repository also includes a few other features which are useful for Go developers:

- A lightweight style guide
- Links to libraries which are commonly used at Kolide
- Links to learning resources outlining some Go best practices

# Install

```
git clone git@github.com:kolide/kit.git $GOPATH/src/github.com/kolide/kit
```

# Documentation

Run `godoc -http=:6060` and then open `http://localhost:6060/pkg/github.com/kolide/kit/` in your browser. You'll see all the available packages in this repository.

# Dependency management and Git workflow.

## Git

At Kolide, we use GitHub for source control. 

* Projects live in the [GOPATH](https://github.com/golang/go/wiki/GOPATH) at the original `$GOPATH/src/github.com/kolide/$repo` path. 
* `github.com/kolide/$repo` is used as the git origin, with your fork being added as a remote. The workflow for a new feature branch becomes:

    ```
    # First you would clone a repo
    git clone git@github.com:kolide/kit.git $GOPATH/src/github.com/kolide/kit
    
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

Historically we've used [`glide`](https://github.com/Masterminds/glide#glide-vendor-package-management-for-golang) to manage Go dependencies, but we've started to adopt [`dep`](https://github.com/golang/dep) for newwer projects.

Using dep requires that you edit the `Gopkg.toml` file with constraints and overrides for the project. See the [oficial docs](https://github.com/golang/dep/blob/master/docs/Gopkg.toml.md) for an up to date guide on the Gopkg file format.
You can run `dep ensure -examples` to see a list of commonly used `dep` commands.

# Go Style Guide

It helps keep development and code review by having general consensus on a set of best practices which we all follow. Our internal style guide is a set of code standards that we try to adhere to whenever possible. Some high-level guidance is:

* Defer to the Go [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments#go-code-review-comments). We largely follow the same conventions in our code.
* Follow these [best practices](https://peter.bourgon.org/go-best-practices-2016/) from Peter Bourgon.
* Avoid package level variables and `init`. Avoiding global state leads to code which is more readable, testable and maintainable. See [this blog](https://peter.bourgon.org/blog/2017/06/09/theory-of-modern-go.html).
* Write tests using the [testify library](https://godoc.org/github.com/stretchr/testify/assert).
* Preferably write your tests as a [table test](https://github.com/golang/go/wiki/TableDrivenTests).
* Use [subtests](https://blog.golang.org/subtests) to run your table driven tests. Subtests provide a way to better handle test failures and and [parallelize](https://rakyll.org/parallelize-test-tables/) tests. Consider the following example test:
    ```go
    func TestAuthenticatedHost(t *testing.T) {
        // set up test dependencies
    	ctx := context.Background()
    	goodNodeKey, err := svc.EnrollAgent(ctx, "foobarbaz", "host123")
    
        // use require if the test cannot continue if the assertion fails
    	require.Nil(t, err)
    	require.NotEmpty(t, goodNodeKey)
    
        // create a []struct for your test cases
    	var authenticatedHostTests = []struct {
    		nodeKey   string
    		shouldErr bool
    	}{
    		{
    			nodeKey:   "invalid",
    			shouldErr: true,
    		},
    		{
    			nodeKey:   "",
    			shouldErr: true,
    		},
    		{
    			nodeKey:   goodNodeKey,
    			shouldErr: false,
    		},
    	}
    
        // use subtests to run through your test cases.
    	for _, tt := range authenticatedHostTests {
    		t.Run("", func(t *testing.T) {
    			var r = struct{ NodeKey string }{NodeKey: tt.nodeKey}
    			_, err = endpoint(context.Background(), r)
    			if tt.shouldErr {
    				assert.IsType(t, osqueryError{}, err)
    			} else {
    				assert.Nil(t, err)
    			}
    		})
    	}
    
    }
    ```

* Use functional options for optional function parameters. [blog](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis), [video](https://www.youtube.com/watch?v=24lFtGHWxAQ)  

Example:  
Let's say you have a `Client` struct, which will implement an API client and has a default timeout of 5 seconds. One way to create the Client would be to write a function like:
```go
NewClient(baseurl *url.URL, timeout time.Duration, debugMode bool) *Client
```

But every time you'll want to add a new configuration parameter, you'll have to make a breaking change to NewClient. A cleaner, more extensible solution is to write it with the following pattern:
```go
// Declare a function type for modifying the client
type Option(*Client)

// WithTimeout sets the timeout on the Client.
func WithTimeout(d time.Duration) Option {
    return func(c *Client) {
        c.timeout = d
    }
}

func Debug() Option {
    return func(c *Client) {
        c.debug = true
    }
}
```

Now you can write the client which will accept a variadic number of option arguments.

```go
NewClient(baseurl *url.URL, opts ...Option) *Client {
    // create a client with some default values.
    client := &Client{
        timeout: 5 * time.Minute,
    }

    // loop through the provided options and override any of the defaults.
    for _, opt := range opts {
        opt(&client)
    }

    return &client
}
```

* Propagate a context through your API. 
The `context` package provides a standard way for managing cancellations and request scoped values in a Go program. When writing server and client code, it is recommended to add `context.Context` as the first argument to your methods.
For example, if you have a function like:

```go
func User(id uint) (*User, error)
```

you should instead write it as:

```go
func User(ctx context.Context, id uint) (*User, error)
```


See the following resources on `context.Context`:
    * https://blog.golang.org/context
    * https://peter.bourgon.org/blog/2016/07/11/context.html
    * [justforfunc video on context use](https://www.youtube.com/watch?v=LSzR0VEraWw&index=1&list=PL64wiCrrxh4Jisi7OcCJIUpguV_f5jGnZ)
    * [GolangUK talk](https://www.youtube.com/watch?v=r4Mlm6qEWRs)
