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

// pass in path and port and script / command to execute as env vars!
// see oauth2_proxy for how they mux env vars, cli options and config file

// const (
// 	path = "/payload"
// 	port = 4567
// )

func main() {

	var port int
	var secret string
	var path string

	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "GHW", 0)
	fs.IntVar(&port, "port", 8080, "tcp port to listen on")
	fs.StringVar(&secret, "secret", "superDuperSecret", "github webhook secret")
	fs.StringVar(&path, "path", "payload", "url path to accept json post request, e.g. /payload")
	fs.StringVar(&file, "file", "", "file to modify so inotify can trigger a real script!")
	fs.Parse(os.Args[1:])

	fmt.Println("Port: ", port)
	fmt.Println("Secret: ", secret)
	fmt.Println("Path: ", path)
	fmt.Println("Starting webhook....")

	hook := github.New(&github.Config{Secret: secret})
	hook.RegisterEvents(HandleMultiple, github.ReleaseEvent, github.PullRequestEvent, github.PushEvent) // Add as many as you want

	err := webhooks.Run(hook, ":"+strconv.Itoa(port), path)
	if err != nil {
		fmt.Println(err)
	}
}

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

	case github.ReleasePayload:
		release := payload.(github.ReleasePayload)
		// Do whatever you want from here...
		fmt.Printf("Release: %+v", release)

	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
		// Do whatever you want from here...
		fmt.Printf("Pull: %+v", pullRequest)
	}
}
