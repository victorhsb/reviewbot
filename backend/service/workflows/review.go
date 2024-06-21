package workflows

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/victorhsb/review-bot/backend/cache"
	"github.com/victorhsb/review-bot/backend/service"
	"github.com/victorhsb/review-bot/backend/service/processor"
)

const (
	ReviewWorkflowName processor.WorkflowName = "review"
)

type ReviewWorkflowState struct {
	UserID    uuid.UUID `json:"userId"`
	ProductID uuid.UUID `json:"productId"`
	ReviewID  uuid.UUID `json:"reviewId"`
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

// Name returns the name of the workflow
func (r ReviewWorkflow) Name() processor.WorkflowName {
	return ReviewWorkflowName
}

// ProcessMessage
func (r *ReviewWorkflow) ProcessMessage(_ context.Context, _ *processor.MessageProcessorState) (*processor.MessageProcessorState, error) {
	panic("not implemented") // TODO: Implement
}

// The review workflow should never claim a message. It should only be triggered by the system.
func (r *ReviewWorkflow) ClaimMessage(ctx context.Context, m *processor.MessageProcessorState) bool {
	return false
}
