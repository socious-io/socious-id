package views

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	authGroup(r)
	usersGroup(r)
	organizationsGroup(r)
	mediaGroup(r)
	verificationsGroup(r)
}
