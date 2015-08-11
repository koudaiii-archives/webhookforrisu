package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidEventFormat = errors.New("Unable to parse event string. Invalid Format.")

type Event struct {
	Owner      string // The username of the owner of the repository
	Repo       string // The name of the repository
	Branch     string // The branch the event took place on
	Commit     string // The head commit hash attached to the event
	Type       string // Can be either "pull_request" or "push"
	BaseOwner  string // For Pull Requests, contains the base owner
	BaseRepo   string // For Pull Requests, contains the base repo
	BaseBranch string // For Pull Requests, contains the base branch
}

// Create a new event from a string, the string format being the same as the one produced by event.String()
func NewEvent(e string) (*Event, error) {
	// Trim whitespace
	e = strings.Trim(e, "\n\t ")

	// Split into lines
	parts := strings.Split(e, "\n")

	// Sanity checking
	if len(parts) != 5 || len(parts) != 8 {
		return nil, ErrInvalidEventFormat
	}
	for _, item := range parts {
		if len(item) < 8 {
			return nil, ErrInvalidEventFormat
		}
	}

	// Fill in values for the event
	event := Event{}
	event.Type = parts[0][8:]
	event.Owner = parts[1][8:]
	event.Repo = parts[2][8:]
	event.Branch = parts[3][8:]
	event.Commit = parts[4][8:]

	// Fill in extra values if it's a pull_request
	if event.Type == "pull_request" {
		if len(parts) != 8 {
			return nil, ErrInvalidEventFormat
		}
		event.BaseOwner = parts[5][8:]
		event.BaseRepo = parts[6][8:]
		event.BaseBranch = parts[7][8:]
	}

	return &event, nil
}

func (e *Event) String() (output string) {
	output += "type:   " + e.Type + "\n"
	output += "owner:  " + e.Owner + "\n"
	output += "repo:   " + e.Repo + "\n"
	output += "branch: " + e.Branch + "\n"
	output += "commit: " + e.Commit + "\n"

	if e.Type == "pull_request" {
		output += "bowner: " + e.BaseOwner + "\n"
		output += "brepo:  " + e.BaseRepo + "\n"
		output += "bbranch:" + e.BaseBranch + "\n"
	}

	return
}

type Server struct {
	Port   int        // Port to listen on. Defaults to 80
	Path   string     // Path to receive on. Defaults to "/build"
	Secret string     // Option secret key for authenticating via HMAC
	Events chan Event // Channel of events. Read from this channel to get push events as they happen.
}

// NewServer is default the Port is set to 80 and the Path is set to `/build`
func NewServer() *Server {
	return &Server{
		Port:   8080,
		Path:   "/build",
		Events: make(chan Event, 10), // buffered to 10 items
	}
}

func main() {

	server := NewServer()
	server.Port = 8080
	server.Secret = "supersecretcode"

	for {
		select {
		case event := <-server.Events:
			fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
		default:
			time.Sleep(100)
		}
	}

}
