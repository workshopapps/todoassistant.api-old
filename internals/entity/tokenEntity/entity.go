package tokenEntity

type TokenRes struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
