package workflow

import (
	"time"

	"github.com/brist-ai/monoid/monoidprotocol"
	"github.com/brist-ai/monoid/workflow/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type HandleRequestArgs struct {
	RequestID string
}

func (w *Workflow) HandleRequestWorkflow(
	ctx workflow.Context,
	args HandleRequestArgs,
) ([]monoidprotocol.MonoidRecord, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 2,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 5,
		},
	}

	records := []monoidprotocol.MonoidRecord{}

	ctx = workflow.WithActivityOptions(ctx, options)

	ac := activity.Activity{}

	err := workflow.ExecuteActivity(ctx, ac.ExecuteRequest, args.RequestID).Get(ctx, &records)

	return records, err
}
