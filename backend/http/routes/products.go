package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/victorhsb/review-bot/backend/service"
)

func RegisterProductRoutes(engine *gin.Engine, svc service.ProductReviewer) {
	v1 := engine.Group("/v1")
	products := v1.Group("/products")
	{
		products.GET("/:id", NewGetProductHandler(svc))
		products.POST("", NewSaveProductHandler(svc))
		products.GET("", NewListProductHandler(svc))
	}
}

func NewGetProductHandler(svc service.ProductReviewer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{"error": "missing id in request path"})
			return
		}

		parsedID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(400, fmt.Errorf("could not parse id: %w", err))
			return
		}

		product, err := svc.GetProduct(ctx, parsedID)
		if err != nil {
			c.JSON(500, gin.H{"error": "could not get product"})
			return
		}

		c.JSON(200, product)
	}
}

func NewListProductHandler(svc service.ProductReviewer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var limit, page int64
		if c.Query("limit") != "" {
			parsedLimit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not parse limit; %s", err.Error())})
				return
			}
			limit = parsedLimit
		}
		if c.Query("page") != "" {
			parsedPage, err := strconv.ParseInt(c.Query("page"), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not parse page; %s", err.Error())})
				return
			}
			page = parsedPage
		}

		products, err := svc.ListProducts(ctx, limit, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not list products; %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

type SaveProductPayload struct {
	Title string `json:"name" binding:"required"`
}

func (p SaveProductPayload) ToModel() (*service.Product, error) {
	return &service.Product{
		Title: p.Title,
	}, nil
}

func NewSaveProductHandler(svc service.ProductReviewer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var payload SaveProductPayload
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(400, gin.H{"error": "invalid payload"})
			return
		}

		product, err := payload.ToModel()
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid payload"})
			return
		}

		if id := c.Param("id"); id != "" {
			parsedID, err := uuid.Parse(id)
			if err != nil {
				c.JSON(400, gin.H{"error": "could not parse id"})
				return
			}
			product.ID = &parsedID
		}

		created, err := svc.SaveProduct(ctx, *product)
		if err != nil {
			c.JSON(500, gin.H{"error": "could not save product"})
			return
		}

		c.Header("X-Created-ID", created.String())
		c.Status(http.StatusAccepted)
	}
}
