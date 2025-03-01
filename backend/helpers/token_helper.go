package helpers

import (
	"backend/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenHelper struct {
	SecretKey string
}

func NewTokenHelper(secretKey string) *TokenHelper {
	return &TokenHelper{
		SecretKey: secretKey,
	}
}

func (t *TokenHelper) GenerateToken(email, username string, createdAt int64, expireAfter time.Duration) (string, error) {
	claims := &models.MemberSignedDetails{
		Email:     email,
		Username:  username,
		CreatedAt: createdAt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireAfter).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(t.SecretKey))
	return token, err
}

func (t *TokenHelper) ValidateToken(tokenString string) (*models.MemberSignedDetails, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.MemberSignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.MemberSignedDetails)
	if !ok {
		return nil, err
	}

	return claims, nil
}
