package routes

import (
	"inventory-management/internal/handler"

	"github.com/gin-gonic/gin"
)

// Layer routes berubah saat komposisi endpoint atau grouping API berubah,
// bukan saat business logic service berubah.
func Register(router *gin.Engine, itemHandler *handler.ItemHandler) {
	apiV1 := router.Group("/api/v1")
	{
		items := apiV1.Group("/items")
		{
			items.POST("", itemHandler.Create)
			items.GET("", itemHandler.List)
			items.GET("/:id", itemHandler.GetByID)
			items.PUT("/:id", itemHandler.Update)
			items.DELETE("/:id", itemHandler.Delete)
		}
	}
}
