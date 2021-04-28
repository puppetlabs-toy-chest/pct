# Contributing

Hi! Thanks for your interest in contributing to the Puppet Development Kit!

Community contributions are essential for keeping Puppet great. We simply can't access the huge number of platforms and myriad configurations for running Puppet. We want to keep it as easy as possible to contribute changes that get things working in your environment. There are a few guidelines that we need contributors to follow so that we can have a chance of keeping on top of things.

We accept pull requests for bug fixes and features where we've discussed the approach in an issue and given the go-ahead for a community member to work on it. We'd also love to hear about ideas for new features as Github Discussion.

Please do:

- Open an issue if things aren't working as expected.
- Open a Github Discussion to propose a significant change.
- Open a pull request to fix a bug.
- Open a pull request to fix documentation about a command.
- Open a pull request for any issue labelled help wanted or good first issue.

## Building the project

Prerequisites:

- Go 1.16+

To build the PDK, run the following command:

```bash
> # on nix
> go build -o pdk
```

```powershell
# on windows
> go build -o pdk.exe
```

To run the new binary:

```bash
> ./pdk
```

To test the PDK, run the following command.

```bash
> go test ./...
```

## Running the project

```bash
> # on nix
> go run cmd/pdk/main.go
```

## Building Cross Platform binaries

Prerequisites:

- Go 1.16+
- goreleaser 0.16+

To build the PDK for more than your current platform, use [`GoReleaser`](https://goreleaser.com/quick-start/#dry-run):

```bash
> goreleaser --snapshot --skip-publish --rm-dist
```

This will ouput a set of binaries in the `dist` folder.

## Submitting a Pull Request

1. Create a new branch `git checkout -b my-branch-name`
1. Make you change, add tests, and ensure tests pass
1. Make sure your commit messages are in the proper format. If the commit addresses an issue filed in the project, start the first line of the commit with the prefix `GH-` and the issue number in parentheses `(GH-111)`. After leave a detailed explanation of your change, so a person in the future can understand what your work does.
1. Submit a pull request (i.e. using the Github commandline: `gh pr create --web`)

## Making Trivial Changes

For changes of a trivial nature, it is not always necessary to create a new ticket. In this case, it is appropriate to start the first line of a commit with one of (docs), (maint), or (packaging) instead of a ticket number.

## Additional Resources

* [Puppet community guidelines](https://puppet.com/community/community-guidelines)
* [Contributor License Agreement](http://cla.puppet.com/)
* [General GitHub documentation](https://help.github.com/)
* [GitHub pull request documentation](https://help.github.com/articles/creating-a-pull-request/)
