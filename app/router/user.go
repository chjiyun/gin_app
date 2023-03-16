package router

import (
	"gin_app/app/controller"
	"github.com/gin-gonic/gin"
)

func (r *Router) User(g *gin.RouterGroup) {
	rg := g.Group("/user")
	{
		rg.POST("/login", controller.Login)
		rg.POST("/register", controller.Register)
		rg.POST("/logout", controller.Logout)
		rg.POST("/resetPassword", controller.ResetPassword)
	}
}
