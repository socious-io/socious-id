package views

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	rootGroup(r)
	authGroup(r)
	usersGroup(r)
	organizationsGroup(r)
	mediaGroup(r)
	verificationsGroup(r)
	impactPointsGroup(r)
	kybVerificationGroup(r)
	referralAchievementsGroup(r)
}
