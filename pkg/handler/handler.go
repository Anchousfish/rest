package handler

import (
	"github.com/Anchousfish/rest/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		s: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
	}
	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getList)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			listItems := lists.Group(":id/items")
			{
				listItems.POST("/", h.createItem)
				listItems.GET("/", h.getAllItems)
			}
		}

		items := lists.Group("/items")
		{
			items.GET("/:id", h.getItem)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}

	}
	return router
}
