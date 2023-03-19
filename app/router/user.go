package router

import (
	"gin_app/app/controller/userController"
	"github.com/gin-gonic/gin"
)

func (r *Router) User(g *gin.RouterGroup) {
	rg := g.Group("/user")
	{
		rg.POST("/login", userController.Login)
		rg.POST("/register", userController.Register)
		rg.POST("/logout", userController.Logout)
		rg.POST("/resetPassword", userController.ResetPassword)
		rg.GET("/page", userController.GetPageUsers)
		rg.GET("/current", userController.GetCurrentUser)
	}
}
