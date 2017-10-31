package main

import (
	"fmt"
	"os"
	"os/user"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
)

// Basic example of how to clone a repository using clone options.
func main() {
	CheckArgs("<url>", "<directory>")
	url := os.Args[1]
	directory := os.Args[2]
	sshKey := "/path/to/key"
	sshPass := "hello"

	currentUser, err := user.Current()
	CheckIfError(err)

	// Clone the given repository to the given directory
	Info("git clone %s %s", url, directory)

	Info("user = ", currentUser.HomeDir)
	// Assuming id_rsa is available at ~/.ssh/id_rsa
	//sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "hello")
	sshAuth, err := ssh.NewPublicKeysFromFile("git", sshKey, sshPass)
	CheckIfError(err)

	branch := "master"
	// cloning a single branch (master) with ssh
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
		SingleBranch:  true,
		Auth:          sshAuth,
	})

	CheckIfError(err)

	// TODO: refresh repository example - https://github.com/src-d/go-git/blob/master/_examples/pull/main.go

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}
