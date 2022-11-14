package main

import (
	"log"
	"os"

	"github.com/brist-ai/monoid/cmd"
	"github.com/brist-ai/monoid/workflow"
	"github.com/brist-ai/monoid/workflow/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	conf := cmd.GetBaseConfig(false, cmd.Models)

	// Create the client object just once per process
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL"),
	})

	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}

	defer c.Close()

	w := worker.New(c, workflow.DockerRunnerQueue, worker.Options{})
	a := activity.Activity{
		Conf: &conf,
	}

	mwf := workflow.Workflow{
		Conf: &conf,
	}

	w.RegisterActivity(a.ValidateDataSiloDef)
	w.RegisterActivity(a.DetectDataSources)
	w.RegisterActivity(a.FindOrCreateJob)
	w.RegisterActivity(a.UpdateJobStatus)

	w.RegisterWorkflow(mwf.ValidateDSWorkflow)
	w.RegisterWorkflow(mwf.DetectDSWorkflow)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
