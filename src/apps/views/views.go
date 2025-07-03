package views

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	rootGroup(r)
	authGroup(r)
	usersGroup(r)
	organizationsGroup(r)
	mediaGroup(r)
	credentialsGroup(r)
	impactPointsGroup(r)
	kybVerificationGroup(r)
	paymentsGroup(r)
	referralsGroup(r)
}
