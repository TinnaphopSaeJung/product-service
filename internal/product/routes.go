package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	router.POST("/product", handler.CreateProduct)
	// router.PATCH("/product/:id", handler.PatchProduct)
}
