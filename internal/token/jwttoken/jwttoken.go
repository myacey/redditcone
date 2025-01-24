package jwttoken

import (
	"encoding/json"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/myacey/redditclone/internal/models"
	"github.com/myacey/redditclone/internal/token"
)

type JWTToken struct {
	secretKey []byte
}

func NewJWTToken(secretKey []byte) token.TokenMaker {
	jwtToken := &JWTToken{secretKey}
	return jwtToken
}

func (t *JWTToken) CreateToken(usr *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user": usr,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tkn.SignedString(t.secretKey)
	if err != nil {
		return "", token.ErrCantCreateToken
	}

	return tokenString, nil
}

func (t *JWTToken) ExtractUserID(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, token.ErrInvalidTokenMethod
		}
		return t.secretKey, nil
	})
	if err != nil {
		return "", err
	}

	userData, err := json.Marshal(claims["user"])
	if err != nil {
		return "", token.ErrNoClaims
	}
	var usr models.User
	err = json.Unmarshal(userData, &usr)
	if err != nil {
		return "", token.ErrNoClaims
	}

	fmt.Printf("TOKEN: %v\n", claims)

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return "", token.ErrCantExtractExpTime
	}
	exp := int64(expFloat)

	expirationTime := time.Unix(exp, 0)
	if expirationTime.Before(time.Now()) {
		return "", token.ErrTokenExpired
	}

	return usr.ID, nil
}
