package router

import (
	"gin_app/app/service/user"
	"github.com/gin-gonic/gin"
)

func (r *Router) User(g *gin.RouterGroup) {
	rg := g.Group("/user")
	{
		rg.POST("/login", user.Login)
		rg.POST("/register", user.Register)
		rg.POST("/logout", user.Logout)
		rg.POST("/resetPassword", user.ResetPassword)
	}
}
