package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
)

var ticket string

func init() {
	flag.StringVar(&ticket, "ticket", "", "JIRA ticket to fetch")
}

func main() {
	flag.Parse()
	if ticket == "" {
		log.Fatal("Please provide a JIRA ticket using the -ticket flag")
	}

	token := os.Getenv("JIRA_TOKEN")
	if token == "" {
		log.Fatal("JIRA_TOKEN environment variable is not set")
	}

	base := "https://<org>.atlassian.net"
	tp := jira.BasicAuthTransport{
		Username: "<email>",
		Password: token,
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}

	issue, _, err := jiraClient.Issue.Get(ticket, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Issue does not exist") {
			log.Fatalf("Issue %s does not exist", ticket)
		}
		log.Fatal(err)
	}

	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)

	// MESOS-3325: Running mesos-slave@0.23 in a container causes slave to be lost after a restart
	// Type: Bug
	// Priority: Critical
}
