package util

import (
	"encoding/json"
	"time"

	"github.com/sidiqPratomo/max-health-backend/config"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserId        int64  `json:"user_id"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	TokenDuration int    `json:"token_duration"`
}

type TokenAuthentication interface {
	CreateAndSign(customClaims JwtCustomClaims, secretKey string) (*string, error)
	ParseAndVerify(signed string, secretKey string) (*JwtCustomClaims, error)
}

type JwtAuthentication struct {
	Config config.Config
	Method *jwt.SigningMethodHMAC
}

func (ja JwtAuthentication) CreateAndSign(customClaims JwtCustomClaims, secretKey string) (*string, error) {
	customClaimsJsonBytes, err := json.Marshal(customClaims)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(ja.Method, jwt.MapClaims{
		"iss":  ja.Config.Issuer,
		"iat":  time.Now(),
		"exp":  time.Now().Add(time.Duration(customClaims.TokenDuration) * time.Minute).Unix(),
		"data": string(customClaimsJsonBytes),
	})

	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &signed, nil
}

func (ja JwtAuthentication) ParseAndVerify(signed string, secretKey string) (*JwtCustomClaims, error) {
	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{ja.Method.Name}),
		jwt.WithIssuer(ja.Config.Issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	customClaims := JwtCustomClaims{}
	data := claims["data"].(string)
	if err := json.Unmarshal([]byte(data), &customClaims); err != nil {
		return nil, err
	}

	return &customClaims, nil
}
