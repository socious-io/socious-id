package auth

import (
	"fmt"
	"math"
	"math/rand/v2"
	"regexp"
	"socious-id/src/apps/models"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterForm struct {
	FirstName *string    `json:"first_name" form:"first_name"`
	LastName  *string    `json:"last_name" form:"last_name"`
	Username  *string    `json:"username" form:"username"`
	Email     string     `json:"email" form:"email" validate:"required,email"`
	Password  *string    `json:"password" form:"password"`
	AvatarID  *uuid.UUID `json:"avatar_id" form:"avatar_id"`
}

type LoginForm struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

type OTPForm struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type SetPasswordForm struct {
	Password string `json:"password" form:"password" validate:"required"`
}

type ChangePasswordForm struct {
	CurrentPassword string `json:"current_password" form:"current_password" validate:"required"`
	SetPasswordForm
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func Signin(id, email string) (map[string]any, error) {
	accessToken, err := GenerateToken(id, email, false)
	if err != nil {
		return nil, err
	}
	refreshToken, err := GenerateToken(id, email, true)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	}, nil
}

func FetchUserByJWT(c *gin.Context) (*models.User, error) {
	tokenStr := c.GetHeader("Authorization")
	splited := strings.Split(tokenStr, " ")
	if len(splited) > 1 {
		tokenStr = splited[1]
	} else {
		tokenStr = splited[0]
	}
	if tokenStr == "" {
		return nil, fmt.Errorf("Authorization header missing")
	}

	claims, err := VerifyToken(tokenStr)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("Invalid token signature")
		}
		return nil, err
	}
	return models.GetUser(uuid.MustParse(claims.ID))
}

func FetchUserBySession(c *gin.Context) (*models.User, error) {
	session := sessions.Default(c)
	id := session.Get("user_id")
	if id == nil {
		return nil, fmt.Errorf("not authorized")
	}
	userID, err := uuid.Parse(id.(string))
	if err != nil {
		return nil, fmt.Errorf("not authorized")
	}
	return models.GetUser(userID)
}

func GenerateUsername(email string) string {
	var username string = email
	var re *regexp.Regexp

	re = regexp.MustCompile("@.*$")
	username = re.ReplaceAllString(username, "")

	re = regexp.MustCompile("[^a-z0-9._-]")
	username = re.ReplaceAllString(username, "-")

	re = regexp.MustCompile("[._-]{2,}")
	username = re.ReplaceAllString(username, "-")

	username = strings.ToLower(username)
	username = username[0:int(math.Min(float64(len(username)), 20))]

	username = username + strconv.Itoa(int(1000+rand.Float64()*9000))

	return username
}
