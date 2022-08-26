package routes

import (
	"AmzonElGalaba/Controller"
	"github.com/gin-gonic/gin"
)

var User Controller.User

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("signup", User.SignUp)
	incomingRoutes.POST("signIn", User.Login)
}
