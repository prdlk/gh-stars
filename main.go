// Command gh-stars is a GitHub CLI extension that lists the repositories a user has starred.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"regexp"

	"github.com/cli/go-gh/v2/pkg/api"
)

type repository struct {
	HTMLURL string `json:"html_url"`
}

var nextLinkRE = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "gh stars:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("gh stars", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: gh stars [flags] [<username>]")
		fmt.Fprintln(fs.Output(), "\nList starred repositories. Defaults to the authenticated user.")
		fmt.Fprintln(fs.Output(), "\nFlags:")
		fs.PrintDefaults()
	}
	sort := fs.String("sort", "created", "sort by `field`: created or updated")
	direction := fs.String("direction", "desc", "sort `direction`: asc or desc")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if fs.NArg() > 1 {
		fs.Usage()
		return errors.New("too many arguments")
	}
	if *sort != "created" && *sort != "updated" {
		return fmt.Errorf("invalid --sort %q: must be created or updated", *sort)
	}
	if *direction != "asc" && *direction != "desc" {
		return fmt.Errorf("invalid --direction %q: must be asc or desc", *direction)
	}

	client, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	path := "user/starred"
	if fs.NArg() == 1 {
		path = "users/" + url.PathEscape(fs.Arg(0)) + "/starred"
	}
	query := url.Values{"per_page": {"100"}, "sort": {*sort}, "direction": {*direction}}
	next := path + "?" + query.Encode()

	for next != "" {
		repos, link, err := fetchPage(ctx, client, next)
		if err != nil {
			return err
		}
		for _, r := range repos {
			if _, err := fmt.Fprintln(out, r.HTMLURL); err != nil {
				return err
			}
		}
		next = ""
		if m := nextLinkRE.FindStringSubmatch(link); m != nil {
			next = m[1]
		}
	}
	return nil
}

func fetchPage(ctx context.Context, client *api.RESTClient, pathOrURL string) ([]repository, string, error) {
	resp, err := client.RequestWithContext(ctx, "GET", pathOrURL, nil)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var repos []repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, "", fmt.Errorf("decoding response: %w", err)
	}
	return repos, resp.Header.Get("Link"), nil
}
