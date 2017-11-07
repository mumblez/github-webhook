package main

import (
	"fmt"
	"os"
	"strconv"

	"strings"

	"github.com/namsral/flag"
	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/github"
)

type Config struct {
	port          int
	secret        string
	path          string
	repo_url      string
	repo_branch   string
	repo_dir      string
	repo_ssh_key  string
	repo_ssh_pass string
}

var cfg Config

func main() {

	//cfg := new(Config)

	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "GHW", 0)
	fs.IntVar(&cfg.port, "port", 8080, "tcp port to listen on")
	fs.StringVar(&cfg.secret, "secret", "superDuperSecret", "github webhook secret")
	fs.StringVar(&cfg.path, "path", "payload", "url path to accept json post request, e.g. /payload")
	fs.StringVar(&cfg.repo_url, "repo_url", "", "git url to refresh e.g. git@github.com:ns/project.git")
	fs.StringVar(&cfg.repo_branch, "repo_branch", "master", "git branch to clone / checkout e.g. master (also the default)")
	fs.StringVar(&cfg.repo_dir, "repo_dir", "", "local directory path for repository")
	fs.StringVar(&cfg.repo_ssh_key, "repo_ssh_key", "", "path to ssh private key file")
	fs.StringVar(&cfg.repo_ssh_pass, "repo_ssh_pass", "", "ssh password to provided private key (optional)")

	fs.Parse(os.Args[1:])

	// debug
	fmt.Println("Port: ", cfg.port)
	fmt.Println("Secret: ", cfg.secret)
	fmt.Println("Path: ", cfg.path)
	fmt.Println("Repository Info:")
	fmt.Println("url: ", cfg.repo_url)
	fmt.Println("branch: ", cfg.repo_branch)
	fmt.Println("dir: ", cfg.repo_dir)
	fmt.Println("ssh key: ", cfg.repo_ssh_key)
	fmt.Println("ssh pass: ", cfg.repo_ssh_pass)
	fmt.Println("Starting webhook....")

	os.Exit(0)

	hook := github.New(&github.Config{Secret: cfg.secret})
	//hook.RegisterEvents(HandleMultiple, github.ReleaseEvent, github.PullRequestEvent, github.PushEvent) // Add as many as you want
	hook.RegisterEvents(HandleMultiple, github.PushEvent) // Add as many as you want

	err := webhooks.Run(hook, ":"+strconv.Itoa(cfg.port), cfg.path)
	if err != nil {
		fmt.Println(err)
	}
}

// GitRefresh
// func GitRefresh(ssh_key, ssh_pass, url, dir string) bool {
// 	fmt.Println('placeholder')
// }

// HandleMultiple handles multiple GitHub events
func HandleMultiple(payload interface{}, header webhooks.Header) {

	fmt.Println("Handling Payload..")

	switch payload.(type) {

	// only handle commit / push events on master!
	case github.PushPayload:
		push := payload.(github.PushPayload)
		// Do whatever you want from here...
		//		fmt.Printf("Push: %+v", push)
		refs := strings.Split(push.Ref, "/")
		branch := refs[len(refs)-1]       // last
		fmt.Printf("Branch: %+v", branch) // ref: refs/heads/[branch]

		// modify file so inotify can react / trigger
		// if err := os.Chtimes("some-filename", time.Now(), time.Now()); err != nil {
		// 	log.Fatal(err)
		// }

		// case github.ReleasePayload:
		// 	release := payload.(github.ReleasePayload)
		// 	// Do whatever you want from here...
		// 	fmt.Printf("Release: %+v", release)
		//
		// case github.PullRequestPayload:
		// 	pullRequest := payload.(github.PullRequestPayload)
		// 	// Do whatever you want from here...
		// 	fmt.Printf("Pull: %+v", pullRequest)
	}
}
