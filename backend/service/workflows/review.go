package workflows

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/victorhsb/review-bot/backend/cache"
	"github.com/victorhsb/review-bot/backend/service"
	"github.com/victorhsb/review-bot/backend/service/processor"
)

const (
	ReviewWorkflowName processor.WorkflowName = "review"

	reviewTriggerMessage string = "Hi! I just saw that your %s has arrived! That is great. In a scale of 1 to 5 (1 being bad and 5 being excelent) how do you rate your product?"
	reviewCommentMessage string = "Thank you! how would you describe your experience with the product?"
	reviewDoneMessage    string = "Thank you for your review! It was very helpful!"
)

const (
	stageRating int = iota
	stageComment
	stageDone
	stageExit
)

type ReviewWorkflowState struct {
	UserID      *uuid.UUID
	ProductID   *uuid.UUID
	ReviewID    *uuid.UUID
	ReviewStage int
}

func NewReviewWorkflow(ctx context.Context, productReviewer service.ProductReviewer, messager service.Messager) *ReviewWorkflow {
	return &ReviewWorkflow{
		cache:           cache.NewCache[ReviewWorkflowState](ctx, 30*time.Minute, 1000),
		productReviewer: productReviewer,
		messager:        messager,
	}
}

// ReviewWorkflow is the service implementation that handles the review process
type ReviewWorkflow struct {
	cache *cache.Cache[ReviewWorkflowState]

	productReviewer service.ProductReviewer
	messager        service.Messager
}

// Trigger is used to trigger the workflow by the system or another workflow.
// It accepts the UUID of the target user and a map of metadata that the workflow might need.
func (r *ReviewWorkflow) Trigger(ctx context.Context, userID *uuid.UUID, metadata map[string]any) (*processor.MessageProcessorState, error) {
	productId, ok := metadata["productId"].(string)
	if !ok {
		return nil, errors.New("product id not found or invalid")
	}

	parsedProductId, err := uuid.Parse(productId)
	if err != nil {
		return nil, fmt.Errorf("could not parse product id; %w", err)
	}

	product, err := r.productReviewer.GetProduct(ctx, parsedProductId)
	if err != nil {
		return nil, fmt.Errorf("could not find product with id %s; %w", productId, err)
	}

	message := service.Message{
		UserID:    userID,
		Direction: service.DirectionSent,
		// TODO: replace with a message template
		Message: fmt.Sprintf(reviewTriggerMessage, product.Title),
	}

	if err := r.messager.SaveMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("could not save message; %w", err)
	}

	r.cache.Set(userID.String(), ReviewWorkflowState{
		ProductID: &parsedProductId,
		UserID:    userID,
	})

	return &processor.MessageProcessorState{
		UserID: userID,
		FootPrint: []service.Message{
			message,
		},
		Workflow: r.Name(),
		Status:   processor.WorkflowStatusNext,
	}, nil
}

// Name returns the name of the workflow
func (r ReviewWorkflow) Name() processor.WorkflowName {
	return ReviewWorkflowName
}

// ProcessMessage
func (r *ReviewWorkflow) ProcessMessage(ctx context.Context, state *processor.MessageProcessorState) (*processor.MessageProcessorState, error) {
	wstate, ok := r.cache.Get(state.UserID.String())
	if !ok {
		return nil, errors.New("could not find workflow state")
	}

	err := r.processMessage(ctx, state, wstate)
	if err != nil {
		var userErr *UserError
		if errors.As(err, &userErr) {
			if err := r.messager.SaveMessage(ctx, service.Message{
				UserID:  state.UserID,
				Message: userErr.Message,
			}); err != nil {
				log.Error().Err(err).Msg("could not send error message to user")
				return nil, err
			}
		}
		log.Error().Err(err).Msg("could not process message")
		return nil, err
	}

	if wstate.ReviewStage == stageExit {
		log.Info().Msg("user exited workflow earlier")
		wstate.ReviewStage = stageDone
	}
	if wstate.ReviewStage == stageDone {
		r.cache.Delete(state.UserID.String())
		state.Status = processor.WorkflowStatusFinished
		return state, nil
	}
	r.cache.Set(state.UserID.String(), *wstate)

	state.Status = processor.WorkflowStatusNext
	return state, nil
}

func (r *ReviewWorkflow) processMessage(ctx context.Context, state *processor.MessageProcessorState, wstate *ReviewWorkflowState) error {
	lastMsg := state.FootPrint[0]
	if lastMsg.Direction == service.DirectionSent {
		// this means somehow this was triggered by a sent message
		return errors.New("processed message was sent by the system, this should not happen")
	}

	if lastMsg.Message == "exit" {
		wstate.ReviewStage = stageExit
		return nil
	}

	switch wstate.ReviewStage {
	case stageRating:
		review, err := r.ParseReview(lastMsg)
		if err != nil {
			return err
		}
		review.ProductID = wstate.ProductID
		review.UserID = wstate.UserID

		created, err := r.productReviewer.SaveProductReview(ctx, *review)
		if err != nil {
			return fmt.Errorf("could not save product review; %w", err)
		}
		wstate.ReviewID = created.ID
		wstate.ReviewStage = stageComment

		err = r.messager.SaveMessage(ctx, service.Message{
			UserID:  state.UserID,
			Message: reviewCommentMessage,
		})
		if err != nil {
			log.Error().Err(err).Msg("could not send comment message to user")
		}
	case stageComment:
		review, err := r.productReviewer.GetProductReview(ctx, *wstate.ReviewID)
		if err != nil {
			return fmt.Errorf("could not find review with id %s; %w", wstate.ReviewID, err)
		}

		review.Review = lastMsg.Message
		if err := r.productReviewer.UpdateProductReview(ctx, *review); err != nil {
			return fmt.Errorf("could not update product review; %w", err)
		}

		err = r.messager.SaveMessage(ctx, service.Message{
			UserID:  state.UserID,
			Message: reviewDoneMessage,
		})
		if err != nil {
			log.Error().Err(err).Msg("could not send comment message to user")
		}
		wstate.ReviewStage = stageDone
	}

	return nil
}

func (r *ReviewWorkflow) ParseReview(msg service.Message) (*service.ProductReview, error) {
	rating, err := strconv.Atoi(strings.TrimSpace(msg.Message))
	if err != nil {
		log.Error().Err(err).Msg("could not parse rating")
		return nil, &UserError{Message: "I couldn't understand your rating. Please, rate the product from 1 to 5 using numerals only", Err: err}
	}

	if rating < 1 || rating > 5 {
		return nil, &UserError{Message: "The rating must be between 1 and 5. nothing more than that. Please, try again."}
	}

	return &service.ProductReview{
		Rating: rating,
	}, nil
}

// The review workflow should never claim a message. It should only be triggered by the system.
func (r *ReviewWorkflow) ClaimMessage(ctx context.Context, m *processor.MessageProcessorState) bool {
	return false
}
