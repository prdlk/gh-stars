# gh-stars

A [GitHub CLI](https://cli.github.com/) extension that lists the repositories a user has starred, one URL per line.

Authentication is handled by `gh` itself (via [go-gh](https://github.com/cli/go-gh)): whatever account `gh auth login` / `GH_TOKEN` / `GITHUB_TOKEN` provides is used, and `GH_HOST` selects a GitHub Enterprise host.

## Install

```sh
gh extension install prdlk/gh-stars
```

Precompiled binaries are published to [Releases](https://github.com/prdlk/gh-stars/releases) for each `v*` tag.

## Usage

```sh
gh stars [flags] [<username>]
```

With no username, lists the authenticated user's stars.

| Flag          | Default   | Values               |
| ------------- | --------- | -------------------- |
| `--sort`      | `created` | `created`, `updated` |
| `--direction` | `desc`    | `asc`, `desc`        |

Examples:

```sh
gh stars                          # your stars, newest first
gh stars octocat                  # octocat's stars
gh stars --direction asc octocat  # oldest first
gh stars | fzf | xargs gh browse  # pick one interactively
```

## Build from source

Requires Go (see `go.mod`).

```sh
go build -o gh-stars .
gh extension install .
```

`gh-stars-linux-amd64.tar.gz` in this repository contains a prebuilt linux/amd64 binary:

```sh
tar -xzf gh-stars-linux-amd64.tar.gz
./gh-stars octocat
```
