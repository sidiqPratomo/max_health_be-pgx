package dto

type AccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func ConvertAccessTokenToResponse(accessToken string) AccessTokenResponse {
	return AccessTokenResponse{AccessToken: accessToken}
}