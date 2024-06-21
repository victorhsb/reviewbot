package workflows

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/victorhsb/review-bot/backend/service"
	"github.com/victorhsb/review-bot/backend/service/processor"
)

var (
	ErrNotAvailable = errors.New("action not available")
)

const CLIWorkflowName processor.WorkflowName = "cli"

type CLIWorkflow struct {
	review     *ReviewWorkflow
	messageSvc service.Messager
}

func (c *CLIWorkflow) Name() processor.WorkflowName {
	return CLIWorkflowName
}

func (c *CLIWorkflow) Trigger(_ context.Context, _ *uuid.UUID, _ map[string]any) (*processor.MessageProcessorState, error) {
	return nil, ErrNotAvailable
}

// ProcessMessage is the method that will be called when a message is received
func (c *CLIWorkflow) ProcessMessage(ctx context.Context, state *processor.MessageProcessorState) (*processor.MessageProcessorState, error) {
	if len(state.FootPrint) == 0 {
		return nil, fmt.Errorf("no message to process")
	}
	msg := state.FootPrint[0]
	parts := strings.Split(strings.TrimSpace(msg.Message), " ")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid commands; expected: !<command> <options...>")
	}

	switch parts[0] {
	case "!review":
		if len(parts) == 2 {
			return c.review.Trigger(ctx, state.UserID, map[string]any{"productId": parts[1]})
		}
		c.messageSvc.SaveMessage(ctx, service.Message{
			UserID:    state.UserID,
			Direction: service.DirectionSent,
			Message:   "invalid command; expected: !review <product_id>",
		})
	case "!echo":
		c.messageSvc.SaveMessage(ctx, service.Message{UserID: state.UserID, Direction: service.DirectionSent, Message: strings.Join(parts[1:], " ")})
	default:
		c.messageSvc.SaveMessage(ctx, service.Message{UserID: state.UserID, Direction: service.DirectionSent, Message: "unknown command"})
	}

	state.Status = processor.WorkflowStatusFinished
	return state, nil
}

// ClaimMessage should return true if the message can be claimed by the processor.
// it is used when the message could not be infered to which workflow it belongs.
func (c *CLIWorkflow) ClaimMessage(ctx context.Context, m *processor.MessageProcessorState) bool {
	if len(m.FootPrint) == 1 {
		return strings.TrimSpace(m.FootPrint[0].Message)[0] == '!'
	}

	return false
}

func NewCLIWorkflow(r *ReviewWorkflow, msg service.Messager) *CLIWorkflow {
	return &CLIWorkflow{review: r, messageSvc: msg}
}
