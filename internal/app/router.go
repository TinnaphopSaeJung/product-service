package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	productModule "product-service/internal/product"
	postgresRepo "product-service/internal/repository/postgres"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "product-service/docs"
)

func NewRouter(dbPool *pgxpool.Pool) *gin.Engine {
	router := gin.Default()

	router.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	productRepository := postgresRepo.NewProductRepository(dbPool)
	productUsecase := productModule.NewUsecase(productRepository)
	productHandler := productModule.NewHandler(productUsecase)

	productModule.RegisterRoutes(router, productHandler)

	return router
}
