package processor

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/victorhsb/review-bot/backend/cache"
	"github.com/victorhsb/review-bot/backend/service"
)

// ProcessorWorkflow is the interface definition for the multiple workflows that can be applied to a message
type ProcessorWorkflow interface {
	Name() WorkflowName
	Trigger(context.Context, uuid.UUID, map[string]any) (*MessageProcessorState, error)
	// ProcessMessage is the method that will be called when a message is received
	ProcessMessage(context.Context, *MessageProcessorState) (*MessageProcessorState, error)
	// ClaimMessage should return true if the message can be claimed by the processor.
	// it is used when the message could not be infered to which workflow it belongs.
	ClaimMessage(context.Context, *MessageProcessorState) bool
}

type WorkflowName string

// MessageProcessor is an interface that defines the methods for the message processor
type MessageProcessor interface {
	ProcessMessage(context.Context, service.Message) error
}

// MessageProcessorState is a message that is queued up for processing.
// it contains the previous messages as a linked list
type MessageProcessorState struct {
	UserID    *uuid.UUID
	FootPrint []service.Message
	Workflow  WorkflowName
}

// LocalMessageProcessor is a service that processes messages inline.
// It handles everything locally
type LocalMessageProcessor struct {
	messagesCache *cache.Cache[MessageProcessorState]

	workflows map[WorkflowName]ProcessorWorkflow
}

// NewLocalMessageProcessor creates a new LocalMessageProcessor
func NewLocalMessageProcessor(ctx context.Context, workflows []ProcessorWorkflow) *LocalMessageProcessor {
	wfMap := make(map[WorkflowName]ProcessorWorkflow)
	for _, w := range workflows {
		wfMap[w.Name()] = w
	}

	return &LocalMessageProcessor{
		messagesCache: cache.NewCache[MessageProcessorState](ctx, time.Minute*30, 1000),
		workflows:     wfMap,
	}
}

func (l *LocalMessageProcessor) runWorkflow(ctx context.Context, workflow ProcessorWorkflow, state *MessageProcessorState) error {
	processedState, err := workflow.ProcessMessage(ctx, state)
	if err != nil {
		return err
	}

	l.messagesCache.Set(state.UserID.String(), *processedState)

	return nil
}

// ProcessMessage processes a message locally
// It hidrates the message with user information and footprints and then calls the appropriate workflow
func (l *LocalMessageProcessor) ProcessMessage(ctx context.Context, msg service.Message) error {
	cached, ok := l.messagesCache.Get(msg.UserID.String())
	if !ok {
		cached = &MessageProcessorState{
			FootPrint: make([]service.Message, 0, 1),
			UserID:    msg.UserID,
		}
	}
	// prepend the new message
	cached.FootPrint = append([]service.Message{msg}, cached.FootPrint...)

	return l.processMessage(ctx, cached)
}

func (l *LocalMessageProcessor) processMessage(ctx context.Context, state *MessageProcessorState) error {
	if workflow, ok := l.workflows[state.Workflow]; ok {
		return l.runWorkflow(ctx, workflow, state)
	}

	for _, w := range l.workflows {
		if w.ClaimMessage(ctx, state) {
			state.Workflow = w.Name()
			return l.runWorkflow(ctx, w, state)
		}
	}

	return nil
}
