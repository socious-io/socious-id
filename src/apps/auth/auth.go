package auth

import (
	"fmt"
	"socious-id/src/apps/models"
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
func GetGoogleToken(code, ref string) (string, error) {
	form := url.Values{
		"code":          {code},
		"client_id":     {Config.oauth.google.id},
		"client_secret": {Config.oauth.google.secret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {fmt.Sprintf("%s/oauth/google", ref)},
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", form)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return getGoogleUserInfo(data.AccessToken)
}

var data struct {
	Email string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName string `json:"given_name"`
}

func getGoogleUserInfo(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.Email, nil
}



async function getGoogleUserInfo(accessToken) {
	const response = await axios.get('', {
	  headers: {
		Authorization: `Bearer ${accessToken}`
	  }
	})
  
	return response.data // This contains user information
  }
  