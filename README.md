# Kolide Kit

Kolide Kit is a collection of Go libraries used in projects at Kolide. This repository also includes a few other features which are useful for Go developers:

- A lightweight style guide
- Links to libraries which are commonly used at Kolide
- Links to learning resources outlining some Go best practices

## Install

```
git clone git@github.com:kolide/kit.git $GOPATH/src/github.com/kolide/kit
```

## Documentation

Run `godoc -http=:6060` and then open `http://localhost:6060/pkg/github.com/kolide/kit/` in your browser. You'll see all the available packages in this repository.

## Style Guide

You will also be able to find Kolide's Go style guide at [styleguide.md](./styleguide.md). We write a lot of Go at Kolide and we like talking about how we can all write better, more consistent Go. In our style guide, we've amalgamated the results of numerous internal discussions and agreed upon best approaches.

## Git and Dependency Workflow

Our development and dependency management workflow is outlined in [workflow.md](./workflow.md). At a high-level, we use [`glide`](https://github.com/Masterminds/glide) most of the time, but we use [`dep`](https://github.com/golang/dep) for newer projects. Our git workflow has us using forks and feature branches extensively.
