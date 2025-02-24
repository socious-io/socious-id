package views

type ConfirmForm struct {
	Confirmed bool `json:"confirmed" form:"confirmed"`
}

type AuthSessionForm struct {
	ClientSecret string `json:"client_secret" form:"client_secret" validate:"required"`
	ClientID     string `json:"client_id" form:"client_id" validate:"required"`
	RedirectURL  string `json:"redirect_url" form:"redirect_url" validate:"required"`
}

type GetTokenForm struct {
	ClientSecret string `json:"client_secret" form:"client_secret" validate:"required"`
	ClientID     string `json:"client_id" form:"client_id" validate:"required"`
	Code         string `json:"code" form:"code" validate:"required"`
}
