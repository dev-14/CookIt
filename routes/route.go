package routes

import (
	"CookIt/controllers"
	middlewares "CookIt/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func PublicEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Generate public endpoints - [ signup] - api/v1/signup

	r.POST("/register", controllers.Register)
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/logout", authMiddleware.LogoutHandler)
}

func AuthenticatedEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r.POST("books/create", controllers.CreateRecipe)

}

func GetRouter(router chan *gin.Engine) {
	gin.ForceConsoleColor()

	r := gin.Default()

	r.Use(cors.Default())
	authMiddleware, _ := middlewares.GetAuthMiddleware()

	// Create a BASE_URL - /api/v1
	v1 := r.Group("/api/v1/")
	PublicEndpoints(v1, authMiddleware)
	AuthenticatedEndpoints(v1.Group("auth"), authMiddleware)
	router <- r
}
