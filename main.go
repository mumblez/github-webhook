package main

import (
	"fmt"
	"os"
	"strconv"

	//	"strings"

	"github.com/namsral/flag"
	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/github"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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

	// Initialise our settings
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "GHW", 0)
	fs.IntVar(&cfg.port, "port", 4567, "tcp port to listen on")
	fs.StringVar(&cfg.secret, "secret", "superDuperSecret", "github webhook secret")
	fs.StringVar(&cfg.path, "path", "payload", "url path to accept json post request, e.g. /payload")
	fs.StringVar(&cfg.repo_url, "repo_url", "", "git url to refresh e.g. git@github.com:ns/project.git")
	fs.StringVar(&cfg.repo_branch, "repo_branch", "master", "git branch to clone / checkout e.g. master (also the default)")
	fs.StringVar(&cfg.repo_dir, "repo_dir", "", "local directory path for repository")
	fs.StringVar(&cfg.repo_ssh_key, "repo_ssh_key", "", "path to ssh private key file")
	fs.StringVar(&cfg.repo_ssh_pass, "repo_ssh_pass", "", "ssh password to provided private key (optional)")

	fs.Parse(os.Args[1:])

	fmt.Println("Starting webhook....")

	hook := github.New(&github.Config{Secret: cfg.secret})
	hook.RegisterEvents(HandleMultiple, github.PushEvent) // Add as many as you want
	err := webhooks.Run(hook, ":"+strconv.Itoa(cfg.port), cfg.path)
	if err != nil {
		fmt.Println(err)
	}
}

// GitRefresh
func GitRefresh() {
	sshAuth, err := ssh.NewPublicKeysFromFile("git", cfg.repo_ssh_key, cfg.repo_ssh_pass)
	if err != nil {
		fmt.Println("Error initialising ssh key and pass")
	}

	// If directory exists just refresh else clone and init
	if _, err := os.Stat(cfg.repo_dir + "/.git"); err == nil {
		fmt.Println("Repository exists, attempting to update...")
		r, err := git.PlainOpen(cfg.repo_dir)
		if err != nil {
			fmt.Println(err)
		}

		w, err := r.Worktree()
		if err != nil {
			fmt.Println(err)
		}

		err = w.Pull(&git.PullOptions{
			RemoteName:    "origin",
			ReferenceName: plumbing.ReferenceName("refs/heads/" + cfg.repo_branch),
			SingleBranch:  true,
			Auth:          sshAuth,
		})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Initial clone of repository...")
		_, err := git.PlainClone(cfg.repo_dir, false, &git.CloneOptions{
			URL:           cfg.repo_url,
			ReferenceName: plumbing.ReferenceName("refs/heads/" + cfg.repo_branch),
			SingleBranch:  true,
			Auth:          sshAuth,
		})
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error cloning repository")
		}

	}

	fmt.Println("Git Refresh complete")
}

// HandleMultiple handles multiple GitHub events
func HandleMultiple(payload interface{}, header webhooks.Header) {

	fmt.Println("Handling Payload..")

	switch payload.(type) {

	// only push events on master!
	case github.PushPayload:
		push := payload.(github.PushPayload)
		//fmt.Printf("Repository info: %+v\n", push.Repository)

		// Refresh only if 'master' branch and repository url matches
		if push.Ref == "refs/heads/"+cfg.repo_branch && push.Repository.SSHURL == cfg.repo_url {
			GitRefresh()
		}

		// modify file so inocron can react / trigger
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
