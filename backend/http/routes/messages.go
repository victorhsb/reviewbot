package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/victorhsb/review-bot/backend/service"
)

func RegisterMessageRoutes(engine *gin.Engine, svc service.Messager) {
	v1 := engine.Group("/v1")
	users := v1.Group("/users")
	{
		users.GET("", NewListUsersHandler(svc))                // list users
		users.GET("/:id", NewGetUserHandler(svc))              // get user
		users.POST("/:id/message", NewSaveMessageHandler(svc)) // save message handler
		users.GET("/:id/messages", NewGetMessagesHandler(svc)) // get messages by user handler
	}
}

func NewListUsersHandler(svc service.Messager) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		users, err := svc.ListUsers(ctx)
		if err != nil {
			log.Error().Err(err).Msg("could not get users")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get users"})
			return
		}

		log.Debug().Msg("ListUsersHandler")
		c.JSON(http.StatusOK, users)
	}
}

func NewGetUserHandler(svc service.Messager) gin.HandlerFunc {
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

		user, err := svc.GetUserByID(ctx, parsedID)
		if err != nil {
			log.Error().Err(err).Msg("could not get user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user"})
			return
		}

		log.Debug().Msg("GetUserHandler")
		c.JSON(http.StatusOK, user)
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
}

func (m NewMessagePayload) ToModel() *service.Message {
	return &service.Message{
		Message:   m.Message,
		Direction: service.DirectionSent,
	}
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

		parsedID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Warn().Err(err).Msg("could not parse user id to uuid")
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse user id"})
			return
		}

		model := message.ToModel()
		model.UserID = &parsedID

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
