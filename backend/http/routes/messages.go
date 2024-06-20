package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/victorhsb/review-bot/backend/service"
)

func RegisterMessageRoutes(engine *gin.Engine, svc service.Messager) {
	v1 := engine.Group("/v1")
	msgs := v1.Group("/messages")
	{
		msgs.GET("/:id", NewGetMessagesHandler(svc)) // get messages by participant handler
		msgs.POST("", NewSaveMessageHandler(svc))    // save message handler
	}
}

func NewGetMessagesHandler(svc service.Messager) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id := c.Param("id")
		if id == "" {
			log.Warn().Msg("missing id in request path")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing id in request path"})
			return
		}

		parsedID, err := uuid.Parse(id)
		if err != nil {
			log.Warn().Err(err).Msg("missing id in request path")
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request id"})
			return
		}

		messages, err := svc.ListMessagesByUserID(ctx, parsedID)
		if err != nil {
			log.Error().Err(err).Msg("could not get messages")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get messages"})
			return
		}

		log.Debug().Any("return", messages).Msg("GetMessagesHandler")
		c.JSON(http.StatusOK, messages)
	}
}

type NewMessagePayload struct {
	Message string `json:"message" binding:"required"`
	UserID  string `json:"userId" binding:"required,uuid"`
}

func (m NewMessagePayload) ToModel() (*service.Message, error) {
	mod := service.Message{
		Message:   m.Message,
		Direction: service.DirectionSent,
	}

	parsedID, err := uuid.Parse(m.UserID)
	if err != nil {
		return nil, fmt.Errorf("could not parse user id; %w", err)
	}
	mod.UserID = parsedID

	return &mod, nil
}

func NewSaveMessageHandler(svc service.Messager) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var message NewMessagePayload
		if err := c.BindJSON(&message); err != nil {
			log.Warn().Err(err).Msg("could not bind body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse request body"})
			return
		}

		model, err := message.ToModel()
		if err != nil {
			log.Warn().Err(err).Msg("could not parse message")
		}

		err = svc.SaveMessage(
			ctx,
			*model,
		)
		if err != nil {
			log.Error().Err(err).Msg("could not save messages")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save message"})
			return
		}

		log.Debug().Msg("SaveMessageHandler")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
