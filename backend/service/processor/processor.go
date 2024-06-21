package processor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/victorhsb/review-bot/backend/cache"
	"github.com/victorhsb/review-bot/backend/service"
)

// ProcessorWorkflow is the interface definition for the multiple workflows that can be applied to a message
// A workflow can be triggered by a message or by the system/other workflow. For this there is the Trigger method and
// the ClaimMessage method.
//
// The Processmessage is the main method that will be called when a message is received.
// It should return the new state of the message processor. It can also return a Processing status to indicate that the
// message needs to be reprocessed with the new state (basically routing the state to another workflow).
type ProcessorWorkflow interface {
	Name() WorkflowName
	// ProcessMessage is the method that will be called when a message is received
	ProcessMessage(ctx context.Context, state *MessageProcessorState) (*MessageProcessorState, error)
	// ClaimMessage should return true if the message can be claimed by the processor.
	// it is used when the message could not be infered to which workflow it belongs.
	ClaimMessage(ctx context.Context, state *MessageProcessorState) bool
	// Trigger is used to trigger the workflow by the system or another workflow.
	// It accepts the UUID of the target user and a map of metadata that the workflow might need.
	Trigger(ctx context.Context, userId *uuid.UUID, metadata map[string]any) (*MessageProcessorState, error)
}

type WorkflowName string

type WorkflowStatus int

const (
	// WorkflowStatusPending means that the message still needs to be processed
	// and the state should not be saved yet.
	WorkflowStatusPending WorkflowStatus = 1 + iota
	// WorkflowStatusNext is the status of a message that has been processed and the state can
	// be saved for further iterations.
	WorkflowStatusNext
	// WorkflowStatusFinished is the status of a message that has been processed and the state can
	// be discarded
	WorkflowStatusFinished
	// WorkflowStatusCorrupted is the status of a message that cannot be processed and should be discarded
	// because it's state is possibly corrupted
	WorkflowStatusCorrupted
)

// MessageProcessorState is a message that is queued up for processing.
// it contains the previous messages as a linked list
type MessageProcessorState struct {
	UserID    *uuid.UUID
	FootPrint []service.Message
	Workflow  WorkflowName
	Status    WorkflowStatus
}

// LocalMessageProcessor is a service that processes messages inline.
// It handles everything locally
type LocalMessageProcessor struct {
	messagesCache *cache.Cache[MessageProcessorState]
	messager      service.Messager

	workflows map[WorkflowName]ProcessorWorkflow
}

var _ service.MessageProcessor = (*LocalMessageProcessor)(nil)

// NewLocalMessageProcessor creates a new LocalMessageProcessor
func NewLocalMessageProcessor(ctx context.Context, workflows ...ProcessorWorkflow) *LocalMessageProcessor {
	wfMap := make(map[WorkflowName]ProcessorWorkflow)
	for _, w := range workflows {
		wfMap[w.Name()] = w
	}

	return &LocalMessageProcessor{
		messagesCache: cache.NewCache[MessageProcessorState](ctx, time.Minute*30, 1000),
		workflows:     wfMap,
	}
}

func (l *LocalMessageProcessor) Use(w ProcessorWorkflow) {
	l.workflows[w.Name()] = w
}

// ProcessMessage implements the MessageProcessor interface
// It gets the user's message history and finds the appropriate workflow to process the message
func (l *LocalMessageProcessor) ProcessMessage(ctx context.Context, msg service.Message) error {
	log.Debug().Any("message", msg).Msg("processing message")
	cached, ok := l.messagesCache.Get(msg.UserID.String())
	if !ok {
		cached = &MessageProcessorState{
			FootPrint: make([]service.Message, 0, 1),
			UserID:    msg.UserID,
		}
	}
	// prepend the new message
	cached.FootPrint = append([]service.Message{msg}, cached.FootPrint...)
	cached.Status = WorkflowStatusPending

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

	log.Debug().Any("message", state.FootPrint[0]).Msg("no workflow found for message")

	return nil
}

func (l *LocalMessageProcessor) runWorkflow(ctx context.Context, workflow ProcessorWorkflow, state *MessageProcessorState) error {
	processedState, err := workflow.ProcessMessage(ctx, state)
	if err != nil {
		return err
	}

	switch processedState.Status {
	case WorkflowStatusPending:
		return l.processMessage(ctx, state)
	case WorkflowStatusNext:
		l.messagesCache.Set(state.UserID.String(), *processedState)
	case WorkflowStatusFinished:
		l.messagesCache.Delete(state.UserID.String())
	case WorkflowStatusCorrupted:
		l.messager.SaveMessage(ctx, service.Message{UserID: state.UserID, Direction: service.DirectionSent, Message: "Sorry, I could not process your message. Try again from the start"})
		l.messagesCache.Delete(state.UserID.String())
	default:
		log.Error().Any("state", processedState).Msg("unknown state")
	}

	return nil
}
